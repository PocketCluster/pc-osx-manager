package main

import "C"
import (
    "time"
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/gravitational/teleport/lib/process"
    "github.com/coreos/etcd/embed"
    "gopkg.in/tylerb/graceful.v1"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/lifecycle"
    "github.com/stkim1/pc-core/event/network"
    "github.com/stkim1/pc-core/event/crash"
    "github.com/stkim1/pc-core/event/operation"
    telesrv "github.com/stkim1/pc-core/extsrv/teleport"
    regisrv "github.com/stkim1/pc-core/extsrv/registry"
    swarmsrv "github.com/stkim1/pc-core/extsrv/swarm"
    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
)

const (
    coreFeedbackSearch = "feedback_search"
    coreFeedbackBeacon = "feedback_beacon"
    coreServiceBeacon  = "service_beacon"
)

func main() {

    mainLifeCycle(func(a *mainLife) {

        var (
            serviceConfig *serviceConfig = nil
            teleProc *process.PocketCoreProcess = nil
            regiProc *regisrv.PocketRegistry = nil
            swarmProc *swarmsrv.Server
            swarmSrv *graceful.Server
            err error = nil

            srvWaiter sync.WaitGroup
        )

        go func(wg *sync.WaitGroup) {
            wg.Wait()
        }(&srvWaiter)

        for e := range a.Events() {
            switch e := a.Filter(e).(type) {

                // APPLICATION LIFECYCLE //

                case lifecycle.Event: {
                    switch e.Crosses(lifecycle.StageDead) {
                        case lifecycle.CrossOn: {
                            log.Debugf("[LIFE] app is Dead %v", e.String())
                        }
                        case lifecycle.CrossOff: {
                            log.Debugf("[LIFE] app is not dead %v", e.String())
                        }
                    }
                    switch e.Crosses(lifecycle.StageAlive) {
                        case lifecycle.CrossOn: {
                            log.Debugf("[LIFE] app is now alive %v", e.String())
                            log.Debugf("[PREP] PREPARING GOLANG CONTEXT")
                            serviceConfig, err = setupServiceConfig()
                            if err != nil {
                                // TODO send error report
                            }
                            FeedSend("successfully initiated engine ..." + serviceConfig.teleConfig.HostUUID)
                        }
                        case lifecycle.CrossOff: {
                            log.Debugf("[LIFE] app is inactive %v", e.String())
                        }
                    }
                    switch e.Crosses(lifecycle.StageVisible) {
                        case lifecycle.CrossOn: {
                            log.Debugf("[LIFE] app is visible %v", e.String())
                        }
                        case lifecycle.CrossOff: {
                            log.Debugf("[LIFE] app is invisible %v", e.String())
                        }
                    }
                    switch e.Crosses(lifecycle.StageFocused) {
                        case lifecycle.CrossOn: {
                            log.Debugf("[LIFE] app is focused %v", e.String())
                        }
                        case lifecycle.CrossOff: {
                            log.Debugf("[LIFE] app is not focused %v", e.String())
                        }
                    }
                }

                // NETWORK EVENT //

                case network.Event: {
                    switch e.NetworkEvent {
                        case network.NetworkChangeInterface: {
                            //log.Debugf(spew.Sdump(e.HostInterfaces))
                            log.Debugf("[NET] %v", e.String())
                            context.MonitorNetworkInterfaces(e.HostInterfaces)
                        }
                        case network.NetworkChangeGateway: {
                            //log.Debugf(spew.Sdump(e.HostGateways))
                            log.Debugf("[NET] %v", e.String())
                            context.MonitorNetworkGateways(e.HostGateways)
                        }
                    }
                }

                // [DEBUG] ARTIFICIAL CRASH //

                case crash.Crash: {
                    switch e.Reason {
                    case crash.CrashEmergentExit: {
                        log.Printf("[CRASH] COCOA SIDE RUNTIME IS DESTORYED. WE NEED TO CLOSE GOLANG SIDE AS WELL. %v", e.String())
                    }
                    default:
                        log.Printf("crash! %v", e.String())
                    }
                }

                // OPERATIONAL COMMAND //

                case operation.Operation: {
                    switch e.Command {

                    /// BEACON ///

                    case operation.CmdBeaconStart: {
                        // TODO : use network interface
                        catcher, err := mcast.NewSearchCatcher("en0")
                        if err != nil {
                            log.Debug(errors.WithStack(err))
                        } else {
                            a.RegisterServiceFunc(func() error {
                                log.Debugf("NewSearchCatcher :: MAIN BEGIN")
                                for {
                                    select {
                                        case <-a.StopChannel(): {
                                            catcher.Close()
                                            log.Debugf("NewSearchCatcher :: MAIN CLOSE")
                                            return nil
                                        }
                                        case r := <-catcher.ChRead: {
                                            log.Debugf("SearchCatcher %v", r.Message)
                                            a.BroadcastEvent(Event{Name:coreFeedbackSearch, Payload:r})
                                        }
                                    }
                                }
                                return nil
                            })
                        }

                        belocat, err := ucast.NewBeaconLocator()
                        if err != nil {
                            log.Debug(errors.WithStack(err))
                        } else {
                            // beacon locator read
                            a.RegisterServiceFunc(func() error {
                                log.Debugf("NewBeaconLocator READ :: MAIN BEGIN")
                                for {
                                    select {
                                        case <- a.StopChannel(): {
                                            belocat.Close()
                                            log.Debugf("NewBeaconLocator READ :: MAIN CLOSE")
                                            return nil
                                        }
                                        case r := <- belocat.ChRead: {
                                            log.Debugf("BeaconLocator READ %v", r.Message)
                                            a.BroadcastEvent(Event{Name:coreFeedbackBeacon, Payload:r.Message})
                                        }
                                    }
                                }
                                return nil
                            })

                            // beacon locator write
                            beaconC := make(chan Event)
                            a.WaitForEvent(coreServiceBeacon, beaconC, make(chan struct{}))
                            a.RegisterServiceFunc(func() error {
                                log.Debugf("NewBeaconLocator WRITE :: MAIN BEGIN")
                                for {
                                    select {
                                        case <- a.StopChannel(): {
                                            log.Debugf("NewBeaconLocator WRITE :: MAIN CLOSE")
                                            return nil
                                        }
                                        case b := <- beaconC: {
                                            bs, ok := b.Payload.(ucast.BeaconSend)
                                            if ok {
                                                log.Debugf("NewBeaconLocator WRITE %v", bs.Host)
                                                belocat.Send(bs.Host, bs.Payload)
                                            }
                                        }
                                    }
                                }
                                return nil
                            })
                        }
                        a.StartServices()
                        log.Debugf("[OP] %v", e.String())
                    }

                    case operation.CmdBeaconStop: {
                        a.StopServices()
                        log.Debugf("[OP] %v", e.String())
                    }

                    /// TELEPORT ///

                    case operation.CmdTeleportStart: {
                        log.Debugf("[OP] %v", e.String())

                        teleProc, err = telesrv.NewTeleportCore(serviceConfig.teleConfig)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        err = teleProc.Start()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                    }
                    case operation.CmdTeleportStop: {
                        err = teleProc.Close()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        err = teleProc.Wait()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        log.Debugf("[OP] %v", e.String())
                    }

                    /// REGISTRY ///

                    case operation.CmdRegistryStart: {
                        log.Debugf("[OP] %v", e.String())
                        regiProc, err = regisrv.NewPocketRegistry(serviceConfig.regConfig)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        srvWaiter.Add(1)
                        err = regiProc.StartOnWaitGroup(&srvWaiter)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                    }
                    case operation.CmdRegistryStop: {
/*
                        err = regiProc.Close()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
*/
                        regiProc.Stop(time.Second)
                        log.Debugf("[OP] %v", e.String())
                    }

                    /// ORCHESTRATION ///

                    case operation.CmdCntrOrchStart: {
                        swarmProc, err = swarmsrv.NewSwarmServer(serviceConfig.swarmConfig)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        srvWaiter.Add(1)
                        swarmSrv, err = swarmProc.ListenAndServeOnWaitGroup(&srvWaiter)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdCntrOrchStop: {
                        log.Debugf("[OP] %v", e.String())
                        go func() {
                            srvWaiter.Wait()
                        }()
                        swarmSrv.Stop(time.Second)
                    }

                    /// STORAGE ///

                    case operation.CmdStorageStart: {
                        log.Debugf("[OP] %v", e.String())
/*
                        etcdProc, err = embed.StartPocketEtcd(serviceConfig.etcdConfig)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        etcdProc.Server.Start()
*/
                        srvWaiter.Add(1)
                        go func() {
                            defer srvWaiter.Done()
                            e, err := embed.StartPocketEtcd(serviceConfig.etcdConfig)
                            if err != nil {
                                log.Debugf(err.Error())
                                return
                            }
                            defer e.Close()
                            select {
                            case <-e.Server.ReadyNotify():
                                log.Printf("Server is ready!")
                            case <-time.After(60 * time.Second):
                                e.Server.Stop() // trigger a shutdown
                                log.Printf("Server took too long to start!")
                            }
                            log.Fatal(<-e.Err())
                        }()
                    }
                    case operation.CmdStorageStop: {
                        log.Debugf("[OP] %v", e.String())
//                        etcdProc.Server.Stop()
                    }


                    case operation.CmdServiceBundleStart: {
                        eventC := make(chan Event)
                        a.WaitForEvent("TEST_EVENT", eventC, make(chan struct{}))

                        a.RegisterServiceFunc(func() error {
                            defer log.Debugf("[TEST SERV 1] -- SERVICE 1 ENDED --")
                            log.Debugf("[TEST SERV 1] test for-select loop started...")

                            for {
                                select {
                                    case <- a.StopChannel():
                                        log.Debugf("[TEST SERV 1] [TEST 1 STOPPING]")
                                        return nil
                                    case <- eventC:
                                        log.Debugf("[TEST SERV 1] new Event received...")
                                }
                            }

                            return nil
                        })

                        a.RegisterServiceFunc(func() error {
                            defer log.Debugf("[TEST SERV 2] -- SERVICE 2 ENDED --")
                            log.Debugf("[TEST SERV 2] test started")

                            for {
                                if a.IsStopped() {
                                    log.Debugf("[TEST SERV 2] [TEST 2 STOPPING]")
                                    return nil
                                }

                                a.BroadcastEvent(Event{Name:"TEST_EVENT"})

                                time.Sleep(time.Second)
                            }

                            return nil
                        })
                        a.StartServices()
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdServiceBundleStop: {
                        a.StopServices()
                        log.Debugf("[OP] %v", e.String())
                    }

                    default:
                        log.Debug("[OP-ERROR] THIS SHOULD NOT HAPPEN %v", e.String())
                    }
                }
            }
        }
    })
}
