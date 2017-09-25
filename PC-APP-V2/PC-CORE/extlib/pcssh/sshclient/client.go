package sshclient

import (
    "fmt"
    "os"
    "time"

    "golang.org/x/crypto/ssh"
    "github.com/pkg/errors"

    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/auth"
    "github.com/gravitational/teleport/lib/client"
    "github.com/gravitational/teleport/lib/defaults"
    "github.com/gravitational/teleport/lib/service"
)

// makeClient takes the command-line configuration and constructs & returns
// a fully configured TeleportClient object
func MakeNewClient(cfg *service.PocketConfig, login, targetHost string) (tc *client.TeleportClient, err error) {
    var labels map[string]string
    fPorts, err := client.ParsePortForwardSpec([]string{})
    if err != nil {
        return nil, err
    }

    id, err := auth.ReadIdentityFromCertStorage(cfg.CoreProperty.CertStorage,
        auth.IdentityID{
            HostUUID: cfg.HostUUID,
            Role: teleport.RoleAdmin})
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // prep client config:
    c := &client.Config{
        Stdout:             os.Stdout,
        Stderr:             os.Stderr,
        Stdin:              os.Stdin,

        // Equal to SetProxy()
        ProxyHostPort:      fmt.Sprintf("127.0.0.1:%d,%d", defaults.HTTPListenPort, defaults.SSHProxyListenPort),
        // Username is the Teleport user's username (to login into proxies)
        Username:           login,
        // SiteName is equivalient to --cluster argument
        SiteName:           "",
        // Target Host to issue SSH command
        Host:               targetHost,
        // SSH Port on a remote SSH host
        HostPort:           int(defaults.SSHServerListenPort),
        // Login on a remote SSH host
        HostLogin:          login,
        Labels:             labels,
        // TTL defines how long a session must be active (in minutes)
        KeyTTL:             time.Minute * time.Duration(defaults.CertDuration / time.Minute),
        // InsecureSkipVerify bypasses verification of HTTPS certificate when talking to web proxy
        InsecureSkipVerify: false,
        SkipLocalAuth:      false,
        AuthMethods:        []ssh.AuthMethod{ssh.PublicKeys(id.KeySigner)},
        LocalForwardPorts:  fPorts,
        // Interactive, when set to true, launches remote command with the terminal attached
        Interactive:        false,
    }
    return client.NewClient(c)
}