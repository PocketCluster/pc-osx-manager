package sshcfg

import (
    "database/sql"
    "fmt"
    "net"
    "os"
    "io/ioutil"
    "path/filepath"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/gravitational/teleport/lib/defaults"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/services"
    "github.com/gravitational/teleport/lib/utils"

    "github.com/cloudflare/cfssl/certdb"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/context"
)

// MakeDefaultConfig creates a new Config structure and populates it with defaults
func MakeMasterConfig(ctx context.HostContext, debug bool) (*service.PocketConfig, error) {
    config := &service.PocketConfig{}
    err := applyMasterDefaults(config, ctx)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    if debug {
        config.Console = ioutil.Discard
        log.Info("Teleport DEBUG output configured")
    } else {
        // TODO : check if we can throw debug info
        config.Console = os.Stdout
        log.Info("Teleport NORMAL cli output configured")
    }
    return config, nil
}

// applyDefaults applies default values to the existing config structure
func applyMasterDefaults(cfg *service.PocketConfig, ctx context.HostContext) error {

    dataDir, err := ctx.ApplicationUserDataDirectory()
    if err != nil {
        return errors.WithStack(err)
    }
    paddr, err := ctx.HostPrimaryAddress()
    if err != nil {
        return errors.WithStack(err)
    }
    cUUID, err := ctx.GetClusterUUID()
    if err != nil {
        return errors.WithStack(err)
    }

    // defaults for the auth service:
    cfg.SeedConfig                   = false
    cfg.Auth.Enabled                 = true
    cfg.AuthServers                  = []utils.NetAddr{*defaults.AuthConnectAddr()}
    cfg.Auth.SSHAddr                 = *defaults.AuthListenAddr()
    cfg.Auth.EventsBackend.Type      = defaults.CoreBackendType
    cfg.Auth.EventsBackend.Params    = sqliteParams(dataDir, defaults.CoreEventsSqliteFile)
    cfg.Auth.KeysBackend.Type        = defaults.CoreBackendType
    cfg.Auth.KeysBackend.Params      = sqliteParams(dataDir, defaults.CoreKeysSqliteFile)
    cfg.Auth.RecordsBackend.Type     = defaults.CoreBackendType
    cfg.Auth.RecordsBackend.Params   = sqliteParams(dataDir, defaults.CoreRecordsSqliteFile)
    defaults.ConfigureLimiter(&cfg.Auth.Limiter)

    // defaults for the SSH proxy service:
    cfg.Proxy.Enabled                 = true
    // disable web ui as it's not necessary
    cfg.Proxy.DisableWebUI            = true
    cfg.Proxy.AssetsDir               = dataDir
    cfg.Proxy.SSHAddr                 = *defaults.ProxyListenAddr()
    cfg.Proxy.WebAddr                 = *defaults.ProxyWebListenAddr()

    cfg.Proxy.ReverseTunnelListenAddr = *defaults.ReverseTunnellListenAddr()
    defaults.ConfigureLimiter(&cfg.Proxy.Limiter)

    // defaults for the SSH service:
    cfg.SSH.Enabled                   = false
    cfg.SSH.Addr                      = *defaults.SSHServerListenAddr()
    cfg.SSH.Shell                     = defaults.DefaultShell
    defaults.ConfigureLimiter(&cfg.SSH.Limiter)

    // global defaults
    cfg.Hostname                      = defaults.CoreHostName
    cfg.HostUUID                      = cUUID
    cfg.DataDir                       = dataDir
    cfg.AdvertiseIP                   = net.ParseIP(paddr)
    return nil
}

func AssignCertStorage(cfg *service.PocketConfig, certStorage certdb.Accessor) {
    cfg.CertStorage = certStorage
}

func AssignDatabaseEngine(cfg *service.PocketConfig, db *sql.DB) {
    cfg.BackendDB = db
}

func AssignCASigner(cfg *service.PocketConfig, caSigner *pcrypto.CaSigner) {
    cfg.CaSigner = caSigner
}

func AssignHostCertAuth(cfg *service.PocketConfig, private, sshCheck []byte, domainName string) {
    cfg.Auth.DomainName = domainName
    hostCA := services.CertAuthority{
        DomainName:      domainName,
        Type:            services.HostCA,
        SigningKeys:     [][]byte{private},
        CheckingKeys:    [][]byte{sshCheck},
    }
    cfg.Auth.Authorities = append(cfg.Auth.Authorities, hostCA)
}

func ValidateMasterConfig(cfg *service.PocketConfig) error {
    if !cfg.Auth.Enabled && !cfg.Proxy.Enabled && cfg.SSH.Enabled {
        return errors.Errorf("[ERR] config: supply at least one of Auth, SSH or Proxy roles")
    }
    if cfg.Auth.DomainName == "" {
        return errors.Errorf("[ERR] config: please supply domain name")
    }
    if cfg.Hostname == "" {
        return errors.Errorf("[ERR] config: please supply core name")
    }
    if cfg.HostUUID == "" {
        return errors.Errorf("[ERR] config: please supply host UUID")
    }
    if cfg.DataDir == "" {
        return errors.Errorf("[ERR] config: please supply data directory")
    }
    if cfg.Console == nil {
        cfg.Console = ioutil.Discard
    }

/*
    (03/25/2017) TLS keys are not necessary now
    if (cfg.Proxy.TLSKey == "" && cfg.Proxy.TLSCert != "") || (cfg.Proxy.TLSKey != "" && cfg.Proxy.TLSCert == "") {
        return errors.Errorf("please supply both TLS key and certificate")
    }
*/

    // TODO : COMBINE with PCrypto CA Cert issuer
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

    if cfg.BackendDB == nil {
        return errors.Errorf("please provide database engine for backend storage")
    }
    if cfg.CertStorage == nil {
        return errors.Errorf("please provide cert storage")
    }
    if cfg.CaSigner == nil {
        return errors.Errorf("please provide CA signer")
    }
    return nil
}

// Generates a string accepted by the SqliteDB driver, like this:
// `{"path": "/var/lib/teleport/records.db"}`
func sqliteParams(storagePath, dbFile string) string {
    return fmt.Sprintf(`{"path": "%s"}`, filepath.Join(storagePath, dbFile))
}
