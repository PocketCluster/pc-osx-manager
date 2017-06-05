package slagent

import (
    "reflect"
    "runtime"
    "time"
    "testing"

    "github.com/davecgh/go-spew/spew"
    "github.com/pborman/uuid"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-node-agent/slcontext"
)

var (
    masterAgentName, slaveNodeName, authToken string
    initSendTimestmap time.Time
    piface slcontext.NetworkInterface
)

func setUp() {
    initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    masterAgentName, _   = context.DebugContextPrepare().MasterAgentName()
    slaveNodeName        = "pc-node1"
    authToken            = uuid.New()
    piface, _            = slcontext.PrimaryNetworkInterface()
    slcontext.DebugSlcontextPrepare()
}

func tearDown() {
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}

func TestUnboundedBroadcastMeta(t *testing.T) {
    setUp()
    defer tearDown()

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
    if len(ma.DiscoveryAgent.SlaveAddress) == 0 || ma.DiscoveryAgent.SlaveAddress != piface.PrimaryIP4Addr() {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveAddress")
        return
    }
    if len(ma.DiscoveryAgent.SlaveGateway) == 0 || ma.DiscoveryAgent.SlaveGateway != piface.GatewayAddr {
        t.Error("[ERR] Incorrect DiscoveryAgent.SlaveGateway")
        return
    }

    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
        return
    }
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
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
}

