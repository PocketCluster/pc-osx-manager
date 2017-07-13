package corereport

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-vbox-comm/utils"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
)

type VBoxCoreTransition int
const (
    VBoxCoreTransitionFail      VBoxCoreTransition = iota
    VBoxCoreTransitionOk
    VBoxCoreTransitionIdle
)

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

/* ---------------------------------------------- Interface Definitions --------------------------------------------- */
type ReporterActionsOnTransition interface {
    OnStateTranstionSuccess(state cpkg.VBoxCoreState, ts time.Time) error
    OnStateTranstionFailure(state cpkg.VBoxCoreState, ts time.Time) error
}

// MasterBeacon is assigned individually for each slave node.
type VBoxCoreReporter interface {
    CurrentState() cpkg.VBoxCoreState
    MakeCoreReporter(timestamp time.Time) ([]byte, error)
    ReadMasterAcknowledgement(metaPackage []byte, timestamp time.Time) error
}

type vboxReporter interface {
    currentState() cpkg.VBoxCoreState

    makeCoreReport(core *coreReporter, ts time.Time) ([]byte, error)
    readMasterAck(core *coreReporter, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error)

    onStateTranstionSuccess(core *coreReporter, ts time.Time) error
    onStateTranstionFailure(core *coreReporter, ts time.Time) error
}

/* ----------------------------------------------- Instance Definitions --------------------------------------------- */
func NewCoreReporter(state cpkg.VBoxCoreState, corePrvkey, corePubkey, masterPubkey []byte) (VBoxCoreReporter, error) {
    var (
        rptr vboxReporter = nil
        encryptor pcrypto.RsaEncryptor = nil
        decryptor pcrypto.RsaDecryptor = nil
        err error = nil
    )
    // check errors first
    if len(corePrvkey) == 0 {
        return nil, errors.Errorf("[ERR] core private key cannot be null")
    }
    if len(corePubkey) == 0 {
        return nil, errors.Errorf("[ERR] core public key cannot be null")
    }
    if state != cpkg.VBoxCoreUnbounded && len(masterPubkey) == 0 {
        return nil, errors.Errorf("[ERR] master public key cannot be null")
    }
    switch state {
        case cpkg.VBoxCoreUnbounded: {
            rptr = stateUnbounded()
        }
        case cpkg.VBoxCoreBindBroken: {
            rptr = stateBindbroken()

            encryptor, err = pcrypto.NewRsaEncryptorFromKeyData(masterPubkey, corePrvkey)
            if err != nil {
                return nil, errors.WithStack(err)
            }
            decryptor, err = pcrypto.NewRsaDecryptorFromKeyData(masterPubkey, corePrvkey)
            if err != nil {
                return nil, errors.WithStack(err)
            }
        }
        default: {
            return nil, errors.Errorf("[ERR] core state should either VBoxCoreUnbounded or VBoxCoreBindBroken")
        }
    }

    return &coreReporter {
        reporter:        rptr,
        privateKey:      corePrvkey,
        publicKey:       corePubkey,
        rsaEncryptor:    encryptor,
        rsaDecryptor:    decryptor,
    }, nil
}

type coreReporter struct {
    reporter                     vboxReporter

    /* ---------------------------------- changing properties to record transaction --------------------------------- */
    // each time we try to make transtion and fail, count goes up.
    transitionActionCount    int

    // last time successfully transitioned state.
    lastTransitionTS         time.Time

    // each time we try to send something, count goes up. This include success/fail altogether.
    txActionCount            int

    // last time transmission takes place. This is to control the frequnecy of transmission
    // !!!IMPORTANT!!! BY NOT SETTING A PARTICULAR VALUE, BY NOT SETTING ANYTHING, WE WILL AUTOMATICALLY EXECUTE
    // TX ACTION ON THE IDLE CYCLE RIGHT AFTER A SUCCESSFUL TRANSITION. SO DO NOT SET ANTYHING IN CONSTRUCTION
    lastTransmissionTS       time.Time

    /* ---------------------------------------- all-states properties ----------------------------------------------- */
    privateKey               []byte
    publicKey                []byte
    rsaEncryptor             pcrypto.RsaEncryptor
    rsaDecryptor             pcrypto.RsaDecryptor
    authToken                string

    // --------------------------------- onSuccess && onFailure external event -----------------------------------------
    eventAction              ReporterActionsOnTransition
}

func (c *coreReporter) CurrentState() cpkg.VBoxCoreState {
    if c.reporter == nil {
        log.Panic("[CRITICAL] vboxReporter cannot be null")
    }
    return c.reporter.currentState()
}

// -------------------------------------- Make Core Reporter for Master ------------------------------------------------
func (c *coreReporter) MakeCoreReporter(timestamp time.Time) ([]byte, error) {
    return c.reporter.makeCoreReport(c, timestamp)
}

/* ------------------------------------- Master Meta Transition Functions ------------------------------------------- */
func (c *coreReporter) transitionTimeout() time.Duration {
    switch c.CurrentState() {
        case cpkg.VBoxCoreBounded: {
            return BoundedTimeout * time.Duration(TxActionLimit)
        }
        default: {
            return UnboundedTimeout * time.Duration(TxActionLimit)
        }
    }
}

