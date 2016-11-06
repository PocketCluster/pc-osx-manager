package ucast

import (
    "net"
    "sync"
    "fmt"
    "log"
    "time"
)

const LOCATOR_RECVD_BUFFER_SIZE int = 16384

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
    buf := make([]byte, LOCATOR_RECVD_BUFFER_SIZE)
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