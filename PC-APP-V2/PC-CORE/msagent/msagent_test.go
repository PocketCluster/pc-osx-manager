package msagent

import (
    "bytes"
    "fmt"
    "time"
    "testing"

    "github.com/davecgh/go-spew/spew"
    "github.com/pborman/uuid"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pcrypto"
)

var (
    masterAgentName, slaveNodeName, slaveUUID string
    initTime time.Time
)

func setUp() {
    context.DebugContextPrepare()
    slcontext.DebugSlcontextPrepare()

    masterAgentName, _ = context.SharedHostContext().MasterAgentName()
    slaveNodeName = "pc-node1"
    slaveUUID = uuid.New()
    initTime, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
}

func tearDown() {
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}

func TestUnboundedInqueryMeta(t *testing.T) {
    setUp()
    defer tearDown()

    paddr, err := context.SharedHostContext().HostPrimaryAddress()
    if err != nil {
        t.Error(err.Error())
        return
    }
    meta, err := TestMasterInquireSlaveRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if meta.MetaVersion != MASTER_META_VERSION {
        t.Errorf("[ERR] Incorrect master meta version")
        return
    }
    if meta.DiscoveryRespond.Version != MASTER_RESPOND_VERSION {
        t.Errorf("[ERR] Incorrect master respond version")
        return
    }
    if meta.DiscoveryRespond.MasterCommandType != COMMAND_SLAVE_IDINQUERY {
        t.Error("[ERR] Incorrect command type. " + COMMAND_SLAVE_IDINQUERY + " is expected")
        return
    }
    if meta.MasterBoundAgent != masterAgentName {
        t.Errorf("[ERR] Incorrect master bound name")
        return
    }
    if meta.DiscoveryRespond.MasterCommandType != COMMAND_SLAVE_IDINQUERY {
        t.Error("[ERR] Master Command is not 'COMMAND_SLAVE_IDINQUERY'")
        return
    }
    // TODO : check respond ip address
    if meta.DiscoveryRespond.MasterAddress != paddr {
        t.Error("[ERR] Master Address is not " + paddr)
        return
    }

    mp, err := PackedMasterMeta(meta)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // http://stackoverflow.com/questions/1098897/what-is-the-largest-safe-udp-packet-size-on-the-internet
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
    if 508 <= len(mp) {
        t.Errorf("[ERR] package meta message size [%d] does not match the expected", len(mp))
        return
    }
    // msgpack verfication
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    if meta.MasterBoundAgent != umeta.MasterBoundAgent {
        t.Errorf("[ERR] Incorrectly unpacked master bound agent")
        return
    }
    if meta.MetaVersion != umeta.MetaVersion {
        t.Errorf("[ERR] Incorrectly unpacked meta version")
        return
    }
    if meta.DiscoveryRespond.Version != umeta.DiscoveryRespond.Version {
        t.Errorf("[ERR] Incorrectly unpacked respond version")
        return
    }
    if meta.DiscoveryRespond.MasterCommandType != umeta.DiscoveryRespond.MasterCommandType {
        t.Errorf("[ERR] Incorrectly unpacked master command")
        return
    }
    if meta.DiscoveryRespond.MasterAddress != umeta.DiscoveryRespond.MasterAddress {
        t.Errorf("[ERR] Incorrectly unpacked master address")
        return
    }
}

func TestMasterDeclarationMeta(t *testing.T) {
    setUp()
    defer tearDown()

    paddr, err := context.SharedHostContext().HostPrimaryAddress()
    if err != nil {
        t.Error(err.Error())
        return
    }
    meta, end, err := TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), initTime)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if meta.MetaVersion != MASTER_META_VERSION {
        t.Error(fmt.Errorf("[ERR] wrong master meta version").Error())
        return
    }
    if meta.StatusCommand.Version != MASTER_COMMAND_VERSION {
        t.Error(fmt.Errorf("[ERR] wrong master command version").Error())
        return
    }
    if meta.MasterBoundAgent != masterAgentName {
        t.Error(fmt.Errorf("[ERR] wrong master agent name").Error())
        return
    }
    if meta.StatusCommand.MasterCommandType != COMMAND_MASTER_DECLARE {
        t.Error("[ERR] Master Command is not 'COMMAND_MASTER_DECLARE'")
        return
    }
    if meta.StatusCommand.MasterAddress != paddr {
        t.Error("[ERR] Incorrect Master ip address")
        return
    }
    // TODO need to check msater timezone
    if !meta.StatusCommand.MasterTimestamp.Equal(end) {
        t.Error("[ERR] Incorrect Master TimeStamp")
        return
    }

    mp, err := PackedMasterMeta(meta)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
    if 508 <= len(mp) {
        t.Errorf("[ERR] Package message length does not match an expectation [%d]", len(mp))
        return
    }

    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    if meta.MasterBoundAgent != umeta.MasterBoundAgent {
        t.Error(fmt.Errorf("[ERR] incorrectly unpacked master bound agent").Error())
        return
    }
    if meta.MetaVersion != umeta.MetaVersion {
        t.Error(fmt.Errorf("[ERR] incorrectly unpacked master meta version").Error())
        return
    }
    if meta.StatusCommand.Version != umeta.StatusCommand.Version {
        t.Error(fmt.Errorf("[ERR] incorrectly unpacked master command version").Error())
        return
    }
    if meta.StatusCommand.MasterCommandType != umeta.StatusCommand.MasterCommandType {
        t.Error(fmt.Errorf("[ERR] incorrectly unpacked master command type").Error())
        return
    }
    if meta.StatusCommand.MasterAddress != umeta.StatusCommand.MasterAddress {
        t.Error(fmt.Errorf("[ERR] incorrectly unpacked master address").Error())
        return
    }
    if !meta.StatusCommand.MasterTimestamp.Equal(umeta.StatusCommand.MasterTimestamp) {
        t.Error(fmt.Errorf("[ERR] incorrectly unpacked master timestamp").Error())
        return
    }
}


