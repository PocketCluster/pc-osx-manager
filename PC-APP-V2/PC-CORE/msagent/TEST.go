package msagent

import (
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-node-agent/crypt"
)

// Let's Suppose you've received an unbounded inquery from a node over multicast net.
func TestSlaveIdentityInqueryRespond() (*PocketMasterAgentMeta, error) {
    msa, err := slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        return nil, err
    }
    psm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        return nil, err
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        return nil, err
    }
    // TODO : we need ways to identify if what this package is
    cmd, err := SlaveIdentityInqueryRespond(usm.DiscoveryAgent)
    if err != nil {
        return nil, err
    }
    return SlaveIdentityInquiryMeta(cmd), nil
}

func TestMasterDeclarationCommand(begin time.Time) (*PocketMasterAgentMeta, time.Time, error) {
    msa, end, err := slagent.TestSlaveAnswerMasterInquiry(begin)
    if err != nil {
        return nil, begin, err
    }
    psm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        return nil, begin, err
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        return nil, begin, err
    }
    end = end.Add(time.Second)
    cmd, err := MasterDeclarationCommand(usm.StatusAgent, end)
    if err != nil {
        return nil, begin, err
    }
    return MasterDeclarationMeta(cmd, crypt.TestMasterPublicKey()), end, nil
}

func TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName string, begin time.Time) (*PocketMasterAgentMeta, time.Time, error) {
    msa, end, err := slagent.TestSlaveKeyExchangeStatus(masterAgentName, begin)
    if err != nil {
        return nil, begin, err
    }
    psm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        return nil, begin, err
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        return nil, begin, err
    }
    // responding commnad
    masterTS := end.Add(time.Second)
    cmd, slvstat, err := ExchangeCryptoKeyAndNameCommand(usm.StatusAgent, slaveNodeName, masterTS)
    if err != nil {
        return nil, begin, err
    }
    meta, err := ExchangeCryptoKeyAndNameMeta(cmd, slvstat, crypt.TestAESKey, crypt.TestAESCryptor, crypt.TestMasterRSACryptor)
    if err != nil {
        return nil, begin, err
    }
    return meta, begin, nil
}

func TestMasterCheckCryptoCommand(masterAgentName, slaveNodeName string, begin time.Time) (*PocketMasterAgentMeta, time.Time, error) {
    msa, end, err := slagent.TestSlaveCheckCryptoStatus(masterAgentName, slaveNodeName, begin)
    if err != nil {
        return nil, begin, err
    }
    psm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        return nil, begin, err
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        return nil, begin, err
    }
    // marshaled, descrypted, slave-status
    mdsa, err := crypt.TestAESCryptor.Decrypt(usm.EncryptedStatus)
    if err != nil {
        return nil, begin, err
    }
    // unmarshaled, slave-status
    ussa, err := slagent.UnpackedSlaveStatus(mdsa)
    if err != nil {
        return nil, begin, err
    }
    // master preperation
    // master crypto check state command
    end = end.Add(time.Second)
    cmd, err := MasterBindReadyCommand(ussa, end)
    if err != nil {
        return nil, begin, err
    }
    meta, err := MasterBindReadyMeta(cmd, crypt.TestAESCryptor)
    if err != nil {
        return nil, begin, err
    }
    return meta, end, nil
}

func TestMasterBoundedStatusCommand(masterAgentName, slaveNodeName string, begin time.Time) (*PocketMasterAgentMeta, time.Time, error) {
    msa, end, err := slagent.TestSlaveBoundedStatus(slaveNodeName, begin)
    if err != nil {
        return nil, begin, err
    }
    psm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        return nil, begin, err
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        return nil, begin, err
    }
    // marshaled, descrypted, slave-status
    mdsa, err := crypt.TestAESCryptor.Decrypt(usm.EncryptedStatus)
    if err != nil {
        return nil, begin, err
    }
    // unmarshaled, slave-status
    ussa, err := slagent.UnpackedSlaveStatus(mdsa)
    if err != nil {
        return nil, begin, err
    }
    // master crypto check state command
    end = end.Add(time.Second)
    cmd, err := BoundedSlaveAckCommand(ussa, end)
    if err != nil {
        return nil, begin, err
    }
    meta, err := BoundedSlaveAckMeta(cmd, crypt.TestAESCryptor)
    if err != nil {
        return nil, begin, err
    }
    return meta, end, nil
}