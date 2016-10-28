package slagent

import (
    "github.com/stkim1/pc-node-agent/crypt"
    "gopkg.in/vmihailenco/msgpack.v2"
    "fmt"
)

type PocketSlaveAgentMeta struct {
    MetaVersion         MetaProtocol                `msgpack:"pc_sl_pm"`
    DiscoveryAgent      *PocketSlaveDiscovery       `msgpack:"pc_sl_ad, inline, omitempty"`
    StatusAgent         *PocketSlaveStatus          `msgpack:"pc_sl_as, inline, omitempty"`
    EncryptedStatus     []byte                      `msgpack:"pc_sl_es, omitempty"`
    SlavePubKey         []byte                      `msgpack:"pc_sl_pk, omitempty"`
}

func PackedSlaveMeta(meta *PocketSlaveAgentMeta) ([]byte, error) {
    return msgpack.Marshal(meta)
}

func UnpackedSlaveMeta(message []byte) (meta *PocketSlaveAgentMeta, err error) {
    err = msgpack.Unmarshal(message, &meta)
    return
}

// --- per-state meta funcs

func UnboundedMasterSearchMeta(agent *PocketSlaveDiscovery) (*PocketSlaveAgentMeta) {
    return &PocketSlaveAgentMeta{
        MetaVersion:    SLAVE_META_VERSION,
        DiscoveryAgent: agent,
    }
}

func AnswerMasterInquiryMeta(agent *PocketSlaveStatus) (*PocketSlaveAgentMeta, error) {
    return &PocketSlaveAgentMeta{
        MetaVersion: SLAVE_META_VERSION,
        StatusAgent: agent,
    }, nil
}

func KeyExchangeMeta(agent *PocketSlaveStatus, pubkey []byte) (*PocketSlaveAgentMeta, error) {
    if pubkey == nil {
        return nil, fmt.Errorf("[ERR] You cannot pass an empty pubkey")
    }
    return &PocketSlaveAgentMeta{
        MetaVersion: SLAVE_META_VERSION,
        StatusAgent: agent,
        SlavePubKey: pubkey,
    }, nil
}


func SlaveBindReadyMeta(agent *PocketSlaveStatus, aescrypto crypt.AESCryptor) (*PocketSlaveAgentMeta, error) {
    mp, err := PackedSlaveStatus(agent)
    if err != nil {
        return nil, err
    }
    crypted, err := aescrypto.Encrypt(mp)
    if err != nil {
        return nil, err
    }
    return &PocketSlaveAgentMeta{
        MetaVersion: SLAVE_META_VERSION,
        EncryptedStatus: crypted,
    }, nil
}

func SlaveBoundedMeta(agent *PocketSlaveStatus, aescrypto crypt.AESCryptor) (*PocketSlaveAgentMeta, error) {
    mp, err := PackedSlaveStatus(agent)
    if err != nil {
        return nil, err
    }
    crypted, err := aescrypto.Encrypt(mp)
    if err != nil {
        return nil, err
    }
    return &PocketSlaveAgentMeta{
        MetaVersion: SLAVE_META_VERSION,
        EncryptedStatus: crypted,
    }, nil
}

func BrokenBindMeta(agent *PocketSlaveDiscovery) (*PocketSlaveAgentMeta) {
    return &PocketSlaveAgentMeta{
        MetaVersion:    SLAVE_META_VERSION,
        DiscoveryAgent: agent,
    }
}
