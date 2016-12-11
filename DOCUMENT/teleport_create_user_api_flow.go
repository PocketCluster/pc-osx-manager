//// ---- BEGIN TELEPORT USER SIGNUP TOKEN ---- ////
type UserCommand struct {
	config        *service.Config
	login         string
	allowedLogins string
	identities    []string
}
cmdUsers := UserCommand{config: cfg}
	case userAdd.FullCommand():
		err = cmdUsers.Add(client)

	// *** IMPORTANT ***
	// connectToAuthService creates a valid client connection to the auth service
	func connectToAuthService(cfg *service.Config) (client *auth.TunClient, err error) {
		// connect to the local auth server by default:
		cfg.Auth.Enabled = true
		if len(cfg.AuthServers) == 0 {
			cfg.AuthServers = []utils.NetAddr{
				*defaults.AuthConnectAddr(),
			}
		}
		// read the host SSH keys and use them to open an SSH connection to the auth service
		i, err := auth.ReadIdentity(cfg.DataDir, auth.IdentityID{Role: teleport.RoleAdmin, HostUUID: cfg.HostUUID})
		if err != nil {
			return nil, trace.Wrap(err)
		}
		client, err = auth.NewTunClient(
			"tctl",
			cfg.AuthServers,
			cfg.HostUUID,
			[]ssh.AuthMethod{ssh.PublicKeys(i.KeySigner)})
		if err != nil {
			return nil, trace.Wrap(err)
		}

		// check connectivity by calling something on a clinet:
		_, err = client.GetDialer()()
		if err != nil {
			utils.Consolef(os.Stderr,
				"Cannot connect to the auth server: %v.\nIs the auth server running on %v?", err, cfg.AuthServers[0].Addr)
			os.Exit(1)
		}
		return client, nil
	}
	client, err := connectToAuthService(cfg)
	if err != nil {
		utils.FatalError(err)
	}

		// TeleportUser is an optional user entry in the database
		type TeleportUser struct {
			// Name is a user name
			Name string `json:"name"`

			// AllowedLogins represents a list of OS users this teleport
			// user is allowed to login as
			AllowedLogins []string `json:"allowed_logins"`

			// OIDCIdentities lists associated OpenID Connect identities
			// that let user log in using externally verified identity
			OIDCIdentities []OIDCIdentity `json:"oidc_identities"`
		}

		// Add creates a new sign-up token and prints a token URL to stdout.
		// A user is not created until he visits the sign-up URL and completes the process
		func (u *UserCommand) Add(client *auth.TunClient) error {
			// if no local logins were specified, default to 'login'
			if u.allowedLogins == "" {
				u.allowedLogins = u.login
			}
			user := services.TeleportUser{
				Name:          u.login,
				AllowedLogins: strings.Split(u.allowedLogins, ","),
			}
			if len(u.identities) != 0 {
				for _, identityVar := range u.identities {
					vals := strings.SplitN(identityVar, ":", 2)
					if len(vals) != 2 {
						return trace.Errorf("bad flag --identity=%v, expected <connector-id>:<email> format", identityVar)
					}
					user.OIDCIdentities = append(user.OIDCIdentities, services.OIDCIdentity{ConnectorID: vals[0], Email: vals[1]})
				}
			}
			token, err := client.CreateSignupToken(&user)
			if err != nil {
				return err
			}
			proxies, err := client.GetProxies()
			if err != nil {
				return trace.Wrap(err)
			}
			hostname := "teleport-proxy"
			if len(proxies) == 0 {
				fmt.Printf("\x1b[1mWARNING\x1b[0m: this Teleport cluster does not have any proxy servers online.\nYou need to start some to be able to login.\n\n")
			} else {
				hostname = proxies[0].Hostname
			}

			// try to auto-suggest the activation link
			_, proxyPort, err := net.SplitHostPort(u.config.Proxy.WebAddr.Addr)
			if err != nil {
				proxyPort = strconv.Itoa(defaults.HTTPListenPort)
			}
			url := web.CreateSignupLink(net.JoinHostPort(hostname, proxyPort), token)
			fmt.Printf("Signup token has been created and is valid for %v seconds. Share this URL with the user:\n%v\n\nNOTE: make sure '%s' is accessible!\n", defaults.MaxSignupTokenTTL.Seconds(), url, hostname)
			return nil
		}

			// User represents teleport or external user
			type User interface {
				// GetName returns user name
				GetName() string
				// GetAllowedLogins returns user's allowed linux logins
				GetAllowedLogins() []string
				// GetIdentities returns a list of connected OIDCIdentities
				GetIdentities() []OIDCIdentity
				// String returns user
				String() string
				// Check checks if all parameters are correct
				Check() error
				// Equals checks if user equals to another
				Equals(other User) bool
			}

			// CreateSignupToken creates one time token for creating account for the user
			// For each token it creates username and hotp generator
			func (c *Client) CreateSignupToken(user services.User) (string, error) {
				if err := user.Check(); err != nil {
					return "", trace.Wrap(err)
				}
				out, err := c.PostJSON(c.Endpoint("signuptokens"), createSignupTokenReq{
					User: user,
				})
				if err != nil {
					return "", trace.Wrap(err)
				}
				var token string
				if err := json.Unmarshal(out.Bytes(), &token); err != nil {
					return "", trace.Wrap(err)
				}
				return token, nil
			}

