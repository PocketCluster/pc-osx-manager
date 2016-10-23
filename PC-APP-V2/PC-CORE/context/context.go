package context

import "sync"

type HostContext interface {
}

type hostContext struct {
    hostInterfaces      *[]*HostNetworkInterface
    hostGateways        *[]*HostNetworkGateway

    primaryInteface     *HostNetworkInterface
    primaryAddress      *HostIPAddress
}

// singleton initialization
var _context *hostContext = nil
var _once sync.Once

func SharedHostContext() (HostContext) {
    return singletonContextInstance()
}

func singletonContextInstance() (*hostContext) {
    _once.Do(func() {
        _context = &hostContext{}
        initializeHostContext(_context)
    })
    return _context
}

func initializeHostContext(ctx *hostContext) {
}

// take network interfaces
func (ctx *hostContext) monitorNetworkInterfaces(interfaces []*HostNetworkInterface) {
    // TODO : we make an assumption that host's primary interface and network addresses are at the same network segment. This could not be the case, we'll look into it v0.1.5
    ctx.hostInterfaces = &interfaces

    for _, iface := range interfaces {
        if iface.IsPrimary {
            ctx.primaryInteface = iface

            for _, addr := range iface.Address {
                if addr.IsPrimary {
                    ctx.primaryAddress = addr
                }
            }
            break
        }
    }

    // this is backup. It selects 1) Wi-Fi interface with 2) an active ip address.
    // Ethernet is going to be obsolete on lots of new macbook. We'll take wifi as default
    if ctx.primaryInteface == nil {
        for _, iface := range interfaces {
            if iface.IsActive && (iface.MediaType == "IEEE80211" || iface.DisplayName == "Wi-Fi") {
                ctx.primaryInteface = iface
                ctx.primaryAddress = iface.Address[0]
                break
            }
        }
    }
}

func (ctx *hostContext) monitorNetworkGateways(gateways []*HostNetworkGateway) {

}

