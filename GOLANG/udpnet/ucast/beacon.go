package ucast

import (
    "net"
    "log"
    "sync"
    "fmt"
)

type (
    BeaconChannel struct {
        closed       bool
        closedCh     chan struct{}
        closeLock    sync.Mutex
        log          *log.Logger

        conn         *net.UDPConn
        ChRead       chan *ChanPkg
        ChWrite      chan *ChanPkg
    }
)

// New constructor of a new server
func NewPocketBeaconChannel(log *log.Logger) (*BeaconChannel, error) {
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: PAGENT_SEND_PORT})
    if err != nil {
        return nil, err
    }
    beacon := &BeaconChannel {
        conn       : conn,
        ChRead     : make(chan *ChanPkg, PC_MAX_COMM_CHAN_CAP),
        ChWrite    : make(chan *ChanPkg, PC_MAX_COMM_CHAN_CAP),
        log        : log,
    }
    go beacon.reader()
    go beacon.writer()
    return beacon, nil
}


// Close is used to cleanup the client
func (c *BeaconChannel) Close() error {
    c.closeLock.Lock()
    defer c.closeLock.Unlock()

    if c.log != nil {
        log.Printf("[INFO] beacon channel closing : %v", *c)
    }

    if c.closed {
        return nil
    }
    c.closed = true

    close(c.closedCh)
    close(c.ChRead)
    close(c.ChWrite)

    if c.conn != nil {
        c.conn.Close()
    }
    return nil
}

func (c *BeaconChannel) reader() {
    var (
        err error
        count int
    )

    for !c.closed {
        pack := &ChanPkg{}
        pack.Pack = make([]byte, PC_MAX_UDP_BUF_SIZE)
        count, pack.Addr, err = c.conn.ReadFromUDP(pack.Pack)
        if err != nil {
            if c.log != nil {
                c.log.Printf("[ERR] beacon channel : Failed to read packet: %v", err)
            }
            continue
        }

        if c.log != nil {
            c.log.Printf("[INFO] %d bytes have been received", count)
        }
        pack.Pack = pack.Pack[:count]
        select {
            case c.ChRead <- pack:
            case <-c.closedCh:
                return
        }
    }
}

func (t *BeaconChannel) writer() {
    for v := range t.ChWrite {
        _, e := t.conn.WriteToUDP(v.Pack, v.Addr)
        if e != nil && t.log != nil {
            t.log.Println(e)
        }
    }
}

func (c *BeaconChannel) Send(targetHost string, buf []byte) error {
    if len(targetHost) == 0 || len(buf) == 0 {
        return fmt.Errorf("[ERR] Cannot send null data to null host")
    }
    targetAddr := &net.UDPAddr{
        IP      : net.ParseIP(targetHost),
        Port    : PAGENT_RECV_PORT,
    }
    c.ChWrite <- &ChanPkg{
        Pack    : buf,
        Addr    : targetAddr,
    }
    return nil
}
