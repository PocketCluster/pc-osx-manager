package discovery

import (
    "time"
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

type SlaveDiscovery interface {
    Unbounded(timestamp *time.Time) (err error)
    Inquired(timestamp *time.Time) (err error)
    KeyExchange(timestamp *time.Time) (err error)
    CryptoCheck(timestamp *time.Time) (err error)
    Bounded(timestamp *time.Time) (err error)
    BindBroken(timestamp *time.Time) (err error)
}

type slaveDiscovery struct {
    *time.Time
    discoveryState         SDState
}

func NewSlaveDiscovery() (sd SlaveDiscovery) {
    sd = &slaveDiscovery{
        discoveryState:SlaveUnbounded,
    }
    return
}
