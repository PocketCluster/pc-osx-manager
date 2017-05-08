package main

import (
    "os"
    "net"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "gopkg.in/vmihailenco/msgpack.v2"

    "github.com/stkim1/pc-node-agent/dhcp"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
)

const (
    nodeServiceSearch  = "service_search"
    nodeServiceBeacon  = "service_beacon"
    nodeFeedbackBeacon = "feedback_beacon"
    nodeFeedbackDHCP   = "feedback_dhcp"
)

func initDhcpListner(app *PocketApplication) error {

    app.RegisterFunc(func () error {
        log.Debugf("[DHCP] dhcp listner started...")

        buf := make([]byte, 20480)
        dhcpEvent := &dhcp.DhcpEvent{}

        // firstly clear off previous socket
        os.Remove(dhcp.DHCPEventSocketPath)
        dhcpListener, err := net.ListenUnix("unix", &net.UnixAddr{dhcp.DHCPEventSocketPath, "unix"})
        if err != nil {
            return errors.WithStack(err)
        }
        defer os.Remove(dhcp.DHCPEventSocketPath)
        defer dhcpListener.Close()

        for {
            conn, err := dhcpListener.AcceptUnix()
            if err != nil {
                log.Error(errors.WithStack(err))
                continue
            }
            count, err := conn.Read(buf)
            if err != nil {
                log.Error(errors.WithStack(err))
                continue
            }
            err = msgpack.Unmarshal(buf[0:count], dhcpEvent)
            if err != nil {
                log.Error(errors.WithStack(err))
                continue
            }

            app.BroadcastEvent(Event{Name:nodeFeedbackDHCP, Payload:dhcpEvent})

            err = conn.Close()
            if err != nil {
                log.Error(errors.WithStack(err))
                continue
            }
        }

        return nil
    })

    return nil
}

func initSearchService(app *PocketApplication) error {
    caster, err := mcast.NewSearchCaster()
    if err != nil {
        return err
    }
    eventsC := make(chan Event)
    app.WaitForEvent(nodeServiceSearch, eventsC, make(chan struct{}))

    app.RegisterFunc(func() error {
        log.Debugf("[SEARCH] starting master serach service...")

        for {
            select {
                case e := <- eventsC: {
                    cm, ok := e.Payload.([]byte)
                    if ok {
                        log.Debugf("[SEARCH] casting message...")
                        err := caster.Send(cm)
                        if err != nil {
                            log.Errorf("[SEARCH] casting error %v", err)
                        }
                    }
                }
            }
        }

        return caster.Close()
    })

    return nil
}

func initBeaconService(app *PocketApplication) error {
    beacon, err := ucast.NewBeaconAgent()
    if err != nil {
        return err
    }
    eventsC := make(chan Event)
    app.WaitForEvent(nodeServiceBeacon, eventsC, make(chan struct{}))

    app.RegisterFunc(func() error {
        log.Debugf("[BEACON] starting beacon service...")
        go func() {
            for v := range beacon.ChRead {
                log.Debugf("[BEACON] message received %v", v)
            }
        }()

        for {
            select {
                case e := <- eventsC: {
                    bs, ok := e.Payload.(ucast.BeaconSend)
                    if ok {
                        log.Debugf("[BEACON] sending message %v", bs)
                        beacon.Send(bs.Host, bs.Payload)
                    }
                }

            }
        }

        return beacon.Close()
    })
    return nil
}

func initAgentService(app *PocketApplication) error {
    beaconC := make(chan Event)
    dhcpC := make(chan Event)
    app.WaitForEvent(nodeFeedbackBeacon, beaconC, make(chan struct{}))
    app.WaitForEvent(nodeFeedbackDHCP, dhcpC, make(chan struct{}))

    app.RegisterFunc(func() error {
        log.Debugf("[AGENT] starting agent service...")
        for {
            select {
                case b := <- beaconC: {
                    log.Debugf("[AGENT] beacon recieved %v", b)
                }
                case d := <- dhcpC: {
                    log.Debugf("[AGENT] dhcp recieved %v", d)
                }
            }
        }
        log.Debugf("[AGENT] finishing agent service...")
        return nil
    })
    return nil
}

func main() {
    var (
        err error = nil
        app *PocketApplication
    )
    log.SetLevel(log.DebugLevel)
    app = NewPocketApplication()

    // dhcp listner
    err = initDhcpListner(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // search service
    err = initSearchService(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // beacon service
    err = initBeaconService(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // agent service
    err = initAgentService(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // application
    err = app.Start()
    if err != nil {
        log.Panic(errors.WithStack(err))
    }
    app.Wait()
}

