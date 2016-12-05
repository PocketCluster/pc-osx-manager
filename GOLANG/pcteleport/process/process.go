package process

import (
    "net"
    "fmt"
    "time"
    "sync"
    "net/http"
    "path/filepath"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/auth"
    "github.com/gravitational/teleport/lib/events"
    "github.com/gravitational/teleport/lib/session"
    "github.com/gravitational/teleport/lib/limiter"
    "github.com/gravitational/teleport/lib/services"
    "github.com/gravitational/teleport/lib/backend"
    "github.com/gravitational/teleport/lib/backend/boltbk"
    "github.com/gravitational/teleport/lib/reversetunnel"
    "github.com/gravitational/teleport/lib/srv"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/web"

    "github.com/stkim1/pcteleport/pcdefaults"
    "github.com/stkim1/pcteleport/tun"
    "golang.org/x/crypto/ssh"
)

// TeleportProcess structure holds the state of the Teleport daemon, controlling
// execution and configuration of the teleport services: ssh, auth and proxy.
type PocketCoreTeleportProcess struct {
    sync.Mutex
    service.Supervisor
    Config *service.Config
    // localAuth has local auth server listed in case if this process
    // has started with auth server role enabled
    localAuth *auth.AuthServer
}

// initAuthStorage initializes the storage backend for the auth. service
func (process *PocketCoreTeleportProcess) initAuthStorage() (backend.Backend, error) {
    cfg := &process.Config.Auth
    var bk backend.Backend
    var err error

    switch cfg.KeysBackend.Type {
    case teleport.BoltBackendType:
        bk, err = boltbk.FromJSON(cfg.KeysBackend.Params)
    default:
        return nil, trace.Errorf("unsupported backend type: %v", cfg.KeysBackend.Type)
    }
    if err != nil {
        return nil, trace.Wrap(err)
    }

    return bk, nil
}

func (process *PocketCoreTeleportProcess) findStaticIdentity(id auth.IdentityID) (*auth.Identity, error) {
    for i := range process.Config.Identities {
        identity := process.Config.Identities[i]
        if identity.ID.Equals(id) {
            return identity, nil
        }
    }
    return nil, trace.NotFound("identity %v not found", &id)
}

// connectToAuthService attempts to login into the auth servers specified in the
// configuration. Returns 'true' if successful
func (process *PocketCoreTeleportProcess) connectToAuthService(role teleport.Role) (*service.Connector, error) {
    id := auth.IdentityID{HostUUID: process.Config.HostUUID, Role: role}
    identity, err := auth.ReadIdentity(process.Config.DataDir, id)
    if err != nil {
        if trace.IsNotFound(err) {
            // try to locate static identity provide in the file
            identity, err = process.findStaticIdentity(id)
            if err != nil {
                return nil, trace.Wrap(err)
            }
            log.Infof("found static identity %v in the config file, writing to disk", &id)
            if err = auth.WriteIdentity(process.Config.DataDir, identity); err != nil {
                return nil, trace.Wrap(err)
            }
        } else {
            return nil, trace.Wrap(err)
        }
    }
    storage := utils.NewFileAddrStorage(
        filepath.Join(process.Config.DataDir, "authservers.json"))

    authUser := identity.Cert.ValidPrincipals[0]
    authClient, err := auth.NewTunClient(
        string(role),
        process.Config.AuthServers,
        authUser,
        []ssh.AuthMethod{ssh.PublicKeys(identity.KeySigner)},
        auth.TunClientStorage(storage),
    )
    // success?
    if err != nil {
        return nil, trace.Wrap(err)
    }
    // try calling a test method via auth api:
    //
    // ??? in case of failure it never gets back here!!!
    dn, err := authClient.GetLocalDomain()
    if err != nil {
        return nil, trace.Wrap(err)
    }
    // success ? we're logged in!
    log.Infof("[Node] %s connected to the cluster '%s'", authUser, dn)
    return &service.Connector{Client: authClient, Identity: identity}, nil
}

func (process *PocketCoreTeleportProcess) setLocalAuth(a *auth.AuthServer) {
    process.Lock()
    defer process.Unlock()
    process.localAuth = a
}

