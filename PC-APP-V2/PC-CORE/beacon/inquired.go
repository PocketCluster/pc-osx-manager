package beacon

import (
    "net"
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/msagent"
)

func inquiredState(oldState *beaconState) BeaconState {
    b := &inquired{}

    b.constState                    = MasterInquired

    b.constTransitionFailureLimit   = TransitionFailureLimit
    b.constTransitionTimeout        = UnboundedTimeout * time.Duration(TxActionLimit)
    b.constTxActionLimit            = TxActionLimit
    b.constTxTimeWindow             = UnboundedTimeout

    b.lastTransitionTS              = time.Now()

    b.timestampTransition           = b.transitionActionWithTimestamp
    b.slaveMetaTransition           = b.inquired

    b.BeaconOnTransitionEvent       = oldState.BeaconOnTransitionEvent
    b.slaveNode                     = oldState.slaveNode
    b.commChan                      = oldState.commChan

    b.slaveLocation                 = nil
    b.slaveStatus                   = oldState.slaveStatus

    return b
}

type inquired struct {
    beaconState
}

func (b *inquired) transitionActionWithTimestamp(masterTimestamp time.Time) error {
    masterPubKey, err := context.SharedHostContext().MasterHostPublicKey()
    if err != nil {
        return errors.WithStack(err)
    }
    if b.slaveStatus == nil {
        return errors.Errorf("[ERR] SlaveStatusAgent is nil. We cannot form a proper response")
    }
    cmd, err := msagent.MasterDeclarationCommand(b.slaveStatus, masterTimestamp)
    if err != nil {
        return errors.WithStack(err)
    }
    meta, err := msagent.MasterDeclarationMeta(cmd, masterPubKey)
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

func (b *inquired) inquired(sender *net.UDPAddr, meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if sender == nil {
        return MasterTransitionIdle, errors.Errorf("[ERR] incorrect slave input. slave address should not be nil when checking identity.")
    }
    if meta.StatusAgent == nil || meta.StatusAgent.Version != slagent.SLAVE_STATUS_VERSION {
        return MasterTransitionFail, errors.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if meta.StatusAgent.SlaveResponse != slagent.SLAVE_SEND_PUBKEY {
        return MasterTransitionIdle, nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    if masterAgentName != meta.MasterBoundAgent {
        return MasterTransitionFail, errors.Errorf("[ERR] Slave reports to incorrect master agent")
    }
    // address check
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
    if b.slaveNode.Arch != meta.StatusAgent.SlaveHardware {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave architecture")
    }
    if len(meta.SlavePubKey) == 0 {
        return MasterTransitionFail, errors.Errorf("[ERR] Inappropriate slave public key")
    }

    // master public key
    masterPrvKey, err := context.SharedHostContext().MasterHostPrivateKey()
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    encryptor, err := pcrypto.NewRsaEncryptorFromKeyData(meta.SlavePubKey, masterPrvKey)
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    b.slaveNode.PublicKey = meta.SlavePubKey
    b.rsaEncryptor = encryptor

    // aeskey & aes encryptor/decryptor
    aesKey := pcrypto.NewAESKey32Byte()
    aesCryptor, err := pcrypto.NewAESCrypto(aesKey)
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    b.aesKey = aesKey
    b.aesCryptor = aesCryptor

    // now assign node name and
    err = b.slaveNode.SanitizeSlave()
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }

    // save status for response generation
    b.slaveStatus = meta.StatusAgent

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}
