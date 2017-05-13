package main

import (
    "net"
    "os"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "gopkg.in/vmihailenco/msgpack.v2"

    "github.com/stkim1/pc-node-agent/dhcp"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
)

import (
    "github.com/davecgh/go-spew/spew"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/locator"
)

const (
    nodeServiceSearch  = "service_search"
    nodeServiceBeacon  = "service_beacon"
    nodeFeedbackBeacon = "feedback_beacon"
    nodeFeedbackDHCP   = "feedback_dhcp"
)

func initDhcpListner(app *PocketApplication) error {
    // firstly clear off previous socket
    os.Remove(dhcp.DHCPEventSocketPath)
    dhcpListener, err := net.ListenUnix("unix", &net.UnixAddr{dhcp.DHCPEventSocketPath, "unix"})
    if err != nil {
        return errors.WithStack(err)
    }

    app.RegisterFunc(func () error {
        log.Debugf("[DHCP] starting dhcp listner...")
        buf := make([]byte, 20480)
        dhcpEvent := &dhcp.DhcpEvent{}

        // TODO : how do we stop this?
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

    app.OnExit(func(payload interface{}) {
        dhcpListener.Close()
        os.Remove(dhcp.DHCPEventSocketPath)
        log.Debugf("[DHCP] close dhcp listner...")
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
                case <- app.stoppedC:
                    return nil
                case e := <-eventsC: {
                    cm, ok := e.Payload.([]byte)
                    if ok {
//                        log.Debugf("[SEARCH] casting message... %v", cm)
                        err := caster.Send(cm)
                        if err != nil {
                            log.Errorf("[SEARCH] casting error %v", err)
                        }
                    }
                }
            }
        }

        return nil
    })

    app.OnExit(func(payload interface{}) {
        caster.Close()
        log.Debugf("[SEARCH] close master serach service...")
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
        for {
            select {
                case <- app.stoppedC:
                    return nil
                case v := <- beacon.ChRead: {
//                    log.Debugf("[BEACON] message received %v", v)
                    app.BroadcastEvent(Event{Name:nodeFeedbackBeacon, Payload:v})
                }
            }
        }
        return nil
    })

    app.RegisterFunc(func() error {
        log.Debugf("[BEACON] starting beacon service...")

        for {
            select {
                case <-app.stoppedC:
                    return nil
                case e := <- eventsC: {
                    bs, ok := e.Payload.(ucast.BeaconSend)
                    if ok {
//                        log.Debugf("[BEACON] sending message %v", bs)
                        beacon.Send(bs.Host, bs.Payload)
                    }
                }
            }
        }

        return nil
    })

    app.OnExit(func(payload interface{}) {
        beacon.Close()
        log.Debugf("[BEACON] close beacon service...")
    })

    return nil
}

func initAgentService(app *PocketApplication) error {
    var (
        beaconC = make(chan Event)
        dhcpC   = make(chan Event)
    )
    app.WaitForEvent(nodeFeedbackBeacon, beaconC, make(chan struct{}))
    app.WaitForEvent(nodeFeedbackDHCP, dhcpC, make(chan struct{}))

    app.RegisterFunc(func() error {
        var (
            unbounded = time.NewTicker(time.Second * 5)
            bounded = time.NewTicker(time.Second * 10)
            context = slcontext.SharedSlaveContext()
            loc locator.SlaveLocator = nil
            err error = nil
        )
        defer unbounded.Stop()
        defer bounded.Stop()

        // setup slave locator
        uuid, err := context.GetSlaveNodeUUID()
        if err == nil && len(uuid) != 0 {
            loc, err = locator.NewSlaveLocator(locator.SlaveBindBroken, nil)
        } else {
            loc, err = locator.NewSlaveLocator(locator.SlaveUnbounded, nil)
        }
        if err != nil {
            return errors.WithStack(err)
        }
        defer loc.Close()

        log.Debugf("[AGENT] starting agent service...")

        for {
            select {
                case <- app.stoppedC:
                    return nil
                case b := <- beaconC: {
                    mp, ok := b.Payload.(ucast.BeaconPack)
                    if ok {
                        mup, err := msagent.UnpackedMasterMeta(mp.Message)
                        if err == nil {
                            log.Debugf("[AGENT-BEACON] RECEIVED\n %v \n %v", spew.Sdump(mp.Address), spew.Sdump(mup))
                        }
                    }
                }
                case d := <- dhcpC: {
                    log.Debugf("[AGENT-DHCP] RECEIVED\n %v", spew.Sdump(d.Payload))
                }
                case <- unbounded.C: {
//                    log.Debugf("[AGENT] unbounded %v", time.Now())
                    pums, err := slagent.SlavePackedUnboundedMasterSearch()
                    if err != nil {
                        log.Debugf("[AGENT-UNBOUNDED] SlavePackedUnboundedMasterSearch error %v", err)
                        continue
                    }
                    app.BroadcastEvent(Event{Name: nodeServiceSearch, Payload:pums})
                }
                case <- bounded.C: {
//                    log.Debugf("[AGENT] bounded %v", time.Now())
                    pbm, err := slagent.SlavePackedBindBrokenSearch("PC-MASTER")
                    if err != nil {
                        log.Debugf("[AGENT-BOUNDED] SlavePackedBindBrokenSearch error %v", err)
                        continue
                    }
                    app.BroadcastEvent(Event{
                        Name: nodeServiceBeacon,
                        Payload: ucast.BeaconSend{
                            Host:"192.168.1.105",
                            Payload:pbm,
                        },
                    })
                }
            }
        }
        return nil
    })

    app.OnExit(func(payload interface{}) {
        log.Debugf("[AGENT] close agent service...")
    })

    return nil
}

func main() {
    var (
        err error = nil
        app *PocketApplication
    )
    log.SetLevel(log.DebugLevel)

    // TODO check user and reject if not root

    // initialize slave context
    slcontext.SharedSlaveContext()
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

