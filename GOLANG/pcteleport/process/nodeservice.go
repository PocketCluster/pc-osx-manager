package process

import (
    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/stkim1/pcteleport/pcconfig"
    "github.com/pborman/uuid"
)

// NewTeleport takes the daemon configuration, instantiates all required services
// and starts them under a supervisor, returning the supervisor object
func NewNodeTeleport(cfg *pcconfig.Config) (*PocketCoreTeleportProcess, error) {
    var err error
    // if there's no host uuid initialized yet, try to read one from the
    // one of the identities
    cfg.HostUUID, err = utils.ReadHostUUID(cfg.DataDir)
    if err != nil {
        if !trace.IsNotFound(err) {
            return nil, trace.Wrap(err)
        }
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

    // if user did not provide auth domain name, use this host UUID
    if cfg.Auth.Enabled && cfg.Auth.DomainName == "" {
        cfg.Auth.DomainName = cfg.HostUUID
    }

    process := &PocketCoreTeleportProcess{
        Supervisor: service.NewSupervisor(),
        Config:     cfg,
    }

    err = process.initSSH();
    if err != nil {
        return nil, err
    }

    return process, nil
}
