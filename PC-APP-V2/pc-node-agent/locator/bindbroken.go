package locator

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
)

type bindbroken struct{
    locatorState
}

func (ls *bindbroken) CurrentState() SlaveLocatingState {
    return SlaveBindBroken
}

func (ls *bindbroken) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
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

    // TODO : broadcast slave meta

    return nil
}

func (ls *bindbroken) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        // if master is wrong version, It's perhaps from different master. we'll skip and wait for another time
        return SlaveTransitionIdle, fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }
    if len(meta.EncryptedMasterRespond) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect encrypted master respond")
    }
    if len(meta.EncryptedAESKey) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect AES key from Master command")
    }
    if len(meta.RsaCryptoSignature) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect RSA signature from Master command")
    }

    aeskey, err := slcontext.SharedSlaveContext().DecryptByRSA(meta.EncryptedAESKey, meta.RsaCryptoSignature)
    if err != nil {
        return SlaveTransitionFail, err
    }
    slcontext.SharedSlaveContext().SetAESKey(aeskey)

    // aes decryption of command
    pckedRsp, err := slcontext.SharedSlaveContext().DecryptByAES(meta.EncryptedMasterRespond)
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

    // set the master ip address
    if len(msRsp.MasterAddress) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect master address")
    }
    slcontext.SharedSlaveContext().SetMasterIP4Address(msRsp.MasterAddress)

    return SlaveTransitionOk, nil
}

func (ls *bindbroken) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return slcontext.SharedSlaveContext().SyncAll()
}

func (ls *bindbroken) onStateTranstionFailure(slaveTimestamp time.Time) error {
    slcontext.SharedSlaveContext().DiscardAESKey()
    return nil
}