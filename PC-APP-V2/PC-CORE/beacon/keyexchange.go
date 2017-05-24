package beacon

import (
    "net"
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/msagent"
)

func keyexchangeState(oldState *beaconState) BeaconState {
    b := &keyexchange{}

    b.constState                    = MasterKeyExchange

    b.constTransitionFailureLimit   = TransitionFailureLimit
    b.constTransitionTimeout        = UnboundedTimeout * time.Duration(TxActionLimit)
    b.constTxActionLimit            = TxActionLimit
    b.constTxTimeWindow             = UnboundedTimeout

    b.lastTransitionTS              = time.Now()

    b.timestampTransition           = b.transitionActionWithTimestamp
    b.slaveMetaTransition           = b.keyExchange
    b.onTransitionSuccess           = b.onStateTranstionSuccess
    b.onTransitionFailure           = b.onStateTranstionFailure

    b.aesKey                        = oldState.aesKey
    b.aesCryptor                    = oldState.aesCryptor
    b.rsaEncryptor                  = oldState.rsaEncryptor
    b.slaveNode                     = oldState.slaveNode
    b.commChan                      = oldState.commChan

    b.slaveLocation                 = nil
    b.slaveStatus                   = oldState.slaveStatus

    return b
}

type keyexchange struct {
    beaconState
}

func (b *keyexchange) transitionActionWithTimestamp(masterTimestamp time.Time) error {
    if b.slaveStatus == nil {
        return errors.Errorf("[ERR] SlaveStatusAgent is nil. We cannot form a proper response")
    }
    cmd, slvstat, err := msagent.ExchangeCryptoKeyAndNameCommand(b.slaveStatus, b.slaveNode.NodeName, b.slaveNode.SlaveUUID, masterTimestamp)
    if err != nil {
        return errors.WithStack(err)
    }
    meta, err := msagent.ExchangeCryptoKeyAndNameMeta(cmd, slvstat, b.aesKey, b.aesCryptor, b.rsaEncryptor)
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

func (b *keyexchange) keyExchange(sender *net.UDPAddr, meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if sender == nil {
        return MasterTransitionIdle, errors.Errorf("[ERR] incorrect slave input. slave address should not be nil when checking important credentials.")
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
    if usm.SlaveResponse != slagent.SLAVE_CHECK_CRYPTO {
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
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave node name beacon [%s] / slave master [%s] ", b.slaveNode.NodeName, usm.SlaveNodeName)
    }
    if b.slaveNode.SlaveUUID != usm.SlaveUUID {
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
    if b.slaveNode.MacAddress != meta.SlaveID {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave MAC address")
    }
    if b.slaveNode.Arch != usm.SlaveHardware {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave architecture")
    }

    // save status for response generation
    b.slaveStatus = usm

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (b *keyexchange) onStateTranstionSuccess(masterTimestamp time.Time) error {
    return nil
}

func (b *keyexchange) onStateTranstionFailure(masterTimestamp time.Time) error {
    return nil
}
