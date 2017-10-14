package ivent

import (
    "net"

    "github.com/stkim1/pc-vbox-comm/masterctrl"
)

// These are the internal events that need to rounted to main packages
const (
    IventNetworkAddressChange    string = "ivent.network.address.change"
    IventBeaconManagerSpawn      string = "ivent.beacon.manager.spawn"
    IventVboxCtrlInstanceSpawn   string = "ivent.vbox.ctrl.instance.spawn"
    IventReportNodeListRequest   string = "ivent.report.node.list.request"
    IventReportNodeListResult    string = "ivent.report.node.list.result"

    // monitor beacon
    IventMonitorNodeReqBeacon    string = "ivent.monitor.node.req.beacon"
    IventMonitorNodeRsltBeacon   string = "ivent.monitor.node.rslt.beacon"

    // pcssh monitor ivent is in its own package

    // monitor orchestration
    IventMonitorNodeReqOrchst    string = "ivent.monitor.node.req.orchst"
    IventMonitorNodeRsltOrchst   string = "ivent.monitor.node.rslt.orchst"
)

// this is to broadcast masterctrl object w/ listener. It's shared with BeaconAgent + VBoxController
type VboxCtrlBrcstObj struct {
    masterctrl.VBoxMasterControl
    net.Listener
}