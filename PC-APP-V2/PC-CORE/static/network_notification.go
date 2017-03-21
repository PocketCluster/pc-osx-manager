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
        addresses[idx] = &context.HostIPAddress{
            Flags       : uint(addr.flags),
            Family      : uint8(addr.family),
            IsPrimary   : bool(addr.is_primary),
            Address     : C.GoString(addr.addr),
            Netmask     : C.GoString(addr.netmask),
            Broadcast   : C.GoString(addr.broadcast),
            Peer        : C.GoString(addr.peer),
        }
    }
    return addresses
}

//export NetworkChangeNotificationInterface
func NetworkChangeNotificationInterface(interfaceArray **C.PCNetworkInterface, length C.uint) {
    log.Debugf("NetworkChangeNotificationInterface")
    return

    var arrayLen int = int(length)
    if arrayLen == 0 || interfaceArray == nil {
        return
    }

    var interfaceSlice []*C.PCNetworkInterface = (*[1 << 10]*C.PCNetworkInterface)(unsafe.Pointer(interfaceArray))[:arrayLen:arrayLen]
    var hostInterfaces []*context.HostNetworkInterface = make([]*context.HostNetworkInterface, arrayLen, arrayLen)
    for idx, iface := range interfaceSlice {

        addresses := convertAddressStruct(iface.address, iface.addrCount)

        hostInterfaces[idx] = &context.HostNetworkInterface{
            WifiPowerOff : bool(iface.wifiPowerOff),
            IsActive     : bool(iface.isActive),
            IsPrimary    : bool(iface.isPrimary),
            AddrCount    : uint(iface.addrCount),
            Address      : addresses,
            BsdName      : C.GoString(iface.bsdName),
            DisplayName  : C.GoString(iface.displayName),
            MacAddress   : C.GoString(iface.macAddress),
            MediaType    : C.GoString(iface.mediaType),
        }
    }
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