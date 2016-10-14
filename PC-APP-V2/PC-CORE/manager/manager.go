package manager

import (
    "github.com/stkim1/pc-node-agent/slagent"
    "time"
)

type MMState int
const (
    MasterUnbounded         MMState = iota
    MasterInquired
    MasterKeyExchange
    MasterCryptoCheck
    MasterBounded
    MasterBindBroken
)

type MMTranstion int
const (
    MasterTransitionFail    MMTranstion = iota
    MasterTransitionOk
    MasterTransitionIdle
)

func (st MMState) String() string {
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

type MasterManagement interface {
    CurrentState() MMState
    TranstionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error)
}