// -- BEGIN SERVER -- //
srv.POST("/v1/signuptokens", httplib.MakeHandler(srv.createSignupToken))

func (s *APIServer) CreateSignupToken(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
	var req *createSignupTokenReqRaw
	if err := httplib.ReadJSON(r, &req); err != nil {
		return nil, trace.Wrap(err)
	}
	user, err := services.GetUserUnmarshaler()(req.User)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	token, err := s.a.CreateSignupToken(user)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return token, nil
}

	func (a *AuthWithRoles) CreateSignupToken(user services.User) (token string, e error) {
		if err := a.permChecker.HasPermission(a.role, ActionCreateSignupToken); err != nil {
			return "", trace.Wrap(err)
		}
		return a.authServer.CreateSignupToken(user)

	}

		// TeleportUser is an optional user entry in the database
		type TeleportUser struct {
			// Name is a user name
			Name string `json:"name"`

			// AllowedLogins represents a list of OS users this teleport
			// user is allowed to login as
			AllowedLogins []string `json:"allowed_logins"`

			// OIDCIdentities lists associated OpenID Connect identities
			// that let user log in using externally verified identity
			OIDCIdentities []OIDCIdentity `json:"oidc_identities"`
		}

		// CreateSignupToken creates one time token for creating account for the user
		// For each token it creates username and hotp generator
		//
		// allowedLogins are linux user logins allowed for the new user to use
		func (s *AuthServer) CreateSignupToken(user services.User) (string, error) {
			if err := user.Check(); err != nil {
				return "", trace.Wrap(err)
			}
			// make sure that connectors actually exist
			for _, id := range user.GetIdentities() {
				if err := id.Check(); err != nil {
					return "", trace.Wrap(err)
				}
				if _, err := s.GetOIDCConnector(id.ConnectorID, false); err != nil {
					return "", trace.Wrap(err)
				}
			}
			// check existing
			_, err := s.GetPasswordHash(user.GetName())
			if err == nil {
				return "", trace.BadParameter("user '%v' already exists", user)
			}

			token, err := utils.CryptoRandomHex(TokenLenBytes)
			if err != nil {
				return "", trace.Wrap(err)
			}

			otp, err := hotp.GenerateHOTP(defaults.HOTPTokenDigits, false)
			if err != nil {
				log.Errorf("[AUTH API] failed to generate HOTP: %v", err)
				return "", trace.Wrap(err)
			}
			otpQR, err := otp.QR("Teleport: " + user.GetName() + "@" + s.AuthServiceName)
			if err != nil {
				return "", trace.Wrap(err)
			}

			otpMarshalled, err := hotp.Marshal(otp)
			if err != nil {
				return "", trace.Wrap(err)
			}

			otpFirstValues := make([]string, defaults.HOTPFirstTokensRange)
			for i := 0; i < defaults.HOTPFirstTokensRange; i++ {
				otpFirstValues[i] = otp.OTP()
			}

			tokenData := services.SignupToken{
				Token: token,
				User: services.TeleportUser{
					Name:           user.GetName(),
					AllowedLogins:  user.GetAllowedLogins(),
					OIDCIdentities: user.GetIdentities()},
				Hotp:            otpMarshalled,
				HotpFirstValues: otpFirstValues,
				HotpQR:          otpQR,
			}

			err = s.UpsertSignupToken(token, tokenData, defaults.MaxSignupTokenTTL)
			if err != nil {
				return "", trace.Wrap(err)
			}

			log.Infof("[AUTH API] created the signup token for %v as %v", user)
			return token, nil
		}
