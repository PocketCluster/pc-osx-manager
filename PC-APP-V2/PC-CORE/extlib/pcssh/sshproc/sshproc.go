// +build darwin
package sshproc

import (
    "fmt"
    "net"
    "path/filepath"
    "time"

    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/auth"
    "github.com/gravitational/teleport/lib/auth/native"
    "github.com/gravitational/teleport/lib/defaults"
    "github.com/gravitational/teleport/lib/events"
    "github.com/gravitational/teleport/lib/session"
    "github.com/gravitational/teleport/lib/limiter"
    "github.com/gravitational/teleport/lib/services"
    "github.com/gravitational/teleport/lib/backend"
    "github.com/gravitational/teleport/lib/backend/sqlitebk"
    "github.com/gravitational/teleport/lib/reversetunnel"
    "github.com/gravitational/teleport/lib/srv"
    "github.com/gravitational/teleport/lib/service"

    "github.com/stkim1/pc-core/extlib/pcssh/sshcfg"
    pervice "github.com/stkim1/pc-core/service"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "golang.org/x/crypto/ssh"
)

const (
    eventProxyIdentity string       = "event.proxy.identity"
    eventAuthorityInstance string   = "event.authority.instance"
    eventAuthorityClientConn string = "event.authority.client.conn"

    ServicePCSSHAuthority string    = "service.pcssh.authority"
    ServicePCSSHConnAdmin string    = "service.pcssh.conn.admin"
    ServicePCSSHServerAuth string   = "service.pcssh.server.auth"

    ServicePCSSHConnProxy string    = "service.pcssh.conn.proxy"
    ServicePCSSHServerProxy string  = "service.pcssh.server.proxy"

)

// IsConnectionProblem returns whether this error is of ConnectionProblemError. This is originated from teleport trace
func isConnectionProblem(e error) bool {
    type ad interface {
        IsConnectionProblemError() bool
    }
    _, ok := e.(ad)
    return ok
}

// NewTeleport takes the daemon configuration, instantiates all required services
// and starts them under a supervisor, returning the supervisor object
func NewEmbeddedMasterProcess(sup pervice.ServiceSupervisor, cfg *service.PocketConfig) (*EmbeddedMasterProcess, error) {
    if err := sshcfg.ValidateMasterConfig(cfg); err != nil {
        return nil, errors.WithMessage(err, "Configuration error")
    }

    // if user started auth and another service (without providing the auth address for
    // that service, the address of the in-process auth will be used
    if cfg.Auth.Enabled && len(cfg.AuthServers) == 0 {
        cfg.AuthServers = []utils.NetAddr{cfg.Auth.SSHAddr}
    }

    // if there are no certificates, use self signed
    process := &EmbeddedMasterProcess{ServiceSupervisor: sup, config: cfg}

    // FIXME : (2017-06-06) this is a good place where CFSSL plugs in
    if cfg.Keygen == nil {
        cfg.Keygen = native.New()
    }
    if err := process.initProxy(); err != nil {
        return nil, err
    }
    if err := process.initAuthService(cfg.Keygen); err != nil {
        return nil, errors.WithStack(err)
    }
    return process, nil
}

// TeleportProcess structure holds the state of the Teleport daemon, controlling
// execution and configuration of the teleport services: ssh, auth and proxy.
type EmbeddedMasterProcess struct {
    pervice.ServiceSupervisor
    config    *service.PocketConfig
}

// connectToAuthService attempts to login into the auth servers specified in the
// configuration. Returns 'true' if successful
func (p *EmbeddedMasterProcess) connectToAuthService(role teleport.Role) (*service.Connector, error) {
    id := auth.IdentityID{HostUUID: p.config.HostUUID, Role: role}
    identity, err := auth.ReadIdentityFromCertStorage(p.config.CertStorage, id)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    authUser := identity.Cert.ValidPrincipals[0]
    authClient, err := auth.NewTunClient(
        string(role),
        p.config.AuthServers,
        authUser,
        []ssh.AuthMethod{ssh.PublicKeys(identity.KeySigner)},
    )
    // success?
    if err != nil {
        return nil, errors.WithStack(err)
    }
    // try calling a test method via auth api:
    //
    // ??? in case of failure it never gets back here!!!
    dn, err := authClient.GetDomainName()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    // success ? we're logged in!
    log.Debugf("[%s] %s connected to the cluster '%s'", role, authUser, dn)
    return &service.Connector{Client: authClient, Identity: identity}, nil
}

