package msagent

import (
    "fmt"
    "time"
    "testing"
    "bytes"

    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/context"
)

var masterAgentName string
var slaveNodeName string

func setup() {
    context.DebugContextPrepare()
    sn, _ := context.SharedHostContext().MasterAgentName()
    masterAgentName = sn
    slaveNodeName = "pc-node1"
}

func destroy() {
    context.DebugContextDestroy()
}

func TestUnboundedInqueryMeta(t *testing.T) {
    setup()
    defer destroy()

    // Let's Suppose you've received an unbounded inquery from a node over multicast net.
    ua, err := slagent.UnboundedMasterDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }
    psm, err := slagent.PackedSlaveMeta(slagent.UnboundedMasterDiscoveryMeta(ua))
    if err != nil {
        t.Error(err.Error())
        return
    }
    //-------------- over master, we've received the message and need to make an inquiry "Who R U"? --------------------
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // TODO : we need ways to identify if what this package is
    cmd, err := SlaveIdentityInqueryRespond(usm.DiscoveryAgent)
    if err != nil {
        t.Error(err.Error())
        return
    }
    meta := SlaveIdentityInquiryMeta(cmd)
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if cmd.MasterCommandType != COMMAND_SLAVE_IDINQUERY {
        t.Error("[ERR] Incorrect command type. " + COMMAND_SLAVE_IDINQUERY + " is expected")
        return
    }
    if 512 <= len(mp) {
        t.Errorf("[ERR] package meta message size [%d] does not match the expected", len(mp))
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
    if meta.DiscoveryRespond.MasterBoundAgent != masterAgentName {
        t.Errorf("[ERR] Incorrect master bound name")
        return
    }
    if meta.DiscoveryRespond.MasterCommandType != COMMAND_SLAVE_IDINQUERY {
        t.Error("[ERR] Master Command is not 'COMMAND_SLAVE_IDINQUERY'")
        return
    }
    // TODO : check respond ip address
    // meta.DiscoveryRespond.MasterAddress

    // msgpack verfication
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
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
    if meta.DiscoveryRespond.MasterBoundAgent != umeta.DiscoveryRespond.MasterBoundAgent {
        t.Errorf("[ERR] Incorrectly unpacked master bound agent")
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
    setup()
    defer destroy()

    // suppose slave agent has answered question who it is
    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }
    agent, err := slagent.AnswerMasterInquiryStatus(timestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    msa, err := slagent.AnswerMasterInquiryMeta(agent)
    if err != nil {
        t.Error(err.Error())
        return
    }
    mpsm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        t.Error(err.Error())
        return
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(mpsm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    timestmap, err = time.Parse(time.RFC3339, "2012-11-01T22:08:42+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }
    cmd, err := MasterDeclarationCommand(usm.StatusAgent, timestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    meta := MasterDeclarationMeta(cmd, crypt.TestMasterPublicKey())

    if meta.MetaVersion != MASTER_META_VERSION {
        t.Error(fmt.Errorf("[ERR] wrong master meta version").Error())
        return
    }
    if meta.StatusCommand.Version != MASTER_COMMAND_VERSION {
        t.Error(fmt.Errorf("[ERR] wrong master command version").Error())
        return
    }
    if meta.StatusCommand.MasterBoundAgent != masterAgentName {
        t.Error(fmt.Errorf("[ERR] wrong master agent name").Error())
        return
    }
    if meta.StatusCommand.MasterCommandType != COMMAND_MASTER_DECLARE {
        t.Error("[ERR] Master Command is not 'COMMAND_MASTER_DECLARE'")
        return
    }

//    TODO need to check msater address, timestamp, timezone
//    if meta.StatusCommand.MasterAddress != "" {
//    }
//    if meta.StatusCommand.MasterTimestamp.String() != "" {
//    }

    mp, err := PackedMasterMeta(meta)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
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
    if meta.StatusCommand.MasterBoundAgent != umeta.StatusCommand.MasterBoundAgent {
        t.Error(fmt.Errorf("[ERR] incorrectly unpacked master bound agent").Error())
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
    setup()
    defer destroy()

    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }
    agent, err := slagent.KeyExchangeStatus(masterAgentName, timestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    msa, err := slagent.KeyExchangeMeta(agent, crypt.TestSlavePublicKey())
    if err != nil {
        t.Error(err.Error())
        return
    }
    psm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        t.Error(err.Error())
        return
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // master preperation
    timestmap, err = time.Parse(time.RFC3339, "2012-11-01T22:08:42+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }
    // encryptor
    rsaenc ,err := crypt.NewEncryptorFromKeyData(crypt.TestSlavePublicKey(), crypt.TestMasterPrivateKey())
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // responding commnad
    cmd, slvstat, err := ExchangeCryptoKeyAndNameCommand(usm.StatusAgent, slaveNodeName, timestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    meta, err := ExchangeCryptoKeyAndNameMeta(cmd, slvstat, crypt.TestAESKey, crypt.TestAESEncryptor, rsaenc)
    if err != nil {
        t.Error(err.Error())
        return
    }
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    if cmd.MasterCommandType != COMMAND_EXCHANGE_CRPTKEY {
        t.Error("[ERR] Master Command is not 'COMMAND_EXCHANGE_CRPTKEY'")
        return
    }
    if 512 <= len(mp) {
        t.Errorf("[ERR] package meta message size [%d] does not match the expected", len(mp))
    }
    if meta.MetaVersion != umeta.MetaVersion {
        t.Errorf("[ERR] package/unpacked meta version differs")
        return
    }
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
    setup()
    defer destroy()

    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }
    agent, err := slagent.CheckSlaveCryptoStatus(masterAgentName, slaveNodeName, timestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    msa, err := slagent.CheckSlaveCryptoMeta(agent, crypt.TestAESEncryptor)
    if err != nil {
        t.Error(err.Error())
        return
    }
    psm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        t.Error(err.Error())
        return
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // marshaled, descrypted, slave-status
    mdsa, err := crypt.TestAESEncryptor.Decrypt(usm.EncryptedStatus)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // unmarshaled, slave-status
    ussa, err := slagent.UnpackedSlaveStatus(mdsa)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // master preperation
    timestmap, err = time.Parse(time.RFC3339, "2012-11-01T22:08:42+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }
    // master crypto check state command
    cmd, err := MasterBindReadyCommand(ussa, timestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    meta, err := MasterBindReadyMeta(cmd, crypt.TestAESEncryptor)
    if err != nil {
        t.Error(err.Error())
        return
    }
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    if cmd.MasterCommandType != COMMAND_MASTER_BIND_READY {
        t.Error("Master Command is not 'COMMAND_MASTER_BIND_READY'")
        return
    }
    if 512 <= len(mp) {
        t.Errorf("[ERR] package meta message size [%d] does not match the expected", len(mp))
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
    setup()
    defer destroy()

    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }
    agent, err := slagent.SlaveBoundedStatus(masterAgentName, timestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    msa, err := slagent.SlaveBoundedMeta(agent, crypt.TestAESEncryptor)
    if err != nil {
        t.Error(err.Error())
        return
    }
    psm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        t.Error(err.Error())
        return
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // marshaled, descrypted, slave-status
    mdsa, err := crypt.TestAESEncryptor.Decrypt(usm.EncryptedStatus)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // unmarshaled, slave-status
    ussa, err := slagent.UnpackedSlaveStatus(mdsa)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // master preperation
    timestmap, err = time.Parse(time.RFC3339, "2012-11-01T22:08:42+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }
    // master crypto check state command
    cmd, err := BoundedSlaveAckCommand(ussa, timestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    meta, err := BoundedSlaveAckMeta(cmd, crypt.TestAESEncryptor)
    if err != nil {
        t.Error(err.Error())
        return
    }
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    if cmd.MasterCommandType != COMMAND_SLAVE_ACK {
        t.Error("[ERR] Master Command is not 'COMMAND_SLAVE_ACK'")
        return
    }
    if 512 <= len(mp) {
        t.Errorf("[ERR] package meta message size [%d] does not match the expected", len(mp))
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

func ExampleBindBrokenMeta() {
}