/*
Copyright 2016 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pcclient

import (
    "bufio"
    "crypto/x509"
    "encoding/pem"
    "fmt"
    "io"
    "io/ioutil"
    "net"
    "os"
    "os/exec"
    "os/signal"
    "os/user"
    "path/filepath"
    "strconv"
    "strings"
    "syscall"

    "github.com/gravitational/teleport/lib/auth/native"
    "github.com/gravitational/teleport/lib/client"
    "github.com/gravitational/teleport/lib/services"
    "github.com/gravitational/teleport/lib/session"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/web"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/agent"
    "golang.org/x/crypto/ssh/terminal"
    "github.com/stkim1/pcteleport/pcdefaults"
)

// TeleportClient is a wrapper around SSH client with teleport specific
// workflow built in
type PocketClient struct {
    client.Config
    localAgent *LocalKeyAgent

    // OnShellCreated gets called when the shell is created. It's
    // safe to keep it nil
    OnShellCreated ShellCreatedCallback
}

// ShellCreatedCallback can be supplied for every teleport client. It will
// be called right after the remote shell is created, but the session
// hasn't begun yet.
//
// It allows clients to cancel SSH action
type ShellCreatedCallback func(shell io.ReadWriteCloser) (exit bool, err error)

func (tc *PocketClient) authMethods() []ssh.AuthMethod {
    return tc.Config.AuthMethods
}

// NewClient creates a TeleportClient object and fully configures it
func NewPocketClient(c *client.Config) (tc *PocketClient, err error) {
    // validate configuration
    if c.Username == "" {
        c.Username = Username()
        log.Infof("no teleport login given. defaulting to %s", c.Username)
    }
    if c.ProxyHost == "" {
        return nil, trace.Errorf("No proxy address specified, missed --proxy flag?")
    }
    if c.HostLogin == "" {
        c.HostLogin = Username()
        log.Infof("no host login given. defaulting to %s", c.HostLogin)
    }
    if c.KeyTTL == 0 {
        c.KeyTTL = pcdefaults.CertDuration
    } else if c.KeyTTL > pcdefaults.MaxCertDuration || c.KeyTTL < pcdefaults.MinCertDuration {
        return nil, trace.Errorf("invalid requested cert TTL")
    }

    tc = &PocketClient{Config: *c}

    // initialize the local agent (auth agent which uses local SSH keys signed by the CA):
    tc.localAgent, err = NewLocalAgent(c.KeysDir, c.Username)
    if err != nil {
        return nil, trace.Wrap(err)
    }

    if tc.Stdout == nil {
        tc.Stdout = os.Stdout
    }
    if tc.Stderr == nil {
        tc.Stderr = os.Stderr
    }
    if tc.Stdin == nil {
        tc.Stdin = os.Stdin
    }
    if tc.HostKeyCallback == nil {
        tc.HostKeyCallback = tc.localAgent.CheckHostSignature
    }
/*
    // sometimes we need to use external auth without using local auth
    // methods, e.g. in automation daemons
    if c.SkipLocalAuth {
        if len(c.AuthMethods) == 0 {
            return nil, trace.BadParameter("SkipLocalAuth is true but no AuthMethods provided")
        }
        return tc, nil
    }

    // we're not going to use ssh agent
    // first, see if we can authenticate with credentials stored in
    // a local SSH agent:
    if sshAgent := connectToSSHAgent(); sshAgent != nil {
        tc.Config.AuthMethods = append(tc.Config.AuthMethods, authMethodFromAgent(sshAgent))
    }
*/
    // then, we'll auth with the local agent keys:
    tc.Config.AuthMethods = append(tc.Config.AuthMethods, authMethodFromAgent(tc.localAgent))

    return tc, nil
}

func (tc *PocketClient) LocalAgent() *LocalKeyAgent {
    return tc.localAgent
}

// getTargetNodes returns a list of node addresses this SSH command needs to
// operate on.
func (tc *PocketClient) getTargetNodes(proxy *ProxyClient) ([]string, error) {
    var (
        err    error
        nodes  []services.Server
        retval = make([]string, 0)
    )
    if tc.Labels != nil && len(tc.Labels) > 0 {
        nodes, err = proxy.FindServersByLabels(tc.Labels)
        if err != nil {
            return nil, trace.Wrap(err)
        }
        for i := 0; i < len(nodes); i++ {
            retval = append(retval, nodes[i].Addr)
        }
    }
    if len(nodes) == 0 {
        retval = append(retval, net.JoinHostPort(tc.Host, strconv.Itoa(tc.HostPort)))
    }
    return retval, nil
}

