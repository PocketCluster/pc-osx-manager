package masterctrl

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/model"

    "github.com/stkim1/pc-vbox-comm/utils"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
)

type VBoxMasterTransition int
const (
    VBoxMasterTransitionFail    VBoxMasterTransition = iota
    VBoxMasterTransitionOk
    VBoxMasterTransitionIdle
)

const (
    TransitionFailureLimit      int           = 3
    // TODO : timeout mechanism for receiving slave meta
    // TransitionTimeout           time.Duration = time.Second * 10

    TxActionLimit               int           = 3
    UnboundedTimeout            time.Duration = time.Second
    BoundedTimeout              time.Duration = time.Second * 3
)

/* ---------------------------------------------- Interface Definitions --------------------------------------------- */
type ControllerActionOnTransition interface {
    OnStateTranstionSuccess(state mpkg.VBoxMasterState, vcore interface{}, ts time.Time) error
    OnStateTranstionFailure(state mpkg.VBoxMasterState, vcore interface{}, ts time.Time) error
}

// MasterBeacon is assigned individually for each slave node.
type VBoxMasterControl interface {
    CurrentState() mpkg.VBoxMasterState

    ReadCoreMetaAndMakeMasterAck(sender interface{}, metaPackage []byte, timestamp time.Time) ([]byte, error)
    CheckTransitionTimeWindow(timestamp time.Time) error

    Shutdown()
}

// this interface is purely internal interface containing only functions, and could be replaced anytime you're to call it
type vboxController interface {
    currentState() mpkg.VBoxMasterState

    readCoreReport(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error)
    makeMasterAck(master *masterControl, ts time.Time) ([]byte, error)

    onStateTranstionSuccess(master *masterControl, ts time.Time) error
    onStateTranstionFailure(master *masterControl, ts time.Time) error
}

/* ----------------------------------------------- Instance Definitions --------------------------------------------- */
func NewVBoxMasterControl(prvkey, pubkey []byte, coreNode *model.CoreNode, eventAction ControllerActionOnTransition) (VBoxMasterControl, error) {
    // TODO check if controller is bounded or unbounded
    var (
        controller vboxController = nil
        encryptor pcrypto.RsaEncryptor = nil
        decryptor pcrypto.RsaDecryptor = nil
        err error = nil
    )
    if prvkey == nil {
        return nil, errors.Errorf("[ERR] private key cannot be null")
    }
    if pubkey == nil {
        return nil, errors.Errorf("[ERR] public key cannot be null")
    }
    if coreNode == nil {
        return nil, errors.Errorf("[ERR] corenode model instance cannot be null")
    }
    _, err = coreNode.GetAuthToken()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // unbounded
    if coreNode.State == model.SNMStateInit && len(coreNode.PublicKey) == 0 {
        controller = stateUnbounded()

    // bind broken
    } else {
        controller = stateBindbroken()
        encryptor, err = pcrypto.NewRsaEncryptorFromKeyData(coreNode.PublicKey, prvkey)
        if err != nil {
            return nil, errors.WithStack(err)
        }
        decryptor, err = pcrypto.NewRsaDecryptorFromKeyData(coreNode.PublicKey, prvkey)
        if err != nil {
            return nil, errors.WithStack(err)
        }
    }

    return &masterControl {
        controller:      controller,
        privateKey:      prvkey,
        publicKey:       pubkey,
        rsaEncryptor:    encryptor,
        rsaDecryptor:    decryptor,
        coreNode:        coreNode,
        eventAction:     eventAction,
    }, nil
}

type masterControl struct {
    controller                  vboxController

    /* ---------------------------------- changing properties to record transaction --------------------------------- */
    // each time we try to make transtion and fail, count goes up.
    transitionActionCount       int

    // last time successfully transitioned state
    lastTransitionTS            time.Time

    txActionCount               int

    // DO NOT SET ANY TIME ON THIS FIELD SO THE FIRST TX ACTION CAN BE DONE WITHIN THE CYCLE
    lastTransmissionTS          time.Time

    /* ---------------------------------------- all-states properties ----------------------------------------------- */
    privateKey                  []byte
    publicKey                   []byte
    rsaEncryptor                pcrypto.RsaEncryptor
    rsaDecryptor                pcrypto.RsaDecryptor
    coreNode                    *model.CoreNode

    // --------------------------------- onSuccess && onFailure external event -----------------------------------------
    eventAction                 ControllerActionOnTransition
}

func (m *masterControl) CurrentState() mpkg.VBoxMasterState {
    if m.controller == nil {
        log.Panic("[CRITICAL] vboxController cannot be null")
    }
    return m.controller.currentState()
}

/* ------------------------------------------ Core Meta Transition Functions ---------------------------------------- */
func (m *masterControl) transitionTimeout() time.Duration {
    switch m.CurrentState() {
        case mpkg.VBoxMasterBounded: {
            return BoundedTimeout * time.Duration(TxActionLimit)
        }
        default: {
            return UnboundedTimeout * time.Duration(TxActionLimit)
        }
    }
}