func (process *PocketCoreTeleportProcess) getLocalAuth() (a *auth.AuthServer) {
    process.Lock()
    defer process.Unlock()
    return process.localAuth
}

// onExit allows individual services to register a callback function which will be
// called when Teleport Process is asked to exit. Usually services terminate themselves
// when the callback is called
func (process *PocketCoreTeleportProcess) onExit(callback func(interface{})) {
    go func() {
        eventC := make(chan service.Event)
        process.WaitForEvent(service.TeleportExitEvent, eventC, make(chan struct{}))
        select {
        case event := <-eventC:
            callback(event.Payload)
        }
    }()
}

// initAuthService can be called to initialize auth server service
func (process *PocketCoreTeleportProcess) initAuthService(authority auth.Authority) error {
    var (
        askedToExit = false
        err         error
    )
    cfg := process.Config
    // Initialize the storage back-ends for keys, events and records
    b, err := process.initAuthStorage()
    if err != nil {
        return trace.Wrap(err)
    }

    // create the audit log, which will be consuming (and recording) all events
    // and record sessions
    var auditLog events.IAuditLog
    if cfg.Auth.NoAudit {
        auditLog = &events.DiscardAuditLog{}
        log.Warn("the audit and session recording are turned off")
    } else {
        auditLog, err = events.NewAuditLog(filepath.Join(cfg.DataDir, "log"))
        if err != nil {
            return trace.Wrap(err)
        }
    }

    // first, create the AuthServer
    authServer, identity, err := auth.Init(auth.InitConfig{
        Backend:         b,
        Authority:       authority,
        DomainName:      cfg.Auth.DomainName,
        AuthServiceName: cfg.Hostname,
        DataDir:         cfg.DataDir,
        HostUUID:        cfg.HostUUID,
        Authorities:     cfg.Auth.Authorities,
        ReverseTunnels:  cfg.ReverseTunnels,
        OIDCConnectors:  cfg.OIDCConnectors,
        Trust:           cfg.Trust,
        Lock:            cfg.Lock,
        Presence:        cfg.Presence,
        Provisioner:     cfg.Provisioner,
        Identity:        cfg.Identity,
        StaticTokens:    cfg.Auth.StaticTokens,
    }, cfg.SeedConfig)
    if err != nil {
        return trace.Wrap(err)
    }
    process.setLocalAuth(authServer)

    // second, create the API Server: it's actually a collection of API servers,
    // each serving requests for a "role" which is assigned to every connected
    // client based on their certificate (user, server, admin, etc)
    sessionService, err := session.New(b)
    if err != nil {
        return trace.Wrap(err)
    }
    apiConf := &auth.APIConfig{
        AuthServer:        authServer,
        SessionService:    sessionService,
        PermissionChecker: auth.NewStandardPermissions(),
        AuditLog:          auditLog,
    }

    limiter, err := limiter.NewLimiter(cfg.Auth.Limiter)
    if err != nil {
        return trace.Wrap(err)
    }

    // Register an SSH endpoint which is used to create an SSH tunnel to send HTTP
    // requests to the Auth API
    var authTunnel *tun.AuthTunnel
    process.RegisterFunc(func() error {
        utils.Consolef(cfg.Console, "[AUTH]  Auth service is starting on %v", cfg.Auth.SSHAddr.Addr)
        authTunnel, err = tun.NewTunnel(
            cfg.Auth.SSHAddr,
            identity.KeySigner,
            apiConf,
            tun.SetLimiter(limiter),
        )
        if err != nil {
            utils.Consolef(cfg.Console, "[AUTH] Error: %v", err)
            return trace.Wrap(err)
        }
        if err := authTunnel.Start(); err != nil {
            if askedToExit {
                log.Infof("[AUTH] Auth Tunnel exited")
                return nil
            }
            utils.Consolef(cfg.Console, "[AUTH] Error: %v", err)
            return trace.Wrap(err)
        }
        return nil
    })

    process.RegisterFunc(func() error {
        // Heart beat auth server presence, this is not the best place for this
        // logic, consolidate it into auth package later
        connector, err := process.connectToAuthService(teleport.RoleAdmin)
        if err != nil {
            return trace.Wrap(err)
        }
        // External integrations rely on this event:
        process.BroadcastEvent(service.Event{Name: service.AuthIdentityEvent, Payload: connector})
        process.onExit(func(payload interface{}) {
            connector.Client.Close()
        })
        return nil
    })

    process.RegisterFunc(func() error {
        srv := services.Server{
            ID:       process.Config.HostUUID,
            Addr:     cfg.Auth.SSHAddr.Addr,
            Hostname: process.Config.Hostname,
        }
        host, port, err := net.SplitHostPort(srv.Addr)
        // advertise-ip is explicitly set:
        if process.Config.AdvertiseIP != nil {
            if err != nil {
                return trace.Wrap(err)
            }
            srv.Addr = fmt.Sprintf("%v:%v", process.Config.AdvertiseIP.String(), port)
        } else {
            // advertise-ip is not set, while the CA is listening on 0.0.0.0? lets try
            // to guess the 'advertise ip' then:
            if net.ParseIP(host).IsUnspecified() {
                ip, err := utils.GuessHostIP()
                if err != nil {
                    log.Warn(err)
                } else {
                    srv.Addr = net.JoinHostPort(ip.String(), port)
                }
            }
            log.Warnf("advertise_ip is not set for this auth server!!! Trying to guess the IP this server can be reached at: %v", srv.Addr)
        }
        // immediately register, and then keep repeating in a loop:
        for !askedToExit {
            err := authServer.UpsertAuthServer(srv, pcdefaults.ServerHeartbeatTTL)
            if err != nil {
                log.Warningf("failed to announce presence: %v", err)
            }
            sleepTime := pcdefaults.ServerHeartbeatTTL/2 + utils.RandomDuration(pcdefaults.ServerHeartbeatTTL/10)
            time.Sleep(sleepTime)
        }
        log.Infof("[AUTH] heartbeat to other auth servers exited")
        return nil
    })

    // execute this when process is asked to exit:
    process.onExit(func(payload interface{}) {
        askedToExit = true
        authTunnel.Close()
        log.Infof("[AUTH] auth service exited")
    })
    return nil
}

