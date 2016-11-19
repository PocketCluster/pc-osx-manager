package slagent

import (
    "time"

    "github.com/stkim1/pcrypto"
)

func TestSlaveUnboundedMasterSearchDiscovery() (*PocketSlaveAgentMeta, error) {
    ua, err := UnboundedMasterDiscovery()
    if err != nil {
        return nil, err
    }
    sm, err := UnboundedMasterDiscoveryMeta(ua)
    if err != nil {
        return nil, err
    }
    return sm, nil
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

func TestSlaveCheckCryptoStatus(masterAgentName, slaveAgentName string, sshkey []byte, aesCryptor pcrypto.AESCryptor, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    sa, err := CheckSlaveCryptoStatus(masterAgentName, slaveAgentName, begin)
    if err != nil {
        return nil, begin, err
    }
    ma, err := CheckSlaveCryptoMeta(sa, sshkey, aesCryptor)
    if err != nil {
        return nil, begin, err
    }
    return ma, begin, nil
}

func TestSlaveBoundedStatus(masterAgentName, slaveNodeName string, aesCryptor pcrypto.AESCryptor, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
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
    bm, err := BrokenBindMeta(ba)
    if err != nil {
        return nil, err
    }
    return bm, nil
}