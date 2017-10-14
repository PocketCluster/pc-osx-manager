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
    IventOrchstInstanceSpawn     string = "ivent.orchst.instance.spawn"
    // external event telling waiters that pcssh proxy instance is up and running
    IventPcsshProxyInstanceSpawn string = "ivent.pcssh.proxy.instance.spawn"


    // request to monitor node status
    IventMonitorNodeReqStatus    string = "ivent.monitor.node.req.status"
    // monitor beacon
    IventMonitorNodeRespBeacon   string = "ivent.monitor.node.resp.beacon"
    // monitor pcssh
    IventMonitorNodeRespPcssh    string = "ivent.monitor.node.resp.pcssh"
    // monitor orchestration
    IventMonitorNodeRespOrchst   string = "ivent.monitor.node.resp.orchst"


    // for package node list
    IventReportNodeListRequest   string = "ivent.report.node.list.request"
    IventReportNodeListResult    string = "ivent.report.node.list.result"
)

// this is to broadcast masterctrl object w/ listener. It's shared with BeaconAgent + VBoxController
type VboxCtrlBrcstObj struct {
    masterctrl.VBoxMasterControl
    net.Listener
}