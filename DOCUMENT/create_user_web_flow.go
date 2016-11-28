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

				/* lib/auth/apiserver.go */
				srv.GET("/v1/signuptokens/:token", httplib.MakeHandler(srv.getSignupTokenData))
				srv.POST("/v1/signuptokens/users", httplib.MakeHandler(srv.createUserWithToken))

				func (s *APIServer) createSignupToken(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
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