func (process *PocketCoreTeleportProcess) initProxyEndpoint(conn *service.Connector) error {
    var (
        askedToExit = true
        err         error
    )
    cfg := process.Config
    proxyLimiter, err := limiter.NewLimiter(cfg.Proxy.Limiter)
    if err != nil {
        return trace.Wrap(err)
    }

    reverseTunnelLimiter, err := limiter.NewLimiter(cfg.Proxy.Limiter)
    if err != nil {
        return trace.Wrap(err)
    }

    tsrv, err := reversetunnel.NewServer(
        cfg.Proxy.ReverseTunnelListenAddr,
        []ssh.Signer{conn.Identity.KeySigner},
        conn.Client,
        reversetunnel.SetLimiter(reverseTunnelLimiter),
        reversetunnel.DirectSite(conn.Identity.Cert.Extensions[utils.CertExtensionAuthority], conn.Client),
    )
    if err != nil {
        return trace.Wrap(err)
    }

    SSHProxy, err := srv.New(cfg.Proxy.SSHAddr,
        cfg.Hostname,
        []ssh.Signer{conn.Identity.KeySigner},
        conn.Client,
        cfg.DataDir,
        nil,
        srv.SetLimiter(proxyLimiter),
        srv.SetProxyMode(tsrv),
        srv.SetSessionServer(conn.Client),
        srv.SetAuditLog(conn.Client),
    )
    if err != nil {
        return trace.Wrap(err)
    }

    // Register reverse tunnel agents pool
    agentPool, err := reversetunnel.NewAgentPool(reversetunnel.AgentPoolConfig{
        HostUUID:    conn.Identity.ID.HostUUID,
        Client:      conn.Client,
        HostSigners: []ssh.Signer{conn.Identity.KeySigner},
    })
    if err != nil {
        return trace.Wrap(err)
    }

    // register SSH reverse tunnel server that accepts connections
    // from remote teleport nodes
    process.RegisterFunc(func() error {
        utils.Consolef(cfg.Console, "[PROXY] Reverse tunnel service is starting on %v", cfg.Proxy.ReverseTunnelListenAddr.Addr)
        if err := tsrv.Start(); err != nil {
            utils.Consolef(cfg.Console, "[PROXY] Error: %v", err)
            return trace.Wrap(err)
        }
        // notify parties that we've started reverse tunnel server
        process.BroadcastEvent(service.Event{Name: service.ProxyReverseTunnelServerEvent, Payload: tsrv})
        tsrv.Wait()
        if askedToExit {
            log.Infof("[PROXY] Reverse tunnel exited")
        }
        return nil
    })

    // Register web proxy server
    var webListener net.Listener

    if !process.Config.Proxy.DisableWebUI {
        process.RegisterFunc(func() error {
            utils.Consolef(cfg.Console, "[PROXY] Web proxy service is starting on %v", cfg.Proxy.WebAddr.Addr)
            webHandler, err := web.NewHandler(
                web.Config{
                    Proxy:       tsrv,
                    AssetsDir:   cfg.Proxy.AssetsDir,
                    AuthServers: cfg.AuthServers[0],
                    DomainName:  cfg.Hostname,
                    ProxyClient: conn.Client,
                    DisableUI:   cfg.Proxy.DisableWebUI,
                })
            if err != nil {
                utils.Consolef(cfg.Console, "[PROXY] starting the web server: %v", err)
                return trace.Wrap(err)
            }
            defer webHandler.Close()

            proxyLimiter.WrapHandle(webHandler)
            process.BroadcastEvent(service.Event{Name: service.ProxyWebServerEvent, Payload: webHandler})

            log.Infof("[PROXY] init TLS listeners")
            webListener, err = utils.ListenTLS(
                cfg.Proxy.WebAddr.Addr,
                cfg.Proxy.TLSCert,
                cfg.Proxy.TLSKey)
            if err != nil {
                return trace.Wrap(err)
            }
            if err = http.Serve(webListener, proxyLimiter); err != nil {
                if askedToExit {
                    log.Infof("[PROXY] web server exited")
                    return nil
                }
                log.Error(err)
            }
            return nil
        })
    } else {
        log.Infof("[WEB] Web UI is disabled")
    }

    // Register ssh proxy server
    process.RegisterFunc(func() error {
        utils.Consolef(cfg.Console, "[PROXY] SSH proxy service is starting on %v", cfg.Proxy.SSHAddr.Addr)
        if err := SSHProxy.Start(); err != nil {
            if askedToExit {
                log.Infof("[PROXY] SSH proxy exited")
                return nil
            }
            utils.Consolef(cfg.Console, "[PROXY] Error: %v", err)
            return trace.Wrap(err)
        }
        return nil
    })

    process.RegisterFunc(func() error {
        log.Infof("[PROXY] starting reverse tunnel agent pool")
        if err := agentPool.Start(); err != nil {
            log.Fatalf("failed to start: %v", err)
            return trace.Wrap(err)
        }
        agentPool.Wait()
        return nil
    })

    // execute this when process is asked to exit:
    process.onExit(func(payload interface{}) {
        tsrv.Close()
        SSHProxy.Close()
        agentPool.Stop()
        if webListener != nil {
            webListener.Close()
        }
        log.Infof("[PROXY] proxy service exited")
    })
    return nil
}

