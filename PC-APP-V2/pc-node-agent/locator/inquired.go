package locator

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
)

func newInquiredState(searchComm SearchTx, beaconComm BeaconTx) LocatorState {
    is := &inquired{}

    is.constState                   = SlaveInquired

    is.constTransitionFailureLimit  = TransitionFailureLimit
    is.constTransitionTimout        = UnboundedTimeout * time.Duration(TxActionLimit)
    is.constTxActionLimit           = TxActionLimit
    is.constTxTimeWindow            = UnboundedTimeout

    is.lastTransitionTS             = time.Now()

    is.timestampTransition          = is.transitionActionWithTimestamp
    is.masterMetaTransition         = is.transitionWithMasterMeta
    is.onTransitionSuccess          = is.onStateTranstionSuccess
    is.onTransitionFailure          = is.onStateTranstionFailure

    is.searchComm                   = searchComm
    is.beaconComm                   = beaconComm
    return is
}

type inquired struct {
    locatorState
}

func (ls *inquired) transitionActionWithTimestamp(slaveTimestamp time.Time) error {
    agent, err := slagent.AnswerMasterInquiryStatus(slaveTimestamp)
    if err != nil {
        return errors.WithStack(err)
    }
    sm, err := slagent.AnswerMasterInquiryMeta(agent)
    if err != nil {
        return errors.WithStack(err)
    }
    pm, err := slagent.PackedSlaveMeta(sm)
    if err != nil {
        return errors.WithStack(err)
    }
    ma, err := slcontext.SharedSlaveContext().GetMasterIP4Address()
    if err != nil {
        return errors.WithStack(err)
    }
    if ls.beaconComm == nil {
        return errors.Errorf("[ERR] Comm Channel is nil")
    }
    return ls.beaconComm.UcastSend(pm, ma)
}

func (ls *inquired) transitionWithMasterMeta(meta *msagent.PocketMasterAgentMeta, slaveTimestamp time.Time) (SlaveLocatingTransition, error) {
    // TODO : 1) check if meta is rightful to be bound
    if meta == nil || meta.MetaVersion != msagent.MASTER_META_VERSION {
        // if master is wrong version, It's perhaps from different master. we'll skip and wait for another time
        return SlaveTransitionIdle, errors.Errorf("[ERR] Null or incorrect version of master meta")
    }
    if meta.StatusCommand == nil || meta.StatusCommand.Version != msagent.MASTER_COMMAND_VERSION {
        return SlaveTransitionFail, errors.Errorf("[ERR] Null or incorrect version of master command")
    }
    if len(meta.MasterPubkey) == 0 {
        return SlaveTransitionFail, errors.Errorf("[ERR] Malformed master command without public key")
    }
    if meta.StatusCommand.MasterCommandType != msagent.COMMAND_MASTER_DECLARE {
        return SlaveTransitionIdle, nil
    }

    // set master agent name
    if err := slcontext.SharedSlaveContext().SetMasterAgent(meta.StatusCommand.MasterBoundAgent); err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }

    // set master public key
    if err := slcontext.SharedSlaveContext().SetMasterPublicKey(meta.MasterPubkey); err != nil {
        return SlaveTransitionFail, errors.WithStack(err)
    }

    return SlaveTransitionOk, nil
}

func (ls *inquired) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return nil
}

func (ls *inquired) onStateTranstionFailure(slaveTimestamp time.Time) error {
    return slcontext.SharedSlaveContext().DiscardAll()
}
