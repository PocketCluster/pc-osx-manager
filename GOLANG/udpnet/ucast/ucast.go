package ucast

import (
    "time"
    "net"
)

const (
    PAGENT_SEND_PORT    = 10060
    PAGENT_RECV_PORT    = 10061
)

//http://stackoverflow.com/questions/1098897/what-is-the-largest-safe-udp-packet-size-on-the-internet
const PC_SAFE_UDP_PKT_SIZE = 508

// for IPv4 network max UDP packet size is 65507
// http://stackoverflow.com/questions/1098897/what-is-the-largest-safe-udp-packet-size-on-the-internet
const PC_MAX_UDP_BUF_SIZE = 65507

// total 256 IP class C client can exists. Since we're to remove router, broadcast, and beacon itself, 253 is the number
const PC_MAX_COMM_CHAN_CAP = 254

type ConnParam struct {
    RecvMessage         chan <- []byte       // Message to recv
    Timeout             time.Duration        // Lookup timeout, default 1 second
}

type (
    ChanPkg struct {
        Pack    []byte
        Addr    *net.UDPAddr
    }
)