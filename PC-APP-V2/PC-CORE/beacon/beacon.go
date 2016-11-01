package beacon

import (
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-core/model"
)

type MasterBeaconState int
const (
    MasterInit              MasterBeaconState = iota
    MasterUnbounded
    MasterInquired
    MasterKeyExchange
    MasterCryptoCheck
    MasterBounded
    MasterBindBroken
    MasterDiscarded
)

type MasterBeaconTranstion int
const (
    MasterTransitionFail    MasterBeaconTranstion = iota
    MasterTransitionOk
    MasterTransitionIdle
)

func (st MasterBeaconState) String() string {
    var state string
    switch st {
        case MasterInit:
            state = "MasterInit"
        case MasterUnbounded:
            state = "MasterUnbounded"
        case MasterInquired:
            state = "MasterInquired"
        case MasterKeyExchange:
            state = "MasterKeyExchange"
        case MasterCryptoCheck:
            state = "MasterCryptoCheck"
        case MasterBounded:
            state = "MasterBounded"
        case MasterBindBroken:
            state = "MasterBindBroken"
        case MasterDiscarded:
            state = "MasterDiscarded"
    }
    return state
}

// MasterBeacon is assigned individually for each slave node.
type MasterBeacon interface {
    CurrentState() MasterBeaconState
    TransitionWithTimestamp(timestamp time.Time) error
    TransitionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error

    AESKey() ([]byte, error)
    AESCryptor() (crypt.AESCryptor, error)
    RSAEncryptor() (crypt.RsaEncryptor, error)

    SlaveNode() *model.SlaveNode
}

