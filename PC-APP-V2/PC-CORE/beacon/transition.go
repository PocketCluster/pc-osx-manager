package beacon

import (
    "fmt"
    "time"

    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-core/context"
)

const allowedTimesOfFailure int = 5
const timeOutWindow time.Duration = time.Second * 10

func NewBeaconForSlaveNode() MasterBeacon {
    return &masterBeacon{
        // the chances where created beacon is tried to transion in later than 10 secs are low. So we'll just initiate the state here.
        lastSuccess    : time.Now(),
        trialFailCount : 0,
        beaconState    : MasterInit,
        slaveNode      : model.NewSlaveNode(),
    }
}

type masterBeacon struct {
    // last time successfully transitioned state
    lastSuccess         time.Time
    // each time we try to make transtion and fail, count goes up.
    trialFailCount      int
    beaconState         MasterBeaconState
    slaveNode           *model.SlaveNode
    aesKey              []byte
    aesCryptor          crypt.AESCryptor
    rsaEncryptor        crypt.RsaEncryptor
}

func (mb *masterBeacon) CurrentState() MasterBeaconState {
    return mb.beaconState
}

func (mb *masterBeacon) AESKey() ([]byte, error) {
    if len(mb.aesKey) == 0 {
        return nil, fmt.Errorf("[ERR] Empty AES Key")
    }
    return mb.aesKey, nil
}
func (mb *masterBeacon) AESCryptor() (crypt.AESCryptor, error) {
    if mb.aesCryptor == nil {
        return nil, fmt.Errorf("[ERR] Null AES cryptor")
    }
    return mb.aesCryptor, nil
}

func (mb *masterBeacon) RSAEncryptor() (crypt.RsaEncryptor, error) {
    if mb.rsaEncryptor == nil {
        return nil, fmt.Errorf("[ERR] Null RSA encryptor")
    }
    return mb.rsaEncryptor, nil
}

func (mb *masterBeacon) SlaveNode() (*model.SlaveNode) {
    // TODO : copy struct that the return value is read-only
    return mb.slaveNode
}

func (mb *masterBeacon) TransitionWithTimestamp(masterTimestamp time.Time) error {
    var err error = nil
    var nextConfirmedState MasterBeaconTransition
    if mb.trialFailCount < allowedTimesOfFailure && masterTimestamp.Sub(mb.lastSuccess) < timeOutWindow {
        nextConfirmedState = MasterTransitionIdle
    } else {
        if allowedTimesOfFailure <= mb.trialFailCount {
            err = fmt.Errorf("[ERR] Transition has failed too many times already")
        } else if timeOutWindow <= masterTimestamp.Sub(mb.lastSuccess) {
            err = fmt.Errorf("[ERR] Slave did not make transition in the given time window " + timeOutWindow.String())
        }
        nextConfirmedState = MasterTransitionFail
    }

    mb.beaconState = stateTransition(mb.beaconState, nextConfirmedState)
    return err
}

func (mb *masterBeacon) TransitionWithSlaveMeta(meta *slagent.PocketSlaveAgentMeta, masterTimestamp time.Time) error {
    var transition MasterBeaconTransition
    var err error = nil

    if meta == nil || meta.MetaVersion != slagent.SLAVE_META_VERSION {
        return fmt.Errorf("[ERR] Null or incorrect version of slave meta")
    }
    if len(meta.SlaveID) == 0 {
        return fmt.Errorf("[ERR] Null or incorrect slave ID")
    }

    switch mb.beaconState {
    case MasterInit: {
        transition, err = mb.beaconInit(meta, masterTimestamp)
        break
    }
    case MasterUnbounded: {
        transition, err = mb.unbounded(meta, masterTimestamp)
        break
    }
    case MasterInquired: {
        transition, err = mb.inquired(meta, masterTimestamp)
        break
    }
    case MasterKeyExchange: {
        transition, err = mb.keyExchange(meta, masterTimestamp)
        break
    }
    case MasterCryptoCheck: {
        transition, err = mb.cryptoCheck(meta, masterTimestamp)
        break
    }
    case MasterBounded: {
        transition, err = mb.bounded(meta, masterTimestamp)
        break
    }
    case MasterBindBroken: {
        transition, err = mb.bindBroken(meta, masterTimestamp)
        break
    }

    case MasterDiscarded: {
        transition, err = MasterTransitionFail, fmt.Errorf("[INFO] discarded beacon should be collected ASAP")
        break
    }
    default:
        transition, err = MasterTransitionFail, fmt.Errorf("[ERR] managmentState should never reach default")
    }

    // this is to apply failed time count and timeout window
    finalStateCandiate := mb.translateStateWithTimeout(transition, masterTimestamp)

    // make transition regardless of the presence of error
    mb.beaconState = stateTransition(mb.beaconState, finalStateCandiate)

    return err
}

