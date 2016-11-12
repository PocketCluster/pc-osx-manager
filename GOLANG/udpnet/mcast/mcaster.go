package mcast

import (
    "net"
    "sync"
    "log"
    "fmt"
    "time"
)

type multiCaster struct {
    ipv4UnicastConn     *net.UDPConn

    closed              bool
    closedCh            chan struct{}
    closeLock           sync.Mutex

    chWrite             chan *CastPkg
    log                 *log.Logger
}

func NewMultiCaster(log *log.Logger) (*multiCaster, error) {
    // Create a IPv4 listener
    uconn4, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
    if err != nil {
        return nil, fmt.Errorf("failed to bind to any unicast udp port : " + err.Error())
    }

    mc := &multiCaster{
        ipv4UnicastConn      : uconn4,
        closedCh             : make(chan struct{}),
        chWrite              : make(chan *CastPkg),//, PC_MCAST_CASTER_CHAN_CAP),
        log                  : log,
    }
    go mc.write()
    return mc, nil
}

// Close is used to cleanup the client
func (mc *multiCaster) Close() error {
    mc.closeLock.Lock()
    defer mc.closeLock.Unlock()

    if mc.closed {
        return nil
    }
    mc.closed = true

    close(mc.closedCh)
    close(mc.chWrite)

    return mc.ipv4UnicastConn.Close()
}

func (mc *multiCaster) write() {
    for cp := range mc.chWrite {
        // TODO : we can timeout if necessary
        _, e := mc.ipv4UnicastConn.WriteToUDP(cp.Message, ipv4McastAddr)
        if e != nil && mc.log != nil {
            mc.log.Println(e)
        }
    }
}

// sendQuery is used to multicast a query out
func (mc *multiCaster) Send(message []byte) error {
    if len(message) == 0 {
        return fmt.Errorf("[ERR] Multicast message is empty")
    }
    cp := &CastPkg{
        Message     : message,
    }
    mc.chWrite <- cp

    time.After(time.Millisecond)
    return
}
