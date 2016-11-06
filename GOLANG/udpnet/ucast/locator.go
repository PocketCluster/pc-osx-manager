package ucast

import (
    "net"
    "sync"
    "fmt"
    "log"
    "time"
)

const (
    PAGENT_SEND_PORT    = 10060
    PAGENT_RECV_PORT    = 10061
)

type ConnParam struct {
    RecvMessage         chan <- []byte       // Message to recv
    Timeout             time.Duration        // Lookup timeout, default 1 second
}

type channel struct {
    ipv4Conn     *net.UDPConn

    closed       bool
    closedCh     chan struct{}
    closeLock    sync.Mutex
}

func NewPocketLocatorChannel() (*channel, error) {
    // Create a IPv4 listener
    urecv4, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: PAGENT_RECV_PORT})
    if err != nil {
        return nil, fmt.Errorf("[ERR] failed to bind to any unicast udp port : " + err.Error())
    }
    c := &channel{
        ipv4Conn    : urecv4,
        closedCh    : make(chan struct{}),
    }
    return c, nil
}

// query is used to perform a lookup and stream results
func (c *channel) Connect(params *ConnParam) error {

    // Start listening for response packets
    msgCh := make(chan []byte, 32)
    go c.recv(c.ipv4Conn, params.RecvMessage)

    // Listen until we reach the timeout
    finish := time.After(params.Timeout)
    for {
        select {
        case resp := <-msgCh:
            params.RecvMessage <- resp
        case <-finish:
            return nil
        }
    }
}

// Close is used to cleanup the client
func (c *channel) Close() error {
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

func (c *channel) Send(targetHost string, q []byte) error {
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
func (c *channel) recv(l *net.UDPConn, msgCh chan <- []byte) {
    if l == nil {
        return
    }
    buf := make([]byte, 65536)
    for !c.closed {
        n, err := l.Read(buf)
        if err != nil {
            log.Printf("[ERR] mdns: Failed to read packet: %v", err)
            continue
        }

        log.Printf("%d bytes have been received %v", n, buf[:n])
        msg := buf[:n]
        select {
        case msgCh <- msg:
        case <-c.closedCh:
            return
        }
    }
}