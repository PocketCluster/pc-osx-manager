package status

import (
    "net"
    "fmt"
    "github.com/stkim1/netifaces"
)

type Interface struct {
    *net.Interface
}

func InterfaceByName(name string) (*Interface, error) {
    iface, err := net.InterfaceByName(name); if err != nil {
        return nil, err
    }
    return &Interface{iface}, nil
}

func (iface *Interface) IP4Addrs() ([]*IP4Addr, error) {
    ifAddrs, err := iface.Addrs()
    if err != nil {
        return nil, err
    }

    var addrs []*IP4Addr
    for _, addr := range ifAddrs {
        switch v := addr.(type) {
        case *net.IPNet:
            if ip4 := v.IP.To4(); ip4 != nil {
                addrs = append(addrs, &IP4Addr{IP:&ip4, IPMask:&v.Mask})
            }
        // TODO : make sure net.IPAddr only represents IP6
        /*
        case *net.IPAddr:
            if ip4 := v.IP.To4(); ip4 != nil {
                addrs = append(addrs, &IP4Addr{IP:&ip4, IPMask:nil})
            }
        */
        }
    }
    if len(addrs) == 0 {
        return nil, fmt.Errorf("[ERR] No IPv4 address is given to interface %s", iface.Name);
    }
    return addrs, nil
}

func (iface *Interface) MacAddress() string {
    return iface.HardwareAddr.String()
}

type IP4Addr struct {
    *net.IP
    *net.IPMask
}

func (a *IP4Addr) IPString() string {
    return a.IP.String()
}

func (a *IP4Addr) IPMaskString() string {
    return a.IPMask.String()
}

func GetDefaultIP4Gateway() (address string, iface string, err error) {
    gw, err := netifaces.FindSystemGateways(); if err != nil {
        return "", "", err
    }
    defer gw.Release()
    address, iface, err = gw.DefaultIP4Gateway(); if err != nil {
        return "", "", err
    }
    return
}