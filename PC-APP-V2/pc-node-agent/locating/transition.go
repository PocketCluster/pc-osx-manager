package locating

import (
    "time"
    "fmt"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slcontext"
)

const allowedTimesOfFailure int = 5
const timeOutWindow time.Duration = time.Second * 5

type slaveDiscovery struct {
    // last time successfully transitioned state
    lastSuccess      time.Time
    // each time we try to make transtion and fail, count goes up.
    trialFailCount   int
    locatingState    SlaveLocatingState
}

func NewSlaveDiscovery() (sd SlaveDiscovery) {
    sd = &slaveDiscovery{
        locatingState   : SlaveUnbounded,
        trialFailCount  : 0,
    }
    return
}

func (sd *slaveDiscovery) CurrentState() SlaveLocatingState {
    return sd.locatingState
}

func (sd *slaveDiscovery) TranstionWithTimestamp(timestamp time.Time) error {

    return nil
}

func (sd *slaveDiscovery) TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) error {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }

    var transition SlaveLocatingTransition
    var err error = nil

    switch sd.locatingState {
    case SlaveUnbounded:
        transition, err = sd.unbounded(meta, slaveTimestamp)

    case SlaveInquired:
        transition, err = sd.inquired(meta, slaveTimestamp)

    case SlaveKeyExchange:
        transition, err = sd.keyExchange(meta, slaveTimestamp)

    case SlaveCryptoCheck:
        transition, err = sd.cryptoCheck(meta, slaveTimestamp)

    case SlaveBounded:
        transition, err = sd.bounded(meta, slaveTimestamp)

    case SlaveBindBroken:
        transition, err = sd.bindBroken(meta, slaveTimestamp)

    default:
        transition, err = SlaveTransitionFail, fmt.Errorf("[ERR] TranstionWithMasterMeta should never reach default")
    }
    // filter out the intermediate transition value with failed count + timestamp
    finalTransitionCandidate := sd.translateStateWithTimeout(transition, slaveTimestamp)

    // finalize locating master beacon state
    sd.locatingState = stateTransition(sd.locatingState, finalTransitionCandidate)

    return err
}

func (sd *slaveDiscovery) translateStateWithTimeout(nextStateCandiate SlaveLocatingTransition, slaveTimestamp time.Time) SlaveLocatingTransition {

    var nextConfirmedState SlaveLocatingTransition
    switch nextStateCandiate {
    case SlaveTransitionOk: {
        sd.lastSuccess = slaveTimestamp
        nextConfirmedState = SlaveTransitionOk
    }
    default: {
        if sd.trialFailCount < allowedTimesOfFailure {
            sd.trialFailCount++
        }

        if sd.trialFailCount < allowedTimesOfFailure && slaveTimestamp.Sub(sd.lastSuccess) < timeOutWindow {
            nextConfirmedState = SlaveTransitionIdle
        } else {
            nextConfirmedState = SlaveTransitionFail
        }
    }
    }
    return nextConfirmedState
}

// -- state evaluation

func (sd *slaveDiscovery) unbounded(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (SlaveLocatingTransition, error) {
    if meta.DiscoveryRespond == nil || meta.DiscoveryRespond.Version != msagent.MASTER_RESPOND_VERSION {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of master response")
    }
    // If command is incorrect, it should not be considered as an error and be ignored, although ignoring shouldn't happen.
    if meta.DiscoveryRespond.MasterCommandType != msagent.COMMAND_SLAVE_IDINQUERY {
        return SlaveTransitionIdle, nil
    }

    return SlaveTransitionOk, nil
}

func (sd *slaveDiscovery) inquired(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (SlaveLocatingTransition, error) {
    // TODO : 1) check if meta is rightful to be bound

    if meta.StatusCommand == nil || meta.StatusCommand.Version != msagent.MASTER_COMMAND_VERSION {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of master command")
    }
    if meta.MasterPubkey == nil {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Malformed master command without public key")
    }
    if meta.StatusCommand.MasterCommandType != msagent.COMMAND_MASTER_DECLARE {
        return SlaveTransitionIdle, nil
    }
    if err := slcontext.SharedSlaveContext().SetMasterAgent(meta.StatusCommand.MasterBoundAgent); err != nil {
        return SlaveTransitionFail, err
    }
    if err := slcontext.SharedSlaveContext().SetMasterPublicKey(meta.MasterPubkey); err != nil {
        return SlaveTransitionFail, err
    }

    return SlaveTransitionOk, nil
}

func (sd *slaveDiscovery) keyExchange(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (SlaveLocatingTransition, error) {

    if len(meta.RsaCryptoSignature) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect RSA signature from Master command")
    }
    if len(meta.EncryptedAESKey) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect AES key from Master command")
    }
    if len(meta.EncryptedMasterCommand) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect encrypted master command")
    }
    if len(meta.EncryptedSlaveStatus) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect slave status from master command")
    }

    aeskey, err := slcontext.SharedSlaveContext().DecryptMessage(meta.EncryptedAESKey, meta.RsaCryptoSignature)
    if err != nil {
        return SlaveTransitionFail, err
    }
    slcontext.SharedSlaveContext().SetAESKey(aeskey)

    // aes decryption of command
    pckedCmd, err := slcontext.SharedSlaveContext().Decrypt(meta.EncryptedMasterCommand)
    if err != nil {
        return SlaveTransitionFail, err
    }
    msCmd, err := msagent.UnpackedMasterCommand(pckedCmd)
    if err != nil {
        return SlaveTransitionFail, err
    }
    msAgent, err := slcontext.SharedSlaveContext().GetMasterAgent()
    if err != nil {
        return SlaveTransitionFail, err
    }
    if msCmd.MasterBoundAgent != msAgent {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Master bound agent is different than current one %s", msAgent)
    }
    if msCmd.Version != msagent.MASTER_COMMAND_VERSION {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Incorrect version of master command")
    }
    // if command is not for exchange key, just ignore
    if msCmd.MasterCommandType != msagent.COMMAND_EXCHANGE_CRPTKEY {
        return SlaveTransitionIdle, nil
    }
    nodeName, err := slcontext.SharedSlaveContext().Decrypt(meta.EncryptedSlaveStatus)
    if err != nil {
        return SlaveTransitionFail, err
    }
    slcontext.SharedSlaveContext().SetSlaveNodeName(string(nodeName))

    return SlaveTransitionOk, nil
}

