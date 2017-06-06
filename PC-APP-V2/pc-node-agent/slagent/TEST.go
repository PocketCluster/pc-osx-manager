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
    agent, err := KeyExchangeStatus(begin)
    if err != nil {
        return nil, begin, err
    }

    ma, err := KeyExchangeMeta(masterAgentName, agent, pubKey)
    if err != nil {
        return nil, begin, err
    }
    return ma, begin, err
}

func TestSlaveCheckCryptoStatus(masterAgentName, slaveAgentName, authToken string, aesCryptor pcrypto.AESCryptor, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    sa, err := CheckSlaveCryptoStatus(slaveAgentName, authToken, begin)
    if err != nil {
        return nil, begin, err
    }
    ma, err := CheckSlaveCryptoMeta(masterAgentName, sa, aesCryptor)
    if err != nil {
        return nil, begin, err
    }
    return ma, begin, nil
}

func TestSlaveBoundedStatus(masterAgentName, slaveNodeName, authToken string, aesCryptor pcrypto.AESCryptor, begin time.Time) (*PocketSlaveAgentMeta, time.Time, error) {
    sa, err := SlaveBoundedStatus(slaveNodeName, authToken, begin)
    if err != nil {
        return nil, begin, err
    }
    ma, err := SlaveBoundedMeta(masterAgentName, sa, aesCryptor)
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