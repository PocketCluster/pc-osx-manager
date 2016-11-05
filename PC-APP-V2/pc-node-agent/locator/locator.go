package locator

import (
    "time"
    "github.com/stkim1/pc-core/msagent"
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

type SlaveLocator interface {
    CurrentState() SlaveLocatingState
    TranstionWithTimestamp(timestamp time.Time) error
    TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) error
    Close() error
}

// On sucess happens at the moment successful state transition takes place
type SlaveLocatorOnStateTransitionSuccess      func (SlaveLocatingState)

// OnIdle happens as locator awaits
type SlaveLocatorOnStateTransitionIdle         func (SlaveLocatingState, time.Time, time.Time, int) bool

// OnFail happens at the moment state transition fails to happen
type SlaveLocatorOnStateTransitionFailure      func (SlaveLocatingState)

// ------ DEFAULT TIMEOUTS ------
const (
    UNBOUNDED_TIMEOUT   = 3 * time.Second
    BOUNDED_TIMEOUT     = 10 * time.Second
)
