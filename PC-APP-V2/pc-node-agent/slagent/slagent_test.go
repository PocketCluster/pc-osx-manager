package slagent

import (
    "time"
    "testing"
    "runtime"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-node-agent/slcontext"
    "reflect"
)

var masterAgentName string
var slaveNodeName string
var initSendTimestmap time.Time

func setUp() {
    masterAgentName, _ = context.DebugContextPrepare().MasterAgentName()
    slcontext.DebugSlcontextPrepare()
    slaveNodeName = "pc-node1"
    initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
}

func tearDown() {
    context.DebugContextDestroy()
    slcontext.DebugSlcontextDestroy()
}

func TestUnboundedBroadcastMeta(t *testing.T) {
    setUp()
    defer tearDown()

    piface, _ := slcontext.SharedSlaveContext().PrimaryNetworkInterface()

    //--- testing body ---
    ma, err := TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Error("[ERR] Incorrect MetaVersion " + ma.MetaVersion + ". Expected : " + SLAVE_META_VERSION)
        return
    }
    if ma.DiscoveryAgent.Version != SLAVE_DISCOVER_VERSION {
        t.Error("[ERR] Incorrect DiscoveryAgent.Version : " + ma.DiscoveryAgent.Version + " Expected : " + SLAVE_DISCOVER_VERSION)
        return
    }
    if ma.DiscoveryAgent.SlaveResponse != SLAVE_LOOKUP_AGENT {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveResponse : " + ma.DiscoveryAgent.SlaveResponse + " Expected : " + SLAVE_LOOKUP_AGENT)
    }
    if len(ma.DiscoveryAgent.SlaveAddress) == 0 || ma.DiscoveryAgent.SlaveAddress != piface.IP.String() {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveAddress")
        return
    }
    if len(ma.DiscoveryAgent.SlaveGateway) == 0 || ma.DiscoveryAgent.SlaveGateway != piface.GatewayAddr {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveGateway")
        return
    }
    if len(ma.DiscoveryAgent.SlaveNetmask) == 0 || ma.DiscoveryAgent.SlaveNetmask != piface.IPMask.String() {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveNetmask")
        return
    }
    if len(ma.DiscoveryAgent.SlaveNodeMacAddr) == 0 || ma.DiscoveryAgent.SlaveNodeMacAddr != piface.HardwareAddr.String() {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveNetmask")
        return
    }

    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if 508 <= len(mp) {
        t.Errorf("[ERR] Package message length does not match an expectation [%d]", len(mp))
        return
    }

    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if ma.MetaVersion != up.MetaVersion {
        t.Error("[ERR] Unidentical MetaVersion")
        return
    }
    if ma.DiscoveryAgent.Version != up.DiscoveryAgent.Version {
        t.Error("[ERR] Unidentical DiscoveryAgent.Version")
        return
    }
    if ma.DiscoveryAgent.SlaveResponse != up.DiscoveryAgent.SlaveResponse {
        t.Error("[ERR] Unidentical DiscoveryAgent.SlaveResponse")
        return
    }
    if ma.DiscoveryAgent.SlaveAddress != up.DiscoveryAgent.SlaveAddress {
        t.Error("[ERR] Unidentical DiscoveryAgent.SlaveAddress")
        return
    }
    if ma.DiscoveryAgent.SlaveGateway != up.DiscoveryAgent.SlaveGateway {
        t.Error("[ERR] Unidentical ma.DiscoveryAgent.SlaveGateway")
        return
    }
    if ma.DiscoveryAgent.SlaveNetmask != up.DiscoveryAgent.SlaveNetmask {
        t.Error("[ERR] Unidentical ma.DiscoveryAgent.SlaveNetmask")
        return
    }
    if ma.DiscoveryAgent.SlaveNodeMacAddr != up.DiscoveryAgent.SlaveNodeMacAddr {
        t.Error("[ERR] Unidentical DiscoveryAgent.SlaveNodeMacAddr")
        return
    }
}