func (sd *slaveDiscovery) cryptoCheck(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (SlaveLocatingTransition, error) {
    if len(meta.EncryptedMasterCommand) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect encrypted master command")
    }
    // aes decryption of command
    pckedCmd, err := slcontext.SharedSlaveContext().Decrypt(meta.EncryptedMasterCommand)
    if err != nil {
        return SlaveTransitionFail, err
    }
    msCmd, err := msagent.UnpackedMasterCommand(pckedCmd)
    if err != nil {
        return SlaveTransitionFail, err
    }
    msAgent, err := slcontext.SharedSlaveContext().GetMasterAgent()
    if err != nil {
        return SlaveTransitionFail, err
    }
    if msCmd.MasterBoundAgent != msAgent {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Master bound agent is different than commissioned one %s", msAgent)
    }
    if msCmd.Version != msagent.MASTER_COMMAND_VERSION {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Incorrect version of master command")
    }
    // if command is not for exchange key, just ignore
    if msCmd.MasterCommandType != msagent.COMMAND_MASTER_BIND_READY {
        return SlaveTransitionIdle, nil
    }

    return SlaveTransitionOk, nil
}

func (sd *slaveDiscovery) bounded(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (SlaveLocatingTransition, error) {

    return SlaveTransitionOk, nil
}

func (sd *slaveDiscovery) bindBroken(meta *msagent.PocketMasterAgentMeta, timestamp time.Time) (SlaveLocatingTransition, error) {
    if len(meta.EncryptedMasterRespond) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect encrypted master respond")
    }
    if len(meta.EncryptedAESKey) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect AES key from Master command")
    }
    if len(meta.RsaCryptoSignature) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect RSA signature from Master command")
    }

    aeskey, err := slcontext.SharedSlaveContext().DecryptMessage(meta.EncryptedAESKey, meta.RsaCryptoSignature)
    if err != nil {
        return SlaveTransitionFail, err
    }
    slcontext.SharedSlaveContext().SetAESKey(aeskey)

    // aes decryption of command
    pckedRsp, err := slcontext.SharedSlaveContext().Decrypt(meta.EncryptedMasterCommand)
    if err != nil {
        return SlaveTransitionFail, err
    }
    msRsp, err := msagent.UnpackedMasterRespond(pckedRsp)
    if err != nil {
        return SlaveTransitionFail, err
    }

    msAgent, err := slcontext.SharedSlaveContext().GetMasterAgent()
    if err != nil {
        return SlaveTransitionFail, err
    }
    if msRsp.MasterBoundAgent != msAgent {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Master bound agent is different than commissioned one %s", msAgent)
    }
    if msRsp.Version != msagent.MASTER_RESPOND_VERSION {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }
    // if command is not for exchange key, just ignore
    if msRsp.MasterCommandType != msagent.COMMAND_RECOVER_BIND {
        return SlaveTransitionIdle, nil
    }

    return SlaveTransitionOk, nil
}


func stateTransition(currState SlaveLocatingState, nextCondition SlaveLocatingTransition) SlaveLocatingState {
    var nextState SlaveLocatingState
    // Succeed to transition to the next
    if  nextCondition == SlaveTransitionOk {
        switch currState {
        case SlaveUnbounded:
            nextState = SlaveInquired
        case SlaveInquired:
            nextState = SlaveKeyExchange
        case SlaveKeyExchange:
            nextState = SlaveCryptoCheck

        case SlaveCryptoCheck:
            fallthrough
        case SlaveBindBroken:
            fallthrough
        case SlaveBounded:
            nextState = SlaveBounded
            break

        default:
            nextState = SlaveUnbounded
        }
        // Fail to transition to the next
    } else if nextCondition == SlaveTransitionFail {
        switch currState {

        case SlaveUnbounded:
            fallthrough
        case SlaveInquired:
            fallthrough
        case SlaveKeyExchange:
            fallthrough
        case SlaveCryptoCheck:
            nextState = SlaveUnbounded
            break

        case SlaveBindBroken:
            fallthrough
        case SlaveBounded:
            nextState = SlaveBindBroken
            break

        default:
            nextState = SlaveUnbounded
        }
        // Idle
    } else {
        nextState = currState
    }
    return nextState
}