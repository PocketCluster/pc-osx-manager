package netifaces
/*
// #include <string.h>
#include <stdlib.h>
#include "netifaces.h"
*/
import "C"
import (
    "syscall"
    "bytes"
//    "unsafe"

    "github.com/pkg/errors"
)

type Gateway struct {
    gateway *C.Gateway
}

// this filters . : a-z 0-9 only
func filteredInterfaceString(input []byte) string {
    var buf bytes.Buffer
    for _, c := range []byte(input) {
        var char int = int(c)
        switch {
            case 48 <= char && char <= 57:
                fallthrough
            case 97 <= char && char <= 122:
                buf.WriteByte(c)
        }
    }
    return buf.String()
}

func filteredIp4String(input []byte) string {
    var buf bytes.Buffer
    for _, c := range []byte(input) {
        var char int = int(c)
        switch {
            case 46 == char:
                fallthrough
            case 48 <= char && char <= 57:
                buf.WriteByte(c)
        }
    }
    return buf.String()
}

// Find all gateways in system
func FindSystemGateways() (*Gateway, error) {
    var gw *Gateway = &Gateway{}
    syserr := C.find_system_gateways(&gw.gateway)
    if syserr != 0 {
        return nil, errors.Errorf("[ERR] Cannot find all system gateways %s", syscall.Errno(syserr).Error())
    }
    return gw, nil
}

// Release all search results
func (g *Gateway) Release() {
    C.release_gateways_info(&g.gateway)
}

// Find the first ip4 default gateway
func (g *Gateway) DefaultIP4Gateway() (string, string, error) {
    var (
        address, ifname string = "", ""
        gw *C.Gateway = nil
    )

    gw = C.find_default_ip4_gw(&g.gateway)
    if gw == nil {
        return "", "", errors.Errorf("[ERR] Cannot find default gateway for IPv4")
    }

    // FIXME : C-String to Go-String conversion with strlen() has cuased erraneous, extra character addition.
    //address = filteredIp4String(C.GoBytes(unsafe.Pointer(gw.addr), C.int(C.strlen(gw.addr))))
    //ifname = filteredInterfaceString(C.GoBytes(unsafe.Pointer(gw.ifname), C.int(C.strlen(gw.ifname))))

    address = filteredIp4String([]byte(C.GoString(gw.addr)))
    ifname = filteredInterfaceString([]byte(C.GoString(gw.ifname)))
    return address, ifname, nil
}

func (g *Gateway) FindGatewayForInterface(iName string) (string, error) {
    var (
        address string = ""
        cIname *C.char = nil
        gw *C.Gateway = nil
    )
    if len(iName) == 0 {
        return "", errors.Errorf("[ERR] invalid interface name for gateway")
    }
    cIname = C.CString(iName)
    defer C.free(cIname)

    gw = C.find_ip4_gw_for_interface(&g.gateway, cIname)
    if gw == nil {
        return "", errors.Errorf("[ERR] Cannot find IPv4 gateway for %s", iName)
    }
    address = filteredIp4String([]byte(C.GoString(gw.addr)))
    return address, nil
}