// SSH connects to a node and, if 'command' is specified, executes the command on it,
// otherwise runs interactive shell
//
// Returns nil if successful, or (possibly) *exec.ExitError
func (tc *PocketClient) SSH(command []string, runLocally bool) error {
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

func (tc *PocketClient) startPortForwarding(nodeClient *NodeClient) error {
    if len(tc.Config.LocalForwardPorts) > 0 {
        for _, fp := range tc.Config.LocalForwardPorts {
            socket, err := net.Listen("tcp", net.JoinHostPort(fp.SrcIP, strconv.Itoa(fp.SrcPort)))
            if err != nil {
                return trace.Wrap(err)
            }
            go nodeClient.listenAndForward(socket, net.JoinHostPort(fp.DestHost, strconv.Itoa(fp.DestPort)))
        }
    }
    return nil
}

// Join connects to the existing/active SSH session
func (tc *PocketClient) Join(sessionID session.ID, input io.Reader) (err error) {
    tc.Stdin = input
    if sessionID.Check() != nil {
        return trace.Errorf("Invalid session ID format: %s", string(sessionID))
    }
    var notFoundErrorMessage = fmt.Sprintf("session '%s' not found or it has ended", sessionID)

    // connect to proxy:
    if !tc.Config.ProxySpecified() {
        return trace.BadParameter("proxy server is not specified")
    }
    proxyClient, err := tc.ConnectToProxy()
    if err != nil {
        return trace.Wrap(err)
    }
    defer proxyClient.Close()
    site, err := proxyClient.ConnectToSite()
    if err != nil {
        return trace.Wrap(err)
    }

    // find the session ID on the site:
    sessions, err := site.GetSessions()
    if err != nil {
        return trace.Wrap(err)
    }
    var session *session.Session
    for _, s := range sessions {
        if s.ID == sessionID {
            session = &s
            break
        }
    }
    if session == nil {
        return trace.NotFound(notFoundErrorMessage)
    }

    // pick the 1st party of the session and use his server ID to connect to
    if len(session.Parties) == 0 {
        return trace.NotFound(notFoundErrorMessage)
    }
    serverID := session.Parties[0].ServerID

    // find a server address by its ID
    nodes, err := site.GetNodes()
    if err != nil {
        return trace.Wrap(err)
    }
    var node *services.Server
    for _, n := range nodes {
        if n.ID == serverID {
            node = &n
            break
        }
    }
    if node == nil {
        return trace.NotFound(notFoundErrorMessage)
    }
    // connect to server:
    fullNodeAddr := node.Addr
    if tc.SiteName != "" {
        fullNodeAddr = fmt.Sprintf("%s@%s", node.Addr, tc.SiteName)
    }
    nc, err := proxyClient.ConnectToNode(fullNodeAddr, tc.Config.HostLogin)
    if err != nil {
        return trace.Wrap(err)
    }
    defer nc.Close()

    // start forwarding ports, if configured:
    tc.startPortForwarding(nc)

    // running shell with a given session means "join" it:
    return tc.runShell(nc, session)
}

// SCP securely copies file(s) from one SSH server to another
func (tc *PocketClient) SCP(args []string, port int, recursive bool) (err error) {
    if len(args) < 2 {
        return trace.Errorf("Need at least two arguments for scp")
    }
    first := args[0]
    last := args[len(args)-1]

    // local copy?
    if !isRemoteDest(first) && !isRemoteDest(last) {
        return trace.BadParameter("making local copies is not supported")
    }

    if !tc.Config.ProxySpecified() {
        return trace.BadParameter("proxy server is not specified")
    }
    log.Infof("Connecting to proxy...")
    proxyClient, err := tc.ConnectToProxy()
    if err != nil {
        return trace.Wrap(err)
    }
    defer proxyClient.Close()

    // gets called to convert SSH error code to tc.ExitStatus
    onError := func(err error) error {
        exitError, _ := trace.Unwrap(err).(*ssh.ExitError)
        if exitError != nil {
            tc.ExitStatus = exitError.ExitStatus()
        }
        return err
    }
    // upload:
    if isRemoteDest(last) {
        login, host, dest := parseSCPDestination(last)
        if login != "" {
            tc.HostLogin = login
        }
        addr := net.JoinHostPort(host, strconv.Itoa(port))

        client, err := proxyClient.ConnectToNode(addr, tc.HostLogin)
        if err != nil {
            return trace.Wrap(err)
        }
        // copy everything except the last arg (that's destination)
        for _, src := range args[:len(args)-1] {
            err = client.Upload(src, dest, tc.Stderr)
            if err != nil {
                return onError(err)
            }
            fmt.Printf("Uploaded %s\n", src)
        }
        // download:
    } else {
        login, host, src := parseSCPDestination(first)
        addr := net.JoinHostPort(host, strconv.Itoa(port))
        if login != "" {
            tc.HostLogin = login
        }
        client, err := proxyClient.ConnectToNode(addr, tc.HostLogin)
        if err != nil {
            return trace.Wrap(err)
        }
        // copy everything except the last arg (that's destination)
        for _, dest := range args[1:] {
            err = client.Download(src, dest, recursive, tc.Stderr)
            if err != nil {
                return onError(err)
            }
            fmt.Printf("Downloaded %s\n", src)
        }
    }
    return nil
}

// parseSCPDestination takes a string representing a remote resource for SCP
// to download/upload, like "user@host:/path/to/resource.txt" and returns
// 3 components of it
func parseSCPDestination(s string) (login, host, dest string) {
    i := strings.IndexRune(s, '@')
    if i > 0 && i < len(s) {
        login = s[:i]
        s = s[i+1:]
    }
    parts := strings.Split(s, ":")
    return login, parts[0], strings.Join(parts[1:], ":")
}

func isRemoteDest(name string) bool {
    return strings.IndexRune(name, ':') >= 0
}

// ListNodes returns a list of nodes connected to a proxy
func (tc *PocketClient) ListNodes() ([]services.Server, error) {
    var err error
    // userhost is specified? that must be labels
    if tc.Host != "" {
        tc.Labels, err = ParseLabelSpec(tc.Host)
        if err != nil {
            return nil, trace.Wrap(err)
        }
    }

    // connect to the proxy and ask it to return a full list of servers
    proxyClient, err := tc.ConnectToProxy()
    if err != nil {
        return nil, trace.Wrap(err)
    }

    defer proxyClient.Close()
    return proxyClient.FindServersByLabels(tc.Labels)
}

// runCommand executes a given bash command on a bunch of remote nodes
func (tc *PocketClient) runCommand(siteName string, nodeAddresses []string, proxyClient *ProxyClient, command []string) error {
    resultsC := make(chan error, len(nodeAddresses))
    for _, address := range nodeAddresses {
        go func(address string) {
            var (
                err         error
                nodeSession *NodeSession
            )
            defer func() {
                resultsC <- err
            }()
            var nodeClient *NodeClient
            nodeClient, err = proxyClient.ConnectToNode(address+"@"+siteName, tc.Config.HostLogin)
            if err != nil {
                fmt.Fprintln(tc.Stderr, err)
                return
            }
            defer nodeClient.Close()

            // run the command on one node:
            if len(nodeAddresses) > 1 {
                fmt.Printf("Running command on %v:\n", address)
            }
            nodeSession, err = newSession(nodeClient, nil, tc.Config.Env, tc.Stdin, tc.Stdout, tc.Stderr)
            if err != nil {
                log.Error(err)
                return
            }
            if err = nodeSession.runCommand(command, tc.OnShellCreated, tc.Config.Interactive); err != nil {
                exitErr, ok := err.(*ssh.ExitError)
                if ok {
                    tc.ExitStatus = exitErr.ExitStatus()
                }
            }
        }(address)
    }
    var lastError error
    for range nodeAddresses {
        if err := <-resultsC; err != nil {
            lastError = err
        }
    }
    return trace.Wrap(lastError)
}

// runShell starts an interactive SSH session/shell.
// sessionID : when empty, creates a new shell. otherwise it tries to join the existing session.
func (tc *PocketClient) runShell(nodeClient *NodeClient, sessToJoin *session.Session) error {
    nodeSession, err := newSession(nodeClient, sessToJoin, tc.Env, tc.Stdin, tc.Stdout, tc.Stderr)
    if err != nil {
        return trace.Wrap(err)
    }
    if err = nodeSession.runShell(tc.OnShellCreated); err != nil {
        return trace.Wrap(err)
    }
    if nodeSession.ExitMsg == "" {
        fmt.Printf("Connection to %s closed from the remote side\n", tc.NodeHostPort())
    } else {
        fmt.Println(nodeSession.ExitMsg)
    }
    return nil
}

// getProxyLogin determines which SSH login to use when connecting to proxy.
func (tc *PocketClient) getProxyLogin() string {
    // we'll fall back to using the target host login
    proxyLogin := tc.Config.HostLogin

    // see if we already have a signed key in the cache, we'll use that instead
    if !tc.Config.SkipLocalAuth {
        keys, err := tc.GetKeys()
        if err == nil && len(keys) > 0 {
            principals := keys[0].Certificate.ValidPrincipals
            if len(principals) > 0 {
                proxyLogin = principals[0]
            }
        }
    }
    return proxyLogin
}

// GetKeys returns a list of stored local keys/certs for this Teleport
// user
func (tc *PocketClient) GetKeys() ([]agent.AddedKey, error) {
    return tc.LocalAgent().GetKeys(tc.Username)
}

// ConnectToProxy dials the proxy server and returns ProxyClient if successful
func (tc *PocketClient) ConnectToProxy() (*ProxyClient, error) {
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

// Logout locates a certificate stored for a given proxy and deletes it
func (tc *PocketClient) Logout() error {
    return trace.Wrap(tc.localAgent.DeleteKey(tc.ProxyHost, tc.Config.Username))
}

// Login logs user in using proxy's local 2FA auth access
// or used OIDC external authentication, it later
// saves the generated credentials into local keystore for future use
func (tc *PocketClient) Login() error {
    // generate a new keypair. the public key will be signed via proxy if our password+HOTP  are legit
    key, err := tc.MakeKey()
    if err != nil {
        return trace.Wrap(err)
    }

    var response *web.SSHLoginResponse
    response, err = tc.directLogin(key.Pub)
    if err != nil {
        return trace.Wrap(err)
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

// Adds a new CA as trusted CA for this client
func (tc *PocketClient) AddTrustedCA(ca *services.CertAuthority) error {
    return tc.LocalAgent().AddHostSignersToCache([]services.CertAuthority{*ca})
}

// MakeKey generates a new unsigned key. It's useless by itself until a
// trusted CA signs it
func (tc *PocketClient) MakeKey() (key *Key, err error) {
    key = &Key{}
    keygen := native.New()
    defer keygen.Close()
    key.Priv, key.Pub, err = keygen.GenerateKeyPair("")
    if err != nil {
        return nil, trace.Wrap(err)
    }
    return key, nil
}

func (tc *PocketClient) AddKey(host string, key *Key) error {
    return tc.localAgent.AddKey(host, tc.Username, key)
}

// directLogin asks for a password + HOTP token, makes a request to CA via proxy
func (tc *PocketClient) directLogin(pub []byte) (*web.SSHLoginResponse, error) {
    httpsProxyHostPort := tc.Config.ProxyHostPort(true)
    certPool := loopbackPool(httpsProxyHostPort)

    // ping the HTTPs endpoint first:
    if err := web.Ping(httpsProxyHostPort, tc.InsecureSkipVerify, certPool); err != nil {
        return nil, trace.Wrap(err)
    }

    // TODO we'll get the password and htopToken rightway
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

// oidcLogin opens browser window and uses OIDC redirect cycle with browser
func (tc *PocketClient) oidcLogin(connectorID string, pub []byte) (*web.SSHLoginResponse, error) {
    log.Infof("oidcLogin start")
    // ask the CA (via proxy) to sign our public key:
    webProxyAddr := tc.Config.ProxyHostPort(true)
    response, err := web.SSHAgentOIDCLogin(webProxyAddr,
        connectorID, pub, tc.KeyTTL, tc.InsecureSkipVerify, loopbackPool(webProxyAddr))
    return response, trace.Wrap(err)
}

// loopbackPool reads trusted CAs if it finds it in a predefined location
// and will work only if target proxy address is loopback
func loopbackPool(proxyAddr string) *x509.CertPool {
    if !utils.IsLoopback(proxyAddr) {
        log.Debugf("not using loopback pool for remote proxy addr: %v", proxyAddr)
        return nil
    }
    log.Debugf("attempting to use loopback pool for local proxy addr: %v", proxyAddr)
    certPool := x509.NewCertPool()

    certPath := filepath.Join(pcdefaults.DataDir, pcdefaults.SelfSignedCertPath)
    pemByte, err := ioutil.ReadFile(certPath)
    if err != nil {
        log.Debugf("could not open any path in: %v", certPath)
        return nil
    }

    for {
        var block *pem.Block
        block, pemByte = pem.Decode(pemByte)
        if block == nil {
            break
        }
        cert, err := x509.ParseCertificate(block.Bytes)
        if err != nil {
            log.Debugf("could not parse cert in: %v, err: %v", certPath, err)
            return nil
        }
        certPool.AddCert(cert)
    }
    log.Debugf("using local pool for loopback proxy: %v, err: %v", certPath, err)
    return certPool
}

// connects to a local SSH agent
func connectToSSHAgent() agent.Agent {
    socketPath := os.Getenv("SSH_AUTH_SOCK")
    if socketPath == "" {
        log.Infof("SSH_AUTH_SOCK is not set. Is local SSH agent running?")
        return nil
    }
    log.Info("socketPath %s", socketPath)
    conn, err := net.Dial("unix", socketPath)
    if err != nil {
        log.Errorf("Failed connecting to local SSH agent via %s", socketPath)
        return nil
    }
    return agent.NewClient(conn)
}

// Username returns the current user's username
func Username() string {
    u, err := user.Current()
    if err != nil {
        utils.FatalError(err)
    }
    return u.Username
}

// AskPasswordAndHOTP prompts the user to enter the password + HTOP 2nd factor
func (tc *PocketClient) AskPasswordAndHOTP() (pwd string, token string, err error) {
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

// passwordFromConsole reads from stdin without echoing typed characters to stdout
func passwordFromConsole() (string, error) {
    fd := syscall.Stdin
    state, err := terminal.GetState(fd)

    // intercept Ctr+C and restore terminal
    sigCh := make(chan os.Signal, 1)
    closeCh := make(chan int)
    if err != nil {
        log.Warnf("failed reading terminal state: %v", err)
    } else {
        signal.Notify(sigCh, syscall.SIGINT)
        go func() {
            select {
            case <-sigCh:
                terminal.Restore(fd, state)
                os.Exit(1)
            case <-closeCh:
            }
        }()
    }
    defer func() {
        close(closeCh)
    }()

    bytes, err := terminal.ReadPassword(fd)
    return string(bytes), err
}

// lineFromConsole reads a line from stdin
func lineFromConsole() (string, error) {
    bytes, _, err := bufio.NewReader(os.Stdin).ReadLine()
    return string(bytes), err
}

// ParseLabelSpec parses a string like 'name=value,"long name"="quoted value"` into a map like
// { "name" -> "value", "long name" -> "quoted value" }
func ParseLabelSpec(spec string) (map[string]string, error) {
    tokens := []string{}
    var openQuotes = false
    var tokenStart, assignCount int
    var specLen = len(spec)
    // tokenize the label spec:
    for i, ch := range spec {
        endOfToken := false
        // end of line?
        if i+1 == specLen {
            i++
            endOfToken = true
        }
        switch ch {
        case '"':
            openQuotes = !openQuotes
        case '=', ',', ';':
            if !openQuotes {
                endOfToken = true
                if ch == '=' {
                    assignCount++
                }
            }
        }
        if endOfToken && i > tokenStart {
            tokens = append(tokens, strings.TrimSpace(strings.Trim(spec[tokenStart:i], `"`)))
            tokenStart = i + 1
        }
    }
    // simple validation of tokenization: must have an even number of tokens (because they're pairs)
    // and the number of such pairs must be equal the number of assignments
    if len(tokens)%2 != 0 || assignCount != len(tokens)/2 {
        return nil, fmt.Errorf("invalid label spec: '%s', should be 'key=value'", spec)
    }
    // break tokens in pairs and put into a map:
    labels := make(map[string]string)
    for i := 0; i < len(tokens); i += 2 {
        labels[tokens[i]] = tokens[i+1]
    }
    return labels, nil
}

func authMethodFromAgent(ag agent.Agent) ssh.AuthMethod {
    return ssh.PublicKeysCallback(ag.Signers)
}

// Executes the given command on the client machine (localhost). If no command is given,
// executes shell
func runLocalCommand(command []string) error {
    if len(command) == 0 {
        user, err := user.Current()
        if err != nil {
            return trace.Wrap(err)
        }
        shell, err := utils.GetLoginShell(user.Username)
        if err != nil {
            return trace.Wrap(err)
        }
        command = []string{shell}
    }
    cmd := exec.Command(command[0], command[1:]...)
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    return cmd.Run()
}