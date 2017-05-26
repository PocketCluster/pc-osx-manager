package locator

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/msagent"
)

type DebugCommChannel struct {
    LastMcastMessage []byte
    LastUcastMessage []byte
    LastUcastHost    string
    MCommCount       int
    UCommCount       int
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

    switch state {
        case SlaveUnbounded: {
            log.Infof("DebugEventReceiver.OnStateTranstionSuccess (%v) | [%v]", ts, state.String())
            return nil
        }
        case SlaveInquired: {
            log.Infof("DebugEventReceiver.OnStateTranstionSuccess (%v) | [%v]", ts, state.String())
            return nil
        }
        case SlaveKeyExchange: {
            log.Infof("DebugEventReceiver.OnStateTranstionSuccess (%v) | [%v]", ts, state.String())
            return nil
        }
        case SlaveCryptoCheck: {
            log.Infof("DebugEventReceiver.OnStateTranstionSuccess (%v) | [%v]", ts, state.String())
            return nil
        }
        case SlaveBounded: {
            log.Infof("DebugEventReceiver.OnStateTranstionSuccess (%v) | [%v]", ts, state.String())
            return nil
        }
        case SlaveBindBroken: {
            log.Infof("DebugEventReceiver.OnStateTranstionSuccess (%v) | [%v]", ts, state.String())
            return nil
        }
    }

    return nil
}

func (d *DebugEventReceiver) OnStateTranstionFailure(state SlaveLocatingState, ts time.Time) error {
    d.LastLocatorState = state
    d.LastSlaveTS = ts
    d.IsSuccess = false

    switch state {
        case SlaveUnbounded: {
            log.Infof("DebugEventReceiver.OnStateTranstionFailure (%v) | [%v]", ts, state.String())
            return nil
        }
        case SlaveInquired: {
            log.Infof("DebugEventReceiver.OnStateTranstionFailure (%v) | [%v]", ts, state.String())
            return nil
        }
        case SlaveKeyExchange: {
            log.Infof("DebugEventReceiver.OnStateTranstionFailure (%v) | [%v]", ts, state.String())
            return nil
        }
        case SlaveCryptoCheck: {
            log.Infof("DebugEventReceiver.OnStateTranstionFailure (%v) | [%v]", ts, state.String())
            return nil
        }
        case SlaveBounded: {
            log.Infof("DebugEventReceiver.OnStateTranstionFailure (%v) | [%v]", ts, state.String())
            return nil
        }
        case SlaveBindBroken: {
            log.Infof("DebugEventReceiver.OnStateTranstionFailure (%v) | [%v]", ts, state.String())
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
