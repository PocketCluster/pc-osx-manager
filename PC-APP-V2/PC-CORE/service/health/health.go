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

type NodeStat struct {
    Name          string        `json:"name"`
    MacAddr       string        `json:"mac"`
    Registered    bool          `json:"reged"`
    Bounded       bool          `json:"bnded"`
    PcsshOn       bool          `json:"pcssh"`
    OrchstOn      bool          `json:"orchst"`
}

type NodeMeta struct {
    Timestamp     int64         `json:"ts"`
    BeaconChecked bool          `json:"-"`
    PcsshChecked  bool          `json:"-"`
    OrchstChecked bool          `json:"-"`
    Stat          []*NodeStat   `json:"stat"`
}

func (nm *NodeMeta) isReadyToReport() bool {
    //return bool(nm.BeaconChecked && nm.PcsshChecked && nm.OrchstChecked)
    return bool(nm.PcsshChecked && nm.OrchstChecked)
}

type TimedStats map[int64]*NodeMeta

func (t TimedStats) removeStatForTimestamp(ts int64) {
    delete(t, ts)
}

func (t TimedStats) cleanRequestBefore(ts int64) {
    if len(t) == 0 {
        return
    }
    var tks []int64 = []int64{}
    for tk := range t {
        tks = append(tks, tk)
    }

    for _, tk := range tks {
        if tk <= ts {
            log.Warnf("[HEALTH] [WARN] deleting old tk %v", tk)
            delete(t, tk)
        }
    }
}

func (t TimedStats) isReadyToRequest() bool {
    return len(t) == 0
}

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
                timedStat    = make(TimedStats, 0)

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
                    case <- nodeBeaconC: {
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
                        meta.PcsshChecked = true

                        if meta.isReadyToReport() {
                            log.Errorf("[HEALTH] <<- (%v) ready to report", md.TimeStamp)
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
                        meta.OrchstChecked = true

                        if meta.isReadyToReport() {
                            log.Errorf("[HEALTH] <<- (%v) ready to report", md.TimeStamp)
                            timedStat.removeStatForTimestamp(md.TimeStamp)
                        }
                    }

                    // service report
                    case ts := <- timer.C: {
                        // report services status
                        var (
                            response = MonitorSystemHealth{}
                            srvStatus = map[string]bool{}
                            tmark = ts.Unix()
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

                        // 3. --- request node status to services ---
                        // clear requests older than 10 sec
                        timedStat.cleanRequestBefore(tmark - int64(10))

                        // unless requested stat report is being cleared, we will not make another request
                        if timedStat.isReadyToRequest() {
                            timedStat[tmark] = &NodeMeta{
                                Timestamp: tmark,
                            }
                            appLife.BroadcastEvent(service.Event{
                                Name: ivent.IventMonitorNodeReqStatus,
                                Payload: tmark,
                            })
                            log.Infof("[HEALTH] ->> (%v) new stat request made", tmark)
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
        service.BindEventWithService(ivent.IventBeaconManagerSpawn,      readyBeaconC),
        service.BindEventWithService(ivent.IventPcsshProxyInstanceSpawn, readyPcsshC),
        service.BindEventWithService(ivent.IventOrchstInstanceSpawn,     readyOrchstC))

    return nil
}
