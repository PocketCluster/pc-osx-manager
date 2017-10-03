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
    CmdDebug0
    CmdDebug1
    CmdDebug2
    CmdDebug3
    CmdDebug4
    CmdDebug5
    CmdDebug6
    CmdDebug7
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
        case CmdDebug0:
            return "CmdDebug0"
        case CmdDebug1:
            return "CmdDebug1"
        case CmdDebug2:
            return "CmdDebug2"
        case CmdDebug3:
            return "CmdDebug3"
        case CmdDebug4:
            return "CmdDebug4"
        case CmdDebug5:
            return "CmdDebug5"
        case CmdDebug6:
            return "CmdDebug6"
        case CmdDebug7:
            return "CmdDebug7"

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