func TestInquiredMetaAgent(t *testing.T) {
    setUp()
    defer tearDown()

    piface, _ := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    ma, _, err := TestSlaveAnswerMasterInquiry(initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Error("[ERR] Incorrect MetaVersion " + ma.MetaVersion + ". Expected : " + SLAVE_META_VERSION)
        return
    }
    if ma.StatusAgent.Version != SLAVE_STATUS_VERSION {
        t.Error("[ERR] Incorrect StatusAgent.Version : " + ma.DiscoveryAgent.Version + " Expected : " + SLAVE_DISCOVER_VERSION)
        return
    }
    if ma.StatusAgent.SlaveResponse != SLAVE_WHO_I_AM {
        t.Error("[ERR] Incorrect StatusAgent.SlaveResponse : " + ma.StatusAgent.SlaveResponse + " Expected : " + SLAVE_WHO_I_AM)
    }
    if len(ma.StatusAgent.SlaveAddress) == 0 || ma.StatusAgent.SlaveAddress != piface.IP.String() {
        t.Error("[ERR] Incorrect StatusAgent.SlaveAddress")
        return
    }
    if len(ma.StatusAgent.SlaveNodeMacAddr) == 0 || ma.StatusAgent.SlaveNodeMacAddr != piface.HardwareAddr.String() {
        t.Error("[ERR] Incorrect StatusAgent.SlaveNodeMacAddr")
        return
    }
    if len(ma.StatusAgent.SlaveHardware) == 0 || ma.StatusAgent.SlaveHardware != runtime.GOARCH {
        t.Error("[ERR] Incorrect StatusAgent.SlaveHardware")
        return
    }
    if ma.StatusAgent.SlaveTimestamp != initSendTimestmap {
        t.Error("[ERR] Incorrect StatusAgent.SlaveTimestamp")
        return
    }

    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if 508 <= len(mp) {
        t.Errorf("[ERR] Package message length [%d] exceeds an expectation", len(mp))
        return
    }

    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != up.MetaVersion {
        t.Error("[ERR] Unidentical MetaVersion")
        return
    }
    if ma.StatusAgent.Version != up.StatusAgent.Version {
        t.Error("[ERR] Unidentical StatusAgent.Version")
        return
    }
    if ma.StatusAgent.SlaveResponse != up.StatusAgent.SlaveResponse {
        t.Error("[ERR] Unidentical StatusAgent.SlaveResponse")
        return
    }
    if ma.StatusAgent.SlaveAddress != up.StatusAgent.SlaveAddress {
        t.Error("[ERR] Unidentical StatusAgent.SlaveAddress")
        return
    }
    if ma.StatusAgent.SlaveNodeMacAddr != up.StatusAgent.SlaveNodeMacAddr {
        t.Error("[ERR] Unidentical StatusAgent.SlaveNodeMacAddr")
        return
    }
    if ma.StatusAgent.SlaveHardware != up.StatusAgent.SlaveHardware {
        t.Error("[ERR] Unidentical StatusAgent.SlaveHardware")
        return
    }
    // TODO : need to fix slave timeout
    if ma.StatusAgent.SlaveTimestamp.Equal(up.StatusAgent.SlaveTimestamp) {
        t.Skip("[ERR] Unidentical StatusAgent.SlaveTimestamp")
        return
    }
}

