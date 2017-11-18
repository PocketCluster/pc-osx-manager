package ivent

/*
 * This package contains internal message event name and struct. Do not import other packages
 */


// These are the internal events that need to rounted to main packages
const (
    // network address change event
    IventNetworkAddressChange    string = "ivent.network.address.change"


    // internal service spawn list
    IventDiscoveryInstanceSpwan  string = "ivent.dicovery.instance.spawn"
    IventRegistryInstanceSpawn   string = "ivent.registry.instance.spawn"
    IventNameServerInstanceSpawn string = "ivent.name.server.instance.spawn"
    IventBeaconManagerSpawn      string = "ivent.beacon.manager.spawn"
    IventPcsshProxyInstanceSpawn string = "ivent.pcssh.proxy.instance.spawn"
    IventOrchstInstanceSpawn     string = "ivent.orchst.instance.spawn"
    IventVboxCtrlInstanceSpawn   string = "ivent.vbox.ctrl.instance.spawn"

    // internal service spawn error for health monitor.
    // Only to be fired from "github.com/stkim1/pc-core/static/main.go" or initialization
    IventInternalSpawnError      string = "ivent.internal.spawn.error"


    // request to monitor node status
    IventMonitorNodeReqStatus    string = "ivent.monitor.node.req.status"
    // monitor beacon
    IventMonitorNodeRespBeacon   string = "ivent.monitor.node.resp.beacon"
    // monitor pcssh
    IventMonitorNodeRespPcssh    string = "ivent.monitor.node.resp.pcssh"
    // monitor orchestration
    IventMonitorNodeRespOrchst   string = "ivent.monitor.node.resp.orchst"
    // stop monitor process reqeust
    IventMonitorStopRequest      string = "ivent.monitor.stop.request"
    // monitor stopped
    IventMonitorStopResult       string = "ivent.monitor.stop.result"

    // for package node list
    IventReportLiveNodesRequest  string = "ivent.report.live.nodes.request"
    IventReportLiveNodesResult   string = "ivent.report.live.nodes.result"
)

// this is to broadcast masterctrl object w/ listener. It's shared with BeaconAgent + VBoxController
type VboxCtrlBrcstObj struct {
    VBoxMasterControl interface{}
    Listener          interface{}
}

// node status info from beacon
type BeaconNodeStatusInfo struct {
    Name          string
    MacAddr       string
    IPAddr        string
    Registered    bool
    Bounded       bool
}

type BeaconNodeStatusMeta struct {
    TimeStamp     int64
    Error         error
    Nodes         []BeaconNodeStatusInfo
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

