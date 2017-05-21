package locator

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/udpnet/ucast"
)
import (
    log "github.com/Sirupsen/logrus"
    "github.com/davecgh/go-spew/spew"
)

type SlaveLocatingState int
const (
    SlaveUnbounded          SlaveLocatingState = iota
    SlaveInquired
    SlaveKeyExchange
    SlaveCryptoCheck
    SlaveBounded
    SlaveBindBroken
)

type SlaveLocatingTransition int
const (
    SlaveTransitionFail     SlaveLocatingTransition = iota
    SlaveTransitionOk
    SlaveTransitionIdle
)

func (st SlaveLocatingState) String() string {
    var state string
    switch st {
    case SlaveUnbounded:
        state = "SlaveUnbounded"
    case SlaveBounded:
        state = "SlaveBounded"
    case SlaveBindBroken:
        state = "SlaveBindBroken"
    case SlaveInquired:
        state = "SlaveInquired"
    case SlaveKeyExchange:
        state = "SlaveKeyExchange"
    case SlaveCryptoCheck:
        state = "SlaveCryptoCheck"
    }
    return state
}

type SearchTx interface {
    McastSend(data []byte) error
}

type SearchTxFunc func(data []byte) error
func (s SearchTxFunc) McastSend(data []byte) error {
    return s(data)
}

type BeaconTx interface {
    UcastSend(target string, data []byte) error
}

type BeaconTxFunc func(target string, data []byte) error
func (b BeaconTxFunc) UcastSend(target string, data []byte) error {
    return b(target, data)
}

type SlaveLocator interface {
    CurrentState() (SlaveLocatingState, error)
    TranstionWithMasterBeacon(bp ucast.BeaconPack, slaveTimestamp time.Time) error
    TranstionWithTimestamp(timestamp time.Time) error
    Close() error

    // TODO : this should be deprecated for testing only
    TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) error
}

type slaveLocator struct {
    state       LocatorState
}

func NewSlaveLocatorWithFunc(state SlaveLocatingState, searchComm SearchTxFunc, beaconComm BeaconTxFunc) (SlaveLocator, error) {
    return NewSlaveLocator(state, searchComm, beaconComm)
}

// New slaveLocator starts only from unbounded or bindbroken
func NewSlaveLocator(state SlaveLocatingState, searchComm SearchTx, beaconComm BeaconTx) (SlaveLocator, error) {
    if searchComm == nil {
        return nil, errors.Errorf("[ERR] MasterSearch cannot be void")
    }
    if beaconComm == nil {
        return nil, errors.Errorf("[ERR] BeaconAgent cannot be void")
    }

    switch state {
        case SlaveUnbounded:
            return &slaveLocator{state: newUnboundedState(searchComm, beaconComm)}, nil
        case SlaveBindBroken:
            return &slaveLocator{state: newBindbrokenState(searchComm, beaconComm)}, nil
    }
    return nil, errors.Errorf("[ERR] SlaveLocator can initiated from SlaveUnbounded or SlaveBindBroken only")
}

func (sl *slaveLocator) CurrentState() (SlaveLocatingState, error) {
    if sl.state == nil {
        return SlaveUnbounded, nil
    }
    return sl.state.CurrentState(), nil
}

func (sl *slaveLocator) TranstionWithMasterBeacon(bp ucast.BeaconPack, slaveTimestamp time.Time) error {
    if sl.state == nil {
        return errors.Errorf("[ERR] LocatorState is nil. Cannot make transition with master meta")
    }
    // (2017-05-21) we're not looking into ucast.BeaconPack.Address for now as Master's interface address might vary
    meta, err := msagent.UnpackedMasterMeta(bp.Message)
    if err != nil {
        return errors.WithStack(err)
    }

    // TODO : should we check MasterBindAgent here?

    log.Debugf("[AGENT-BEACON] RECEIVED\n %v \n %v", spew.Sdump(bp.Address), spew.Sdump(meta))

    sl.state, err = sl.state.MasterMetaTransition(meta, slaveTimestamp)
    return err
}

func (sl *slaveLocator) TranstionWithTimestamp(slaveTimestamp time.Time) error {
    if sl.state == nil {
        return errors.Errorf("[ERR] LocatorState is nil. Cannot make transition with master meta")
    }
    var err error
    sl.state, err = sl.state.TimestampTransition(slaveTimestamp)
    return err
}

func (sl *slaveLocator) Close() error {
    // TODO : TO BE CONTINUED
    return nil
}

