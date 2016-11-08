package state

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
)

type keyexchange struct{
    LocatorState
}

func (ls *keyexchange) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    slctx := slcontext.SharedSlaveContext()

    masterAgentName, err := slctx.GetMasterAgent()
    if err != nil {
        return err
    }
    agent, err := slagent.KeyExchangeStatus(masterAgentName, slaveTimestamp)
    if err != nil {
        return err
    }
    _, err = slagent.KeyExchangeMeta(agent, slctx.GetPublicKey())
    if err != nil {
        return err
    }
    _, err = slcontext.SharedSlaveContext().GetMasterIP4Address()
    if err != nil {
        return err
    }

    // TODO : send answer to master

    return nil
}

func (ls *keyexchange) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (locator.SlaveLocatingTransition, error) {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        // if master is wrong version, It's perhaps from different master. we'll skip and wait for another time
        return locator.SlaveTransitionIdle, fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }
    if len(meta.RsaCryptoSignature) == 0 {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect RSA signature from Master command")
    }
    if len(meta.EncryptedAESKey) == 0 {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect AES key from Master command")
    }
    if len(meta.EncryptedMasterCommand) == 0 {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect encrypted master command")
    }
    if len(meta.EncryptedSlaveStatus) == 0 {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect slave status from master command")
    }

    aeskey, err := slcontext.SharedSlaveContext().DecryptByRSA(meta.EncryptedAESKey, meta.RsaCryptoSignature)
    if err != nil {
        return locator.SlaveTransitionFail, err
    }
    slcontext.SharedSlaveContext().SetAESKey(aeskey)

    // aes decryption of command
    pckedCmd, err := slcontext.SharedSlaveContext().DecryptByAES(meta.EncryptedMasterCommand)
    if err != nil {
        return locator.SlaveTransitionFail, err
    }
    msCmd, err := msagent.UnpackedMasterCommand(pckedCmd)
    if err != nil {
        return locator.SlaveTransitionFail, err
    }

    msAgent, err := slcontext.SharedSlaveContext().GetMasterAgent()
    if err != nil {
        return locator.SlaveTransitionFail, err
    }
    if msCmd.MasterBoundAgent != msAgent {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Master bound agent is different than current one %s", msAgent)
    }
    if msCmd.Version != msagent.MASTER_COMMAND_VERSION {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Incorrect version of master command")
    }
    // if command is not for exchange key, just ignore
    if msCmd.MasterCommandType != msagent.COMMAND_EXCHANGE_CRPTKEY {
        return locator.SlaveTransitionIdle, nil
    }
    // set slave node name
    nodeName, err := slcontext.SharedSlaveContext().DecryptByAES(meta.EncryptedSlaveStatus)
    if err != nil {
        return locator.SlaveTransitionFail, err
    }
    slcontext.SharedSlaveContext().SetSlaveNodeName(string(nodeName))

    return locator.SlaveTransitionOk, nil
}

func (ls *keyexchange) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return nil
}

func (ls *keyexchange) onStateTranstionFailure(slaveTimestamp time.Time) error {
    return slcontext.SharedSlaveContext().DiscardAll()
}
