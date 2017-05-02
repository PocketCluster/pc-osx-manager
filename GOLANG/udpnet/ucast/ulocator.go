package ucast

import (
    "net"
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

type LocatorChannel struct {
    closedCh     chan struct{}
    closeLock    sync.Mutex

    conn         *net.UDPConn
    waiter       *sync.WaitGroup
    ChRead       chan BeaconPack
    chWrite      chan BeaconPack
}

// New constructor of a new server
func NewPocketLocatorChannel(srvWaiter *sync.WaitGroup) (*LocatorChannel, error) {
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: PAGENT_RECV_PORT})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    locator := &LocatorChannel{
        conn:        conn,
        waiter:      srvWaiter,
        ChRead:      make(chan BeaconPack, PC_UCAST_LOCATOR_CHAN_CAP),
        chWrite:     make(chan BeaconPack, PC_UCAST_LOCATOR_CHAN_CAP),
        closedCh:    make(chan struct{}),
    }
    go locator.reader()
    go locator.writer()
    return locator, nil
}

// Close is used to cleanup the client
func (lc *LocatorChannel) Close() error {
    lc.closeLock.Lock()
    defer lc.closeLock.Unlock()

    _, isOpen := <- lc.closedCh
    if !isOpen {
        return nil
    }

    log.Debugf("[INFO] locator channel closing : %v", *lc)

    close(lc.closedCh)
    close(lc.ChRead)
    close(lc.chWrite)

    return lc.conn.Close()
}

func (lc *LocatorChannel) reader() {
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
    defer lc.waiter.Done()
    defer ticker.Stop()

    for {
        select {
            case <- lc.closedCh:
                return

            case <- ticker.C:
                terr = lc.conn.SetReadDeadline(time.Now().Add(readTimeout))
                if terr != nil {
                    continue
                }
                count, addr, err = lc.conn.ReadFromUDP(buff)
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
                lc.ChRead <- pack
        }
    }
}

func (lc *LocatorChannel) writer() {
    defer lc.waiter.Done()
    for {
        select {
            case <- lc.closedCh:
                return

            case v := <-lc.chWrite:
                if len(v.Message) == 0 {
                    continue
                }
                _, err := lc.conn.WriteToUDP(v.Message, &v.Address)
                if err != nil {
                    log.Info(err)
                }
            }
    }
}

func (lc *LocatorChannel) Send(targetHost string, buf []byte) error {
    if len(targetHost) == 0 || len(buf) == 0 {
        return errors.Errorf("[ERR] Cannot send null data to null host")
    }
    lc.chWrite <- BeaconPack{
        Message: buf,
        Address: net.UDPAddr {
            IP:      net.ParseIP(targetHost),
            Port:    PAGENT_SEND_PORT,
        },
    }

    // TODO : find ways to remove this. We'll wait artificially for now (v0.1.4)
    //time.After(time.Millisecond)
    return nil
}
