package ucast

import "time"

const (
    PAGENT_SEND_PORT    = 10060
    PAGENT_RECV_PORT    = 10061
)

type ConnParam struct {
    RecvMessage         chan <- []byte       // Message to recv
    Timeout             time.Duration        // Lookup timeout, default 1 second
}

