package discovery

import (
    "time"
    "github.com/stkim1/pc-core/msagent"
)

type SDState int
const (
    SlaveUnbounded         SDState = iota
    SlaveBounded
    SlaveBindBroken
    SlaveInquired
    SlaveKeyExchange
    SlaveCryptoCheck
    SlaveDiscarded
)

type SDTranstion int
const (
    SlaveTransitionFail    SDTranstion = iota
    SlaveTransitionOk
    SlaveTransitionIdle
)

func (st SDState) String() string {
    var state string
    switch st {
    case SlaveUnbounded:
        state = "SlaveUnbounded"
    case SlaveBounded:
        state = "SlaveBounded"
    case SlaveBindBroken:
        state = "SlaveBindBroken"
    case SlaveInquired:
        state = "SlaveInquired"
    case SlaveKeyExchange:
        state = "SlaveKeyExchange"
    case SlaveCryptoCheck:
        state = "SlaveCryptoCheck"
    case SlaveDiscarded:
        state = "SlaveDiscarded"
    }
    return state
}

type SlaveDiscovery interface {
    CurrentState() SDState
    TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta) (func (timestamp time.Time) (error))
}

func NewSlaveDiscovery() (sd SlaveDiscovery) {
    sd = &slaveDiscovery{
        discoveryState:SlaveUnbounded,
    }
    return
}

