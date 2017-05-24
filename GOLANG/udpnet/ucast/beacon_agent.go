package ucast

import (
    "net"
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

type BeaconAgent struct {
    isClosed     bool
    waiter       sync.WaitGroup
    chClosed     chan bool

    conn         *net.UDPConn
    ChRead       chan BeaconPack
    chWrite      chan BeaconPack
}

// New constructor of a new server
func NewBeaconAgent() (*BeaconAgent, error) {
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: POCKET_AGENT_PORT})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    conn.SetReadBuffer(PC_MAX_UCAST_UDP_BUF_SIZE)
    conn.SetWriteBuffer(PC_MAX_UCAST_UDP_BUF_SIZE)
    agent := &BeaconAgent{
        isClosed:    false,
        chClosed:    make(chan bool),

        conn:        conn,
        ChRead:      make(chan BeaconPack),
        chWrite:     make(chan BeaconPack),
    }
    agent.waiter.Add(2)
    go agent.reader()
    go agent.writer()
    return agent, nil
}

// Close is used to cleanup the client
func (bc *BeaconAgent) Close() error {
    if bc.isClosed {
        return nil
    }
    bc.isClosed = true
    err := bc.conn.Close()

    close(bc.chClosed)
    bc.waiter.Wait()
    close(bc.ChRead)
    close(bc.chWrite)
    return err
}

func (bc *BeaconAgent) reader() {
    var (
        buff []byte          = make([]byte, PC_MAX_UCAST_UDP_BUF_SIZE)
        addr *net.UDPAddr    = nil
        err error            = nil
        count int            = 0
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

    for !bc.isClosed {
        // Set a deadline for reading. Read operation will fail if no data
        // is received after deadline.
        //bc.conn.SetReadDeadline(time.Now().Add(readTimeout))

        count, addr, err = bc.conn.ReadFromUDP(buff)
        if err != nil {
            continue
        }
        adr := copyUDPAddr(addr)
        msg := make([]byte, count)
        copy(msg, buff[:count])
        bc.ChRead <- BeaconPack{Address:adr,Message:msg}
    }

    bc.waiter.Done()
}

func (bc *BeaconAgent) writer() {
    defer bc.waiter.Done()

    for {
        select {
            case <- bc.chClosed: {
                return
            }
            case v, ok := <- bc.chWrite: {
                if !ok {
                    return
                }
                if len(v.Message) == 0 {
                    continue
                }
                _, err := bc.conn.WriteToUDP(v.Message, &v.Address)
                if err != nil {
                    log.Info(err)
                }
            }
        }
    }
}

func (bc *BeaconAgent) Send(targetHost string, buf []byte) error {
    if len(targetHost) == 0 || len(buf) == 0 {
        return errors.Errorf("[ERR] Cannot send null data to null host")
    }
    if bc.isClosed {
        return nil
    }

    select {
        case <- bc.chClosed: {
            return nil
        }
        case  bc.chWrite <- BeaconPack{
            Address: net.UDPAddr {
                IP:      net.ParseIP(targetHost),
                Port:    POCKET_LOCATOR_PORT,
            },
            Message: buf,
        }: {
            // TODO : find ways to remove this. We'll wait artificially for now (v0.1.4)
            time.After(time.Millisecond)
        }
    }

    return nil
}
