// +build darwin
package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_interface_status_with_callback,-U,_gateway_status_with_callback

#include "SCNetworkTypes.h"
#include "PCInterfaceTypes.h"

*/
import "C"
import (
    "unsafe"

    "github.com/stkim1/pc-core/context"
    log "github.com/Sirupsen/logrus"
    "github.com/davecgh/go-spew/spew"
)

func convertAddressStruct(addrArray **C.SCNIAddress, addrCount C.uint) ([]*context.HostIPAddress) {
    var (
        addrLength uint                    = uint(addrCount)
        addrSlice []*C.SCNIAddress         = (*[1 << 10]*C.SCNIAddress)(unsafe.Pointer(addrArray))[:addrLength:addrLength]
        addresses []*context.HostIPAddress = make([]*context.HostIPAddress, addrLength, addrLength)
    )
    for idx, addr := range addrSlice {
        var (
            address, netmask, broadcast, peer = "", "", "", ""
        )
        if addr.addr != nil {
            address = C.GoString(addr.addr)
        }
        if addr.netmask != nil {
            netmask = C.GoString(addr.netmask)
        }
        if addr.broadcast != nil {
            broadcast = C.GoString(addr.broadcast)
        }
        if addr.peer != nil {
            peer = C.GoString(addr.peer)
        }
        addresses[idx] = &context.HostIPAddress{
            Flags:        uint(addr.flags),
            Family:       uint8(addr.family),
            IsPrimary:    bool(addr.is_primary),
            Address:      address,
            Netmask:      netmask,
            Broadcast:    broadcast,
            Peer:         peer,
        }
    }
    return addresses
}

//export NetworkChangeNotificationInterface
func NetworkChangeNotificationInterface(interfaceArray **C.PCNetworkInterface, length C.uint) {
    var arrayLen int = int(length)
    if arrayLen == 0 || interfaceArray == nil {
        return
    }
    var (
        interfaceSlice []*C.PCNetworkInterface = (*[1 << 10]*C.PCNetworkInterface)(unsafe.Pointer(interfaceArray))[:arrayLen:arrayLen]
        hostInterfaces []*context.HostNetworkInterface = make([]*context.HostNetworkInterface, arrayLen, arrayLen)
    )
    for idx, iface := range interfaceSlice {
        var (
            addresses []*context.HostIPAddress = nil
            bsdName, displayName, macAddress, mediaType = "", "", "", ""
        )
        if iface.address != nil && 0 < uint(iface.addrCount) {
            addresses = convertAddressStruct(iface.address, iface.addrCount)
        }
        if iface.bsdName != nil {
            bsdName = C.GoString(iface.bsdName)
        }
        if iface.displayName != nil {
            displayName = C.GoString(iface.displayName)
        }
        if iface.macAddress != nil {
            macAddress = C.GoString(iface.macAddress)
        }
        if iface.mediaType != nil {
            mediaType = C.GoString(iface.mediaType)
        }
        hostInterfaces[idx] = &context.HostNetworkInterface{
            WifiPowerOff : bool(iface.wifiPowerOff),
            IsActive     : bool(iface.isActive),
            IsPrimary    : bool(iface.isPrimary),
            AddrCount    : uint(iface.addrCount),
            Address      : addresses,
            BsdName      : bsdName,
            DisplayName  : displayName,
            MacAddress   : macAddress,
            MediaType    : mediaType,
        }
    }
    log.Debugf(spew.Sdump(hostInterfaces))
    context.MonitorNetworkInterfaces(hostInterfaces)
}

//export NetworkChangeNotificationGateway
func NetworkChangeNotificationGateway(gatewayArray **C.SCNIGateway, length C.uint) {
    var arrayLen int = int(length)
    if arrayLen == 0 || gatewayArray == nil {
        return
    }
    var (
        gatewaySlice []*C.SCNIGateway = (*[1 << 10]*C.SCNIGateway)(unsafe.Pointer(gatewayArray))[:arrayLen:arrayLen]
        hostGateways []*context.HostNetworkGateway = make([]*context.HostNetworkGateway, arrayLen, arrayLen)
    )
    for idx, gw := range gatewaySlice {
        hostGateways[idx] = &context.HostNetworkGateway{
            Family:       uint8(gw.family),
            IsDefault:    bool(gw.is_default),
            IfaceName:    C.GoString(gw.ifname),
            Address:      C.GoString(gw.addr),
        }
    }
    log.Debugf(spew.Sdump(hostGateways))
    context.MonitorNetworkGateways(hostGateways)
}