// -- END SERVER -- //

			// CreateSignupLink generates and returns a URL which is given to a new
			// user to complete registration with Teleport via Web UI
			func CreateSignupLink(hostPort string, token string) string {
				return "https://" + hostPort + "/web/newuser/" + token
			}

//// ---- END TELEPORT USER SIGNUP TOKEN ---- ////





//// ---- BEGIN TELEPORT WEB USER INVITE ---- ////
newUser: '/web/newuser/:inviteToken',
	invitePath: '/v1/webapi/users/invites/:inviteToken',
	    getInviteUrl: function getInviteUrl(inviteToken) {
	      return formatPattern(cfg.api.invitePath, { inviteToken: inviteToken });
	    },
// -- BEGIN SERVER -- //
			h.GET("/webapi/users/invites/:token", httplib.MakeHandler(h.renderUserInvite))

			// renderUserInvite is called to show user the new user invitation page
			//
			// GET /v1/webapi/users/invites/:token
			//
			// Response:
			//
			// {"invite_token": "token", "user": "alex", qr: "base64-encoded-qr-code image"}
			//
			//
			func (m *Handler) renderUserInvite(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
				token := p[0].Value
				user, QRCodeBytes, _, err := m.auth.GetUserInviteInfo(token)
				if err != nil {
					return nil, trace.Wrap(err)
				}

				return &renderUserInviteResponse{
					InviteToken: token,
					User:        user,
					QR:          QRCodeBytes,
				}, nil
			}
// -- END SERVER -- //

//// ---- END TELEPORT WEB USER INVITE ---- ////





//// ---- BEGIN TELEPORT SERVER USER CREATION ---- ////
createUserPath: '/v1/webapi/users',

// -- BEGIN WEB API -- // (lib/web/web.go)

h.POST("/webapi/users", httplib.MakeHandler(h.createNewUser))

// createNewUser creates new user entry based on the invite token
//
// POST /v1/webapi/users
//
// {"invite_token": "unique invite token", "pass": "user password", "second_factor_token": "valid second factor token"}
//
// Sucessful response: (session cookie is set)
//
// {"type": "bearer", "token": "bearer token", "user": "alex", "expires_in": 20}
func (m *Handler) createNewUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
	var req *createNewUserReq
	if err := httplib.ReadJSON(r, &req); err != nil {
		return nil, trace.Wrap(err)
	}
	sess, err := m.auth.CreateNewUser(req.InviteToken, req.Pass, req.SecondFactorToken)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	ctx, err := m.auth.ValidateSession(sess.Username, sess.ID)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	if err := SetSession(w, sess.Username, sess.ID); err != nil {
		return nil, trace.Wrap(err)
	}
	return NewSessionResponse(ctx)
}

	func (s *sessionCache) CreateNewUser(token, password, hotpToken string) (*auth.Session, error) {
		method, err := auth.NewSignupTokenAuth(token)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		clt, err := auth.NewTunClient("web.create-user", s.authServers, "tokenAuth", method)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		defer clt.Close()
		sess, err := clt.CreateUserWithToken(token, password, hotpToken)
		return sess, trace.Wrap(err)
	}

		// NewTunClient returns an instance of new HTTP client to Auth server API
		// exposed over SSH tunnel, so client  uses SSH credentials to dial and authenticate
		//  - purpose is mostly for debuggin, like "web client" or "reverse tunnel client"
		//  - authServers: list of auth servers in this cluster (they are supposed to be in sync)
		//  - authMethods: how to authenticate (via cert, web passwowrd, etc)
		//  - opts : functional arguments for further extending
		func NewTunClient(purpose string,
			authServers []utils.NetAddr,
			user string,
			authMethods []ssh.AuthMethod,
			opts ...TunClientOption) (*TunClient, error) {
			if user == "" {
				return nil, trace.BadParameter("SSH connection requires a valid username")
			}
			tc := &TunClient{
				purpose:           purpose,
				user:              user,
				staticAuthServers: authServers,
				authMethods:       authMethods,
				closeC:            make(chan struct{}),
			}
			for _, o := range opts {
				o(tc)
			}
			log.Debugf("newTunClient(%s) with auth: %v", purpose, authServers)

			clt, err := NewClient("http://stub:0", tc.Dial)
			if err != nil {
				return nil, trace.Wrap(err)
			}
			tc.Client = *clt

			// use local information about auth servers if it's available
			if tc.addrStorage != nil {
				cachedAuthServers, err := tc.addrStorage.GetAddresses()
				if err != nil {
					log.Infof("unable to load the auth server cache: %v", err)
				} else {
					tc.setAuthServers(cachedAuthServers)
				}
			}
			return tc, nil
		}

		// CreateUserWithToken creates account with provided token and password.
		// Account username and hotp generator are taken from token data.
		// Deletes token after account creation.
		func (c *Client) CreateUserWithToken(token, password, hotpToken string) (*Session, error) {
			out, err := c.PostJSON(c.Endpoint("signuptokens", "users"), createUserWithTokenReq{
				Token:     token,
				Password:  password,
				HOTPToken: hotpToken,
			})
			if err != nil {
				return nil, trace.Wrap(err)
			}
			var sess *Session
			if err := json.Unmarshal(out.Bytes(), &sess); err != nil {
				return nil, trace.Wrap(err)
			}
			return sess, nil
		}

			// PostJSON is a generic method that issues http POST request to the server
			func (c *Client) PostJSON(
				endpoint string, val interface{}) (*roundtrip.Response, error) {
				return httplib.ConvertResponse(c.Client.PostJSON(endpoint, val))
			}