// initAuthStorage initializes the storage backend for the auth. service
func (p *EmbeddedMasterProcess) initAuthStorage() (backend.Backend, error) {
    var (
        cfg = &p.config.Auth
        bk backend.Backend = nil
        err error = nil
    )

    switch cfg.KeysBackend.Type {
        case teleport.SQLiteBackendType:
            // when backend is sqlite, we use db instance rather
            bk, err = sqlitebk.NewBackendFromDB(p.config.BackendDB)
        default:
            return nil, errors.Errorf("unsupported backend type: %v", cfg.KeysBackend.Type)
    }
    if err != nil {
        return nil, errors.WithStack(err)
    }

    return bk, nil
}

// initAuthService can be called to initialize auth server service
func (p *EmbeddedMasterProcess) initAuthService(authority auth.Authority) error {
    var (
        cfg = p.config
    )

    // updating the auth server presence
    var (
        asrvEventC = make(chan pervice.Event)
        authConnC  = make(chan pervice.Event)
    )
    p.RegisterServiceWithFuncs(
        ServicePCSSHServerAuth,
        func() error {
            var (
                authServer *auth.AuthServer = nil
                asrvOk, authConnOk bool = false, false
            )

            // prepare auth server
            srv := services.Server{
                ID:       p.config.HostUUID,
                Addr:     cfg.Auth.SSHAddr.Addr,
                Hostname: p.config.Hostname,
            }
            host, port, err := net.SplitHostPort(srv.Addr)
            // advertise-ip is explicitly set:
            if p.config.AdvertiseIP != nil {
                if err != nil {
                    return errors.WithStack(err)
                }
                srv.Addr = fmt.Sprintf("%v:%v", p.config.AdvertiseIP.String(), port)
            } else {
                // advertise-ip is not set, while the CA is listening on 0.0.0.0? lets try
                // to guess the 'advertise ip' then:
                if net.ParseIP(host).IsUnspecified() {
                    ip, err := utils.GuessHostIP()
                    if err != nil {
                        log.Debug(err)
                    } else {
                        srv.Addr = net.JoinHostPort(ip.String(), port)
                    }
                }
                log.Debugf("advertise_ip is not set for this auth server!!! Trying to guess the IP this server can be reached at: %v", srv.Addr)
            }

            // immediately register, and then keep repeating in a loop:
            for {
                select {
                    case <-p.StopChannel(): {
                        log.Debugf("[AUTH] heartbeat to other auth servers exited")
                        return nil
                    }
                    // waiting for authserver to come up
                    case ae := <-asrvEventC: {
                        authServer, asrvOk = ae.Payload.(*auth.AuthServer)
                        if asrvOk {
                            log.Debugf("[AUTH] AuthServer instance delivery succeed")
                        } else {
                            return errors.Errorf("[AUTH] AuthServer instance delivery failed")
                        }
                    }
                    // waiting for authority client connection to come up
                    case <- authConnC: {
                        authConnOk = true
                        log.Debugf("[AUTH] authServer client connection succeed")
                    }
                    default: {
                        if asrvOk && authConnOk {
                            err := authServer.UpsertAuthServer(srv, defaults.ServerHeartbeatTTL)
                            if err != nil {
                                log.Debugf("failed to announce presence: %v", err)
                            }
                            sleepTime := defaults.ServerHeartbeatTTL / 2 + utils.RandomDuration(defaults.ServerHeartbeatTTL / 10)
                            time.Sleep(sleepTime)
                        }
                    }
                }
            }
        },
        pervice.BindEventWithService(eventAuthorityInstance, asrvEventC),
        pervice.BindEventWithService(eventAuthorityClientConn, authConnC))

    // Register an SSH endpoint which is used to create an SSH tunnel to send HTTP requests to the Auth API
    p.RegisterServiceWithFuncs(
        ServicePCSSHAuthority,
        func() error {
            log.Debugf("[AUTH] Auth service is starting on %v", cfg.Auth.SSHAddr.Addr)

            var (
                authTunnel *auth.AuthTunnel
                err error = nil
            )

            // Initialize the storage back-ends for keys, events and records
            bkEnd, err := p.initAuthStorage()
            if err != nil {
                return errors.WithStack(err)
            }

            // TODO : Disable audit for release
            // create the audit log, which will be consuming (and recording) all events
            // and record sessions
            var auditLog events.IAuditLog
            if cfg.Auth.NoAudit {
                auditLog = &events.DiscardAuditLog{}
                log.Debugf("the audit and session recording are turned off")
            } else {
                auditLog, err = events.NewAuditLog(filepath.Join(cfg.DataDir, "log"))
                if err != nil {
                    return errors.WithStack(err)
                }
            }

            // first, create the AuthServer
            authServer, identity, err := auth.PocketAuthInit(
                auth.InitConfig{
                    Backend:         bkEnd,
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
                },
                p.config.CertStorage,
                cfg.SeedConfig)
            if err != nil {
                return errors.WithStack(err)
            }
            p.BroadcastEvent(pervice.Event{Name:eventAuthorityInstance, Payload:authServer})

            // second, create the API Server: it's actually a collection of API servers,
            // each serving requests for a "role" which is assigned to every connected
            // client based on their certificate (user, server, admin, etc)
            sessionService, err := session.New(bkEnd)
            if err != nil {
                return errors.WithStack(err)
            }
            apiConf := &auth.APIConfig{
                AuthServer:        authServer,
                SessionService:    sessionService,
                PermissionChecker: auth.NewStandardPermissions(),
                AuditLog:          auditLog,
                CertSigner:        p.config.CaSigner,
                CertStorage:       p.config.CertStorage,
            }

            limiter, err := limiter.NewLimiter(cfg.Auth.Limiter)
            if err != nil {
                return errors.WithStack(err)
            }
            authTunnel, err = auth.NewTunnel(
                cfg.Auth.SSHAddr,
                identity.KeySigner,
                apiConf,
                auth.SetLimiter(limiter),
            )
            if err != nil {
                log.Debugf("[AUTH] Error: %v", err)
                return errors.WithStack(err)
            }
            // broadcast authTunnel to close it later
            err = authTunnel.Start()
            if err != nil {
                log.Debugf("[AUTH] Auth Tunnel start error: %v", err)
            }
            <- p.StopChannel()
            authTunnel.Close()
            log.Debugf("[AUTH] Auth Tunnel exited %v", err)
            return errors.WithStack(err)
        })

    // Heart beat auth server presence, this is not the best place for this logic. Consolidate it into auth package later
    p.RegisterServiceWithFuncs(
        ServicePCSSHConnAdmin,
        func() error {
            var (
                retryTime = defaults.ServerHeartbeatTTL / 3
                role = teleport.RoleAdmin
            )
            for {
                connector, err := p.connectToAuthService(role)
                if err == nil {
                    log.Debugf("[%v] connected successfully.", role)
                    p.BroadcastEvent(pervice.Event{Name:eventAuthorityClientConn})

                    // wait for service closure
                    <- p.StopChannel()
                    err = connector.Client.Close()
                    log.Debugf("[%v] connection closed. Error : %v", role, err)
                    return errors.WithStack(err)
                }
                if err != nil {
                    log.Debugf("[%v] failed to connect to auth server. Error : %v", role, err)
                    time.Sleep(retryTime)
                    continue
                }
            }
        })

    return nil
}

