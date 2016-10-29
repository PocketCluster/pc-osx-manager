package beacon

import (
    "time"

    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/slagent"
)

func TestMasterIdentityInqueryRespond() (meta *msagent.PocketMasterAgentMeta, err error) {
    // ------------- Let's Suppose you've sent an unbounded inquery from a node over multicast net ---------------------
    ua, err := slagent.UnboundedMasterDiscovery()
    if err != nil {
        return
    }
    usm := slagent.UnboundedMasterDiscoveryMeta(ua)
    cmd, err := msagent.SlaveIdentityInqueryRespond(usm.DiscoveryAgent)
    if err != nil {
        return
    }
    meta = msagent.SlaveIdentityInquiryMeta(cmd)
    return
}

func TestMasterIdentityFixationRespond(begin time.Time) (meta *msagent.PocketMasterAgentMeta, end time.Time, err error) {
    agent, err := slagent.AnswerMasterInquiryStatus(begin)
    if err != nil {
        return
    }
    msa, err := slagent.AnswerMasterInquiryMeta(agent)
    if err != nil {
        return
    }
    // --- over master side
    end = begin.Add(time.Second)
    cmd, err := msagent.MasterDeclarationCommand(msa.StatusAgent, end)
    if err != nil {
        return
    }
    meta = msagent.MasterDeclarationMeta(cmd, crypt.TestMasterPublicKey())
    return
}

func TestMasterKeyExchangeCommand(masterBoundAgentName, slaveNodeName string, begin time.Time) (meta *msagent.PocketMasterAgentMeta, end time.Time, err error) {
    agent, err := slagent.KeyExchangeStatus(masterBoundAgentName, begin)
    if err != nil {
        return
    }
    sam, err := slagent.KeyExchangeMeta(agent, crypt.TestSlavePublicKey())
    if err != nil {
        return
    }
    // --- over master side ---
    end = begin.Add(time.Second)
    // encryptor
    rsaenc ,err := crypt.NewEncryptorFromKeyData(sam.SlavePubKey, crypt.TestMasterPrivateKey())
    if err != nil {
        return
    }
    // responding commnad
    cmd, slvstat, err := msagent.ExchangeCryptoKeyAndNameCommand(sam.StatusAgent, slaveNodeName, end)
    if err != nil {
        return
    }
    meta, err = msagent.ExchangeCryptoKeyAndNameMeta(cmd, slvstat, crypt.TestAESKey, crypt.TestAESEncryptor, rsaenc)
    return
}

func TestMasterCryptoCheckCommand(masterBoundAgentName, slaveNodeName string, begin time.Time) (meta *msagent.PocketMasterAgentMeta, end time.Time, err error) {
    agent, err := slagent.CheckSlaveCryptoStatus(masterBoundAgentName, slaveNodeName, begin)
    if err != nil {
        return
    }
    msa, err := slagent.CheckSlaveCryptoMeta(agent, crypt.TestAESEncryptor)
    if err != nil {
        return
    }
    //-------------- over master, we've received the message ----------------------
    mdsa, err := crypt.TestAESEncryptor.Decrypt(msa.EncryptedStatus)
    if err != nil {
        return
    }
    // unmarshaled, slave-status
    ussa, err := slagent.UnpackedSlaveStatus(mdsa)
    if err != nil {
        return
    }
    // master preperation
    end = begin.Add(time.Second)
    cmd, err := msagent.MasterBindReadyCommand(ussa, end)
    if err != nil {
        return
    }
    meta, err = msagent.MasterBindReadyMeta(cmd, crypt.TestAESEncryptor)
    return
}

func TestMasterBrokenBindRecoveryCommand(masterBoundAgentName string) (meta *msagent.PocketMasterAgentMeta, err error) {
    agent, err := slagent.BrokenBindDiscovery(masterBoundAgentName)
    if err != nil {
        return
    }
    sam := slagent.BrokenBindMeta(agent)
    //-------------- over master, we've received the message ----------------------
    // master preperation
    cmd, err := msagent.BrokenBindRecoverRespond(sam.DiscoveryAgent)
    if err != nil {
        return
    }
    // encryptor
    rsaenc ,err := crypt.NewEncryptorFromKeyData(crypt.TestSlavePublicKey(), crypt.TestMasterPrivateKey())
    if err != nil {
        return
    }
    meta, err = msagent.BrokenBindRecoverMeta(cmd, crypt.TestAESKey, crypt.TestAESEncryptor, rsaenc)
    return
}