// -- END WEB API -- // (lib/web/web.go)

// -- BEGIN SERVER -- // (lib/auth/apiserver.go)
/*				
	// signup token data is utterly useless
	srv.GET("/v1/signuptokens/:token", httplib.MakeHandler(srv.getSignupTokenData))

	// getSignupTokenData auth API method creates a new sign-up token for adding a new user
	func (s *APIServer) getSignupTokenData(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
		token := p[0].Value

		user, QRImg, hotpFirstValues, err := s.a.GetSignupTokenData(token)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		return &getSignupTokenDataResponse{
			User:            user,
			QRImg:           QRImg,
			HotpFirstValues: hotpFirstValues,
		}, nil
	}
*/

srv.POST("/v1/signuptokens/users", httplib.MakeHandler(srv.createUserWithToken))

type createUserWithTokenReq struct {
	Token     string `json:"token"`
	Password  string `json:"password"`
	HOTPToken string `json:"hotp_token"`
}

func (s *APIServer) createUserWithToken(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
	var req *createUserWithTokenReq
	if err := httplib.ReadJSON(r, &req); err != nil {
		return nil, trace.Wrap(err)
	}
	sess, err := s.a.CreateUserWithToken(req.Token, req.Password, req.HOTPToken)
	if err != nil {
		log.Error(err)
		return nil, trace.Wrap(err)
	}
	return sess, nil
}

	func (a *AuthWithRoles) CreateUserWithToken(token, password, hotpToken string) (*Session, error) {
		if err := a.permChecker.HasPermission(a.role, ActionCreateUserWithToken); err != nil {
			return nil, trace.Wrap(err)
		}
		return a.authServer.CreateUserWithToken(token, password, hotpToken)
	}

		// CreateUserWithToken creates account with provided token and password.
		// Account username and hotp generator are taken from token data.
		// Deletes token after account creation.
		func (s *AuthServer) CreateUserWithToken(token, password, hotpToken string) (*Session, error) {
			err := s.AcquireLock("signuptoken"+token, time.Hour)
			if err != nil {
				return nil, trace.Wrap(err)
			}

			defer func() {
				err := s.ReleaseLock("signuptoken" + token)
				if err != nil {
					log.Errorf(err.Error())
				}
			}()

			tokenData, err := s.GetSignupToken(token)
			if err != nil {
				return nil, trace.Wrap(err)
			}

			otp, err := hotp.Unmarshal(tokenData.Hotp)
			if err != nil {
				return nil, trace.Wrap(err)
			}

			ok := otp.Scan(hotpToken, defaults.HOTPFirstTokensRange)
			if !ok {
				return nil, trace.BadParameter("wrong HOTP token")
			}

			_, _, err = s.UpsertPassword(tokenData.User.GetName(), []byte(password))
			if err != nil {
				return nil, trace.Wrap(err)
			}

			// apply user allowed logins
			if err = s.UpsertUser(&tokenData.User); err != nil {
				return nil, trace.Wrap(err)
			}

			err = s.UpsertHOTP(tokenData.User.GetName(), otp)
			if err != nil {
				return nil, trace.Wrap(err)
			}

			log.Infof("[AUTH] created new user: %v", &tokenData.User)

			if err = s.DeleteSignupToken(token); err != nil {
				return nil, trace.Wrap(err)
			}

			sess, err := s.NewWebSession(tokenData.User.GetName())
			if err != nil {
				return nil, trace.Wrap(err)
			}

			err = s.UpsertWebSession(tokenData.User.GetName(), sess, WebSessionTTL)
			if err != nil {
				return nil, trace.Wrap(err)
			}

			sess.WS.Priv = nil
			return sess, nil
		}

	func (s *sessionCache) ValidateSession(user, sid string) (*SessionContext, error) {
		ctx, err := s.getContext(user, sid)
		if err == nil {
			return ctx, nil
		}
		log.Debugf("ValidateSession(%s, %s)", user, sid)
		method, err := auth.NewWebSessionAuth(user, []byte(sid))
		if err != nil {
			return nil, trace.Wrap(err)
		}
		// Note: do not close this auth API client now. It will exist inside of "session context"
		clt, err := auth.NewTunClient("web.session-user", s.authServers, user, method)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		sess, err := clt.GetWebSessionInfo(user, sid)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		c := &SessionContext{
			clt:    clt,
			user:   user,
			sess:   sess,
			parent: s,
		}
		c.Entry = log.WithFields(log.Fields{
			"user": user,
			"sess": sess.ID[:4],
		})

		out, err := s.insertContext(user, sid, c, auth.WebSessionTTL)
		if err != nil {
			// this means that someone has just inserted the context, so
			// close our extra context and return
			if trace.IsAlreadyExists(err) {
				log.Infof("just created, returning the existing one")
				defer c.Close()
				return out, nil
			}
			return nil, trace.Wrap(err)
		}
		return out, nil
	}
