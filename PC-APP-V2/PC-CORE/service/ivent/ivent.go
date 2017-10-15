package ivent

/*
 * This package contains internal message event name and struct. Do not import other packages
 */


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
    VBoxMasterControl interface{}
    Listener          interface{}
}

// node status info from orchestration
type EngineStatusInfo struct {
    Name          string
    ID            string
    IP            string
    Addr          string
}

type EngineStatusMeta struct {
    TimeStamp     int64
    Error         error
    Engines       []EngineStatusInfo
}

// node status info from Pcssh
type PcsshNodeStatusInfo struct {
    HostName      string
    ID            string
    Addr          string
    HasSession    bool
}

type PcsshNodeStatusMeta struct {
    TimeStamp     int64
    Error         error
    Nodes         []PcsshNodeStatusInfo
}
