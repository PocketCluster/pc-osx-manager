package corereport

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-vbox-comm/utils"
)

type VBoxCoreState int
const (
    VBoxCoreUnbounded           VBoxCoreState = iota
    VBoxCoreBounded
    VBoxCoreBindBroken
)

type VBoxCoreTransition int
const (
    VBoxCoreTransitionFail      VBoxCoreTransition = iota
    VBoxCoreTransitionOk
    VBoxCoreTransitionIdle
)

func (s VBoxCoreState) String() string {
    var state string
    switch s {
        case VBoxCoreUnbounded:
            state = "VBoxCoreUnbounded"
        case VBoxCoreBounded:
            state = "VBoxCoreBounded"
        case VBoxCoreBindBroken:
            state = "VBoxCoreBindBroken"
    }
    return state
}

const (
    TransitionFailureLimit      int           = 3

    // TODO : timeout mechanism for receiving master meta
    // Currently (v0.1.4), there is no timeout mechanism implemented for receiving master meta (i.e. if not master
    // response for a certian amount of time, state goes to failure mode). Instead, (TxActionLimit * Unbounded or Bounded)
    // times out the crrent state. When TxAction does not work, state will stall. We'll reinvestigate in the future
    //TransitionTimeout      time.Duration = time.Second * 10

    TxActionLimit               int           = 3
    UnboundedTimeout            time.Duration = time.Second
    BoundedTimeout              time.Duration = time.Second * 3
)

type CommChannel interface {
    //McastSend(data []byte) error
    UcastSend(target string, data []byte) error
}

type CommChannelFunc func(target string, data []byte) error
func (c CommChannelFunc) UcastSend(target string, data []byte) error {
    return c(target, data)
}

type ReporterActionsOnTransition interface {
    OnStateTranstionSuccess(state VBoxCoreState, ts time.Time) error
    OnStateTranstionFailure(state VBoxCoreState, ts time.Time) error
}

// MasterBeacon is assigned individually for each slave node.
type VBoxCoreReporter interface {
    CurrentState() VBoxCoreState
    TransitionWithCoreMeta(sender interface{}, metaPackage []byte, ts time.Time) error
    TransitionWithTimestamp(ts time.Time) error
}

type vboxReporter interface {
    currentState() VBoxCoreState

    transitionWithMasterMeta(core *coreReporter, sender interface{}, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error)
    transitionWithTimeStamp(core *coreReporter, ts time.Time) error

    onStateTranstionSuccess(core *coreReporter, ts time.Time) error
    onStateTranstionFailure(core *coreReporter, ts time.Time) error
}

type coreReporter struct {
    reporter                     vboxReporter

    /* ---------------------------------- changing properties to record transaction --------------------------------- */
    // each time we try to make transtion and fail, count goes up.
    transitionActionCount       int

    // last time successfully transitioned state.
    lastTransitionTS            time.Time

    // each time we try to send something, count goes up. This include success/fail altogether.
    txActionCount               int

    // last time transmission takes place. This is to control the frequnecy of transmission
    // !!!IMPORTANT!!! BY NOT SETTING A PARTICULAR VALUE, BY NOT SETTING ANYTHING, WE WILL AUTOMATICALLY EXECUTE
    // TX ACTION ON THE IDLE CYCLE RIGHT AFTER A SUCCESSFUL TRANSITION. SO DO NOT SET ANTYHING IN CONSTRUCTION
    lastTransmissionTS          time.Time

    /* ---------------------------------------- all-states properties ----------------------------------------------- */
    publicKey                   []byte
    privateKey                  []byte
    rsaEncryptor                pcrypto.RsaEncryptor
    rsaDecryptor                pcrypto.RsaDecryptor
    authToken                   string

    // --------------------------------- onSuccess && onFailure external event -----------------------------------------
    ReporterActionsOnTransition

    // -------------------------------------  Communication Channel ----------------------------------------------------
    CommChannel

}

func (c *coreReporter) CurrentState() VBoxCoreState {
    if c.reporter == nil {
        log.Panic("[CRITICAL] vboxReporter cannot be null")
    }
    return c.reporter.currentState()
}

/* ------------------------------------- Master Meta Transition Functions ------------------------------------------- */
func (c *coreReporter) transitionTimeout() time.Duration {
    switch c.CurrentState() {
        case VBoxCoreBounded: {
            return BoundedTimeout * time.Duration(TxActionLimit)
        }
        default: {
            return UnboundedTimeout * time.Duration(TxActionLimit)
        }
    }
}

func stateTransition(currentState VBoxCoreState, transitCondition VBoxCoreTransition) VBoxCoreState {
    var nextState VBoxCoreState

    switch transitCondition {
        // successfully transition to the next
        case VBoxCoreTransitionOk: {
            nextState = VBoxCoreBounded
        }

        // failed to transit
        case VBoxCoreTransitionFail: {
            switch currentState {
                case VBoxCoreUnbounded: {
                    nextState = VBoxCoreUnbounded
                }
                default: {
                    nextState = VBoxCoreBindBroken
                }
            }
        }

        // idle
        default: {
            nextState = currentState
        }
    }

    return nextState
}

