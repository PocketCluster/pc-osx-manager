package ucast

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-node-agent/service"
)

const (
    EventBeaconNodeLocationSend string    = "event.beacon.node.location.send"
    EventBeaconNodeLocationReceive string = "event.beacon.node.location.receive"
    iventBeaconNodeServiceClose string    = "ivent.beacon.node.service.close"
)

type BeaconAgent struct {
    service.AppSupervisor
    conn *net.UDPConn
}

// New constructor of a new server
func NewBeaconAgent(aup service.AppSupervisor) (*BeaconAgent, error) {
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: POCKET_AGENT_PORT})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    conn.SetReadBuffer(PC_MAX_UCAST_UDP_BUF_SIZE)
    conn.SetWriteBuffer(PC_MAX_UCAST_UDP_BUF_SIZE)
    agent := &BeaconAgent{
        AppSupervisor:    aup,
        conn:             conn,
    }
    agent.read()
    agent.write()
    return agent, nil
}

// Close is used to cleanup the client
func (b *BeaconAgent) Close() error {
    b.BroadcastEvent(service.Event{Name:iventBeaconNodeServiceClose})
    return nil
}

func (b *BeaconAgent) read() {
    var(
        closedC = make(chan service.Event)
    )
    b.RegisterServiceWithFuncs(
        func() error {
            var (
                buff []byte          = make([]byte, PC_MAX_UCAST_UDP_BUF_SIZE)
                addr *net.UDPAddr    = nil
                err error            = nil
                count int            = 0
            )
            for {
                select {
                    case <- b.StopChannel(): {
                        err = b.conn.Close()
                        return errors.WithStack(err)
                    }
                    case <- closedC: {
                        err = b.conn.Close()
                        return errors.WithStack(err)
                    }
                    default: {
                        // Set a deadline for reading. Read operation will fail if no data
                        // is received after deadline.
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
                                Name:       EventBeaconNodeLocationReceive,
                                Payload:    BeaconPack{
                                    Address:    adr,
                                    Message:    msg},
                            },
                        )
                    }
                }
            }
        },
        func(_ func(interface{})) error {
            return nil
        },
        service.BindEventWithService(iventBeaconNodeServiceClose, closedC))
}

func (b *BeaconAgent) write() {
    var(
        eventC  chan service.Event = make(chan service.Event)
        closedC chan service.Event = make(chan service.Event)
    )
    b.RegisterServiceWithFuncs(
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
                            Port:    POCKET_LOCATOR_PORT,
                        }
                        _, err := b.conn.WriteToUDP(bs.Payload, addr)
                        if err != nil {
                            log.Debugf("[ERR] beacon send error %v", errors.WithStack(err))
                        }
                    }
                }
            }
        },
        func(_ func(interface{})) error {
            return nil
        },
        service.BindEventWithService(EventBeaconNodeLocationSend, eventC),
        service.BindEventWithService(iventBeaconNodeServiceClose, closedC))
}