func stateTransition(currentState cpkg.VBoxCoreState, transitCondition VBoxCoreTransition) cpkg.VBoxCoreState {
    var nextState cpkg.VBoxCoreState

    switch transitCondition {
        // successfully transition to the next
        case VBoxCoreTransitionOk: {
            nextState = cpkg.VBoxCoreBounded
        }

        // failed to transit
        case VBoxCoreTransitionFail: {
            switch currentState {
                case cpkg.VBoxCoreUnbounded: {
                    nextState = cpkg.VBoxCoreUnbounded
                }
                default: {
                    nextState = cpkg.VBoxCoreBindBroken
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

func finalizeTransitionWithMeta(core *coreReporter, nextStateCandiate VBoxCoreTransition, timestamp time.Time) VBoxCoreTransition {
    var nextConfirmedState VBoxCoreTransition
    switch nextStateCandiate {
        // TODO : need to think about how to reset variables
        case VBoxCoreTransitionOk: {
            // reset transition action count / timestamp to 0
            core.transitionActionCount = 0
            core.lastTransitionTS = timestamp
            core.txActionCount = 0
            core.lastTransmissionTS = timestamp
            nextConfirmedState = VBoxCoreTransitionOk
        }
        default: {
            if core.transitionActionCount < TransitionFailureLimit {
                core.transitionActionCount++
            }

            if core.transitionActionCount < TransitionFailureLimit && timestamp.Sub(core.lastTransitionTS) < core.transitionTimeout() {
                nextConfirmedState = VBoxCoreTransitionIdle
            } else {
                nextConfirmedState = VBoxCoreTransitionFail
            }
        }
    }
    return nextConfirmedState
}

func runOnTransitionEvents(core *coreReporter, newState, oldState cpkg.VBoxCoreState, transition VBoxCoreTransition, timestamp time.Time) error {
    var (
        ierr, oerr error = nil, nil
    )
    if core.reporter == nil {
        log.Panic("[CRITICAL] vboxReporter cannot be null")
    }
    if newState != oldState {
        switch transition {
            case VBoxCoreTransitionOk: {
                ierr = core.reporter.onStateTranstionSuccess(core, timestamp)

                if core.eventAction != nil {
                    oerr = core.eventAction.OnStateTranstionSuccess(core.CurrentState(), timestamp)
                }
                // TODO : we need to a way to formalize this
                return utils.SummarizeErrors(ierr, oerr)
            }

            case VBoxCoreTransitionFail: {
                ierr = core.reporter.onStateTranstionFailure(core, timestamp)

                if core.eventAction != nil {
                    oerr = core.eventAction.OnStateTranstionFailure(core.CurrentState(), timestamp)
                }
                // TODO : we need to a way to formalize this
                return utils.SummarizeErrors(ierr, oerr)
            }
        }
    }
    return nil
}

func newReporterForState(reporter vboxReporter, newState, oldState cpkg.VBoxCoreState) vboxReporter {
    var (
        newReporter vboxReporter = nil
        err error = nil
    )
    if newState == oldState {
        return reporter
    }

    switch newState {
        case cpkg.VBoxCoreUnbounded: {
            newReporter = stateUnbounded()
        }
        case cpkg.VBoxCoreBounded: {
            newReporter = stateBounded()
        }
        case cpkg.VBoxCoreBindBroken:
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

/* ----------------------------------------- Timestamp Transition Functions ----------------------------------------- */
func (c *coreReporter) txTimeWindow() time.Duration {
    switch c.CurrentState() {
        case cpkg.VBoxCoreBounded: {
            return BoundedTimeout
        }
        default: {
            return UnboundedTimeout
        }
    }
}

// TODO : need to think about how to reset variables
func stateTransitionWithTimestamp(core *coreReporter, timestamp time.Time) (VBoxCoreTransition, error) {
    if core.txActionCount < TxActionLimit {

        // if tx timeout window is smaller than time delta (T_1 - T_0), don't do anything!!! just skip!
        if core.txTimeWindow() < timestamp.Sub(core.lastTransmissionTS) {

            // since an action is taken, the action counter goes up regardless of error
            core.txActionCount++
            // we'll reset slave action timestamp
            core.lastTransmissionTS = timestamp
        }

        return VBoxCoreTransitionIdle, nil
    }
    // this is failure. the fact that this is called indicate that we're ready to move to failure state
    return VBoxCoreTransitionFail, errors.Errorf("[ERR] transmission count has exceeded a given limit")
}

func (c *coreReporter) ReadMasterAcknowledgement(metaPackage []byte, timestamp time.Time) error {
    var (
        newState, oldState cpkg.VBoxCoreState = c.CurrentState(), c.CurrentState()
        transition VBoxCoreTransition
        transErr, eventErr error = nil, nil
    )
    if c.reporter == nil {
        log.Panic("[CRITICAL] vboxReporter cannot be null")
    }

    // when there is ack from master...
    if len(metaPackage) != 0 {
        var (
            tempTransition VBoxCoreTransition
        )
        tempTransition, transErr = c.reporter.readMasterAck(c, metaPackage, timestamp)

        // finalize intermediate transition with failed count + timestamp
        transition = finalizeTransitionWithMeta(c, tempTransition, timestamp)

    } else {
        transition, transErr = stateTransitionWithTimestamp(c, timestamp)
    }

    // finalize core reporter state
    newState = stateTransition(oldState, transition)

    // execute event lisenter
    eventErr = runOnTransitionEvents(c, newState, oldState, transition, timestamp)

    // assign vbox reporter for new state
    c.reporter = newReporterForState(c.reporter, newState, oldState)

    return utils.SummarizeErrors(transErr, eventErr)
}