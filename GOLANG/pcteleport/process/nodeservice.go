package process

import (
    "github.com/gravitational/teleport/lib/service"
    "github.com/stkim1/pcteleport/pcconfig"
)

// NewTeleport takes the daemon configuration, instantiates all required services
// and starts them under a supervisor, returning the supervisor object
func NewNodeTeleport(cfg *pcconfig.Config) (*PocketCoreTeleportProcess, error) {
    var err error
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
