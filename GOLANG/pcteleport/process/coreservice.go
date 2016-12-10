package process

import (
    "os"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/auth/native"

    "github.com/stkim1/pcteleport/pcconfig"
    "github.com/pborman/uuid"
)

// NewTeleport takes the daemon configuration, instantiates all required services
// and starts them under a supervisor, returning the supervisor object
func NewCoreTeleport(cfg *pcconfig.Config) (*PocketCoreTeleportProcess, error) {
    if err := pcconfig.ValidateCoreConfig(cfg); err != nil {
        return nil, trace.Wrap(err, "Configuration error")
    }

    // create the data directory if it's missing
    _, err := os.Stat(cfg.DataDir)
    if os.IsNotExist(err) {
        err := os.MkdirAll(cfg.DataDir, os.ModeDir|0700)
        if err != nil {
            return nil, trace.Wrap(err)
        }
    }

    // if there's no host uuid initialized yet, try to read one from the
    // one of the identities
    cfg.HostUUID, err = utils.ReadHostUUID(cfg.DataDir)
    if err != nil {
/*
        TODO : need to look into IsNotFound Error to see what really happens
        if !trace.IsNotFound(err) {
            return nil, trace.Wrap(err)
        }
*/
        if len(cfg.Identities) != 0 {
            cfg.HostUUID = cfg.Identities[0].ID.HostUUID
            log.Infof("[INIT] taking host uuid from first identity: %v", cfg.HostUUID)
        } else {
            cfg.HostUUID = uuid.New()
            log.Infof("[INIT] generating new host UUID: %v", cfg.HostUUID)
        }
        if err := utils.WriteHostUUID(cfg.DataDir, cfg.HostUUID); err != nil {
            return nil, trace.Wrap(err)
        }
    }

    // if user started auth and another service (without providing the auth address for
    // that service, the address of the in-process auth will be used
    if cfg.Auth.Enabled && len(cfg.AuthServers) == 0 {
        cfg.AuthServers = []utils.NetAddr{cfg.Auth.SSHAddr}
    }

    // if user did not provide auth domain name, use this host UUID
    if cfg.Auth.Enabled && cfg.Auth.DomainName == "" {
        cfg.Auth.DomainName = cfg.HostUUID
    }

    // try to login into the auth service:

    // if there are no certificates, use self signed
    process := &PocketCoreTeleportProcess{
        Supervisor: service.NewSupervisor(),
        Config:     cfg,
    }

    serviceStarted := false

    if cfg.Auth.Enabled {
        if cfg.Keygen == nil {
            cfg.Keygen = native.New()
        }
        if err := process.initAuthService(cfg.Keygen); err != nil {
            return nil, trace.Wrap(err)
        }
        serviceStarted = true
    }

    if cfg.Proxy.Enabled {
        if err := process.initProxy(); err != nil {
            return nil, err
        }
        serviceStarted = true
    }

    if !serviceStarted {
        return nil, trace.Errorf("all services failed to start")
    }

    return process, nil
}