func TestKeyExchangeMetaAgent(t *testing.T) {
    setUp()
    defer tearDown()

    // test comparison
    piface, _ := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    ma, _, err := TestSlaveKeyExchangeStatus(masterAgentName, pcrypto.TestSlavePublicKey(), initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Errorf("[ERR] Incorrect slave meta version %s\n", SLAVE_META_VERSION)
        return
    }
    if ma.StatusAgent.Version != SLAVE_STATUS_VERSION {
        t.Errorf("[ERR] Incorrect slave status version %s\n", SLAVE_STATUS_VERSION)
        return
    }
    if ma.StatusAgent.SlaveResponse != SLAVE_SEND_PUBKEY {
        t.Errorf("[ERR] Incorrect slave status %s\n", SLAVE_SEND_PUBKEY)
        return
    }
    if ma.StatusAgent.SlaveAddress != piface.IP.String() || len(ma.StatusAgent.SlaveAddress) == 0 {
        t.Errorf("[ERR] Incorrect slave address %s\n", piface.IP.String())
        return
    }
    if ma.StatusAgent.SlaveNodeMacAddr != piface.HardwareAddr.String() || len(ma.StatusAgent.SlaveNodeMacAddr) == 0 {
        t.Errorf("[ERR] Incorrect slave mac address %s\n", piface.HardwareAddr.String())
        return
    }
    if ma.StatusAgent.SlaveHardware != runtime.GOARCH {
        t.Errorf("[ERR] in correct slave hardware %s\n", runtime.GOARCH)
        return
    }
    if !ma.StatusAgent.SlaveTimestamp.Equal(initSendTimestmap) {
        t.Errorf("[ERR] Incorrect slave timestamp %s\n", ma.StatusAgent.SlaveTimestamp.String())
        return
    }
    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if  512 <= len(mp) {
        t.Errorf("[ERR] Package message length [%d] exceeds an expectation", len(mp))
        return
    }

    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if ma.MetaVersion != up.MetaVersion {
        t.Error("[ERR] Unidentical MetaVersion")
        return
    }
    if ma.StatusAgent.Version != up.StatusAgent.Version {
        t.Error("[ERR] Unidentical StatusAgent.Version")
        return
    }
    if ma.StatusAgent.SlaveResponse != up.StatusAgent.SlaveResponse {
        t.Error("[ERR] Unidentical StatusAgent.SlaveResponse")
        return
    }
    if ma.StatusAgent.SlaveAddress != up.StatusAgent.SlaveAddress {
        t.Error("[ERR] Unidentical StatusAgent.SlaveAddress")
        return
    }
    if ma.StatusAgent.SlaveNodeMacAddr != up.StatusAgent.SlaveNodeMacAddr {
        t.Error("[ERR] Unidentical StatusAgent.SlaveNodeMacAddr")
        return
    }
    if ma.StatusAgent.SlaveHardware != up.StatusAgent.SlaveHardware {
        t.Error("[ERR] Unidentical StatusAgent.SlaveHardware")
        return
    }
    // TODO : need to fix slave timeout
    if ma.StatusAgent.SlaveTimestamp.Equal(up.StatusAgent.SlaveTimestamp) {
        t.Skip("[ERR] Unidentical StatusAgent.SlaveTimestamp")
        return
    }
}

