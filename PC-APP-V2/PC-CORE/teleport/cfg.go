package teleport

import (
    "os"
    "fmt"
    "path/filepath"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/teleport/lib/service"
    "github.com/stkim1/pc-core/context"
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
    hostname, err := os.Hostname()
    if err != nil {
        hostname = "localhost"
        log.Errorf("Failed to determine hostname: %v", err)
    }

    appDataDir, err := context.ApplicationUserDataDirectory()
    if err != nil {
        log.Errorf("Failed to determine hostname: %v", err)
    }
    dataDir := appDataDir + "/teleport"
    // check if the path exists and make it if absent
    if _, err := os.Stat(dataDir); err != nil {
        if os.IsNotExist(err) {
            os.MkdirAll(dataDir, os.ModeDir|0700);
        }
    }

    cfg.SeedConfig = false

    // defaults for the auth service:
    cfg.Auth.Enabled = true
    cfg.Auth.SSHAddr = *AuthListenAddr()
    cfg.Auth.EventsBackend.Type = BackendType
    cfg.Auth.EventsBackend.Params = dbParams(dataDir, EventsSqliteFile)
    cfg.Auth.KeysBackend.Type = BackendType
    cfg.Auth.KeysBackend.Params = dbParams(dataDir, KeysSqliteFile)
    cfg.Auth.RecordsBackend.Type = BackendType
    cfg.Auth.RecordsBackend.Params = dbParams(dataDir, RecordsSqliteFile)
    ConfigureLimiter(&cfg.Auth.Limiter)

    // defaults for the SSH proxy service:
    cfg.Proxy.Enabled = true
    cfg.Proxy.AssetsDir = dataDir
    cfg.Proxy.SSHAddr = *ProxyListenAddr()
    cfg.Proxy.WebAddr = *ProxyWebListenAddr()
    cfg.Proxy.ReverseTunnelListenAddr = *ReverseTunnellListenAddr()
    ConfigureLimiter(&cfg.Proxy.Limiter)

    // defaults for the SSH service:
    cfg.SSH.Enabled = true
    cfg.SSH.Addr = *SSHServerListenAddr()
    cfg.SSH.Shell = DefaultShell
    ConfigureLimiter(&cfg.SSH.Limiter)

    // global defaults
    cfg.Hostname = hostname
    cfg.DataDir = dataDir
    cfg.Console = os.Stdout
}

