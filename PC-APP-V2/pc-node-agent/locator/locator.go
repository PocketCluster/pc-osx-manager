package locator

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-core/msagent"
    "bytes"
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

type CommChannel interface {
    McastSend(data []byte) error
    UcastSend(data []byte, target string) error
}

type SlaveLocator interface {
    CurrentState() (SlaveLocatingState, error)
    TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) error
    TranstionWithTimestamp(timestamp time.Time) error
    Close() error
}

type slaveLocator struct {
    state       LocatorState
}

// New slaveLocator starts only from unbounded or bindbroken
func NewSlaveLocator(state SlaveLocatingState, comm CommChannel) (SlaveLocator, error) {
    if comm == nil {
        return nil, fmt.Errorf("[ERR] communication channel cannot be void")
    }

    switch state {
    case SlaveUnbounded:
        return &slaveLocator{state: newUnboundedState(comm)}, nil
    case SlaveBindBroken:
        return &slaveLocator{state: newBindbrokenState(comm)}, nil
    default:
        return nil, fmt.Errorf("[ERR] SlaveLocator can initiated from SlaveUnbounded or SlaveBindBroken only")
    }
}

func (sl *slaveLocator) CurrentState() (SlaveLocatingState, error) {
    if sl.state == nil {
        return SlaveUnbounded, nil
    }
    return sl.state.CurrentState(), nil
}

func (sl *slaveLocator) TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) error {
    if sl.state == nil {
        return fmt.Errorf("[ERR] LocatorState is nil. Cannot make transition with master meta")
    }
    var err error
    sl.state, err = sl.state.MasterMetaTransition(meta, slaveTimestamp)
    return err
}

func (sl *slaveLocator) TranstionWithTimestamp(slaveTimestamp time.Time) error {
    if sl.state == nil {
        return fmt.Errorf("[ERR] LocatorState is nil. Cannot make transition with master meta")
    }
    var err error
    sl.state, err = sl.state.TimestampTransition(slaveTimestamp)
    return err
}

func (sl *slaveLocator) Close() error {
    // TODO : TO BE CONTINUED
    return nil
}

type opError struct {
    TransitionError         error
    EventError              error
}

func (oe *opError) Error() string {
    var errStr bytes.Buffer

    if oe.TransitionError != nil {
        errStr.WriteString(oe.TransitionError.Error())
    }

    if oe.EventError != nil {
        errStr.WriteString(oe.EventError.Error())
    }
    return errStr.String()
}

func summarizeErrors(transErr error, eventErr error) error {
    if transErr == nil && eventErr == nil {
        return nil
    }
    return &opError{TransitionError: transErr, EventError: eventErr}
}
