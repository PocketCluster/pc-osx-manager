package mcast

import (
    "net"
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

type SearchCaster struct {
    isClosed    bool

    waiter      *sync.WaitGroup
    conn        *net.UDPConn
    chWrite     chan CastPack
}

func NewSearchCaster(waiter *sync.WaitGroup) (*SearchCaster, error) {
    // Create a IPv4 listener
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    conn.SetReadBuffer(PC_MAX_MCAST_UDP_BUF_SIZE)
    conn.SetWriteBuffer(PC_MAX_MCAST_UDP_BUF_SIZE)
    mc := &SearchCaster{
        waiter:     waiter,
        conn:       conn,
        chWrite:    make(chan CastPack),//, PC_MCAST_CASTER_CHAN_CAP),
    }
    waiter.Add(1)
    go mc.write()
    return mc, nil
}

// Close is used to cleanup the client
func (mc *SearchCaster) Close() error {
    if mc.isClosed {
        return nil
    }

    mc.isClosed = true
    close(mc.chWrite)
    return errors.WithStack(mc.conn.Close())
}

func (mc *SearchCaster) write() {
    defer mc.waiter.Done()

    for cp := range mc.chWrite {
        // TODO : we can timeout if necessary
        _, e := mc.conn.WriteToUDP(cp.Message, ipv4McastAddr)
        if e != nil {
            log.Info(errors.WithStack(e))
        }
    }
}

// sendQuery is used to multicast a query out
func (mc *SearchCaster) Send(message []byte) error {
    if len(message) == 0 {
        return errors.Errorf("[ERR] multicast message is empty")
    }
    mc.chWrite <- CastPack{
        Message: message,
    }

    time.After(time.Millisecond)
    return nil
}
