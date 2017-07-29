package main

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/service"
    "github.com/stkim1/pc-node-agent/pcssh/sshproc"
)

import (
    "github.com/davecgh/go-spew/spew"
)

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
                        app.BroadcastEvent(service.Event{Name:sshproc.EventNodeSSHServiceStop})
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
        service.BindEventWithService(iventNodeDHCPFeedback, dhcpC),
    )

    return nil
}