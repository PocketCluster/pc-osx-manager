package beacon

import (
    "fmt"
    "time"

    "github.com/stkim1/pc-core/model"
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
    b.onTransitionSuccess           = b.onStateTranstionSuccess
    b.onTransitionFailure           = b.onStateTranstionFailure

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
    masterPubKey, err := context.SharedHostContext().MasterPublicKey()
    if err != nil {
        return err
    }
    if b.slaveStatus == nil {
        return fmt.Errorf("[ERR] SlaveStatusAgent is nil. We cannot form a proper response")
    }
    cmd, err := msagent.MasterDeclarationCommand(b.slaveStatus, masterTimestamp)
    if err != nil {
        return err
    }
    meta := msagent.MasterDeclarationMeta(cmd, masterPubKey)
    pm, err := msagent.PackedMasterMeta(meta)
    if err != nil {
        return err
    }
    if b.commChan == nil {
        fmt.Errorf("[ERR] Communication channel is null. This should never happen")
    }
    return b.commChan.UcastSend(pm, b.slaveNode.IP4Address)
}

func (b *inquired) inquired(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if meta.StatusAgent == nil || meta.StatusAgent.Version != slagent.SLAVE_STATUS_VERSION {
        return MasterTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if meta.StatusAgent.SlaveResponse != slagent.SLAVE_SEND_PUBKEY {
        return MasterTransitionIdle, nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return MasterTransitionFail, err
    }
    if masterAgentName != meta.StatusAgent.MasterBoundAgent {
        return MasterTransitionFail, fmt.Errorf("[ERR] Slave reports to incorrect master agent")
    }
    if b.slaveNode.IP4Address != meta.StatusAgent.SlaveAddress {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if meta.SlaveID != meta.StatusAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if b.slaveNode.MacAddress != meta.StatusAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if b.slaveNode.Arch != meta.StatusAgent.SlaveHardware {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave architecture")
    }
    if len(meta.SlavePubKey) == 0 {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave public key")
    }

    // master public key
    masterPrvKey, err := context.SharedHostContext().MasterPrivateKey()
    if err != nil {
        return MasterTransitionFail, err
    }
    encryptor, err := pcrypto.NewRsaEncryptorFromKeyData(meta.SlavePubKey, masterPrvKey)
    if err != nil {
        return MasterTransitionFail, err
    }
    b.slaveNode.PublicKey = meta.SlavePubKey
    b.rsaEncryptor = encryptor

    // aeskey & aes encryptor/decryptor
    aesKey := pcrypto.NewAESKey32Byte()
    aesCryptor, err := pcrypto.NewAESCrypto(aesKey)
    if err != nil {
        return MasterTransitionFail, err
    }
    b.aesKey = aesKey
    b.aesCryptor = aesCryptor

    nodeName, err := model.FindSlaveNameCandiate()
    if err != nil {
        return MasterTransitionFail, err
    }
    b.slaveNode.NodeName = nodeName

    // save status for response generation
    b.slaveStatus = meta.StatusAgent

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (b *inquired) onStateTranstionSuccess(masterTimestamp time.Time) error {
    return nil
}

func (b *inquired) onStateTranstionFailure(masterTimestamp time.Time) error {
    return nil
}
