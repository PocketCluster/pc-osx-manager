package pcauth

import (
    "fmt"
    "os"
    "time"

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
    role        teleport.Role
    alog        events.IAuditLog
    signer      pcrypto.CaSigner
}

func (a *authWithRoles) issueSignedCertificateWithToken(req *signedCertificateReq) (*auth.PackedKeys, error) {
    if err := a.permChecker.HasPermission(a.role, auth.ActionRegisterUsingToken); err != nil {
        return nil, trace.Wrap(err)
    }
    return issueSignedCertificateWithToken(a, req.Token, req.Hostname, req.HostUUID, req.IPAddress, req.Role)
}

// issueSignedCertificateWithToken adds a new signed certificate for a node to the PocketCluster using previously issued token.
// A node must also request a specific role (and the role must match one of the roles the token was generated for).
//
// If a token was generated with a TTL, it gets enforced (can't register new nodes after TTL expires)
// If a token was generated with a TTL=0, it means it's a single-use token and it gets destroyed
// after a successful registration.
func issueSignedCertificateWithToken(a *authWithRoles, token, hostname, hostUUID, ipAddress string, role teleport.Role) (*auth.PackedKeys, error) {
    log.Infof("[AUTH] Node `%v` requests a signed certificate", hostname)
    if len(hostname) == 0 {
        return nil, trace.BadParameter("Hostname cannot be empty")
    }
    if len(hostUUID) == 0 {
        return nil, trace.BadParameter("HostID cannot be empty")
    }
    if err := role.Check(); err != nil {
        return nil, trace.Wrap(err)
    }
    // make sure the token is valid:
    roles, err := a.authServer.ValidateToken(token)
    if err != nil {
        msg := fmt.Sprintf("`%v` cannot join the cluster as %s. Token error: %v", hostname, role, err)
        log.Warnf("[AUTH] %s", msg)
        return nil, trace.AccessDenied(msg)
    }
    // make sure the caller is requested wthe role allowed by the token:
    if !roles.Include(role) {
        msg := fmt.Sprintf("'%v' cannot join the cluster, the token does not allow '%s' role", hostname, role)
        log.Warningf("[AUTH] %s", msg)
        return nil, trace.BadParameter(msg)
    }
    if !checkTokenTTL(a.authServer, token) {
        return nil, trace.AccessDenied("'%v' cannot join the cluster. The token has expired", hostname)
    }
    // generate & return the node cert:
    keys, err := generateSignedCertificate(a, hostname, hostUUID, ipAddress)
    if err != nil {
        return nil, trace.Wrap(err)
    }
    utils.Consolef(os.Stdout, "[AUTH] Signed Certificate for Node `%v` issued", hostname)
    return keys, nil
}

// enforceTokenTTL deletes the given token if it's TTL is over. Returns 'false'
// if this token cannot be used
func checkTokenTTL(s *auth.AuthServer, token string) bool {
    // look at the tokens in the token storage
    tok, err := s.Provisioner.GetToken(token)
    if err != nil {
        log.Warn(err)
        return true
    }
    // s.clock is replaced with time.Now()
    now := time.Now().UTC()
    if tok.Expires.Before(now) {
        if err = s.DeleteToken(token); err != nil {
            log.Error(err)
        }
        return false
    }
    return true
}

// generateSignedCertificate generates private key and certificate signed
// by the host certificate authority, listing the role of this server
func generateSignedCertificate(a *authWithRoles, hostname, hostUUID, ipAddress string) (*auth.PackedKeys, error) {
    _, nodeKey, _, err := pcrypto.GenerateStrongKeyPair()
    if err != nil {
        return nil, trace.Wrap(err)
    }

    // TODO : check if signed cert for this uuid exists. If does, return the value


    // we always append authority's domain to resulting node name,
    // that's how we make sure that nodes are uniquely identified/found
    // in cases when we have multiple environments/organizations
    //fqdn := fmt.Sprintf("%s.%s", hostID, s.DomainName)
    // TODO : fix parameter input values
    c, err := a.signer.GenerateSignedCertificate(hostname, ipAddress, nodeKey)
    if err != nil {
        log.Warningf("[AUTH] Node `%v` cannot join: cert generation error. %v", hostname, err)
        return nil, trace.Wrap(err)
    }
    return &auth.PackedKeys{
        Key:  nodeKey,
        Cert: c,
    }, nil
}
