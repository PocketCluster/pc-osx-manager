package locator

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/slcontext"
)

func newBoundedState() LocatorState {
    bs := &bounded{}

    bs.constState                   = SlaveBounded

    bs.constTransitionFailureLimit  = TransitionFailureLimit
    bs.constTransitionTimout        = TransitionTimeout
    bs.constTxActionLimit           = TxActionLimit
    bs.constTxTimeWindow = BoundedTimeout

    bs.timestampTransition          = bs.transitionActionWithTimestamp
    bs.masterMetaTransition         = bs.transitionWithMasterMeta
    bs.onTransitionSuccess          = bs.onStateTranstionSuccess
    bs.onTransitionFailure          = bs.onStateTranstionFailure
    return bs
}

type bounded struct{
    locatorState
}

func (ls *bounded) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    // we'll reset TX action count to 0 and now so successful tx action can happen infinitely
    // we need to reset the counter here than receiver
    ls.txActionCount = 0

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

    // TODO : send answer to master

    return SlaveTransitionOk, nil
}

func (ls *bounded) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return nil
}

func (ls *bounded) onStateTranstionFailure(slaveTimestamp time.Time) error {
    slcontext.SharedSlaveContext().DiscardAESKey()
    return nil
}
