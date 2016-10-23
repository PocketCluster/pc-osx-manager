package context

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -Wl,-U,_interface_status_with_callback,-U,_interface_status_with_gocall,-U,_gateway_status_with_callback,-U,_gateway_status_with_gocall  -framework Cocoa

#include "SCNetworkTypes.h"
#include "PCInterfaceTypes.h"

*/
import "C"
import (
    "unsafe"
    "fmt"
)

//export passGatewayArray
func passGatewayArray(gatewayArray **C.SCNIGateway, length C.uint) C.bool {
    var arrayLen int = int(length)
    if arrayLen == 0 || gatewayArray == nil {
        return C.bool(true)
    }

    var gatewaySlice []*C.SCNIGateway = (*[1 << 10]*C.SCNIGateway)(unsafe.Pointer(gatewayArray))[:arrayLen:arrayLen]
    var hostGateways []*HostNetworkGateway = make([]*HostNetworkGateway, arrayLen, arrayLen)
    for idx, gw := range gatewaySlice {
        fmt.Printf("gw address %s\n", C.GoString(gw.addr))
        hostGateways[idx] = &HostNetworkGateway{
            Family      : uint8(gw.family),
            IsDefault   : bool(gw.is_default),
            IfaceName   : C.GoString(gw.ifname),
            Address     : C.GoString(gw.addr),
        }
    }
    singletonContextInstance().monitorNetworkGateways(hostGateways)

    return C.bool(true)
}

func FindSystemGatewayStatus() {
    C.gateway_status_with_gocall()
}


func convertAddressStruct(addrArray **C.SCNIAddress, addrCount C.uint) (addresses []*HostIPAddress) {
    if addrCount == 0 || addrArray == nil {
        addresses = nil
    }

    var addrLength uint = uint(addrCount)
    var addrSlice []*C.SCNIAddress = (*[1 << 10]*C.SCNIAddress)(unsafe.Pointer(addrArray))[:addrLength:addrLength]
    addresses = make([]*HostIPAddress, addrLength, addrLength)
    for idx, addr := range addrSlice {
        addresses[idx] = &HostIPAddress{
            Flags       : uint(addr.flags),
            Family      : uint8(addr.family),
            IsPrimary   : bool(addr.is_primary),
            Address     : C.GoString(addr.addr),
            Netmask     : C.GoString(addr.netmask),
            Broadcast   : C.GoString(addr.broadcast),
            Peer        : C.GoString(addr.peer),
        }
    }
    return
}

//export passInterfaceArray
func passInterfaceArray(interfaceArray **C.PCNetworkInterface, length C.uint) C.bool {
    var arrayLen int = int(length)
    if arrayLen == 0 || interfaceArray == nil {
        return C.bool(true)
    }

    var interfaceSlice []*C.PCNetworkInterface = (*[1 << 10]*C.PCNetworkInterface)(unsafe.Pointer(interfaceArray))[:arrayLen:arrayLen]
    var hostInterfaces []*HostNetworkInterface = make([]*HostNetworkInterface, arrayLen, arrayLen)
    for idx, iface := range interfaceSlice {

        addresses := convertAddressStruct(iface.address, iface.addrCount)

        hostInterfaces[idx] = &HostNetworkInterface{
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
    singletonContextInstance().monitorNetworkInterfaces(hostInterfaces)
    return C.bool(true)
}

func FindSystemInterfaceStatus() {
    C.interface_status_with_gocall()
}
