package masterctrl

import (
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/model"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    "github.com/stkim1/pc-vbox-comm/utils"
    "github.com/stkim1/pcrypto"
)

type VBoxMasterTransition int
const (
    VBoxMasterTransitionFail    VBoxMasterTransition = iota
    VBoxMasterTransitionOk
    VBoxMasterTransitionIdle
)

const (
    TransitionFailureLimit      int           = 3
    TxActionLimit               int           = 3

    UnboundedTimeout            time.Duration = time.Second
    BoundedTimeout              time.Duration = time.Second
)

/* ---------------------------------------------- Interface Definitions --------------------------------------------- */
type ControllerActionOnTransition interface {
    OnStateTranstionSuccess(state mpkg.VBoxMasterState, vcore interface{}, ts time.Time) error
    OnStateTranstionFailure(state mpkg.VBoxMasterState, vcore interface{}, ts time.Time) error
}

// MasterBeacon is assigned individually for each slave node.
type VBoxMasterControl interface {
    CurrentState() mpkg.VBoxMasterState
    GetCoreNode() *model.CoreNode

    SetMasterIPv4ExternalAddress(addr string) error
    ClearMasterIPv4ExternalAddress() error

    ReadCoreMetaAndMakeMasterAck(sender interface{}, metaPackage []byte, timestamp time.Time) ([]byte, error)
    HandleCoreDisconnection(timestamp time.Time) error
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
func NewVBoxMasterControl(clusterID, extIP4Addr string, prvkey, pubkey []byte, coreNode *model.CoreNode, eventAction ControllerActionOnTransition) (VBoxMasterControl, error) {
    // TODO check if controller is bounded or unbounded
    var (
        encryptor pcrypto.RsaEncryptor = nil
        decryptor pcrypto.RsaDecryptor = nil
        err error = nil
    )
    if len(clusterID) == 0 {
        return nil, errors.Errorf("[ERR] invalid cluster id")
    }
    if len(extIP4Addr) == 0 {
        return nil, errors.Errorf("[ERR] invalid external ip4 address")
    }
    if prvkey == nil {
        return nil, errors.Errorf("[ERR] private key cannot be null")
    }
    if pubkey == nil {
        return nil, errors.Errorf("[ERR] public key cannot be null")
    }
    if coreNode == nil {
        return nil, errors.Errorf("[ERR] corenode model instance cannot be null")
    }
    encryptor, err = pcrypto.NewRsaEncryptorFromKeyData(coreNode.PublicKey, prvkey)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    decryptor, err = pcrypto.NewRsaDecryptorFromKeyData(coreNode.PublicKey, prvkey)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    return &masterControl {
        controller:      stateBindbroken(),
        clusterID:       clusterID,
        extIP4Addr:      extIP4Addr,
        rsaEncryptor:    encryptor,
        rsaDecryptor:    decryptor,
        coreNode:        coreNode,
        eventAction:     eventAction,
    }, nil
}

type masterControl struct {
    sync.Mutex
    controller                  vboxController

    /* ---------------------------------- changing properties to record transaction --------------------------------- */
    // each time we try to make transtion and fail, count goes up.
    transitionActionCount       int

    // last time successfully transitioned state
    lastTransitionTS            time.Time

    /* ---------------------------------------- all-states properties ----------------------------------------------- */
    clusterID                   string
    extIP4Addr                  string
    rsaEncryptor                pcrypto.RsaEncryptor
    rsaDecryptor                pcrypto.RsaDecryptor
    coreNode                    *model.CoreNode

    // --------------------------------- onSuccess && onFailure external event -----------------------------------------
    eventAction                 ControllerActionOnTransition
}

func (m *masterControl) CurrentState() mpkg.VBoxMasterState {
    m.Lock()
    defer m.Unlock()

    if m.controller == nil {
        log.Panic("[CRITICAL] vboxController cannot be null")
    }
    return m.controller.currentState()
}

func (m *masterControl) GetCoreNode() *model.CoreNode {
    return m.coreNode
}

func (m *masterControl) SetMasterIPv4ExternalAddress(addr string) error {
    m.Lock()
    defer m.Unlock()

    if len(addr) == 0 {
        return errors.Errorf("[ERR] cannot assign invalid external IPv4 address")
    }

    m.extIP4Addr = addr
    return nil
}

func (m *masterControl) ClearMasterIPv4ExternalAddress() error {
    m.Lock()
    defer m.Unlock()

    m.extIP4Addr = ""
    return nil
}

// this is a helper internal method
func (m *masterControl) getMasterIPv4ExternalAddress() string {
    // as master external ipv4 is accessed, we need to synchronize.
    m.Lock()
    defer m.Unlock()

    return m.extIP4Addr
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
            nextState = mpkg.VBoxMasterBounded
        }

        // failed to transit
        case VBoxMasterTransitionFail: {
            nextState = mpkg.VBoxMasterBindBroken
        }

        // idle
        default: {
            nextState = currentState
        }
    }

    return nextState
}

