package msagent

import (
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pcrypto"
)

// Let's Suppose you've received an unbounded inquery from a node over multicast net.
func TestMasterInquireSlaveRespond() (*PocketMasterAgentMeta, error) {
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

func TestMasterAgentDeclarationCommand(masterPubKey []byte, begin time.Time) (*PocketMasterAgentMeta, time.Time, error) {
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
    return MasterDeclarationMeta(cmd, masterPubKey), end, nil
}

func TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName string, slavePubKey []byte, aesKey []byte, aesCryptor pcrypto.AESCryptor, rsaEncryptor pcrypto.RsaEncryptor, begin time.Time) (*PocketMasterAgentMeta, time.Time, error) {
    msa, end, err := slagent.TestSlaveKeyExchangeStatus(masterAgentName, slavePubKey, begin)
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
    meta, err := ExchangeCryptoKeyAndNameMeta(cmd, slvstat, aesKey, aesCryptor, rsaEncryptor)
    if err != nil {
        return nil, begin, err
    }
    return meta, begin, nil
}

func TestMasterCheckCryptoCommand(masterAgentName, slaveNodeName string, slaveSSHKey []byte, aesCryptor pcrypto.AESCryptor, begin time.Time) (*PocketMasterAgentMeta, time.Time, error) {
    msa, end, err := slagent.TestSlaveCheckCryptoStatus(masterAgentName, slaveNodeName, slaveSSHKey, aesCryptor, begin)
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
    mdsa, err := aesCryptor.DecryptByAES(usm.EncryptedStatus)
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
    meta, err := MasterBindReadyMeta(cmd, aesCryptor)
    if err != nil {
        return nil, begin, err
    }
    return meta, end, nil
}

func TestMasterBoundedStatusCommand(masterAgentName, slaveNodeName string, aesCryptor pcrypto.AESCryptor, begin time.Time) (*PocketMasterAgentMeta, time.Time, error) {
    msa, end, err := slagent.TestSlaveBoundedStatus(masterAgentName, slaveNodeName, aesCryptor, begin)
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
    mdsa, err := aesCryptor.DecryptByAES(usm.EncryptedStatus)
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
    meta, err := BoundedSlaveAckMeta(cmd, aesCryptor)
    if err != nil {
        return nil, begin, err
    }
    return meta, end, nil
}

func TestMasterBrokenBindRecoveryCommand(masterBoundAgentName string, aesKey []byte, aesCryptor pcrypto.AESCryptor, rsaEncryptor pcrypto.RsaEncryptor) (meta *PocketMasterAgentMeta, err error) {
    sam, err := slagent.TestSlaveBindBroken(masterBoundAgentName)
    if err != nil {
        return
    }
    psm, err := slagent.PackedSlaveMeta(sam)
    if err != nil {
        return nil, err
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        return nil, err
    }
    // master preperation
    cmd, err := BrokenBindRecoverRespond(usm.DiscoveryAgent)
    if err != nil {
        return
    }
    // encryptor
    return BrokenBindRecoverMeta(cmd, aesKey, aesCryptor, rsaEncryptor)
}