func (mb *masterBeacon) translateStateWithTimeout(nextStateCandiate MasterBeaconTransition, masterTimestamp time.Time) MasterBeaconTransition {
    var nextConfirmedState MasterBeaconTransition

    switch nextStateCandiate {
        // As MasterTransitionOk does not check timewindow, it could grant an infinite timewindow to make transition.
        // This is indeed intented as it will give us a chance to handle racing situations. Plus, TransitionWithTimestamp()
        // should have squashed suspected beacons and that's the role of TransitionWithTimestamp()
        case MasterTransitionOk: {
            mb.lastSuccess = masterTimestamp
            mb.trialFailCount = 0
            nextConfirmedState = MasterTransitionOk
            break
        }
        default:{
            if mb.trialFailCount < allowedTimesOfFailure {
                mb.trialFailCount++
            }

            if mb.trialFailCount < allowedTimesOfFailure && masterTimestamp.Sub(mb.lastSuccess) < timeOutWindow {
                nextConfirmedState = MasterTransitionIdle
            } else {
                nextConfirmedState = MasterTransitionFail
            }
        }
    }

    return nextConfirmedState
}

func (mb *masterBeacon) beaconInit(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if meta.DiscoveryAgent == nil || meta.DiscoveryAgent.Version != slagent.SLAVE_DISCOVER_VERSION {
        return MasterTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of slave discovery")
    }
    // if slave isn't looking for agent, then just return. this is not for this state.
    if meta.DiscoveryAgent.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return MasterTransitionIdle, nil
    }
    if len(meta.DiscoveryAgent.MasterBoundAgent) != 0 {
        return MasterTransitionIdle, fmt.Errorf("[ERR] Incorrect slave bind. Slave should not be bound to a master when it looks for joining")
    }
    if len(meta.DiscoveryAgent.SlaveAddress) != 0 {
        mb.slaveNode.IP4Address = meta.DiscoveryAgent.SlaveAddress
    } else {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave node address")
    }
    if len(meta.DiscoveryAgent.SlaveGateway) != 0 {
        mb.slaveNode.IP4Gateway = meta.DiscoveryAgent.SlaveGateway
    } else {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave node gateway")
    }
    if len(meta.DiscoveryAgent.SlaveNetmask) != 0 {
        mb.slaveNode.IP4Netmask = meta.DiscoveryAgent.SlaveNetmask
    } else {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave node netmask")
    }
    if meta.SlaveID != meta.DiscoveryAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if len(meta.DiscoveryAgent.SlaveNodeMacAddr) != 0 {
        mb.slaveNode.MacAddress = meta.DiscoveryAgent.SlaveNodeMacAddr
    } else {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave MAC address")
    }

    return MasterTransitionOk, nil
}

