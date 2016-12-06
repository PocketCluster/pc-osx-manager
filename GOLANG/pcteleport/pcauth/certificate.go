package pcauth

import (
    "encoding/json"
    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/auth"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/stkim1/pcrypto"
)

// RequestSignedCertificate is used by auth service clients (other services, like proxy or SSH) when a new node joins
// the cluster
func RequestSignedCertificate(dataDir, token, hostname, ip4Addr string, id auth.IdentityID, servers []utils.NetAddr) error {
    tok, err := readToken(token)
    if err != nil {
        return trace.Wrap(err)
    }
    method, err := auth.NewTokenAuth(id.HostUUID, tok)
    if err != nil {
        return trace.Wrap(err)
    }

    client, err := auth.NewTunClient(
        "auth.client.request",
        servers,
        id.HostUUID,
        method)
    if err != nil {
        return trace.Wrap(err)
    }
    defer client.Close()

    keys, err := requestSignedCertificateWithToken(client, tok, id.HostUUID, hostname, ip4Addr, id.Role)
    if err != nil {
        return trace.Wrap(err)
    }
    return writeKeys(dataDir, id, keys.Key, keys.Cert)
}

// requestSignedCertificateWithToken calls the auth service API to register a new node via registration token which has
// been previously issued via GenerateToken
func requestSignedCertificateWithToken(c *auth.TunClient, token, hostID, hostname, ip4Addr string, role teleport.Role) (*auth.PackedKeys, error) {
    out, err := c.PostJSON(c.Endpoint(pcHostToken, pcHostID, pcHostName, pcHostIp4Addr, pcHostRole),
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
    var keys auth.PackedKeys
    if err := json.Unmarshal(out.Bytes(), &keys); err != nil {
        return nil, trace.Wrap(err)
    }
    return &keys, nil
}

// createSignedCertificate generates private key and certificate signed
// by the host certificate authority, listing the role of this server
func createSignedCertificate(signer pcrypto.CaSigner, hostID, hostname, ipAddress string) (*auth.PackedKeys, error) {
    // TODO : check if signed cert for this uuid exists. If does, return the value

    _, nodeKey, _, err := pcrypto.GenerateStrongKeyPair()
    if err != nil {
        return nil, trace.Wrap(err)
    }
    c, err := signer.GenerateSignedCertificate(hostname, ipAddress, nodeKey)
    if err != nil {
        log.Warningf("[AUTH] Node `%v` cannot receive a signed certificate : cert generation error. %v", hostname, err)
        return nil, trace.Wrap(err)
    }
    return &auth.PackedKeys{
        Key:  nodeKey,
        Cert: c,
    }, nil
}
