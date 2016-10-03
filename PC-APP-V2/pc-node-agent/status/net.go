package status

import (
    "net"
    "github.com/stkim1/netifaces"
)

func Interfaces() ([]net.Interface, error) {
    return net.Interfaces()
}

func InterfaceByName(name string) (*net.Interface, error) {
    return net.InterfaceByName(name)
}

func IP4Addrs(iface *net.Interface) ([]*IP4Addr, error) {
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
        case *net.IPAddr:
            if ip4 := v.IP.To4(); ip4 != nil {
                addrs = append(addrs, &IP4Addr{IP:&ip4, IPMask:nil})
            }
        }

    }
    return addrs, nil
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