package pcauth

import (
    "encoding/json"
    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/auth"

    "github.com/gravitational/trace"
    "github.com/stkim1/pcteleport/pcconfig"
)

// RequestSignedCertificate is used by auth service clients (other services, like proxy or SSH) when a new node joins
// the cluster
func RequestSignedCertificate(cfg *pcconfig.Config, id auth.IdentityID, token string) error {
    tok, err := readToken(token)
    if err != nil {
        return trace.Wrap(err)
    }
    method, err := auth.NewTokenAuth(id.HostUUID, tok)
    if err != nil {
        return trace.Wrap(err)
    }

    var servers []utils.NetAddr = cfg.AuthServers
    client, err := auth.NewTunClient(
        "auth.client.cert.reqsigned",
        servers,
        id.HostUUID,
        method)
    if err != nil {
        return trace.Wrap(err)
    }
    defer client.Close()

    keys, err := requestSignedCertificateWithToken(client, tok, id.HostUUID, cfg.Hostname, cfg.IP4Addr, id.Role)
    if err != nil {
        return trace.Wrap(err)
    }
    return writeDockerKeyAndCert(cfg, keys)
}

// requestSignedCertificateWithToken calls the auth service API to register a new node via registration token which has
// been previously issued via GenerateToken
func requestSignedCertificateWithToken(c *auth.TunClient, token, hostID, hostname, ip4Addr string, role teleport.Role) (*packedAuthKeyCert, error) {
    out, err := c.PostJSON(apiEndpoint(PocketCertificate, PocketRequestSigned),
        signedCertificateReq{
            Token:      token,
            HostID:     hostID,
            Hostname:   hostname,
            IP4Addr:    ip4Addr,
            Role:       role,
        })
    if err != nil {
        return nil, trace.Wrap(err)
    }
    var keys packedAuthKeyCert
    if err := json.Unmarshal(out.Bytes(), &keys); err != nil {
        return nil, trace.Wrap(err)
    }
    return &keys, nil
}
