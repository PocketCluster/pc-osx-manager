package mcast

import (
    "net"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/event/operation"
)

const (
    EventBeaconCoreSearchReceive string      = "event.beacon.core.search.receive"
    iventBeaconCoreSearchCatcherClose string = "ivent.beacon.core.search.catcher.close"
)

type SearchCatcher struct {
    service.ServiceSupervisor
    conn *net.UDPConn
}

func NewSearchCatcher(sup service.ServiceSupervisor, niface string) (*SearchCatcher, error) {
    iface, err := net.InterfaceByName(niface)
    if err != nil {
        log.Error(err)
        return nil, errors.WithStack(err)
    }
    conn, err := net.ListenMulticastUDP("udp4", iface, ipv4McastAddr)
    if err != nil {
        return nil, errors.Errorf("[ERR] failed to bind to any multicast udp port", err)
    }
    conn.SetReadBuffer(PC_MAX_MCAST_UDP_BUF_SIZE)
    listener := &SearchCatcher{
        ServiceSupervisor:    sup,
        conn:                 conn,
    }
    listener.read()
    return listener, nil
}

// Close is used to cleanup the client
func (s *SearchCatcher) Close() error {
    s.BroadcastEvent(service.Event{Name:iventBeaconCoreSearchCatcherClose})
    return nil
}

func (s *SearchCatcher) read() {
    var (
        closedC = make(chan service.Event)
    )
    s.RegisterServiceWithFuncs(
        operation.ServiceBeaconCatcher,
        func() error {
            var (
                buff []byte          = make([]byte, PC_MAX_MCAST_UDP_BUF_SIZE)
                addr *net.UDPAddr    = nil
                err error            = nil
                count int            = 0
            )
            for {
                select {
                    case <-s.StopChannel(): {
                        err = s.conn.Close()
                        return errors.WithStack(err)
                    }
                    case <- closedC: {
                        err = s.conn.Close()
                        return errors.WithStack(err)
                    }
                    default: {
                        /*** Set a deadline for reading. Read operation will fail if no data is received after deadline. ***/
                        err = s.conn.SetReadDeadline(time.Now().Add(readTimeout))
                        if err != nil {
                            continue
                        }

                        count, addr, err = s.conn.ReadFromUDP(buff)
                        if err != nil {
                            continue
                        }
                        if count == 0 {
                            continue
                        }
                        adr := copyUDPAddr(addr)
                        msg := make([]byte, count)
                        copy(msg, buff[:count])
                        s.BroadcastEvent(
                            service.Event{
                                Name:       EventBeaconCoreSearchReceive,
                                Payload:    CastPack{
                                    Address:    adr,
                                    Message:    msg},
                            },
                        )
                    }
                }
            }

            return nil
        },
        service.BindEventWithService(iventBeaconCoreSearchCatcherClose, closedC))
}