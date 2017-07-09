package vmagent

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/utils"
)

type VBoxMasterState int
const (
    VBoxMasterKeyExchange       VBoxMasterState = iota
    VBoxMasterBounded
    VBoxMasterBindBroken
)

type VBoxMasterTransition int
const (
    VBoxMasterTransitionFail    VBoxMasterTransition = iota
    VBoxMasterTransitionOk
    VBoxMasterTransitionIdle
)

func (s VBoxMasterState) String() string {
    var state string
    switch s {
        case VBoxMasterKeyExchange:
            state = "VBoxMasterKeyExchange"
        case VBoxMasterBounded:
            state = "VBoxMasterBounded"
        case VBoxMasterBindBroken:
            state = "VBoxMasterBindBroken"
    }
    return state
}

const (
    TransitionFailureLimit      int           = 6
    // TODO : timeout mechanism for receiving slave meta
    // TransitionTimeout           time.Duration = time.Second * 10

    TxActionLimit               int           = 6
    UnboundedTimeout            time.Duration = time.Second
    BoundedTimeout              time.Duration = time.Second
)

type CommChannel interface {
    //McastSend(data []byte) error
    UcastSend(target string, data []byte) error
}

type CommChannelFunc func(target string, data []byte) error
func (c CommChannelFunc) UcastSend(target string, data []byte) error {
    return c(target, data)
}

type ControlOnTransitionEvent interface {
    OnStateTranstionSuccess(state VBoxMasterState, vcore interface{}, ts time.Time) error
    OnStateTranstionFailure(state VBoxMasterState, vcore interface{}, ts time.Time) error
}

// MasterBeacon is assigned individually for each slave node.
type VBoxMasterControl interface {
    CurrentState() VBoxMasterState
    TransitionWithTimestamp(timestamp time.Time) error
    TransitionWithCoreMeta(sender interface{}, metaPackage []byte, timestamp time.Time) error
    Shutdown()
}

// this interface is purely internal interface containing only functions, and could be replaced anytime you're to call it
type vboxController interface {
    currentState() VBoxMasterState

    transitionWithCoreMeta(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error)
    transitionWithTimeStamp(master *masterControl, ts time.Time) error

    onStateTranstionSuccess(master *masterControl, ts time.Time) error
    onStateTranstionFailure(master *masterControl, ts time.Time) error
}

type masterControl struct {
    controller                  vboxController

    /* ---------------------------------- changing properties to record transaction --------------------------------- */
    // each time we try to make transtion and fail, count goes up.
    transitionFailureCount      int

    // last time successfully transitioned state
    lastTransitionTS            time.Time

    txActionCount               int

    // DO NOT SET ANY TIME ON THIS FIELD SO THE FIRST TX ACTION CAN BE DONE WITHIN THE CYCLE
    lastTransmissionTS          time.Time

    /* ---------------------------------------- all-states properties ----------------------------------------------- */
    rsaEncryptor                pcrypto.RsaEncryptor
    rsaDecryptor                pcrypto.RsaDecryptor
    coreNode                    interface{}

    // --------------------------------- onSuccess && onFailure external event -----------------------------------------
    ControlOnTransitionEvent
}

func (m *masterControl) CurrentState() VBoxMasterState {
    if m.controller == nil {
        log.Panic("[CRITICAL] vboxController cannot be null")
    }
    return m.controller.currentState()
}

/* ------------------------------------------ Core Meta Transition Functions ---------------------------------------- */
func (m *masterControl) transitionTimeout() time.Duration {
    switch m.CurrentState() {
        case VBoxMasterBounded: {
            return BoundedTimeout * time.Duration(TxActionLimit)
        }
        default: {
            return UnboundedTimeout * time.Duration(TxActionLimit)
        }
    }
}

