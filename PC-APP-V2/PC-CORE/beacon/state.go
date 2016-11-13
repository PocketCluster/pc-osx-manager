package beacon

import (
    "time"
    "fmt"
    "log"
    "bytes"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/model"
)

const (
    TransitionFailureLimit      uint          = 5
    // TODO : timeout mechanism for receiving slave meta
    TransitionTimeout           time.Duration = time.Second * 10

    TxActionLimit               uint          = 5
    UnboundedTimeout            time.Duration = time.Second * 3
    BoundedTimeout              time.Duration = time.Second * 10
)

type transitionWithSlaveMeta          func (meta *slagent.PocketSlaveAgentMeta, masterTimestamp time.Time) (MasterBeaconTransition, error)

type transitionActionWithTimestamp    func (masterTimestamp time.Time) error

type onStateTranstionSuccess          func (masterTimestamp time.Time) error

type onStateTranstionFailure          func (masterTimestamp time.Time) error

type BeaconState interface {
    CurrentState() MasterBeaconState
    TransitionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, masterTimestamp time.Time) (BeaconState, error)
    TransitionWithTimestamp(masterTimestamp time.Time) (BeaconState, error)
    Close() error
}

type beaconState struct {
    /* -------------------------------------- given, constant field ------------------------------------------------- */
    constState                  MasterBeaconState

    constTransitionFailureLimit uint

    constTransitionTimeout      time.Duration

    constTxActionLimit          uint

    constTxTimeWindow           time.Duration

    /* ---------------------------------- changing properties to record transaction --------------------------------- */
    // each time we try to make transtion and fail, count goes up.
    transitionFailureCount      uint

    // last time successfully transitioned state
    lastTransitionTS            time.Time

    txActionCount               uint

    lastTransmissionTS          time.Time

    /* ----------------------------------------- transition functions ----------------------------------------------- */
    slaveMetaTransition         transitionWithSlaveMeta

    // timestamp transition func
    timestampTransition         transitionActionWithTimestamp

    // onSuccess
    onTransitionSuccess         onStateTranstionSuccess

    // onFailure
    onTransitionFailure         onStateTranstionFailure

    /* ---------------------------------------- all-states properties ----------------------------------------------- */
    slaveNode                   *model.SlaveNode
    aesKey                      []byte
    aesCryptor                  pcrypto.AESCryptor
    rsaEncryptor                pcrypto.RsaEncryptor
    commChan                    CommChannel
}

// properties
func (b *beaconState) CurrentState() MasterBeaconState {
    return b.constState
}

func (b *beaconState) transitionFailureLimit() uint {
    return b.constTransitionFailureLimit
}

func (b *beaconState) transitionTimeout() time.Duration {
    return b.constTransitionTimeout
}

func (b *beaconState) txActionLimit() uint {
    return b.constTxActionLimit
}

func (b *beaconState) txTimeWindow() time.Duration {
    return b.constTxTimeWindow
}

func (b *beaconState) AESKey() ([]byte, error) {
    if len(b.aesKey) == 0 {
        return nil, fmt.Errorf("[ERR] Empty AES Key")
    }
    return b.aesKey, nil
}
func (b *beaconState) AESCryptor() (pcrypto.AESCryptor, error) {
    if b.aesCryptor == nil {
        return nil, fmt.Errorf("[ERR] Null AES cryptor")
    }
    return b.aesCryptor, nil
}

func (b *beaconState) RSAEncryptor() (pcrypto.RsaEncryptor, error) {
    if b.rsaEncryptor == nil {
        return nil, fmt.Errorf("[ERR] Null RSA encryptor")
    }
    return b.rsaEncryptor, nil
}

func (b *beaconState) SlaveNode() (*model.SlaveNode) {
    // TODO : copy struct that the return value is read-only
    return b.slaveNode
}

/* ------------------------------------------------ Helper Functions ------------------------------------------------ */
// close func pointers and delegates to help GC
func (b *beaconState) Close() error {
    b.slaveMetaTransition    = nil
    b.timestampTransition    = nil
    b.onTransitionSuccess    = nil
    b.onTransitionFailure    = nil

    b.aesKey                 = nil
    b.aesCryptor             = nil
    b.rsaEncryptor           = nil
    b.slaveNode              = nil
    b.commChan               = nil

    return nil
}

/* ------------------------------------------ Meta Transition Functions --------------------------------------------- */
func stateTransition(currState MasterBeaconState, nextCondition MasterBeaconTransition) MasterBeaconState {
    var nextState MasterBeaconState
    // successfully transition to the next
    if nextCondition == MasterTransitionOk {
        switch currState {
        case MasterInit:
            nextState = MasterUnbounded
        case MasterUnbounded:
            nextState = MasterInquired
        case MasterInquired:
            nextState = MasterKeyExchange
        case MasterKeyExchange:
            nextState = MasterCryptoCheck

        case MasterCryptoCheck:
            fallthrough
        case MasterBounded:
            fallthrough
        case MasterBindBroken:
            nextState = MasterBounded
            break

        case MasterDiscarded:
            nextState = currState
        }
        // failed to transit
    } else if nextCondition == MasterTransitionFail {
        switch currState {

        case MasterInit:
            fallthrough
        case MasterUnbounded:
            fallthrough
        case MasterInquired:
            fallthrough
        case MasterKeyExchange:
            fallthrough
        case MasterCryptoCheck:
            nextState = MasterDiscarded
            break

        case MasterBounded:
            fallthrough
        case MasterBindBroken:
            nextState = MasterBindBroken
            break

        case MasterDiscarded:
            nextState = currState
        }
        // idle
    } else  {
        nextState = currState
    }
    return nextState
}

