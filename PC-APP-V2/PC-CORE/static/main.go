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
    "github.com/stkim1/pc-core/rasker"
    "github.com/stkim1/pc-core/rasker/install"
    "github.com/stkim1/pc-core/rasker/pkgtask"
    "github.com/stkim1/pc-core/rasker/regnode"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/initcheck"
    "github.com/stkim1/pc-core/route/list"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/container"
    "github.com/stkim1/pc-core/service/dns"
    "github.com/stkim1/pc-core/service/health"
    "github.com/stkim1/pc-core/service/ivent"
    "github.com/stkim1/pc-core/service/master"
    "github.com/stkim1/pc-core/service/vbox"
    "github.com/stkim1/pc-core/vboxglue"
)

// for upgrade
import (
    ctx "context"
    tclient "github.com/gravitational/teleport/lib/client"
    tdefs "github.com/gravitational/teleport/lib/defaults"
)

func main() {

    appLifeCycle(func(appLife *appMainLife) {

        var (
            appCfg        *config.ServiceConfig = nil
            err           error                 = nil
            vboxCore      vboxglue.VBoxGlue     = nil
            IsContextInit, IsNetworkInit        = false, false
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
                            if !IsContextInit {
                                IsContextInit = true
                                // this should happen only once in the lifetime of an application.
                                // so we'll initialize our context here to have safe operation

                                // this needs to be initialized before service loop initiated
                                context.SharedHostContext()
                                log.Debugf("[LIFE] context creation")

                                // -- initial service path registration ---
                                // by this time the response feeder should be initialized on frontend
                                initcheck.InitApplicationCheck(appLife, theFeeder)
                                pkgtask.InitPackageLifeCycle(rasker.RouteTasker{
                                    ServiceSupervisor: appLife.ServiceSupervisor,
                                    Router: appLife.Router},
                                    theFeeder)
                                pkgtask.InitPackageKillCycle(rasker.RouteTasker{
                                    ServiceSupervisor: appLife.ServiceSupervisor,
                                    Router: appLife.Router},
                                    theFeeder)

                                // package list
                                list.InitRouthPathListAvailable(appLife, theFeeder)
                                list.InitRouthPathListInstalled(appLife, theFeeder)

                                // node registration
                                regnode.InitNodeRegisterCycle(rasker.RouteTasker{
                                    ServiceSupervisor: appLife.ServiceSupervisor,
                                    Router: appLife.Router},
                                    theFeeder)
                                regnode.InitNodeRegisterStop(rasker.RouteTasker{
                                    ServiceSupervisor: appLife.ServiceSupervisor,
                                    Router: appLife.Router},
                                    theFeeder)
                                regnode.InitNodeRegisterCanidate(rasker.RouteTasker{
                                    ServiceSupervisor: appLife.ServiceSupervisor,
                                    Router: appLife.Router},
                                    theFeeder)

                                log.Debugf("[LIFE] service path registration")

                                // (2017/11/10) we ought to have frontend check engine response, but then it complicated network monitoring.
                                // So, we'll just give more time for context to be initialized for now.
                                err := reportContextInit(appLife, theFeeder)
                                if err != nil {
                                    log.Errorf("[SYSCONTEXT] error in reporting network init %v", err.Error())
                                }

                                // make sure initialization happens only once
                                log.Debugf("[LIFE] app is now created, fully initialized %v", e.String())
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
                            updated := context.SharedHostContext().UpdateNetworkInterfaces(e.HostInterfaces)

                            // notify frontend to initiate the next move.
                            if !IsNetworkInit {
                                IsNetworkInit = true

                                err := reportNetworkInit(appLife, theFeeder)
                                if err != nil {
                                    log.Errorf("[SYSNET] error in reporting network init %v", err.Error())
                                }
                            }

                            // services should be running before receiving event. Otherwise, service will not start
                            // TODO check if service is running
                            isSrvRun := len(appLife.ServiceList()) != 0
                            if isSrvRun && updated {
                                log.Debugf("[SYSNET] network address change event triggered")
                                appLife.BroadcastEvent(service.Event{Name:ivent.IventNetworkAddressChange})
                            }
                            log.Debugf("[SYSNET] %v", e.String())
                        }
                        case network.NetworkChangeGateway: {
                            context.SharedHostContext().UpdateNetworkGateways(e.HostGateways)
                            log.Debugf("[NET] %v", e.String())
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

                        // check if this is the first run
                        var isFirstTimeRun = context.SharedHostContext().CheckIsFistTimeExecution()

                        /*
                         * health monitor monitors every error and internal service spawn external listeners
                         * Any error happens in initializing internal service is critical one.
                         *
                         * (health monitor doesn't return any error)
                         * (NODEP netchange, DEP beacon, orchst, pcssh, vbox, discovery, registry)
                         */
                        health.InitSystemHealthMonitor(appLife, theFeeder)

                        // config should be setup after acquiring ip address on wifi
                        // This should run only once
                        appCfg, err = config.InitServiceConfig()
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        // --- acquire critical information ---
                        // TODO check all the status before start
                        // get primary interface bsd name
                        iname, err := context.SharedHostContext().HostPrimaryInterfaceShortName()
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        // get cluster id
                        cid, err := context.SharedHostContext().MasterAgentName()
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        // --- services ---
                        // storage service
                        // (NODEP netchange, NODEP services)
                        err = container.InitDiscoveryService(appLife, appCfg.ETCD)
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        // registry service
                        // (NODEP netchange, NODEP services)
                        err = container.InitRegistryService(appLife, appCfg.REG)
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        // internal name service
                        // (NODEP netchange, DEP master beacon service)
                        err = dns.InitPocketNameService(appLife, cid)
                        if err != nil {
                            log.Debug(err)
                            continue
                        }

                        // orcst service
                        // (NODEP netchange, DEP master beacon service)
                        err = container.InitOrchstService(appLife)
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        // master beacon service
                        // (DEP netchange, DEP vboxcontrol + teleport service)
                        err = master.InitMasterBeaconService(appLife, cid, appCfg.PCSSH)
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        if isFirstTimeRun {
                            // base user and vbox core setup
                            // (DEP teleport service)
                            err = setupBaseUsersWithVboxCore(appLife, appCfg.PCSSH)
                            if err != nil {
                                log.Error(err)
                                appLife.BroadcastEvent(service.Event{
                                    Name:ivent.IventInternalSpawnError,
                                    Payload:err})
                                continue
                            }
                        }

                        // teleport service
                        // (DEP netchange, NODEP services)
                        // TODO : need to hold teleport instance from GC -> not necessary as it embeds service instance???
                        _, err = sshproc.NewEmbeddedMasterProcess(appLife.ServiceSupervisor, appCfg.PCSSH)
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        // --- internal listeners ---
                        // vboxcontrol service comes first (as it's internal network listener)
                        // (DEP netchange, NODEP service)
                        err = vbox.InitVboxCoreReportService(appLife, cid)
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        // up until this point everything has to be executed fast as other services are waiting with timeout.
                        // after this, we can take time to start up, and open external listener once everyone is all ready.

                        // --- additional route path event ---
                        install.InitRoutePathInstallPackage(
                            rasker.RouteTasker{
                                ServiceSupervisor: appLife.ServiceSupervisor,
                                Router: appLife.Router},
                            theFeeder, appCfg.PCSSH)

                        // --- external listeners ---
                        // beacon locator service needs to initiated after master beacon
                        // (NODEP netchange, NODEP service)
                        // TODO : need to hold beacon instance from GC -> not necessary as it embeds service instance???
                        _, err = ucast.NewBeaconLocator(appLife.ServiceSupervisor)
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        // search catcher service needs to initiated after master beacon
                        // (DEP netchange, NODEP services)
                        // TODO : need to hold beacon instance from GC -> not necessary as it embeds service instance???
                        _, err = mcast.NewSearchCatcher(appLife.ServiceSupervisor, iname)
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }

                        // start services
                        appLife.StartServices()

                        // if this is the first time run, then wait for setup completed and start vbox core
                        if isFirstTimeRun {
                            reportC := make(chan service.Event)
                            if err := appLife.BindDiscreteEvent(ivent.IventSetupUsersAndVboxCore, reportC); err != nil {
                                log.Error(err)
                                appLife.BroadcastEvent(service.Event{
                                    Name:    ivent.IventInternalSpawnError,
                                    Payload: err})
                                continue
                            }
                            sdone := <-reportC
                            appLife.UntieDiscreteEvent(ivent.IventSetupUsersAndVboxCore)
                            if err, ok := sdone.Payload.(error); !ok || err != nil {
                                log.Error("unable launch vbox core due to setup error")
                                appLife.BroadcastEvent(service.Event{
                                    Name:    ivent.IventInternalSpawnError,
                                    Payload: err})
                                continue
                            }
                        }
                        vboxCore, err = startVboxCore()
                        if err != nil {
                            log.Error(err)
                            appLife.BroadcastEvent(service.Event{
                                Name:ivent.IventInternalSpawnError,
                                Payload:err})
                            continue
                        }
                        log.Debugf("[OP] %v", e.String())
                    }

                    case operation.CmdBaseServiceStop: {
                        err := stopHealthMonitor(appLife)
                        if err != nil {
                            log.Errorf("[ERROR] unable to terminate app %v", err.Error())
                        }
                        // as we close service and shutdown vbox core, we'll have chance to clean up docker-compose task
                        err = appLife.StopServices()
                        if err != nil {
                            log.Errorf("[ERROR] unable to close core node %v", err.Error())
                        }
                        err = stopVboxCore(vboxCore)
                        if err != nil {
                            log.Errorf("[ERROR] unable to close core node %v", err.Error())
                        }
                        err = model.CloseRecordGate()
                        if err != nil {
                            log.Errorf("[ERROR] error in closing storage %v", err.Error())
                        }
                        err = feedShutdownReadySignal(theFeeder)
                        if err != nil {
                            log.Errorf("[ERROR] error in closing storage %v", err.Error())
                        }
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdClusterShutdown: {
                        err := shutdownSlaveNodes(appLife, appCfg.PCSSH)
                        if err != nil {
                            log.Errorf("[ERROR] unable to shutdown cluster %v", err.Error())
                        }
                        err = stopHealthMonitor(appLife)
                        if err != nil {
                            log.Errorf("[ERROR] unable to terminate app %v", err.Error())
                        }
                        // as we close service and shutdown vbox core, we'll have chance to clean up docker-compose task
                        err = appLife.StopServices()
                        if err != nil {
                            log.Errorf("[ERROR] unable to close core node %v", err.Error())
                        }
                        err = stopVboxCore(vboxCore)
                        if err != nil {
                            log.Errorf("[ERROR] unable to close core node %v", err.Error())
                        }
                        err = model.CloseRecordGate()
                        if err != nil {
                            log.Errorf("[ERROR] error in closing storage %v", err.Error())
                        }
                        err = feedShutdownReadySignal(theFeeder)
                        if err != nil {
                            log.Errorf("[ERROR] error in closing storage %v", err.Error())
                        }
                        log.Debugf("[OP] %v", e.String())
                    }

                    /// STORAGE ///

                    case operation.CmdStorageStart: {
                        err = container.InitDiscoveryService(appLife, appCfg.ETCD)
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
                            // warn user and reset additional changes
                            if chgd {
                                log.Errorf("core node setting has changed. discard additional settings")
                                err = vcore.DiscardMachineSettings()
                                if err != nil {
                                    // unable to discard changes. abort startup
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
                        // update process
                        roots, err := model.FindUserMetaWithLogin("root")
                        if err != nil {
                            log.Error(err.Error())
                            continue
                        }

                        clt, err := tclient.MakeNewClient(appCfg.PCSSH, "root", "pc-node1")
                        if err != nil {
                            log.Error(err.Error())
                            continue
                        }

                        err = clt.APISCP(ctx.TODO(), []string{"/Users/almightykim/temp/update.sh", "root@pc-node1:/opt/pocket/bin/update.sh"}, roots[0].Password, tdefs.SSHServerListenPort, false, false)
                        if err != nil {
                            log.Errorf("ERROR : %v", err.Error())
                        }
                        err = clt.APISCP(ctx.TODO(), []string{"/Users/almightykim/temp/pocketd.update", "root@pc-node1:/opt/pocket/bin/pocketd.update"}, roots[0].Password, tdefs.SSHServerListenPort, false, false)
                        if err != nil {
                            log.Errorf("ERROR : %v", err.Error())
                            // exit with the same exit status as the failed command:
                            if clt.ExitStatus != 0 {
                            } else {
                            }
                        }
                        err = clt.APISSH(ctx.TODO(), []string{"/bin/bash", "/opt/pocket/bin/update.sh"}, roots[0].Password, false)
                        if err != nil {
                            log.Errorf("ERROR : %v", err.Error())
                        }
                        clt.Logout()

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