func stateTransition(currentState VBoxMasterState, transitCondition VBoxMasterTransition) VBoxMasterState {
    var nextState VBoxMasterState

    switch transitCondition {
        // successfully transition to the next
        case VBoxMasterTransitionOk: {
            nextState = VBoxMasterBounded
        }

        // failed to transit
        case VBoxMasterTransitionFail: {
            switch currentState {
                case VBoxMasterKeyExchange: {
                    nextState = VBoxMasterKeyExchange
                }
                default: {
                    nextState = VBoxMasterBindBroken
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

func finalizeStateTransitionWithTimeout(master *masterControl, nextStateCandiate VBoxMasterTransition, masterTimestamp time.Time) VBoxMasterTransition {
    var nextConfirmedState VBoxMasterTransition

    switch nextStateCandiate {
        // As MasterTransitionOk does not check timewindow, it could grant an infinite timewindow to make transition.
        // This is indeed intented as it will give us a chance to handle racing situations. Plus, TransitionWithTimestamp()
        // should have squashed suspected beacons and that's the role of TransitionWithTimestamp()
        case VBoxMasterTransitionOk: {
            master.lastTransitionTS = masterTimestamp
            master.transitionFailureCount = 0
            nextConfirmedState = VBoxMasterTransitionOk
            break
        }
        default: {
            if master.transitionFailureCount < TransitionFailureLimit {
                master.transitionFailureCount++
            }

            if master.transitionFailureCount < TransitionFailureLimit && masterTimestamp.Sub(master.lastTransitionTS) < master.transitionTimeout() {
                nextConfirmedState = VBoxMasterTransitionIdle
            } else {
                nextConfirmedState = VBoxMasterTransitionFail
            }
        }
    }

    return nextConfirmedState
}


func runOnTransitionEvents(master *masterControl, newState, oldState VBoxMasterState, transition VBoxMasterTransition, masterTimestamp time.Time) error {
    var (
        ierr, oerr error = nil, nil
    )
    if master.controller == nil {
        log.Panic("[CRITICAL] vboxController cannot be null")
    }
    if newState != oldState {
        switch transition {
            case VBoxMasterTransitionOk: {
                ierr = master.controller.onStateTranstionSuccess(master, masterTimestamp)

                if master.ControlOnTransitionEvent != nil {
                    oerr = master.OnStateTranstionSuccess(master.CurrentState(), master.coreNode, masterTimestamp)
                }
                // TODO : we need to a way to formalize this
                return utils.SummarizeErrors(ierr, oerr)
            }

            case VBoxMasterTransitionFail: {
                ierr = master.controller.onStateTranstionFailure(master, masterTimestamp)

                if master.ControlOnTransitionEvent != nil {
                    oerr = master.OnStateTranstionFailure(master.CurrentState(), master.coreNode, masterTimestamp)
                }
                // TODO : we need to a way to formalize this
                return utils.SummarizeErrors(ierr, oerr)
            }
        }
    }
    return nil
}

func newControllerForState(ctrl vboxController, newState, oldState VBoxMasterState) vboxController {
    var (
        newBeaconState vboxController = nil
        err error = nil
    )
    if newState == oldState {
        return ctrl
    }

    switch newState {
        case VBoxMasterKeyExchange: {
            newBeaconState = stateKeyexchange()
        }
        case VBoxMasterBounded: {
            newBeaconState = stateBounded()
        }

        case VBoxMasterBindBroken:
            fallthrough
        default: {
            newBeaconState = stateBindbroken()
            if err != nil {
                // this should never happen
                log.Panic(errors.WithStack(err))
            }
        }
    }

    return newBeaconState
}

func (m *masterControl) TransitionWithCoreMeta(sender interface{}, metaPackage []byte, timestamp time.Time) error {
    var (
        newState, oldState VBoxMasterState = m.CurrentState(), m.CurrentState()
        transitionCandidate, finalTransition VBoxMasterTransition
        transErr, eventErr error = nil, nil
    )
    if m.controller == nil {
        log.Panic("[CRITICAL] slaveMetaTransition func cannot be null")
    }

    transitionCandidate, transErr = m.controller.transitionWithCoreMeta(m, sender, metaPackage, timestamp)

    // this is to apply failed time count and timeout window
    finalTransition = finalizeStateTransitionWithTimeout(m, transitionCandidate, timestamp)

    // finalize master beacon state
    newState = stateTransition(oldState, finalTransition)

    // execute on events
    eventErr = runOnTransitionEvents(m, newState, oldState, finalTransition, timestamp)

    // assign vbox controller for state
    m.controller = newControllerForState(m.controller, newState, oldState)

    // return combined errors
    return utils.SummarizeErrors(transErr, eventErr)
}

/* ----------------------------------------- Timestamp Transition Functions ----------------------------------------- */

func (m *masterControl) txTimeWindow() time.Duration {
    switch m.CurrentState() {
        case VBoxMasterBounded: {
            return BoundedTimeout
        }
        default: {
            return UnboundedTimeout
        }
    }
}

func stateTransitionWithTimestamp(master *masterControl, timestamp time.Time) (VBoxMasterTransition, error) {
    var transErr error = nil

    if master.txActionCount < TxActionLimit {

        // if tx timeout window is smaller than time delta (T_1 - T_0), don't do anything!!! just skip!
        if master.txTimeWindow() < timestamp.Sub(master.lastTransmissionTS) {

            transErr = master.controller.transitionWithTimeStamp(master, timestamp)
            // since an action is taken, the action counter goes up regardless of error
            master.txActionCount++
            // we'll reset slave action timestamp
            master.lastTransmissionTS = timestamp
        }

        return VBoxMasterTransitionIdle, transErr
    }

    // this is failure. the fact that this is called indicate that we're ready to move to failure state
    return VBoxMasterTransitionFail, errors.Errorf("[ERR] Transmission count has exceeded a given limit")
}

func (m *masterControl) TransitionWithTimestamp(timestamp time.Time) error {
    var (
        newState, oldState VBoxMasterState = m.CurrentState(), m.CurrentState()
        transition VBoxMasterTransition
        transErr, eventErr error = nil, nil
    )
    if m.controller == nil {
        log.Panic("[CRITICAL] slaveMetaTransition func cannot be null")
    }

    transition, transErr = stateTransitionWithTimestamp(m, timestamp)

    // finalize state
    newState = stateTransition(oldState, transition)

    // event
    eventErr = runOnTransitionEvents(m, newState, oldState, transition, timestamp)

    // assign vbox controller for state
    m.controller = newControllerForState(m.controller, newState, oldState)

    // return combined errors
    return utils.SummarizeErrors(transErr, eventErr)
}

func (m *masterControl) Shutdown() {

}