func (mb *masterBeacon) unbounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if meta.StatusAgent == nil || meta.StatusAgent.Version != slagent.SLAVE_STATUS_VERSION {
        return MasterTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of slave status")
    }
    // check if slave response is what we look for
    if meta.StatusAgent.SlaveResponse != slagent.SLAVE_WHO_I_AM {
        return MasterTransitionIdle, nil
    }
    if meta.SlaveID != meta.StatusAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if mb.slaveNode.IP4Address != meta.StatusAgent.SlaveAddress {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if mb.slaveNode.MacAddress != meta.StatusAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if len(meta.StatusAgent.SlaveHardware) != 0 {
        mb.slaveNode.Arch = meta.StatusAgent.SlaveHardware
    } else {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave architecture")
    }

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (mb *masterBeacon) inquired(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
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
    if mb.slaveNode.IP4Address != meta.StatusAgent.SlaveAddress {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if meta.SlaveID != meta.StatusAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if mb.slaveNode.MacAddress != meta.StatusAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if mb.slaveNode.Arch != meta.StatusAgent.SlaveHardware {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave architecture")
    }
    if len(meta.SlavePubKey) != 0 {
        masterPrvKey, err := context.SharedHostContext().MasterPrivateKey()
        if err != nil {
            return MasterTransitionFail, err
        }
        encryptor, err := crypt.NewEncryptorFromKeyData(meta.SlavePubKey, masterPrvKey)
        if err != nil {
            return MasterTransitionFail, err
        }
        mb.slaveNode.PublicKey = meta.SlavePubKey
        mb.rsaEncryptor = encryptor
    } else {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave public key")
    }

    aesKey := crypt.NewAESKey32Byte()
    aesCryptor, err := crypt.NewAESCrypto(aesKey)
    if err != nil {
        return MasterTransitionFail, err
    }
    mb.aesKey = aesKey
    mb.aesCryptor = aesCryptor

    nodeName, err := model.FindSlaveNameCandiate()
    if err != nil {
        return MasterTransitionFail, err
    }
    mb.slaveNode.NodeName = nodeName

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (mb *masterBeacon) keyExchange(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if len(meta.EncryptedStatus) == 0 {
        return MasterTransitionFail, fmt.Errorf("[ERR] Null encrypted slave status")
    }
    if mb.aesCryptor == nil {
        return MasterTransitionFail, fmt.Errorf("[ERR] AES Cryptor is null. This should not happen")
    }
    if mb.aesKey == nil {
        return MasterTransitionFail, fmt.Errorf("[ERR] AES Key is null. This should not happen")
    }
    plain, err := mb.aesCryptor.Decrypt(meta.EncryptedStatus)
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
    if usm.SlaveResponse != slagent.SLAVE_CHECK_CRYPTO {
        return MasterTransitionIdle, nil
    }
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return MasterTransitionFail, err
    }
    if masterAgentName != usm.MasterBoundAgent {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect master agent name from slave")
    }
    if mb.slaveNode.NodeName != usm.SlaveNodeName {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave node name beacon [%s] / slave master [%s] ", mb.slaveNode.NodeName, usm.SlaveNodeName)
    }
    if mb.slaveNode.IP4Address != usm.SlaveAddress {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if meta.SlaveID != usm.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if mb.slaveNode.MacAddress != usm.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if mb.slaveNode.Arch != usm.SlaveHardware {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave architecture")
    }

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (mb *masterBeacon) cryptoCheck(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if len(meta.EncryptedStatus) == 0 {
        return MasterTransitionFail, fmt.Errorf("[ERR] Null encrypted slave status")
    }
    if mb.aesCryptor == nil {
        return MasterTransitionFail, fmt.Errorf("[ERR] AES Cryptor is null. This should not happen")
    }
    if mb.aesKey == nil {
        return MasterTransitionFail, fmt.Errorf("[ERR] AES Key is null. This should not happen")
    }
    plain, err := mb.aesCryptor.Decrypt(meta.EncryptedStatus)
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
    if mb.slaveNode.NodeName != usm.SlaveNodeName {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave master agent")
    }
    if mb.slaveNode.IP4Address != usm.SlaveAddress {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if meta.SlaveID != usm.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if mb.slaveNode.MacAddress != usm.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if mb.slaveNode.Arch != usm.SlaveHardware {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave architecture")
    }

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (mb *masterBeacon) bounded(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if len(meta.EncryptedStatus) == 0 {
        return MasterTransitionFail, fmt.Errorf("[ERR] Null encrypted slave status")
    }
    if mb.aesCryptor == nil {
        return MasterTransitionFail, fmt.Errorf("[ERR] AES Cryptor is null. This should not happen")
    }
    if mb.aesKey == nil {
        return MasterTransitionFail, fmt.Errorf("[ERR] AES Key is null. This should not happen")
    }
    plain, err := mb.aesCryptor.Decrypt(meta.EncryptedStatus)
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
    if mb.slaveNode.NodeName != usm.SlaveNodeName {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave master agent")
    }
    if mb.slaveNode.IP4Address != usm.SlaveAddress {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if meta.SlaveID != usm.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if mb.slaveNode.MacAddress != usm.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave MAC address")
    }
    if mb.slaveNode.Arch != usm.SlaveHardware {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave architecture")
    }

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func (mb *masterBeacon) bindBroken(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
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
    if mb.slaveNode.IP4Address != meta.DiscoveryAgent.SlaveAddress {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave ip address")
    }
    if mb.slaveNode.IP4Gateway != meta.DiscoveryAgent.SlaveGateway {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave gateway address")
    }
    if mb.slaveNode.IP4Netmask != meta.DiscoveryAgent.SlaveNetmask {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave netmask address")
    }
    if meta.SlaveID != meta.DiscoveryAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if mb.slaveNode.MacAddress != meta.DiscoveryAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Incorrect slave MAC address")
    }

    // TODO : for now (v0.1.4), we'll not check slave timestamp. the validity (freshness) will be looked into.
    return MasterTransitionOk, nil
}

func stateTransition(currState MasterBeaconState, nextCondition MasterBeaconTransition) MasterBeaconState {
    var nextState MasterBeaconState
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
                fallthrough
            case MasterBounded:
                fallthrough
            case MasterBindBroken:
                nextState = MasterBounded
                break

            case MasterDiscarded:
                nextState = currState
        }
        // failed to transit
    } else if nextCondition == MasterTransitionFail {
        switch currState {

            case MasterInit:
                fallthrough
            case MasterUnbounded:
                fallthrough
            case MasterInquired:
                fallthrough
            case MasterKeyExchange:
                fallthrough
            case MasterCryptoCheck:
                nextState = MasterDiscarded
                break

            case MasterBounded:
                fallthrough
            case MasterBindBroken:
                nextState = MasterBindBroken
                break

            case MasterDiscarded:
                nextState = currState
        }
        // idle
    } else  {
        nextState = currState
    }
    return nextState
}
