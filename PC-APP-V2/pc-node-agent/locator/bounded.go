package locator

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/slcontext"
)

type bounded struct{
    locatorState
}

func (ls *bounded) CurrentState() SlaveLocatingState {
    return SlaveBounded
}

func (ls *bounded) txTimeout() time.Duration {
    return BoundedTimeout
}

func (ls *bounded) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
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
    sa, err := slagent.SlaveBoundedStatus(masterAgentName, slaveAgentName, slaveTimestamp)
    if err != nil {
        return err
    }
    _, err = slagent.SlaveBoundedMeta(sa, aesCryptor)
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

func (ls *bounded) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        // if master is wrong version, It's perhaps from different master. we'll skip and wait for another time
        return SlaveTransitionIdle, fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }
    return SlaveTransitionOk, nil
}

func (ls *bounded) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return nil
}

func (ls *bounded) onStateTranstionFailure(slaveTimestamp time.Time) error {
    slcontext.SharedSlaveContext().DiscardAESKey()
    return nil
}
