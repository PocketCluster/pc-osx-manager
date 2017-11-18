package health

import (
    "encoding/json"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/defaults"
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
    Nodes         []NodeStat    `json:"nodes"`
}

func newNodeMetaWithTS(ts int64) *NodeStatMeta {
    return &NodeStatMeta{
        Timestamp: ts,
        Nodes:     []NodeStat{},
    }
}

func (nm *NodeStatMeta) isReadyToReport() bool {
    return bool(nm.BeaconChecked && nm.PcsshChecked && nm.OrchstChecked)
}

func (nm *NodeStatMeta) updateBeaconStatus(bMeta ivent.BeaconNodeStatusMeta) {
    nm.BeaconChecked = true

    update_beacon:
    for _, bn := range bMeta.Nodes {
        // for core node (core node might not have ip address due to internal pcssh issue)
        if bn.Name == defaults.PocketClusterCoreName {
            for i := range nm.Nodes {
                if nm.Nodes[i].Name == defaults.PocketClusterCoreName {
                    nm.Nodes[i].MacAddr    = bn.MacAddr
                    nm.Nodes[i].IPAddr     = bn.IPAddr
                    nm.Nodes[i].Registered = bn.Registered
                    nm.Nodes[i].Bounded    = bn.Bounded
                    continue update_beacon
                }
            }

        } else {
            for i := range nm.Nodes {
                if bn.Name == nm.Nodes[i].Name && bn.IPAddr == nm.Nodes[i].IPAddr {
                    nm.Nodes[i].MacAddr    = bn.MacAddr
                    nm.Nodes[i].Registered = bn.Registered
                    nm.Nodes[i].Bounded    = bn.Bounded
                    continue update_beacon
                }
            }
        }
        // given beacon node not found. so let's add
        nm.Nodes = append(nm.Nodes, NodeStat{
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
        // for core node (core node might not have ip address due to internal pcssh issue)
        if pn.HostName == defaults.PocketClusterCoreName {
            for i := range nm.Nodes {
                if nm.Nodes[i].Name == defaults.PocketClusterCoreName { // && pn.Addr == defaults.PocketClusterCodeInteralAddr (we don't check address for now)
                    nm.Nodes[i].PcsshOn = true
                    continue update_pcssh
                }
            }
            // given pcssh core node not found. let's add w/o address
            nm.Nodes = append(nm.Nodes, NodeStat{
                Name:    pn.HostName,
                PcsshOn: true,
            })

        } else {
            for i := range nm.Nodes {
                if pn.HostName == nm.Nodes[i].Name && pn.Addr == nm.Nodes[i].IPAddr {
                    nm.Nodes[i].PcsshOn = true
                    continue update_pcssh
                }
            }
            // given pcssh node not found. so let's add
            nm.Nodes = append(nm.Nodes, NodeStat{
                Name:    pn.HostName,
                IPAddr:  pn.Addr,
                PcsshOn: true,
            })
        }
    }
}

func (nm *NodeStatMeta) updateOrchstStatus(oMeta ivent.EngineStatusMeta) {
    nm.OrchstChecked = true

    update_orchst:
    for _, oe := range oMeta.Engines {
        // for core node (core node might not have ip address due to internal pcssh issue)
        if oe.Name == defaults.PocketClusterCoreName {
            for i := range nm.Nodes {
                if nm.Nodes[i].Name == defaults.PocketClusterCoreName {
                    nm.Nodes[i].IPAddr   = oe.IP
                    nm.Nodes[i].OrchstOn = true
                    continue update_orchst
                }
            }

        } else {
            for i := range nm.Nodes {
                if oe.Name == nm.Nodes[i].Name && oe.IP == nm.Nodes[i].IPAddr {
                    nm.Nodes[i].OrchstOn = true
                    continue update_orchst
                }
            }
        }
        // given orchst node not found. so let's add
        nm.Nodes = append(nm.Nodes, NodeStat{
            Name:     oe.Name,
            IPAddr:   oe.IP,
            OrchstOn: true,
        })
    }
}

func (nm *NodeStatMeta) buildReport(checkCoreError bool) ([]byte, error) {
    var (
        resp = route.ReponseMessage{
            "node-stat": {
                "status": true,
                "ts":     nm.Timestamp,
                "nodes":  nm.Nodes,
            },
        }
    )
    // find core node and build error if core is not normal
    if checkCoreError {
        var (
            cFound = false
            err error = nil
        )
        for _, ns := range nm.Nodes {
            if ns.Name == defaults.PocketClusterCoreName {
                // core node is found
                cFound = true
                // we might want to count ip address but that's to restrictive. Let's only count what pc-master sees
                if !(ns.Registered && ns.Bounded && ns.PcsshOn && ns.OrchstOn) {
                    err = errors.Errorf("[HEALTH-CRITICAL] core node has an issue. Registered[%v], Bounded[%v] PcsshOn[%v] OrchstOn[%v]",
                        ns.Registered, ns.Bounded, ns.PcsshOn, ns.OrchstOn)
                }
                break
            }
        }
        // core node is not found. this is even more serious issue
        if !cFound {
            err = errors.Errorf("[HEALTH-CRITICAL] core node not found")
        }
        // include error if exists. this is critical
        if err != nil {
            resp["node-stat"]["status"] = false
            resp["node-stat"]["error"] = err.Error()
        }
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

func reportNodeStats(meta *NodeStatMeta, fdr route.ResponseFeeder, rpath string, checkCoreError bool) error {
    data, err := meta.buildReport(checkCoreError)
    if err != nil {
        return err
    }
    return fdr.FeedResponseForGet(rpath, string(data))
}