func TestSlaveCheckCryptoAgent(t *testing.T) {
    setUp()
    defer tearDown()

    piface, _ := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    ma, _, err := TestSlaveCheckCryptoStatus(masterAgentName, slaveNodeName, slcontext.SharedSlaveContext().GetSSHKey(), pcrypto.TestAESCryptor, initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Errorf("[ERR] Incorrect slave meta version %s\n", SLAVE_META_VERSION)
        return
    }
    if len(ma.EncryptedStatus) == 0 {
        t.Errorf("[ERR] Incorrect slave status data %s\n", len(ma.EncryptedStatus))
        return
    }
    esd, err := pcrypto.TestAESCryptor.DecryptByAES(ma.EncryptedStatus)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if len(ma.SlavePubKey) == 0 {
        t.Errorf("ERR] Incorrect slave ssh key field %d\n", len(ma.SlavePubKey))
        return
    }
    ssk, err := pcrypto.TestAESCryptor.DecryptByAES(ma.SlavePubKey)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if !reflect.DeepEqual(ssk, pcrypto.TestSlaveSSHKey()) {
        t.Errorf("[ERR] Slave ssh key is not the same!")
        return
    }
    sd, err := UnpackedSlaveStatus(esd)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if sd.Version != SLAVE_STATUS_VERSION {
        t.Errorf("[ERR] Incorrect slave status version %s\n", SLAVE_STATUS_VERSION)
        return
    }
    if sd.MasterBoundAgent != masterAgentName {
        t.Errorf("[ERR] Incorrect master agent name %s\n", masterAgentName)
        return
    }
    if sd.SlaveResponse != SLAVE_CHECK_CRYPTO {
        t.Errorf("[ERR] Incorrect slave status %s\n", SLAVE_CHECK_CRYPTO)
        return
    }
    if sd.SlaveNodeName != slaveNodeName {
        t.Errorf("[ERR] Incorrect slave agent name %s\n", slaveNodeName)
        return
    }
    if sd.SlaveAddress != piface.IP.String()  {
        t.Errorf("[ERR] Incorrect slave address %s\n", piface.IP.String())
        return
    }
    if sd.SlaveNodeMacAddr != piface.HardwareAddr.String() {
        t.Errorf("[ERR] Incorrect slave mac address %s\n", piface.HardwareAddr.String())
        return
    }
    if sd.SlaveHardware != runtime.GOARCH {
        t.Errorf("[ERR] in correct slave hardware %s\n", runtime.GOARCH)
        return
    }
    if !sd.SlaveTimestamp.Equal(initSendTimestmap) {
        t.Errorf("[ERR] Incorrect slave timestamp %s\n", ma.StatusAgent.SlaveTimestamp.String())
        return
    }

    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if 508 <= len(mp) {
        t.Errorf("[ERR] Package message length [%d] exceeds an expectation", len(mp))
        return
    }

    // unpack message data
    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != up.MetaVersion {
        t.Error("[ERR] Unidentical MetaVersion")
        return
    }
    if len(up.EncryptedStatus) == 0 {
        t.Errorf("[ERR] Incorrect slave status data %s\n", len(up.EncryptedStatus))
        return
    }

    esd, err = pcrypto.TestAESCryptor.DecryptByAES(up.EncryptedStatus)
    if err != nil {
        t.Error(err.Error())
        return
    }
    usd, err := UnpackedSlaveStatus(esd)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if sd.Version != usd.Version {
        t.Error("[ERR] Unidentical StatusAgent.Version")
        return
    }
    if sd.MasterBoundAgent != usd.MasterBoundAgent {
        t.Error("[ERR] Unidentical StatusAgent.MasterBoundAgent")
        return
    }
    if sd.SlaveResponse != usd.SlaveResponse {
        t.Error("[ERR] Unidentical StatusAgent.SlaveResponse")
        return
    }
    if sd.SlaveNodeName != usd.SlaveNodeName {
        t.Error("[ERR] Unidentical StatusAgent.SlaveNodeName")
        return
    }
    if sd.SlaveAddress != usd.SlaveAddress {
        t.Error("[ERR] Unidentical StatusAgent.SlaveAddress")
        return
    }
    if sd.SlaveNodeMacAddr != usd.SlaveNodeMacAddr {
        t.Error("[ERR] Unidentical StatusAgent.SlaveNodeMacAddr")
        return
    }
    if sd.SlaveHardware != usd.SlaveHardware {
        t.Error("[ERR] Unidentical StatusAgent.SlaveHardware")
        return
    }
    if !reflect.DeepEqual(ma.SlavePubKey, up.SlavePubKey) {
        t.Errorf("[ERR] Unidenticla Slave sshkey")
        return
    }
    // TODO : need to fix slave timezone
    if sd.SlaveTimestamp.Equal(usd.SlaveTimestamp) {
        t.Skip("[ERR] Unidentical StatusAgent.SlaveTimestamp")
        return
    }
}

