package slagent

import (
    "fmt"

    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext"
)

type PocketSlaveAgentMeta struct {
    MetaVersion        MetaProtocol                `msgpack:"s_pm"`
    SlaveID            string                      `msgpack:"s_id"`
    DiscoveryAgent     *PocketSlaveDiscovery       `msgpack:"s_ad, inline, omitempty"`
    StatusAgent        *PocketSlaveStatus          `msgpack:"s_as, inline, omitempty"`
    EncryptedStatus    []byte                      `msgpack:"s_es, omitempty"`
    SlavePubKey        []byte                      `msgpack:"s_pk, omitempty"`
    //EncryptedSlaveSSHKey []byte                    `msgpack:"pc_sl_sh, omitempty"`
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


func CheckSlaveCryptoMeta(agent *PocketSlaveStatus, sshKey []byte, aescrypto pcrypto.AESCryptor) (*PocketSlaveAgentMeta, error) {
    mp, err := PackedSlaveStatus(agent)
    if err != nil {
        return nil, err
    }
    crypted, err := aescrypto.EncryptByAES(mp)
    if err != nil {
        return nil, err
    }
    if len(sshKey) == 0 {
        return nil, fmt.Errorf("[ERR] Cannot send empty sshkey")
    }
    cryptedSshKey, err := aescrypto.EncryptByAES(sshKey)
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
        // FIXME : this is wrong at many levels but if we're to extend one more field, the packet size will exceed 508 @ keyexchange state.
        SlavePubKey     : cryptedSshKey,
        //EncryptedSlaveSSHKey : cryptedSshKey,
    }, nil
}

func SlaveBoundedMeta(agent *PocketSlaveStatus, aescrypto pcrypto.AESCryptor) (*PocketSlaveAgentMeta, error) {
    mp, err := PackedSlaveStatus(agent)
    if err != nil {
        return nil, err
    }
    crypted, err := aescrypto.EncryptByAES(mp)
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
