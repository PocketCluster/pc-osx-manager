package ucast

import (
    "net"
    "sync"
    "fmt"
    "log"
    //"time"
    "time"
)

/*
type locatorChannel struct {
    ipv4Conn     *net.UDPConn

    closed       bool
    closedCh     chan struct{}
    closeLock    sync.Mutex
}

func NewPocketLocatorChannel() (*locatorChannel, error) {
    // Create a IPv4 listener
    urecv4, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: PAGENT_RECV_PORT})
    if err != nil {
        return nil, fmt.Errorf("[ERR] failed to bind to any unicast udp port : " + err.Error())
    }
    c := &locatorChannel{
        ipv4Conn    : urecv4,
        closedCh    : make(chan struct{}),
    }
    return c, nil
}

// query is used to perform a lookup and stream results
func (c *locatorChannel) Connect(param *ConnParam) error {

    // Start listening for response packets
    go c.recv(c.ipv4Conn, param.RecvMessage)

    // Listen until we reach the timeout
    finish := time.After(param.Timeout)
    for {
        select {
        case <-finish:
            return nil
        }
    }
}

// Close is used to cleanup the client
func (c *locatorChannel) Close() error {
    c.closeLock.Lock()
    defer c.closeLock.Unlock()

    if c.closed {
        return nil
    }
    c.closed = true

    log.Printf("[INFO] mdns: Closing client %v", *c)
    close(c.closedCh)

    if c.ipv4Conn != nil {
        c.ipv4Conn.Close()
    }
    return nil
}

func (c *locatorChannel) Send(targetHost string, q []byte) error {
    targetAddr := &net.UDPAddr{
        IP:   net.ParseIP(targetHost),
        Port: PAGENT_SEND_PORT,
    }
    buf := q
    if c.ipv4Conn != nil {
        c.ipv4Conn.WriteToUDP(buf, targetAddr)
    }
    return nil
}

// recv is used to receive until we get a shutdown
func (c *locatorChannel) recv(l *net.UDPConn, msgChan chan <- []byte) {
    if l == nil {
        return
    }
    buf := make([]byte, PC_MAX_UDP_BUF_SIZE)
    for !c.closed {
        n, err := l.Read(buf)
        if err != nil {
            log.Printf("[ERR] mdns: Failed to read packet: %v", err)
            continue
        }

        log.Printf("%d bytes have been received %v", n, buf[:n])
        msg := buf[:n]
        select {
        case msgChan <- msg:
        case <-c.closedCh:
            return
        }
    }
}
*/

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
        ChRead     : make(chan *ChanPkg, PC_LOCATOR_CHAN_CAP),
        chWrite    : make(chan *ChanPkg, PC_LOCATOR_CHAN_CAP),
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
        pack.Pack = make([]byte, PC_MAX_UDP_BUF_SIZE)
        count, pack.Addr, err = lc.conn.ReadFromUDP(pack.Pack)
        if err != nil {
            if lc.log != nil {
                lc.log.Printf("[ERR] locator channel : Failed to read packet: %v", err)
            }
            continue
        }

        if lc.log != nil {
            lc.log.Printf("[INFO] %d bytes have been received", count)
        }
        pack.Pack = pack.Pack[:count]
        select {
        case lc.ChRead <- pack:
        case <-lc.closedCh:
            return
        }
    }
}

func (lc *locatorChannel) writer() {
    for v := range lc.chWrite {
        _, e := lc.conn.WriteToUDP(v.Pack, v.Addr)
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
        Pack    : buf,
        Addr    : targetAddr,
    }

    // TODO : find ways to remove this. We'll wait artificially for now (v0.1.4)
    time.Sleep(time.Microsecond * 100)
    return nil
}
