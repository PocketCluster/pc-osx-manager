package slagent

import (
    "github.com/stkim1/pc-node-agent/crypt"
    "gopkg.in/vmihailenco/msgpack.v2"
    "fmt"
)

type PocketSlaveAgentMeta struct {
    MetaVersion         MetaProtocol                `msgpack:"pc_sl_pm"`
    DiscoveryAgent      *PocketSlaveDiscovery  `msgpack:"pc_sl_ad, inline, omitempty"`
    StatusAgent         *PocketSlaveStatus     `msgpack:"pc_sl_as, inline, omitempty"`
    EncryptedStatus     []byte                      `msgpack:"pc_sl_es, omitempty"`
    SlavePubKey         []byte                      `msgpack:"pc_sl_pk, omitempty"`
}

func PackedSlaveMeta(meta *PocketSlaveAgentMeta) ([]byte, error) {
    return msgpack.Marshal(meta)
}

func UnpackedSlaveMeta(message []byte) (*PocketSlaveAgentMeta, error) {
    var meta *PocketSlaveAgentMeta
    err := msgpack.Unmarshal(message, &meta)
    if err != nil {
        return nil, err
    }
    return meta, nil
}

// --- per-state meta funcs

func UnboundedMasterSearchMeta(agent *PocketSlaveDiscovery) (*PocketSlaveAgentMeta) {
    return &PocketSlaveAgentMeta{
        MetaVersion:    SLAVE_META_VERSION,
        DiscoveryAgent: agent,
    }
}

func AnswerMasterInquiryMeta(agent *PocketSlaveStatus) (meta *PocketSlaveAgentMeta, err error) {
    meta = &PocketSlaveAgentMeta{
        MetaVersion: SLAVE_META_VERSION,
        StatusAgent: agent,
    }
    err = nil
    return
}

func KeyExchangeMeta(agent *PocketSlaveStatus, pubkey []byte) (meta *PocketSlaveAgentMeta, err error) {
    if pubkey == nil {
        err = fmt.Errorf("[ERR] You cannot pass an empty pubkey")
        return
    }
    meta = &PocketSlaveAgentMeta{
        MetaVersion: SLAVE_META_VERSION,
        StatusAgent: agent,
        SlavePubKey: pubkey,
    }
    err = nil
    return
}


func SlaveBindReadyMeta(agent *PocketSlaveStatus, aescrypto crypt.AESCryptor) (meta *PocketSlaveAgentMeta, err error) {
    mp, err := PackedSlaveStatus(agent)
    if err != nil {
        return nil, err
    }
    crypted, err := aescrypto.Encrypt(mp)
    if err != nil {
        return nil, err
    }
    meta = &PocketSlaveAgentMeta{
        MetaVersion: SLAVE_META_VERSION,
        EncryptedStatus: crypted,
    }
    err = nil
    return
}

func SlaveBoundedMeta(agent *PocketSlaveStatus, aescrypto crypt.AESCryptor) (meta *PocketSlaveAgentMeta, err error) {
    mp, err := PackedSlaveStatus(agent)
    if err != nil {
        return nil, err
    }
    crypted, err := aescrypto.Encrypt(mp)
    if err != nil {
        return nil, err
    }
    meta = &PocketSlaveAgentMeta{
        MetaVersion: SLAVE_META_VERSION,
        EncryptedStatus: crypted,
    }
    err = nil
    return
}

func BrokenBindMeta(agent *PocketSlaveDiscovery) (*PocketSlaveAgentMeta) {
    return &PocketSlaveAgentMeta{
        MetaVersion:    SLAVE_META_VERSION,
        DiscoveryAgent: agent,
    }
}
