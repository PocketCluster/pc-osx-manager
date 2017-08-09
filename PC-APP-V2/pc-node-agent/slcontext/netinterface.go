package slcontext

import (
    "net"
    "strings"

    "github.com/pkg/errors"
    "github.com/stkim1/findgate"
)

type NetworkInterface struct {
    net.Interface
    HardwareAddr    string
    GatewayAddr     string
    IP4Address      []string
}

func (n *NetworkInterface) PrimaryIP4Addr() string {
    // (2017-05-15) We'll only take the first ip address for now
    return n.IP4Address[0]
}

// --- Network ---
type ip4addr struct {
    *net.IP
    *net.IPMask
}

func ip4Address(iface *net.Interface) ([]*ip4addr, error) {
    ifAddrs, err := iface.Addrs()
    if err != nil {
        return nil, err
    }

    var addrs []*ip4addr
    for _, addr := range ifAddrs {
        switch v := addr.(type) {
        case *net.IPNet:
            if ip4 := v.IP.To4(); ip4 != nil {
                addrs = append(addrs, &ip4addr{IP:&ip4, IPMask:&v.Mask})
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
        return nil, errors.Errorf("[ERR] No IPv4 address is given to interface %s", iface.Name);
    }
    return addrs, nil
}

func ip4AddrsToStringList(iface net.Interface) ([]string, error) {
    addresses, err := iface.Addrs()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    var addrstr []string
    for _, addr := range addresses {
        switch v := addr.(type) {
        case *net.IPNet:
            if ip4 := v.IP.To4(); ip4 != nil {
                addrstr = append(addrstr, addr.String())
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
    if len(addrstr) == 0 {
        return nil, errors.Errorf("[ERR] no IPv4 address is given to interface %s", iface.Name);
    }
    return addrstr, nil
}

func PrimaryNetworkInterface() (NetworkInterface, error) {
    var (
        err error                     = nil
        gateway *findgate.IPv4Gateway = nil
        ifcFound bool                 = false
        ipaddrs []string
        ifaces []net.Interface
        priface net.Interface
    )

    gateway, err = findgate.DefaultIPv4Gateway()
    if err != nil {
        return NetworkInterface{}, errors.WithStack(err)
    }

    // This loop is to fix wrong interface name on RPI "eth0 + random string" issue
    ifaces, err = net.Interfaces()
    for _, i := range ifaces {
        if strings.HasPrefix(gateway.Interface, i.Name) {
            priface = i
            ifcFound = true
            break
        }
    }
    if !ifcFound {
        return NetworkInterface{}, errors.Errorf("[ERR] primary interface is not found")
    }

    ipaddrs, err = ip4AddrsToStringList(priface)
    if err != nil {
        return NetworkInterface{}, errors.WithStack(err)
    }

    return NetworkInterface {
        Interface:       priface,
        HardwareAddr:    priface.HardwareAddr.String(),
        GatewayAddr:     gateway.Address,
        IP4Address:      ipaddrs,
    }, nil
}
