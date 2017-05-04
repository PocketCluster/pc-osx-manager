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
    closeLock    sync.Mutex

    conn         *net.UDPConn
    waiter       *sync.WaitGroup
    ChRead       chan BeaconPack
    chWrite      chan BeaconPack
}

// New constructor of a new server
func NewBeaconAgent(waiter *sync.WaitGroup) (*BeaconAgent, error) {
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: PAGENT_SEND_PORT})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    conn.SetReadBuffer(PC_MAX_UCAST_UDP_BUF_SIZE)
    conn.SetWriteBuffer(PC_MAX_UCAST_UDP_BUF_SIZE)
    beacon := &BeaconAgent{
        isClosed:    false,

        conn:       conn,
        waiter:     waiter,
        ChRead:     make(chan BeaconPack, PC_UCAST_BEACON_CHAN_CAP),
        chWrite:    make(chan BeaconPack, PC_UCAST_BEACON_CHAN_CAP),
    }
    waiter.Add(2)
    go beacon.reader()
    go beacon.writer()
    return beacon, nil
}

// Close is used to cleanup the client
func (bc *BeaconAgent) Close() error {
    bc.closeLock.Lock()
    defer bc.closeLock.Unlock()

    if bc.isClosed {
        return nil
    }
    bc.isClosed = true
    close(bc.ChRead)
    close(bc.chWrite)
    return bc.conn.Close()
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
    defer bc.waiter.Done()

    for !bc.isClosed {
        // Set a deadline for reading. Read operation will fail if no data
        // is received after deadline.
        //bc.conn.SetReadDeadline(time.Now().Add(readTimeout))

        count, addr, err = bc.conn.ReadFromUDP(buff)
        if err != nil {
            continue
        }
        if count == 0 {
            continue
        }
        adr := copyUDPAddr(addr)
        msg := make([]byte, count)
        copy(msg, buff[:count])
        pack := BeaconPack{
            Address:    adr,
            Message:    msg,
        }
        bc.ChRead <- pack
    }
}

func (bc *BeaconAgent) writer() {
    defer bc.waiter.Done()

    for v := range bc.chWrite {
        if len(v.Message) == 0 {
            continue
        }
        _, err := bc.conn.WriteToUDP(v.Message, &v.Address)
        if err != nil {
            log.Info(err)
        }
    }

}

func (bc *BeaconAgent) Send(targetHost string, buf []byte) error {
    if len(targetHost) == 0 || len(buf) == 0 {
        return errors.Errorf("[ERR] Cannot send null data to null host")
    }
    bc.chWrite <- BeaconPack{
        Message: buf,
        Address: net.UDPAddr {
            IP:      net.ParseIP(targetHost),
            Port:    PAGENT_SEND_PORT,
        },
    }

    // TODO : find ways to remove this. We'll wait artificially for now (v0.1.4)
    time.After(time.Millisecond)
    return nil
}
