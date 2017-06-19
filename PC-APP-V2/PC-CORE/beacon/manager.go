package beacon

import (
    "fmt"
    "sync"
    "time"

    "github.com/docker/docker/pkg/discovery"
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

type BeaconEventNotification interface {
    BeaconEventPrepareJoin(slave *model.SlaveNode) error
    BeaconEventResurrect(slaves []model.SlaveNode) error
    BeaconEventTranstion(state MasterBeaconState, slave *model.SlaveNode, ts time.Time, transOk bool) error
    BeaconEventDiscard(slave *model.SlaveNode) error
    BeaconEventShutdown() error
}

func NewBeaconManagerWithFunc(cid string, noti BeaconEventNotification, comm CommChannelFunc) (BeaconManger, error) {
    return NewBeaconManager(cid, noti, comm)
}

func NewBeaconManager(cid string, noti BeaconEventNotification, comm CommChannel) (BeaconManger, error) {
    var (
        beacons []MasterBeacon = []MasterBeacon{}
        bm *beaconManger = nil
        nodes []model.SlaveNode = nil
        mb MasterBeacon = nil
        err error = nil
    )
    if comm == nil {
        return nil, errors.Errorf("[ERR] comm channel cannot be nil")
    }
    if noti == nil {
        return nil, errors.Errorf("[ERR] notification receiver cannot be nil")
    }

    bm = &beaconManger {
        clusterID:       cid,
        notiReceiver:    noti,
        commChannel:     comm,
    }

    // respawn nodes
    nodes, err = model.FindAllSlaveNode()
    if err != nil && err != model.NoItemFound{
        return nil, errors.WithStack(err)
    }

    var nLen int = len(nodes)
    for i := 0; i < nLen; i ++ {
        n := &(nodes[i])
        switch n.State {
            case model.SNMStateJoined: {
                mb, err = NewMasterBeacon(MasterBindBroken, n, comm, bm)
                if err != nil {
                    return nil, errors.WithStack(err)
                }
            }
        }
        beacons = append(beacons, mb)
    }
    bm.beaconList = beacons

    // we'll ignore error message for now
    noti.BeaconEventResurrect(nodes)

    return bm, nil
}

type BeaconManger interface {
    TransitionWithBeaconData(beaconD ucast.BeaconPack, ts time.Time) error
    TransitionWithSearchData(searchD mcast.CastPack, ts time.Time) error
    TransitionWithTimestamp(ts time.Time) error
    Shutdown() error
    // swarm discovery backend
    discovery.Backend
}

// We might not need a locking mechanism as "select" statement will choose only "one input" at a time.
type beaconManger struct {
    sync.Mutex
    clusterID         string
    notiReceiver      BeaconEventNotification
    commChannel       CommChannel
    beaconList        []MasterBeacon
    // swarm heartbeat
    swarmHeartbeat    time.Duration
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

    log.Debugf("[SEARCH-RX] %v\n%v ", searchD.Address.IP.String(), spew.Sdump(usm))

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

// --- Swarm Discovery Methods --- //

// Initialize the discovery with URIs, a heartbeat, a ttl and optional settings.
func (b *beaconManger) Initialize(_ string, hb time.Duration, _ time.Duration, _ map[string]string) error {
    b.swarmHeartbeat = hb
    return nil
}

// Watch the discovery for entry changes.
// Returns a channel that will receive changes or an error.
// Providing a non-nil stopCh can be used to stop watching.
func (b *beaconManger) Watch(stopCh <-chan struct{}) (<-chan discovery.Entries, <-chan error) {
    var (
        ch = make(chan discovery.Entries)
        errCh = make(chan error)
        ticker = time.NewTicker(b.swarmHeartbeat)
    )

    go func(bm *beaconManger) {
        defer close(errCh)
        defer close(ch)

        // Send the initial entries if available.
        var (
            values []string = findBoundedNodesForSwarm(bm)
            currentEntries, newEntries discovery.Entries
            err error = nil
        )

        if len(values) > 0 {
            currentEntries, err = discovery.CreateEntries(values)
        }
        if err != nil {
            errCh <- err
        } else if currentEntries != nil {
            ch <- currentEntries
        }

        // Periodically send updates.
        for {
            select {
                case <-ticker.C: {
                    values = findBoundedNodesForSwarm(bm)
                    newEntries, err = discovery.CreateEntries(values)
                    if err != nil {
                        errCh <- err
                        continue
                    }
                    // Check if the file has really changed.
                    if !newEntries.Equals(currentEntries) {
                        ch <- newEntries
                    }
                    currentEntries = newEntries
                }

                case <-stopCh: {
                    ticker.Stop()
                    return
                }
            }
        }
    }(b)

    return ch, errCh
}

// Register to the discovery.
func (b *beaconManger) Register(string) error {
    log.Debugf("(INFO) this should not work as registration of new address is not done by swarm")
    return nil
}

// --- BeaconOnTransitionEvent methods --- //

// state transition success from
func (b *beaconManger) OnStateTranstionSuccess(state MasterBeaconState, slave *model.SlaveNode, ts time.Time) error {
    // we'll ignore error message for now
    b.notiReceiver.BeaconEventTranstion(state, slave, ts, true)
    return nil
}

// state transition failure from
func (b *beaconManger) OnStateTranstionFailure(state MasterBeaconState, slave *model.SlaveNode, ts time.Time) error {
    // we'll ignore error message for now
    b.notiReceiver.BeaconEventTranstion(state, slave, ts, false)
    return nil
}

// --- private static methods --- //

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
            b.notiReceiver.BeaconEventDiscard(bc.SlaveNode())
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
        findName = func(mbl []MasterBeacon, authToken, nName string) bool {
            var bLen = len(mbl)
            for i := 0; i < bLen; i++ {
                mb := mbl[i]
                if mb.SlaveNode().AuthToken == authToken {
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
        if !findName(b.beaconList, s.AuthToken, cname) {
            s.NodeName = cname
            return
        }
        ci++
    }
}

func findBoundedNodesForSwarm(b *beaconManger) []string {
    b.Lock()
    defer b.Unlock()

    const (
        dockerPort string = "2375"
    )

    var (
        addr string = ""
        err error = nil
        nodeList []string = []string{}
        bLen int = len(b.beaconList)
    )

    for i := 0; i < bLen; i++ {
        bc := b.beaconList[i]
        if bc.CurrentState() == MasterBounded {
            addr, err = bc.SlaveNode().IP4AddrString()
            if err != nil {
                continue
            }
            nodeList = append(nodeList, fmt.Sprintf("%s:%s", addr, dockerPort))
        }
    }
    return nodeList
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
    b.commChannel = nil
    // ignore error for now
    b.notiReceiver.BeaconEventShutdown()
    b.notiReceiver = nil
}
