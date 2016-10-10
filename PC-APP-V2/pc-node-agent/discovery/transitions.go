package discovery

import (
    "time"
    "fmt"
    "github.com/stkim1/pc-core/msagent"
)

func stateTransition(currState SDState, nextCondition func() (SDTranstion)) (nextState SDState, err error) {
    var transtion SDTranstion = nextCondition()

    // Succeed to transition to the next
    if  transtion == SlaveTransitionOk {
        switch currState {
        case SlaveUnbounded:
            nextState = SlaveInquired
        case SlaveInquired:
            nextState = SlaveKeyExchange
        case SlaveKeyExchange:
            nextState = SlaveCryptoCheck
        case SlaveCryptoCheck:
            nextState = SlaveBounded
        case SlaveBindBroken:
            nextState = SlaveBounded
        case SlaveBounded:
            nextState = currState
        default:
            err = fmt.Errorf("[PANIC] 'nextCondition is true and hit default' cannot happen")
        }
        // Fail to transition to the next
    } else if transtion == SlaveTransitionFail {
        switch currState {
        case SlaveUnbounded:
            nextState = SlaveUnbounded
        case SlaveInquired:
            nextState = SlaveUnbounded
        case SlaveKeyExchange:
            nextState = SlaveUnbounded
        case SlaveCryptoCheck:
            nextState = SlaveUnbounded
        case SlaveBindBroken:
            nextState = SlaveBindBroken
        case SlaveBounded:
            nextState = SlaveBindBroken
        default:
            err = fmt.Errorf("[PANIC] 'nextCondition is false and hit default' cannot happen")
        }
        // Idle
    } else {
        nextState = currState
    }
    return
}

type slaveDiscovery struct {
    lastSuccess            time.Time
    discoveryState         SDState
}

func (sd *slaveDiscovery) CurrentState() SDState {
    return sd.discoveryState
}

func (sd *slaveDiscovery) TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta) (func (time.Time) error) {

    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        return func (time.Time) error {
            return fmt.Errorf("[ERR] Null or incorrect version of master meta")
        }
    }

    switch sd.discoveryState {
    case SlaveUnbounded: {
        if meta.DiscoveryRespond == nil || meta.DiscoveryRespond.Version != msagent.MASTER_DISCOVERY_VERSION {
            return func (time.Time) error {
                return fmt.Errorf("[ERR] Null or incorrect version of master discovery response %s", meta.MetaVersion)
            }
        }
        // If command is incorrect, it should not be considered as an error and be ignored, although ignoring shouldn't happen.
        if meta.DiscoveryRespond.MasterCommandType == msagent.COMMAND_WHO_R_U {
            return sd.unbounded
        } else {
            return func (time.Time) error {
                return nil
            }
        }
    }
    case SlaveInquired:
    case SlaveKeyExchange:
    case SlaveCryptoCheck:
    case SlaveBindBroken:
    case SlaveBounded:
    default:
    }



    return func (timestamp time.Time) error {

        return nil
    }
}




func (sd *slaveDiscovery) unbounded(timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) Inquired(timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) KeyExchange(timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) CryptoCheck(timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) Bounded(timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) BindBroken(timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}