// initProxy gets called if teleport runs with 'proxy' role enabled.
// this means it will do two things:
//    1. serve a web UI
//    2. proxy SSH connections to nodes running with 'node' role
//    3. take care of revse tunnels
func (process *PocketCoreTeleportProcess) initProxy() error {
    // TODO : (11/28/2016) TLS Certificate should be generated in pc-core context initiation
    /*
    // if no TLS key was provided for the web UI, generate a self signed cert
    if process.Config.Proxy.TLSKey == "" && !process.Config.Proxy.DisableWebUI {
        //err := initSelfSignedHTTPSCert(process.Config)
        var err error = nil
        if err != nil {
            return trace.Wrap(err)
        }
    }
    */

    process.RegisterWithAuthServer(
        process.Config.Token, teleport.RoleProxy,
        service.ProxyIdentityEvent)

    process.RegisterFunc(func() error {
        eventsC := make(chan service.Event)
        process.WaitForEvent(service.ProxyIdentityEvent, eventsC, make(chan struct{}))

        event := <-eventsC
        log.Infof("[SSH] received %v", &event)
        conn, ok := (event.Payload).(*service.Connector)
        if !ok {
            return trace.BadParameter("unsupported connector type: %T", event.Payload)
        }
        return trace.Wrap(process.initProxyEndpoint(conn))
    })
    return nil
}


