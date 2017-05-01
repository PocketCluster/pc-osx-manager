package mcast

import (
    "net"
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

type multiListener struct {
    ipv4mconn    *net.UDPConn

    closed       bool
    closedCh     chan struct{}
    closeLock    sync.Mutex

    ChRead       chan *CastPkg
}

func NewMcastListener(niface string) (*multiListener, error) {
    iface, err := net.InterfaceByName(niface)
    if err != nil {
        log.Error(err)
        return nil, errors.WithStack(err)
    }
    mconn4, err := net.ListenMulticastUDP("udp4", iface, ipv4McastAddr)
    if err != nil {
        log.Errorf("[ERR] failed to bind to any multicast udp port %v", err)
        return nil, errors.Errorf("[ERR] failed to bind to any multicast udp port", err)
    }
    listener := &multiListener{
        ipv4mconn    : mconn4,
        closedCh     : make(chan struct{}),
        ChRead       : make(chan *CastPkg, PC_MCAST_LISTENER_CHAN_CAP),
    }
    go listener.read()
    return listener, nil
}

// Close is used to cleanup the client
func (ml *multiListener) Close() error {
    ml.closeLock.Lock()
    defer ml.closeLock.Unlock()

    if ml.closed {
        return nil
    }
    ml.closed = true

    close(ml.closedCh)
    close(ml.ChRead)

    if ml.ipv4mconn != nil {
        ml.ipv4mconn.Close()
    }
    return nil
}


// recv is used to receive until we get a shutdown
func (ml *multiListener) read() {
    var (
        err error
        count int
    )
    for !ml.closed {
        pack := &CastPkg{}
        pack.Message = make([]byte, PC_MAX_MCAST_UDP_BUF_SIZE)
        count, pack.Address, err = ml.ipv4mconn.ReadFromUDP(pack.Message)
        if err != nil {
            log.Info("[INFO] beacon channel : Failed to read packet: %v", err)
            continue
        }

        pack.Message = pack.Message[:count]
        select {
        case ml.ChRead <- pack:
        case <- ml.closedCh:
            return
        }
    }
}
