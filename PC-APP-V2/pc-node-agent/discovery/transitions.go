package discovery

import (
    "time"
    "fmt"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slcontext"
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
    slaveContext           slcontext.PocketSlaveContext
    lastSuccess            time.Time
    discoveryState         SDState
}

func NewSlaveDiscovery(context slcontext.PocketSlaveContext) (sd SlaveDiscovery) {
    sd = &slaveDiscovery{
        slaveContext: context,
        discoveryState:SlaveUnbounded,
    }
    return
}

func (sd *slaveDiscovery) CurrentState() SDState {
    return sd.discoveryState
}

func (sd *slaveDiscovery) TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) error {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }

    switch sd.discoveryState {
    case SlaveUnbounded:
        return sd.unbounded(meta, timestamp)

    case SlaveInquired:
        return sd.inquired(meta, timestamp)

    case SlaveKeyExchange:
        return sd.keyExchange(meta, timestamp)

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

// -- state evaluation

func (sd *slaveDiscovery) unbounded(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (err error) {
    if meta.DiscoveryRespond == nil || meta.DiscoveryRespond.Version != msagent.MASTER_RESPOND_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of master response")
    }
    // If command is incorrect, it should not be considered as an error and be ignored, although ignoring shouldn't happen.
    if meta.DiscoveryRespond.MasterCommandType != msagent.COMMAND_WHO_R_U {
        return nil
    }
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
    // TODO : 1) check if meta is rightful to be bound

    if meta.StatusCommand == nil || meta.StatusCommand.Version != msagent.MASTER_COMMAND_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of master command")
    }
    if meta.MasterPubkey == nil {
        return fmt.Errorf("[ERR] Malformed master command without public key")
    }
    if meta.StatusCommand.MasterCommandType != msagent.COMMAND_SEND_PUBKEY {
        return nil
    }
    if err = sd.slaveContext.SetMasterAgent(meta.StatusCommand.MasterBoundAgent); err != nil {
        return
    }
    if err = sd.slaveContext.SetMasterPublicKey(meta.MasterPubkey); err != nil {
        return
    }

    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) keyExchange(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (err error) {

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

    aeskey, err := sd.slaveContext.DecryptMessage(meta.EncryptedAESKey, meta.RsaCryptoSignature)
    if err != nil {
        return
    }
    sd.slaveContext.SetAESKey(aeskey)

    // aes decryption of command
    pckedCmd, err := sd.slaveContext.Decrypt(meta.EncryptedMasterCommand)
    if err != nil {
        return
    }
    msCmd, err := msagent.UnpackedMasterCommand(pckedCmd)
    if err != nil {
        return
    }
    msAgent, err := sd.slaveContext.GetMasterAgent()
    if err != nil {
        return
    }
    if msCmd.MasterBoundAgent != msAgent {
        return fmt.Errorf("[ERR] Master bound agent is different than current one %s", msAgent)
    }
    if msCmd.Version != msagent.MASTER_COMMAND_VERSION {
        return fmt.Errorf("[ERR] Incorrect version of master command")
    }
    // if command is not for exchange key, just ignore
    if msCmd.MasterCommandType != msagent.COMMAND_SEND_AES {
        return nil
    }
    nodeName, err := sd.slaveContext.Decrypt(meta.EncryptedSlaveStatus)
    if err != nil {
        return
    }
    sd.slaveContext.SetSlaveNodeName(string(nodeName))

    // let's make transition
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
