package main

import (
    "net"
    "os"
    "time"

    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/gravitational/teleport/embed"
    tervice "github.com/gravitational/teleport/lib/service"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/pc-node-agent/dhcp"
    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/service"
)

import (
    "github.com/davecgh/go-spew/spew"
)

const (
    nodeServiceSearch  = "service_search"
    nodeServiceBeacon  = "service_beacon"
    nodeFeedbackBeacon = "feedback_beacon"
    nodeFeedbackDHCP   = "feedback_dhcp"
    nodeTeleportStart  = "teleport_start"
    nodeTeleportStop   = "teleport_stop"
)

func initDhcpListner(app service.AppSupervisor) error {
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

            app.BroadcastEvent(service.Event{Name:nodeFeedbackDHCP, Payload:dhcpEvent})

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

func initSearchService(app service.AppSupervisor) error {
    caster, err := mcast.NewSearchCaster()
    if err != nil {
        return err
    }
    eventsC := make(chan service.Event)
    app.WaitForEvent(nodeServiceSearch, eventsC, make(chan struct{}))

    app.RegisterFunc(func() error {
        log.Debugf("[SEARCH] starting master serach service...")

        for {
            select {
                case <- app.StopChannel():
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

func initBeaconService(app service.AppSupervisor) error {
    beacon, err := ucast.NewBeaconAgent()
    if err != nil {
        return err
    }
    eventsC := make(chan service.Event)
    app.WaitForEvent(nodeServiceBeacon, eventsC, make(chan struct{}))

    app.RegisterFunc(func() error {
        for {
            select {
                case <- app.StopChannel():
                    return nil
                case v := <- beacon.ChRead: {
                    app.BroadcastEvent(service.Event{Name:nodeFeedbackBeacon, Payload:v})
                }
            }
        }
        return nil
    })

    app.RegisterFunc(func() error {
        log.Debugf("[BEACON] starting beacon service...")

        for {
            select {
                case <-app.StopChannel():
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

func initTeleportNodeService(app service.AppSupervisor) error {
    var (
        nodeProc *embed.EmbeddedNodeProcess = nil
        startC = make(chan service.Event)
        stopC = make(chan service.Event)
        startTeleport = func () error {
            maddr, err := slcontext.SharedSlaveContext().GetMasterIP4Address()
            if err != nil {
                log.Errorf(err.Error())
                return errors.WithStack(err)
            }
            slid, err := slcontext.SharedSlaveContext().GetSlaveNodeUUID()
            if err != nil {
                log.Errorf(err.Error())
                return errors.WithStack(err)
            }
            cfg, err := tervice.MakeNodeConfig(maddr, slid, true)
            if err != nil {
                log.Errorf(err.Error())
                return errors.WithStack(err)
            }
            nodeProc, err = embed.NewEmbeddedNodeProcess(app, cfg)
            if err != nil {
                return errors.WithStack(err)
            }
            err = nodeProc.StartNodeSSH()
            if err != nil {
                log.Errorf(err.Error())
                errors.WithStack(err)
            }

            log.Debugf("\n\n(INFO) teleport node started success!\n")
            return nil
        }
    )
    app.WaitForEvent(nodeTeleportStart, startC, make(chan struct{}))
    app.WaitForEvent(nodeTeleportStop,   stopC, make(chan struct{}))

    app.RegisterFunc(func() error{

        for {
            select {
                case <-startC: {
                    return startTeleport()
                }
                case <-stopC: {
                    return startTeleport()
                }
            }
        }

        return nil
    })

    return nil
}

func initAgentService(app service.AppSupervisor) error {
    var (

        beaconC = make(chan service.Event)
        dhcpC   = make(chan service.Event)

        searchTx = func(data []byte) error {
            log.Debugf("[SEARCH-TX] %v", time.Now())
            app.BroadcastEvent(service.Event{Name: nodeServiceSearch, Payload:data})
            return nil
        }

        beaconTx = func(target string, data []byte) error {
            log.Debugf("[BEACON-TX] %v TO : %v", time.Now(), target)
            app.BroadcastEvent(service.Event{
                Name: nodeServiceBeacon,
                Payload: ucast.BeaconSend{
                    Host:target,
                    Payload:data,
                },
            })
            return nil
        }

        transitEvent = func (state locator.SlaveLocatingState, ts time.Time, transOk bool) error {
            if transOk {
                log.Debugf("(INFO) [%v] BeaconEventTranstion -> %v | SUCCESS ", ts, state.String())
            } else {
                log.Debugf("(INFO) [%v] BeaconEventTranstion -> %v | FAILED ", ts, state.String())
            }

            if transOk {
                switch state {
                    case locator.SlaveUnbounded: {
                        return nil
                    }
                    case locator.SlaveInquired: {
                        return nil
                    }
                    case locator.SlaveKeyExchange: {
                        return nil
                    }
                    case locator.SlaveCryptoCheck: {
                        app.BroadcastEvent(service.Event{Name:nodeTeleportStart})
                        return nil
                    }
                    case locator.SlaveBounded: {
                        return nil
                    }
                    case locator.SlaveBindBroken: {
                        app.BroadcastEvent(service.Event{Name:nodeTeleportStart})
                        return nil
                    }
                }

            } else {
                switch state {
                    case locator.SlaveUnbounded: {
                        return nil
                    }
                    case locator.SlaveInquired: {
                        return nil
                    }
                    case locator.SlaveKeyExchange: {
                        return nil
                    }
                    case locator.SlaveCryptoCheck: {
                        return nil
                    }
                    case locator.SlaveBounded: {
                        return nil
                    }
                    case locator.SlaveBindBroken: {
                        return nil
                    }
                }
            }

            return nil
        }

        serviceFunc = func() error {
            var (
                timer = time.NewTicker(time.Second)
                context = slcontext.SharedSlaveContext()
                loc locator.SlaveLocator = nil
                locState locator.SlaveLocatingState = locator.SlaveUnbounded
                err error = nil
            )

            // setup slave locator
            uuid, err := context.GetSlaveNodeUUID()
            if err == nil && len(uuid) != 0 {
                locState = locator.SlaveBindBroken
            } else {
                locState = locator.SlaveUnbounded
            }
            loc, err = locator.NewSlaveLocatorWithFunc(locState, searchTx, beaconTx, transitEvent)
            if err != nil {
                return errors.WithStack(err)
            }
            defer loc.Shutdown()
            defer timer.Stop()

            log.Debugf("[AGENT] starting agent service...")

            for {
                select {
                    case <- app.StopChannel():
                        return nil
                    case b := <- beaconC: {
                        mp, ok := b.Payload.(ucast.BeaconPack)
                        if ok {
                            err = loc.TranstionWithMasterBeacon(mp, time.Now())
                            if err != nil {
                                log.Debug(err.Error())
                            }
                        }
                    }
                    case d := <- dhcpC: {
                        log.Debugf("[DHCP] RECEIVED\n %v", spew.Sdump(d.Payload))
                    }
                    case <- timer.C: {
                        err = loc.TranstionWithTimestamp(time.Now())
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                    }
                }
            }
            return nil
        }
    )

    app.WaitForEvent(nodeFeedbackBeacon, beaconC, make(chan struct{}))
    app.WaitForEvent(nodeFeedbackDHCP, dhcpC, make(chan struct{}))
    app.RegisterFunc(serviceFunc)

    app.OnExit(func(payload interface{}) {
        log.Debugf("[AGENT] close agent service...")
    })

    return nil
}

func main() {
    var (
        err error = nil
        app service.AppSupervisor
    )
    log.SetLevel(log.DebugLevel)

    // TODO check user and reject if not root

    // initialize slave context
    slcontext.SharedSlaveContext()
    app = service.NewAppSupervisor()

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

    // teleport management
    err = initTeleportNodeService(app)
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

