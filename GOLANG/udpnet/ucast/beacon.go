package ucast

import (
    "net"
    "log"
    "sync"
    "fmt"
    "time"
)

type BeaconChannel struct {
    closed       bool
    closedCh     chan struct{}
    closeLock    sync.Mutex
    log          *log.Logger

    conn         *net.UDPConn
    ChRead       chan *ChanPkg
    chWrite      chan *ChanPkg
}

// New constructor of a new server
func NewPocketBeaconChannel(log *log.Logger) (*BeaconChannel, error) {
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: PAGENT_SEND_PORT})
    if err != nil {
        return nil, err
    }
    beacon := &BeaconChannel {
        conn       : conn,
        ChRead     : make(chan *ChanPkg, PC_BEACON_CHAN_CAP),
        chWrite    : make(chan *ChanPkg, PC_BEACON_CHAN_CAP),
        closedCh   : make(chan struct{}),
        log        : log,
    }
    go beacon.reader()
    go beacon.writer()
    return beacon, nil
}

// Close is used to cleanup the client
func (bc *BeaconChannel) Close() error {
    bc.closeLock.Lock()
    defer bc.closeLock.Unlock()

    if bc.log != nil {
        log.Printf("[INFO] beacon channel closing : %v", *bc)
    }

    if bc.closed {
        return nil
    }
    bc.closed = true

    close(bc.closedCh)
    close(bc.ChRead)
    close(bc.chWrite)

    if bc.conn != nil {
        bc.conn.Close()
    }
    return nil
}

func (bc *BeaconChannel) reader() {
    var (
        err error
        count int
    )

    for !bc.closed {
        pack := &ChanPkg{}
        pack.Pack = make([]byte, PC_MAX_UDP_BUF_SIZE)
        count, pack.Addr, err = bc.conn.ReadFromUDP(pack.Pack)
        if err != nil {
            if bc.log != nil {
                bc.log.Printf("[ERR] beacon channel : Failed to read packet: %v", err)
            }
            continue
        }

        if bc.log != nil {
            bc.log.Printf("[INFO] %d bytes have been received", count)
        }
        pack.Pack = pack.Pack[:count]
        select {
            case bc.ChRead <- pack:
            case <-bc.closedCh:
                return
        }
    }
}

func (bc *BeaconChannel) writer() {
    for v := range bc.chWrite {
        _, e := bc.conn.WriteToUDP(v.Pack, v.Addr)
        if e != nil && bc.log != nil {
            bc.log.Println(e)
        }
    }
}

func (bc *BeaconChannel) Send(targetHost string, buf []byte) error {
    if len(targetHost) == 0 || len(buf) == 0 {
        return fmt.Errorf("[ERR] Cannot send null data to null host")
    }
    targetAddr := &net.UDPAddr{
        IP      : net.ParseIP(targetHost),
        Port    : PAGENT_RECV_PORT,
    }
    bc.chWrite <- &ChanPkg{
        Pack    : buf,
        Addr    : targetAddr,
    }

    // TODO : find ways to remove this. We'll wait artificially for now (v0.1.4)
    time.Sleep(time.Microsecond * 100)
    return nil
}
