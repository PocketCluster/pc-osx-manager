package health

import (
    "encoding/json"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/extlib/pcssh/sshproc"
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
                timer        = time.NewTicker(time.Second)
                // node status (beacon, pcssh, orchst) will be coalesced into one report
                rpNodeStat   = routepath.RpathMonitorNodeStatus()
                rpSrvStat    = routepath.RpathMonitorServiceStatus()
                failtimeout  = time.NewTicker(time.Minute)
                readyMarker  = map[string]bool{
                    ivent.IventBeaconManagerSpawn:        false,
                    sshproc.EventPCSSHServerProxyStarted: false,
                    ivent.IventOrchstInstanceSpawn:       false,
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
                        readyMarker[sshproc.EventPCSSHServerProxyStarted] = true
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
                        data, err := json.Marshal(map[string]interface{} {"nodestat" : nodes})
                        if err != nil {
                            log.Debugf(err.Error())
                            continue
                        }
                        err = feeder.FeedResponseForGet(rpNodeStat, string(data))
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                    }

                    // monitoring pcssh
                    case <- nodePcsshC: {
                    }

                    // monitoring orchst
                    case <- nodeOrchstC: {
                    }

                    // service report
                    case <- timer.C: {
                        // report services status
                        var (
                            response = MonitorSystemHealth{}
                            srvStatus = map[string]bool{}
                        )

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
                    }
                }
            }

            return nil
        },
        service.BindEventWithService(ivent.IventMonitorNodeRsltBeacon, nodeBeaconC),
        service.BindEventWithService(sshproc.EventPCSSHNodeListResult, nodePcsshC),
        service.BindEventWithService(ivent.IventMonitorNodeRsltOrchst, nodeOrchstC),

        // service readiness checker
        service.BindEventWithService(ivent.IventBeaconManagerSpawn,    readyBeaconC),
        service.BindEventWithService(sshproc.EventPCSSHServerProxyStarted, readyPcsshC),
        service.BindEventWithService(ivent.IventOrchstInstanceSpawn,   readyOrchstC))

    return nil
}
