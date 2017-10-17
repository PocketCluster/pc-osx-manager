package health

import (
    "encoding/json"
    "strings"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/service/ivent"
)

type NodeStat struct {
    Name          string        `json:"name"`
    MacAddr       string        `json:"mac"`
    IPAddr        string        `json:"-"`
    Registered    bool          `json:"rgstd"`
    Bounded       bool          `json:"bound"`
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
        "node-stat": {
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
