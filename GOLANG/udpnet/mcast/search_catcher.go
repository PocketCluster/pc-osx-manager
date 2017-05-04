package mcast

import (
    "net"
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "time"
)

type SearchCatcher struct {
    isClosed     bool
    closeLock    sync.Mutex

    waiter       *sync.WaitGroup
    conn         *net.UDPConn
    ChRead       chan CastPack
}

func NewSearchCatcher(niface string, waiter *sync.WaitGroup) (*SearchCatcher, error) {
    iface, err := net.InterfaceByName(niface)
    if err != nil {
        log.Error(err)
        return nil, errors.WithStack(err)
    }
    conn, err := net.ListenMulticastUDP("udp4", iface, ipv4McastAddr)
    if err != nil {
        return nil, errors.Errorf("[ERR] failed to bind to any multicast udp port", err)
    }
    conn.SetReadBuffer(PC_MAX_MCAST_UDP_BUF_SIZE)
    listener := &SearchCatcher{
        waiter:    waiter,
        conn:      conn,
        ChRead:    make(chan CastPack, PC_MCAST_LISTENER_CHAN_CAP),
    }
    waiter.Add(1)
    go listener.read()
    return listener, nil
}

// Close is used to cleanup the client
func (sl *SearchCatcher) Close() error {
    sl.closeLock.Lock()
    defer sl.closeLock.Unlock()

    if sl.isClosed {
        return nil
    }
    sl.isClosed = true
    close(sl.ChRead)
    return sl.conn.Close()
}

func (sl *SearchCatcher) read() {
    var (
        buff []byte            = make([]byte, PC_MAX_MCAST_UDP_BUF_SIZE)
        addr *net.UDPAddr      = nil
        err error              = nil
        count int              = 0
    )

    copyUDPAddr := func(adr *net.UDPAddr) net.UDPAddr {
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
    defer sl.waiter.Done()

    for !sl.isClosed {
        // Set a deadline for reading. Read operation will fail if no data
        // is received after deadline.
        sl.conn.SetReadDeadline(time.Now().Add(readTimeout))

        count, addr, err = sl.conn.ReadFromUDP(buff)
        if err != nil {
            continue
        }
        if count == 0 {
            continue
        }
        adr := copyUDPAddr(addr)
        msg := make([]byte, count)
        copy(msg, buff[:count])
        pack := CastPack{
            Address:    adr,
            Message:    msg,
        }
        sl.ChRead <- pack
    }
}