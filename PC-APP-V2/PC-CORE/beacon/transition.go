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
        beaconState: MasterUnbounded,
    }
}

type masterBeacon struct {
    lastSuccess         time.Time
    beaconState         MasterBeaconState
}

func (mb *masterBeacon) CurrentState() MasterBeaconState {
    return mb.beaconState
}

func (mb *masterBeacon) TranstionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta == nil || meta.MetaVersion != slagent.SLAVE_META_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave meta")
    }
    switch mb.beaconState {
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
        return fmt.Errorf("[ERR] managmentState cannot reach default")
    }
}

func (mb *masterBeacon) unbounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta.DiscoveryAgent == nil || meta.DiscoveryAgent.Version != slagent.SLAVE_DISCOVER_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave discovery")
    }
    // if slave isn't looking for agent, then just return. this is not for this state.
    if meta.DiscoveryAgent.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return nil
    }
    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) inquired(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta.StatusAgent == nil || meta.StatusAgent.Version != slagent.SLAVE_STATUS_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if meta.StatusAgent.SlaveResponse != slagent.SLAVE_WHO_I_AM {
        return nil
    }

    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) keyExchange(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta.StatusAgent == nil || meta.StatusAgent.Version != slagent.SLAVE_STATUS_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if meta.StatusAgent.SlaveResponse != slagent.SLAVE_SEND_PUBKEY {
        return nil
    }
    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) cryptoCheck(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta.StatusAgent == nil || meta.StatusAgent.Version != slagent.SLAVE_STATUS_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if meta.StatusAgent.SlaveResponse != slagent.SLAVE_CHECK_CRYPTO {
        return nil
    }
    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) bounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta.StatusAgent == nil || meta.StatusAgent.Version != slagent.SLAVE_STATUS_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if meta.StatusAgent.SlaveResponse != slagent.SLAVE_REPORT_STATUS {
        return nil
    }
    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) bindBroken(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    return nil
}
