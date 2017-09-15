package health

import (
    "encoding/json"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
)

func InitSystemHealthMonitor(appLife route.Router) error {
    var (
        unregNodeC = make(chan service.Event)
        regNodeC = make(chan service.Event)
    )
    appLife.RegisterServiceWithFuncs(
        operation.ServiceMonitorSystemHealth,
        func() error {

            type MonitorSystemHealth map[string]interface{}

            var (
                timer = time.NewTicker(time.Second)
                rpUnregNode = routepath.RpathMonitorNodeUnregistered()
                rpRegNode = routepath.RpathMonitorNodeRegistered()
                rpSrvStat = routepath.RpathMonitorServiceStatus()
            )

            for {
                select {
                case <- appLife.StopChannel(): {
                    timer.Stop()
                    return nil
                }
                    // monitoring unregistered nodes
                case re := <- unregNodeC: {
                    nodes, ok := re.Payload.([]map[string]string)
                    if !ok {
                        log.Debugf("[ERR] invalid unregistered node list type")
                        continue
                    }
                    data, err := json.Marshal(map[string]interface{} {"unregistered" : nodes})
                    if err != nil {
                        log.Debugf(err.Error())
                        continue
                    }
                    err = FeedResponseForGet(rpUnregNode, string(data))
                    if err != nil {
                        log.Debugf(err.Error())
                    }
                }
                    // monitoring registered nodes
                case re := <- regNodeC: {
                    nodes, ok := re.Payload.([]map[string]string)
                    if !ok {
                        log.Debugf("[ERR] invalid registered node list type")
                        continue
                    }
                    data, err := json.Marshal(map[string]interface{} {"registered" : nodes})
                    if err != nil {
                        log.Debugf(err.Error())
                        continue
                    }
                    err = FeedResponseForGet(rpRegNode, string(data))
                    if err != nil {
                        log.Debugf(err.Error())
                    }
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
                    err = FeedResponseForGet(rpSrvStat, string(data))
                    if err != nil {
                        log.Debugf(err.Error())
                    }
                }
                }
            }

            return nil
        },
        service.BindEventWithService(ivent.IventMonitorUnregisteredNode, unregNodeC),
        service.BindEventWithService(ivent.IventMonitorRegisteredNode,   regNodeC))

    return nil
}
