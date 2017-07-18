package findgate

import (
    "bufio"
    "encoding/binary"
    "fmt"
    "os"
    "strings"
    "strconv"
    "net"

    "github.com/pkg/errors"
)

const (
    route_tbl_path string = "/proc/net/route"
    line_separater string = "\t"                // field separator
)

/*
 * Iface    Destination   Gateway    Flags   RefCnt   Use   Metric   Mask       MTU   Window   IRTT
 * enp0s3   00000000      0202000A   0003    0        0     0        00000000   0     0        0
 * enp0s3   0002000A      00000000   0001    0        0     0        00FFFFFF   0     0        0
 * enp0s8   0001A8C0      00000000   0001    0        0     0        00FFFFFF   0     0        0
 */
const (
    field_interface     int = iota
    field_destination
    field_gateway                               // field containing hex gateway address
    field_flags
    field_ref_count
    field_use
    field_metric
    field_mask
    field_mtu
    field_window
    field_irtt
)

type Gateway struct {
    AddrMask     net.IPNet
    Address      string
    Mask         string
    Interface    string
}

func hexAddressToIPv4(gwHexAddr string) (net.IP, string, error) {
    var (
        gwIP net.IP = make(net.IP, net.IPv4len)
        err error = nil
        hexAddr int64
    )

    // cast hex address to uint32
    hexAddr, err = strconv.ParseInt(fmt.Sprintf("0x%s",gwHexAddr), 0, 64)
    if err != nil {
        return gwIP, "", errors.WithStack(err)
    }

    // make net.IP address from uint32
    binary.LittleEndian.PutUint32(gwIP, uint32(hexAddr))

    // format net.IP to dotted ipV4 string
    return gwIP, net.IP(gwIP).String(), nil
}

func hexMaskToIPv4(gwHexMask string) (net.IPMask, string, error) {
    var (
        gwMask net.IPMask = make(net.IPMask, net.IPv4len)
        err error = nil
        hexMask int64
    )

    // cast hex net mask to uint32
    hexMask, err = strconv.ParseInt(fmt.Sprintf("0x%s",gwHexMask), 0, 64)
    if err != nil {
        return gwMask, "", errors.WithStack(err)
    }

    // make net.IP address from uint32
    binary.LittleEndian.PutUint32(gwMask, uint32(hexMask))

    // format net.IP to dotted ipV4 string
    return gwMask, net.IP(gwMask).String(), nil
}

func DefaultIPv4Gateway() (*Gateway, error){
    var (
        scanner *bufio.Scanner = nil
        file *os.File = nil
        err error = nil
        tokens []string = nil
        gwIP net.IP
        gwMask net.IPMask
        addr, mask, iface string
    )

    file, err = os.Open(route_tbl_path)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    defer file.Close()

    // jump to line containing the agteway address
    scanner = bufio.NewScanner(file)
    scanner.Scan()

    for scanner.Scan() {
        // get the default gateway address
        tokens = strings.Split(scanner.Text(), line_separater)

        // get field containing gateway address
        gwIP, addr, err = hexAddressToIPv4(tokens[field_gateway])
        if err != nil {
            return nil, errors.WithStack(err)
        }

        // get gateway interface
        iface = strings.TrimSpace(tokens[field_interface])

        // get gateway mask
        gwMask, mask, err = hexMaskToIPv4(tokens[field_mask])
        if err != nil {
            return nil, errors.WithStack(err)
        }

        return &Gateway {
            AddrMask:   net.IPNet {
                IP:     gwIP,
                Mask:   gwMask,
            },
            Address:    addr,
            Mask:       mask,
            Interface:  iface,
        }, nil
    }

    return nil, errors.Errorf("[ERR] should have returned gateway info already")
}

func AllIPv4Gateways() (map[string][]Gateway, error) {
    var (
        gwList = map[string][]Gateway{}
        gwIP net.IP
        gwMask net.IPMask
        scanner *bufio.Scanner = nil
        file *os.File = nil
        err error = nil
        tokens []string = nil
        addr, mask, iface string
    )

    file, err = os.Open(route_tbl_path)
    if err != nil {
        return gwList, errors.WithStack(err)
    }
    defer file.Close()

    // jump to line containing the agteway address
    scanner = bufio.NewScanner(file)
    scanner.Scan()

    for scanner.Scan() {
        // get the default gateway address
        tokens = strings.Split(scanner.Text(), line_separater)

        // get field containing gateway address
        gwIP, addr, err = hexAddressToIPv4(tokens[field_gateway])
        if err != nil {
            continue
        }

        // get gateway interface
        iface = strings.TrimSpace(tokens[field_interface])

        // get gateway mask
        gwMask, mask, err = hexMaskToIPv4(tokens[field_mask])
        if err != nil {
            continue
        }

        // allocate & append gw instance
        gwList[iface] = append(gwList[iface],
            Gateway{
                AddrMask:   net.IPNet{
                    IP:     gwIP,
                    Mask:   gwMask,
                },
                Address:    addr,
                Mask:       mask,
                Interface:  iface,
            })
    }

    return gwList, nil
}

func FindIPv4GatewayWithInterface(iName string) (*Gateway, error) {
    return nil, nil
}