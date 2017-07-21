package crcontext

import (
    "net"
    "strings"

    "github.com/pkg/errors"
    "github.com/stkim1/findgate"
)

const (
    InternalNetworkDevice string = "eth0"
    ExternalNetworkDevice string = "eth1"
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
        // (2017-05-15) We'll only take the first ip address for now
        IP4Address:      ipaddrs,
    }, nil
}

func findNetworkInterface(iName string) (NetworkInterface, error) {
    var (
        err error                     = nil
        ifcFound bool                 = false
        gateways []findgate.IPv4Gateway
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
    gateways, err = findgate.FindIPv4GatewayWithInterface(iName)
    if err != nil {
        return netStatus, errors.WithStack(err)
    }
    // TODO : find the default gateway for this interface. Most likely, the first gateway is the default gateway for the given interface
    netStatus.GatewayAddr = gateways[0].Address

    // --- results ---
    return netStatus, nil
}

func InternalNetworkInterface() (NetworkInterface, error) {
    return findNetworkInterface(InternalNetworkDevice)
}

func ExternalNetworkInterface() (NetworkInterface, error) {
    return findNetworkInterface(ExternalNetworkDevice)
}