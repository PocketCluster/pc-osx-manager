package locator

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
)

func newCryptocheckState() LocatorState {
    cc := &cryptocheck{}
    cc.timestampTransition     = cc.transitionActionWithTimestamp
    cc.masterMetaTransition    = cc.transitionWithMasterMeta
    cc.onTransitionSuccess     = cc.onStateTranstionSuccess
    cc.onTransitionFailure     = cc.onStateTranstionFailure
    return cc
}

type cryptocheck struct{
    locatorState
}

func (ls *cryptocheck) CurrentState() SlaveLocatingState {
    return SlaveCryptoCheck
}

func (ls *cryptocheck) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    slctx := slcontext.SharedSlaveContext()

    masterAgentName, err := slctx.GetMasterAgent()
    if err != nil {
        return err
    }
    slaveAgentName, err := slctx.GetSlaveNodeName()
    if err != nil {
        return err
    }
    aesCryptor, err := slctx.AESCryptor()
    if err != nil {
        return err
    }
    sa, err := slagent.CheckSlaveCryptoStatus(masterAgentName, slaveAgentName, slaveTimestamp)
    if err != nil {
        return err
    }
    _, err = slagent.CheckSlaveCryptoMeta(sa, aesCryptor)
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

func (ls *cryptocheck) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        // if master is wrong version, It's perhaps from different master. we'll skip and wait for another time
        return SlaveTransitionIdle, fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }
    if len(meta.EncryptedMasterCommand) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect encrypted master command")
    }
    // aes decryption of command
    pckedCmd, err := slcontext.SharedSlaveContext().DecryptByAES(meta.EncryptedMasterCommand)
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

func (ls *cryptocheck) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return slcontext.SharedSlaveContext().SyncAll()
}

func (ls *cryptocheck) onStateTranstionFailure(slaveTimestamp time.Time) error {
    return slcontext.SharedSlaveContext().DiscardAll()
}
