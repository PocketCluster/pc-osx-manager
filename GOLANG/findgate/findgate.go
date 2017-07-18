// +build linux

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
    route_tbl_path  = "/proc/net/route"
    line  = 1    // line containing the gateway addr. (first line: 0)
    sep   = "\t" // field separator
    field = 2    // field containing hex gateway address (first field: 0)
)

func DefaultGatewayWithInterface() (string, string, error){
    var (
        scanner *bufio.Scanner = nil
        file *os.File = nil
        err error = nil
        tokens []string = nil
        gatewayHex string
        hexAddr int64
        intAddr uint32
    )

    file, err = os.Open(route_tbl_path)
    if err != nil {
        return "", "", errors.WithStack(err)
    }
    defer file.Close()

    scanner = bufio.NewScanner(file)
    for scanner.Scan() {

        // jump to line containing the agteway address
        for i := 0; i < line; i++ {
            scanner.Scan()
        }

        // get field containing gateway address
        tokens = strings.Split(scanner.Text(), sep)
        gatewayHex = "0x" + tokens[field]

        // cast hex address to uint32
        hexAddr, err = strconv.ParseInt(gatewayHex, 0, 64)
        if err != nil {
            continue
        }
        intAddr = uint32(hexAddr)

        // make net.IP address from uint32
        ipd32 := make(net.IP, 4)
        binary.LittleEndian.PutUint32(ipd32, intAddr)
        fmt.Printf("%T --> %[1]v\n", ipd32)

        // format net.IP to dotted ipV4 string
        ip := net.IP(ipd32).String()
        fmt.Printf("%T --> %[1]v\n", ip)

        // exit scanner
        break
    }

    return "", "", nil
}

func FindGatewayForInterface(iName string) (string, error) {

    return "", nil
}