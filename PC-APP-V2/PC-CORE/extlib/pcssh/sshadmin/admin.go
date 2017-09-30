package sshadmin

import (
    "fmt"
    "time"

    "golang.org/x/crypto/ssh"
    "github.com/pkg/errors"

    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/auth"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/services"
)

const (
    MaxInvitationTLL time.Duration = (time.Minute * 5)
    MinInvitationTLL time.Duration = time.Minute
)

// generates an invitation token which can be used to add another SSH node to a cluster
func GenerateNodeInviationWithTTL(client *auth.TunClient, ttl time.Duration) (string, error) {
    roles, err := teleport.ParseRoles("node")
    if err != nil {
        return "", errors.WithStack(err)
    }

    // adjust ttl
    if ttl < MinInvitationTLL || MaxInvitationTLL < ttl {
        ttl = MaxInvitationTLL
    }
    return client.GenerateToken(roles, ttl)
}

// connectToAuthService creates a valid client connection to the auth service
func OpenAdminClientWithAuthService(cfg *service.PocketConfig) (client *auth.TunClient, err error) {
    id, err := auth.ReadIdentityFromCertStorage(cfg.CoreProperty.CertStorage,
        auth.IdentityID{
            HostUUID: cfg.HostUUID,
            Role: teleport.RoleAdmin})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    authUser := id.Cert.ValidPrincipals[0]
    client, err = auth.NewTunClient(
        "embed.admin-client",
        cfg.AuthServers,
        authUser,
        []ssh.AuthMethod{ssh.PublicKeys(id.KeySigner)},
    )
    if err != nil {
        return nil, errors.WithStack(err)
    }
    // check connectivity by calling something on a clinet:
    _, err = client.GetDialer()()
    if err != nil {
        return nil, errors.WithMessage(err, fmt.Sprintf("Cannot connect to the auth server: %v.\nIs the auth server running on %v?", err, cfg.AuthServers[0].Addr))
    }
    return client, nil
}

func CreateTeleportUser(client *auth.TunClient, user, pass string) error {
    var (
        u = &services.TeleportUser{
            Name:          user,
            AllowedLogins: []string{user},
        }
    )
    token, err := client.CreateSignupToken(u)
    if err != nil {
        return errors.WithStack(err)
    }
    hotpToken, err := auth.RequestHOTPforSignupToken(client, token)
    if err != nil {
        return errors.WithStack(err)
    }
    _, err = client.CreateUserWithToken(token, pass, hotpToken[0])
    if err != nil {
        return errors.WithStack(err)
    }
    return nil
}