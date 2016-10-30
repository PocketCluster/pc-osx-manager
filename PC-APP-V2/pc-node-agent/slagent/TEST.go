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

func TestSlaveKeyExchangeStatus(masterAgentName string, pubKey []byte, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    agent, err := KeyExchangeStatus(masterAgentName, begin)
    if err != nil {
        return nil, begin, err
    }

    ma, err := KeyExchangeMeta(agent, pubKey)
    if err != nil {
        return nil, begin, err
    }
    return ma, begin, err
}

func TestSlaveCheckCryptoStatus(masterAgentName, slaveAgentName string, aesCryptor crypt.AESCryptor, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    sa, err := CheckSlaveCryptoStatus(masterAgentName, slaveAgentName, begin)
    if err != nil {
        return nil, begin, err
    }
    ma, err := CheckSlaveCryptoMeta(sa, aesCryptor)
    if err != nil {
        return nil, begin, err
    }
    return ma, begin, nil
}

func TestSlaveBoundedStatus(slaveNodeName string, aesCryptor crypt.AESCryptor, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    masterAgentName, err := context.SharedHostContext().MasterAgentName()
    if err != nil {
        return nil, begin, err
    }
    sa, err := SlaveBoundedStatus(masterAgentName, slaveNodeName, begin)
    if err != nil {
        return nil, begin, err
    }
    ma, err := SlaveBoundedMeta(sa, aesCryptor)
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