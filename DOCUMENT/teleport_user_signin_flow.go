// onSSH executes 'tsh ssh' command
func onSSH(cf *CLIConf) {
	tc, err := makeClient(cf)
	if err != nil {
		utils.FatalError(err)
	}
	tc.Stdin = os.Stdin
	if err = tc.SSH(cf.RemoteCommand, cf.LocalExec); err != nil {
		// exit with the same exit status as the failed command:
		if tc.ExitStatus != 0 {
			os.Exit(tc.ExitStatus)
		} else {
			utils.FatalError(err)
		}
	}
}
	// makeClient takes the command-line configuration and constructs & returns
	// a fully configured TeleportClient object
	func makeClient(cf *CLIConf) (tc *client.TeleportClient, err error) {
		// apply defults
		if cf.NodePort == 0 {
			cf.NodePort = defaults.SSHServerListenPort
		}
		if cf.MinsToLive == 0 {
			cf.MinsToLive = int32(defaults.CertDuration / time.Minute)
		}

		// split login & host
		hostLogin := cf.NodeLogin
		var labels map[string]string
		if cf.UserHost != "" {
			parts := strings.Split(cf.UserHost, "@")
			if len(parts) > 1 {
				hostLogin = parts[0]
				cf.UserHost = parts[1]
			}
			// see if remote host is specified as a set of labels
			if strings.Contains(cf.UserHost, "=") {
				labels, err = client.ParseLabelSpec(cf.UserHost)
				if err != nil {
					return nil, err
				}
			}
		}
		fPorts, err := client.ParsePortForwardSpec(cf.LocalForwardPorts)
		if err != nil {
			return nil, err
		}

		// prep client config:
		c := &client.Config{
			Stdout:             os.Stdout,
			Stderr:             os.Stderr,
			Stdin:              os.Stdin,
			Username:           cf.Username,
			ProxyHost:          cf.Proxy,
			Host:               cf.UserHost,
			HostPort:           int(cf.NodePort),
			HostLogin:          hostLogin,
			Labels:             labels,
			KeyTTL:             time.Minute * time.Duration(cf.MinsToLive),
			InsecureSkipVerify: cf.InsecureSkipVerify,
			LocalForwardPorts:  fPorts,
			ConnectorID:        cf.ExternalAuth,
			SiteName:           cf.SiteName,
			Interactive:        cf.Interactive,
		}
		return client.NewClient(c)
	}

	// SSH connects to a node and, if 'command' is specified, executes the command on it,
	// otherwise runs interactive shell
	//
	// Returns nil if successful, or (possibly) *exec.ExitError
	func (tc *TeleportClient) SSH(command []string, runLocally bool) error {
		// connect to proxy first:
		if !tc.Config.ProxySpecified() {
			return trace.BadParameter("proxy server is not specified")
		}
		proxyClient, err := tc.ConnectToProxy()
		if err != nil {
			return trace.Wrap(err)
		}
		defer proxyClient.Close()
		siteInfo, err := proxyClient.getSite()
		if err != nil {
			return trace.Wrap(err)
		}
		// which nodes are we executing this commands on?
		nodeAddrs, err := tc.getTargetNodes(proxyClient)
		if err != nil {
			return trace.Wrap(err)
		}
		if len(nodeAddrs) == 0 {
			return trace.BadParameter("no target host specified")
		}
		// more than one node for an interactive shell?
		// that can't be!
		if len(nodeAddrs) != 1 {
			fmt.Printf(
				"\x1b[1mWARNING\x1b[0m: multiple nodes match the label selector. Picking %v (first)\n",
				nodeAddrs[0])
		}
		nodeClient, err := proxyClient.ConnectToNode(nodeAddrs[0]+"@"+siteInfo.Name, tc.Config.HostLogin)
		if err != nil {
			return trace.Wrap(err)
		}
		// proxy local ports (forward incoming connections to remote host ports)
		tc.startPortForwarding(nodeClient)

		// local execution?
		if runLocally {
			if len(tc.Config.LocalForwardPorts) == 0 {
				fmt.Println("Executing command locally without connecting to any servers. This makes no sense.")
			}
			return runLocalCommand(command)
		}
		// execute command(s) or a shell on remote node(s)
		if len(command) > 0 {
			return tc.runCommand(siteInfo.Name, nodeAddrs, proxyClient, command)
		}
		return tc.runShell(nodeClient, nil)
	}

		// ConnectToProxy dials the proxy server and returns ProxyClient if successful
		func (tc *TeleportClient) ConnectToProxy() (*ProxyClient, error) {
			proxyAddr := tc.Config.ProxyHostPort(false)
			sshConfig := &ssh.ClientConfig{
				User:            tc.getProxyLogin(),
				HostKeyCallback: tc.HostKeyCallback,
			}

			log.Infof("[CLIENT] connecting to proxy %v with host login '%v'", proxyAddr, sshConfig.User)

			// try to authenticate using every non interactive auth method we have:
			for _, m := range tc.authMethods() {
				sshConfig.Auth = []ssh.AuthMethod{m}
				proxyClient, err := ssh.Dial("tcp", proxyAddr, sshConfig)
				if err != nil {
					if utils.IsHandshakeFailedError(err) {
						continue
					}
					return nil, trace.Wrap(err)
				}
				log.Infof("[CLIENT] successfully authenticated with %v", proxyAddr)
				return &ProxyClient{
					Client:          proxyClient,
					proxyAddress:    proxyAddr,
					hostKeyCallback: sshConfig.HostKeyCallback,
					authMethods:     tc.authMethods(),
					hostLogin:       tc.Config.HostLogin,
					siteName:        tc.Config.SiteName,
				}, nil
			}
			// we have exhausted all auth existing auth methods and local login
			// is disabled in configuration
			if tc.Config.SkipLocalAuth {
				return nil, trace.BadParameter("failed to authenticate with proxy %v", proxyAddr)
			}
			// if we get here, it means we failed to authenticate using stored keys
			// and we need to ask for the login information
			err := tc.Login()
			if err != nil {
				// we need to communicate directly to user here,
				// otherwise user will see endless loop with no explanation
				if trace.IsTrustError(err) {
					fmt.Printf("Refusing to connect to untrusted proxy %v without --insecure flag\n", proxyAddr)
				}
				return nil, trace.Wrap(err)
			}
			log.Debugf("Received a new set of keys from %v", proxyAddr)
			// After successfull login we have local agent updated with latest
			// and greatest auth information, try it now
			sshConfig.Auth = []ssh.AuthMethod{authMethodFromAgent(tc.localAgent)}
			proxyClient, err := ssh.Dial("tcp", proxyAddr, sshConfig)
			if err != nil {
				return nil, trace.Wrap(err)
			}
			log.Debugf("Successfully authenticated with %v", proxyAddr)
			return &ProxyClient{
				Client:          proxyClient,
				proxyAddress:    proxyAddr,
				hostKeyCallback: sshConfig.HostKeyCallback,
				authMethods:     tc.authMethods(),
				hostLogin:       tc.Config.HostLogin,
				siteName:        tc.Config.SiteName,
			}, nil
		}
		
			// Login logs user in using proxy's local 2FA auth access
			// or used OIDC external authentication, it later
			// saves the generated credentials into local keystore for future use
			func (tc *TeleportClient) Login() error {
				// generate a new keypair. the public key will be signed via proxy if our password+HOTP  are legit
				key, err := tc.MakeKey()
				if err != nil {
					return trace.Wrap(err)
				}

				var response *web.SSHLoginResponse
				if tc.ConnectorID == "" {
					response, err = tc.directLogin(key.Pub)
					if err != nil {
						return trace.Wrap(err)
					}
				} else {
					response, err = tc.oidcLogin(tc.ConnectorID, key.Pub)
					if err != nil {
						return trace.Wrap(err)
					}
					// in this case identity is returned by the proxy
					tc.Username = response.Username
				}
				key.Cert = response.Cert
				// save the key:
				if err = tc.localAgent.AddKey(tc.ProxyHost, tc.Config.Username, key); err != nil {
					return trace.Wrap(err)
				}
				// save the list of CAs we trust to the cache file
				err = tc.localAgent.AddHostSignersToCache(response.HostSigners)
				if err != nil {
					return trace.Wrap(err)
				}

				// get site info:
				proxy, err := tc.ConnectToProxy()
				if err != nil {
					return trace.Wrap(err)
				}
				site, err := proxy.getSite()
				if err != nil {
					return trace.Wrap(err)
				}
				tc.SiteName = site.Name
				return nil
			}
				// directLogin asks for a password + HOTP token, makes a request to CA via proxy
				func (tc *TeleportClient) directLogin(pub []byte) (*web.SSHLoginResponse, error) {
					httpsProxyHostPort := tc.Config.ProxyHostPort(true)
					certPool := loopbackPool(httpsProxyHostPort)

					// ping the HTTPs endpoint first:
					if err := web.Ping(httpsProxyHostPort, tc.InsecureSkipVerify, certPool); err != nil {
						return nil, trace.Wrap(err)
					}

					password, hotpToken, err := tc.AskPasswordAndHOTP()
					if err != nil {
						return nil, trace.Wrap(err)
					}

					// ask the CA (via proxy) to sign our public key:
					response, err := web.SSHAgentLogin(httpsProxyHostPort,
						tc.Config.Username,
						password,
						hotpToken,
						pub,
						tc.KeyTTL,
						tc.InsecureSkipVerify,
						certPool)

					return response, trace.Wrap(err)
				}
					// AskPasswordAndHOTP prompts the user to enter the password + HTOP 2nd factor
					func (tc *TeleportClient) AskPasswordAndHOTP() (pwd string, token string, err error) {
						fmt.Printf("Enter password for Teleport user %v:\n", tc.Config.Username)
						pwd, err = passwordFromConsole()
						if err != nil {
							fmt.Println(err)
							return "", "", trace.Wrap(err)
						}

						fmt.Printf("Enter your HOTP token:\n")
						token, err = lineFromConsole()
						if err != nil {
							fmt.Println(err)
							return "", "", trace.Wrap(err)
						}
						return pwd, token, nil
					}