func finalizeTransitionWithMeta(core *coreReporter, nextStateCandiate VBoxCoreTransition, coreTimestamp time.Time) VBoxCoreTransition {
    var nextConfirmedState VBoxCoreTransition
    switch nextStateCandiate {
        case VBoxCoreTransitionOk: {
            // reset transition action count / timestamp to 0
            core.transitionActionCount = 0
            core.lastTransitionTS = coreTimestamp
            core.txActionCount = 0
            nextConfirmedState = VBoxCoreTransitionOk
        }
        default: {
            if core.transitionActionCount < TransitionFailureLimit {
                core.transitionActionCount++
            }

            if core.transitionActionCount < TransitionFailureLimit && coreTimestamp.Sub(core.lastTransitionTS) < core.transitionTimeout() {
                nextConfirmedState = VBoxCoreTransitionIdle
            } else {
                nextConfirmedState = VBoxCoreTransitionFail
            }
        }
    }
    return nextConfirmedState
}

func runOnTransitionEvents(core *coreReporter, newState, oldState VBoxCoreState, transition VBoxCoreTransition, ts time.Time) error {
    var (
        ierr, oerr error = nil, nil
    )
    if core.reporter == nil {
        log.Panic("[CRITICAL] vboxReporter cannot be null")
    }
    if newState != oldState {
        switch transition {
            case VBoxCoreTransitionOk: {
                ierr = core.reporter.onStateTranstionSuccess(core, ts)

                if core.ReporterActionsOnTransition != nil {
                    oerr = core.OnStateTranstionSuccess(core.CurrentState(), ts)
                }
                // TODO : we need to a way to formalize this
                return utils.SummarizeErrors(ierr, oerr)
            }

            case VBoxCoreTransitionFail: {
                ierr = core.reporter.onStateTranstionFailure(core, ts)

                if core.ReporterActionsOnTransition != nil {
                    oerr = core.OnStateTranstionFailure(core.CurrentState(), ts)
                }
                // TODO : we need to a way to formalize this
                return utils.SummarizeErrors(ierr, oerr)
            }
        }
    }
    return nil
}

func newReporterForState(reporter vboxReporter, newState, oldState VBoxCoreState) vboxReporter {
    var (
        newReporter vboxReporter = nil
        err error = nil
    )
    if newState == oldState {
        return reporter
    }

    switch newState {
        case VBoxCoreUnbounded: {
            newReporter = stateUnbounded()
        }
        case VBoxCoreBounded: {
            newReporter = stateBounded()
        }
        case VBoxCoreBindBroken:
            fallthrough
        default: {
            newReporter = stateBindbroken()
            if err != nil {
                // this should never happen
                log.Panic(errors.WithStack(err))
            }
        }
    }

    return newReporter
}

func (c *coreReporter) TransitionWithCoreMeta(sender interface{}, metaPackage []byte, timestamp time.Time) error {
    var (
        transitionCandidate, finalTransition VBoxCoreTransition
        transErr, eventErr error = nil, nil
        newState, oldState VBoxCoreState = c.CurrentState(), c.CurrentState()
    )
    if c.reporter == nil {
        log.Panic("[CRITICAL] vboxReporter cannot be null")
    }

    transitionCandidate, transErr = c.reporter.transitionWithMasterMeta(c, sender, metaPackage, timestamp)

    // filter out the intermediate transition value with failed count + timestamp
    finalTransition = finalizeTransitionWithMeta(c, transitionCandidate, timestamp)

    // finalize core reporter state
    newState = stateTransition(oldState, finalTransition)

    // execute on event
    eventErr = runOnTransitionEvents(c, newState, oldState, finalTransition, timestamp)

    // assign vbox reporter for new state
    c.reporter = newReporterForState(c.reporter, newState, oldState)

    return utils.SummarizeErrors(transErr, eventErr)
}

/* ----------------------------------------- Timestamp Transition Functions ----------------------------------------- */
func (c *coreReporter) txTimeWindow() time.Duration {
    switch c.CurrentState() {
        case VBoxCoreBounded: {
            return BoundedTimeout
        }
        default: {
            return UnboundedTimeout
        }
    }
}

func stateTransitionWithTimestamp(core *coreReporter, timestamp time.Time) (VBoxCoreTransition, error) {
    var transErr error = nil

    if core.txActionCount < TxActionLimit {

        // if tx timeout window is smaller than time delta (T_1 - T_0), don't do anything!!! just skip!
        if core.txTimeWindow() < timestamp.Sub(core.lastTransmissionTS) {

            transErr = core.reporter.transitionWithTimeStamp(core, timestamp)
            // since an action is taken, the action counter goes up regardless of error
            core.txActionCount++
            // we'll reset slave action timestamp
            core.lastTransmissionTS = timestamp
        }

        return VBoxCoreTransitionIdle, transErr
    }
    // this is failure. the fact that this is called indicate that we're ready to move to failure state
    return VBoxCoreTransitionFail, errors.Errorf("[ERR] transmission count has exceeded a given limit")
}

func (c *coreReporter) TransitionWithTimestamp(timestamp time.Time) error {
    var (
        newState, oldState VBoxCoreState = c.CurrentState(), c.CurrentState()
        transition VBoxCoreTransition
        transErr, eventErr error = nil, nil
    )
    if c.reporter == nil {
        log.Panic("[CRITICAL] vboxReporter cannot be null")
    }

    transition, transErr = stateTransitionWithTimestamp(c, timestamp)
    // now idle action condition has failed, and we need to make transition to FAILTURE state


    // finalize locating master beacon state
    newState = stateTransition(oldState, transition)

    // execute event lisenter
    eventErr = runOnTransitionEvents(c, newState, oldState, transition, timestamp)

    // assign vbox reporter for new state
    c.reporter = newReporterForState(c.reporter, newState, oldState)

    return utils.SummarizeErrors(transErr, eventErr)
}