package mcast

import (
    "net"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-node-agent/service"
)

const (
    EventBeaconNodeSearchSend string        = "event.beacon.node.search.send"
    iventBeaconNodeSearchCasterClose string = "ivent.beacon.node.search.caster.close"
)

type SearchCaster struct {
    service.AppSupervisor
    conn *net.UDPConn
}

func NewSearchCaster(aup service.AppSupervisor) (*SearchCaster, error) {
    // Create a IPv4 listener
    conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
    if err != nil {
        return nil, errors.WithStack(err)
    }
    conn.SetWriteBuffer(PC_MAX_MCAST_UDP_BUF_SIZE)
    mc := &SearchCaster{
        AppSupervisor:    aup,
        conn:             conn,
    }
    mc.write()
    return mc, nil
}

// Close is used to cleanup the client
func (s *SearchCaster) Close() error {
    s.BroadcastEvent(service.Event{Name:iventBeaconNodeSearchCasterClose})
    return nil
}

func (s *SearchCaster) write() {
    var(
        eventC  chan service.Event = make(chan service.Event)
        closedC chan service.Event = make(chan service.Event)
    )
    s.RegisterServiceWithFuncs(
        func() error {
            var (
                err error = nil
            )
            for {
                select {
                    case <- s.StopChannel(): {
                        err = s.conn.Close()
                        return errors.WithStack(err)
                    }
                    case <- closedC: {
                        err = s.conn.Close()
                        return errors.WithStack(err)
                    }
                    case e := <- eventC: {
                        cpkg, ok := e.Payload.(CastPack)
                        if !ok {
                            log.Debugf("[WARN] invalid SearchCaster type")
                            continue
                        }
                        if len(cpkg.Message) == 0 {
                            log.Debugf("[WARN] empty SearchCaster message")
                            continue
                        }
                        _, err = s.conn.WriteToUDP(cpkg.Message, ipv4McastAddr)
                        if err != nil {
                            log.Debugf("[WARN] SearchCaster transmit. Error : %v", errors.WithStack(err))
                        }
                    }
                }
            }
        },
        func(_ func(interface{})) error {
            return nil
        },
        service.BindEventWithService(EventBeaconNodeSearchSend,        eventC),
        service.BindEventWithService(iventBeaconNodeSearchCasterClose, closedC))
}
