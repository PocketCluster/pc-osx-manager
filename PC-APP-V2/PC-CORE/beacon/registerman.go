package beacon

import (
    "net"
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/model"
)
import (
    "github.com/davecgh/go-spew/spew"
)

const (
    unregNodeRefreshPeriod time.Duration = time.Duration(UnboundedTimeout * 4)
)

type RegisterManger interface {
    MonitoringMasterSearchData(searchD mcast.CastPack, ts time.Time) error
    UnregisteredNodeList(ts time.Time) ([]map[string]string, error)
    RegisterMonitoredNodes(ts time.Time) error
    GuideNodeRegistrationWithBeacon(beaconD ucast.BeaconPack, ts time.Time) error
}

type monitorMeta struct {
    addedTS time.Time
    net.UDPAddr
    *slagent.PocketSlaveAgentMeta
}

type registerManager struct {
    sync.Mutex

    isRegisteringNode bool
    nodeList          []monitorMeta
    *beaconManger
}

func NewNodeRegisterManager(master BeaconManger) (RegisterManger, error) {
    bm, ok := master.(*beaconManger)
    if !ok {
        return nil, errors.Errorf("invalid meanager type")
    }
    return &registerManager{
        isRegisteringNode: false,
        beaconManger:      bm,
        nodeList:          make([]monitorMeta, 0),
        }, nil
}

/*
func (b *beaconManger) TransitionWithSearchData(searchD mcast.CastPack, ts time.Time) error {
    var (
        bcFound bool            = false
        mc MasterBeacon         = nil
        err error               = nil
        slave *model.SlaveNode  = nil
        usm *slagent.PocketSlaveAgentMeta = nil
        state MasterBeaconState = MasterBounded
    )

    usm, err = slagent.UnpackedSlaveMeta(searchD.Message)
    if err != nil {
        return errors.WithStack(err)
    }

    //log.Debugf("[SEARCH-RX] %v\n%v ", searchD.Address.IP.String(), spew.Sdump(usm))
    log.Debugf("[SEARCH-RX] %v\n%v ", searchD.Address.IP.String())

    // this packet looks for something else
    if len(usm.MasterBoundAgent) != 0 && usm.MasterBoundAgent != b.clusterID {
        log.Debugf("[SEARCH-RX] this packet belong to other master | usm.DiscoveryAgent.MasterBoundAgent %v | b.clusterID %v", usm.MasterBoundAgent, b.clusterID)
        return nil
    }

    // remove discarded beacon
    pruneBeaconList(b)

    // check if beacon for this packet exists
    var bLen int = len(b.beaconList)
    for i := 0; i < bLen; i++  {
        bc := b.beaconList[i]
        if bc.SlaveNode().SlaveID == usm.SlaveID {

            // this beacons are created and waiting for an input
            state = bc.CurrentState()
            if state == MasterInit || state == MasterBindBroken {
                log.Debugf("[SEARCH-NODE-FOUND] (%s | %s) ", bc.SlaveNode().SlaveID, bc.CurrentState().String())
                return bc.TransitionWithSlaveMeta(&searchD.Address, usm, ts)
            }

            // if beacon is not in searching state, then mark and we've found target
            bcFound = true
            break
        }
    }

    // since we've not found, create new beacon
    if !bcFound {
        slave = model.NewSlaveNode(b)
        // we'll ignore message for now
        b.notiReceiver.BeaconEventPrepareJoin(slave)
        mc, err = NewMasterBeacon(MasterInit, slave, b.commChannel, b)
        if err != nil {
            return errors.WithStack(err)
        }
        insertMasterBeacon(b, mc)
        return mc.TransitionWithSlaveMeta(&searchD.Address, usm, ts)
    }

    return errors.Errorf("[SEARCH-ERR] TransitionWithSearchData reaches at the end. *this should never happen, and might be a malicious attempt*")
}
*/

func (r *registerManager) MonitoringMasterSearchData(searchD mcast.CastPack, ts time.Time) error {
    usm, err := slagent.UnpackedSlaveMeta(searchD.Message)
    if err != nil {
        // (Ignore) there are way too many unpackable packages.
        return nil
    }
    if len(usm.MasterBoundAgent) != 0 {
        // this is registered to a master
        return nil
    }

    log.Debugf("[REGISTER-RX] %v\n%v ", searchD.Address.IP.String(), spew.Sdump(usm))

    // remove discarded beacon
    pruneBeaconList(r.beaconManger)

    // check if total node slots are filled
    var bLen int = len(r.beaconManger.beaconList)
    if defaults.TotalPossibleSlaveNodeCount <= bLen {
        return errors.Errorf("[REGISTER-RX] cannot exceed total possible node count")
    }
    // check if beacon for this packet exists
    for i := 0; i < bLen; i++  {
        bc := r.beaconManger.beaconList[i]
        if bc.SlaveNode().SlaveID == usm.SlaveID {
            return errors.Errorf("[REGISTER-RX] node %v should not exist", usm.SlaveID)
        }
    }

    // protect nodelist, isNodeFilled
    r.Lock()
    defer r.Unlock()

    if r.isRegisteringNode {
        return errors.Errorf("[REGISTER-RX] node registration is underway. please try another batch after finishing this round.")
    }

    /*
     * monitoring takes all the available, unregistered node on the net, refreshing the list for 4 unbounded packet frame (12sec)
     */
    var (
        nl = make([]monitorMeta, 0)
        nLen = len(r.nodeList)
    )
    // take only fresh, unidentical ones
    for i := 0; i < nLen; i++ {
        n := r.nodeList[i]
        if ts.Sub(n.addedTS) < unregNodeRefreshPeriod && n.SlaveID != usm.SlaveID {
            nl = append(nl, n)
        }
    }
    // now add this packet
    r.nodeList = append(nl, monitorMeta {
        addedTS: ts,
        UDPAddr: searchD.Address,
        PocketSlaveAgentMeta: usm,
    })
    return nil
}


