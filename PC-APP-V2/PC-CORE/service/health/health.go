package health

import (
    "encoding/json"
    "strings"
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
    IPAddr        string        `json:"-"`
    Registered    bool          `json:"reged"`
    Bounded       bool          `json:"bnded"`
    PcsshOn       bool          `json:"pcssh"`
    OrchstOn      bool          `json:"orchst"`
}

type NodeStatMeta struct {
    Timestamp     int64         `json:"ts"`
    BeaconChecked bool          `json:"-"`
    PcsshChecked  bool          `json:"-"`
    OrchstChecked bool          `json:"-"`
    Nodes         []*NodeStat   `json:"nodes"`
}

func newNodeMetaWithTS(ts int64) *NodeStatMeta {
    return &NodeStatMeta{
        Timestamp: ts,
        Nodes:     []*NodeStat{},
    }
}

func (nm *NodeStatMeta) isReadyToReport() bool {
    return bool(nm.BeaconChecked && nm.PcsshChecked && nm.OrchstChecked)
}

func (nm *NodeStatMeta) updateBeaconStatus(bMeta ivent.BeaconNodeStatusMeta) {
    nm.BeaconChecked = true

    update_beacon:
    for _, bn := range bMeta.Nodes {
        for i, _ := range nm.Nodes {
            ns := nm.Nodes[i]
            if strings.HasPrefix(bn.Name, ns.Name) && bn.IPAddr == ns.IPAddr {
                ns.MacAddr    = bn.MacAddr
                ns.Registered = bn.Registered
                ns.Bounded    = bn.Bounded
                continue update_beacon
            }
        }

        // given pcssh node not found. so let's add
        nm.Nodes = append(nm.Nodes, &NodeStat{
            Name:        bn.Name,
            MacAddr:     bn.MacAddr,
            IPAddr:      bn.IPAddr,
            Registered:  bn.Registered,
            Bounded:     bn.Bounded,
        })
    }
}

func (nm *NodeStatMeta) updatePcsshStatus(pMeta ivent.PcsshNodeStatusMeta) {
    nm.PcsshChecked = true
    // https://github.com/golang/go/wiki/SliceTricks#additional-tricks
    // nl := nm.Nodes[:0]

    update_pcssh:
    for _, pn := range pMeta.Nodes {
        for i, _ := range nm.Nodes {
            ns := nm.Nodes[i]
            if strings.HasPrefix(pn.HostName, ns.Name) && pn.Addr == ns.IPAddr {
                ns.PcsshOn = true
                continue update_pcssh
            }
        }
        // given pcssh node not found. so let's add
        nm.Nodes = append(nm.Nodes, &NodeStat{
            Name:    pn.HostName,
            IPAddr:  pn.Addr,
            PcsshOn: true,
        })
    }
}

func (nm *NodeStatMeta) updateOrchstStatus(oMeta ivent.EngineStatusMeta) {
    nm.OrchstChecked = true

    update_orchst:
    for _, oe := range oMeta.Engines {
        for i, _ := range nm.Nodes {
            ns := nm.Nodes[i]
            if strings.HasPrefix(oe.Name, ns.Name) && oe.IP == ns.IPAddr {
                ns.OrchstOn = true
                continue update_orchst
            }
        }
        // given pcssh node not found. so let's add
        nm.Nodes = append(nm.Nodes, &NodeStat{
            Name:     oe.Name,
            IPAddr:   oe.IP,
            OrchstOn: true,
        })
    }
}

func (nm *NodeStatMeta) buildReport() ([]byte, error) {
    resp := route.ReponseMessage{
        "nodestat": {
            "status": true,
            "ts":    nm.Timestamp,
            "nodes": nm.Nodes,
        },
    }
    return json.Marshal(resp)
}

type TimedStats map[int64]*NodeStatMeta

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

func readyChecker(marker map[string]bool) bool {
    for k := range marker {
        if !marker[k] {
            return false
        }
    }
    return true
}

func reportNodeStats(meta *NodeStatMeta, fdr route.ResponseFeeder, rpath string) error {
    data, err := meta.buildReport()
    if err != nil {
        return err
    }
    return fdr.FeedResponseForGet(rpath, string(data))
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

            var (
                // node status (beacon, pcssh, orchst) will be coalesced into one report
                rpNodeStat   = routepath.RpathMonitorNodeStatus()
                rpSrvStat    = routepath.RpathMonitorServiceStatus()

                // timers
                // service checker
                sChkTimer   = time.NewTicker(time.Second * 2)
                // node status checker. node status checking frequent than 10 sec puts stress in system that
                // other services miss catching important signals such as stop.
                nStatTimer  = time.NewTicker(time.Second * 10)
                failtimeout = time.NewTicker(time.Minute)

                // stat collector
                timedStat   = make(TimedStats, 0)

                // service ready checks
                readyMarker  = map[string]bool{
                    ivent.IventBeaconManagerSpawn:      false,
                    ivent.IventPcsshProxyInstanceSpawn: false,
                    ivent.IventOrchstInstanceSpawn:     false,
                }
            )

            // monitor pre-requisite services with timeout
            for {
                select {
                    case <- failtimeout.C: {
                        failtimeout.Stop()
                        sChkTimer.Stop()
                        nStatTimer.Stop()
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
                        sChkTimer.Stop()
                        nStatTimer.Stop()
                        return nil
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
                            "srvcstats": {
                                "status": true,
                                "stats": srvStatus,
                            },
                        }
                        data, err := json.Marshal(resp)
                        if err != nil {
                            log.Debugf(err.Error())
                        }
                        err = feeder.FeedResponseForGet(rpSrvStat, string(data))
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
                        log.Infof("[HEALTH] beacon %v", md)

                        if meta.isReadyToReport() {
                            log.Errorf("[HEALTH] <<- (%v) ready to report", md.TimeStamp)
                            err := reportNodeStats(meta, feeder, rpNodeStat)
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
                        log.Infof("[HEALTH] pcssh %v", md)

                        if meta.isReadyToReport() {
                            log.Errorf("[HEALTH] <<- (%v) ready to report", md.TimeStamp)
                            err := reportNodeStats(meta, feeder, rpNodeStat)
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
                        log.Infof("[HEALTH] orchst %v", md)

                        if meta.isReadyToReport() {
                            log.Errorf("[HEALTH] <<- (%v) ready to report", md.TimeStamp)
                            err := reportNodeStats(meta, feeder, rpNodeStat)
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
        service.BindEventWithService(ivent.IventBeaconManagerSpawn,      readyBeaconC),
        service.BindEventWithService(ivent.IventPcsshProxyInstanceSpawn, readyPcsshC),
        service.BindEventWithService(ivent.IventOrchstInstanceSpawn,     readyOrchstC))

    return nil
}
