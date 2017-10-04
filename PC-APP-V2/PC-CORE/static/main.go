// +build darwin
package main

import "C"
import (
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"

    "github.com/stkim1/pc-core/config"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/event/lifecycle"
    "github.com/stkim1/pc-core/event/network"
    "github.com/stkim1/pc-core/event/crash"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/extlib/pcssh/sshproc"
    "github.com/stkim1/pc-core/extlib/pcssh/sshadmin"
    "github.com/stkim1/pc-core/model"
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
    "github.com/stkim1/pc-core/vboxglue"
)

func main() {

    appLifeCycle(func(appLife *appMainLife) {

        var (
            appCfg        *config.ServiceConfig = nil
            err           error                 = nil
            vboxCore      vboxglue.VBoxGlue     = nil
            IsContextInit bool                  = false
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
                            if !IsContextInit {
                                log.Debugf("[LIFE] initialize ApplicationContext...")
                                // this needs to be initialized before service loop initiated
                                context.SharedHostContext()
                                initcheck.InitRoutePathServices(appLife, theFeeder)
                                install.InitInstallListRouthPath(appLife, theFeeder)

                                IsContextInit = true
                            }

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
                        appCfg, err = config.InitServiceConfig()
                        if err != nil {
                            // TODO send error report
                            log.Debugf("[LIFE] CRITICAL ERROR %v", err)
                            continue
                        }
                        // TODO : move this initializer after network initiated
                        install.InitInstallPackageRoutePath(appLife, theFeeder, appCfg.PCSSH)

                        // TODO check all the status before start

                        // --- acquire informations ---
                        // get cluster id
                        cid, err := context.SharedHostContext().MasterAgentName()
                        if err != nil {
                            log.Debug(err)
                            continue
                        }
                        // get primary interface bsd name
                        iname, err := context.SharedHostContext().HostPrimaryInterfaceShortName()
                        if err != nil {
                            log.Debug(err)
                            continue
                        }

                        // --- role service sequence ---
                        // storage service
                        // (NODEP netchange, NODEP services)
                        err = container.InitStorageServie(appLife, appCfg.ETCD)
                        if err != nil {
                            log.Debug(err)
                            continue
                        }

                        // registry service
                        // (NODEP netchange, NODEP services)
                        err = container.InitRegistryService(appLife, appCfg.REG)
                        if err != nil {
                            log.Debug(err)
                            continue
                        }

                        // search catcher service
                        // (DEP netchange, NODEP services)
                        // TODO : need to hold beacon instance from GC -> not necessary as it embeds service instance???
                        _, err = mcast.NewSearchCatcher(appLife.ServiceSupervisor, iname)
                        if err != nil {
                            log.Debug(err)
                            continue
                        }
                        // beacon locator service
                        // (NODEP netchange, NODEP service)
                        // TODO : need to hold beacon instance from GC -> not necessary as it embeds service instance???
                        _, err = ucast.NewBeaconLocator(appLife.ServiceSupervisor)
                        if err != nil {
                            log.Debug(err)
                            continue
                        }

                        // internal name service
                        // (NODEP netchange, DEP master beacon service)
                        err = dns.InitPocketNameService(appLife, cid)
                        if err != nil {
                            log.Debug(err)
                            continue
                        }

                        // swarm service
                        // (NODEP netchange, DEP master beacon service)
                        err = container.InitSwarmService(appLife)
                        if err != nil {
                            log.Debug(err)
                            continue
                        }

                        // master beacon service
                        // (DEP netchange, DEP vboxcontrol + teleport service)
                        err = master.InitMasterBeaconService(appLife, cid, appCfg.PCSSH)
                        if err != nil {
                            log.Debug(err)
                            continue
                        }

                        // vboxcontrol service
                        // (DEP netchange, NODEP service)
                        err = vbox.InitVboxCoreReportService(appLife, cid)
                        if err != nil {
                            log.Debug(err)
                            continue
                        }

                        // teleport service
                        // (DEP netchange, NODEP services)
                        // TODO : need to hold teleport instance from GC -> not necessary as it embeds service instance???
                        _, err = sshproc.NewEmbeddedMasterProcess(appLife.ServiceSupervisor, appCfg.PCSSH)
                        if err != nil {
                            log.Debug(err)
                            continue
                        }

                        err = health.InitSystemHealthMonitor(appLife, theFeeder)
                        if err != nil {
                            log.Debug(err)
                            continue
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
                            continue
                        }
                        appLife.StartServices()
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdStorageStop: {
                        appLife.StopServices()
                        log.Debugf("[OP] %v", e.String())
                    }

                    /// DEBUG ///

                    case operation.CmdDebug0: {
                        // setup users
                        {
                            cli, err := sshadmin.OpenAdminClientWithAuthService(appCfg.PCSSH)
                            if err != nil {
                                log.Error(err.Error())
                            }
                            roots, err := model.FindUserMetaWithLogin("root")
                            if err != nil {
                                log.Error(err.Error())
                            }
                            err = sshadmin.CreateTeleportUser(cli, "root", roots[0].Password)
                            if err != nil {
                                log.Error(err.Error())
                            }
                            uname, err := context.SharedHostContext().LoginUserName()
                            if err != nil {
                                log.Error(err.Error())
                            }
                            lusers, err := model.FindUserMetaWithLogin(uname)
                            if err != nil {
                                log.Error(err.Error())
                            }
                            err = sshadmin.CreateTeleportUser(cli, uname, lusers[0].Password)
                            if err != nil {
                                log.Error(err.Error())
                            }
                        }
                        log.Debugf("[OP] %v", e.String())
                    }

                    case operation.CmdDebug1: {
                        // setup vbox
                        {
                            cid, err := context.SharedHostContext().MasterAgentName()
                            if err != nil {
                                log.Debug(err.Error())
                                continue
                            }

                            err = vboxglue.BuildVboxCoreDisk(cid, appCfg.PCSSH)
                            if err != nil {
                                log.Debug(err.Error())
                                continue
                            }

                            vcore, err := vboxglue.NewGOVboxGlue()
                            if err != nil {
                                log.Debug(err.Error())
                                continue
                            }

                            err = vboxglue.CreateNewMachine(vcore)
                            if err != nil {
                                log.Debug(err.Error())
                                vcore.Close()
                                continue
                            }

                            // shutoff vbox core. very unlikely
                            if !vcore.IsMachineSafeToStart() {
                                err := vboxglue.EmergencyStop(vcore, defaults.PocketClusterCoreName)
                                if err != nil {
                                    log.Debug(err.Error())
                                    vcore.Close()
                                    continue
                                }
                            }

                            // then start back up
                            err = vcore.StartMachine()
                            if err != nil {
                                log.Debug(err.Error())
                                vboxCore.Close()
                                continue
                            }
                            vboxCore = vcore
                        }
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdDebug2: {
                        // start machine
                        {
                            vcore, err := vboxglue.NewGOVboxGlue()
                            if err != nil {
                                log.Debug(err.Error())
                                continue
                            }
                            err = vcore.FindMachineByNameOrID(defaults.PocketClusterCoreName)
                            if err != nil {
                                log.Debug(err.Error())
                                vboxCore.Close()
                                continue
                            }

                            // force shutoff vbox core
                            if !vcore.IsMachineSafeToStart() {
                                err := vboxglue.EmergencyStop(vcore, defaults.PocketClusterCoreName)
                                if err != nil {
                                    log.Debug(err.Error())
                                    vboxCore.Close()
                                    continue
                                }
                            }

                            // check if machine setting changed
                            chgd, err := vcore.IsMachineSettingChanged()
                            if err != nil {
                                log.Debug(err.Error())
                                vcore.Close()
                                continue
                            }
                            // warn user and abort boot procedure
                            if chgd {
                                log.Errorf("core node setting has changed. abort boot procedure")
                                // reset the option again (2017/10/03 : this is not working as of now)
                                _ = vboxglue.ResetExistingMachine(vcore)
                                vcore.Close()
                                continue
                            }

                            // then start back up
                            err = vcore.StartMachine()
                            if err != nil {
                                log.Debug(err.Error())
                                vboxCore.Close()
                                continue
                            }
                            vboxCore = vcore
                        }
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdDebug3: {
                        // stop machine
                        {
                            // this is case where previous run or user has acticated pc-core
                            if vboxCore == nil {
                                vcore, err := vboxglue.NewGOVboxGlue()
                                if err != nil {
                                    log.Debug(err.Error())
                                    continue
                                }
                                err = vcore.FindMachineByNameOrID(defaults.PocketClusterCoreName)
                                if err != nil {
                                    log.Debug(err.Error())
                                    vcore.Close()
                                    continue
                                }

                                if !vcore.IsMachineSafeToStart() {
                                    err := vboxglue.EmergencyStop(vcore, defaults.PocketClusterCoreName)
                                    if err != nil {
                                        log.Debug(err.Error())
                                    }
                                }
                                err = vcore.Close()
                                if err != nil {
                                    log.Debug(err.Error())
                                }

                            } else {
                                // normal start and stop procedure
                                err := vboxCore.AcpiStopMachine()
                                if err != nil {
                                    log.Debug(err.Error())
                                }
                                err = vboxCore.Close()
                                if err != nil {
                                    log.Debug(err.Error())
                                }
                                vboxCore = nil
                            }
                        }

                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdDebug4: {
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdDebug5: {
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdDebug6: {
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdDebug7: {
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
