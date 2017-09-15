package ivent

// These are the internal events that need to rounted to main packages
const (
    IventNetworkAddressChange    string = "ivent.network.address.change"
    IventBeaconManagerSpawn      string = "ivent.beacon.manager.spawn"
    IventVboxCtrlInstanceSpawn   string = "ivent.vbox.ctrl.instance.spawn"
    IventMonitorRegisteredNode   string = "ivent.monitor.registered.node"
    IventMonitorUnregisteredNode string = "ivent.monitor.unregistered.node"
)
