package state

import (
    "time"

    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-core/msagent"
)

const (
    TranstionFailureLimit uint          = 5

    TransitionTimeout     time.Duration = time.Second * 10


    TxActionLimit         uint          = 5

    UnboundedTimeout      time.Duration = time.Second * 3

    BoundedTimeout        time.Duration = time.Second * 10
)

type opError struct {
    TransitionError         error
    EventError              error
}

func (oe *opError) Error() string {
    var errStr string = ""

    if oe.TransitionError != nil {
        errStr += oe.TransitionError.Error()
    }

    if oe.EventError != nil {
        errStr += oe.EventError.Error()
    }
    return errStr
}

func summarizeErrors(transErr error, eventErr error) *opError {
    if transErr == nil && eventErr == nil {
        return nil
    }
    return &opError{TransitionError: transErr, EventError: eventErr}
}

type LocatorState struct {
    // each time we try to make transtion and fail, count goes up.
    transitionActionCount uint

    // last time successfully transitioned state.
    lastTransitionTS      time.Time

    // each time we try to send something, count goes up. This include success/fail altogether.
    txActionCount         uint

    // last time transmission takes place. This is to control the frequnecy of transmission
    lastTxTS              time.Time

    locatingState         locator.SlaveLocatingState
}

// property functions
func (ls *LocatorState) LocatingState() locator.SlaveLocatingState {
    return ls.locatingState
}

func (ls *LocatorState) transtionFailureLimit() uint {
    return TranstionFailureLimit
}

func (ls *LocatorState) transitionTimeout() time.Duration {
    return TransitionTimeout
}

func (ls *LocatorState) txActionLimit() uint {
    return TxActionLimit
}

func (ls *LocatorState) txTimeout() time.Duration {
    return UnboundedTimeout
}


// -- STATE TRANSITION
func stateTransition(currState locator.SlaveLocatingState, nextCondition locator.SlaveLocatingTransition) locator.SlaveLocatingState {
    var nextState locator.SlaveLocatingState
    // Succeed to transition to the next
    if  nextCondition == locator.SlaveTransitionOk {
        switch currState {
            case locator.SlaveUnbounded:
                nextState = locator.SlaveInquired
            case locator.SlaveInquired:
                nextState = locator.SlaveKeyExchange
            case locator.SlaveKeyExchange:
                nextState = locator.SlaveCryptoCheck

            case locator.SlaveCryptoCheck:
                fallthrough
            case locator.SlaveBindBroken:
                fallthrough
            case locator.SlaveBounded:
                nextState = locator.SlaveBounded
                break

            default:
                nextState = locator.SlaveUnbounded
        }
        // Fail to transition to the next
    } else if nextCondition == locator.SlaveTransitionFail {
        switch currState {
            case locator.SlaveUnbounded:
                fallthrough
            case locator.SlaveInquired:
                fallthrough
            case locator.SlaveKeyExchange:
                fallthrough
            case locator.SlaveCryptoCheck:
                nextState = locator.SlaveUnbounded
                break

            case locator.SlaveBindBroken:
                fallthrough
            case locator.SlaveBounded:
                nextState = locator.SlaveBindBroken
                break

            default:
                nextState = locator.SlaveUnbounded
        }
        // Idle
    } else {
        nextState = currState
    }
    return nextState
}

// --- STATE TRANSITION ---
func (ls *LocatorState) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (locator.SlaveLocatingTransition, error) {
    return locator.SlaveTransitionIdle, nil
}

func finalizeTransitionWithTimeout(ls *LocatorState, nextStateCandiate locator.SlaveLocatingTransition, slaveTimestamp time.Time) locator.SlaveLocatingTransition {
    var nextConfirmedState locator.SlaveLocatingTransition
    switch nextStateCandiate {
        case locator.SlaveTransitionOk: {
            // reset transition action count / timestamp to 0
            ls.lastTransitionTS = slaveTimestamp
            ls.transitionActionCount = 0

            // since
            ls.lastTxTS = slaveTimestamp
            ls.txActionCount = 0
            nextConfirmedState = locator.SlaveTransitionOk
        }
        default: {
            if ls.transitionActionCount < ls.transtionFailureLimit() {
                ls.transitionActionCount++
            }

            if ls.transitionActionCount < ls.transtionFailureLimit() && slaveTimestamp.Sub(ls.lastTransitionTS) < ls.transitionTimeout() {
                nextConfirmedState = locator.SlaveTransitionIdle
            } else {
                nextConfirmedState = locator.SlaveTransitionFail
            }
        }
    }
    return nextConfirmedState
}

func (ls *LocatorState) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return nil
}

func (ls *LocatorState) onStateTranstionFailure(slaveTimestamp time.Time) error {
    return nil
}

func executeOnTransitionEvents(ls *LocatorState, newState, oldState locator.SlaveLocatingState, transition locator.SlaveLocatingTransition, slaveTimestamp time.Time) error {
    if newState != oldState {
        switch transition {
            case locator.SlaveTransitionOk:
                return ls.onStateTranstionSuccess(slaveTimestamp)
            case locator.SlaveTransitionFail: {
                return ls.onStateTranstionFailure(slaveTimestamp)
            }
        }
    }
    return nil
}

func (ls *LocatorState) MasterMetaTranstion(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (locator.SlaveLocatingState, locator.SlaveLocatingTransition, error) {
    var (
        transition locator.SlaveLocatingTransition
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

    return ls.locatingState, finalTransitionCandidate, summarizeErrors(transErr, eventErr)
}

// --- TRANSMISSION CONTROL
func (ls *LocatorState) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    return nil
}

func checkTxStateWithTime(ls *LocatorState, slaveTimestamp time.Time) locator.SlaveLocatingTransition {
    if ls.txActionCount < ls.txActionLimit() {
        // if tx timeout window is smaller than time delta (T_1 - T_0), don't do anything!!! just skip!
        if ls.txTimeout() < slaveTimestamp.Sub(ls.lastTransitionTS) {
            ls.txActionCount++
        }
        return locator.SlaveTransitionIdle
    }
    // this is failure. the fact that this is called indicate that we're ready to move to failure state
    return locator.SlaveTransitionFail
}

func (ls *LocatorState) TimestampTranstion(slaveTimestamp time.Time) (locator.SlaveLocatingState, locator.SlaveLocatingTransition, error) {
    var (
        transition locator.SlaveLocatingTransition
        transErr error = nil
        eventErr error = nil
    )

    transition = checkTxStateWithTime(ls, slaveTimestamp)

    if transition == locator.SlaveTransitionIdle {
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

    return ls.locatingState, transition, summarizeErrors(transErr, eventErr)
}
