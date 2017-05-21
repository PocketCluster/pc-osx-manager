package locator

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/msagent"
)

type DebugCommChannel struct {
    LastMcastMessage []byte
    LastUcastMessage []byte
    LastUcastHost    string
    MCommCount       uint
    UCommCount       uint
}

func (dc *DebugCommChannel) McastSend(data []byte) error {
    dc.LastMcastMessage = data
    dc.MCommCount++
    return nil
}

func (dc *DebugCommChannel) UcastSend(target string, data []byte) error {
    dc.LastUcastMessage = data
    dc.LastUcastHost = target
    dc.UCommCount++
    return nil
}

// TODO : (2017-05-21) This interface should be deprecated for test only
func (sl *slaveLocator) TranstionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) error {
    if sl.state == nil {
        return errors.Errorf("[ERR] LocatorState is nil. Cannot make transition with master meta")
    }
    var err error
    sl.state, err = sl.state.MasterMetaTransition(meta, slaveTimestamp)
    return err
}
