package slagent

import (
    "github.com/stkim1/pc-node-agent/crypt"
    "gopkg.in/vmihailenco/msgpack.v2"
)

type PocketSlaveAgentMeta struct {
    MetaVersion         MetaProtocol                `msgpack:"pc_sl_pm"`
    StatusAgent         []byte                      `msgpack:"pc_sl_as", omitempty`
    DiscoveryAgent      *PocketSlaveDiscoveryAgent  `msgpack:"pc_sl_ad", inline, omitempty`
}

func DiscoveryMetaAgent(agent *PocketSlaveDiscoveryAgent) (*PocketSlaveAgentMeta) {
    return &PocketSlaveAgentMeta{
        MetaVersion:    SLAVE_META_VERSION,
        DiscoveryAgent: agent,
    }
}

func StatusMetaAgent(agent *PocketSlaveStatusAgent, aescrypto crypt.AESCryptor) (meta *PocketSlaveAgentMeta, err error) {
    mp, err := msgpack.Marshal(agent)
    if err != nil {
        return nil, err
    }
    crypted, err := aescrypto.Encrypt(mp)
    if err != nil {
        return nil, err
    }
    meta = &PocketSlaveAgentMeta{
        MetaVersion: SLAVE_META_VERSION,
        StatusAgent: crypted,
    }
    err = nil
    return
}