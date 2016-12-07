package pcconfig


import (
    "os"
    "io/ioutil"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/trace"

    "github.com/stkim1/pcteleport/pcdefaults"
    "github.com/stkim1/pc-node-agent/slcontext"
)

// MakeDefaultConfig creates a new Config structure and populates it with defaults
func MakeNodeTeleportConfig() (config *Config) {
    config = &Config{}
    applyNodeDefaults(config, slcontext.SharedSlaveContext())
    return config
}

// applyDefaults applies default values to the existing config structure
func applyNodeDefaults(cfg *Config, context slcontext.PocketSlaveContext) {
    var (
        hostname, dataDir string = "", ""
        err error = nil
    )

    hostname, err = os.Hostname()
    if err != nil {
        hostname = "localhost"
        log.Errorf("Failed to determine hostname: %v", err)
    }

    // dataDir should have been created before pcteleport is executed
    dataDir = context.SlaveConfigPath()
    cfg.SeedConfig = false

    // defaults for the auth service:
    cfg.AuthServers = []utils.NetAddr{*pcdefaults.AuthConnectAddr()}
    cfg.Auth.Enabled = false
    cfg.Auth.SSHAddr = *pcdefaults.AuthListenAddr()
    cfg.Auth.EventsBackend.Type = pcdefaults.CoreBackendType
    cfg.Auth.EventsBackend.Params = dbParams(dataDir, pcdefaults.CoreEventsSqliteFile)
    cfg.Auth.KeysBackend.Type = pcdefaults.CoreBackendType
    cfg.Auth.KeysBackend.Params = dbParams(dataDir, pcdefaults.CoreKeysSqliteFile)
    cfg.Auth.RecordsBackend.Type = pcdefaults.CoreBackendType
    cfg.Auth.RecordsBackend.Params = dbParams(dataDir, pcdefaults.CoreRecordsSqliteFile)
    pcdefaults.ConfigureLimiter(&cfg.Auth.Limiter)

    // defaults for the SSH proxy service:
    cfg.Proxy.Enabled = false
    // disable web ui as it's not necessary
    cfg.Proxy.DisableWebUI = false
    cfg.Proxy.AssetsDir = dataDir
    cfg.Proxy.SSHAddr = *pcdefaults.ProxyListenAddr()
    cfg.Proxy.WebAddr = *pcdefaults.ProxyWebListenAddr()

    cfg.Proxy.ReverseTunnelListenAddr = *pcdefaults.ReverseTunnellListenAddr()
    pcdefaults.ConfigureLimiter(&cfg.Proxy.Limiter)

    // defaults for the SSH service:
    cfg.SSH.Enabled = true
    cfg.SSH.Addr = *pcdefaults.SSHServerListenAddr()
    cfg.SSH.Shell = pcdefaults.DefaultShell
    pcdefaults.ConfigureLimiter(&cfg.SSH.Limiter)

    // global defaults
    cfg.Hostname = hostname
    cfg.DataDir = dataDir
    // TODO remove Stdout
    cfg.Console = os.Stdout

    // get current network interface address
    iface, err := context.PrimaryNetworkInterface()
    if err != nil {
        // TODO if this keeps fail, we'll enforce to get current interface
        log.Errorf("Failed to determine network interface: %v", err)
    }
    cfg.IP4Addr = iface.IP.String()
    cfg.KeyCertDir = context.SlaveKeyAndCertPath()
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