func (b *beaconState) translateStateWithTimeout(nextStateCandiate MasterBeaconTransition, masterTimestamp time.Time) MasterBeaconTransition {
    var nextConfirmedState MasterBeaconTransition

    switch nextStateCandiate {
        // As MasterTransitionOk does not check timewindow, it could grant an infinite timewindow to make transition.
        // This is indeed intented as it will give us a chance to handle racing situations. Plus, TransitionWithTimestamp()
        // should have squashed suspected beacons and that's the role of TransitionWithTimestamp()
        case MasterTransitionOk: {
            b.lastTransitionTS = masterTimestamp
            b.transitionFailureCount = 0
            nextConfirmedState = MasterTransitionOk
            break
        }
        default:{
            if b.transitionFailureCount < b.transitionFailureLimit() {
                b.transitionFailureCount++
            }

            if b.transitionFailureCount < b.transitionFailureLimit() && masterTimestamp.Sub(b.lastTransitionTS) < b.transitionTimeout() {
                nextConfirmedState = MasterTransitionIdle
            } else {
                nextConfirmedState = MasterTransitionFail
            }
        }
    }

    return nextConfirmedState
}

func runOnTransitionEvents(ls *beaconState, newState, oldState MasterBeaconState, transition MasterBeaconTransition, masterTimestamp time.Time) error {
    if newState != oldState {
        switch transition {
            case MasterTransitionOk:
                if ls.onTransitionSuccess != nil {
                    return ls.onTransitionSuccess(masterTimestamp)
                }
            case MasterTransitionFail: {
                if ls.onTransitionFailure != nil {
                    return ls.onTransitionFailure(masterTimestamp)
                }
            }
        }
    }
    return nil
}

func newBeaconForState(b* beaconState, newState, oldState MasterBeaconState) BeaconState {
    if newState == oldState {
        return b
    }

    var newBeaconState BeaconState = nil
    var err error = nil
    switch newState {
        case MasterInit:
            newBeaconState = beaconinitState(b.commChan)

        case MasterUnbounded:
            newBeaconState = unboundedState(b)

        case MasterInquired:
            newBeaconState = inquiredState(b)

        case MasterKeyExchange:
            newBeaconState = keyexchangeState(b)

        case MasterCryptoCheck:
            newBeaconState = cryptocheckState(b)

        case MasterBounded:
            newBeaconState = boundedState(b)

        case MasterBindBroken:
            newBeaconState, err = bindbrokenState(b.slaveNode, b.commChan)
            if err != nil {
                // this should never happen
                log.Panic(err.Error())
            }

        case MasterDiscarded:
            // TODO : implement master discarded
    }
    return newBeaconState
}

func (b *beaconState) TransitionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, masterTimestamp time.Time) (BeaconState, error) {
    var (
        newState, oldState MasterBeaconState = b.CurrentState(), b.CurrentState()
        transitionCandidate, finalTransition MasterBeaconTransition
        transErr, eventErr error = nil, nil
    )
    if b.slaveMetaTransition == nil {
        log.Panic("[CRITICAL] slaveMetaTransition func cannot be null")
    }

    if meta == nil || meta.MetaVersion != slagent.SLAVE_META_VERSION {
        return nil, fmt.Errorf("[ERR] Null or incorrect version of slave meta")
    }
    if len(meta.SlaveID) == 0 {
        return nil, fmt.Errorf("[ERR] Null or incorrect slave ID")
    }

    transitionCandidate, transErr = b.slaveMetaTransition(meta, masterTimestamp)

    // this is to apply failed time count and timeout window
    finalTransition = b.translateStateWithTimeout(transitionCandidate, masterTimestamp)

    // finalize master beacon state
    newState = stateTransition(oldState, finalTransition)

    // execute on events
    eventErr = runOnTransitionEvents(b, newState, oldState, finalTransition, masterTimestamp)

    return newBeaconForState(b, newState, oldState), summarizeErrors(transErr, eventErr)
}

/* ----------------------------------------- Timestamp Transition Functions ----------------------------------------- */
func (b *beaconState) TransitionWithTimestamp(masterTimestamp time.Time) (BeaconState, error) {
    var err error = nil
    var nextConfirmedState MasterBeaconTransition
    if b.transitionFailureCount < TransitionFailureLimit && masterTimestamp.Sub(b.lastTransitionTS) < TransitionTimeout {
        nextConfirmedState = MasterTransitionIdle
    } else {
        if TransitionFailureLimit <= b.transitionFailureCount {
            err = fmt.Errorf("[ERR] Transition has failed too many times already")
        } else if TransitionTimeout <= masterTimestamp.Sub(b.lastTransitionTS) {
            err = fmt.Errorf("[ERR] Slave did not make transition in the given time window " + TransitionTimeout.String())
        }
        nextConfirmedState = MasterTransitionFail
    }

    b.constState = stateTransition(b.constState, nextConfirmedState)
    return nil, err
}

/* ================================================= Operation Error ================================================ */
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
