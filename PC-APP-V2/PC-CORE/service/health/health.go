package health

import (
    "encoding/json"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
)

func InitSystemHealthMonitor(appLife service.ServiceSupervisor, feeder route.ResponseFeeder) {
    var (
        nodeBeaconC = make(chan service.Event)
        nodePcsshC  = make(chan service.Event)
        nodeOrchstC = make(chan service.Event)

        // service readiness checker
        readyDscvryC = make(chan service.Event)
        readyRgstryC = make(chan service.Event)
        readyNameC   = make(chan service.Event)
        readyBeaconC = make(chan service.Event)
        readyPcsshC  = make(chan service.Event)
        readyOrchstC = make(chan service.Event)
        readyVboxC   = make(chan service.Event)
        innerErrC    = make(chan service.Event)
        stopMonitorC = make(chan service.Event)
    )

    appLife.RegisterServiceWithFuncs(
        operation.ServiceMonitorSystemHealth,
        func() error {

            var (
                // node status (beacon, pcssh, orchst) will be coalesced into one report
                rpNodeStat   = routepath.RpathMonitorNodeStatus()
                rpSrvcTimeup = routepath.RpathNotiSrvcOnlineTimeup()
                rpSrvcStat   = routepath.RpathMonitorServiceStatus()

                // --- timers ---
                // service checker
                sChkTimer   = time.NewTicker(time.Second * 2)
                // node status checker. node status checking frequent than 10 sec puts stress in system that
                // other services miss catching important signals such as stop.
                nStatTimer  = time.NewTicker(time.Second * 10)
                // this is to wait timer for other services to start. Especially for this timeout, we'll give 90 secs
                // this also works as a timeout for core node to boot. make sure core node boot in 90 sec
                failTimeout = time.NewTicker(time.Second * 90)
                // app start timeup counter. This should only be triggered after 1 minute
                nodeOnlineTimeup *time.Ticker = nil

                // stat collector
                timedStat   = make(TimedStats, 0)

                // service ready checks
                readyMarker  = map[string]bool{
                    ivent.IventDiscoveryInstanceSpwan:  false,
                    ivent.IventRegistryInstanceSpawn:   false,
                    ivent.IventNameServerInstanceSpawn: false,
                    ivent.IventBeaconManagerSpawn:      false,
                    ivent.IventPcsshProxyInstanceSpawn: false,
                    ivent.IventOrchstInstanceSpawn:     false,
                    ivent.IventVboxCtrlInstanceSpawn:   false,
                }

                // ignore core node error. the only critical error for node report is core node death.
                // make sure core node error reported to front-end after node oneline ticker fires
                checkCoreError = false

                // stop monitor request
                shouldStop = false
            )

            // monitor pre-requisite services with timeout
            for {
                select {
                    case <- failTimeout.C: {
                        failTimeout.Stop()
                        sChkTimer.Stop()
                        nStatTimer.Stop()

                        data, err := json.Marshal(route.ReponseMessage{
                            "srvc-timeup": {
                                "status": false,
                                "error": "unable to start internal services on time",
                            }})
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                        err = feeder.FeedResponseForGet(rpSrvcTimeup, string(data))
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                        return errors.Errorf("[HEALTH] fail to start health service")
                    }
                    // !!! Any error happens in initializing internal service is critical one. !!!
                    // provide feedback upon receiving one and stop application
                    case ee := <-innerErrC: {
                        failTimeout.Stop()
                        sChkTimer.Stop()
                        nStatTimer.Stop()

                        irr, ok := ee.Payload.(error)
                        if !ok {
                            log.Error("[HEALTH] invalid internal error message type")
                        }
                        data, err := json.Marshal(route.ReponseMessage{
                            "srvc-timeup": {
                                "status": false,
                                "error": irr.Error(),
                            }})
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                        err = feeder.FeedResponseForGet(rpSrvcTimeup, string(data))
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                        return errors.Errorf("[HEALTH] fail to start health service")
                    }

                    // discovery
                    case <- readyDscvryC: {
                        log.Infof("[HEALTH] discovery ready")
                        readyMarker[ivent.IventDiscoveryInstanceSpwan] = true
                        if readyChecker(readyMarker) {
                            goto monstart
                        }
                    }
                    // registry
                    case <- readyRgstryC: {
                        log.Infof("[HEALTH] registry ready")
                        readyMarker[ivent.IventRegistryInstanceSpawn] = true
                        if readyChecker(readyMarker) {
                            goto monstart
                        }
                    }
                    // named
                    case <- readyNameC: {
                        log.Infof("[HEALTH] named ready")
                        readyMarker[ivent.IventNameServerInstanceSpawn] = true
                        if readyChecker(readyMarker) {
                            goto monstart
                        }
                    }
                    // beacon
                    case <- readyBeaconC: {
                        log.Infof("[HEALTH] beacon ready")
                        readyMarker[ivent.IventBeaconManagerSpawn] = true
                        if readyChecker(readyMarker) {
                            goto monstart
                        }
                    }
                    // pcssh
                    case <- readyPcsshC: {
                        log.Infof("[HEALTH] pcssh ready")
                        readyMarker[ivent.IventPcsshProxyInstanceSpawn] = true
                        if readyChecker(readyMarker) {
                            goto monstart
                        }
                    }
                    // orchst
                    case <- readyOrchstC: {
                        log.Infof("[HEALTH] orchst ready")
                        readyMarker[ivent.IventOrchstInstanceSpawn] = true
                        if readyChecker(readyMarker) {
                            goto monstart
                        }
                    }
                    // vbox
                    case <- readyVboxC: {
                        log.Infof("[HEALTH] vbox ready")
                        readyMarker[ivent.IventVboxCtrlInstanceSpawn] = true
                        if readyChecker(readyMarker) {
                            goto monstart
                        }
                    }
                }
            }

            monstart:
            failTimeout.Stop()
            // notify frontend that all services started as intented
            data, err := json.Marshal(route.ReponseMessage{"srvc-timeup": {"status": true}})
            if err != nil {
                log.Debugf(err.Error())
            }
            err = feeder.FeedResponseForGet(rpSrvcTimeup, string(data))
            if err != nil {
                log.Debugf(err.Error())
            }
            nodeOnlineTimeup = time.NewTicker(time.Minute)
            log.Infof("[HEALTH] all required services are ready")

            for {
                select {
                    // service halt
                    case <- appLife.StopChannel(): {
                        sChkTimer.Stop()
                        nStatTimer.Stop()
                        nodeOnlineTimeup.Stop()
                        return nil
                    }
                    // stop monitoring
                    case <- stopMonitorC: {
                        shouldStop = true
                        sChkTimer.Stop()
                        nStatTimer.Stop()
                        nodeOnlineTimeup.Stop()

                        // we return now
                        if timedStat.isReadyToRequest() {
                            log.Info("[HEALTH] stop monitoring...")
                            appLife.BroadcastEvent(service.Event{Name:ivent.IventMonitorStopResult})
                            return nil
                        }
                    }

                    // node should have been all online timeup (should fire only once forfrontend to prep)
                    case <- nodeOnlineTimeup.C: {
                        // shoul not nullify start timeup. It will crash
                        nodeOnlineTimeup.Stop()
                        // do not ignore core node error after this point
                        checkCoreError = true

                        data, err := json.Marshal(route.ReponseMessage{"node-timeup": {"status": true}})
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                        err = feeder.FeedResponseForGet(routepath.RpathNotiNodeOnlineTimeup(), string(data))
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                        log.Infof("[HEALTH] node online due is up!")
                    }

                    // report services status
                    case <- sChkTimer.C: {
                        var (
                            srvStatus = map[string]bool{}
                        )

                        // 1. report service status
                        sl := appLife.ServiceList()
                        for i, _ := range sl {
                            s := sl[i]
                            srvStatus[s.Tag()] = s.IsRunning()
                        }
                        resp := route.ReponseMessage{
                            "srvc-stat": {
                                "status": true,
                                "srvcs": srvStatus,
                            },
                        }
                        data, err := json.Marshal(resp)
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                        err = feeder.FeedResponseForGet(rpSrvcStat, string(data))
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                    }

                    // request node status to services
                    case ts := <- nStatTimer.C: {
                        var (
                            tmark = ts.Unix()
                        )

                        // clear requests older than certain time period
                        timedStat.cleanRequestBefore(tmark - int64(30))

                        // unless requested stat report is being cleared, we will not make another request
                        if timedStat.isReadyToRequest() {
                            timedStat[tmark] = newNodeMetaWithTS(tmark)
                            appLife.BroadcastEvent(service.Event{
                                Name: ivent.IventMonitorNodeReqStatus,
                                Payload: tmark,
                            })
                            log.Infof("[HEALTH] ->> (%v) new stat request made", tmark)
                        }
                    }

                    // monitoring beacon
                    case re := <- nodeBeaconC: {
                        md, ok := re.Payload.(ivent.BeaconNodeStatusMeta)
                        if !ok {
                            log.Errorf("[HEALTH] [ERR] cannot fetch node status from pcssh w/ invalid data %v", md)
                            continue
                        }
                        if md.Error != nil {
                            log.Errorf(md.Error.Error())
                            continue
                        }
                        meta, ok := timedStat[md.TimeStamp]
                        if !ok {
                            log.Errorf("[HEALTH] [ERR] timestamp %v record for reported stat does not exists", md.TimeStamp)
                            continue
                        }
                        meta.updateBeaconStatus(md)

                        if meta.isReadyToReport() {
                            if shouldStop {
                                log.Info("[HEALTH] stop monitoring...")
                                appLife.BroadcastEvent(service.Event{Name:ivent.IventMonitorStopResult})
                                return nil
                            }
                            log.Errorf("[HEALTH] <<- (%v) ready to report", md.TimeStamp)
                            err := reportNodeStats(meta, feeder, rpNodeStat, checkCoreError)
                            if err != nil {
                                log.Errorf("[HEALTH] [ERR] unable to report node stat %v", err)
                            }
                            timedStat.removeStatForTimestamp(md.TimeStamp)
                        }
                    }

                    // monitoring pcssh
                    case re := <- nodePcsshC: {
                        md, ok := re.Payload.(ivent.PcsshNodeStatusMeta)
                        if !ok {
                            log.Errorf("[HEALTH] [ERR] cannot fetch node status from pcssh w/ invalid data %v", md)
                            continue
                        }
                        if md.Error != nil {
                            log.Errorf(md.Error.Error())
                            continue
                        }
                        meta, ok := timedStat[md.TimeStamp]
                        if !ok {
                            log.Errorf("[HEALTH] [ERR] timestamp %v record for reported stat does not exists", md.TimeStamp)
                            continue
                        }
                        meta.updatePcsshStatus(md)
                        log.Infof("[HEALTH] PCSSH META %v", md)

                        if meta.isReadyToReport() {
                            if shouldStop {
                                log.Info("[HEALTH] stop monitoring...")
                                appLife.BroadcastEvent(service.Event{Name:ivent.IventMonitorStopResult})
                                return nil
                            }
                            log.Errorf("[HEALTH] <<- (%v) ready to report", md.TimeStamp)
                            err := reportNodeStats(meta, feeder, rpNodeStat, checkCoreError)
                            if err != nil {
                                log.Errorf("[HEALTH] [ERR] unable to report node stat %v", err)
                            }
                            timedStat.removeStatForTimestamp(md.TimeStamp)
                        }
                    }

                    // monitoring orchst
                    case re := <- nodeOrchstC: {
                        md, ok := re.Payload.(ivent.EngineStatusMeta)
                        if !ok {
                            log.Errorf("[HEALTH] [ERR] cannot fetch node status from orchst w/ invalid data %v", md)
                            continue
                        }
                        if md.Error != nil {
                            log.Errorf(md.Error.Error())
                            continue
                        }
                        meta, ok := timedStat[md.TimeStamp]
                        if !ok {
                            log.Errorf("[HEALTH] [ERR] timestamp %v record for reported stat does not exists", md.TimeStamp)
                            continue
                        }
                        meta.updateOrchstStatus(md)

                        if meta.isReadyToReport() {
                            if shouldStop {
                                log.Info("[HEALTH] stop monitoring...")
                                appLife.BroadcastEvent(service.Event{Name:ivent.IventMonitorStopResult})
                                return nil
                            }
                            log.Errorf("[HEALTH] <<- (%v) ready to report", md.TimeStamp)
                            err := reportNodeStats(meta, feeder, rpNodeStat, checkCoreError)
                            if err != nil {
                                log.Errorf("[HEALTH] [ERR] unable to report node stat %v", err)
                            }
                            timedStat.removeStatForTimestamp(md.TimeStamp)
                        }
                    }
                }
            }

            return nil
        },
        service.BindEventWithService(ivent.IventMonitorNodeRespBeacon,   nodeBeaconC),
        service.BindEventWithService(ivent.IventMonitorNodeRespPcssh,    nodePcsshC),
        service.BindEventWithService(ivent.IventMonitorNodeRespOrchst,   nodeOrchstC),

        // service readiness checker
        service.BindEventWithService(ivent.IventDiscoveryInstanceSpwan,  readyDscvryC),
        service.BindEventWithService(ivent.IventRegistryInstanceSpawn,   readyRgstryC),
        service.BindEventWithService(ivent.IventNameServerInstanceSpawn, readyNameC),
        service.BindEventWithService(ivent.IventBeaconManagerSpawn,      readyBeaconC),
        service.BindEventWithService(ivent.IventPcsshProxyInstanceSpawn, readyPcsshC),
        service.BindEventWithService(ivent.IventOrchstInstanceSpawn,     readyOrchstC),
        service.BindEventWithService(ivent.IventVboxCtrlInstanceSpawn,   readyVboxC),

        // internal error collector
        service.BindEventWithService(ivent.IventInternalSpawnError,      innerErrC),
        service.BindEventWithService(ivent.IventMonitorStopRequest,      stopMonitorC))
}
