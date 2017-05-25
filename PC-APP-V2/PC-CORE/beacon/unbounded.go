package beacon

import (
    "net"
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/msagent"
)

func unboundedState(oldState *beaconState) BeaconState {
    b := &unbounded{}

    b.constState                    = MasterUnbounded

    b.constTransitionFailureLimit   = TransitionFailureLimit
    b.constTransitionTimeout        = UnboundedTimeout * time.Duration(TxActionLimit)
    b.constTxActionLimit            = TxActionLimit
    b.constTxTimeWindow             = UnboundedTimeout

    b.lastTransitionTS              = time.Now()

    b.timestampTransition           = b.transitionActionWithTimestamp
    b.slaveMetaTransition           = b.unbounded
    b.onTransitionSuccess           = b.onStateTranstionSuccess
    b.onTransitionFailure           = b.onStateTranstionFailure

    b.BeaconOnTransitionEvent       = oldState.BeaconOnTransitionEvent
    b.slaveNode                     = oldState.slaveNode
    b.commChan                      = oldState.commChan

    b.slaveLocation                 = oldState.slaveLocation
    b.slaveStatus                   = nil

    return b
}

type unbounded struct {
    beaconState
}

func (b *unbounded) transitionActionWithTimestamp(masterTimestamp time.Time) error {
    if b.slaveLocation == nil {
        return errors.Errorf("[ERR] SlaveDiscoveryAgent is nil. We cannot form a proper response %v", b.constState.String())
    }
    cmd, err := msagent.SlaveIdentityInqueryRespond(b.slaveLocation)
    if err != nil {
        return errors.WithStack(err)
    }
    meta, err := msagent.SlaveIdentityInquiryMeta(cmd)
    if err != nil {
        return errors.WithStack(err)
    }
    pm, err := msagent.PackedMasterMeta(meta)
    if err != nil {
        return errors.WithStack(err)
    }
    addr, err := b.slaveNode.IP4AddrString()
    if err != nil {
        return errors.WithStack(err)
    }
    if b.commChan == nil {
        errors.Errorf("[ERR] Communication channel is null. This should never happen")
    }
    return b.commChan.UcastSend(addr, pm)
}

func (b *unbounded) unbounded(sender *net.UDPAddr, meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if sender == nil {
        return MasterTransitionIdle, errors.Errorf("[ERR] incorrect slave input. slave address should not be nil when unbounded.")
    }
    if meta.StatusAgent == nil || meta.StatusAgent.Version != slagent.SLAVE_STATUS_VERSION {
        return MasterTransitionFail, errors.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if meta.StatusAgent.SlaveResponse != slagent.SLAVE_WHO_I_AM {
        return MasterTransitionIdle, nil
    }
    // check address
    addr, err := b.slaveNode.IP4AddrString()
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    if addr != sender.IP.String() {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave ip address")
    }
    if b.slaveNode  .MacAddress != meta.SlaveID {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave MAC address")
    }
    // slave hardware architecture
    if len(meta.StatusAgent.SlaveHardware) == 0 {
        return MasterTransitionFail, errors.Errorf("[ERR] Inappropriate slave architecture")
    }
    b.slaveNode.Arch = meta.StatusAgent.SlaveHardware

    // save status for response generation
    b.slaveStatus = meta.StatusAgent

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (b *unbounded) onStateTranstionSuccess(masterTimestamp time.Time) error {
    return nil
}

func (b *unbounded) onStateTranstionFailure(masterTimestamp time.Time) error {
    return nil
}
