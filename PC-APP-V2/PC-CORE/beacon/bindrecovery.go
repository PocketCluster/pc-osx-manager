package beacon

import (
    "time"
    "fmt"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-core/context"
)

func bindrecoveryState(oldState *beaconState) BeaconState {
    b := &bindrecovery{}

    b.constState                    = MasterBindRecovery

    b.constTransitionFailureLimit   = TransitionFailureLimit
    b.constTransitionTimeout        = UnboundedTimeout * time.Duration(TxActionLimit)
    b.constTxActionLimit            = TxActionLimit
    b.constTxTimeWindow             = UnboundedTimeout

    b.lastTransitionTS              = time.Now()

    b.timestampTransition           = b.transitionActionWithTimestamp
    b.slaveMetaTransition           = b.transitionWithSlaveMeta
    b.onTransitionSuccess           = b.onStateTranstionSuccess
    b.onTransitionFailure           = b.onStateTranstionFailure

    b.aesKey                        = oldState.aesKey
    b.aesCryptor                    = oldState.aesCryptor
    b.rsaEncryptor                  = oldState.rsaEncryptor
    b.slaveNode                     = oldState.slaveNode
    b.commChan                      = oldState.commChan

    b.slaveLocation                 = nil
    // this status is generated from bindbroken
    b.slaveStatus                   = oldState.slaveStatus

    return b
}

type bindrecovery struct {
    beaconState
}

func (b *bindrecovery) transitionActionWithTimestamp(masterTimestamp time.Time) error {
    // master preperation
    if b.slaveLocation == nil {
        return fmt.Errorf("[ERR] SlaveDiscoveryAgent is nil. We cannot form a proper response")
    }
    cmd, err := msagent.BrokenBindRecoverRespond(b.slaveLocation)
    if err != nil {
        return err
    }
    // meta
    meta, err := msagent.BrokenBindRecoverMeta(cmd, b.aesKey, b.aesCryptor, b.rsaEncryptor)
    if err != nil {
        return err
    }
    pm, err := msagent.PackedMasterMeta(meta)
    if err != nil {
        return err
    }
    if b.commChan == nil {
        fmt.Errorf("[ERR] Communication channel is null. This should never happen")
    }
    return b.commChan.UcastSend(pm, b.slaveNode.IP4Address)
}

func (b *bindrecovery) transitionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, masterTimestamp time.Time) (MasterBeaconTransition, error) {
    if len(meta.EncryptedStatus) == 0 {
        return MasterTransitionFail, fmt.Errorf("[ERR] Null encrypted slave status")
    }
    if b.aesCryptor == nil {
        return MasterTransitionFail, fmt.Errorf("[ERR] AES Cryptor is null. This should not happen")
    }
    if b.aesKey == nil {
        return MasterTransitionFail, fmt.Errorf("[ERR] AES Key is null. This should not happen")
    }
    plain, err := b.aesCryptor.DecryptByAES(meta.EncryptedStatus)
    if err != nil {
        return MasterTransitionFail, err
    }
    usm, err := slagent.UnpackedSlaveStatus(plain)
    if err != nil {
        return MasterTransitionFail, err
    }
    if usm == nil || usm.Version != slagent.SLAVE_STATUS_VERSION {
        return MasterTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if usm.SlaveResponse != slagent.SLAVE_REPORT_STATUS {
        return MasterTransitionIdle, nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return MasterTransitionFail, err
    }
    if masterAgentName != usm.MasterBoundAgent {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect master agent name from slave")
    }
    if b.slaveNode.NodeName != usm.SlaveNodeName {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave master agent")
    }
    if b.slaveNode.IP4Address != usm.SlaveAddress {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if meta.SlaveID != usm.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if b.slaveNode.MacAddress != usm.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if b.slaveNode.Arch != usm.SlaveHardware {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave architecture")
    }

    // this status comes from slavenode. Save status for response generation
    b.slaveStatus = usm

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (b *bindrecovery) onStateTranstionSuccess(masterTimestamp time.Time) error {
    return nil
}

func (b *bindrecovery) onStateTranstionFailure(masterTimestamp time.Time) error {
    return nil
}

