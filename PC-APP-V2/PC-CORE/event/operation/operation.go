package operation

import (
    "fmt"
)

type CommandType int32

const (
    // Context control : This opens/closes database
    CmdContextOpen      = iota
    CmdContextClose

    // Beacon control
    CmdBeaconStart
    CmdBeaconStop

    // Teleport control
    CmdTeleportStart
    CmdTeleportStop

    // ETCD control
    CmdStorageStart
    CmdStroageStop

    // swarm control
    CmdCntrOrchStart
    CmdCntrOrchStop

    // Registry control
    CmdImageRegistryStart
    CmdImageRegistryStop
)

func (c CommandType) String() string {
    switch c {
        case CmdContextOpen:
            return "CmdContextOpen"
        case CmdContextClose:
            return "CmdContextClose"
        case CmdBeaconStart:
            return "CmdBeaconStart"
        case CmdBeaconStop:
            return "CmdBeaconStop"
        case CmdTeleportStart:
            return "CmdTeleportStart"
        case CmdTeleportStop:
            return "CmdTeleportStop"
        case CmdStorageStart:
            return "CmdStorageStart"
        case CmdStroageStop:
            return "CmdStroageStop"
        case CmdCntrOrchStart:
            return "CmdCntrOrchStart"
        case CmdCntrOrchStop:
            return "CmdCntrOrchStop"
        case CmdImageRegistryStart:
            return "CmdImageRegistryStart"
        case CmdImageRegistryStop:
            return "CmdImageRegistryStop"
        default:
            return fmt.Sprintf("CommandType(%d)", c)
    }
}

type Operation struct {
    Command    CommandType
}

func (o *Operation) String() string {
    return o.Command.String()
}