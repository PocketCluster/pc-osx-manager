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

func NewBeaconManagerWithFunc(comm CommChannelFunc) BeaconManger {
    return NewBeaconManager(comm)
}

func NewBeaconManager(comm CommChannel) BeaconManger {
    // TODO find existing slave models and attach them
    return &beaconManger {
        commChannel:    comm,
        beaconList:     []MasterBeacon{},
    }
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
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(beaconD.Message)
    if err != nil {
        return errors.WithStack(err)
    }

    log.Debugf("[BEACON] %v", spew.Sdump(usm))
    model.FindSlaveNode("slave_id = ?", usm.SlaveID)
    return nil
}

func (b *beaconManger) TransitionWithSearchData(searchD mcast.CastPack) error {
    var (
        bcFound bool = false
        mc MasterBeacon = nil
        ts time.Time = time.Now()
    )

    usm, err := slagent.UnpackedSlaveMeta(searchD.Message)
    if err != nil {
        return errors.WithStack(err)
    }

    log.Debugf("[SEARCH] FROM %v\n%v ", searchD.Address.IP.String(), spew.Sdump(usm))

    for _, bc := range b.beaconList {
        if bc.SlaveNode().MacAddress == usm.SlaveID {

            // this beacon's state should be identified
            err = bc.TransitionWithSlaveMeta(nil, usm, ts)

            bcFound = true
            break
        }
    }
    if !bcFound {
        mc, err = NewMasterBeacon(MasterInit, nil, b.commChannel)
        if err != nil {
            return errors.WithStack(err)
        }
        b.beaconList = append(b.beaconList, mc)
        err = mc.TransitionWithSlaveMeta(nil, usm, ts)
        if err != nil {
            return errors.WithStack(err)
        }
    }
    return nil
}

func (b *beaconManger) TransitionWithTimestamp(ts time.Time) error {

    return nil
}

func (b *beaconManger) Shutdown() error {
    b.commChannel = nil
    return nil
}
