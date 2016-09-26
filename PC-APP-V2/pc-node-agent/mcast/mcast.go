package mcast

import (
    "net"
    "sync"
    "log"
    "fmt"
    "time"
)

const (
    ipv4mdns              = "224.0.0.251"
    ipv6mdns              = "ff02::fb"
    mdnsPort              = 5353
)

var (
    ipv4Addr = &net.UDPAddr{
        IP:   net.ParseIP(ipv4mdns),
        Port: mdnsPort,
    }
    ipv6Addr = &net.UDPAddr{
        IP:   net.ParseIP(ipv6mdns),
        Port: mdnsPort,
    }
)

// QueryParam is used to customize how a Lookup is performed
type QueryParam struct {
    SendMessage         []byte               // Message to send
    RecvMessage         chan <- []byte       // Message to recv
    Timeout             time.Duration        // Lookup timeout, default 1 second
    Interface           *net.Interface       // Multicast interface to use
}

// DefaultParams is used to return a default set of QueryParam's
func DefaultParams(sendMessage []byte, recvMessage chan <- []byte) *QueryParam {
    return &QueryParam{
        SendMessage:         sendMessage,
        RecvMessage:         recvMessage,
        Timeout:             time.Second,
    }
}

// Client provides a Multi-Cast to identify PocketCluster Manager
type client struct {
    ipv4UnicastConn     *net.UDPConn
    ipv4MulticastConn   *net.UDPConn

    closed              bool
    closedCh            chan struct{}
    closeLock           sync.Mutex
}

// NewClient creates a new Multi-Cast Client that can be used to identify PocketCluster Manager
func NewClient() (*client, error) {
    // Create a IPv4 listener
    uconn4, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
    if err != nil {
        return nil, fmt.Errorf("failed to bind to any unicast udp port : " + err.Error())
    }
    mconn4, err := net.ListenMulticastUDP("udp4", nil, ipv4Addr)
    if err != nil {
        return nil, fmt.Errorf("failed to bind to any multicast udp port : " + err.Error())
    }
    c := &client{
        ipv4MulticastConn: mconn4,
        ipv4UnicastConn:   uconn4,
        closedCh:          make(chan struct{}),
    }
    return c, nil
}

// Close is used to cleanup the client
func (c *client) Close() error {
    c.closeLock.Lock()
    defer c.closeLock.Unlock()

    if c.closed {
        return nil
    }
    c.closed = true

    log.Printf("[INFO] mdns: Closing client %v", *c)
    close(c.closedCh)

    if c.ipv4UnicastConn != nil {
        c.ipv4UnicastConn.Close()
    }
    if c.ipv4MulticastConn != nil {
        c.ipv4MulticastConn.Close()
    }
    return nil
}

// query is used to perform a lookup and stream results
func (c *client) Query(params *QueryParam) error {

    // Start listening for response packets
    msgCh := make(chan []byte, 32)
    //go c.recv(c.ipv4UnicastConn, msgCh)
    //go c.recv(c.ipv4MulticastConn, msgCh)
    // FIXME : if msgCh is used, some are missed out
    go c.recv(c.ipv4UnicastConn, params.RecvMessage)
    go c.recv(c.ipv4MulticastConn, params.RecvMessage)

    // Send the query
    if err := c.send(params.SendMessage); err != nil {
        return err
    }

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

// sendQuery is used to multicast a query out
func (c *client) send(q []byte) error {
    //TODO : imeplemnt an interface that follows `dns.Msg`
/*
    buf, err := q.Pack()
    if err != nil {
        return err
    }
*/
    buf := q
    if c.ipv4UnicastConn != nil {
        c.ipv4UnicastConn.WriteToUDP(buf, ipv4Addr)
    }
    return nil
}

// recv is used to receive until we get a shutdown
func (c *client) recv(l *net.UDPConn, msgCh chan <- []byte) {
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
        //TODO : Implement an interface that follows below
/*
        msg := new(dns.Msg)
        if err := msg.Unpack(buf[:n]); err != nil {
            log.Printf("[ERR] mdns: Failed to unpack packet: %v", err)
            continue
        }
*/
        select {
            case msgCh <- msg:
            case <-c.closedCh:
                return
        }
    }
}