// registerWithAuthServer uses one time provisioning token obtained earlier from the server to get a pair of SSH keys
// signed by Auth server host certificate authority
func (p *EmbeddedMasterProcess) registerWithAuthServer(token string, role teleport.Role, eventName string) {
    var (
        cfg        = p.config
        identityID = auth.IdentityID{Role: role, HostUUID: cfg.HostUUID}
        eventC     = make(chan pervice.Event)
    )
    // this means the server has not been initialized yet, we are starting the registering client that attempts to
    // connect to the auth server and provision the keys
    p.RegisterServiceWithFuncs(
        ServicePCSSHConnProxy,
        func() error {
            var (
                retryTime = defaults.ServerHeartbeatTTL / 3
                authServer *auth.AuthServer = nil
            )
            // wait for AuthServer to come up
            ae := <- eventC
            asrv, ok := ae.Payload.(*auth.AuthServer)
            if ok {
                authServer = asrv
                log.Debugf("[%v] AuthServer instance delivery succeed", role)
            } else {
                return errors.Errorf("[%v] AuthServer instance delivery failed", role)
            }

            for {
                connector, err := p.connectToAuthService(role)
                if err == nil {
                    log.Debugf("[%v] connected successfully. Broadcast the connection with %s", role, eventName)
                    p.BroadcastEvent(pervice.Event{Name: eventName, Payload: connector})

                    // wait for service closure
                    <- p.StopChannel()
                    err = connector.Client.Close()
                    log.Debugf("[%v] connection closed. Error : %v", role.String(), err)
                    return errors.WithStack(err)
                }
                if isConnectionProblem(err) {
                    log.Debugf("[%v] connecting from %v to auth server: %v", role, cfg.HostUUID, err)
                    time.Sleep(retryTime)
                    continue
                }
                // Auth service is on the same host, no need to go though the invitation procedure
                // This is the place where proxy certificate is generated and stored
                log.Debugf("[%s] this server has local Auth server started, using it to add role to the cluster", role.String())
                err = auth.LocalRegisterWithCertStorage(authServer, cfg.CertStorage, identityID)
                if err != nil {
                    log.Debugf("[%v] failed to join the cluster: %v", role, err)
                    time.Sleep(retryTime)
                } else {
                    log.Debugf("[%v] Successfully registered with the cluster", role)
                    continue
                }
            }
            return nil
        },
        pervice.BindEventWithService(eventAuthorityInstance, eventC))
}

