package state

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
)

type inquired struct {
    LocatorState
}

func (ls *inquired) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    agent, err := slagent.AnswerMasterInquiryStatus(slaveTimestamp)
    if err != nil {
        return err
    }
    _, err = slagent.AnswerMasterInquiryMeta(agent)
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

func (ls *inquired) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (locator.SlaveLocatingTransition, error) {
    // TODO : 1) check if meta is rightful to be bound
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        // if master is wrong version, It's perhaps from different master. we'll skip and wait for another time
        return locator.SlaveTransitionIdle, fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }
    if meta.StatusCommand == nil || meta.StatusCommand.Version != msagent.MASTER_COMMAND_VERSION {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of master command")
    }
    if len(meta.MasterPubkey) == 0 {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Malformed master command without public key")
    }
    if meta.StatusCommand.MasterCommandType != msagent.COMMAND_MASTER_DECLARE {
        return locator.SlaveTransitionIdle, nil
    }

    // set master agent name
    if err := slcontext.SharedSlaveContext().SetMasterAgent(meta.StatusCommand.MasterBoundAgent); err != nil {
        return locator.SlaveTransitionFail, err
    }

    // set master public key
    if err := slcontext.SharedSlaveContext().SetMasterPublicKey(meta.MasterPubkey); err != nil {
        return locator.SlaveTransitionFail, err
    }

    return locator.SlaveTransitionOk, nil
}

func (ls *inquired) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return nil
}

func (ls *inquired) onStateTranstionFailure(slaveTimestamp time.Time) error {
    return slcontext.SharedSlaveContext().DiscardAll()
}
