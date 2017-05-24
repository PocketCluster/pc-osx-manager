package beacon

import (
    "net"
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pcrypto"
)

func bindbrokenState(slaveNode *model.SlaveNode, comm CommChannel, event BeaconOnTransitionEvent) (BeaconState, error) {
    b := &bindbroken{}

    b.constState                    = MasterBindBroken

    b.constTransitionFailureLimit   = TransitionFailureLimit
    b.constTransitionTimeout        = UnboundedTimeout * time.Duration(TxActionLimit)
    b.constTxActionLimit            = TxActionLimit
    b.constTxTimeWindow             = UnboundedTimeout

    b.lastTransitionTS              = time.Now()

    b.timestampTransition           = b.transitionActionWithTimestamp
    b.slaveMetaTransition           = b.bindBroken

    b.BeaconOnTransitionEvent       = event
    b.slaveNode                     = slaveNode
    b.commChan                      = comm

    b.slaveLocation                 = nil
    b.slaveStatus                   = nil

    // aeskey & aes encryptor/decryptor
    aesKey := pcrypto.NewAESKey32Byte()
    aesCryptor, err := pcrypto.NewAESCrypto(aesKey)
    if err != nil {
        b.Close()
        return nil, errors.Errorf("[ERR] cannot create AES cyprtor " + err.Error())
    }
    b.aesKey = aesKey
    b.aesCryptor = aesCryptor

    // set RSA encryptor
    masterPrvKey, err := context.SharedHostContext().MasterHostPrivateKey()
    if err != nil {
        b.Close()
        return nil, errors.WithStack(err)
    }
    if len(slaveNode.PublicKey) == 0 {
        b.Close()
        return nil, errors.Errorf("[ERR] Cannot bind a slave without its public key. This only happens when user has deleted master database")
    }
    encryptor, err := pcrypto.NewRsaEncryptorFromKeyData(slaveNode.PublicKey, masterPrvKey)
    if err != nil {
        b.Close()
        return nil, errors.WithStack(err)
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

func (b *bindbroken) bindBroken(sender *net.UDPAddr, meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if sender == nil {
        return MasterTransitionIdle, errors.Errorf("[ERR] incorrect slave input. slave address should not be nil when receiving multicast in broken bind.")
    }
    if meta.DiscoveryAgent == nil || meta.DiscoveryAgent.Version != slagent.SLAVE_DISCOVER_VERSION {
        return MasterTransitionFail, errors.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // if slave isn't looking for agent, then just return. this is not for this state.
    if meta.DiscoveryAgent.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return MasterTransitionIdle, nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    // since this node isn't looking for us, sliently ignore this request
    if masterAgentName != meta.MasterBoundAgent {
        return MasterTransitionIdle, nil
    }

    // check mac address
    if b.slaveNode.MacAddress != meta.SlaveID {
        return MasterTransitionFail, errors.Errorf("[ERR] Incorrect slave MAC address")
    }
    // slave ip address
    // TODO : (2015-05-16) we're not checking ip + subnet eligivility for now
    addr, err := model.IP4AddrToString(meta.DiscoveryAgent.SlaveAddress)
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    if addr != sender.IP.String() {
        return MasterTransitionFail, errors.Errorf("[ERR] malicious slave ip address.")
    }
    b.slaveNode.IP4Address = meta.DiscoveryAgent.SlaveAddress
    // slave ip gateway
    if len(meta.DiscoveryAgent.SlaveGateway) == 0 {
        return MasterTransitionFail, errors.Errorf("[ERR] Inappropriate slave node gateway")
    }
    b.slaveNode.IP4Gateway = meta.DiscoveryAgent.SlaveGateway

    // save slave discovery agent
    b.slaveLocation = meta.DiscoveryAgent

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}