// initProxy gets called if teleport runs with 'proxy' role enabled to proxy SSH connections to nodes running with
// 'node' role
func (p *EmbeddedMasterProcess) initProxy() error {
    var (
        cfg     = p.config
        eventsC = make(chan pervice.Event)
    )
    p.RegisterServiceWithFuncs(
        ServicePCSSHServerProxy,
        func() error {
            // wait for client connection
            event := <-eventsC
            log.Debugf("[SSH] received %v", &event)
            conn, ok := (event.Payload).(*service.Connector)
            if !ok {
                return errors.Errorf("unsupported connector type: %T", event.Payload)
            }

            // setup reverse tunnel
            reverseTunnelLimiter, err := limiter.NewLimiter(cfg.Proxy.Limiter)
            if err != nil {
                return errors.WithStack(err)
            }
            tsrv, err := reversetunnel.NewServer(
                cfg.Proxy.ReverseTunnelListenAddr,
                []ssh.Signer{conn.Identity.KeySigner},
                conn.Client,
                reversetunnel.SetLimiter(reverseTunnelLimiter),
                reversetunnel.DirectSite(conn.Identity.Cert.Extensions[utils.CertExtensionAuthority], conn.Client),
            )
            if err != nil {
                return errors.WithStack(err)
            }

            // setup ssh proxy server
            proxyLimiter, err := limiter.NewLimiter(cfg.Proxy.Limiter)
            if err != nil {
                return errors.WithStack(err)
            }
            SSHProxy, err := srv.NewPocketSSHServer(cfg.Proxy.SSHAddr,
                cfg.Hostname,
                cfg.HostUUID,
                []ssh.Signer{conn.Identity.KeySigner},
                conn.Client,
                // TODO : we need to supply proxy ip
                cfg.AdvertiseIP,
                srv.SetLimiter(proxyLimiter),
                srv.SetProxyMode(tsrv),
                srv.SetSessionServer(conn.Client),
                srv.SetAuditLog(conn.Client),
            )
            if err != nil {
                return errors.WithStack(err)
            }
            err = SSHProxy.Start()
            log.Debugf("[PROXY] SSH proxy service is starting on %v. Error : %v", cfg.Proxy.SSHAddr.Addr, err)

            // wait for exit
            <- p.StopChannel()
            err = SSHProxy.Close()
            log.Debugf("[PROXY] SSH proxy exited. Error: %v", err)
            err = tsrv.Close()
            log.Debugf("[PROXY] ReverseTunnel exited. Error: %v", err)
            return errors.WithStack(err)
        },
        pervice.BindEventWithService(eventProxyIdentity, eventsC))

    p.registerWithAuthServer(p.config.Token, teleport.RoleProxy, eventProxyIdentity)

    return nil
}
