package beacon

import (
    "fmt"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pcrypto"
)

func bindbrokenState(slaveNode *model.SlaveNode, comm CommChannel) (BeaconState, error) {
    b := &bindbroken{}

    b.constState                    = MasterBindBroken

    b.constTransitionFailureLimit   = TransitionFailureLimit
    b.constTransitionTimeout        = UnboundedTimeout * time.Duration(TxActionLimit)
    b.constTxActionLimit            = TxActionLimit
    b.constTxTimeWindow             = UnboundedTimeout

    b.lastTransitionTS              = time.Now()

    b.timestampTransition           = b.transitionActionWithTimestamp
    b.slaveMetaTransition           = b.bindBroken
    b.onTransitionSuccess           = b.onStateTranstionSuccess
    b.onTransitionFailure           = b.onStateTranstionFailure

    b.slaveNode                     = slaveNode
    b.commChan                      = comm

    b.slaveLocation                 = nil
    b.slaveStatus                   = nil

    // aeskey & aes encryptor/decryptor
    aesKey := pcrypto.NewAESKey32Byte()
    aesCryptor, err := pcrypto.NewAESCrypto(aesKey)
    if err != nil {
        b.Close()
        return nil, fmt.Errorf("[ERR] cannot create AES cyprtor " + err.Error())
    }
    b.aesKey = aesKey
    b.aesCryptor = aesCryptor

    // set RSA encryptor
    masterPrvKey, err := context.SharedHostContext().MasterHostPrivateKey()
    if err != nil {
        b.Close()
        return nil, err
    }
    if len(slaveNode.PublicKey) == 0 {
        b.Close()
        return nil, fmt.Errorf("[ERR] Cannot bind a slave without its public key. This only happens when user has deleted master database")
    }
    encryptor, err := pcrypto.NewRsaEncryptorFromKeyData(slaveNode.PublicKey, masterPrvKey)
    if err != nil {
        b.Close()
        return nil, err
    }
    b.rsaEncryptor = encryptor

    return b, nil
}

type bindbroken struct {
    beaconState
}

func (b *bindbroken) transitionActionWithTimestamp(masterTimestamp time.Time) error {
    // bindbroken state waits indifinitely for the rite slave meta comes in so that it can move on.
    // Here, we'll reset action count all the time.
    b.txActionCount = 0
    return nil
}

func (b *bindbroken) bindBroken(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if meta.DiscoveryAgent == nil || meta.DiscoveryAgent.Version != slagent.SLAVE_DISCOVER_VERSION {
        return MasterTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // if slave isn't looking for agent, then just return. this is not for this state.
    if meta.DiscoveryAgent.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return MasterTransitionIdle, nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return MasterTransitionFail, err
    }
    // since this node isn't looking for us, sliently ignore this request
    if masterAgentName != meta.DiscoveryAgent.MasterBoundAgent {
        return MasterTransitionIdle, nil
    }
    if b.slaveNode.IP4Address != meta.DiscoveryAgent.SlaveAddress {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if b.slaveNode.IP4Gateway != meta.DiscoveryAgent.SlaveGateway {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave gateway address")
    }
    if b.slaveNode.IP4Netmask != meta.DiscoveryAgent.SlaveNetmask {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave netmask address")
    }
    if meta.SlaveID != meta.DiscoveryAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if b.slaveNode.MacAddress != meta.DiscoveryAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave MAC address")
    }

    // save discovery agent for respond generation
    discovery, err := slagent.ConvertBindAttemptDiscoveryAgent(meta.DiscoveryAgent, b.slaveNode.NodeName, b.slaveNode.Arch)
    if err != nil {
        return MasterTransitionFail, err
    }
    b.slaveStatus = discovery

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (b *bindbroken) onStateTranstionSuccess(masterTimestamp time.Time) error {
    return nil
}

func (b *bindbroken) onStateTranstionFailure(masterTimestamp time.Time) error {
    return nil
}
