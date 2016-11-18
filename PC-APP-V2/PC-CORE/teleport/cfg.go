package teleport

import (
    "os"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/teleport/lib/defaults"
    "github.com/gravitational/teleport/lib/service"
    "github.com/stkim1/pc-core/context"
)
// MakeDefaultConfig creates a new Config structure and populates it with defaults
func MakePocketTeleportConfig() (config *service.Config) {
    config = &service.Config{}
    applyDefaults(config, context.SharedHostContext())
    return config
}

// applyDefaults applies default values to the existing config structure
func applyDefaults(cfg *service.Config, context context.HostContext) {
    hostname, err := os.Hostname()
    if err != nil {
        hostname = "localhost"
        log.Errorf("Failed to determine hostname: %v", err)
    }
    cfg.SeedConfig = false

    // defaults for the auth service:
    cfg.Auth.Enabled = true
    cfg.Auth.SSHAddr = *defaults.AuthListenAddr()
    cfg.Auth.EventsBackend.Type = defaults.BackendType
    //cfg.Auth.EventsBackend.Params = boltParams(defaults.DataDir, defaults.EventsBoltFile)
    cfg.Auth.KeysBackend.Type = defaults.BackendType
    //cfg.Auth.KeysBackend.Params = boltParams(defaults.DataDir, defaults.KeysBoltFile)
    cfg.Auth.RecordsBackend.Type = defaults.BackendType
    //cfg.Auth.RecordsBackend.Params = boltParams(defaults.DataDir, defaults.RecordsBoltFile)
    defaults.ConfigureLimiter(&cfg.Auth.Limiter)

    // defaults for the SSH proxy service:
    cfg.Proxy.Enabled = true
    cfg.Proxy.AssetsDir = defaults.DataDir
    cfg.Proxy.SSHAddr = *defaults.ProxyListenAddr()
    cfg.Proxy.WebAddr = *defaults.ProxyWebListenAddr()
    cfg.Proxy.ReverseTunnelListenAddr = *defaults.ReverseTunnellListenAddr()
    defaults.ConfigureLimiter(&cfg.Proxy.Limiter)

    // defaults for the SSH service:
    cfg.SSH.Enabled = true
    cfg.SSH.Addr = *defaults.SSHServerListenAddr()
    cfg.SSH.Shell = defaults.DefaultShell
    defaults.ConfigureLimiter(&cfg.SSH.Limiter)

    // global defaults
    cfg.Hostname = hostname
    cfg.DataDir = defaults.DataDir
    cfg.Console = os.Stdout
}

