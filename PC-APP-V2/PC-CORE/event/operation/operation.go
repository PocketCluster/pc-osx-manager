package operation

import (
    "fmt"
)

type CommandType int32

const(
    ServiceBeaconCatcher             string = "service.beacon.catcher"
    ServiceBeaconLocationRead        string = "service.beacon.location.read"
    ServiceBeaconLocationWrite       string = "service.beacon.location.write"
    ServiceBeaconMaster              string = "service.beacon.master"
    ServiceOrchestrationServer       string = "service.orchestration.server"
    ServiceOrchestrationOperation    string = "service.orchestration.operation"
    ServiceStorageProcess            string = "service.storage.process"
    ServiceContainerRegistry         string = "service.container.registry"
    ServiceInternalNodeNameServer    string = "service.internal.node.name.server"
    ServiceInternalNodeNameOperation string = "service.internal.node.name.operation"
    ServiceVBoxMasterControl         string = "service.vbox.master.control"
    ServiceVBoxMasterListener        string = "service.vbox.master.listener"
    ServiceMonitorSystemHealth       string = "service.monitor.system.health"
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