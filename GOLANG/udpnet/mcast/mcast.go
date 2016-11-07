package mcast

import (
    "net"
    "time"
)

const (
    ipv4mdns   = "224.0.0.251"
    ipv6mdns   = "ff02::fb"
    mdnsPort   = 5353

    // for IPv4 network max UDP packet size is 65507
    // http://stackoverflow.com/questions/1098897/what-is-the-largest-safe-udp-packet-size-on-the-internet
    PC_MAX_MCAST_UDP_BUF_SIZE int = 65507

    // total 256 IP class C client can exists. Since we're to remove router, broadcast, and beacon itself, 253 is the number
    PC_MCAST_LISTENER_CHAN_CAP int = 254

    // locator channel capacitor doesn't need to be big. just big enough to hold communication with beacon
    PC_MCAST_CASTER_CHAN_CAP int = 16
)

var (
    ipv4McastAddr = &net.UDPAddr{
        IP:   net.ParseIP(ipv4mdns),
        Port: mdnsPort,
    }
    ipv6McastAddr = &net.UDPAddr{
        IP:   net.ParseIP(ipv6mdns),
        Port: mdnsPort,
    }
)

type (
    CastPkg struct {
        Message     []byte
        Address     *net.UDPAddr
        Timeout     time.Duration
    }
)