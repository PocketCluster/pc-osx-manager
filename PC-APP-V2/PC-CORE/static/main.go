package main

import "C"
import (
    log "github.com/Sirupsen/logrus"
    tembed "github.com/gravitational/teleport/embed"
    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/lifecycle"
    "github.com/stkim1/pc-core/event/network"
    "github.com/stkim1/pc-core/event/crash"
    "github.com/stkim1/pc-core/event/operation"
)

func main() {

    appLifeCycle(func(a *appMainLife) {

        var (
            config *serviceConfig = nil
            err error = nil
        )

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
                            config, err = setupServiceConfig()
                            if err != nil {
                                // TODO send error report
                                log.Debugf("[LIFE] CRITICAL ERROR %v", err)
                            } else {
                                FeedSend("[LIFE] successfully initiated engine ..." + config.teleConfig.HostUUID)
                            }
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

                    /// BASE OPERATION ///

                    case operation.CmdBaseServiceStart: {

                        cid, err := context.SharedHostContext().MasterAgentName()
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // name service
                        err = initPocketNameService(a, cid)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // storage service
                        err = initStorageServie(a, config.etcdConfig)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // registry
                        err = initRegistryService(a, config.regConfig)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // teleport service added
                        // TODO : need to hold teleport instance from GC
                        _, err = tembed.NewEmbeddedCoreProcess(a.ServiceSupervisor, config.teleConfig)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // beacon service added
                        // TODO : use network interface
                        // TODO : need to catcher instance from GC
                        _, err = mcast.NewSearchCatcher(a.ServiceSupervisor, "en0")
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // TODO : need to hold beacon instance from GC
                        _, err = ucast.NewBeaconLocator(a.ServiceSupervisor)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        err = initSwarmService(a)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        err = initMasterBeaconService(a, cid, config.teleConfig)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        a.StartServices()
                        log.Debugf("[OP] %v", e.String())
                    }

                    case operation.CmdBaseServiceStop: {
                        a.StopServices()
                        log.Debugf("[OP] %v", e.String())
                    }

                    /// STORAGE ///

                    case operation.CmdStorageStart: {
                        err = initStorageServie(a, config.etcdConfig)
                        if err != nil {
                            log.Debug(err)
                            return
                        }
                        a.StartServices()
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdStorageStop: {
                        a.StopServices()
                        log.Debugf("[OP] %v", e.String())
                    }

                    case operation.CmdTeleportRootAdd: {
                        cid, err := context.SharedHostContext().MasterAgentName()
                        if err != nil {
                            log.Debug(err)
                        }
                        err = buildVboxCoreDisk(cid, config.teleConfig)
                        if err != nil {
                            log.Debug(err)
                        }
                    }
                    case operation.CmdTeleportUserAdd: {
                        cid, err := context.SharedHostContext().MasterAgentName()
                        if err != nil {
                            log.Debug(err)
                        }

                        err = initVboxCoreReportService(a, cid)
                        if err != nil {
                            log.Debug(err)
                        }
                    }

                    /// DEBUG ///

                    case operation.CmdDebug: {
                        sl := a.ServiceList()
                        for i, _ := range sl {
                            s := sl[i]
                            log.Debugf("[SERVICE] %s, %v", s.Tag(), s.IsRunning())
                        }
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
