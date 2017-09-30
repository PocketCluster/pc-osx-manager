package main

import (
    "path/filepath"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/auth"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/services"
    "github.com/gravitational/trace"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pcteleport/pcconfig"
    "golang.org/x/crypto/ssh"
    "github.com/stkim1/pcteleport/pcauth"
)

// connectToAuthService creates a valid client connection to the auth service
func connectToAuthService(cfg *pcconfig.Config) (client *auth.TunClient, err error) {
    identity, err := auth.ReadIdentity(cfg.DataDir, auth.IdentityID{HostUUID: cfg.HostUUID, Role: teleport.RoleAdmin})
    if err != nil {
        return nil, trace.Wrap(err)
    }
    authUser := identity.Cert.ValidPrincipals[0]
    client, err = auth.NewTunClient(
        "pc-osx-manager",
        cfg.AuthServers,
        authUser,
        []ssh.AuthMethod{ssh.PublicKeys(identity.KeySigner)},
    )
    if err != nil {
        return nil, trace.Wrap(err)
    }
    // check connectivity by calling something on a clinet:
    _, err = client.GetDialer()()
    if err != nil {
        return nil, trace.Wrap(err, "Cannot connect to the auth server: %v.\nIs the auth server running on %v?", err, cfg.AuthServers[0].Addr)
    }
    return client, nil
}

func createSignupToken(client *auth.TunClient, login string) (string, error) {
    user := services.TeleportUser{
        Name:          login,
        AllowedLogins: []string{login},
    }
    return client.CreateSignupToken(&user)
}

func create_user() {
    context.DebugContextPrepare()
    cfg := pcconfig.MakeCoreTeleportConfig(true)
    clt, err := connectToAuthService(cfg)
    if err != nil {
        log.Error(err.Error())
        return
    }
    defer clt.Close()
    token, err := createSignupToken(clt, "root")
    if err != nil {
        log.Error(err.Error())
        return
    }
    // ed61fba1372a29469e544c6a5aaf5082
    log.Infof("Signup Token %v", token)
    hotpToken, err := pcauth.RequestHOTPforSignupToken(clt, token)
    if err != nil {
        log.Error(err.Error())
        return
    }
    for _, h := range hotpToken {
        log.Info(h)
    }

    _, err = clt.CreateUserWithToken(token, "1524rmfo", hotpToken[0])
    if err != nil {
        log.Error(err.Error())
        return
    }
}

func create_user_old() {
    context.DebugContextPrepare()
    cfg := pcconfig.MakeCoreTeleportConfig(true)

    id := auth.IdentityID{HostUUID: cfg.HostUUID, Role: teleport.RoleAdmin}
    identity, err := auth.ReadIdentity(cfg.DataDir, id)
    if err != nil {
        log.Error(err.Error())
        return
    }
    storage := utils.NewFileAddrStorage(filepath.Join(cfg.DataDir, "authservers.json"))
    authUser := identity.Cert.ValidPrincipals[0]

    clt, err := auth.NewTunClient(
        string(teleport.RoleAdmin),
        cfg.AuthServers,
        authUser,
        []ssh.AuthMethod{ssh.PublicKeys(identity.KeySigner)},
        auth.TunClientStorage(storage),
    )
    // success?
    if err != nil {
        log.Error(err.Error())
        return
    }
    // try calling a test method via auth api:
    //
    // ??? in case of failure it never gets back here!!!
    dn, err := clt.GetLocalDomain()
    if err != nil {
        log.Error(err.Error())
        return
    }
    // success ? we're logged in!
    log.Infof("[Node] %s connected to the cluster '%s'", authUser, dn)
    //return &service.Connector{Client: authClient, Identity: identity}, nil

/*
    priv, pub, err := clt.GenerateKeyPair("")
    if err != nil {
        log.Error(err.Error())
    }
    // should NOT be able to generate a user cert without basic HTTP auth
    cert, err := clt.GenerateUserCert(pub, "user1", time.Hour)
    if err != nil {
        log.Error(err.Error())
    }
*/

    hotpURL, _, err := clt.UpsertPassword("root", []byte("1524rmfo"))
    if err != nil {
        log.Error(err.Error())
    }
    log.Infof("hotpURL %s", hotpURL)

    user := &services.TeleportUser{Name: "root", AllowedLogins: []string{"root"}}
    err = clt.UpsertUser(user)
    if err != nil {
        log.Error(err.Error())
    }
    clt.DeleteUser("root")
    clt.Close()
}

func delet_user(login string) {
    context.DebugContextPrepare()
    cfg := pcconfig.MakeCoreTeleportConfig(true)
    clt, err := connectToAuthService(cfg)
    if err != nil {
        log.Error(err.Error())
        return
    }
    defer clt.Close()
    clt.DeleteUser("root")
}

func main() {
    delet_user("root")
    create_user()
}
