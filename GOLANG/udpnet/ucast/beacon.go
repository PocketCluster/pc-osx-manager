package ucast

import (
    "net"
    "sync"
    "fmt"
    "log"
    "time"
)

type beaconChannel struct {
    ipv4Conn     *net.UDPConn

    closed       bool
    closedCh     chan struct{}
    closeLock    sync.Mutex
}

func NewPocketBeaconChannel() (*beaconChannel, error) {
    // Create a IPv4 listener
    urecv4, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: PAGENT_SEND_PORT})
    if err != nil {
        return nil, fmt.Errorf("[ERR] failed to bind to any unicast udp port : " + err.Error())
    }
    c := &beaconChannel{
        ipv4Conn    : urecv4,
        closedCh    : make(chan struct{}),
    }
    return c, nil
}

// query is used to perform a lookup and stream results
func (c *beaconChannel) Connect(param *ConnParam) error {

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
func (c *beaconChannel) Close() error {
    c.closeLock.Lock()
    defer c.closeLock.Unlock()

    if c.closed {
        return nil
    }
    c.closed = true

    log.Printf("[INFO] beacon channel closing : %v", *c)
    close(c.closedCh)

    if c.ipv4Conn != nil {
        c.ipv4Conn.Close()
    }
    return nil
}

func (c *beaconChannel) Send(targetHost string, buf []byte) error {
    targetAddr := &net.UDPAddr{
        IP      : net.ParseIP(targetHost),
        Port    : PAGENT_RECV_PORT,
    }
    if c.ipv4Conn != nil {
        c.ipv4Conn.WriteToUDP(buf, targetAddr)
    }
    return nil
}

// recv is used to receive until we get a shutdown
func (c *beaconChannel) recv(l *net.UDPConn, msgCh chan <- []byte) {
    if l == nil {
        return
    }
    buf := make([]byte, PC_MAX_UDP_BUF_SIZE)
    for !c.closed {
        n, err := l.Read(buf)
        if err != nil {
            log.Printf("[ERR] beacon channel : Failed to read packet: %v", err)
            continue
        }

        log.Printf("[INFO] %d bytes have been received %v", n, buf[:n])
        msg := buf[:n]
        select {
        case msgCh <- msg:
        case <-c.closedCh:
            return
        }
    }
}