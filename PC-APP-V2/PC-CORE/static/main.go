package main

import "C"
import (
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"

    "github.com/stkim1/pc-core/config"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/lifecycle"
    "github.com/stkim1/pc-core/event/network"
    "github.com/stkim1/pc-core/event/crash"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/extlib/pcssh/sshproc"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/install"
    "github.com/stkim1/pc-core/route/initcheck"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/container"
    "github.com/stkim1/pc-core/service/dns"
    "github.com/stkim1/pc-core/service/health"
    "github.com/stkim1/pc-core/service/ivent"
    "github.com/stkim1/pc-core/service/master"
    "github.com/stkim1/pc-core/service/vbox"
)

func main() {

    appLifeCycle(func(appLife *appMainLife) {

        var (
            appCfg *config.ServiceConfig = nil
            err    error                 = nil
        )

        for e := range appLife.Events() {
            switch e := appLife.Filter(e).(type) {

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
                            // this should happen only once in the lifetime of an application.
                            // so we'll initialize our context here to have safe operation

                            log.Debugf("[LIFE] initialize ApplicationContext...")
                            // this needs to be initialized before service loop initiated
                            context.SharedHostContext()
                            initcheck.InitRoutePathServices(appLife, theFeeder)
                            install.InitInstallRoutePath(appLife, theFeeder)

                            log.Debugf("[LIFE] app is now created, fully initialized %v", e.String())
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
                            isSrvRun := len(appLife.ServiceList()) != 0
                            updated := context.SharedHostContext().UpdateNetworkInterfaces(e.HostInterfaces)

                            // services should be running before receiving event. Otherwise, service will not start
                            if isSrvRun && updated {
                                log.Debugf("[NET] network address change event triggered")
                                appLife.BroadcastEvent(service.Event{Name:ivent.IventNetworkAddressChange})
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

                // OPERATIONAL ROUTE //

                case route.Request: {
                    err := appLife.Dispatch(e)
                    if err != nil {
                        log.Debugf("[ROUTE] ERROR %v", err)
                    }
                }

                // OPERATIONAL COMMAND //

                case operation.Operation: {
                    switch e.Command {

                    /// BASE OPERATION ///

                    case operation.CmdBaseServiceStart: {

                        // config should be setup after acquiring ip address on wifi
                        // This should run only once
                        appCfg, err = config.SetupServiceConfig()
                        if err != nil {
                            // TODO send error report
                            log.Debugf("[LIFE] CRITICAL ERROR %v", err)
                            return
                        }

                        // TODO check all the status before start

                        // --- acquire informations ---
                        // get cluster id
                        cid, err := context.SharedHostContext().MasterAgentName()
                        if err != nil {
                            log.Debug(err)
                            return
                        }
                        // get primary interface bsd name
                        iname, err := context.SharedHostContext().HostPrimaryInterfaceShortName()
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // --- role service sequence ---
                        // storage service
                        // (NODEP netchange, NODEP services)
                        err = container.InitStorageServie(appLife, appCfg.ETCD)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // registry service
                        // (NODEP netchange, NODEP services)
                        err = container.InitRegistryService(appLife, appCfg.REG)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // search catcher service
                        // (DEP netchange, NODEP services)
                        // TODO : need to hold beacon instance from GC -> not necessary as it embeds service instance???
                        _, err = mcast.NewSearchCatcher(appLife.ServiceSupervisor, iname)
                        if err != nil {
                            log.Debug(err)
                            return
                        }
                        // beacon locator service
                        // (NODEP netchange, NODEP service)
                        // TODO : need to hold beacon instance from GC -> not necessary as it embeds service instance???
                        _, err = ucast.NewBeaconLocator(appLife.ServiceSupervisor)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // internal name service
                        // (NODEP netchange, DEP master beacon service)
                        err = dns.InitPocketNameService(appLife, cid)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // swarm service
                        // (NODEP netchange, DEP master beacon service)
                        err = container.InitSwarmService(appLife)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // master beacon service
                        // (DEP netchange, DEP vboxcontrol + teleport service)
                        err = master.InitMasterBeaconService(appLife, cid, appCfg.PCSSH)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // vboxcontrol service
                        // (DEP netchange, NODEP service)
                        err = vbox.InitVboxCoreReportService(appLife, cid)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        // teleport service
                        // (DEP netchange, NODEP services)
                        // TODO : need to hold teleport instance from GC -> not necessary as it embeds service instance???
                        _, err = sshproc.NewEmbeddedMasterProcess(appLife.ServiceSupervisor, appCfg.PCSSH)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        err = health.InitSystemHealthMonitor(appLife, theFeeder)
                        if err != nil {
                            log.Debug(err)
                            return
                        }

                        appLife.StartServices()
                        log.Debugf("[OP] %v", e.String())
                    }

                    case operation.CmdBaseServiceStop: {
                        appLife.StopServices()
                        log.Debugf("[OP] %v", e.String())
                    }

                    /// STORAGE ///

                    case operation.CmdStorageStart: {
                        err = container.InitStorageServie(appLife, appCfg.ETCD)
                        if err != nil {
                            log.Debug(err)
                            return
                        }
                        appLife.StartServices()
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdStorageStop: {
                        appLife.StopServices()
                        log.Debugf("[OP] %v", e.String())
                    }

                    case operation.CmdTeleportRootAdd: {
                    }
                    case operation.CmdTeleportUserAdd: {
                        log.Debugf("[OP] %v", e.String())
                    }

                    /// DEBUG ///

                    case operation.CmdDebug: {
                        cid, err := context.SharedHostContext().MasterAgentName()
                        if err != nil {
                            log.Debug(err)
                        }
                        err = vbox.BuildVboxCoreDisk(cid, appCfg.PCSSH)
                        if err != nil {
                            log.Debug(err)
                        }
                        err = vbox.BuildVboxMachine()
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
