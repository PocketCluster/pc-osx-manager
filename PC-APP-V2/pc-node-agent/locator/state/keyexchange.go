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

func (ls *keyexchange) executeTranslateMasterMetaWithTimestamp(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (locator.SlaveLocatingTransition, error) {
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

    aeskey, err := slcontext.SharedSlaveContext().DecryptMessage(meta.EncryptedAESKey, meta.RsaCryptoSignature)
    if err != nil {
        return locator.SlaveTransitionFail, err
    }
    slcontext.SharedSlaveContext().SetAESKey(aeskey)

    // aes decryption of command
    pckedCmd, err := slcontext.SharedSlaveContext().Decrypt(meta.EncryptedMasterCommand)
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
    nodeName, err := slcontext.SharedSlaveContext().Decrypt(meta.EncryptedSlaveStatus)
    if err != nil {
        return locator.SlaveTransitionFail, err
    }
    slcontext.SharedSlaveContext().SetSlaveNodeName(string(nodeName))

    return locator.SlaveTransitionOk, nil
}

func (ls *keyexchange) executeStateTxWithTimestamp(slaveTimestamp time.Time) error {
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

    // TODO : send answer to master

    return nil
}
