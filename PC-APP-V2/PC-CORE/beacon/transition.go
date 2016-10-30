package beacon

import (
    "fmt"
    "time"

    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-core/context"
)

func stateTransition(currState MasterBeaconState, nextCondition MasterBeaconTranstion) (nextState MasterBeaconState, err error) {
    // successfully transition to the next
    if nextCondition == MasterTransitionOk {
        switch currState {
        case MasterInit:
            nextState = MasterUnbounded
        case MasterUnbounded:
            nextState = MasterInquired
        case MasterInquired:
            nextState = MasterKeyExchange
        case MasterKeyExchange:
            nextState = MasterCryptoCheck
        case MasterCryptoCheck:
            nextState = MasterBounded
        case MasterBounded:
            nextState = MasterBounded
        case MasterBindBroken:
            nextState = MasterBounded
        default:
            err = fmt.Errorf("[ERR] 'nextCondition is true and hit default' cannot happen")
        }
        // failed to transit
    } else if nextCondition == MasterTransitionFail {
        switch currState {
        case MasterInit:
            nextState = MasterDiscarded
        case MasterUnbounded:
            nextState = MasterDiscarded
        case MasterInquired:
            nextState = MasterDiscarded
        case MasterKeyExchange:
            nextState = MasterDiscarded
        case MasterCryptoCheck:
            nextState = MasterDiscarded
        case MasterBounded:
            nextState = MasterBindBroken
        case MasterBindBroken:
            nextState = MasterBindBroken
        default:
            err = fmt.Errorf("[ERR] 'nextCondition is true and hit default' cannot happen")
        }
        // idle
    } else  {
        nextState = currState
    }
    return
}

func NewBeaconForSlaveNode() MasterBeacon {
    return &masterBeacon{
        beaconState    : MasterUnbounded,
        slaveNode      : &model.SlaveNode{},
    }
}

type masterBeacon struct {
    lastSuccess         time.Time
    beaconState         MasterBeaconState
    slaveNode           *model.SlaveNode
    aesKey              []byte
    aesCryptor          crypt.AESCryptor
    rsaEncryptor        crypt.RsaEncryptor
    rsaDecryptor        crypt.RsaDecryptor
}

func (mb *masterBeacon) CurrentState() MasterBeaconState {
    return mb.beaconState
}

func (mb *masterBeacon) TranstionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta == nil || meta.MetaVersion != slagent.SLAVE_META_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave meta")
    }
    switch mb.beaconState {
    case MasterInit:
        return mb.beaconInit(meta, timestamp)

    case MasterUnbounded:
        return mb.unbounded(meta, timestamp)

    case MasterInquired:
        return mb.inquired(meta, timestamp)

    case MasterKeyExchange:
        return mb.keyExchange(meta, timestamp)

    case MasterCryptoCheck:
        return mb.cryptoCheck(meta, timestamp)

    case MasterBounded:
        return mb.bounded(meta, timestamp)

    case MasterBindBroken:
        return mb.bindBroken(meta, timestamp)

    default:
        return fmt.Errorf("[ERR] managmentState cannot reach default")
    }
}

