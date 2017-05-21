package beacon

import (
    "fmt"
    "sync"
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

func NewBeaconManagerWithFunc(cid string, comm CommChannelFunc) (BeaconManger, error) {
    return NewBeaconManager(cid, comm)
}

func NewBeaconManager(cid string, comm CommChannel) (BeaconManger, error) {
    var (
        beacons []MasterBeacon = []MasterBeacon{}
        nodes []model.SlaveNode = nil
        mb MasterBeacon = nil
        err error = nil
    )

    nodes, err = model.FindAllSlaveNode()
    if err != nil && err != model.NoItemFound{
        return nil, errors.WithStack(err)
    }

    var nLen int = len(nodes)
    for i := 0; i < nLen; i ++ {
        n := &(nodes[i])
        switch n.State {
            case model.SNMStateJoined: {
                mb, err = NewMasterBeacon(MasterBindBroken, n, comm)
                if err != nil {
                    return nil, errors.WithStack(err)
                }
            }
        }
        beacons = append(beacons, mb)
    }

    return &beaconManger {
        clusterID:      cid,
        commChannel:    comm,
        beaconList:     beacons,
    }, nil
}

type BeaconManger interface {
    TransitionWithBeaconData(beaconD ucast.BeaconPack, ts time.Time) error
    TransitionWithSearchData(searchD mcast.CastPack, ts time.Time) error
    TransitionWithTimestamp(ts time.Time) error
    Shutdown() error
}

// We might not need a locking mechanism as "select" statement will choose only "one input" at a time.
type beaconManger struct {
    sync.Mutex
    clusterID      string
    commChannel    CommChannel
    beaconList     []MasterBeacon
}

func (b *beaconManger) Sanitize(s *model.SlaveNode) error {
    assignSlaveNodeName(b, s)
    return nil
}

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

    log.Debugf("[BEACON] FROM %v\n%v", beaconD.Address.IP.String(), spew.Sdump(usm))

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

    return nil
}

func (b *beaconManger) TransitionWithSearchData(searchD mcast.CastPack, ts time.Time) error {
    var (
        bcFound bool            = false
        mc MasterBeacon         = nil
        err error               = nil
        usm *slagent.PocketSlaveAgentMeta = nil
        state MasterBeaconState = MasterBounded
    )

    usm, err = slagent.UnpackedSlaveMeta(searchD.Message)
    if err != nil {
        return errors.WithStack(err)
    }

    log.Debugf("[SEARCH] FROM %v\n%v ", searchD.Address.IP.String(), spew.Sdump(usm))

    // this packet looks for something else
    if len(usm.MasterBoundAgent) != 0 && usm.MasterBoundAgent != b.clusterID {
        log.Debugf("[SEARCH] this packet belong to other master | usm.DiscoveryAgent.MasterBoundAgent %v | b.clusterID %v", usm.MasterBoundAgent, b.clusterID)
        return nil
    }

    // remove discarded beacon
    pruneBeaconList(b)

    // check if beacon for this packet exists
    var bLen int = len(b.beaconList)
    for i := 0; i < bLen; i++  {
        bc := b.beaconList[i]
        if bc.SlaveNode().MacAddress == usm.SlaveID {

            // this beacons are created and waiting for an input
            state = bc.CurrentState()
            if state == MasterInit || state == MasterBindBroken {
                return bc.TransitionWithSlaveMeta(&searchD.Address, usm, ts)
            }

            // if beacon is not in searching state, then mark and we've found target
            bcFound = true
            break
        }
    }

    // since we've not found, create new beacon
    if !bcFound {
        mc, err = NewMasterBeacon(MasterInit, model.NewSlaveNode(b), b.commChannel)
        if err != nil {
            return errors.WithStack(err)
        }
        insertMasterBeacon(b, mc)
        return mc.TransitionWithSlaveMeta(&searchD.Address, usm, ts)
    }

    return errors.Errorf("[ERR] TransitionWithSearchData reaches at the end. *this should never happen, and might be a malicious attempt*")
}

func (b *beaconManger) TransitionWithTimestamp(ts time.Time) error {
    var err error = nil

    // remove discarded beacon
    pruneBeaconList(b)

    // check if beacon for this packet exists
    var bLen int = len(b.beaconList)
    for i := 0; i < bLen; i++  {
        bc := b.beaconList[i]
        err = bc.TransitionWithTimestamp(ts)
        if err != nil {
            log.Debugf(err.Error())
        }
    }

    return nil
}

func (b *beaconManger) Shutdown() error {
    shutdownMasterBeacons(b)
    b.commChannel = nil
    return nil
}

func pruneBeaconList(b *beaconManger) {
    b.Lock()
    defer b.Unlock()

    var (
        activeBC []MasterBeacon = []MasterBeacon{}
        bLen int = len(b.beaconList)
    )

    for i := 0; i < bLen; i++ {
        bc := b.beaconList[i]
        if bc.CurrentState() == MasterDiscarded {
            bc.Shutdown()
        } else {
            activeBC = append(activeBC, bc)
        }
    }
    b.beaconList = activeBC
}

func insertMasterBeacon(b *beaconManger, m MasterBeacon) {
    b.Lock()
    defer b.Unlock()

    b.beaconList = append(b.beaconList, m)
}

func assignSlaveNodeName(b *beaconManger, s *model.SlaveNode) {
    b.Lock()
    defer b.Unlock()

    var (
        ci int = 0
        cname string = ""
        findName = func(mbl []MasterBeacon, nUUID, nName string) bool {
            var bLen = len(mbl)
            for i := 0; i < bLen; i++ {
                mb := mbl[i]
                if mb.SlaveNode().SlaveUUID == nUUID {
                    continue
                }
                switch mb.CurrentState() {
                case MasterDiscarded:
                    continue
                default:
                    if mb.SlaveNode().NodeName == cname {
                        return true
                    }
                }
            }
            return false
        }
    )

    for {
        cname = fmt.Sprintf("pc-node%d", ci + 1)
        if !findName(b.beaconList, s.SlaveUUID, cname) {
            s.NodeName = cname
            return
        }
        ci++
    }
}

func shutdownMasterBeacons(b *beaconManger) {
    b.Lock()
    defer b.Unlock()

    var (
        bLen = len(b.beaconList)
        mb MasterBeacon = nil
    )

    for i := 0; i < bLen; i++ {
        mb = b.beaconList[i]
        mb.Shutdown()
    }
    // assign new slice to prevent nil crash
    b.beaconList = []MasterBeacon{}
}
