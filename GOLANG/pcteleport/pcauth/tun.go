/*
Copyright 2015 Gravitational, Inc.

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

package pcauth

import (
    "encoding/json"
    "fmt"
    "io"
    "net"
    "net/http"
    "os"

    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/limiter"
    "github.com/gravitational/teleport/lib/services"
    "github.com/gravitational/teleport/lib/sshutils"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/auth"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/agent"
)

// AuthTunnel listens on TCP/IP socket and accepts SSH connections. It then establishes
// an SSH tunnell which HTTP requests travel over. In other words, the Auth Service API
// runs on HTTP-via-SSH-tunnel.
//
// Use auth.TunClient to connect to AuthTunnel
type AuthTunnel struct {
    // authServer implements the "beef" of the Auth service
    authServer         *auth.AuthServer
    config             *auth.APIConfig

    // sshServer implements the nuts & bolts of serving an SSH connection
    // to create a tunnel
    sshServer          *sshutils.Server
    hostSigner         ssh.Signer
    hostCertChecker    ssh.CertChecker
    userCertChecker    ssh.CertChecker
    limiter            *limiter.Limiter
}

// ServerOption is the functional argument passed to the server
type ServerOption func(s *AuthTunnel) error

// SetLimiter sets rate and connection limiter for auth tunnel server
func SetLimiter(limiter *limiter.Limiter) ServerOption {
    return func(s *AuthTunnel) error {
        s.limiter = limiter
        return nil
    }
}

// NewTunnel creates a new SSH tunnel server which is not started yet.
// This is how "site API" (aka "auth API") is served: by creating
// an "tunnel server" which serves HTTP via SSH.
func NewTunnel(addr utils.NetAddr,
    hostSigner ssh.Signer,
    apiConf *auth.APIConfig,
    opts ...ServerOption) (tunnel *AuthTunnel, err error) {

    tunnel = &AuthTunnel{
        authServer: apiConf.AuthServer,
        config:     apiConf,
    }
    tunnel.limiter, err = limiter.NewLimiter(limiter.LimiterConfig{})
    if err != nil {
        return nil, trace.Wrap(err)
    }
    // apply functional options:
    for _, o := range opts {
        if err := o(tunnel); err != nil {
            return nil, err
        }
    }
    // create an SSH server and assign the tunnel to be it's "new SSH channel handler"
    tunnel.sshServer, err = sshutils.NewServer(
        teleport.ComponentAuth,
        addr,
        tunnel,
        []ssh.Signer{hostSigner},
        sshutils.AuthMethods{
            Password:  tunnel.passwordAuth,
            PublicKey: tunnel.keyAuth,
        },
        sshutils.SetLimiter(tunnel.limiter),
    )
    if err != nil {
        return nil, err
    }
    tunnel.userCertChecker = ssh.CertChecker{IsAuthority: tunnel.isUserAuthority}
    tunnel.hostCertChecker = ssh.CertChecker{IsAuthority: tunnel.isHostAuthority}
    return tunnel, nil
}

func (s *AuthTunnel) Addr() string {
    return s.sshServer.Addr()
}

func (s *AuthTunnel) Start() error {
    return s.sshServer.Start()
}

func (s *AuthTunnel) Close() error {
    if s != nil && s.sshServer != nil {
        return s.sshServer.Close()
    }
    return nil
}

// HandleNewChan implements NewChanHandler interface: it gets called every time a new SSH
// connection is established
func (s *AuthTunnel) HandleNewChan(_ net.Conn, sconn *ssh.ServerConn, nch ssh.NewChannel) {
    log.Debugf("[AUTH] new channel request for %v from %v", nch.ChannelType(), sconn.RemoteAddr())
    cht := nch.ChannelType()
    switch cht {

    // New connection to the Auth API via SSH:
    case auth.ReqDirectTCPIP:
        if !s.haveExt(sconn, auth.ExtHost, auth.ExtWebSession, auth.ExtWebPassword) {
            nch.Reject(
                ssh.UnknownChannelType,
                fmt.Sprintf("register clients can not TCPIP: %v", cht))
            return
        }
        req, err := sshutils.ParseDirectTCPIPReq(nch.ExtraData())
        if err != nil {
            log.Errorf("[AUTH] failed to parse request data: %v, err: %v",
                string(nch.ExtraData()), err)
            nch.Reject(ssh.UnknownChannelType,
                "failed to parse direct-tcpip request")
            return
        }
        sshCh, _, err := nch.Accept()
        if err != nil {
            log.Infof("[AUTH] could not accept channel (%s)", err)
            return
        }
        go s.onAPIConnection(sconn, sshCh, req)

    case auth.ReqWebSessionAgent:
        // this is a protective measure, so web requests can be only done
        // if have session ready
        if !s.haveExt(sconn, auth.ExtWebSession) {
            nch.Reject(
                ssh.UnknownChannelType,
                fmt.Sprintf("don't have web session for: %v", cht))
            return
        }
        ch, _, err := nch.Accept()
        if err != nil {
            log.Infof("[AUTH] could not accept channel (%s)", err)
            return
        }
        go s.handleWebAgentRequest(sconn, ch)

    default:
        nch.Reject(ssh.UnknownChannelType, fmt.Sprintf(
            "unknown channel type: %v", cht))
    }
}

// isHostAuthority is called during checking the client key, to see if the signing
// key is the real host CA authority key.
func (s *AuthTunnel) isHostAuthority(auth ssh.PublicKey) bool {
    key, err := s.authServer.GetCertAuthority(services.CertAuthID{DomainName: s.authServer.DomainName, Type: services.HostCA}, false)
    if err != nil {
        log.Errorf("failed to retrieve user authority key, err: %v", err)
        return false
    }
    checkers, err := key.Checkers()
    if err != nil {
        log.Errorf("failed to parse CA keys: %v", err)
        return false
    }
    for _, checker := range checkers {
        if sshutils.KeysEqual(checker, auth) {
            return true
        }
    }
    return false
}

// isUserAuthority is called during checking the client key, to see if the signing
// key is the real user CA authority key.
func (s *AuthTunnel) isUserAuthority(auth ssh.PublicKey) bool {
    keys, err := s.getTrustedCAKeys(services.UserCA)
    if err != nil {
        log.Errorf("failed to retrieve trusted keys, err: %v", err)
        return false
    }
    for _, k := range keys {
        if sshutils.KeysEqual(k, auth) {
            return true
        }
    }
    return false
}

func (s *AuthTunnel) getTrustedCAKeys(CertType services.CertAuthType) ([]ssh.PublicKey, error) {
    cas, err := s.authServer.GetCertAuthorities(CertType, false)
    if err != nil {
        return nil, err
    }
    out := []ssh.PublicKey{}
    for _, ca := range cas {
        checkers, err := ca.Checkers()
        if err != nil {
            return nil, trace.Wrap(err)
        }
        out = append(out, checkers...)
    }
    return out, nil
}

func (s *AuthTunnel) haveExt(sconn *ssh.ServerConn, ext ...string) bool {
    if sconn.Permissions == nil {
        return false
    }
    for _, e := range ext {
        if sconn.Permissions.Extensions[e] != "" {
            return true
        }
    }
    return true
}

func (s *AuthTunnel) handleWebAgentRequest(sconn *ssh.ServerConn, ch ssh.Channel) {
    defer ch.Close()

    if sconn.Permissions.Extensions[auth.ExtRole] != string(teleport.RoleWeb) {
        log.Errorf("role %v doesn't have permission to request agent",
            sconn.Permissions.Extensions[auth.ExtRole])
        return
    }

    ws, err := s.authServer.GetWebSession(sconn.User(), sconn.Permissions.Extensions[auth.ExtWebSession])
    if err != nil {
        log.Errorf("session error: %v", err)
        return
    }

    priv, err := ssh.ParseRawPrivateKey(ws.WS.Priv)
    if err != nil {
        log.Errorf("session error: %v", err)
        return
    }

    pub, _, _, _, err := ssh.ParseAuthorizedKey(ws.WS.Pub)
    if err != nil {
        log.Errorf("session error: %v", err)
        return
    }

    cert, ok := pub.(*ssh.Certificate)
    if !ok {
        log.Errorf("session error, not a cert: %T", pub)
        return
    }
    addedKey := agent.AddedKey{
        PrivateKey:       priv,
        Certificate:      cert,
        Comment:          "web-session@teleport",
        LifetimeSecs:     0,
        ConfirmBeforeUse: false,
    }
    newKeyAgent := agent.NewKeyring()
    if err := newKeyAgent.Add(addedKey); err != nil {
        log.Errorf("failed to add: %v", err)
        return
    }
    if err := agent.ServeAgent(newKeyAgent, ch); err != nil && err != io.EOF {
        log.Errorf("Serve agent err: %v", err)
    }
}

// onAPIConnection accepts an incoming SSH connection via TCP/IP and forwards
// it to the local auth server which listens on local UNIX pipe
//
func (s *AuthTunnel) onAPIConnection(sconn *ssh.ServerConn, sshChan ssh.Channel, req *sshutils.DirectTCPIPReq) {
    defer sconn.Close()

    // retreive the role from thsi connection's permissions (make sure it's a valid role)
    role := teleport.Role(sconn.Permissions.Extensions[auth.ExtRole])
    if err := role.Check(); err != nil {
        log.Errorf(err.Error())
        return
    }

    api := auth.NewAPIServer(s.config, role)
    // Since PocketCluster API is an addition to existing api, we'll handle normal request in NotFound functions
    pcapi := NewPocketAPIServer(s.config, role, func(w http.ResponseWriter, r *http.Request){
        // TODO : log
        log.Infof("[AUTH] PocketCluster API does not exists %v", r.RequestURI)
        api.ServeHTTP(w, r)
    })

    socket := fakeSocket{
        closed:      make(chan int),
        connections: make(chan net.Conn),
    }

    go func() {
        connection := &fakeSSHConnection{
            remoteAddr: sconn.RemoteAddr(),
            sshChan:    sshChan,
            closed:     make(chan int),
        }
        // fakesocket.Accept() will pick it up:
        socket.connections <- connection

        // wait for the connection wrapper to close, so we'll close
        // the fake socket, causing http.Serve() below to stop
        <-connection.closed
        socket.Close()
    }()

    // serve HTTP API via this SSH connection until it gets closed:
    http.Serve(&socket, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // take SSH client name and pass it to HTTP API via HTTP Auth
        r.SetBasicAuth(sconn.User(), "")
        pcapi.ServeHTTP(w, r)
    }))
}

func (s *AuthTunnel) keyAuth(
    conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {

    log.Infof("[AUTH] keyAuth: %v->%v, user=%v", conn.RemoteAddr(), conn.LocalAddr(), conn.User())
    cert, ok := key.(*ssh.Certificate)
    if !ok {
        return nil, trace.Errorf("ERROR: Server doesn't support provided key type")
    }

    if cert.CertType == ssh.HostCert {
        err := s.hostCertChecker.CheckHostKey(conn.User(), conn.RemoteAddr(), key)
        if err != nil {
            log.Warningf("conn(%v->%v, user=%v) ERROR: failed auth user %v, err: %v",
                conn.RemoteAddr(), conn.LocalAddr(), conn.User(), conn.User(), err)
            return nil, err
        }
        if err := s.hostCertChecker.CheckCert(conn.User(), cert); err != nil {
            log.Warningf("conn(%v->%v, user=%v) ERROR: Failed to authorize user %v, err: %v",
                conn.RemoteAddr(), conn.LocalAddr(), conn.User(), conn.User(), err)
            return nil, trace.Wrap(err)
        }
        perms := &ssh.Permissions{
            Extensions: map[string]string{
                auth.ExtHost: conn.User(),
                auth.ExtRole: cert.Permissions.Extensions[utils.CertExtensionRole],
            },
        }
        return perms, nil
    }
    // we are assuming that this is a user cert
    if err := s.userCertChecker.CheckCert(conn.User(), cert); err != nil {
        log.Warningf("conn(%v->%v, user=%v) ERROR: Failed to authorize user %v, err: %v",
            conn.RemoteAddr(), conn.LocalAddr(), conn.User(), conn.User(), err)
        return nil, trace.Wrap(err)
    }
    // we are not using cert extensions for User certificates because of OpenSSH bug
    // https://bugzilla.mindrot.org/show_bug.cgi?id=2387
    perms := &ssh.Permissions{
        Extensions: map[string]string{
            auth.ExtHost: conn.User(),
            auth.ExtRole: string(teleport.RoleUser),
        },
    }
    return perms, nil
}

func (s *AuthTunnel) passwordAuth(
    conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
    var ab *authBucket
    if err := json.Unmarshal(password, &ab); err != nil {
        return nil, err
    }

    log.Infof("[AUTH] login attempt: user '%v' type '%v'", conn.User(), ab.Type)

    switch ab.Type {
    case auth.AuthWebPassword:
        if err := s.authServer.CheckPassword(conn.User(), ab.Pass, ab.HotpToken); err != nil {
            log.Warningf("password auth error: %#v", err)
            return nil, trace.Wrap(err)
        }
        perms := &ssh.Permissions{
            Extensions: map[string]string{
                auth.ExtWebPassword: "<password>",
                auth.ExtRole:        string(teleport.RoleUser),
            },
        }
        log.Infof("[AUTH] password authenticated user: '%v'", conn.User())
        return perms, nil
    case auth.AuthWebSession:
        // we use extra permissions mechanism to keep the connection data
        // after authorization, in this case the session
        perms := &ssh.Permissions{
            Extensions: map[string]string{
                auth.ExtWebSession: string(ab.Pass),
                auth.ExtRole:       string(teleport.RoleWeb),
            },
        }
        if _, err := s.authServer.GetWebSession(conn.User(), string(ab.Pass)); err != nil {
            return nil, trace.Errorf("session resume error: %v", trace.Wrap(err))
        }
        log.Infof("[AUTH] session authenticated user: '%v'", conn.User())
        return perms, nil
    // when a new server tries to use the auth API to register in the cluster,
    // it will use the token as a passowrd (happens only once during registration):
    case auth.AuthToken:
        _, err := s.authServer.ValidateToken(string(ab.Pass))
        if err != nil {
            log.Errorf("token validation error: %v", err)
            return nil, trace.Wrap(err, fmt.Sprintf("invalid token for: %v", ab.User))
        }
        perms := &ssh.Permissions{
            Extensions: map[string]string{
                auth.ExtToken: string(password),
                auth.ExtRole:  string(teleport.RoleProvisionToken),
            }}
        utils.Consolef(os.Stdout, "[AUTH] Successfully accepted token for %v", conn.User())
        return perms, nil
    case auth.AuthSignupToken:
        _, err := s.authServer.GetSignupToken(string(ab.Pass))
        if err != nil {
            return nil, trace.Errorf("token validation error: %v", trace.Wrap(err))
        }
        perms := &ssh.Permissions{
            Extensions: map[string]string{
                auth.ExtToken: string(password),
                auth.ExtRole:  string(teleport.RoleSignup),
            }}
        log.Infof("[AUTH] session authenticated prov. token: '%v'", conn.User())
        return perms, nil
    default:
        return nil, trace.Errorf("unsupported auth method: '%v'", ab.Type)
    }
}

// authBucket uses password to transport app-specific user name and
// auth-type in addition to the password to support auth
type authBucket struct {
    User      string `json:"user"`
    Type      string `json:"type"`
    Pass      []byte `json:"pass"`
    HotpToken string `json:"hotpToken"`
}
