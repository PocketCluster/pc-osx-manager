package msagent

import (
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
    return msgpack.Marshal(meta)
}

func UnpackedMasterMeta(message []byte) (meta *PocketMasterAgentMeta, err error) {
    err = msgpack.Unmarshal(message, &meta)
    return
}

// --- per-state meta function

func SlaveIdentityInquiryMeta(respond *PocketMasterRespond) (meta *PocketMasterAgentMeta) {
    meta = &PocketMasterAgentMeta {
        MetaVersion         :MASTER_META_VERSION,
        DiscoveryRespond    :respond,
    }
    return
}

func MasterDeclarationMeta(command *PocketMasterCommand, pubkey []byte) (meta *PocketMasterAgentMeta) {
    meta = &PocketMasterAgentMeta {
        MetaVersion         :MASTER_META_VERSION,
        StatusCommand       :command,
        MasterPubkey        :pubkey,
    }
    return
}

// AES key is encrypted with RSA for async encryption scheme, and rest of data, EncryptedMasterCommand &
// EncryptedSlaveStatus, are encrypted with AES
func ExchangeCryptoKeyAndNameMeta(command *PocketMasterCommand, status *slagent.PocketSlaveStatus, aeskey []byte, aescrypto pcrypto.AESCryptor, rsacrypto pcrypto.RsaEncryptor) (meta *PocketMasterAgentMeta, err error) {
    // marshal command
    mc, err := PackedMasterCommand(command)
    if err != nil {
        return
    }
    // encrypt the marshaled command with AES
    encryptedCommand, err := aescrypto.EncryptByAES(mc)
    if err != nil {
        return
    }

    //TODO : since including encrypted status bloats the final meta packet size to 633, we're here to omit it and put encrypted slave name instead. this should later be looked into again
/*
    // marshal status
    ms, err := msgpack.Marshal(status)
    if err != nil {
        return
    }
    // encrypt the marshaled status with AES
    encryptedStatus, err := aescrypto.Encrypt(ms)
    if err != nil {
        return
    }
*/
    // encrypted slave name with AES
    encryptedSlaveName, err := aescrypto.EncryptByAES([]byte(status.SlaveNodeName))
    if err != nil {
        return
    }
    // encrypt the AES key with RSA
    encryptedAES, AESsignature, err := rsacrypto.EncryptByRSA(aeskey)
    if err != nil {
        return
    }
    meta = &PocketMasterAgentMeta {
        MetaVersion             :MASTER_META_VERSION,
        EncryptedMasterCommand  :encryptedCommand,
        EncryptedSlaveStatus    :encryptedSlaveName, //encryptedStatus,
        EncryptedAESKey         :encryptedAES,
        RsaCryptoSignature      :AESsignature,
    }
    return
}

func MasterBindReadyMeta(command *PocketMasterCommand, aescrypto pcrypto.AESCryptor) (meta *PocketMasterAgentMeta, err error) {
    // marshal command
    mc, err := PackedMasterCommand(command)
    if err != nil {
        return
    }
    // encrypt the marshaled command with AES
    encryptedCommand, err := aescrypto.EncryptByAES(mc)
    if err != nil {
        return
    }
    meta = &PocketMasterAgentMeta {
        MetaVersion             :MASTER_META_VERSION,
        EncryptedMasterCommand  :encryptedCommand,
    }
    return
}

func BoundedSlaveAckMeta(command *PocketMasterCommand, aescrypto pcrypto.AESCryptor) (meta *PocketMasterAgentMeta, err error) {
    // marshal command
    mc, err := PackedMasterCommand(command)
    if err != nil {
        return
    }
    // encrypt the marshaled command with AES
    encryptedCommand, err := aescrypto.EncryptByAES(mc)
    if err != nil {
        return
    }
    meta = &PocketMasterAgentMeta {
        MetaVersion             :MASTER_META_VERSION,
        EncryptedMasterCommand  :encryptedCommand,
    }
    return
}

func BrokenBindRecoverMeta(respond *PocketMasterRespond, aeskey []byte, aescrypto pcrypto.AESCryptor, rsacrypto pcrypto.RsaEncryptor) (meta *PocketMasterAgentMeta, err error) {
    // marshal command
    mr, err := PackedMasterRespond(respond)
    if err != nil {
        return
    }
    // encrypt the marshaled command with AES
    er, err := aescrypto.EncryptByAES(mr)
    if err != nil {
        return
    }
    // encrypt the AES key with RSA
    ea, as, err := rsacrypto.EncryptByRSA(aeskey)
    if err != nil {
        return
    }
    meta = &PocketMasterAgentMeta {
        MetaVersion             :MASTER_META_VERSION,
        EncryptedMasterRespond  :er,
        EncryptedAESKey         :ea,
        RsaCryptoSignature      :as,
    }
    return
}