// -- END SERVER -- //

//// ---- END TELEPORT SERVER USER APPROVAL ---- ////


srv.POST("/v1/users", httplib.MakeHandler(srv.upsertUser))

func (s *APIServer) upsertUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
	var req *upsertUserReqRaw
	if err := httplib.ReadJSON(r, &req); err != nil {
		return nil, trace.Wrap(err)
	}
	user, err := services.GetUserUnmarshaler()(req.User)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	err = s.a.UpsertUser(user)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return message(fmt.Sprintf("'%v' user upserted", user.GetName())), nil
}

	func (a *AuthWithRoles) UpsertUser(u services.User) error {
		if err := a.permChecker.HasPermission(a.role, ActionUpsertUser); err != nil {
			return trace.Wrap(err)
		}
		return a.authServer.UpsertUser(u)
	}

		// UpsertUser updates parameters about user
		func (s *IdentityService) UpsertUser(user services.User) error {
			if !cstrings.IsValidUnixUser(user.GetName()) {
				return trace.BadParameter("'%v is not a valid unix username'", user.GetName())
			}

			for _, l := range user.GetAllowedLogins() {
				if !cstrings.IsValidUnixUser(l) {
					return trace.BadParameter("'%v is not a valid unix username'", l)
				}
			}
			for _, i := range user.GetIdentities() {
				if err := i.Check(); err != nil {
					return trace.Wrap(err)
				}
			}
			data, err := json.Marshal(user)
			if err != nil {
				return trace.Wrap(err)
			}

			err = s.backend.UpsertVal([]string{"web", "users", user.GetName()}, "params", []byte(data), backend.Forever)
			if err != nil {
				return trace.Wrap(err)
			}
			return nil
		}