func finalizeStateTransitionWithTimeout(master *masterControl, nextStateCandiate VBoxMasterTransition, timestamp time.Time) VBoxMasterTransition {
    var nextConfirmedState VBoxMasterTransition

    switch nextStateCandiate {
        // As MasterTransitionOk does not check timewindow, it could grant an infinite timewindow to make transition.
        // This is indeed intented as it will give us a chance to handle racing situations. Plus, CheckTransitionTimeWindow()
        // should have squashed suspected beacons and that's the role of CheckTransitionTimeWindow()
        // TODO : need to think about how to reset variables
        case VBoxMasterTransitionOk: {
            master.transitionActionCount = 0
            master.lastTransitionTS = timestamp
            nextConfirmedState = VBoxMasterTransitionOk
            break
        }
        default: {
            if master.transitionActionCount < TransitionFailureLimit {
                master.transitionActionCount++
            }

            if master.transitionActionCount < TransitionFailureLimit && timestamp.Sub(master.lastTransitionTS) < master.transitionTimeout() {
                nextConfirmedState = VBoxMasterTransitionIdle
            } else {
                nextConfirmedState = VBoxMasterTransitionFail
            }
        }
    }

    return nextConfirmedState
}

func runOnTransitionEvents(master *masterControl, newState, oldState mpkg.VBoxMasterState, transition VBoxMasterTransition, timestamp time.Time) error {
    var (
        ierr, oerr error = nil, nil
    )
    if master.controller == nil {
        log.Panic("[CRITICAL] vboxController cannot be null")
    }
    if newState != oldState {
        switch transition {
            case VBoxMasterTransitionOk: {
                ierr = master.controller.onStateTranstionSuccess(master, timestamp)

                if master.eventAction != nil {
                    oerr = master.eventAction.OnStateTranstionSuccess(master.CurrentState(), master.coreNode, timestamp)
                }
                // TODO : we need to a way to formalize this
                return utils.SummarizeErrors(ierr, oerr)
            }

            case VBoxMasterTransitionFail: {
                ierr = master.controller.onStateTranstionFailure(master, timestamp)

                if master.eventAction != nil {
                    oerr = master.eventAction.OnStateTranstionFailure(master.CurrentState(), master.coreNode, timestamp)
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
    )
    if newState == oldState {
        return ctrl
    }

    switch newState {
        case mpkg.VBoxMasterBounded: {
            newController = stateBounded()
        }
        default: {
            newController = stateBindbroken()
        }
    }

    return newController
}

func (m *masterControl) ReadCoreMetaAndMakeMasterAck(sender interface{}, metaPackage []byte, timestamp time.Time) ([]byte, error) {
    var (
        oldState mpkg.VBoxMasterState = m.CurrentState()
        newState mpkg.VBoxMasterState = mpkg.VBoxMasterBindBroken
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
    m.Lock()
    m.controller = newControllerForState(m.controller, newState, oldState)
    m.Unlock()

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

func (m *masterControl) HandleCoreDisconnection(timestamp time.Time) error {
    var (
        oldState mpkg.VBoxMasterState = m.CurrentState()
        newState mpkg.VBoxMasterState = mpkg.VBoxMasterBindBroken
        transErr, eventErr error = nil, nil
    )
    if m.controller == nil {
        log.Panic("[CRITICAL] vboxController func cannot be null")
    }

    // finalize state
    newState = stateTransition(oldState, VBoxMasterTransitionFail)

    // event
    eventErr = runOnTransitionEvents(m, newState, oldState, VBoxMasterTransitionFail, timestamp)

    // assign vbox controller for state
    m.Lock()
    m.controller = newControllerForState(m.controller, newState, oldState)
    m.Unlock()

    // return combined errors
    return utils.SummarizeErrors(transErr, eventErr)
}
