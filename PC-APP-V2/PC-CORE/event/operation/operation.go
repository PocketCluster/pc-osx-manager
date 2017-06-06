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
    CmdStorageStop

    // Registry control
    CmdRegistryStart
    CmdRegistryStop

    // debug bundle start & stop
    CmdServiceBundleStart
    CmdServiceBundleStop

    // debug add node, root, & user
    CmdTeleportNodeAdd
    CmdTeleportRootAdd
    CmdTeleportUserAdd
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
        case CmdStorageStop:
            return "CmdStroageStop"
        case CmdRegistryStart:
            return "CmdImageRegistryStart"
        case CmdRegistryStop:
            return "CmdImageRegistryStop"
        case CmdServiceBundleStart:
            return "CmdServiceBundleStart"
        case CmdServiceBundleStop:
            return "CmdServiceBundleStop"
        case CmdTeleportNodeAdd:
            return "CmdTeleportNodeAdd"
        case CmdTeleportRootAdd:
            return "CmdTeleportRootAdd"
        case CmdTeleportUserAdd:
            return "CmdTeleportUserAdd"

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