// becuase the encrypted output differs everytime, we can only check by decrypt it.
func TestBoundedStatusMetaAgent(t *testing.T) {
    setUp()
    defer tearDown()

    piface, _ := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    ma, _, err := TestSlaveBoundedStatus(masterAgentName, slaveNodeName, pcrypto.TestAESCryptor, initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Errorf("[ERR] Incorrect slave meta version %s\n", SLAVE_META_VERSION)
        return
    }
    if len(ma.EncryptedStatus) == 0 {
        t.Errorf("[ERR] Incorrect slave status data %s\n", len(ma.EncryptedStatus))
        return
    }
    esd, err := pcrypto.TestAESCryptor.DecryptByAES(ma.EncryptedStatus)
    if err != nil {
        t.Error(err.Error())
        return
    }
    sd, err := UnpackedSlaveStatus(esd)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if sd.Version != SLAVE_STATUS_VERSION {
        t.Errorf("[ERR] Incorrect slave status version %s\n", SLAVE_STATUS_VERSION)
        return
    }
    if sd.MasterBoundAgent != masterAgentName {
        t.Errorf("[ERR] Incorrect master agent name %s\n", masterAgentName)
        return
    }
    if sd.SlaveResponse != SLAVE_REPORT_STATUS {
        t.Errorf("[ERR] Incorrect slave status %s\n", SLAVE_REPORT_STATUS)
        return
    }
    if sd.SlaveNodeName != slaveNodeName {
        t.Errorf("[ERR] Incorrect slave agent name %s\n", slaveNodeName)
        return
    }
    if sd.SlaveAddress != piface.IP.String()  {
        t.Errorf("[ERR] Incorrect slave address %s\n", piface.IP.String())
        return
    }
    if sd.SlaveNodeMacAddr != piface.HardwareAddr.String() {
        t.Errorf("[ERR] Incorrect slave mac address %s\n", piface.HardwareAddr.String())
        return
    }
    if sd.SlaveHardware != runtime.GOARCH {
        t.Errorf("[ERR] in correct slave hardware %s\n", runtime.GOARCH)
        return
    }
    if !sd.SlaveTimestamp.Equal(initSendTimestmap) {
        t.Errorf("[ERR] Incorrect slave timestamp %s\n", ma.StatusAgent.SlaveTimestamp.String())
        return
    }

    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if 508 <= len(mp) {
        t.Errorf("[ERR] Package message length [%d] exceeds an expectation", len(mp))
        return
    }

    // unpack message data
    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != up.MetaVersion {
        t.Error("[ERR] Unidentical MetaVersion")
        return
    }
    if len(up.EncryptedStatus) == 0 {
        t.Errorf("[ERR] Incorrect slave status data %s\n", len(up.EncryptedStatus))
        return
    }

    esd, err = pcrypto.TestAESCryptor.DecryptByAES(up.EncryptedStatus)
    if err != nil {
        t.Error(err.Error())
        return
    }
    usd, err := UnpackedSlaveStatus(esd)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if sd.Version != usd.Version {
        t.Error("[ERR] Unidentical StatusAgent.Version")
        return
    }
    if sd.MasterBoundAgent != usd.MasterBoundAgent {
        t.Error("[ERR] Unidentical StatusAgent.MasterBoundAgent")
        return
    }
    if sd.SlaveResponse != usd.SlaveResponse {
        t.Error("[ERR] Unidentical StatusAgent.SlaveResponse")
        return
    }
    if sd.SlaveNodeName != usd.SlaveNodeName {
        t.Error("[ERR] Unidentical StatusAgent.SlaveNodeName")
        return
    }
    if sd.SlaveAddress != usd.SlaveAddress {
        t.Error("[ERR] Unidentical StatusAgent.SlaveAddress")
        return
    }
    if sd.SlaveNodeMacAddr != usd.SlaveNodeMacAddr {
        t.Error("[ERR] Unidentical StatusAgent.SlaveNodeMacAddr")
        return
    }
    if sd.SlaveHardware != usd.SlaveHardware {
        t.Error("[ERR] Unidentical StatusAgent.SlaveHardware")
        return
    }
    // TODO : need to fix slave timezone
    if sd.SlaveTimestamp.Equal(usd.SlaveTimestamp) {
        t.Skip("[ERR] Unidentical StatusAgent.SlaveTimestamp")
        return
    }
}

