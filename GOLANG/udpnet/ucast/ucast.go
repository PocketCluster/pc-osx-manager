package ucast

import "time"

const (
    PAGENT_SEND_PORT    = 10060
    PAGENT_RECV_PORT    = 10061
)

//http://stackoverflow.com/questions/1098897/what-is-the-largest-safe-udp-packet-size-on-the-internet
const PC_SAFE_UDP_PKT_SIZE = 508

// for IPv4 network max UDP packet size is 65507
// http://stackoverflow.com/questions/1098897/what-is-the-largest-safe-udp-packet-size-on-the-internet
const PC_MAX_UDP_BUF_SIZE = 65507

type ConnParam struct {
    RecvMessage         chan <- []byte       // Message to recv
    Timeout             time.Duration        // Lookup timeout, default 1 second
}

