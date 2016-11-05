package slagent

import (
    "fmt"

    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext"
)

type PocketSlaveAgentMeta struct {
    MetaVersion        MetaProtocol                `msgpack:"pc_sl_pm"`
    SlaveID            string                      `msgpack:"pc_sl_id"`
    DiscoveryAgent     *PocketSlaveDiscovery       `msgpack:"pc_sl_ad, inline, omitempty"`
    StatusAgent        *PocketSlaveStatus          `msgpack:"pc_sl_as, inline, omitempty"`
    EncryptedStatus    []byte                      `msgpack:"pc_sl_es, omitempty"`
    SlavePubKey        []byte                      `msgpack:"pc_sl_pk, omitempty"`
}

func PackedSlaveMeta(meta *PocketSlaveAgentMeta) ([]byte, error) {
    return msgpack.Marshal(meta)
}

func UnpackedSlaveMeta(message []byte) (meta *PocketSlaveAgentMeta, err error) {
    err = msgpack.Unmarshal(message, &meta)
    return
}

// --- per-state meta funcs

func UnboundedMasterDiscoveryMeta(agent *PocketSlaveDiscovery) (*PocketSlaveAgentMeta, error) {
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveAgentMeta{
        MetaVersion     : SLAVE_META_VERSION,
        SlaveID         : piface.HardwareAddr.String(),
        DiscoveryAgent  : agent,
    }, nil
}

func AnswerMasterInquiryMeta(agent *PocketSlaveStatus) (*PocketSlaveAgentMeta, error) {
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveAgentMeta{
        MetaVersion     : SLAVE_META_VERSION,
        SlaveID         : piface.HardwareAddr.String(),
        StatusAgent     : agent,
    }, nil
}

func KeyExchangeMeta(agent *PocketSlaveStatus, pubkey []byte) (*PocketSlaveAgentMeta, error) {
    if pubkey == nil {
        return nil, fmt.Errorf("[ERR] You cannot pass an empty pubkey")
    }
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveAgentMeta{
        MetaVersion     : SLAVE_META_VERSION,
        SlaveID         : piface.HardwareAddr.String(),
        StatusAgent     : agent,
        SlavePubKey     : pubkey,
    }, nil
}


func CheckSlaveCryptoMeta(agent *PocketSlaveStatus, aescrypto pcrypto.AESCryptor) (*PocketSlaveAgentMeta, error) {
    mp, err := PackedSlaveStatus(agent)
    if err != nil {
        return nil, err
    }
    crypted, err := aescrypto.Encrypt(mp)
    if err != nil {
        return nil, err
    }
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveAgentMeta{
        MetaVersion     : SLAVE_META_VERSION,
        SlaveID         : piface.HardwareAddr.String(),
        EncryptedStatus : crypted,
    }, nil
}

func SlaveBoundedMeta(agent *PocketSlaveStatus, aescrypto pcrypto.AESCryptor) (*PocketSlaveAgentMeta, error) {
    mp, err := PackedSlaveStatus(agent)
    if err != nil {
        return nil, err
    }
    crypted, err := aescrypto.Encrypt(mp)
    if err != nil {
        return nil, err
    }
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveAgentMeta{
        MetaVersion     : SLAVE_META_VERSION,
        SlaveID         : piface.HardwareAddr.String(),
        EncryptedStatus : crypted,
    }, nil
}

func BrokenBindMeta(agent *PocketSlaveDiscovery) (*PocketSlaveAgentMeta, error) {
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveAgentMeta{
        MetaVersion     : SLAVE_META_VERSION,
        SlaveID         : piface.HardwareAddr.String(),
        DiscoveryAgent  : agent,
    }, nil
}
