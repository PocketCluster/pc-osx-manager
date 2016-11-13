package beacon

import (
    "fmt"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/context"
)

func boundedState(oldState *beaconState) MasterBeacon {
    b := &bounded{}

    b.constState                    = MasterInit

    b.constTransitionFailureLimit   = TransitionFailureLimit
    b.constTransitionTimeout        = UnboundedTimeout * time.Duration(TxActionLimit)
    b.constTxActionLimit            = TxActionLimit
    b.constTxTimeWindow             = UnboundedTimeout

    b.lastTransitionTS              = time.Now()

    b.timestampTransition           = b.transitionActionWithTimestamp
    b.slaveMetaTransition           = b.bounded
    b.onTransitionSuccess           = b.onStateTranstionSuccess
    b.onTransitionFailure           = b.onStateTranstionFailure

    b.aesKey                        = oldState.aesKey
    b.aesCryptor                    = oldState.aesCryptor
    b.rsaEncryptor                  = oldState.rsaEncryptor
    b.slaveNode                     = oldState.slaveNode
    b.commChan                      = oldState.commChan

    return b
}

type bounded struct {
    beaconState
}

func (b *bounded) transitionActionWithTimestamp() error {
    return nil
}

func (b *bounded) bounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
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

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (b *bounded) onStateTranstionSuccess(slaveTimestamp time.Time) error {
    return nil
}

func (b *bounded) onStateTranstionFailure(slaveTimestamp time.Time) error {
    return nil
}
