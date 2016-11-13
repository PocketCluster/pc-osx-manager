package locator

import (
    "time"

    "github.com/stkim1/pc-core/msagent"
    "log"
    "fmt"
)

const (
    TransitionFailureLimit uint          = 5

    // TODO : timeout mechanism for receiving master meta
    // Currently (v0.1.4), there is no timeout mechanism implemented for receiving master meta (i.e. if not master
    // response for a certian amount of time, state goes to failure mode). Instead, (TxActionLimit * Unbounded or Bounded)
    // times out the crrent state. When TxAction does not work, state will stall. We'll reinvestigate in the future
    //TransitionTimeout      time.Duration = time.Second * 10

    TxActionLimit          uint          = 5
    UnboundedTimeout       time.Duration = time.Second * 3
    BoundedTimeout         time.Duration = time.Second * 10
)

type transitionWithMasterMeta       func (meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error)

type transitionActionWithTimestamp  func (slaveTimestamp time.Time) error

type onStateTranstionSuccess        func (slaveTimestamp time.Time) error

type onStateTranstionFailure        func (slaveTimestamp time.Time) error

type LocatorState interface {
    CurrentState() SlaveLocatingState
    MasterMetaTransition(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (LocatorState, error)
    TimestampTransition(slaveTimestamp time.Time) (LocatorState, error)
}

type locatorState struct {
    /* -------------------------------------- given, constant field ------------------------------------------------- */
    // this is given state that will not change
    constState                  SlaveLocatingState

    // transition failure
    constTransitionFailureLimit uint

    // transition timeout
    constTransitionTimout       time.Duration

    // transmission limit
    constTxActionLimit          uint

    // unbounded timeout
    constTxTimeWindow           time.Duration

    /* ---------------------------------- changing properties to record transaction --------------------------------- */
    // each time we try to make transtion and fail, count goes up.
    transitionActionCount       uint

    // last time successfully transitioned state.
    lastTransitionTS            time.Time

    // each time we try to send something, count goes up. This include success/fail altogether.
    txActionCount               uint

    // last time transmission takes place. This is to control the frequnecy of transmission
    // !!!IMPORTANT!!! BY NOT SETTING A PARTICULAR VALUE, BY NOT SETTING ANYTHING, WE WILL AUTOMATICALLY EXECUTE
    // TX ACTION ON THE IDLE CYCLE RIGHT AFTER A SUCCESSFUL TRANSITION. SO DO NOT SET ANTYHING IN CONSTRUCTION
    lastTxTS                    time.Time

    /* ----------------------------------------- transition functions ----------------------------------------------- */
    // master transition func
    masterMetaTransition        transitionWithMasterMeta

    // timestamp transition func
    timestampTransition         transitionActionWithTimestamp

    // onSuccess
    onTransitionSuccess         onStateTranstionSuccess

    // onFailure
    onTransitionFailure         onStateTranstionFailure

    /* ---------------------------------------- Communication Channel ----------------------------------------------- */
    commChannel                 CommChannel
}

// property functions
func (ls *locatorState) CurrentState() SlaveLocatingState {
    return ls.constState
}

func (ls *locatorState) transitionFailureLimit() uint {
    return ls.constTransitionFailureLimit
}

func (ls *locatorState) transitionTimeout() time.Duration {
    return ls.constTransitionTimout
}

func (ls *locatorState) txActionLimit() uint {
    return ls.constTxActionLimit
}

func (ls *locatorState) txTimeWindow() time.Duration {
    return ls.constTxTimeWindow
}

/* ------------------------------------------------ Helper Functions ------------------------------------------------ */
// close func pointers and delegates to help GC
func (ls *locatorState) Close() error {
    ls.masterMetaTransition  = nil
    ls.timestampTransition   = nil
    ls.onTransitionSuccess   = nil
    ls.onTransitionFailure   = nil
    ls.commChannel           = nil
    return nil
}

// -- STATE TRANSITION
func newLocatorStateForState(ls *locatorState, newState, oldState SlaveLocatingState) LocatorState {
    if newState == oldState {
        return ls
    }

    var comm CommChannel = ls.commChannel
    var newLocatorState LocatorState = nil
    switch newState {
        case SlaveUnbounded:
            newLocatorState = newUnboundedState(comm)

        case SlaveInquired:
            newLocatorState = newInquiredState(comm)

        case SlaveKeyExchange:
            newLocatorState = newKeyexchangeState(comm)

        case SlaveCryptoCheck:
            newLocatorState = newCryptocheckState(comm)

        case SlaveBounded:
            newLocatorState = newBoundedState(comm)

        case SlaveBindBroken:
            newLocatorState = newBindbrokenState(comm)

        default:
            newLocatorState = newUnboundedState(comm)
    }
    // invalidate old LocatorState CommChannel for GC
    ls.Close()
    return newLocatorState
}

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
            if ls.transitionActionCount < ls.transitionFailureLimit() {
                ls.transitionActionCount++
            }

            if ls.transitionActionCount < ls.transitionFailureLimit() && slaveTimestamp.Sub(ls.lastTransitionTS) < ls.transitionTimeout() {
                nextConfirmedState = SlaveTransitionIdle
            } else {
                nextConfirmedState = SlaveTransitionFail
            }
        }
    }
    return nextConfirmedState
}

