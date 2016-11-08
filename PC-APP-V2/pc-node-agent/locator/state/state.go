package state

import (
    "time"

    "github.com/stkim1/pc-core/msagent"
)

const (
    TranstionFailureLimit uint          = 5
    TransitionTimeout     time.Duration = time.Second * 10

    TxActionLimit         uint          = 5
    UnboundedTimeout      time.Duration = time.Second * 3
    BoundedTimeout        time.Duration = time.Second * 10
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

type LocatorState interface {
    CurrentState() SlaveLocatingState
    MasterMetaTranstion(timestamp time.Time) (SlaveLocatingTransition, error)
    TimestampTranstion(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (SlaveLocatingTransition, error)
}

type locatorState struct {
    // each time we try to make transtion and fail, count goes up.
    transitionActionCount uint

    // last time successfully transitioned state.
    lastTransitionTS      time.Time

    // each time we try to send something, count goes up. This include success/fail altogether.
    txActionCount         uint

    // last time transmission takes place. This is to control the frequnecy of transmission
    lastTxTS              time.Time

    locatingState         SlaveLocatingState
}

// property functions
func (ls *locatorState) CurrentState() SlaveLocatingState {
    return ls.locatingState
}

func (ls *locatorState) transtionFailureLimit() uint {
    return TranstionFailureLimit
}

func (ls *locatorState) transitionTimeout() time.Duration {
    return TransitionTimeout
}

func (ls *locatorState) txActionLimit() uint {
    return TxActionLimit
}

func (ls *locatorState) txTimeout() time.Duration {
    return UnboundedTimeout
}

// -- STATE TRANSITION
func stateTransition(currState SlaveLocatingState, nextCondition SlaveLocatingTransition) SlaveLocatingState {
    var nextState SlaveLocatingState
    // Succeed to transition to the next
    if  nextCondition == SlaveTransitionOk {
        switch currState {
            case SlaveUnbounded:
                nextState = SlaveInquired
            case SlaveInquired:
                nextState = SlaveKeyExchange
            case SlaveKeyExchange:
                nextState = SlaveCryptoCheck

            case SlaveCryptoCheck:
                fallthrough
            case SlaveBindBroken:
                fallthrough
            case SlaveBounded:
                nextState = SlaveBounded
                break

            default:
                nextState = SlaveUnbounded
        }
        // Fail to transition to the next
    } else if nextCondition == SlaveTransitionFail {
        switch currState {
            case SlaveUnbounded:
                fallthrough
            case SlaveInquired:
                fallthrough
            case SlaveKeyExchange:
                fallthrough
            case SlaveCryptoCheck:
                nextState = SlaveUnbounded
                break

            case SlaveBindBroken:
                fallthrough
            case SlaveBounded:
                nextState = SlaveBindBroken
                break

            default:
                nextState = SlaveUnbounded
        }
        // Idle
    } else {
        nextState = currState
    }
    return nextState
}

// --- STATE TRANSITION ---
func (ls *locatorState) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    return SlaveTransitionIdle, nil
}

func finalizeTransitionWithTimeout(ls *locatorState, nextStateCandiate SlaveLocatingTransition, slaveTimestamp time.Time) SlaveLocatingTransition {
    var nextConfirmedState SlaveLocatingTransition
    switch nextStateCandiate {
        case SlaveTransitionOk: {
            // reset transition action count / timestamp to 0
            ls.lastTransitionTS = slaveTimestamp
            ls.transitionActionCount = 0

            // since
            ls.lastTxTS = slaveTimestamp
            ls.txActionCount = 0
            nextConfirmedState = SlaveTransitionOk
        }
        default: {
            if ls.transitionActionCount < ls.transtionFailureLimit() {
                ls.transitionActionCount++
            }

            if ls.transitionActionCount < ls.transtionFailureLimit() && slaveTimestamp.Sub(ls.lastTransitionTS) < ls.transitionTimeout() {
                nextConfirmedState = SlaveTransitionIdle
            } else {
                nextConfirmedState = SlaveTransitionFail
            }
        }
    }
    return nextConfirmedState
}

func (ls *locatorState) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return nil
}

func (ls *locatorState) onStateTranstionFailure(slaveTimestamp time.Time) error {
    return nil
}

func executeOnTransitionEvents(ls *locatorState, newState, oldState SlaveLocatingState, transition SlaveLocatingTransition, slaveTimestamp time.Time) error {
    if newState != oldState {
        switch transition {
            case SlaveTransitionOk:
                return ls.onStateTranstionSuccess(slaveTimestamp)
            case SlaveTransitionFail: {
                return ls.onStateTranstionFailure(slaveTimestamp)
            }
        }
    }
    return nil
}

func (ls *locatorState) MasterMetaTranstion(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    var (
        transition SlaveLocatingTransition
        transErr error = nil
        eventErr error = nil
    )

    transition, transErr = ls.transitionWithMasterMeta(meta, slaveTimestamp)

    // filter out the intermediate transition value with failed count + timestamp
    finalTransitionCandidate := finalizeTransitionWithTimeout(ls, transition, slaveTimestamp)

    oldState := ls.locatingState
    // finalize locating master beacon state
    newState := stateTransition(ls.locatingState, finalTransitionCandidate)
    // fianalize state change
    ls.locatingState = newState

    // execute event lisenter
    eventErr = executeOnTransitionEvents(ls, newState, oldState, finalTransitionCandidate, slaveTimestamp)

    return finalTransitionCandidate, summarizeErrors(transErr, eventErr)
}

// --- TRANSMISSION CONTROL
func (ls *locatorState) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    return nil
}

func checkTxStateWithTime(ls *locatorState, slaveTimestamp time.Time) SlaveLocatingTransition {
    if ls.txActionCount < ls.txActionLimit() {
        // if tx timeout window is smaller than time delta (T_1 - T_0), don't do anything!!! just skip!
        if ls.txTimeout() < slaveTimestamp.Sub(ls.lastTransitionTS) {
            ls.txActionCount++
        }
        return SlaveTransitionIdle
    }
    // this is failure. the fact that this is called indicate that we're ready to move to failure state
    return SlaveTransitionFail
}

func (ls *locatorState) TimestampTranstion(slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    var (
        transition SlaveLocatingTransition
        transErr error = nil
        eventErr error = nil
    )

    transition = checkTxStateWithTime(ls, slaveTimestamp)

    if transition == SlaveTransitionIdle {
        transErr = ls.transitionActionWithTimestamp(slaveTimestamp)

        // since an action is taken, the action counter goes up regardless of error
        ls.txActionCount++
    } else {
        // now idle action condition has failed, and we need to make transition to FAILTURE state

        oldState := ls.locatingState
        // finalize locating master beacon state
        newState := stateTransition(ls.locatingState, transition)
        // fianalize state change
        ls.locatingState = newState

        // execute event lisenter
        eventErr = executeOnTransitionEvents(ls, newState, oldState, transition, slaveTimestamp)
    }

    return transition, summarizeErrors(transErr, eventErr)
}
