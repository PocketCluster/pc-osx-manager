package main

import (
    "net"
    "os"
    "time"

    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/gravitational/teleport/embed"
    sysd "github.com/coreos/go-systemd/dbus"
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
    dockerServiceUnit  = "docker.service"
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
    )
    app.WaitForEvent(nodeTeleportStart, startC, make(chan struct{}))
    app.WaitForEvent(nodeTeleportStop,   stopC, make(chan struct{}))

    app.RegisterFunc(func() error{

        for {
            select {
                case <- app.StopChannel():
                    return nil

                case noti := <-startC: {
                    // restart teleport
                    cfg, err := tervice.MakeNodeConfig(slcontext.SharedSlaveContext(), true)
                    if err != nil {
                        log.Errorf(err.Error())
                        continue
                    }
                    nodeProc, err = embed.NewEmbeddedNodeProcess(app, cfg)
                    if err != nil {
                        log.Errorf(err.Error())
                        continue
                    }
                    // execute docker engine cert acquisition before SSH node start
                    state, ok := noti.Payload.(locator.SlaveLocatingState)
                    if ok && state == locator.SlaveCryptoCheck {
                        // TODO : create a waitforevent channel and restart docker engine accordingly
                        err = nodeProc.AcquireEngineCertificate(slcontext.DockerEnvironemtPostProcess)
                        if err != nil {
                            log.Errorf(err.Error())
                        }
                    }
                    err = nodeProc.StartNodeSSH()
                    if err != nil {
                        log.Errorf(err.Error())
                        continue
                    }
                    log.Debugf("\n\n(INFO) teleport node started success!\n")

                    continue

                    // restart docker engine
                    // TODO : FIX /opt/gopkg/src/github.com/godbus/dbus/conn.go:345 send on closed channel
                    conn, err := sysd.NewSystemdConnection()
                    if err != nil {
                        log.Errorf(err.Error())
                    } else {
                        did, err := conn.RestartUnit(dockerServiceUnit, "replace", nil)
                        if err != nil {
                            log.Errorf(err.Error())
                        } else {
                            conn.Close()
                            log.Debugf("\n\n(INFO) docker engin restart success! ID %d\n", did)
                        }
                    }
                }

                case <-stopC: {
                    err := nodeProc.Close()
                    if err != nil {
                        log.Errorf(err.Error())
                        continue
                    }
                    nodeProc = nil
                    log.Debugf("\n\n(INFO) teleport node stop success!\n")
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
                switch state {
                    case locator.SlaveCryptoCheck:
                        fallthrough
                    case locator.SlaveBindBroken: {
                        app.BroadcastEvent(service.Event{Name:nodeTeleportStart, Payload:state})
                        return nil
                    }
                    case locator.SlaveBounded: {
                        return nil
                    }
                    default: {
                        return nil
                    }
                }

            } else {
                log.Debugf("(INFO) [%v] BeaconEventTranstion -> %v | FAILED ", ts, state.String())
                switch state {
                    case locator.SlaveBounded: {
                        app.BroadcastEvent(service.Event{Name:nodeTeleportStop, Payload:state})
                        return nil
                    }
                    default: {
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
            authToken, err := context.GetSlaveAuthToken()
            if err == nil && len(authToken) != 0 {
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
                    case <- app.StopChannel(): {
                        return nil
                    }
                    case <- timer.C: {
                        err = loc.TranstionWithTimestamp(time.Now())
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                    }
                    case evt := <-beaconC: {
                        mp, mk := evt.Payload.(ucast.BeaconPack)
                        if mk {
                            err = loc.TranstionWithMasterBeacon(mp, time.Now())
                            if err != nil {
                                log.Debug(err.Error())
                            }
                        }
                    }
                    case dvt := <- dhcpC: {
                        log.Debugf("[DHCP] RECEIVED\n %v", spew.Sdump(dvt.Payload))
                    }
                }
            }
            return nil
        }
    )

    app.WaitForEvent(nodeFeedbackBeacon, beaconC, make(chan struct{}))
    app.WaitForEvent(nodeFeedbackDHCP,     dhcpC, make(chan struct{}))
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

    // DNS service
    err = initDNSService(app)
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

