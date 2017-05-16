package slagent

import (
    "time"

    "github.com/stkim1/pcrypto"
)

func TestSlaveUnboundedMasterSearchDiscovery() (*PocketSlaveAgentMeta, error) {
    sm, err := UnboundedMasterDiscoveryMeta()
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

func TestSlaveCheckCryptoStatus(masterAgentName, slaveAgentName, slaveUUID string, aesCryptor pcrypto.AESCryptor, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    sa, err := CheckSlaveCryptoStatus(masterAgentName, slaveAgentName, slaveUUID, begin)
    if err != nil {
        return nil, begin, err
    }
    ma, err := CheckSlaveCryptoMeta(sa, aesCryptor)
    if err != nil {
        return nil, begin, err
    }
    return ma, begin, nil
}

func TestSlaveBoundedStatus(masterAgentName, slaveNodeName, slaveUUID string, aesCryptor pcrypto.AESCryptor, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    sa, err := SlaveBoundedStatus(masterAgentName, slaveNodeName, slaveUUID, begin)
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
    bm, err := BrokenBindMeta(masterAgentName)
    if err != nil {
        return nil, err
    }
    return bm, nil
}