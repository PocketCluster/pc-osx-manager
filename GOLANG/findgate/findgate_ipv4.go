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
    ipv4_route_tbl_path string = "/proc/net/route"
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

// Reference : /usr/include/linux/route.h
const (
    rtf_up          uint16 = 0x0001    /* route usable                 */
    rtf_gateway     uint16 = 0x0002    /* destination is a gateway     */
    rtf_host        uint16 = 0x0004    /* host entry (net otherwise)   */
    rtf_reinstate   uint16 = 0x0008    /* reinstate route after tmout  */
    rtf_dynamic     uint16 = 0x0010    /* created dyn. (by redirect)   */
    rtf_modified    uint16 = 0x0020    /* modified dyn. (by redirect)  */
    rtf_mtu         uint16 = 0x0040    /* specific MTU for this route  */
    rtf_mss         uint16 = rtf_mtu   /* Compatibility :-(            */
    rtf_window      uint16 = 0x0080    /* per route window clamping    */
    rtf_irtt        uint16 = 0x0100    /* Initial round trip time      */
    rtf_reject      uint16 = 0x0200    /* Reject route                 */
)

type IPv4Gateway struct {
    AddrMask     net.IPNet
    Flag         uint16
    Address      string
    Mask         string
    Interface    string
}

func (i *IPv4Gateway) IsUsable() bool {
    return (i.Flag & rtf_up) != 0
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

func hexFlagToUint16(gwHexFlag string) (uint16, error) {
    var (
        err error = nil
        hexFlag int64
    )

    // cast hex net mask to uint16
    hexFlag, err = strconv.ParseInt(fmt.Sprintf("0x%s",gwHexFlag), 0, 64)
    if err != nil {
        return 0, errors.WithStack(err)
    }

    // format int64 to unsigned short
    return uint16(hexFlag), nil
}

func DefaultIPv4Gateway() (*IPv4Gateway, error){
    var (
        scanner *bufio.Scanner = nil
        file *os.File = nil
        err error = nil
        tokens []string = nil
        gwIP net.IP
        gwMask net.IPMask
        addr, mask, iface, entry string
        flag uint16 = 0
    )

    file, err = os.Open(ipv4_route_tbl_path)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    defer file.Close()

    // jump to line containing the agteway address
    scanner = bufio.NewScanner(file)
    scanner.Scan()

    for scanner.Scan() {
        // get the default gateway address
        entry = strings.TrimSpace(scanner.Text())
        if len(entry) == 0 {
            continue
        }
        tokens = strings.Split(entry, line_separater)

        // get gateway interface
        iface = strings.TrimSpace(tokens[field_interface])

        // get field containing gateway address
        gwIP, addr, err = hexAddressToIPv4(tokens[field_gateway])
        if err != nil {
            return nil, errors.WithStack(err)
        }

        // get gateway mask
        gwMask, mask, err = hexMaskToIPv4(tokens[field_mask])
        if err != nil {
            return nil, errors.WithStack(err)
        }

        // get gateway flag
        flag, err = hexFlagToUint16(tokens[field_flags])
        if err != nil {
            return nil, errors.WithStack(err)
        }

        return &IPv4Gateway{
            AddrMask:   net.IPNet {
                IP:     gwIP,
                Mask:   gwMask,
            },
            Flag:       flag,
            Address:    addr,
            Mask:       mask,
            Interface:  iface,
        }, nil
    }

    return nil, errors.Errorf("[ERR] should have returned gateway info already")
}

func AllIPv4Gateways() (map[string][]IPv4Gateway, error) {
    var (
        gwList = map[string][]IPv4Gateway{}
        gwIP net.IP
        gwMask net.IPMask
        scanner *bufio.Scanner = nil
        file *os.File = nil
        err error = nil
        tokens []string = nil
        addr, mask, iface, entry string
        flag uint16 = 0
    )

    file, err = os.Open(ipv4_route_tbl_path)
    if err != nil {
        return gwList, errors.WithStack(err)
    }
    defer file.Close()

    // jump to line containing the agteway address
    scanner = bufio.NewScanner(file)
    scanner.Scan()

    for scanner.Scan() {
        // get the default gateway address
        entry = strings.TrimSpace(scanner.Text())
        if len(entry) == 0 {
            continue
        }
        tokens = strings.Split(entry, line_separater)

        // get gateway interface
        iface = strings.TrimSpace(tokens[field_interface])

        // get field containing gateway address
        gwIP, addr, err = hexAddressToIPv4(tokens[field_gateway])
        if err != nil {
            continue
        }

        // get gateway mask
        gwMask, mask, err = hexMaskToIPv4(tokens[field_mask])
        if err != nil {
            continue
        }

        // get gateway flag
        flag, err = hexFlagToUint16(tokens[field_flags])
        if err != nil {
            return nil, errors.WithStack(err)
        }

        // allocate & append gw instance
        gwList[iface] = append(gwList[iface],
            IPv4Gateway{
                AddrMask:   net.IPNet{
                    IP:     gwIP,
                    Mask:   gwMask,
                },
                Flag:       flag,
                Address:    addr,
                Mask:       mask,
                Interface:  iface,
            })
    }

    return gwList, nil
}

func FindIPv4GatewayWithInterface(iName string) ([]IPv4Gateway, error) {
    var (
        gwList = []IPv4Gateway{}
        gwIP net.IP
        gwMask net.IPMask
        scanner *bufio.Scanner = nil
        file *os.File = nil
        err error = nil
        tokens []string = nil
        addr, mask, iface, entry string
        flag uint16 = 0
    )

    if len(iName) == 0 {
        return gwList, errors.Errorf("[ERR] invalid interface name")
    }

    file, err = os.Open(ipv4_route_tbl_path)
    if err != nil {
        return gwList, errors.WithStack(err)
    }
    defer file.Close()

    // jump to line containing the agteway address
    scanner = bufio.NewScanner(file)
    scanner.Scan()

    for scanner.Scan() {
        // get the default gateway address
        entry = strings.TrimSpace(scanner.Text())
        if len(entry) == 0 {
            continue
        }
        tokens = strings.Split(entry, line_separater)

        // get gateway interface & match with given name
        iface = strings.TrimSpace(tokens[field_interface])
        if iface != iName {
            continue
        }

        // get field containing gateway address
        gwIP, addr, err = hexAddressToIPv4(tokens[field_gateway])
        if err != nil {
            continue
        }

        // get gateway mask
        gwMask, mask, err = hexMaskToIPv4(tokens[field_mask])
        if err != nil {
            continue
        }

        // get gateway flag
        flag, err = hexFlagToUint16(tokens[field_flags])
        if err != nil {
            return nil, errors.WithStack(err)
        }

        // allocate & append gw instance
        gwList = append(gwList,
            IPv4Gateway{
                AddrMask:   net.IPNet{
                    IP:     gwIP,
                    Mask:   gwMask,
                },
                Flag:       flag,
                Address:    addr,
                Mask:       mask,
                Interface:  iface,
            })
    }

    // check if we've found any
    if len(gwList) == 0 {
        return gwList, errors.Errorf("[ERR] cannot find a gateway for interface %s", iName)
    }
    return gwList, nil
}