func (r *registerManager) UnregisteredNodeList(ts time.Time) ([]map[string]string, error) {
    var list = make([]map[string]string, 0)

    // protect nodelist
    r.Lock()
    defer r.Unlock()

    var nLen int = len(r.nodeList)
    for i := 0; i < nLen; i++ {
        n := r.nodeList[i]
        list = append(list,
            map[string]string{
                "mac":  n.SlaveID,
                "addr": n.UDPAddr.IP.String(),
            })
    }

    return list, nil
}

func (r *registerManager) RegisterMonitoredNodes(ts time.Time) error {
    // protect nodelist, isNodeFilled
    r.Lock()
    defer r.Unlock()

    // no more registration
    r.isRegisteringNode = true

    // select the ones to be registered
    var nLen = len(r.nodeList)
    for i := 0; i < nLen; i++ {
        n := r.nodeList[i]

        slave := model.NewSlaveNode(r.beaconManger)
        // we'll ignore message for now
        r.beaconManger.notiReceiver.BeaconEventPrepareJoin(slave)
        mc, err := NewMasterBeacon(MasterInit, slave, r.beaconManger.commChannel, r.beaconManger)
        if err != nil {
            log.Errorf("[REGISTER-TX] %v", err.Error())
            continue
        }
        insertMasterBeacon(r.beaconManger, mc)
        err = mc.TransitionWithSlaveMeta(&n.UDPAddr, n.PocketSlaveAgentMeta, ts)
        if err != nil {
            log.Errorf("[REGISTER-TX] %v", err.Error())
        }
    }

    return nil
}

/*
func (b *beaconManger) TransitionWithBeaconData(beaconD ucast.BeaconPack, ts time.Time) error {
    var (
        err error = nil
        usm *slagent.PocketSlaveAgentMeta = nil
    )

    // suppose we've sort out what this is.
    usm, err = slagent.UnpackedSlaveMeta(beaconD.Message)
    if err != nil {
        return errors.WithStack(err)
    }

    log.Debugf("[BEACON-RX] %v\n%v", beaconD.Address.IP.String(), spew.Sdump(usm))

    // this packet looks for something else
    if len(usm.MasterBoundAgent) != 0 && usm.MasterBoundAgent != b.clusterID {
        return nil
    }

    // remove discarded beacon
    pruneBeaconList(b)

    // check if beacon for this packet exists
    var bLen int = len(b.beaconList)
    for i := 0; i < bLen; i++  {
        bc := b.beaconList[i]
        if bc.SlaveNode().SlaveID == usm.SlaveID {
            switch bc.CurrentState() {
                case MasterInit:
                    fallthrough
                case MasterBindBroken:
                    fallthrough
                case MasterDiscarded: {
                    log.Debugf("[BEACON-ERR] (%s):[%s] We've found beacon for this packet, but they are not in proper mode.", bc.CurrentState().String(), bc.SlaveNode().SlaveID)
                    return nil
                }
                default: {
                    return bc.TransitionWithSlaveMeta(&beaconD.Address, usm, ts)
                }
            }
        }
    }

    return nil
}
*/

func (r *registerManager) GuideNodeRegistrationWithBeacon(beaconD ucast.BeaconPack, ts time.Time) error {
    // (Ignore) there are way too many unpackable packages.
    usm, err := slagent.UnpackedSlaveMeta(beaconD.Message)
    if err != nil {
        return nil
    }
    // this packet looks for something else
    if len(usm.MasterBoundAgent) != 0 && usm.MasterBoundAgent != r.beaconManger.clusterID {
        return nil
    }

    log.Debugf("[REGISTER-RX] %v\n%v", beaconD.Address.IP.String(), spew.Sdump(usm))

    // remove discarded beacon
    pruneBeaconList(r.beaconManger)

    // check if beacon for this packet exists
    var bLen int = len(r.beaconManger.beaconList)
    for i := 0; i < bLen; i++  {
        bc := r.beaconManger.beaconList[i]
        if bc.SlaveNode().SlaveID == usm.SlaveID {
            switch bc.CurrentState() {
                case MasterDiscarded:
                    fallthrough
                // should be in registration
                case MasterInit:
                    fallthrough
                // should be in recovery
                case MasterBindBroken:
                    fallthrough
                // should be in bind
                case MasterBindRecovery: {
                    return errors.Errorf("[REGISTER-RX] Node (%v|%v|%v) in illegal state", usm.SlaveID, bc.CurrentState().String(), beaconD.Address.IP.String())
                }

                // we need to monitor this to make sure the node we try to bind has been bound successfully
                case MasterBounded: {
                    log.Debugf("[REGISTER-RX] Node (%v|%v|%v) check if bound ok", usm.SlaveID, bc.CurrentState().String(), beaconD.Address.IP.String())
                }

                default: {
                    return bc.TransitionWithSlaveMeta(&beaconD.Address, usm, ts)
                }
            }
        }
    }

    return errors.Errorf("[REGISTER-RX] Node(%v|%v) unregistered node with same cluster id. *should never happen*", usm.SlaveID, beaconD.Address.IP.String())
}