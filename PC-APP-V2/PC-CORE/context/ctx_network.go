package context

import (
    "github.com/pkg/errors"
)

type HostContextNetwork interface {
    HostPrimaryAddress() (string, error)
    HostDefaultGatewayAddress() (string, error)
    RefreshNetworkInterfaces(interfaces []*HostNetworkInterface)
    RefreshNetworkGateways(gateways []*HostNetworkGateway)
}

type HostIPAddress struct {
    Flags               uint
    Family              uint8
    IsPrimary           bool
    Address             string
    Netmask             string
    Broadcast           string
    Peer                string
}

type HostNetworkInterface struct {
    WifiPowerOff        bool
    IsActive            bool
    IsPrimary           bool
    AddrCount           uint
    Address             []*HostIPAddress
    BsdName             string
    DisplayName         string
    MacAddress          string
    MediaType           string
}

type HostNetworkGateway struct {
    Family              uint8
    IsDefault           bool
    IfaceName           string
    Address             string
}

type hostNetwork struct {
    hostInterfaces               *[]*HostNetworkInterface
    hostGateways                 *[]*HostNetworkGateway

    primaryInteface              *HostNetworkInterface
    primaryAddress               *HostIPAddress
    primaryGateway               *HostNetworkGateway
}

// take network interfaces
func (ctx *hostContext) RefreshNetworkInterfaces(interfaces []*HostNetworkInterface) {
    ctx.Lock()
    defer ctx.Unlock()

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

func (ctx *hostContext) RefreshNetworkGateways(gateways []*HostNetworkGateway) {
    ctx.Lock()
    defer ctx.Unlock()

    for _, gw := range gateways {
        if gw.IsDefault {
            ctx.primaryGateway = gw
        }
    }
    return
}

func (ctx *hostContext) HostPrimaryAddress() (string, error) {
    addr := ctx.primaryAddress
    if addr != nil {
        return addr.Address, nil
    }

    return "", errors.Errorf("[ERR] No primary address has been found")
}

func (ctx *hostContext) HostDefaultGatewayAddress() (string, error) {
    gateway := ctx.primaryGateway
    if gateway != nil {
        return gateway.Address, nil
    }

    return "", errors.Errorf("[ERR] No default gateway is found")
}