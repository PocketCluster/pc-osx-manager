package beacon

import (
    "fmt"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
)

func stateTransition(currState MasterBeaconState, nextCondition MasterBeaconTranstion) (nextState MasterBeaconState, err error) {
    // successfully transition to the next
    if nextCondition == MasterTransitionOk {
        switch currState {
        case MasterUnbounded:
            nextState = MasterInquired
        case MasterInquired:
            nextState = MasterKeyExchange
        case MasterKeyExchange:
            nextState = MasterCryptoCheck
        case MasterCryptoCheck:
            nextState = MasterBounded
        case MasterBounded:
            nextState = MasterBounded
        case MasterBindBroken:
            nextState = MasterBounded
        default:
            err = fmt.Errorf("[ERR] 'nextCondition is true and hit default' cannot happen")
        }
    // failed to transit
    } else if nextCondition == MasterTransitionFail {
        switch currState {
        case MasterUnbounded:
            nextState = MasterUnbounded
        case MasterInquired:
            nextState = MasterUnbounded
        case MasterKeyExchange:
            nextState = MasterUnbounded
        case MasterCryptoCheck:
            nextState = MasterUnbounded
        case MasterBounded:
            nextState = MasterBounded
        case MasterBindBroken:
            nextState = MasterBindBroken
        default:
            err = fmt.Errorf("[ERR] 'nextCondition is true and hit default' cannot happen")
        }
    // idle
    } else  {
        nextState = currState
    }
    return
}

func NewBeaconForSlaveNode() MasterBeacon {
    return &masterBeacon{
        managmentState: MasterUnbounded,
    }
}

type masterBeacon struct {
    lastSuccess    time.Time
    managmentState MasterBeaconState
}

func (mb *masterBeacon) CurrentState() MasterBeaconState {
    return mb.managmentState
}

func (mb *masterBeacon) TranstionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    switch mb.managmentState {
    case MasterUnbounded:
        return mb.unbounded(meta, timestamp)
    case MasterInquired:
        return mb.inquired(meta, timestamp)
    case MasterKeyExchange:
        return mb.keyExchange(meta, timestamp)
    case MasterCryptoCheck:
        return mb.cryptoCheck(meta, timestamp)
    case MasterBounded:
        return mb.bounded(meta, timestamp)
    case MasterBindBroken:
        return mb.bindBroken(meta, timestamp)
    default:
        err = fmt.Errorf("[ERR] managmentState cannot default")
    }
    return
}

func (mm *masterBeacon) unbounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}

func (mm *masterBeacon) inquired(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}

func (mm *masterBeacon) keyExchange(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}

func (mm *masterBeacon) cryptoCheck(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}

func (mm *masterBeacon) bounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}

func (mm *masterBeacon) bindBroken(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}