func TestExecKeyExchangeMeta(t *testing.T) {
    setUp()
    defer tearDown()

    meta, _, err := TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName, slaveUUID, pcrypto.TestSlavePublicKey(), pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor, initTime)
    if err != nil {
        t.Error(err.Error())
        return
    }
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        t.Error(err.Error())
        return
    }
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
    if 508 <= len(mp) {
        t.Errorf("[ERR] package meta message size [%d] does not match the expected", len(mp))
    }
    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    if meta.MetaVersion != umeta.MetaVersion {
        t.Errorf("[ERR] package/unpacked meta version differs")
        return
    }
/*
    if meta.StatusCommand.MasterCommandType != COMMAND_EXCHANGE_CRPTKEY {
        t.Error("[ERR] Master Command is not 'COMMAND_EXCHANGE_CRPTKEY'")
        return
    }
*/
    if bytes.Compare(meta.EncryptedMasterCommand, umeta.EncryptedMasterCommand) != 0{
        t.Errorf("[ERR] package/unpacked encrypted command differs")
        return
    }
    if bytes.Compare(meta.EncryptedSlaveStatus, umeta.EncryptedSlaveStatus) != 0{
        t.Errorf("[ERR] package/unpacked encrypted slave response differs")
        return
    }
    if bytes.Compare(meta.EncryptedAESKey, umeta.EncryptedAESKey) != 0{
        t.Errorf("[ERR] package/unpacked encrypted aes key differs")
        return
    }
    if bytes.Compare(meta.RsaCryptoSignature, umeta.RsaCryptoSignature) != 0{
        t.Errorf("[ERR] package/unpacked encryption signature differs")
        return
    }
}

func TestSendCryptoCheckMeta(t *testing.T) {
    setUp()
    defer tearDown()

    meta, _, err := TestMasterCheckCryptoCommand(masterAgentName, slaveNodeName, slaveUUID, pcrypto.TestAESCryptor, initTime)
    if err != nil {
        t.Error(err.Error())
        return
    }
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        t.Error(err.Error())
        return
    }
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
    if 508 <= len(mp) {
        t.Errorf("[ERR] package meta message size [%d] does not match the expected", len(mp))
    }
    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    if meta.MetaVersion != umeta.MetaVersion {
        t.Errorf("[ERR] package/unpacked meta version differs")
        return
    }
    if bytes.Compare(meta.EncryptedMasterCommand, umeta.EncryptedMasterCommand) != 0{
        t.Errorf("[ERR] package/unpacked encrypted command differs")
        return
    }
}

func TestBoundedStatusMeta(t *testing.T) {
    setUp()
    defer tearDown()

    meta, _, err := TestMasterBoundedStatusCommand(masterAgentName, slaveNodeName, slaveUUID, pcrypto.TestAESCryptor, initTime)
    if err != nil {
        t.Error(err.Error())
        return
    }
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        t.Error(err.Error())
        return
    }
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
    if 508 <= len(mp) {
        t.Errorf("[ERR] package meta message size [%d] does not match the expected", len(mp))
        return
    }
    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    if meta.MetaVersion != umeta.MetaVersion {
        t.Errorf("[ERR] package/unpacked meta version differs")
        return
    }
    if bytes.Compare(meta.EncryptedMasterCommand, umeta.EncryptedMasterCommand) != 0{
        t.Errorf("[ERR] package/unpacked encrypted command differs")
        return
    }
}

func TestBindBrokenMeta(t *testing.T) {
    setUp()
    defer tearDown()

    meta, err := TestMasterBrokenBindRecoveryCommand(masterAgentName, pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        t.Error(err.Error())
        return
    }
    t.Logf("[INFO] size of packed meta %d\n package content %s", len(mp), spew.Sdump(mp))
    if 508 <= len(mp) {
        t.Errorf("[ERR] package meta message size [%d] does not match the expected", len(mp))
        return
    }
    // verification step
    _, err = UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

}