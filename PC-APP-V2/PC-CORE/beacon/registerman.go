package beacon

import (
    "fmt"
    "sync"
    "time"

    "github.com/docker/docker/pkg/discovery"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/davecgh/go-spew/spew"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/pc-vbox-comm/masterctrl"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    "github.com/stkim1/pc-node-agent/slagent"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/service/ivent"
)

type RegisterManger interface {
    MonitoringMasterSearchData(searchD mcast.CastPack, ts time.Time) error
}

type registerManager struct {
    *beaconManger
}

func NewNodeRegisterManager(master BeaconManger) (RegisterManger, error) {
    bm, ok := master.(*beaconManger)
    if !ok {
        return nil, errors.Errorf("invalid meanager type")
    }
    return &registerManager{bm}, nil
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

    log.Debugf("[SEARCH-RX] %v\n%v ", searchD.Address.IP.String(), spew.Sdump(usm))

    // this doesn't belong the master
    if len(usm.MasterBoundAgent) != 0 && usm.MasterBoundAgent != r.clusterID {
        log.Debugf("[SEARCH-RX] this packet belong to other master | usm.DiscoveryAgent.MasterBoundAgent %v | b.clusterID %v", usm.MasterBoundAgent, b.clusterID)
        return nil
    }
    return nil
}

func (r *registerManager) RegisterNewSlaveNode() error {

    return nil
}