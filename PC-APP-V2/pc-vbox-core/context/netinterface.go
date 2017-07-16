package context

import (
    "net"
    "strings"

    "github.com/pkg/errors"
    "github.com/stkim1/netifaces"
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
        gateway *netifaces.Gateway    = nil
        ifcFound bool                 = false
        gwaddr, gwiface string
        ipaddrs []string
        ifaces []net.Interface
        priface net.Interface
    )
    gateway, err = netifaces.FindSystemGateways()
    if err != nil {
        return NetworkInterface{}, errors.WithStack(err)
    }
    defer gateway.Release()

    gwaddr, gwiface, err = gateway.DefaultIP4Gateway()
    if err != nil {
        return NetworkInterface{}, errors.WithStack(err)
    }
    if len(gwaddr) == 0 || len(gwiface) == 0 {
        return NetworkInterface{}, errors.Errorf("[ERR] inappropriate gateway address or interface")
    }

    // This loop is to fix wrong interface name on RPI "eth0 + random string" issue
    ifaces, err = net.Interfaces()
    for _, i := range ifaces {
        if strings.HasPrefix(gwiface, i.Name) {
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
        GatewayAddr:     gwaddr,
        // (2017-05-15) We'll only take the first ip address for now
        IP4Address:      ipaddrs,
    }, nil
}

func findNetworkInterface(iName string) (NetworkInterface, error) {
    var (
        err error                     = nil
        gw *netifaces.Gateway         = nil
        ifcFound bool                 = false
        gwaddr string
        ipaddrs []string
        ifaces []net.Interface
        netIface net.Interface
        netStatus NetworkInterface    = NetworkInterface{}
    )
    if len(iName) == 0 {
        return netStatus, errors.Errorf("[ERR] inappropriate interface name")
    }


    // --- find network interface hardware status ---
    ifaces, err = net.Interfaces()
    for _, i := range ifaces {
        if i.Name == iName {
            netIface = i
            ifcFound = true
            break
        }
    }
    if !ifcFound {
        return netStatus, errors.Errorf("[ERR] invalid interface name to find")
    }
    netStatus.Interface    = netIface
    netStatus.HardwareAddr = netIface.HardwareAddr.String()


    // --- find ip address status ---
    ipaddrs, err = ip4AddrsToStringList(netIface)
    if err != nil {
        return netStatus, errors.WithStack(err)
    }
    netStatus.IP4Address   = ipaddrs


    // --- find gateway information ---
    gw, err = netifaces.FindSystemGateways()
    if err != nil {
        return netStatus, errors.WithStack(err)
    }
    defer gw.Release()

    gwaddr, err = gw.FindGatewayForInterface(iName)
    if err != nil {
        return netStatus, errors.WithStack(err)
    }
    if len(gwaddr) != 0 {
        netStatus.GatewayAddr = gwaddr
    }

    // --- results ---
    return netStatus, nil
}

func FirstNetworkInterface() (NetworkInterface, error) {
    return findNetworkInterface("eth0")
}

func SecondNetworkInterface() (NetworkInterface, error) {
    return findNetworkInterface("eth1")
}