package sshproc

import (
    "time"

    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/auth"
    "github.com/gravitational/teleport/lib/defaults"
    "github.com/gravitational/teleport/lib/limiter"
    "github.com/gravitational/teleport/lib/srv"
    "github.com/gravitational/teleport/lib/service"
    pervice "github.com/stkim1/pc-node-agent/service"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "golang.org/x/crypto/ssh"
)

const (
    // EventNodeSSHServiceStop is to be generated when node's service should stop
    EventNodeSSHServiceStop string = "event.node.ssh.service.stop"
    // eventNodeSSHIdentity is generated when node's identity has been received
    eventNodeSSHIdentity string    = "event.node.ssh.identity"
)

// IsConnectionProblem returns whether this error is of ConnectionProblemError. This is originated from teleport trace
func isConnectionProblem(e error) bool {
    type ad interface {
        IsConnectionProblemError() bool
    }
    _, ok := e.(ad)
    return ok
}

// NewTeleport takes the daemon configuration, instantiates all required services, but does not run service
func NewEmbeddedNodeProcess(sup pervice.AppSupervisor, cfg *service.PocketConfig) (*EmbeddedNodeProcess, error) {
    return &EmbeddedNodeProcess{
        AppSupervisor:    sup,
        config:           cfg,
    }, nil
}

// TeleportProcess structure holds the state of the Teleport daemon, controlling
// execution and configuration of the teleport services: ssh, auth and proxy.
type EmbeddedNodeProcess struct {
    pervice.AppSupervisor
    config *service.PocketConfig
}

func (p *EmbeddedNodeProcess) Close() error {
    p.BroadcastEvent(pervice.Event{Name:EventNodeSSHServiceStop})
    return nil
}

// connectToAuthService attempts to login into the auth servers specified in the
// configuration. Returns 'true' if successful
func (p *EmbeddedNodeProcess) connectToAuthService(role teleport.Role) (*service.Connector, error) {
    var (
        cfg = p.config
        id = auth.IdentityID{HostUUID: p.config.HostUUID, Role: role}
    )

    identity, err := auth.NodeReadIdentityFromFile(cfg.NodeSSHPrivateKeyFile, cfg.NodeSSHCertificateFile, id)
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
    log.Debugf("[Node] %s connected to the cluster '%s'", authUser, dn)
    return &service.Connector{Client: authClient, Identity: identity}, nil
}

// RegisterWithAuthServer uses one time provisioning token obtained earlier
// from the server to get a pair of SSH keys signed by Auth server host
// certificate authority
func (p *EmbeddedNodeProcess) registerWithAuthServer(token string, role teleport.Role, eventName string) {
    var (
        cfg              = p.config
        stopEventC       = make(chan pervice.Event)
        identityID       = auth.IdentityID {
            Role:        role,
            HostUUID:    cfg.HostUUID,
            NodeName:    cfg.Hostname,
        }
    )

    // this means the server has not been initialized yet, we are starting
    // the registering client that attempts to connect to the auth server
    // and provision the keys
    p.RegisterServiceWithFuncs(
        func() error {
            var (
                retryTime = defaults.ServerHeartbeatTTL / 3
            )
            for {
                connector, err := p.connectToAuthService(role)
                if err == nil {
                    p.BroadcastEvent(pervice.Event{Name:eventName, Payload:connector})
                    select {
                        case <- p.StopChannel(): {
                            err = connector.Client.Close()
                            log.Debugf("[%v] connection closed. Error : %v ", role, err)
                        }
                        case <- stopEventC: {
                            err = connector.Client.Close()
                            log.Debugf("[%v] connection closed. Error : %v ", role, err)
                        }
                    }
                    return errors.WithStack(err)
                }
                if isConnectionProblem(err) {
                    log.Debugf("[%v] connecting from %v to auth server: %v ", role, cfg.HostUUID ,err)
                    time.Sleep(retryTime)
                    continue
                }
                // We haven't connected yet, so we expect the token to exist
                // TODO when it's necessary to bring local connectivity on OSX, we'll do following
                // 1) bring in LocalAuth connectivity or
                // 2) combine PocketCoreTeleportProcess & PocketNodeTeleportProcess together

                // Auth server is remote, so we need a provisioning token.
                // !!! Since we need token only at the initialization, empty token check need to be here.!!!
                if token == "" {
                    return errors.Errorf("[%v] must join a cluster and needs a provisioning token", role)
                }
                log.Debugf("[%v] joining the cluster with a token %v", role, token)
                err = auth.NodeRegister(cfg.NodeSSHPrivateKeyFile, cfg.NodeSSHCertificateFile, token, identityID, cfg.AuthServers)
                if err != nil {
                    log.Debugf("[%v] failed to join the cluster: %v", role, err)
                    time.Sleep(retryTime)
                } else {
                    log.Debugf("[%v] successfully registered to the cluster", role)
                    continue
                }
            }
        },
        func(_ func(interface{})) error {
            return nil
        },
        pervice.BindEventWithService(EventNodeSSHServiceStop, stopEventC))
}

