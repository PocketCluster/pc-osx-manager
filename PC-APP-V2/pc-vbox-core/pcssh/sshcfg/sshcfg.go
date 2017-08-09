package sshcfg

import (
    "net"
    "os"
    "io/ioutil"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/gravitational/teleport/lib/defaults"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/stkim1/pc-vbox-core/crcontext"
)

// MakeDefaultConfig creates a new Config structure and populates it with defaults
func MakeCoreConfig(ctx crcontext.PocketCoreContext, debug bool) (*service.PocketConfig, error) {
    config := &service.PocketConfig{}
    err := applyCoreDefaults(config, ctx, debug)
    return config, err
}

func applyCoreDefaults(cfg *service.PocketConfig, ctx crcontext.PocketCoreContext, debug bool) error {
    authServerAddr, err := ctx.CoreAuthServerAddr()
    if err != nil {
        return errors.WithStack(err)
    }

    addr, err := utils.ParseHostPortAddr(authServerAddr, int(defaults.AuthListenPort))
    if err != nil {
        return errors.WithStack(err)
    }
    log.Infof("Using auth server: %v", addr.FullAddress())
    // dataDir should have been created before pcteleport is executed
    dataDir := ctx.CoreConfigPath()
    // check if the path exists and report error if absent
    if _, err := os.Stat(dataDir); err != nil {
        return errors.WithStack(err)
    }
    keyCertDir := ctx.CoreKeyAndCertPath()
    // check if the path exists and report error if absent
    if _, err := os.Stat(keyCertDir); err != nil {
        return errors.WithStack(err)
    }
    log.Printf("DataDir: %v, KeyCertDir %v", dataDir, keyCertDir)

    // global defaults
    nodeName                   := ctx.CoreNodeName()
    hostUUID                   := ctx.CoreNodeUUID()
    authToken                  := ctx.CoreAuthToken()

    // defaults for the auth service:
    cfg.SeedConfig             = false

    cfg.Auth.Enabled           = false
    cfg.AuthServers            = []utils.NetAddr{*addr}
    cfg.Auth.SSHAddr           = *defaults.AuthListenAddr()
    cfg.ApplyToken(authToken)

    // defaults for the SSH service:
    cfg.SSH.Enabled            = true
    cfg.SSH.Addr               = *defaults.SSHServerListenAddr()
    cfg.SSH.Shell              = defaults.DefaultShell
    defaults.ConfigureLimiter(&cfg.SSH.Limiter)

    cfg.Hostname               = nodeName
    cfg.HostUUID               = hostUUID
    cfg.DataDir                = dataDir

    cfg.AdvertiseIP            = net.ParseIP(crcontext.CoreSSHAdvertiseAddr())
    cfg.NodeSSHCertificateFile = ctx.CoreSSHCertificateFileName()
    cfg.NodeSSHPrivateKeyFile  = ctx.CoreSSHPrivateKeyFileName()

    if debug {
        cfg.Console     = os.Stdout
        utils.InitLoggerDebug()
        log.Info("Teleport DEBUG output configured")
    } else {
        cfg.Console     = ioutil.Discard
        utils.InitLoggerCLI()
        log.Info("Teleport NORMAL output configured")
    }
    return nil
}

func ValidateCoreConfig(cfg *service.PocketConfig) error {
    if !cfg.Auth.Enabled && !cfg.SSH.Enabled && !cfg.Proxy.Enabled {
        return errors.Errorf("config: supply at least one of Auth, SSH or Proxy roles")
    }

    if cfg.DataDir == "" {
        return errors.Errorf("config: please supply data directory")
    }

    if cfg.Console == nil {
        cfg.Console = ioutil.Discard
    }

    if (cfg.Proxy.TLSKey == "" && cfg.Proxy.TLSCert != "") || (cfg.Proxy.TLSKey != "" && cfg.Proxy.TLSCert == "") {
        return errors.Errorf("please supply both TLS key and certificate")
    }

    if len(cfg.AuthServers) == 0 {
        return errors.Errorf("auth_servers is empty")
    }
    for i := range cfg.Auth.Authorities {
        if err := cfg.Auth.Authorities[i].Check(); err != nil {
            return errors.WithStack(err)
        }
    }
    for _, tun := range cfg.ReverseTunnels {
        if err := tun.Check(); err != nil {
            return errors.WithStack(err)
        }
    }

    return nil
}