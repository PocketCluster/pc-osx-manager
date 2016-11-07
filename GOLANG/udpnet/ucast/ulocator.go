package ucast

import (
    "net"
    "sync"
    "fmt"
    "log"
    "time"
)

type locatorChannel struct {
    closed       bool
    closedCh     chan struct{}
    closeLock    sync.Mutex
    log          *log.Logger

    conn         *net.UDPConn
    ChRead       chan *ChanPkg
    chWrite      chan *ChanPkg
}

// New constructor of a new server
func NewPocketLocatorChannel(log *log.Logger) (*locatorChannel, error) {
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: PAGENT_RECV_PORT})
    if err != nil {
        return nil, err
    }
    locator := &locatorChannel {
        conn       : conn,
        ChRead     : make(chan *ChanPkg, PC_UCAST_LOCATOR_CHAN_CAP),
        chWrite    : make(chan *ChanPkg, PC_UCAST_LOCATOR_CHAN_CAP),
        closedCh   : make(chan struct{}),
        log        : log,
    }
    go locator.reader()
    go locator.writer()
    return locator, nil
}

// Close is used to cleanup the client
func (lc *locatorChannel) Close() error {
    lc.closeLock.Lock()
    defer lc.closeLock.Unlock()

    if lc.log != nil {
        log.Printf("[INFO] locator channel closing : %v", *lc)
    }

    if lc.closed {
        return nil
    }
    lc.closed = true

    close(lc.closedCh)
    close(lc.ChRead)
    close(lc.chWrite)

    if lc.conn != nil {
        lc.conn.Close()
    }
    return nil
}

func (lc *locatorChannel) reader() {
    var (
        err error
        count int
    )

    for !lc.closed {
        pack := &ChanPkg{}
        pack.Message = make([]byte, PC_MAX_UCAST_UDP_BUF_SIZE)
        count, pack.Address, err = lc.conn.ReadFromUDP(pack.Message)
        if err != nil {
            if lc.log != nil {
                lc.log.Printf("[ERR] locator channel : Failed to read packet: %v", err)
            }
            continue
        }

        if lc.log != nil {
            lc.log.Printf("[INFO] %d bytes have been received", count)
        }
        pack.Message = pack.Message[:count]
        select {
        case lc.ChRead <- pack:
        case <-lc.closedCh:
            return
        }
    }
}

func (lc *locatorChannel) writer() {
    for v := range lc.chWrite {
        _, e := lc.conn.WriteToUDP(v.Message, v.Address)
        if e != nil && lc.log != nil {
            lc.log.Println(e)
        }
    }
}

func (lc *locatorChannel) Send(targetHost string, buf []byte) error {
    if len(targetHost) == 0 || len(buf) == 0 {
        return fmt.Errorf("[ERR] Cannot send null data to null host")
    }
    targetAddr := &net.UDPAddr{
        IP      : net.ParseIP(targetHost),
        Port    : PAGENT_SEND_PORT,
    }
    lc.chWrite <- &ChanPkg{
        Message    : buf,
        Address    : targetAddr,
    }

    // TODO : find ways to remove this. We'll wait artificially for now (v0.1.4)
    time.After(time.Millisecond)
    return nil
}
