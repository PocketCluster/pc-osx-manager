package context

import (
    "reflect"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/davecgh/go-spew/spew"
)

type HostContextNetwork interface {
    HostPrimaryAddress() (string, error)
    HostPrimaryInterfaceShortName() (string, error)
    HostPrimaryInterfaceFullName() (string, error)
    HostDefaultGatewayAddress() (string, error)
    UpdateNetworkInterfaces(interfaces []*HostNetworkInterface) bool
    UpdateNetworkGateways(gateways []*HostNetworkGateway)
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

func (ctx *hostContext) HostPrimaryAddress() (string, error) {
    ctx.Lock()
    defer ctx.Unlock()

    addr := ctx.primaryAddress
    if addr == nil {
        return "", errors.Errorf("[ERR] no primary address has been found")
    }

    return addr.Address, nil
}

func (ctx *hostContext) HostPrimaryInterfaceShortName() (string, error) {
    ctx.Lock()
    defer ctx.Unlock()

    iface := ctx.primaryInteface
    if iface == nil {
        return "", errors.Errorf("[ERR] no primary interface has been found")
    }
    return iface.BsdName, nil
}

func (ctx *hostContext) HostPrimaryInterfaceFullName() (string, error) {
    ctx.Lock()
    defer ctx.Unlock()

    iface := ctx.primaryInteface
    if iface == nil {
        return "", errors.Errorf("[ERR] no primary interface has been found")
    }
    return iface.DisplayName, nil
}

func (ctx *hostContext) HostDefaultGatewayAddress() (string, error) {
    ctx.Lock()
    defer ctx.Unlock()

    gateway := ctx.primaryGateway
    if gateway != nil {
        return gateway.Address, nil
    }

    return "", errors.Errorf("[ERR] No default gateway is found")
}

// take network interfaces
func (ctx *hostContext) UpdateNetworkInterfaces(interfaces []*HostNetworkInterface) bool {
    ctx.Lock()
    defer ctx.Unlock()

    var (
        addr     *HostIPAddress = ctx.primaryAddress
        adrFound bool           = false
    )

    //log.Debugf(spew.Sdump(interfaces))
/*
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
*/

    // It selects 1) Wi-Fi interface with 2) an active ip address.
    // Ethernet is going to be obsolete on lots of new macbook. We'll take wifi as default
    for _, iface := range interfaces {
        if iface.IsActive && (iface.MediaType == "IEEE80211" || iface.DisplayName == "Wi-Fi") {
            ctx.primaryInteface = iface
            for _, addr := range iface.Address {
                if addr.IsPrimary {
                    ctx.primaryAddress = addr
                    adrFound = true
                }
            }
            if !adrFound && 0 < len(iface.Address) {
                ctx.primaryAddress = iface.Address[0]
                adrFound = true
            }
            break
        }
    }

    // no address found
    if !adrFound {
        ctx.primaryAddress = nil
    }

    log.Debugf(spew.Sdump(ctx.primaryAddress))
    log.Debugf(spew.Sdump(ctx.primaryInteface))

    return !reflect.DeepEqual(addr, ctx.primaryAddress)
}

func (ctx *hostContext) UpdateNetworkGateways(gateways []*HostNetworkGateway) {
    ctx.Lock()
    defer ctx.Unlock()

    //log.Debugf(spew.Sdump(gateways))
    for _, gw := range gateways {
        if gw.IsDefault {
            ctx.primaryGateway = gw
        }
    }
    return
}
