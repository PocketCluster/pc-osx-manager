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

func (sd *slaveDiscovery) TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) error {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }

    switch sd.discoveryState {
    case SlaveUnbounded: {
        if meta.DiscoveryRespond == nil || meta.DiscoveryRespond.Version != msagent.MASTER_RESPOND_VERSION {
            return fmt.Errorf("[ERR] Null or incorrect version of master response")
        }
        // If command is incorrect, it should not be considered as an error and be ignored, although ignoring shouldn't happen.
        if meta.DiscoveryRespond.MasterCommandType == msagent.COMMAND_WHO_R_U {
            return sd.unbounded(meta, timestamp)
        } else {
            return nil
        }
    }
    case SlaveInquired: {
        if meta.StatusCommand == nil || meta.StatusCommand.Version != msagent.MASTER_COMMAND_VERSION {
            return fmt.Errorf("[ERR] Null or incorrect version of master command")
        }
        if meta.MasterPubkey == nil {
            return fmt.Errorf("[ERR] Malformed master command without public key")
        }
        if meta.StatusCommand.MasterCommandType == msagent.COMMAND_SEND_PUBKEY {
            return sd.inquired(meta, timestamp)
        } else {
            return nil
        }
    }
    case SlaveKeyExchange: {
        if meta.StatusCommand == nil || meta.StatusCommand.Version != msagent.MASTER_COMMAND_VERSION {
            return fmt.Errorf("[ERR] Null or incorrect version of master command")
        }
        if len(meta.RsaCryptoSignature) == 0 {
            return fmt.Errorf("[ERR] Null or incorrect RSA signature from Master command")
        }
        if len(meta.EncryptedAESKey) == 0 {
            return fmt.Errorf("[ERR] Null or incorrect AES key from Master command")
        }
        if len(meta.EncryptedMasterCommand) == 0 {
            return fmt.Errorf("[ERR] Null or incorrect encrypted master command")
        }
        if len(meta.EncryptedSlaveStatus) == 0 {
            return fmt.Errorf("[ERR] Null or incorrect slave status from master command")
        }
/*
        if meta.StatusCommand.MasterCommandType == msagent.COMMAND_SEND_AES {
            return sd.inquired(meta, timestamp)
        } else {
            return nil
        }
*/
        return nil
    }

    case SlaveCryptoCheck: {
        if meta.StatusCommand == nil || meta.StatusCommand.Version != msagent.MASTER_COMMAND_VERSION {
            return fmt.Errorf("[ERR] Null or incorrect version of master command")
        }
        if len(meta.EncryptedMasterCommand) == 0 {
            return fmt.Errorf("[ERR] Null or incorrect encrypted master command")
        }
        return nil
    }
    case SlaveBounded: {
        if meta.StatusCommand == nil || meta.StatusCommand.Version != msagent.MASTER_COMMAND_VERSION {
            return fmt.Errorf("[ERR] Null or incorrect version of master command")
        }
        if len(meta.EncryptedMasterCommand) == 0 {
            return fmt.Errorf("[ERR] Null or incorrect encrypted master command")
        }
        return nil
    }
    case SlaveBindBroken: {

    }
    }
    return fmt.Errorf("[ERR] TranstionWithMasterMeta should never reach default")
}


func (sd *slaveDiscovery) unbounded(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) inquired(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (err error) {
    // TODO : 1) check if meta is rightful to be bound, 2) Save Master name, 3) Save master key
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) keyExchange(timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) cryptoCheck(timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) bounded(timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) bindBroken(timestamp time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}
