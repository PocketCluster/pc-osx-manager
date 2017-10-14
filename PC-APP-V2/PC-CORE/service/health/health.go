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

func InitSystemHealthMonitor(appLife service.ServiceSupervisor, feeder route.ResponseFeeder) error {
    var (
        nodeBeaconC = make(chan service.Event)
        nodePcsshC  = make(chan service.Event)
        nodeOrchstC = make(chan service.Event)

        // service readiness checker
        readyBeaconC = make(chan service.Event)
        readyPcsshC  = make(chan service.Event)
        readyOrchstC = make(chan service.Event)
    )

    appLife.RegisterServiceWithFuncs(
        operation.ServiceMonitorSystemHealth,
        func() error {

            type MonitorSystemHealth map[string]interface{}

            var (
                // node status (beacon, pcssh, orchst) will be coalesced into one report
                //rpNodeStat   = routepath.RpathMonitorNodeStatus()
                rpSrvStat    = routepath.RpathMonitorServiceStatus()

                // timers
                timer        = time.NewTicker(time.Second * 2)
                failtimeout  = time.NewTicker(time.Minute)
                lastTS       = time.Now()

                // service ready checks
                readyMarker  = map[string]bool{
                    ivent.IventBeaconManagerSpawn:      false,
                    ivent.IventPcsshProxyInstanceSpawn: false,
                    ivent.IventOrchstInstanceSpawn:     false,
                }
                readyChecker = func(marker map[string]bool) bool {
                    for k := range marker {
                        if !marker[k] {
                            return false
                        }
                    }
                    return true
                }
            )

            // monitor pre-requisite services with timeout
            for {
                select {
                    case <- failtimeout.C: {
                        failtimeout.Stop()
                        timer.Stop()
                        return errors.Errorf("[HEALTH] fail to start health service")
                    }
                    case <- readyBeaconC: {
                        log.Infof("[HEALTH] beacon ready")
                        readyMarker[ivent.IventBeaconManagerSpawn] = true
                        if readyChecker(readyMarker) {
                            goto monstart
                        }
                    }
                    case <- readyPcsshC: {
                        log.Infof("[HEALTH] pcssh ready")
                        readyMarker[ivent.IventPcsshProxyInstanceSpawn] = true
                        if readyChecker(readyMarker) {
                            goto monstart
                        }
                    }
                    case <- readyOrchstC: {
                        log.Infof("[HEALTH] orchst ready")
                        readyMarker[ivent.IventOrchstInstanceSpawn] = true
                        if readyChecker(readyMarker) {
                            goto monstart
                        }
                    }
                }
            }

            monstart:
            failtimeout.Stop()
            log.Infof("[HEALTH] all required services are ready")

            for {
                select {
                    // service halt
                    case <- appLife.StopChannel(): {
                        timer.Stop()
                        return nil
                    }

                    // monitoring beacon
                    case re := <- nodeBeaconC: {
                        nodes, ok := re.Payload.([]map[string]string)
                        if !ok {
                            log.Debugf("[ERR] invalid unregistered node list type")
                            continue
                        }
                        _, err := json.Marshal(map[string]interface{} {"nodestat" : nodes})
                        if err != nil {
                            log.Debugf(err.Error())
                            continue
                        }
/*
                        err = feeder.FeedResponseForGet(rpNodeStat, string(data))
                        if err != nil {
                            log.Debugf(err.Error())
                        }
*/
                    }

                    // monitoring pcssh
                    case re := <- nodePcsshC: {
                        md, ok := re.Payload.(ivent.NodeStatusMeta)
                        if !ok {
                            er, ok := re.Payload.(error)
                            if ok {
                                log.Errorf("cannot fetch node status from pcssh %v", er)
                            } else {
                                log.Errorf("cannot fetch node status from pcssh w/ invalid data %v", md)
                            }
                            continue
                        }
                        if lastTS.After(md.TimeStamp) {
                            log.Errorf("invalid timestamp from pcssh last.ts %v | md.ts %v", lastTS, md.TimeStamp)
                            continue
                        }
                    }

                    // monitoring orchst
                    case re := <- nodeOrchstC: {
                        md, ok := re.Payload.(ivent.EngineStatusMeta)
                        if !ok {
                            er, ok := re.Payload.(error)
                            if ok {
                                log.Errorf("cannot fetch node status from orchst %v", er)
                            } else {
                                log.Errorf("cannot fetch node status from orchst w/ invalid data %v", md)
                            }
                            continue
                        }
                        if lastTS.After(md.TimeStamp) {
                            log.Errorf("invalid timestamp from orchst last.ts %v | md.ts %v", lastTS, md.TimeStamp)
                            continue
                        }
                    }

                    // service report
                    case ts := <- timer.C: {
                        // report services status
                        var (
                            response = MonitorSystemHealth{}
                            srvStatus = map[string]bool{}
                        )

                        // 1. report service status
                        sl := appLife.ServiceList()
                        for i, _ := range sl {
                            s := sl[i]
                            srvStatus[s.Tag()] = s.IsRunning()
                        }
                        response["services"] = srvStatus
                        data, err := json.Marshal(response)
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                        err = feeder.FeedResponseForGet(rpSrvStat, string(data))
                        if err != nil {
                            log.Debugf(err.Error())
                        }

                        // 2. report node status


                        // 3. request node status to services
                        appLife.BroadcastEvent(service.Event{
                            Name: ivent.IventMonitorNodeReqStatus,
                            Payload: ts,
                        })

                        // record ts
                        lastTS = ts
                    }
                }
            }

            return nil
        },
        service.BindEventWithService(ivent.IventMonitorNodeRespBeacon,   nodeBeaconC),
        service.BindEventWithService(ivent.IventMonitorNodeRespPcssh,    nodePcsshC),
        service.BindEventWithService(ivent.IventMonitorNodeRespOrchst,   nodeOrchstC),

        // service readiness checker
        service.BindEventWithService(ivent.IventBeaconManagerSpawn,      readyBeaconC),
        service.BindEventWithService(ivent.IventPcsshProxyInstanceSpawn, readyPcsshC),
        service.BindEventWithService(ivent.IventOrchstInstanceSpawn,     readyOrchstC))

    return nil
}
