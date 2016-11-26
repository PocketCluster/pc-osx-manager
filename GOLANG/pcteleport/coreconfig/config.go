package coreconfig

import (
    "os"
    "fmt"
    "io/ioutil"
    "path/filepath"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/trace"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pcteleport/pcdefaults"
)

// MakeDefaultConfig creates a new Config structure and populates it with defaults
func MakePocketTeleportConfig() (config *service.Config) {
    config = &service.Config{}
    applyDefaults(config, context.SharedHostContext())
    return config
}

// Generates a string accepted by the BoltDB driver, like this:
// `{"path": "/var/lib/teleport/records.db"}`
func dbParams(storagePath, dbFile string) string {
    return fmt.Sprintf(`{"path": "%s"}`, filepath.Join(storagePath, dbFile))
}

// applyDefaults applies default values to the existing config structure
func applyDefaults(cfg *service.Config, context context.HostContext) {
    var (
        hostname, appDataDir, dataDir string = "", "", ""
        err error = nil
    )

    hostname, err = os.Hostname()
    if err != nil {
        hostname = "localhost"
        log.Errorf("Failed to determine hostname: %v", err)
    }

    appDataDir, err = context.ApplicationUserDataDirectory()
    if err != nil {
        log.Errorf("Failed to determine hostname: %v", err)
    }
    dataDir = appDataDir + "/teleport"
    // check if the path exists and make it if absent
    if _, err := os.Stat(dataDir); err != nil {
        if os.IsNotExist(err) {
            os.MkdirAll(dataDir, os.ModeDir|0700);
        }
    }

    cfg.SeedConfig = false

    // defaults for the auth service:
    cfg.AuthServers = []utils.NetAddr{*pcdefaults.AuthConnectAddr()}
    cfg.Auth.Enabled = true
    cfg.Auth.SSHAddr = *pcdefaults.AuthListenAddr()
    cfg.Auth.EventsBackend.Type = pcdefaults.CoreBackendType
    cfg.Auth.EventsBackend.Params = dbParams(dataDir, pcdefaults.CoreEventsSqliteFile)
    cfg.Auth.KeysBackend.Type = pcdefaults.CoreBackendType
    cfg.Auth.KeysBackend.Params = dbParams(dataDir, pcdefaults.CoreKeysSqliteFile)
    cfg.Auth.RecordsBackend.Type = pcdefaults.CoreBackendType
    cfg.Auth.RecordsBackend.Params = dbParams(dataDir, pcdefaults.CoreRecordsSqliteFile)
    pcdefaults.ConfigureLimiter(&cfg.Auth.Limiter)

    // defaults for the SSH proxy service:
    cfg.Proxy.Enabled = true
    // disable web ui as it's not necessary
    //cfg.Proxy.DisableWebUI = true
    cfg.Proxy.AssetsDir = dataDir
    cfg.Proxy.SSHAddr = *pcdefaults.ProxyListenAddr()
    cfg.Proxy.WebAddr = *pcdefaults.ProxyWebListenAddr()

    cfg.Proxy.ReverseTunnelListenAddr = *pcdefaults.ReverseTunnellListenAddr()
    pcdefaults.ConfigureLimiter(&cfg.Proxy.Limiter)

    // defaults for the SSH service:
    cfg.SSH.Enabled = false
    cfg.SSH.Addr = *pcdefaults.SSHServerListenAddr()
    cfg.SSH.Shell = pcdefaults.DefaultShell
    pcdefaults.ConfigureLimiter(&cfg.SSH.Limiter)

    // global defaults
    cfg.Hostname = hostname
    cfg.DataDir = dataDir
    cfg.Console = os.Stdout
}

func ValidateConfig(cfg *service.Config) error {
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