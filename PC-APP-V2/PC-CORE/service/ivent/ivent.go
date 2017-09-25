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
    IventMonitorRegisteredNode   string = "ivent.monitor.registered.node"
    IventMonitorUnregisteredNode string = "ivent.monitor.unregistered.node"
)

// this is to broadcast masterctrl object w/ listener. It's shared with BeaconAgent + VBoxController
type VboxCtrlBrcstObj struct {
    masterctrl.VBoxMasterControl
    net.Listener
}