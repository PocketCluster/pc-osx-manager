package msagent

import (
    "github.com/pkg/errors"
    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slagent"
)

type PocketMasterAgentMeta struct {
    MetaVersion               MetaProtocol            `msgpack:"m_pm"`
    DiscoveryRespond          *PocketMasterRespond    `msgpack:"m_dr, inline, omitempty"`
    StatusCommand             *PocketMasterCommand    `msgpack:"m_sc, inline, omitempty"`
    EncryptedMasterCommand    []byte                  `msgpack:"m_ec, omitempty"`
    EncryptedSlaveStatus      []byte                  `msgpack:"m_es, omitempty"`
    MasterPubkey              []byte                  `msgpack:"m_pk, omitempty"`
    EncryptedAESKey           []byte                  `msgpack:"m_ak, omitempty"`
    RsaCryptoSignature        []byte                  `msgpack:"m_sg, omitempty"`
    EncryptedMasterRespond    []byte                  `msgpack:"m_er, omitempty"`
}


func PackedMasterMeta(meta *PocketMasterAgentMeta) ([]byte, error) {
    m, err := msgpack.Marshal(meta)
    return m, errors.WithStack(err)
}

func UnpackedMasterMeta(message []byte) (meta *PocketMasterAgentMeta, err error) {
    err = errors.WithStack(msgpack.Unmarshal(message, &meta))
    return
}

// --- per-state meta function

func SlaveIdentityInquiryMeta(respond *PocketMasterRespond) (meta *PocketMasterAgentMeta) {
    meta = &PocketMasterAgentMeta {
        MetaVersion:         MASTER_META_VERSION,
        DiscoveryRespond:    respond,
    }
    return
}

func MasterDeclarationMeta(command *PocketMasterCommand, pubkey []byte) (meta *PocketMasterAgentMeta) {
    meta = &PocketMasterAgentMeta {
        MetaVersion:         MASTER_META_VERSION,
        StatusCommand:       command,
        MasterPubkey:        pubkey,
    }
    return
}

// AES key is encrypted with RSA for async encryption scheme, and rest of data, EncryptedMasterCommand &
// EncryptedSlaveStatus, are encrypted with AES
func ExchangeCryptoKeyAndNameMeta(command *PocketMasterCommand, slaveIdentity *slagent.PocketSlaveIdentity, aeskey []byte, aescrypto pcrypto.AESCryptor, rsacrypto pcrypto.RsaEncryptor) (*PocketMasterAgentMeta, error) {
    // marshal command
    mc, err := PackedMasterCommand(command)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    // encrypt the marshaled command with AES
    encryptedCommand, err := aescrypto.EncryptByAES(mc)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // packed slave name & uuid
    pslid, err := slagent.PackPocketSlaveIdentity(slaveIdentity)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    // encrypt the marshaled status with AES
    eslid, err := aescrypto.EncryptByAES(pslid)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // encrypt the AES key with RSA
    encryptedAES, AESsignature, err := rsacrypto.EncryptByRSA(aeskey)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketMasterAgentMeta {
        MetaVersion:               MASTER_META_VERSION,
        EncryptedMasterCommand:    encryptedCommand,
        EncryptedSlaveStatus:      eslid,
        EncryptedAESKey:           encryptedAES,
        RsaCryptoSignature:        AESsignature,
    }, nil
}

func MasterBindReadyMeta(command *PocketMasterCommand, aescrypto pcrypto.AESCryptor) (*PocketMasterAgentMeta, error) {
    // marshal command
    mc, err := PackedMasterCommand(command)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    // encrypt the marshaled command with AES
    encryptedCommand, err := aescrypto.EncryptByAES(mc)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketMasterAgentMeta {
        MetaVersion:              MASTER_META_VERSION,
        EncryptedMasterCommand:   encryptedCommand,
    }, nil
}

func BoundedSlaveAckMeta(command *PocketMasterCommand, aescrypto pcrypto.AESCryptor) (*PocketMasterAgentMeta, error) {
    // marshal command
    mc, err := PackedMasterCommand(command)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    // encrypt the marshaled command with AES
    encryptedCommand, err := aescrypto.EncryptByAES(mc)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketMasterAgentMeta {
        MetaVersion             :MASTER_META_VERSION,
        EncryptedMasterCommand  :encryptedCommand,
    }, nil
}

func BrokenBindRecoverMeta(respond *PocketMasterRespond, aeskey []byte, aescrypto pcrypto.AESCryptor, rsacrypto pcrypto.RsaEncryptor) (*PocketMasterAgentMeta, error) {
    // marshal command
    mr, err := PackedMasterRespond(respond)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    // encrypt the marshaled command with AES
    er, err := aescrypto.EncryptByAES(mr)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    // encrypt the AES key with RSA
    ea, as, err := rsacrypto.EncryptByRSA(aeskey)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketMasterAgentMeta {
        MetaVersion:               MASTER_META_VERSION,
        EncryptedMasterRespond:    er,
        EncryptedAESKey:           ea,
        RsaCryptoSignature:        as,
    }, nil
}
