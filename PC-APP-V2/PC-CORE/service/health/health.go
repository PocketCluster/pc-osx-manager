package health

import (
    "encoding/json"
    "time"

    log "github.com/Sirupsen/logrus"
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
            )

            // TODO : we need a trigger to make sure beacon/ pcssh/ orchst all have started

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
        service.BindEventWithService(ivent.IventMonitorNodeRsltOrchst, nodeOrchstC))

    return nil
}
