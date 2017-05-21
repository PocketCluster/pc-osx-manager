package slagent

import (
    "github.com/pkg/errors"
    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext"
)

type PocketSlaveAgentMeta struct {
    MetaVersion        MetaProtocol                `msgpack:"s_pm"`
    MasterBoundAgent   string                      `msgpack:"m_ba, omitempty"`
    SlaveID            string                      `msgpack:"s_id"`
    DiscoveryAgent     *PocketSlaveDiscovery       `msgpack:"s_ad, inline, omitempty"`
    StatusAgent        *PocketSlaveStatus          `msgpack:"s_as, inline, omitempty"`
    EncryptedStatus    []byte                      `msgpack:"s_es, omitempty"`
    SlavePubKey        []byte                      `msgpack:"s_pk, omitempty"`
}

func PackedSlaveMeta(meta *PocketSlaveAgentMeta) ([]byte, error) {
    pm, err := msgpack.Marshal(meta)
    return pm, errors.WithStack(err)
}

func UnpackedSlaveMeta(message []byte) (meta *PocketSlaveAgentMeta, err error) {
    err = errors.WithStack(msgpack.Unmarshal(message, &meta))
    return
}

// --- per-state meta funcs

func UnboundedMasterDiscoveryMeta() (*PocketSlaveAgentMeta, error) {
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketSlaveAgentMeta{
        MetaVersion:       SLAVE_META_VERSION,
        SlaveID:           piface.HardwareAddr,
        DiscoveryAgent:    &PocketSlaveDiscovery {
            Version:             SLAVE_DISCOVER_VERSION,
            SlaveResponse:       SLAVE_LOOKUP_AGENT,
            SlaveAddress:        piface.PrimaryIP4Addr(),
            SlaveGateway:        piface.GatewayAddr,
        },
    }, nil
}

func AnswerMasterInquiryMeta(agent *PocketSlaveStatus) (*PocketSlaveAgentMeta, error) {
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketSlaveAgentMeta{
        MetaVersion:         SLAVE_META_VERSION,
        SlaveID:             piface.HardwareAddr,
        StatusAgent:         agent,
    }, nil
}

func KeyExchangeMeta(master string, agent *PocketSlaveStatus, pubkey []byte) (*PocketSlaveAgentMeta, error) {
    if pubkey == nil {
        return nil, errors.Errorf("[ERR] You cannot pass an empty pubkey")
    }
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketSlaveAgentMeta{
        MetaVersion:         SLAVE_META_VERSION,
        MasterBoundAgent:    master,
        SlaveID:             piface.HardwareAddr,
        StatusAgent:         agent,
        SlavePubKey:         pubkey,
    }, nil
}

func CheckSlaveCryptoMeta(master string, agent *PocketSlaveStatus, aescrypto pcrypto.AESCryptor) (*PocketSlaveAgentMeta, error) {
    mp, err := PackedSlaveStatus(agent)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    encrypted, err := aescrypto.EncryptByAES(mp)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketSlaveAgentMeta{
        MetaVersion:         SLAVE_META_VERSION,
        MasterBoundAgent:    master,
        SlaveID:             piface.HardwareAddr,
        EncryptedStatus:     encrypted,
    }, nil
}

func SlaveBoundedMeta(master string, agent *PocketSlaveStatus, aescrypto pcrypto.AESCryptor) (*PocketSlaveAgentMeta, error) {
    mp, err := PackedSlaveStatus(agent)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    encrypted, err := aescrypto.EncryptByAES(mp)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketSlaveAgentMeta{
        MetaVersion:         SLAVE_META_VERSION,
        MasterBoundAgent:    master,
        SlaveID:             piface.HardwareAddr,
        EncryptedStatus:     encrypted,
    }, nil
}

func BrokenBindMeta(master string) (*PocketSlaveAgentMeta, error) {
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketSlaveAgentMeta{
        MetaVersion:         SLAVE_META_VERSION,
        MasterBoundAgent:    master,
        SlaveID:             piface.HardwareAddr,
        DiscoveryAgent:      &PocketSlaveDiscovery {
            Version:             SLAVE_DISCOVER_VERSION,
            SlaveResponse:       SLAVE_LOOKUP_AGENT,
            SlaveAddress:        piface.PrimaryIP4Addr(),
            SlaveGateway:        piface.GatewayAddr,
        },
    }, nil
}
