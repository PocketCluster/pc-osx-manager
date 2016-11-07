package state

import (
    "time"

    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/slcontext"
)

type bounded struct{
    LocatorState
}

func (ls *bounded) txTimeout() time.Duration {
    return BoundedTimeout
}

func (ls *bounded) executeTranslateMasterMetaWithTimestamp(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (locator.SlaveLocatingTransition, error) {
    return locator.SlaveTransitionOk, nil
}

func (ls *bounded) executeStateTxWithTimestamp(slaveTimestamp time.Time) error {
    slctx := slcontext.SharedSlaveContext()

    masterAgentName, err := slctx.GetMasterAgent()
    if err != nil {
        return err
    }
    slaveAgentName, err := slctx.GetMasterAgent()
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

    // TODO : send answer to master
    return nil
}
