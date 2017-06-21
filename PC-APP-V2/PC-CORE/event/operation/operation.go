package operation

import (
    "fmt"
)

type CommandType int32

const(
    ServiceBeaconCatcher string          = "service.beacon.catcher"
    ServiceBeaconLocationRead string     = "service.beacon.location.read"
    ServiceBeaconLocationWrite string    = "service.beacon.location.write"
    ServiceBeaconMaster string           = "service.beacon.master"
    ServiceSwarmEmbeddedServer string    = "service.swarm.embedded.server"
    ServiceSwarmEmbeddedOperation string = "service.swarm.embedded.operation"
    ServiceStorageProcess string         = "service.storage.process"
    ServiceContainerRegistry string      = "service.container.registry"
)

const (
    // Base Service start
    CmdBaseServiceStart     = iota
    CmdBaseServiceStop

    // ETCD control
    CmdStorageStart
    CmdStorageStop

    // debug add node, root, & user
    CmdTeleportRootAdd
    CmdTeleportUserAdd

    // Debug control
    CmdDebug
)

func (c CommandType) String() string {
    switch c {
        case CmdBaseServiceStart:
            return "CmdBaseServiceStart"
        case CmdBaseServiceStop:
            return "CmdBaseServiceStop"
        case CmdStorageStart:
            return "CmdStorageStart"
        case CmdStorageStop:
            return "CmdStroageStop"
        case CmdTeleportRootAdd:
            return "CmdTeleportRootAdd"
        case CmdTeleportUserAdd:
            return "CmdTeleportUserAdd"
        case CmdDebug:
            return "CmdDebug"

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