// RegisterWithAuthServer uses one time provisioning token obtained earlier
// from the server to get a pair of SSH keys signed by Auth server host
// certificate authority
func (process *PocketCoreTeleportProcess) RegisterWithAuthServer(token string, role teleport.Role, eventName string) {
    cfg := process.Config
    identityID := auth.IdentityID{Role: role, HostUUID: cfg.HostUUID}

    // this means the server has not been initialized yet, we are starting
    // the registering client that attempts to connect to the auth server
    // and provision the keys
    var authClient *auth.TunClient
    process.RegisterFunc(func() error {
        retryTime := pcdefaults.ServerHeartbeatTTL / 3
        for {
            connector, err := process.connectToAuthService(role)
            if err == nil {
                process.BroadcastEvent(service.Event{Name: eventName, Payload: connector})
                authClient = connector.Client
                return nil
            }
            if trace.IsConnectionProblem(err) {
                utils.Consolef(cfg.Console, "[%v] connecting to auth server: %v", role, err)
                time.Sleep(retryTime)
                continue
            }
            if !trace.IsNotFound(err) {
                return trace.Wrap(err)
            }
            //  we haven't connected yet, so we expect the token to exist
            if process.getLocalAuth() != nil {
                // Auth service is on the same host, no need to go though the invitation
                // procedure
                log.Debugf("[Node] this server has local Auth server started, using it to add role to the cluster")
                err = auth.LocalRegister(cfg.DataDir, identityID, process.getLocalAuth())
            } else {
                // Auth server is remote, so we need a provisioning token
                if token == "" {
                    return trace.BadParameter("%v must join a cluster and needs a provisioning token", role)
                }
                log.Infof("[Node] %v joining the cluster with a token %v", role, token)
                err = auth.Register(cfg.DataDir, token, identityID, cfg.AuthServers)
            }
            if err != nil {
                utils.Consolef(cfg.Console, "[%v] failed to join the cluster: %v", role, err)
                time.Sleep(retryTime)
            } else {
                utils.Consolef(cfg.Console, "[%v] Successfully registered with the cluster", role)
                continue
            }
        }
    })

    process.onExit(func(interface{}) {
        if authClient != nil {
            authClient.Close()
        }
    })
}


/*
// initSelfSignedHTTPSCert generates and self-signs a TLS key+cert pair for https connection
// to the proxy server.
func initSelfSignedHTTPSCert(cfg *Config) (err error) {
	log.Warningf("[CONFIG] NO TLS Keys provided, using self signed certificate")

	keyPath := filepath.Join(cfg.DataDir, defaults.SelfSignedKeyPath)
	certPath := filepath.Join(cfg.DataDir, defaults.SelfSignedCertPath)

	cfg.Proxy.TLSKey = keyPath
	cfg.Proxy.TLSCert = certPath

	// return the existing pair if they ahve already been generated:
	_, err = tls.LoadX509KeyPair(certPath, keyPath)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return trace.Wrap(err, "unrecognized error reading certs")
	}
	log.Warningf("[CONFIG] Generating self signed key and cert to %v %v", keyPath, certPath)

	creds, err := utils.GenerateSelfSignedCert([]string{cfg.Hostname, "localhost"})
	if err != nil {
		return trace.Wrap(err)
	}

	if err := ioutil.WriteFile(keyPath, creds.PrivateKey, 0600); err != nil {
		return trace.Wrap(err, "error writing key PEM")
	}
	if err := ioutil.WriteFile(certPath, creds.Cert, 0600); err != nil {
		return trace.Wrap(err, "error writing key PEM")
	}
	return nil
}
 */
