package mcast

import (
    "net"
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

type SearchCaster struct {
    conn        *net.UDPConn
    chWrite     chan CastPack
    chClosed    chan bool
}

func NewSearchCaster() (*SearchCaster, error) {
    // Create a IPv4 listener
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    conn.SetWriteBuffer(PC_MAX_MCAST_UDP_BUF_SIZE)
    mc := &SearchCaster{
        conn:       conn,
        chClosed:   make(chan bool),
        chWrite:    make(chan CastPack),//, PC_MCAST_CASTER_CHAN_CAP),
    }
    go mc.write()
    return mc, nil
}

// Close is used to cleanup the client
func (mc *SearchCaster) Close() error {
    close(mc.chClosed)
    close(mc.chWrite)
    return errors.WithStack(mc.conn.Close())
}

func (mc *SearchCaster) write() {
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
    select {
        case <-mc.chClosed:
            return nil
        case mc.chWrite <- CastPack{
            Message: message,
        }:
        time.After(time.Millisecond)
    }
    return nil
}
