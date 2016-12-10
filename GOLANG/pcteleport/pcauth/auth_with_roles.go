package pcauth

import (
    "fmt"
    "os"

    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/auth"
    "github.com/gravitational/teleport/lib/session"
    "github.com/gravitational/teleport/lib/events"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/stkim1/pcrypto"
)

type authWithRoles struct {
    authServer  *auth.AuthServer
    permChecker auth.PermissionChecker
    sessions    session.Service
    caSigner    *pcrypto.CaSigner
    role        teleport.Role
    alog        events.IAuditLog
}

func (a *authWithRoles) issueSignedCertificateWithToken(req *signedCertificateReq) (*packedAuthKeyCert, error) {
    // TODO : add action perm for requesting signed certificate
    if err := a.permChecker.HasPermission(a.role, auth.ActionRegisterUsingToken); err != nil {
        return nil, trace.Wrap(err)
    }
    return issueSignedCertificateWithToken(a, req)
}

// issueSignedCertificateWithToken adds a new signed certificate for a node to the PocketCluster using previously issued token.
// A node must also request a specific role (and the role must match one of the roles the token was generated for).
//
// If a token was generated with a TTL, it gets enforced (can't register new nodes after TTL expires)
// If a token was generated with a TTL=0, it means it's a single-use token and it gets destroyed
// after a successful registration.
func issueSignedCertificateWithToken(a *authWithRoles, req *signedCertificateReq) (*packedAuthKeyCert, error) {
    if len(req.Hostname) == 0 {
        return nil, trace.BadParameter("Hostname cannot be empty")
    }
    if len(req.HostID) == 0 {
        return nil, trace.BadParameter("HostID cannot be empty")
    }
    log.Infof("[AUTH] Node `%v`[%v] requests a signed certificate", req.Hostname, req.HostID)
    if err := req.Role.Check(); err != nil {
        return nil, trace.Wrap(err)
    }
    // make sure the token is valid:
    roles, err := a.authServer.ValidateToken(req.Token)
    if err != nil {
        msg := fmt.Sprintf("`%v` cannot receive a signed certificate with %s. Token error: %v", req.Hostname, req.Role, err)
        log.Warnf("[AUTH] %s", msg)
        return nil, trace.AccessDenied(msg)
    }
    // make sure the caller is requested wthe role allowed by the token:
    if !roles.Include(req.Role) {
        msg := fmt.Sprintf("'%v' cannot receive a signed certificate, the token does not allow '%s' role", req.Hostname, req.Role)
        log.Warningf("[AUTH] %s", msg)
        return nil, trace.BadParameter(msg)
    }
    if !checkTokenTTL(a.authServer, req.Token) {
        return nil, trace.AccessDenied("'%v' cannot cannot receive a signed certificate. The token has expired", req.Hostname)
    }
    // generate & return the node cert:
    keys, err := createSignedCertificate(a.caSigner, req)
    if err != nil {
        return nil, trace.Wrap(err)
    }
    utils.Consolef(os.Stdout, "[AUTH] A signed Certificate for Node `%v` is issued", req.Hostname)
    return keys, nil
}

type packedAuthKeyCert struct {
    Auth []byte `json:"auth"`
    Key  []byte `json:"key"`
    Cert []byte `json:"cert"`
}

// createSignedCertificate generates private key and certificate signed
// by the host certificate authority, listing the role of this server
func createSignedCertificate(caSigner *pcrypto.CaSigner, req *signedCertificateReq) (*packedAuthKeyCert, error) {
    // TODO : check if signed cert for this uuid exists. If does, return the value

    a := caSigner.CertificateAuthority()
    _, k, _, err := pcrypto.GenerateStrongKeyPair()
    if err != nil {
        return nil, trace.Wrap(err)
    }
    c, err := caSigner.GenerateSignedCertificate(req.Hostname, req.IP4Addr, k)
    if err != nil {
        log.Warningf("[AUTH] Node `%v` cannot receive a signed certificate : cert generation error. %v", req.Hostname, err)
        return nil, trace.Wrap(err)
    }
    return &packedAuthKeyCert{
        Auth: a,
        Key:  k,
        Cert: c,
    }, nil
}
