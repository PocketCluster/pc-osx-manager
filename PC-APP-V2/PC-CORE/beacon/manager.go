package beacon

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-node-agent/slagent"
)
import (
    "github.com/davecgh/go-spew/spew"
)

func NewBeaconManagerWithFunc(comm CommChannelFunc) (BeaconManger, error) {
    return NewBeaconManager(comm)
}

func NewBeaconManager(comm CommChannel) (BeaconManger, error) {
    var (
        beacons []MasterBeacon = []MasterBeacon{}
        nodes []model.SlaveNode = nil
        mb MasterBeacon = nil
        err error = nil
    )
    nodes, err = model.FindAllSlaveNode()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    for _, n := range nodes {
        mb, err = NewMasterBeacon(MasterBindBroken, &n, comm)
        if err != nil {
            return nil, errors.WithStack(err)
        }
        beacons = append(beacons, mb)
    }
    return &beaconManger {
        commChannel:    comm,
        beaconList:     beacons,
    }, nil
}

type BeaconManger interface {
    TransitionWithBeaconData(beaconD ucast.BeaconPack) error
    TransitionWithSearchData(searchD mcast.CastPack) error
    TransitionWithTimestamp(ts time.Time) error
    Shutdown() error
}

// We might not need a locking mechanism as "select" statement will choose only "one input" at a time.
type beaconManger struct {
    commChannel  CommChannel
    beaconList   []MasterBeacon
}

func (b *beaconManger) TransitionWithBeaconData(beaconD ucast.BeaconPack) error {
    var (
        err error               = nil
        usm *slagent.PocketSlaveAgentMeta = nil
        ts time.Time            = time.Now()
        activeBC []MasterBeacon = []MasterBeacon{}
    )

    // suppose we've sort out what this is.
    usm, err = slagent.UnpackedSlaveMeta(beaconD.Message)
    if err != nil {
        return errors.WithStack(err)
    }

    log.Debugf("[BEACON] %v", spew.Sdump(usm))

    // this packet looks for something else
    if len(usm.DiscoveryAgent.MasterBoundAgent) != 0 && usm.DiscoveryAgent.MasterBoundAgent != "current master" {
        return nil
    }

    // remove discarded beacon
    for _, bc := range b.beaconList {
        if bc.CurrentState() != MasterDiscarded {
            activeBC = append(activeBC, bc)
        }
    }
    b.beaconList = activeBC

    // check if beacon for this packet exists
    for _, bc := range b.beaconList {
        if bc.SlaveNode().MacAddress == usm.SlaveID {
            switch bc.CurrentState() {
                case MasterInit:
                    fallthrough
                case MasterBindBroken:
                    fallthrough
                case MasterDiscarded: {
                    log.Debugf("We've found beacon for this packet, but they are not in proper mode. Let's stop")
                    return nil
                }
                default: {
                    return bc.TransitionWithSlaveMeta(&beaconD.Address, usm, ts)
                }
            }
        }
    }

    model.FindSlaveNode("slave_id = ?", usm.SlaveID)
    return nil
}

func (b *beaconManger) TransitionWithSearchData(searchD mcast.CastPack) error {
    var (
        bcFound bool            = false
        mc MasterBeacon         = nil
        err error               = nil
        usm *slagent.PocketSlaveAgentMeta = nil
        state MasterBeaconState = MasterBounded
        ts time.Time            = time.Now()
        activeBC []MasterBeacon = []MasterBeacon{}
    )

    usm, err = slagent.UnpackedSlaveMeta(searchD.Message)
    if err != nil {
        return errors.WithStack(err)
    }

    log.Debugf("[SEARCH] FROM %v\n%v ", searchD.Address.IP.String(), spew.Sdump(usm))

    // this packet looks for something else
    if len(usm.DiscoveryAgent.MasterBoundAgent) != 0 && usm.DiscoveryAgent.MasterBoundAgent != "current master" {
        return nil
    }

    // remove discarded beacon
    for _, bc := range b.beaconList {
        if bc.CurrentState() != MasterDiscarded {
            activeBC = append(activeBC, bc)
        }
    }
    b.beaconList = activeBC

    // check if beacon for this packet exists
    for _, bc := range b.beaconList {
        if bc.SlaveNode().MacAddress == usm.SlaveID {

            // this beacons are created and waiting for an input
            state = bc.CurrentState()
            if state == MasterInit || state == MasterBindBroken {
                return bc.TransitionWithSlaveMeta(nil, usm, ts)
            }

            // if beacon is not in searching state, then mark and we've found target
            bcFound = true
            break
        }
    }

    // since we've not found, create new beacon
    if !bcFound {
        mc, err = NewMasterBeacon(MasterInit, nil, b.commChannel)
        if err != nil {
            return errors.WithStack(err)
        }
        b.beaconList = append(b.beaconList, mc)
        return mc.TransitionWithSlaveMeta(nil, usm, ts)
    }

    return errors.Errorf("[ERR] TransitionWithSearchData reaches at the end. *this should never happen*")
}

func (b *beaconManger) TransitionWithTimestamp(ts time.Time) error {
    var (
        err error               = nil
        activeBC []MasterBeacon = []MasterBeacon{}
    )

    // check if beacon for this packet exists
    for _, bc := range b.beaconList {
        err = bc.TransitionWithTimestamp(ts)
        if err != nil {
            log.Debugf(err.Error())
        }
    }

    // remove discarded beacon
    for _, bc := range b.beaconList {
        if bc.CurrentState() != MasterDiscarded {
            activeBC = append(activeBC, bc)
        }
    }
    b.beaconList = activeBC

    return nil
}

func (b *beaconManger) Shutdown() error {
    b.commChannel = nil
    return nil
}