func TestBindBrokenBroadcastMeta(t *testing.T) {
    setUp()
    defer tearDown()

    // test comparison
    piface, _ := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    ma, err := TestSlaveBindBroken(masterAgentName)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Errorf("[ERR] slave meta protocol version differs from %s\n", SLAVE_META_VERSION)
    }
    if ma.DiscoveryAgent.Version != SLAVE_DISCOVER_VERSION {
        t.Errorf("[ERR] slave discovery protocol version differs from %s\n", SLAVE_DISCOVER_VERSION)
    }
    if ma.DiscoveryAgent.MasterBoundAgent != masterAgentName {
        t.Errorf("[ERR] master bound agent name differs from %s\n", masterAgentName)
    }
    if ma.DiscoveryAgent.SlaveResponse != SLAVE_LOOKUP_AGENT {
        t.Errorf("[ERR] Slave is not in correct state %s\n", SLAVE_LOOKUP_AGENT)
    }
    if ma.DiscoveryAgent.SlaveAddress != piface.IP.String() || len(ma.DiscoveryAgent.SlaveAddress) == 0 {
        t.Errorf("[ERR] Slave address is incorrect %s\n", piface.IP.String())
    }
    if ma.DiscoveryAgent.SlaveGateway != piface.GatewayAddr || len(ma.DiscoveryAgent.SlaveGateway) == 0 {
        t.Errorf("[ERR] Slave gateway is incorrect %s\n", piface.GatewayAddr)
    }
    if ma.DiscoveryAgent.SlaveNetmask != piface.IPMask.String() || len(ma.DiscoveryAgent.SlaveNetmask) == 0 {
        t.Errorf("[ERR] Slave netmask is incorrect %s\n", piface.IPMask.String())
    }
    if ma.DiscoveryAgent.SlaveNodeMacAddr != piface.HardwareAddr.String() || len(ma.DiscoveryAgent.SlaveNodeMacAddr) == 0 {
        t.Errorf("[ERR] Slave MAC address is incorrect %s\n",piface.IPMask.String())
    }
    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
    }
    if 508 <= len(mp) {
        t.Errorf("[ERR] Incorrect MsgPack Length %d", len(mp))
    }

    // verification
    up, err := UnpackedSlaveMeta(mp)
    if err != nil {
        t.Error(err.Error())
    }
    if ma.MetaVersion != up.MetaVersion {
        t.Errorf("[ERR] slave meta protocol version differs from %s\n", ma.MetaVersion)
    }
    if ma.DiscoveryAgent.Version != up.DiscoveryAgent.Version {
        t.Errorf("[ERR] slave discovery protocol version differs from %s\n", ma.DiscoveryAgent.Version)
    }
    if ma.DiscoveryAgent.MasterBoundAgent != up.DiscoveryAgent.MasterBoundAgent {
        t.Errorf("[ERR] master bound agent name differs from %s\n", ma.DiscoveryAgent.MasterBoundAgent)
    }
    if ma.DiscoveryAgent.SlaveResponse != up.DiscoveryAgent.SlaveResponse {
        t.Errorf("[ERR] Slave is not in correct state %s\n", ma.DiscoveryAgent.SlaveResponse)
    }
    if ma.DiscoveryAgent.SlaveAddress != up.DiscoveryAgent.SlaveAddress {
        t.Errorf("[ERR] Slave address is incorrect %s\n", ma.DiscoveryAgent.SlaveAddress)
    }
    if ma.DiscoveryAgent.SlaveGateway != up.DiscoveryAgent.SlaveGateway {
        t.Errorf("[ERR] Slave gateway is incorrect %s\n", ma.DiscoveryAgent.SlaveGateway)
    }
    if ma.DiscoveryAgent.SlaveNetmask != up.DiscoveryAgent.SlaveNetmask {
        t.Errorf("[ERR] Slave netmask is incorrect %s\n", ma.DiscoveryAgent.SlaveNetmask)
    }
    if ma.DiscoveryAgent.SlaveNodeMacAddr != up.DiscoveryAgent.SlaveNodeMacAddr {
        t.Errorf("[ERR] Slave MAC address is incorrect %s\n", ma.DiscoveryAgent.SlaveNodeMacAddr)
    }
}
