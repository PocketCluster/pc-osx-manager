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
    "github.com/stkim1/pc-node-agent/utils/dhcp"
    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/service"
)

import (
    "github.com/davecgh/go-spew/spew"
)

const (
    nodeFeedbackDHCP   = "feedback_dhcp"
    dockerServiceUnit  = "docker.service"

    servicePcsshInit   = "service.pcssh.init"
    servicePcsshStart  = "service.pcssh.start"
)

func initDhcpListner(app service.AppSupervisor) error {
    // firstly clear off previous socket
    os.Remove(dhcp.DHCPEventSocketPath)
    dhcpListener, err := net.ListenUnix("unix", &net.UnixAddr{dhcp.DHCPEventSocketPath, "unix"})
    if err != nil {
        return errors.WithStack(err)
    }

    app.RegisterServiceWithFuncs(
        func () error {
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
        },
        func(_ func(interface{})) error {
            dhcpListener.Close()
            os.Remove(dhcp.DHCPEventSocketPath)
            log.Debugf("[DHCP] close dhcp listner...")
            return nil
        },
    )

    return nil
}

func initTeleportNodeService(app service.AppSupervisor) error {
    app.RegisterNamedServiceWithFuncs(
        servicePcsshInit,
        func() error{
            var (
                pcsshNode *embed.EmbeddedNodeProcess = nil
                err error = nil
            )
            // restart teleport
            cfg, err := tervice.MakeNodeConfig(slcontext.SharedSlaveContext(), true)
            if err != nil {
                return errors.WithStack(err)
            }
            pcsshNode, err = embed.NewEmbeddedNodeProcess(app, cfg)
            if err != nil {
                log.Errorf(err.Error())
                return errors.WithStack(err)
            }

            // execute docker engine cert acquisition before SSH node start
            // TODO : create a waitforevent channel and restart docker engine accordingly
            err = pcsshNode.AcquireEngineCertificate(slcontext.DockerEnvironemtPostProcess)
            if err != nil {
                return errors.WithStack(err)
            }

            err = pcsshNode.StartNodeSSH()
            if err != nil {
                return errors.WithStack(err)
            }
            log.Debugf("\n\n(INFO) teleport node started success!\n")

            return nil

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

            return nil
        },
        func(_ func(interface{})) error {
            return nil
        })

    app.RegisterNamedServiceWithFuncs(
        servicePcsshStart,
        func() error{
            var (
                pcsshNode *embed.EmbeddedNodeProcess = nil
                err error = nil
            )
            // restart teleport
            cfg, err := tervice.MakeNodeConfig(slcontext.SharedSlaveContext(), true)
            if err != nil {
                return errors.WithStack(err)
            }
            pcsshNode, err = embed.NewEmbeddedNodeProcess(app, cfg)
            if err != nil {
                return errors.WithStack(err)
            }

            err = pcsshNode.StartNodeSSH()
            if err != nil {
                return errors.WithStack(err)
            }
            log.Debugf("\n\n(INFO) teleport node started success!\n")

            return nil

            // restart docker engine
            // TODO : FIX /opt/gopkg/src/github.com/godbus/dbus/conn.go:345 send on closed channel
            conn, err := sysd.NewSystemdConnection()
            if err != nil {
                log.Errorf(err.Error())
            } else {
                did, err := conn.RestartUnit(dockerServiceUnit, "replace", nil)
                if err != nil {
                    return errors.WithStack(err)
                } else {
                    conn.Close()
                    log.Debugf("\n\n(INFO) docker engin restart success! ID %d\n", did)
                }
            }

            return nil
        },
        func(_ func(interface{})) error {
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
            app.BroadcastEvent(
                service.Event{
                    Name:       mcast.EventBeaconNodeSearchSend,
                    Payload:    mcast.CastPack{
                        Message:    data,
                    },
                })
            return nil
        }

        beaconTx = func(target string, data []byte) error {
            log.Debugf("[BEACON-TX] %v TO : %v", time.Now(), target)
            app.BroadcastEvent(
                service.Event{
                    Name:       ucast.EventBeaconNodeLocationSend,
                    Payload:    ucast.BeaconSend{
                        Host:       target,
                        Payload:    data,
                    },
                })
            return nil
        }

        transitEvent = func (state locator.SlaveLocatingState, ts time.Time, transOk bool) error {
            if transOk {
                log.Debugf("(INFO) [%v] BeaconEventTranstion -> %v | SUCCESS ", ts, state.String())
                switch state {
                    case locator.SlaveCryptoCheck: {
                        app.RunNamedService(servicePcsshInit)
                        return nil
                    }
                    case locator.SlaveBindBroken: {
                        app.RunNamedService(servicePcsshStart)
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
                        app.BroadcastEvent(service.Event{Name:embed.EventNodeSSHServiceStop})
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
        exitFunc = func(_ func(interface{})) error {
            log.Debugf("[AGENT] close agent service...")
            return nil
        }
    )

    app.RegisterServiceWithFuncs(
        serviceFunc,
        exitFunc,
        service.BindEventWithService(ucast.EventBeaconNodeLocationReceive, beaconC),
        service.BindEventWithService(nodeFeedbackDHCP, dhcpC),
    )

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
    _, err = mcast.NewSearchCaster(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // beacon service
    _, err = ucast.NewBeaconAgent(app)
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

