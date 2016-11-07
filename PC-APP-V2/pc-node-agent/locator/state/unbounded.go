package state

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-node-agent/slagent"
)

type unbounded struct{
    LocatorState
}

func (ls *unbounded) executeTranslateMasterMetaWithTimestamp(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (locator.SlaveLocatingTransition, error) {
    if meta.DiscoveryRespond == nil || meta.DiscoveryRespond.Version != msagent.MASTER_RESPOND_VERSION {
        return locator.SlaveTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of master response")
    }
    // If command is incorrect, it should not be considered as an error and be ignored, although ignoring shouldn't happen.
    if meta.DiscoveryRespond.MasterCommandType != msagent.COMMAND_SLAVE_IDINQUERY {
        return locator.SlaveTransitionIdle, nil
    }

    return locator.SlaveTransitionOk, nil
}

func (ls *unbounded) executeStateTxWithTimestamp(slaveTimestamp time.Time) error {
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

