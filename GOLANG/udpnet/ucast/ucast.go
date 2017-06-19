package ucast

import (
    "net"
    "time"
)

const (
    POCKET_LOCATOR_PORT = 10060
    POCKET_AGENT_PORT   = 10061

    //http://stackoverflow.com/questions/1098897/what-is-the-largest-safe-udp-packet-size-on-the-internet
    PC_SAFE_UDP_PKT_SIZE = 508

    // for IPv4 network max UDP packet size is 65507
    // http://stackoverflow.com/questions/1098897/what-is-the-largest-safe-udp-packet-size-on-the-internet
    PC_MAX_UCAST_UDP_BUF_SIZE = 65507

    // total 256 IP class C client can exists. Since we're to remove router, broadcast, and beacon itself, 253 is the number
    PC_UCAST_BEACON_CHAN_CAP = 254

    // locator channel capacitor doesn't need to be big. just big enough to hold communication with beacon
    PC_UCAST_LOCATOR_CHAN_CAP = 4

    readTimeout = time.Second * 3
)

type BeaconPack struct {
    Address    net.UDPAddr
    Message    []byte
}

type BeaconSend struct {
    Host       string
    Payload    []byte
}

func copyUDPAddr(adr *net.UDPAddr) net.UDPAddr {
    lenIP := len(adr.IP)
    ip := make([]byte, lenIP)
    copy(ip, adr.IP)
    zone := string([]byte(adr.Zone))

    return net.UDPAddr {
        IP:     ip,
        Port:   adr.Port,
        Zone:   zone,
    }
}
