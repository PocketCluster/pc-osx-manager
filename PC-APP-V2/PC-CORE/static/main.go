package main

import "C"
import (
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/lifecycle"
    "github.com/stkim1/pc-core/event/network"
    "github.com/stkim1/pc-core/event/crash"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/extlib/pcssh/sshproc"
    "github.com/stkim1/pc-core/service"
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
                            log.Debugf("[NET] %v", e.String())

                            // TODO check if service is running
                            isSrvRun := len(a.ServiceList()) != 0
                            updated := context.SharedHostContext().UpdateNetworkInterfaces(e.HostInterfaces)

                            // services should be running before receiving event. Otherwise, service will not start
                            if isSrvRun && updated {
                                log.Debugf("[NET] network address change event triggered")
                                a.BroadcastEvent(service.Event{Name:iventNetworkAddressChange})
                            }
                        }
                        case network.NetworkChangeGateway: {
                            log.Debugf("[NET] %v", e.String())
                            context.SharedHostContext().UpdateNetworkGateways(e.HostGateways)
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

                        // TODO check all the status before start

                        // config should be setup after acquiring ip address on wifi
                        // This should run only once
                        config, err = setupServiceConfig()
                        if err != nil {
                            // TODO send error report
                            log.Debugf("[LIFE] CRITICAL ERROR %v", err)
                            return
                        } else {
                            FeedSend("[LIFE] SUCCESSFULLY INITIATED ENGINE " + config.teleConfig.HostUUID)
                        }

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
                        _, err = sshproc.NewEmbeddedMasterProcess(a.ServiceSupervisor, config.teleConfig)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // beacon service added
                        iname, err := context.SharedHostContext().HostPrimaryInterfaceShortName()
                        if err != nil {
                            log.Debug(err)
                            return
                        }
                        _, err = mcast.NewSearchCatcher(a.ServiceSupervisor, iname)
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

                        err = initVboxCoreReportService(a, cid)
                        if err != nil {
                            log.Debug(err)
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
                    }
                    case operation.CmdTeleportUserAdd: {
                        log.Debugf("[OP] %v", e.String())
                    }

                    /// DEBUG ///

                    case operation.CmdDebug: {
/*
                        sl := a.ServiceList()
                        for i, _ := range sl {
                            s := sl[i]
                            log.Debugf("[SERVICE] %s, %v", s.Tag(), s.IsRunning())
                        }
*/
                        cid, err := context.SharedHostContext().MasterAgentName()
                        if err != nil {
                            log.Debug(err)
                        }
                        err = buildVboxCoreDisk(cid, config.teleConfig)
                        if err != nil {
                            log.Debug(err)
                        }
                        err = buildVboxMachine(a)
                        if err != nil {
                            log.Debugf("vbox operation error %v", err)
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
