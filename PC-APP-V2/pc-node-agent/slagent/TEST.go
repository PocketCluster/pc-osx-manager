package slagent

import (
    "time"

    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-core/context"
)

func TestSlaveUnboundedMasterSearchDiscovery() (*PocketSlaveAgentMeta, error) {
    ua, err := UnboundedMasterDiscovery()
    if err != nil {
        return nil, err
    }
    return UnboundedMasterDiscoveryMeta(ua), nil
}

func TestSlaveAnswerMasterInquiry(begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    agent, err := AnswerMasterInquiryStatus(begin)
    if err != nil {
        return nil, begin, err
    }
    ma, err := AnswerMasterInquiryMeta(agent)
    if err != nil {
        return nil, begin, err
    }
    return ma, begin, nil
}

func TestSlaveKeyExchangeStatus(masterAgentName string, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    agent, err := KeyExchangeStatus(masterAgentName, begin)
    if err != nil {
        return nil, begin, err
    }

    ma, err := KeyExchangeMeta(agent, crypt.TestSlavePublicKey())
    if err != nil {
        return nil, begin, err
    }
    return ma, begin, err
}

func TestCheckSlaveCryptoStatus(masterAgentName, slaveAgentName string, aesEnc crypt.AESCryptor, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    sa, err := CheckSlaveCryptoStatus(masterAgentName, slaveAgentName, begin)
    if err != nil {
        return nil, begin, err
    }
    ma, err := CheckSlaveCryptoMeta(sa, aesEnc)
    if err != nil {
        return nil, begin, err
    }
    return ma, begin, nil
}

func TestSlaveBoundedStatus(slaveNodeName string, aesEnc crypt.AESCryptor, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return nil, begin, err
    }
    sa, err := SlaveBoundedStatus(masterAgentName, slaveNodeName, begin)
    if err != nil {
        return nil, begin, err
    }
    ma, err := SlaveBoundedMeta(sa, aesEnc)
    if err != nil {
        return nil, begin, err
    }
    return ma, begin, nil
}

func TestSlaveBindBroken(masterAgentName string) (*PocketSlaveAgentMeta, error) {
    ba, err := BrokenBindDiscovery(masterAgentName)
    if err != nil {
        return nil, err
    }
    return BrokenBindMeta(ba), nil
}