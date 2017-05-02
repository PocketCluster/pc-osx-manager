package ucast

import (
    "net"
    "time"
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

type BeaconChannel struct {
    closedCh     chan struct{}
    closeLock    sync.Mutex

    conn         *net.UDPConn
    waiter       *sync.WaitGroup
    ChRead       chan BeaconPack
    chWrite      chan BeaconPack
}

// New constructor of a new server
func NewPocketBeaconChannel(waiter *sync.WaitGroup) (*BeaconChannel, error) {
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: PAGENT_SEND_PORT})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    beacon := &BeaconChannel {
        conn:       conn,
        waiter:     waiter,
        ChRead:     make(chan BeaconPack, PC_UCAST_BEACON_CHAN_CAP),
        chWrite:    make(chan BeaconPack, PC_UCAST_BEACON_CHAN_CAP),
        closedCh:   make(chan struct{}),
    }
    go beacon.reader()
    go beacon.writer()
    return beacon, nil
}

// Close is used to cleanup the client
func (bc *BeaconChannel) Close() error {
    bc.closeLock.Lock()
    defer bc.closeLock.Unlock()

    _, isOpen := <- bc.closedCh
    if !isOpen {
        return nil
    }

    log.Debugf("[INFO] locator channel closing : %v", *bc)

    close(bc.closedCh)
    close(bc.ChRead)
    close(bc.chWrite)

    return bc.conn.Close()
}

func (bc *BeaconChannel) reader() {
    var (
        buff []byte          = make([]byte, PC_MAX_UCAST_UDP_BUF_SIZE)
        addr *net.UDPAddr    = nil
        terr, err error      = nil, nil
        count int            = 0
        ticker *time.Ticker  = time.NewTicker(readTimeout)
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
    defer bc.waiter.Done()
    defer ticker.Stop()

    for {
        select {
        case <- bc.closedCh:
            return

        case <- ticker.C:
            terr = bc.conn.SetReadDeadline(time.Now().Add(readTimeout))
            if terr != nil {
                continue
            }
            count, addr, err = bc.conn.ReadFromUDP(buff)
            if err != nil {
                log.Infof("[INFO] locator channel : Failed to read packet: %v", err)
                continue
            }
            if count == 0 {
                log.Infof("[INFO] empty message. ignore")
                continue
            }
            adr := copyUDPAddr(addr)
            msg := make([]byte, count)
            copy(msg, buff[:count])
            pack := BeaconPack{
                Address:    adr,
                Message:    msg,
            }
            bc.ChRead <- pack
        }
    }
}

func (bc *BeaconChannel) writer() {
    defer bc.waiter.Done()
    for {
        select {
        case <- bc.closedCh:
            return

        case v := <-bc.chWrite:
            if len(v.Message) == 0 {
                continue
            }
            _, err := bc.conn.WriteToUDP(v.Message, &v.Address)
            if err != nil {
                log.Info(err)
            }
        }
    }

}

func (bc *BeaconChannel) Send(targetHost string, buf []byte) error {
    if len(targetHost) == 0 || len(buf) == 0 {
        return errors.Errorf("[ERR] Cannot send null data to null host")
    }
    bc.chWrite <- BeaconPack{
        Message: buf,
        Address: net.UDPAddr {
            IP:      net.ParseIP(targetHost),
            Port:    PAGENT_SEND_PORT,
        },
    }

    // TODO : find ways to remove this. We'll wait artificially for now (v0.1.4)
    time.After(time.Millisecond)
    return nil
}
