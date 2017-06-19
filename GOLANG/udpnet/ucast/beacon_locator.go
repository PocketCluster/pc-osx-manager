package ucast

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/event/operation"
)

const (
    EventBeaconCoreLocationReceive string = "event.beacon.core.location.receive"
    EventBeaconCoreLocationSend string    = "event.beacon.core.location.send"
    iventBeaconCoreServiceClose string    = "ivent.beacon.core.service.close"
)

type BeaconLocator struct {
    service.ServiceSupervisor
    conn *net.UDPConn
}

// New constructor of a new server
func NewBeaconLocator(sup service.ServiceSupervisor) (*BeaconLocator, error) {
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: POCKET_LOCATOR_PORT})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    conn.SetReadBuffer(PC_MAX_UCAST_UDP_BUF_SIZE)
    conn.SetWriteBuffer(PC_MAX_UCAST_UDP_BUF_SIZE)
    locator := &BeaconLocator{
        ServiceSupervisor:    sup,
        conn:                 conn,
    }
    locator.read()
    locator.write()
    return locator, nil
}

// Close is used to cleanup the client
func (lc *BeaconLocator) Close() error {
    lc.BroadcastEvent(service.Event{Name:iventBeaconCoreServiceClose})
    return nil
}

func (b *BeaconLocator) read() {
    var (
        closedC chan service.Event = make(chan service.Event)
    )
    b.RegisterServiceWithFuncs(
        operation.ServiceBeaconLocationRead,
        func() error {

            var (
                buff []byte          = make([]byte, PC_MAX_UCAST_UDP_BUF_SIZE)
                addr *net.UDPAddr    = nil
                err error            = nil
                count int            = 0
            )

            for {
                select {
                    case <-b.StopChannel(): {
                        err = b.conn.Close()
                        return errors.WithStack(err)
                    }
                    case <- closedC: {
                        err = b.conn.Close()
                        return errors.WithStack(err)
                    }
                    default: {
                        /*** Set a deadline for reading. Read operation will fail if no data is received after deadline. ***/
                        err = b.conn.SetReadDeadline(time.Now().Add(readTimeout))
                        if err != nil {
                            continue
                        }
                        count, addr, err = b.conn.ReadFromUDP(buff)
                        if err != nil {
                            continue
                        }
                        if count == 0 {
                            continue
                        }
                        adr := copyUDPAddr(addr)
                        msg := make([]byte, count)
                        copy(msg, buff[:count])
                        b.BroadcastEvent(
                            service.Event{
                                Name:       EventBeaconCoreLocationReceive,
                                Payload:    BeaconPack{
                                    Address:    adr,
                                    Message:    msg},
                            })
                    }
                }
            }
        },
        service.BindEventWithService(iventBeaconCoreServiceClose, closedC))
}

func (b *BeaconLocator) write() {
    var (
        eventC  chan service.Event = make(chan service.Event)
        closedC chan service.Event = make(chan service.Event)
    )
    b.RegisterServiceWithFuncs(
        operation.ServiceBeaconLocationWrite,
        func() error {
            for {
                select {
                    case <- b.StopChannel(): {
                        return nil
                    }
                    case <- closedC: {
                        return nil
                    }
                    case e := <- eventC: {
                        bs, ok := e.Payload.(BeaconSend)
                        if !ok {
                            log.Debugf("[WARN] invalid BeaconSend type")
                            continue
                        }
                        if len(bs.Host) == 0 {
                            log.Debugf("[WARN] unable to Beacon invalid host address")
                            continue
                        }
                        if len(bs.Payload) == 0 {
                            log.Debugf("[WARN] empty BeaconSend payload")
                            continue
                        }
                        addr := &net.UDPAddr{
                            IP:      net.ParseIP(bs.Host),
                            Port:    POCKET_AGENT_PORT,
                        }
                        _, err := b.conn.WriteToUDP(bs.Payload, addr)
                        if err != nil {
                            log.Debugf("[ERR] beacon send error %v", errors.WithStack(err))
                        }
                    }
                }
            }
        },
        service.BindEventWithService(EventBeaconCoreLocationSend, eventC),
        service.BindEventWithService(iventBeaconCoreServiceClose,  closedC))
}
