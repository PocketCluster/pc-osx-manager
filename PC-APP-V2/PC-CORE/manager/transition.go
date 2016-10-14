package manager

import (
    "fmt"
    "time"
    "github.com/stkim1/pc-node-agent/slagent"
)

func stateTransition(currState MMState, nextCondition MMTranstion) (nextState MMState, err error) {
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

func NewManagerForSlaveNode(slaveNode *slagent.PocketSlaveAgentMeta) (manager MasterManagement) {
    return &masterManagement{
        managmentState: MasterUnbounded,
    }
}

type masterManagement struct {
    lastSuccess       time.Time
    managmentState    MMState
}

func (mm *masterManagement) CurrentState() MMState {
    return mm.managmentState
}

func (mm *masterManagement) TranstionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    switch mm.managmentState {
    case MasterUnbounded:
        return mm.unbounded(meta, timestamp)
    case MasterInquired:
        return mm.inquired(meta, timestamp)
    case MasterKeyExchange:
        return mm.keyExchange(meta, timestamp)
    case MasterCryptoCheck:
        return mm.cryptoCheck(meta, timestamp)
    case MasterBounded:
        return mm.bounded(meta, timestamp)
    case MasterBindBroken:
        return mm.bindBroken(meta, timestamp)
    default:
        err = fmt.Errorf("[ERR] managmentState cannot default")
    }
    return
}

func (mm *masterManagement) unbounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}

func (mm *masterManagement) inquired(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}

func (mm *masterManagement) keyExchange(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}

func (mm *masterManagement) cryptoCheck(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}

func (mm *masterManagement) bounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}

func (mm *masterManagement) bindBroken(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (err error) {
    return
}
