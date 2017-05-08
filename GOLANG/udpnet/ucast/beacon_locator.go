package ucast

import (
    "net"
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

type BeaconLocator struct {
    isClosed     bool
    waiter       sync.WaitGroup
    chClosed     chan bool

    conn         *net.UDPConn
    ChRead       chan BeaconPack
    chWrite      chan BeaconPack
}

// New constructor of a new server
func NewBeaconLocator() (*BeaconLocator, error) {
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: PAGENT_RECV_PORT})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    conn.SetReadBuffer(PC_MAX_UCAST_UDP_BUF_SIZE)
    conn.SetWriteBuffer(PC_MAX_UCAST_UDP_BUF_SIZE)
    locator := &BeaconLocator{
        isClosed:    false,
        chClosed:    make(chan bool),
        conn:        conn,
        ChRead:      make(chan BeaconPack),
        chWrite:     make(chan BeaconPack),
    }
    locator.waiter.Add(2)
    go locator.read()
    go locator.write()
    return locator, nil
}

// Close is used to cleanup the client
func (lc *BeaconLocator) Close() error {
    if lc.isClosed {
        return nil
    }
    lc.isClosed = true
    err := lc.conn.Close()

    // broad case close action
    close(lc.chClosed)
    // wait...
    lc.waiter.Wait()
    // then close channels
    close(lc.ChRead)
    close(lc.chWrite)
    return err
}

func (lc *BeaconLocator) read() {
    var (
        buff []byte            = make([]byte, PC_MAX_UCAST_UDP_BUF_SIZE)
        addr *net.UDPAddr      = nil
        err error              = nil
        count int              = 0
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

    for !lc.isClosed {
        /*** Set a deadline for reading. Read operation will fail if no data is received after deadline. ***/
        //lc.conn.SetReadDeadline(time.Now().Add(readTimeout))

        count, addr, err = lc.conn.ReadFromUDP(buff)
        if err != nil {
            continue
        }
        adr := copyUDPAddr(addr)
        msg := make([]byte, count)
        copy(msg, buff[:count])
        lc.ChRead <- BeaconPack{Address:adr,Message:msg}
    }

    lc.waiter.Done()
}

func (lc *BeaconLocator) write() {
    defer lc.waiter.Done()

    for {
        select {
            case <- lc.chClosed: {
                return
            }
            case v, ok := <- lc.chWrite: {
                if !ok {
                    return
                }
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
}

func (lc *BeaconLocator) Send(targetHost string, buf []byte) error {
    if len(targetHost) == 0 {
        return errors.Errorf("[ERR] BeaconLocator.Send() :: cannot send data to null host")
    }
    if len(buf) == 0 {
        return errors.Errorf("[ERR] BeaconLocator.Send() :: cannot send null data")
    }
    if lc.isClosed {
        return nil
    }

    select {
        case <-lc.chClosed: {
            return nil
        }
        case lc.chWrite <- BeaconPack{
            Message: buf,
            Address: net.UDPAddr{
                IP:      net.ParseIP(targetHost),
                Port:    PAGENT_SEND_PORT,
            },
        }: {
            // TODO : find ways to remove this. We'll wait artificially for now (v0.1.4)
            time.After(time.Millisecond)
        }
    }

    return nil
}
