package pcconfig

import (
    "errors"
    "os"
    "io/ioutil"
    "path/filepath"

    "github.com/gravitational/teleport/lib/utils"
    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"

    "github.com/stkim1/pcteleport/pcdefaults"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slcontext/config"
)

// MakeDefaultConfig creates a new Config structure and populates it with defaults
func MakeNodeTeleportConfig(authServerAddr, authToken string, debug bool) (*Config, error) {
    config := &Config{}
    err := applyNodeDefaults(config, slcontext.SharedSlaveContext(), authServerAddr, authToken, debug)
    return config, err
}

// applyDefaults applies default values to the existing config structure
func applyNodeDefaults(cfg *Config, context slcontext.PocketSlaveContext, authServerAddr, authToken string, debug bool) error {
    addr, err := utils.ParseHostPortAddr(authServerAddr, int(pcdefaults.AuthListenPort))
    if err != nil {
        return trace.Wrap(err)
    }
    log.Infof("Using auth server: %v", addr.FullAddress())
    if len(authToken) == 0 {
        return trace.Wrap(errors.New("[ERR] Invalid AuthToken"))
    }
    // dataDir should have been created before pcteleport is executed
    dataDir := context.SlaveConfigPath()
    // check if the path exists and report error if absent
    if _, err := os.Stat(dataDir); err != nil {
        return trace.Wrap(err)
    }
    keyCertDir := context.SlaveKeyAndCertPath()
    // check if the path exists and report error if absent
    if _, err := os.Stat(keyCertDir); err != nil {
        return trace.Wrap(err)
    }
    log.Printf("DataDir: %v, KeyCertDir %v", dataDir, keyCertDir)

    // global defaults
    nodeName, err := context.GetSlaveNodeName()
    if err != nil {
        return trace.Wrap(err)
    }
    // get current network interface address
    iface, err := context.PrimaryNetworkInterface()
    if err != nil {
        // TODO if this keeps fail, we'll enforce to get current interface
        log.Errorf("Failed to determine network interface: %v", err)
        return trace.Wrap(err)
    }
    // TODO : read host UUID from slave context
    // if there's no host uuid initialized yet, try to read one from the
    // one of the identities
    hostUUID, err := context.GetSlaveNodeUUID()
    if err != nil {
        log.Errorf(err.Error())
        return trace.Wrap(err)
    }

    // defaults for the auth service:
    cfg.SeedConfig      = false

    cfg.Auth.Enabled    = false
    cfg.AuthServers     = []utils.NetAddr{*addr}
    cfg.Auth.SSHAddr    = *pcdefaults.AuthListenAddr()
    cfg.ApplyToken(authToken)

/*
    cfg.Auth.EventsBackend.Type = pcdefaults.CoreBackendType
    cfg.Auth.EventsBackend.Params = dbParams(dataDir, pcdefaults.CoreEventsSqliteFile)
    cfg.Auth.KeysBackend.Type = pcdefaults.CoreBackendType
    cfg.Auth.KeysBackend.Params = dbParams(dataDir, pcdefaults.CoreKeysSqliteFile)
    cfg.Auth.RecordsBackend.Type = pcdefaults.CoreBackendType
    cfg.Auth.RecordsBackend.Params = dbParams(dataDir, pcdefaults.CoreRecordsSqliteFile)
    pcdefaults.ConfigureLimiter(&cfg.Auth.Limiter)

    // defaults for the SSH proxy service:
    cfg.Proxy.Enabled = false
    cfg.Proxy.DisableWebUI = false
    cfg.Proxy.AssetsDir = dataDir
    cfg.Proxy.SSHAddr = *pcdefaults.ProxyListenAddr()
    cfg.Proxy.WebAddr = *pcdefaults.ProxyWebListenAddr()
    cfg.Proxy.ReverseTunnelListenAddr = *pcdefaults.ReverseTunnellListenAddr()
    pcdefaults.ConfigureLimiter(&cfg.Proxy.Limiter)
*/

    // defaults for the SSH service:
    cfg.SSH.Enabled     = true
    cfg.SSH.Addr        = *pcdefaults.SSHServerListenAddr()
    cfg.SSH.Shell       = pcdefaults.DefaultShell
    pcdefaults.ConfigureLimiter(&cfg.SSH.Limiter)

    cfg.Hostname        = nodeName
    cfg.DataDir         = dataDir

    cfg.IP4Addr         = iface.IP.String()
    cfg.DockerAuthFile  = filepath.Join(keyCertDir, config.SlaveDockerAuthFileName)
    cfg.DockerKeyFile   = filepath.Join(keyCertDir, config.SlaveDockerKeyFileName)
    cfg.DockerCertFile  = filepath.Join(keyCertDir, config.SlaveDockerCertFileName)

    cfg.HostUUID        = hostUUID

    // if user did not provide auth domain name, use this host UUID
    if cfg.Auth.Enabled && cfg.Auth.DomainName == "" {
        log.Info("cfg.Auth.DomainName set to UUID")
        cfg.Auth.DomainName = cfg.HostUUID
    }

    if debug {
        cfg.Console     = os.Stdout
        utils.InitLoggerDebug()
        trace.SetDebug(true)
        log.Info("Teleport DEBUG output configured")
    } else {
        cfg.Console     = ioutil.Discard
        utils.InitLoggerCLI()
        log.Info("Teleport NORMAL output configured")
    }
    return nil
}

func ValidateNodeConfig(cfg *Config) error {
    if !cfg.Auth.Enabled && !cfg.SSH.Enabled && !cfg.Proxy.Enabled {
        return trace.BadParameter(
            "config: supply at least one of Auth, SSH or Proxy roles")
    }

    if cfg.DataDir == "" {
        return trace.BadParameter("config: please supply data directory")
    }

    if cfg.Console == nil {
        cfg.Console = ioutil.Discard
    }

    if (cfg.Proxy.TLSKey == "" && cfg.Proxy.TLSCert != "") || (cfg.Proxy.TLSKey != "" && cfg.Proxy.TLSCert == "") {
        return trace.BadParameter("please supply both TLS key and certificate")
    }

    if len(cfg.AuthServers) == 0 {
        return trace.BadParameter("auth_servers is empty")
    }
    for i := range cfg.Auth.Authorities {
        if err := cfg.Auth.Authorities[i].Check(); err != nil {
            return trace.Wrap(err)
        }
    }
    for _, tun := range cfg.ReverseTunnels {
        if err := tun.Check(); err != nil {
            return trace.Wrap(err)
        }
    }

    return nil
}