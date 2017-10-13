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
    IventMonitorNodeBeacon       string = "ivent.monitor.node.beacon"
    IventMonitorNodePcssh        string = "ivent.monitor.node.pcssh"
    IventMonitorNodeOrchst       string = "ivent.monitor.node.orchst"
    IventReportNodeListRequest   string = "ivent.report.node.list.request"
    IventReportNodeListResult    string = "ivent.report.node.list.result"
)

// this is to broadcast masterctrl object w/ listener. It's shared with BeaconAgent + VBoxController
type VboxCtrlBrcstObj struct {
    masterctrl.VBoxMasterControl
    net.Listener
}