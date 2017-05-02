package ucast

import (
    "net"
    "time"
)

const (
    PAGENT_SEND_PORT    = 10060
    PAGENT_RECV_PORT    = 10061

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
    Message     []byte
    Address     net.UDPAddr
}