func executeOnTransitionEvents(ls *locatorState, newState, oldState SlaveLocatingState, transition SlaveLocatingTransition, slaveTimestamp time.Time) error {
    if newState != oldState {
        switch transition {
            case SlaveTransitionOk:
                if ls.onTransitionSuccess != nil {
                    return ls.onTransitionSuccess(slaveTimestamp)
                }
            case SlaveTransitionFail: {
                if ls.onTransitionFailure != nil {
                    return ls.onTransitionFailure(slaveTimestamp)
                }
            }
        }
    }
    return nil
}

func (ls *locatorState) MasterMetaTransition(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (LocatorState, error) {
    var (
        transition SlaveLocatingTransition
        transErr, eventErr error = nil, nil
        newState, oldState SlaveLocatingState = ls.CurrentState(), ls.CurrentState()
    )
    if ls.masterMetaTransition == nil {
        log.Panic("[PANIC] MASTER META TRANSTION SHOULD HAVE BEEN SETUP PROPERLY")
    }

    transition, transErr = ls.masterMetaTransition(meta, slaveTimestamp)

    // filter out the intermediate transition value with failed count + timestamp
    finalTransitionCandidate := finalizeTransitionWithTimeout(ls, transition, slaveTimestamp)

    // finalize locating master beacon state
    newState = stateTransition(oldState, finalTransitionCandidate)

    // execute event lisenter
    eventErr = executeOnTransitionEvents(ls, newState, oldState, finalTransitionCandidate, slaveTimestamp)

    return newLocatorStateForState(ls, newState, oldState), summarizeErrors(transErr, eventErr)
}

// --- TRANSMISSION CONTROL
func runTxStateActionWithTimestamp(ls *locatorState, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    var transErr error = nil

    if ls.txActionCount < ls.txActionLimit() {

        // if tx timeout window is smaller than time delta (T_1 - T_0), don't do anything!!! just skip!
        if ls.txTimeWindow() < slaveTimestamp.Sub(ls.lastTxTS) {

            transErr = ls.timestampTransition(slaveTimestamp)
            // since an action is taken, the action counter goes up regardless of error
            ls.txActionCount++
            // we'll reset slave action timestamp
            ls.lastTxTS = slaveTimestamp
        }

        return SlaveTransitionIdle, transErr
    }
    // this is failure. the fact that this is called indicate that we're ready to move to failure state
    return SlaveTransitionFail, fmt.Errorf("[ERR] Transmission count has exceeded a given limit")
}

func (ls *locatorState) TimestampTransition(slaveTimestamp time.Time) (LocatorState, error) {
    var (
        newState, oldState SlaveLocatingState = ls.CurrentState(), ls.CurrentState()
        transition SlaveLocatingTransition
        transErr, eventErr error = nil, nil
    )
    if ls.timestampTransition == nil {
        log.Panic("[PANIC] TIMESTAMP TRANSITION SHOULD HAVE BEEN SETUP PROPERLY")
    }

    transition, transErr = runTxStateActionWithTimestamp(ls, slaveTimestamp)
    // now idle action condition has failed, and we need to make transition to FAILTURE state
    if transition == SlaveTransitionFail {

        // finalize locating master beacon state
        newState = stateTransition(oldState, transition)

        // execute event lisenter
        eventErr = executeOnTransitionEvents(ls, newState, oldState, transition, slaveTimestamp)
    }

    return newLocatorStateForState(ls, newState, oldState), summarizeErrors(transErr, eventErr)
}
