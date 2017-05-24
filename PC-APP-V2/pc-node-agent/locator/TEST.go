package locator

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slcontext"
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

type DebugEventReceiver struct {
    LastLocatorState    SlaveLocatingState
    LastSlaveTS         time.Time
    IsSuccess           bool
}

func (d *DebugEventReceiver) OnStateTranstionSuccess(state SlaveLocatingState, ts time.Time) error {
    d.LastLocatorState = state
    d.LastSlaveTS = ts
    d.IsSuccess = true

    log.Infof("DebugEventReceiver.OnStateTranstionSuccess %v", state.String())

    switch state {

        case SlaveUnbounded: {
            // nothing to do for unbounded -> inquired state failure
            return nil
        }
        case SlaveInquired: {
            return nil
        }
        case SlaveKeyExchange: {
            return nil
        }
        case SlaveCryptoCheck: {
            // here we'll save all the detail and save it to disk
            err := slcontext.SharedSlaveContext().SyncAll()
            if err != nil {
                return errors.WithStack(err)
            }
            err = slcontext.SharedSlaveContext().SaveConfiguration()
            return errors.WithStack(err)
        }
        case SlaveBounded: {
            return nil
        }

        case SlaveBindBroken: {
            return slcontext.SharedSlaveContext().SyncAll()
        }
    }

    return nil
}

func (d *DebugEventReceiver) OnStateTranstionFailure(state SlaveLocatingState, ts time.Time) error {
    d.LastLocatorState = state
    d.LastSlaveTS = ts
    d.IsSuccess = false

    log.Infof("DebugEventReceiver.OnStateTranstionFailure %v", state.String())

    switch state {

        case SlaveUnbounded: {
            return slcontext.SharedSlaveContext().DiscardAll()
        }
        case SlaveInquired: {
            return slcontext.SharedSlaveContext().DiscardAll()
        }
        case SlaveKeyExchange: {
            return slcontext.SharedSlaveContext().DiscardAll()
        }
        case SlaveCryptoCheck: {
            return slcontext.SharedSlaveContext().DiscardAll()
        }
        case SlaveBounded: {
            slcontext.SharedSlaveContext().DiscardAESKey()
            return nil
        }
        case SlaveBindBroken: {
            slcontext.SharedSlaveContext().DiscardAESKey()
            return nil
        }
    }
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
