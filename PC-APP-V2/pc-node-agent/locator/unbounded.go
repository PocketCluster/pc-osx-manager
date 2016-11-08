package locator

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/slcontext"
)

type unbounded struct{
    locatorState
}

func (ls *unbounded) CurrentState() SlaveLocatingState {
    return SlaveUnbounded
}

func (ls *unbounded) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    ua, err := slagent.UnboundedMasterDiscovery()
    if err != nil {
        return err
    }
    _, err = slagent.UnboundedMasterDiscoveryMeta(ua)
    if err != nil {
        return err
    }

    // TODO : broadcast slave meta

    return nil
}

func (ls *unbounded) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        // if master is wrong version, It's perhaps from different master. we'll skip and wait for another time
        return SlaveTransitionIdle, fmt.Errorf("[ERR] Null or incorrect version of master meta")
    }
    if meta.DiscoveryRespond == nil || meta.DiscoveryRespond.Version != msagent.MASTER_RESPOND_VERSION {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of master response")
    }
    // If command is incorrect, it should not be considered as an error and be ignored, although ignoring shouldn't happen.
    if meta.DiscoveryRespond.MasterCommandType != msagent.COMMAND_SLAVE_IDINQUERY {
        return SlaveTransitionIdle, nil
    }
    // set the master ip address
    if len(meta.DiscoveryRespond.MasterAddress) == 0 {
        return SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect master address")
    }
    slcontext.SharedSlaveContext().SetMasterIP4Address(meta.DiscoveryRespond.MasterAddress)

    return SlaveTransitionOk, nil
}

func (ls *unbounded) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    // nothing to do for unbounded -> inquired state failure
    return nil
}

func (ls *unbounded) onStateTranstionFailure(slaveTimestamp time.Time) error {
    return slcontext.SharedSlaveContext().DiscardAll()
}
