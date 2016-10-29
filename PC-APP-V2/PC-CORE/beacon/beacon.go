package beacon

import (
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
)

type MasterBeaconState int
const (
    MasterUnbounded         MasterBeaconState = iota
    MasterInquired
    MasterKeyExchange
    MasterCryptoCheck
    MasterBounded
    MasterBindBroken
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
    }
    return state
}

type MasterBeacon interface {
    CurrentState() MasterBeaconState
    TranstionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error)
}