func (mb *masterBeacon) beaconInit(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta.DiscoveryAgent == nil || meta.DiscoveryAgent.Version != slagent.SLAVE_DISCOVER_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave discovery")
    }
    // if slave isn't looking for agent, then just return. this is not for this state.
    if meta.DiscoveryAgent.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return nil
    }
    if len(meta.DiscoveryAgent.SlaveAddress) != 0 {
        mb.slaveNode.IP4Address = meta.DiscoveryAgent.SlaveAddress
    } else {
        return fmt.Errorf("[ERR] Inappropriate slave node address")
    }
    if len(meta.DiscoveryAgent.SlaveGateway) != 0 {
        mb.slaveNode.IP4Gateway = meta.DiscoveryAgent.SlaveGateway
    } else {
        return fmt.Errorf("[ERR] Inappropriate slave node gateway")
    }
    if len(meta.DiscoveryAgent.SlaveNetmask) != 0 {
        mb.slaveNode.IP4Netmask = meta.DiscoveryAgent.SlaveNetmask
    } else {
        return fmt.Errorf("[ERR] Inappropriate slave node netmask")
    }
    if len(meta.DiscoveryAgent.SlaveNodeMacAddr) != 0 {
        mb.slaveNode.MacAddress = meta.DiscoveryAgent.SlaveNodeMacAddr
    } else {
        return fmt.Errorf("[ERR] Inappropriate slave MAC address")
    }

    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) unbounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta.StatusAgent == nil || meta.StatusAgent.Version != slagent.SLAVE_STATUS_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if meta.StatusAgent.SlaveResponse != slagent.SLAVE_WHO_I_AM {
        return nil
    }
    if mb.slaveNode.IP4Address != meta.StatusAgent.SlaveAddress {
        return fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if mb.slaveNode.MacAddress != meta.StatusAgent.SlaveNodeMacAddr {
        return fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if len(meta.StatusAgent.SlaveHardware) != 0 {
        mb.slaveNode.Arch = meta.StatusAgent.SlaveHardware
    } else {
        return fmt.Errorf("[ERR] Inappropriate slave architecture")
    }

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) inquired(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta.StatusAgent == nil || meta.StatusAgent.Version != slagent.SLAVE_STATUS_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if meta.StatusAgent.SlaveResponse != slagent.SLAVE_SEND_PUBKEY {
        return nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return err
    }
    if masterAgentName != meta.StatusAgent.MasterBoundAgent {
        return fmt.Errorf("[ERR] Slave reports to incorrect master agent")
    }
    if mb.slaveNode.IP4Address != meta.StatusAgent.SlaveAddress {
        return fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if mb.slaveNode.MacAddress != meta.StatusAgent.SlaveNodeMacAddr {
        return fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if mb.slaveNode.Arch != meta.StatusAgent.SlaveHardware {
        return fmt.Errorf("[ERR] Incorrect slave architecture")
    }
    if len(meta.SlavePubKey) != 0 {
        // TODO fix PublicKey to []byte
        // mb.slaveNode.PublicKey = meta.SlavePubKey
        
        //TODO : build RSA enc/decryptor
        //mb.rsaDecryptor = crypt.RsaEncryptor()
        //mb.rsaDecryptor = crypt.RsaDecryptor()
    } else {
        return fmt.Errorf("[ERR] Inappropriate slave public key")
    }
    aesKey := crypt.NewAESKey32Byte()
    aesCryptor, err := crypt.NewAESCrypto(aesKey)
    if err != nil {
        return err
    }
    mb.aesKey = aesKey
    mb.aesCryptor = aesCryptor

    // TODO : generate node name
    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) keyExchange(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if len(meta.EncryptedStatus) == 0 {
        return fmt.Errorf("[ERR] Null encrypted slave status")
    }
    if mb.aesCryptor == nil {
        return fmt.Errorf("[ERR] AES Cryptor is null. This should not happen")
    }
    if mb.aesKey == nil {
        return fmt.Errorf("[ERR] AES Key is null. This should not happen")
    }
    plain, err := mb.aesCryptor.Decrypt(meta.EncryptedStatus)
    if err != nil {
        return err
    }
    usm, err := slagent.UnpackedSlaveStatus(plain)
    if err != nil {
        return err
    }
    if usm == nil || usm.Version != slagent.SLAVE_STATUS_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if usm.SlaveResponse != slagent.SLAVE_CHECK_CRYPTO {
        return nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return err
    }
    if masterAgentName != usm.MasterBoundAgent {
        return fmt.Errorf("[ERR] Incorrect master agent name from slave")
    }
    if mb.slaveNode.NodeName != usm.SlaveNodeName {
        return fmt.Errorf("[ERR] Incorrect slave master agent")
    }
    if mb.slaveNode.IP4Address != usm.SlaveAddress {
        return fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if mb.slaveNode.MacAddress != usm.SlaveNodeMacAddr {
        return fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if mb.slaveNode.Arch != usm.SlaveHardware {
        return fmt.Errorf("[ERR] Incorrect slave architecture")
    }

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) cryptoCheck(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if len(meta.EncryptedStatus) == 0 {
        return fmt.Errorf("[ERR] Null encrypted slave status")
    }
    if mb.aesCryptor == nil {
        return fmt.Errorf("[ERR] AES Cryptor is null. This should not happen")
    }
    if mb.aesKey == nil {
        return fmt.Errorf("[ERR] AES Key is null. This should not happen")
    }
    plain, err := mb.aesCryptor.Decrypt(meta.EncryptedStatus)
    if err != nil {
        return err
    }
    usm, err := slagent.UnpackedSlaveStatus(plain)
    if err != nil {
        return err
    }
    if usm == nil || usm.Version != slagent.SLAVE_STATUS_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if usm.SlaveResponse != slagent.SLAVE_REPORT_STATUS {
        return nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return err
    }
    if masterAgentName != usm.MasterBoundAgent {
        return fmt.Errorf("[ERR] Incorrect master agent name from slave")
    }
    if mb.slaveNode.NodeName != usm.SlaveNodeName {
        return fmt.Errorf("[ERR] Incorrect slave master agent")
    }
    if mb.slaveNode.IP4Address != usm.SlaveAddress {
        return fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if mb.slaveNode.MacAddress != usm.SlaveNodeMacAddr {
        return fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if mb.slaveNode.Arch != usm.SlaveHardware {
        return fmt.Errorf("[ERR] Incorrect slave architecture")
    }

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) bounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if len(meta.EncryptedStatus) == 0 {
        return fmt.Errorf("[ERR] Null encrypted slave status")
    }
    if mb.aesCryptor == nil {
        return fmt.Errorf("[ERR] AES Cryptor is null. This should not happen")
    }
    if mb.aesKey == nil {
        return fmt.Errorf("[ERR] AES Key is null. This should not happen")
    }
    plain, err := mb.aesCryptor.Decrypt(meta.EncryptedStatus)
    if err != nil {
        return err
    }
    usm, err := slagent.UnpackedSlaveStatus(plain)
    if err != nil {
        return err
    }
    if usm == nil || usm.Version != slagent.SLAVE_STATUS_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if usm.SlaveResponse != slagent.SLAVE_REPORT_STATUS {
        return nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return err
    }
    if masterAgentName != usm.MasterBoundAgent {
        return fmt.Errorf("[ERR] Incorrect master agent name from slave")
    }
    if mb.slaveNode.NodeName != usm.SlaveNodeName {
        return fmt.Errorf("[ERR] Incorrect slave master agent")
    }
    if mb.slaveNode.IP4Address != usm.SlaveAddress {
        return fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if mb.slaveNode.MacAddress != usm.SlaveNodeMacAddr {
        return fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if mb.slaveNode.Arch != usm.SlaveHardware {
        return fmt.Errorf("[ERR] Incorrect slave architecture")
    }

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}

func (mb *masterBeacon) bindBroken(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) error {
    if meta.DiscoveryAgent == nil || meta.DiscoveryAgent.Version != slagent.SLAVE_DISCOVER_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // if slave isn't looking for agent, then just return. this is not for this state.
    if meta.DiscoveryAgent.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return err
    }
    // since this node isn't looking for us, sliently ignore this request
    if masterAgentName != meta.DiscoveryAgent.MasterBoundAgent {
        return nil
    }
    if mb.slaveNode.IP4Address != meta.DiscoveryAgent.SlaveAddress {
        return fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if mb.slaveNode.IP4Gateway != meta.DiscoveryAgent.SlaveGateway {
        return fmt.Errorf("[ERR] Incorrect slave gateway address")
    }
    if mb.slaveNode.IP4Netmask != meta.DiscoveryAgent.SlaveNetmask {
        return fmt.Errorf("[ERR] Incorrect slave netmask address")
    }
    if mb.slaveNode.MacAddress != meta.DiscoveryAgent.SlaveNodeMacAddr {
        return fmt.Errorf("[ERR] Incorrect slave MAC address")
    }

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    state, err := stateTransition(mb.beaconState, MasterTransitionOk)
    if err != nil {
        return err
    }
    mb.beaconState = state
    return nil
}