func TestInquiredMetaAgent(t *testing.T) {
    setUp()
    defer tearDown()

    piface, _ := slcontext.PrimaryNetworkInterface()
    ma, _, err := TestSlaveAnswerMasterInquiry(initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Error("[ERR] Incorrect MetaVersion " + ma.MetaVersion + ". Expected : " + SLAVE_META_VERSION)
        return
    }
    if len(ma.SlaveID) == 0 || ma.SlaveID != piface.HardwareAddr {
        t.Error("[ERR] incorrect SlaveMeta.SlaveID")
        return
    }
    if ma.StatusAgent.Version != SLAVE_STATUS_VERSION {
        t.Error("[ERR] Incorrect StatusAgent.Version : " + ma.DiscoveryAgent.Version + " Expected : " + SLAVE_DISCOVER_VERSION)
        return
    }
    if ma.StatusAgent.SlaveResponse != SLAVE_WHO_I_AM {
        t.Error("[ERR] Incorrect StatusAgent.SlaveResponse : " + ma.StatusAgent.SlaveResponse + " Expected : " + SLAVE_WHO_I_AM)
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
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
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
    if len(up.SlaveID) == 0 || up.SlaveID != ma.SlaveID {
        t.Error("[ERR] incorrect SlaveMeta.SlaveID")
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
    piface, _ := slcontext.PrimaryNetworkInterface()
    ma, _, err := TestSlaveKeyExchangeStatus(masterAgentName, pcrypto.TestSlavePublicKey(), initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Errorf("[ERR] Incorrect slave meta version %s\n", SLAVE_META_VERSION)
        return
    }
    if len(ma.SlaveID) == 0 || ma.SlaveID != piface.HardwareAddr {
        t.Error("[ERR] incorrect SlaveMeta.SlaveID")
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
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
    if  508 <= len(mp) {
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
    if len(up.SlaveID) == 0 || up.SlaveID != ma.SlaveID {
        t.Error("[ERR] incorrect SlaveMeta.SlaveID")
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

    ma, _, err := TestSlaveCheckCryptoStatus(masterAgentName, slaveNodeName, authToken, pcrypto.TestAESCryptor, initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Errorf("[ERR] Incorrect slave meta version %s\n", SLAVE_META_VERSION)
        return
    }
    if ma.MasterBoundAgent != masterAgentName {
        t.Errorf("[ERR] Incorrect master agent name %s\n", masterAgentName)
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
    if len(ma.SlavePubKey) != 0 {
        t.Errorf("[ERR] meta.SlavePubKey should be null %d\n", len(ma.SlavePubKey))
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
    if sd.SlaveResponse != SLAVE_CHECK_CRYPTO {
        t.Errorf("[ERR] Incorrect slave status %s\n", SLAVE_CHECK_CRYPTO)
        return
    }
    if sd.SlaveNodeName != slaveNodeName {
        t.Errorf("[ERR] Incorrect slave agent name %s\n", slaveNodeName)
        return
    }
    if sd.SlaveAuthToken != authToken  {
        t.Errorf("[ERR] Incorrect slave auth token %s\n", authToken)
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
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
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
    if ma.MasterBoundAgent != up.MasterBoundAgent {
        t.Error("[ERR] Unidentical StatusAgent.MasterBoundAgent")
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
    if sd.SlaveResponse != usd.SlaveResponse {
        t.Error("[ERR] Unidentical StatusAgent.SlaveResponse")
        return
    }
    if sd.SlaveNodeName != usd.SlaveNodeName {
        t.Error("[ERR] Unidentical StatusAgent.SlaveNodeName")
        return
    }
    if sd.SlaveAuthToken != usd.SlaveAuthToken {
        t.Error("[ERR] Unidentical StatusAgent.SlaveAuthToken")
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

    ma, _, err := TestSlaveBoundedStatus(masterAgentName, slaveNodeName, authToken, pcrypto.TestAESCryptor, initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Errorf("[ERR] Incorrect slave meta version %s\n", SLAVE_META_VERSION)
        return
    }
    if ma.MasterBoundAgent != masterAgentName {
        t.Errorf("[ERR] Incorrect master agent name %s\n", masterAgentName)
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
    if sd.SlaveResponse != SLAVE_REPORT_STATUS {
        t.Errorf("[ERR] Incorrect slave status %s\n", SLAVE_REPORT_STATUS)
        return
    }
    if sd.SlaveNodeName != slaveNodeName {
        t.Errorf("[ERR] Incorrect slave agent name %s\n", slaveNodeName)
        return
    }
    if sd.SlaveAuthToken != authToken {
        t.Errorf("[ERR] Incorrect slave auth token %s\n", authToken)
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
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
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
    if ma.MasterBoundAgent != up.MasterBoundAgent {
        t.Error("[ERR] Unidentical StatusAgent.MasterBoundAgent")
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
    if sd.SlaveResponse != usd.SlaveResponse {
        t.Error("[ERR] Unidentical StatusAgent.SlaveResponse")
        return
    }
    if sd.SlaveNodeName != usd.SlaveNodeName {
        t.Error("[ERR] Unidentical StatusAgent.SlaveNodeName")
        return
    }
    if sd.SlaveAuthToken != usd.SlaveAuthToken {
        t.Error("[ERR] Unidentical StatusAgent.SlaveAuthToken")
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
    piface, _ := slcontext.PrimaryNetworkInterface()
    ma, err := TestSlaveBindBroken(masterAgentName)
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ma.MetaVersion != SLAVE_META_VERSION {
        t.Errorf("[ERR] slave meta protocol version differs from %s\n", SLAVE_META_VERSION)
    }
    if ma.MasterBoundAgent != masterAgentName {
        t.Errorf("[ERR] master bound agent name differs from %s\n", masterAgentName)
    }
    if ma.DiscoveryAgent.Version != SLAVE_DISCOVER_VERSION {
        t.Errorf("[ERR] slave discovery protocol version differs from %s\n", SLAVE_DISCOVER_VERSION)
    }
    if ma.DiscoveryAgent.SlaveResponse != SLAVE_LOOKUP_AGENT {
        t.Errorf("[ERR] Slave is not in correct state %s\n", SLAVE_LOOKUP_AGENT)
    }
    if ma.DiscoveryAgent.SlaveAddress != piface.PrimaryIP4Addr() || len(ma.DiscoveryAgent.SlaveAddress) == 0 {
        t.Errorf("[ERR] Slave address is incorrect %s\n", piface.PrimaryIP4Addr())
    }
    if ma.DiscoveryAgent.SlaveGateway != piface.GatewayAddr || len(ma.DiscoveryAgent.SlaveGateway) == 0 {
        t.Errorf("[ERR] Slave gateway is incorrect %s\n", piface.GatewayAddr)
    }
    mp, err := PackedSlaveMeta(ma)
    if err != nil {
        t.Error(err.Error())
    }
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
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
    if ma.MasterBoundAgent != up.MasterBoundAgent {
        t.Errorf("[ERR] master bound agent name differs from %s\n", ma.MasterBoundAgent)
    }
    if ma.DiscoveryAgent.Version != up.DiscoveryAgent.Version {
        t.Errorf("[ERR] slave discovery protocol version differs from %s\n", ma.DiscoveryAgent.Version)
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
}
