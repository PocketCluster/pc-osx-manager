package beacon

import (
    "net"
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/msagent"
)

func boundedState(oldState *beaconState) BeaconState {
    b := &bounded{}

    b.constState                    = MasterBounded

    b.constTransitionFailureLimit   = TransitionFailureLimit
    b.constTransitionTimeout        = BoundedTimeout * time.Duration(TxActionLimit)
    b.constTxActionLimit            = TxActionLimit
    b.constTxTimeWindow             = BoundedTimeout

    b.lastTransitionTS              = time.Now()

    b.timestampTransition           = b.transitionActionWithTimestamp
    b.slaveMetaTransition           = b.bounded
    b.onTransitionSuccess           = b.onStateTranstionSuccess
    b.onTransitionFailure           = b.onStateTranstionFailure

    b.BeaconOnTransitionEvent       = oldState.BeaconOnTransitionEvent
    b.aesKey                        = oldState.aesKey
    b.aesCryptor                    = oldState.aesCryptor
    b.rsaEncryptor                  = oldState.rsaEncryptor
    b.slaveNode                     = oldState.slaveNode
    b.commChan                      = oldState.commChan

    b.slaveLocation                 = nil
    b.slaveStatus                   = oldState.slaveStatus

    return b
}

type bounded struct {
    beaconState
}

func (b *bounded) transitionActionWithTimestamp(masterTimestamp time.Time) error {
    if b.slaveStatus == nil {
        return errors.Errorf("[ERR] SlaveStatusAgent is nil. We cannot form a proper response")
    }
    cmd, err := msagent.BoundedSlaveAckCommand(b.slaveStatus, masterTimestamp)
    if err != nil {
        return errors.WithStack(err)
    }
    meta, err := msagent.BoundedSlaveAckMeta(cmd, b.aesCryptor)
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

func (b *bounded) bounded(sender *net.UDPAddr, meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if sender == nil {
        return MasterTransitionIdle, errors.Errorf("[ERR] incorrect slave input. slave address should not be nil when pertaining bind.")
    }
    if len(meta.EncryptedStatus) == 0 {
        return MasterTransitionFail, errors.Errorf("[ERR] Null encrypted slave status")
    }
    if b.aesCryptor == nil {
        return MasterTransitionFail, errors.Errorf("[ERR] AES Cryptor is null. This should not happen")
    }
    if b.aesKey == nil {
        return MasterTransitionFail, errors.Errorf("[ERR] AES Key is null. This should not happen")
    }
    plain, err := b.aesCryptor.DecryptByAES(meta.EncryptedStatus)
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    usm, err := slagent.UnpackedSlaveStatus(plain)
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    if usm == nil || usm.Version != slagent.SLAVE_STATUS_VERSION {
        return MasterTransitionFail, errors.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if usm.SlaveResponse != slagent.SLAVE_REPORT_STATUS {
        return MasterTransitionIdle, nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    if masterAgentName != meta.MasterBoundAgent {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect master agent name from slave")
    }
    if b.slaveNode.NodeName != usm.SlaveNodeName {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave node name")
    }
    if b.slaveNode.AuthToken != usm.SlaveAuthToken {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave UUID")
    }
    // check address
    addr, err := b.slaveNode.IP4AddrString()
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    if addr != sender.IP.String() {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave ip address")
    }
    if b.slaveNode.SlaveID != meta.SlaveID {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave MAC address")
    }
    if b.slaveNode.Hardware != usm.SlaveHardware {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave architecture")
    }

    // save status for response generation
    b.slaveStatus = usm

    // (2016-11-16) We'll reset TX action count to 0 and now so successful tx action can happen infinitely
    // We need to reset the counter here when correct slave meta comes in
    // It is b/c when succeeded in confirming with slave, we should be able to keep receiving slave meta

    b.txActionCount = 0

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (b *bounded) onStateTranstionSuccess(masterTimestamp time.Time) error {
    return nil
}

func (b *bounded) onStateTranstionFailure(masterTimestamp time.Time) error {
    return nil
}
