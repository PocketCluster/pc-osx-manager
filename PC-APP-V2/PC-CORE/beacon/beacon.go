package beacon

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/model"
    "net"
)

type MasterBeaconState int
const (
    MasterInit              MasterBeaconState = iota
    MasterUnbounded
    MasterInquired
    MasterKeyExchange
    MasterCryptoCheck
    MasterBounded
    MasterBindRecovery
    MasterBindBroken
    MasterDiscarded
)

type MasterBeaconTransition int
const (
    MasterTransitionFail    MasterBeaconTransition = iota
    MasterTransitionOk
    MasterTransitionIdle
)

func (st MasterBeaconState) String() string {
    var state string
    switch st {
        case MasterInit:
            state = "MasterInit"
        case MasterUnbounded:
            state = "MasterUnbounded"
        case MasterInquired:
            state = "MasterInquired"
        case MasterKeyExchange:
            state = "MasterKeyExchange"
        case MasterCryptoCheck:
            state = "MasterCryptoCheck"
        case MasterBounded:
            state = "MasterBounded"
        case MasterBindRecovery:
            state = "MasterRecovery"
        case MasterBindBroken:
            state = "MasterBindBroken"
        case MasterDiscarded:
            state = "MasterDiscarded"
    }
    return state
}

type CommChannel interface {
    //McastSend(data []byte) error
    UcastSend(target string, data []byte) error
}

type CommChannelFunc func(target string, data []byte) error
func (c CommChannelFunc) UcastSend(target string, data []byte) error {
    return c(target, data)
}

// MasterBeacon is assigned individually for each slave node.
type MasterBeacon interface {
    CurrentState() MasterBeaconState
    TransitionWithTimestamp(timestamp time.Time) error
    TransitionWithSlaveMeta(sender *net.UDPAddr, meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error
    Shutdown()
    SlaveNode() *model.SlaveNode
}

func NewMasterBeaconWithFunc(state MasterBeaconState, slaveNode *model.SlaveNode, comm CommChannelFunc) (MasterBeacon, error) {
    return NewMasterBeacon(state, slaveNode, comm)
}

func NewMasterBeacon(state MasterBeaconState, slaveNode *model.SlaveNode, comm CommChannel) (MasterBeacon, error) {
    if slaveNode == nil {
        return nil, errors.Errorf("[ERR] slavenode cannot be nil")
    }
    if comm == nil {
        return nil, errors.Errorf("[ERR] communication channel cannot be nil")
    }

    switch state {
        case MasterInit:
            return &masterBeacon{state:beaconinitState(slaveNode, comm)}, nil

        case MasterBindBroken:
            if slaveNode == nil {
                return nil, errors.Errorf("[ERR] Slavenode cannot be nil")
            }
            bstate, err := bindbrokenState(slaveNode, comm)
            if err != nil {
                return nil, errors.WithStack(err)
            }
            return &masterBeacon{state:bstate}, nil
    }
    return nil, errors.Errorf("[ERR] MasterBeacon can initiated from MasterInit or MasterBindBroken only")
}

type masterBeacon struct {
    state       BeaconState
}

func (mb *masterBeacon) CurrentState() MasterBeaconState {
    return mb.state.CurrentState()
}

func (mb *masterBeacon) TransitionWithTimestamp(timestamp time.Time) error {
    if mb.state == nil {
        return errors.Errorf("[ERR] BeaconState is nil. Cannot make transition with master timestamp")
    }
    var err error = nil
    mb.state, err = mb.state.TransitionWithTimestamp(timestamp)
    return err
}

func (mb *masterBeacon) TransitionWithSlaveMeta(sender *net.UDPAddr, meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if mb.state == nil {
        return errors.Errorf("[ERR] BeaconState is nil. Cannot make transition with master meta")
    }
    var err error = nil
    mb.state, err = mb.state.TransitionWithSlaveMeta(sender, meta, timestamp)
    return errors.WithStack(err)
}

func (mb *masterBeacon) SlaveNode() *model.SlaveNode {
    return mb.state.SlaveNode()
}

func (mb *masterBeacon) Shutdown() {
    mb.state.Close()
}
