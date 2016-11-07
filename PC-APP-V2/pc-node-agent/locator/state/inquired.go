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

func (ls *inquired) executeTranslateMasterMetaWithTimestamp(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (locator.SlaveLocatingTransition, error) {
    // TODO : 1) check if meta is rightful to be bound

    if meta.StatusCommand == nil || meta.StatusCommand.Version != msagent.MASTER_COMMAND_VERSION {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of master command")
    }
    if meta.MasterPubkey == nil {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Malformed master command without public key")
    }
    if meta.StatusCommand.MasterCommandType != msagent.COMMAND_MASTER_DECLARE {
        return locator.SlaveTransitionIdle, nil
    }
    if err := slcontext.SharedSlaveContext().SetMasterAgent(meta.StatusCommand.MasterBoundAgent); err != nil {
        return locator.SlaveTransitionFail, err
    }
    if err := slcontext.SharedSlaveContext().SetMasterPublicKey(meta.MasterPubkey); err != nil {
        return locator.SlaveTransitionFail, err
    }

    return locator.SlaveTransitionOk, nil
}

func (ls *inquired) executeStateTxWithTimestamp(slaveTimestamp time.Time) error {
    agent, err := slagent.AnswerMasterInquiryStatus(slaveTimestamp)
    if err != nil {
        //return false
    }
    _, err = slagent.AnswerMasterInquiryMeta(agent)
    if err != nil {
        //return false
    }
    // TODO : send answer to master

    return nil
}
