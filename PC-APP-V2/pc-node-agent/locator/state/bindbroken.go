package state

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
)

type bindbroken struct{
    LocatorState
}

func (ls *bindbroken) executeTranslateMasterMetaWithTimestamp(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (locator.SlaveLocatingTransition, error) {
    if len(meta.EncryptedMasterRespond) == 0 {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect encrypted master respond")
    }
    if len(meta.EncryptedAESKey) == 0 {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect AES key from Master command")
    }
    if len(meta.RsaCryptoSignature) == 0 {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect RSA signature from Master command")
    }

    aeskey, err := slcontext.SharedSlaveContext().DecryptMessage(meta.EncryptedAESKey, meta.RsaCryptoSignature)
    if err != nil {
        return locator.SlaveTransitionFail, err
    }
    slcontext.SharedSlaveContext().SetAESKey(aeskey)

    // aes decryption of command
    pckedRsp, err := slcontext.SharedSlaveContext().Decrypt(meta.EncryptedMasterCommand)
    if err != nil {
        return locator.SlaveTransitionFail, err
    }
    msRsp, err := msagent.UnpackedMasterRespond(pckedRsp)
    if err != nil {
        return locator.SlaveTransitionFail, err
    }

    msAgent, err := slcontext.SharedSlaveContext().GetMasterAgent()
    if err != nil {
        return locator.SlaveTransitionFail, err
    }
    if msRsp.MasterBoundAgent != msAgent {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Master bound agent is different than commissioned one %s", msAgent)
    }
    if msRsp.Version != msagent.MASTER_RESPOND_VERSION {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }
    // if command is not for exchange key, just ignore
    if msRsp.MasterCommandType != msagent.COMMAND_RECOVER_BIND {
        return locator.SlaveTransitionIdle, nil
    }

    return locator.SlaveTransitionOk, nil
}

func (ls *bindbroken) executeStateTxWithTimestamp(slaveTimestamp time.Time) error {
    slctx := slcontext.SharedSlaveContext()

    masterAgentName, err := slctx.GetMasterAgent()
    if err != nil {
        return nil
    }
    ba, err := slagent.BrokenBindDiscovery(masterAgentName)
    if err != nil {
        return nil
    }
    _, err = slagent.BrokenBindMeta(ba)
    if err != nil {
        return nil
    }

    // TODO : send answer to master

    return nil
}