func (p *EmbeddedNodeProcess) StartNodeSSH() error {
    var (
        cfg           = p.config
        eventsC       = make(chan pervice.Event)
        stopEventC    = make(chan pervice.Event)
    )

    p.RegisterServiceWithFuncs(
        func() error {
            var (
                sshServer *srv.Server = nil
            )
            event := <-eventsC
            log.Debugf("[SSH] received %v", &event)
            conn, ok := (event.Payload).(*service.Connector)
            if !ok {
                return errors.Errorf("unsupported connector type: %T", event.Payload)
            }

            limiter, err := limiter.NewLimiter(cfg.SSH.Limiter)
            if err != nil {
                return errors.WithStack(err)
            }

            sshServer, err = srv.NewPocketSSHServer(cfg.SSH.Addr,
                cfg.Hostname,
                cfg.HostUUID,
                []ssh.Signer{conn.Identity.KeySigner},
                conn.Client,
                cfg.AdvertiseIP,
                srv.SetLimiter(limiter),
                srv.SetShell(cfg.SSH.Shell),
                srv.SetAuditLog(conn.Client),
                srv.SetSessionServer(conn.Client),
                srv.SetLabels(cfg.SSH.Labels, cfg.SSH.CmdLabels),
            )
            if err != nil {
                return errors.WithStack(err)
            }

            err = sshServer.Start()
            if err != nil {
                return errors.WithStack(err)
            }
            log.Debugf("[SSH] Service is starting on %v. Error : %v", cfg.SSH.Addr.Addr, err)

            select {
                case <- p.StopChannel(): {
                    err = sshServer.Close()
                    log.Debugf("[SSH] node service exited. Error : %v", err)
                }
                case <- stopEventC: {
                    err = sshServer.Close()
                    log.Debugf("[SSH] node service exited. Error : %v", err)
                }
            }
            return errors.WithStack(err)
        },
        func(_ func(interface{})) error {
            return nil
        },
        pervice.BindEventWithService(eventNodeSSHIdentity, eventsC),
        pervice.BindEventWithService(EventNodeSSHServiceStop, stopEventC))

    // register & generate a signed ssh pub/prv key set. This should be executed after receiver funcs are registered
    p.registerWithAuthServer(p.config.Token, teleport.RoleNode, eventNodeSSHIdentity)

    return nil
}

/// --- DOCKER ENGINE CERTIFICATE ACQUISITION --- ///
// AcquireEngineCertificate uses one time provisioning token obtained earlier from the server to get a pair of Docker
// Engine keys signed by Auth server host certificate authority
func (p *EmbeddedNodeProcess) AcquireEngineCertificate(onSucessAct func(certPack *auth.PocketResponseAuthKeyCert) error) error {
    var (
        cfg     = p.config
        eventsC = make(chan pervice.Event)
        role    = teleport.RoleNode
        token   = p.config.Token
    )
    // Auth server is remote, so we need a provisioning token
    if token == "" {
        return errors.Errorf("%v must request a signed certificate and needs a provisioning token", role)
    }

    // this means the server has not been initialized yet, we are starting the registering client that attempts to
    // connect to the auth server and provision the keys
    p.RegisterServiceWithFuncs(
        func() error {
            var retryTime = defaults.ServerHeartbeatTTL / 3

            // we're to wait until SSH successfully connects to master
            _ = <-eventsC
            log.Debugf("[Node] %v requesting a signed certificate with a token %v | UUID %v : ", role, token, cfg.HostUUID)
            // start request signed certificate
            for {
                keys, err := auth.RequestSignedCertificate(
                    &auth.PocketRequestBase {
                        AuthServers:    cfg.AuthServers,
                        Role:           role,
                        Hostname:       cfg.Hostname,
                        HostUUID:       cfg.HostUUID,
                        AuthToken:      token,
                    })
                if err != nil {
                    log.Debugf("[%v] failed to receive a signed certificate : %v", role, err)
                    time.Sleep(retryTime)
                } else {
                    if onSucessAct != nil {
                        err = onSucessAct(keys)
                        if err != nil {
                            return err
                        }
                    }
                    log.Debugf("[%v] Successfully received a signed certificate & finished subsequent action", role)
                    return nil
                }
            }
        },
        func(_ func(interface{})) error {
            return nil
        },
        // we're to wait until SSH successfully connects to master
        pervice.BindEventWithService(eventNodeSSHIdentity, eventsC),
    )

    return nil
}

/// --- NODE OPERATION PARAM ACQUISITION --- ///
// AcquireEngineCertificate uses one time provisioning token obtained earlier from the server to get a pair of Docker
// Engine keys signed by Auth server host certificate authority
func (p *EmbeddedNodeProcess) AcquireUserIdentity(onSucessAct func(user *auth.PocketResponseUserIdentity) error) error {
    var (
        cfg     = p.config
        eventsC = make(chan pervice.Event)
        role    = teleport.RoleNode
        token   = p.config.Token
    )
    // Auth server is remote, so we need a provisioning token
    if token == "" {
        return errors.Errorf("%v must request a signed certificate and needs a provisioning token", role)
    }

    // this means the server has not been initialized yet, we are starting the registering client that attempts to
    // connect to the auth server and provision the keys
    p.RegisterServiceWithFuncs(
        func() error {
            var retryTime = defaults.ServerHeartbeatTTL / 3

            // we're to wait until SSH successfully connects to master
            _ = <-eventsC
            log.Debugf("[Node] %v requesting a signed certificate with a token %v | UUID %v : ", role, token, cfg.HostUUID)
            // start request signed certificate
            for {
                user, err := auth.RequestUserIdentity(
                    &auth.PocketRequestBase {
                        AuthServers:     cfg.AuthServers,
                        Role:            role,
                        Hostname:        cfg.Hostname,
                        HostUUID:        cfg.HostUUID,
                        AuthToken:       token,
                    })
                if err != nil {
                    log.Debugf("[%v] failed to receive a signed certificate : %v", role, err)
                    time.Sleep(retryTime)
                } else {
                    if onSucessAct != nil {
                        err = onSucessAct(user)
                        if err != nil {
                            return err
                        }
                    }
                    log.Debugf("[%v] Successfully received a signed certificate & finished subsequent action", role)
                    return nil
                }
            }
        },
        func(_ func(interface{})) error {
            return nil
        },
        // we're to wait until SSH successfully connects to master
        pervice.BindEventWithService(eventNodeSSHIdentity, eventsC),
    )

    return nil
}