func stateTransition(currentState mpkg.VBoxMasterState, transitCondition VBoxMasterTransition) mpkg.VBoxMasterState {
    var nextState mpkg.VBoxMasterState

    switch transitCondition {
        // successfully transition to the next
        case VBoxMasterTransitionOk: {
            switch currentState {
                case mpkg.VBoxMasterUnbounded: {
                    nextState = mpkg.VBoxMasterKeyExchange
                }
                default: {
                    nextState = mpkg.VBoxMasterBounded
                }
            }
        }

        // failed to transit
        case VBoxMasterTransitionFail: {
            switch currentState {
                case mpkg.VBoxMasterUnbounded:
                    fallthrough
                case mpkg.VBoxMasterKeyExchange: {
                    nextState = mpkg.VBoxMasterUnbounded
                }
                default: {
                    nextState = mpkg.VBoxMasterBindBroken
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
        // This is indeed intented as it will give us a chance to handle racing situations. Plus, CheckTransitionTimeWindow()
        // should have squashed suspected beacons and that's the role of CheckTransitionTimeWindow()
        case VBoxMasterTransitionOk: {
            master.lastTransitionTS = masterTimestamp
            master.transitionActionCount = 0
            nextConfirmedState = VBoxMasterTransitionOk
            break
        }
        default: {
            if master.transitionActionCount < TransitionFailureLimit {
                master.transitionActionCount++
            }

            if master.transitionActionCount < TransitionFailureLimit && masterTimestamp.Sub(master.lastTransitionTS) < master.transitionTimeout() {
                nextConfirmedState = VBoxMasterTransitionIdle
            } else {
                nextConfirmedState = VBoxMasterTransitionFail
            }
        }
    }

    return nextConfirmedState
}

func runOnTransitionEvents(master *masterControl, newState, oldState mpkg.VBoxMasterState, transition VBoxMasterTransition, masterTimestamp time.Time) error {
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

                if master.eventAction != nil {
                    oerr = master.eventAction.OnStateTranstionSuccess(master.CurrentState(), master.coreNode, masterTimestamp)
                }
                // TODO : we need to a way to formalize this
                return utils.SummarizeErrors(ierr, oerr)
            }

            case VBoxMasterTransitionFail: {
                ierr = master.controller.onStateTranstionFailure(master, masterTimestamp)

                if master.eventAction != nil {
                    oerr = master.eventAction.OnStateTranstionFailure(master.CurrentState(), master.coreNode, masterTimestamp)
                }
                // TODO : we need to a way to formalize this
                return utils.SummarizeErrors(ierr, oerr)
            }
        }
    }
    return nil
}

func newControllerForState(ctrl vboxController, newState, oldState mpkg.VBoxMasterState) vboxController {
    var (
        newController vboxController = nil
        err error = nil
    )
    if newState == oldState {
        return ctrl
    }

    switch newState {
        case mpkg.VBoxMasterUnbounded: {
            newController = stateUnbounded()
        }
        case mpkg.VBoxMasterKeyExchange: {
            newController = stateKeyexchange()
        }
        case mpkg.VBoxMasterBounded: {
            newController = stateBounded()
        }
        case mpkg.VBoxMasterBindBroken:
            fallthrough
        default: {
            newController = stateBindbroken()
            if err != nil {
                // this should never happen
                log.Panic(errors.WithStack(err))
            }
        }
    }

    return newController
}

func (m *masterControl) ReadCoreMetaAndMakeMasterAck(sender interface{}, metaPackage []byte, timestamp time.Time) ([]byte, error) {
    var (
        newState, oldState mpkg.VBoxMasterState = m.CurrentState(), m.CurrentState()
        transitionCandidate, finalTransition VBoxMasterTransition
        transErr, eventErr, buildErr, tErr error = nil, nil, nil, nil
        masterAck []byte = nil
    )
    if m.controller == nil {
        log.Panic("[CRITICAL] vboxController func cannot be null")
    }

    // ------------------------------------------- read core status ----------------------------------------------------
    transitionCandidate, transErr = m.controller.readCoreReport(m, sender, metaPackage, timestamp)

    // this is to apply failed time count and timeout window
    finalTransition = finalizeStateTransitionWithTimeout(m, transitionCandidate, timestamp)

    // finalize master controller state
    newState = stateTransition(oldState, finalTransition)

    // execute on events
    eventErr = runOnTransitionEvents(m, newState, oldState, finalTransition, timestamp)

    // assign vbox controller for new state
    m.controller = newControllerForState(m.controller, newState, oldState)

    // -------------------------------------- make master acknowledgement ----------------------------------------------
    masterAck, buildErr = m.controller.makeMasterAck(m, timestamp)

    // return combined errors
    tErr = utils.SummarizeErrors(transErr, eventErr)
    return masterAck, errors.WithStack(utils.SummarizeErrors(tErr, buildErr))
}

/* ----------------------------------------- Timestamp Transition Functions ----------------------------------------- */
func (m *masterControl) txTimeWindow() time.Duration {
    switch m.CurrentState() {
        case mpkg.VBoxMasterBounded: {
            return BoundedTimeout
        }
        default: {
            return UnboundedTimeout
        }
    }
}

func stateCheckTransitionTimeWindow(master *masterControl, timestamp time.Time) (VBoxMasterTransition, error) {
    if master.txActionCount < TxActionLimit {
        // if tx timeout window is smaller than time delta (T_1 - T_0), don't do anything!!! just skip!
        if master.txTimeWindow() < timestamp.Sub(master.lastTransmissionTS) {
            // since an action is taken, the action counter goes up regardless of error
            master.txActionCount++
            // we'll reset slave action timestamp
            master.lastTransmissionTS = timestamp
        }
        return VBoxMasterTransitionIdle, nil
    }

    // this is failure. the fact that this is called indicate that we're ready to move to failure state
    return VBoxMasterTransitionFail, errors.Errorf("[ERR] transmission count has exceeded a given limit")
}

func (m *masterControl) CheckTransitionTimeWindow(timestamp time.Time) error {
    var (
        newState, oldState mpkg.VBoxMasterState = m.CurrentState(), m.CurrentState()
        transition VBoxMasterTransition
        transErr, eventErr error = nil, nil
    )
    if m.controller == nil {
        log.Panic("[CRITICAL] vboxController func cannot be null")
    }

    transition, transErr = stateCheckTransitionTimeWindow(m, timestamp)

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
