package msagent

import (
    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-node-agent/slagent"
)

type PocketMasterAgentMeta struct {
    MetaVersion            MetaProtocol                       `msgpack:"pc_ms_pm"`
    DiscoveryRespond       *PocketMasterDiscoveryRespond      `msgpack:"pc_ms_dr, inline, omitempty"`
    StatusCommand          *PocketMasterStatusCommand         `msgpack:"pc_ms_sc, inline, omitempty"`
    EncryptedMasterCommand []byte                             `msgpack:"pc_ms_ec, omitempty"`
    EncryptedSlaveStatus   []byte                             `msgpack:"pc_ms_es, omitempty"`
    MasterPubkey           []byte                             `msgpack:"pc_ms_pk, omitempty"`
    EncryptedAESKey        []byte                             `msgpack:"pc_ms_ak, omitempty"`
    RsaCryptoSignature     []byte                             `msgpack:"pc_ms_sg, omitempty"`
}


func PackedMasterMeta(meta *PocketMasterAgentMeta) ([]byte, error) {
    return msgpack.Marshal(meta)
}

func UnpackedMasterMeta(message []byte) (*PocketMasterAgentMeta, error) {
    var meta *PocketMasterAgentMeta
    err := msgpack.Unmarshal(message, &meta)
    if err != nil {
        return nil, err
    }
    return meta, nil
}

func UnboundedInqueryMeta(respond *PocketMasterDiscoveryRespond) (meta *PocketMasterAgentMeta) {
    meta = &PocketMasterAgentMeta{
        MetaVersion         :MASTER_META_VERSION,
        DiscoveryRespond    :respond,
    }
    return
}

func IdentityInqueryMeta(command *PocketMasterStatusCommand, pubkey []byte) (meta *PocketMasterAgentMeta) {
    meta = &PocketMasterAgentMeta{
        MetaVersion         :MASTER_META_VERSION,
        StatusCommand       :command,
        MasterPubkey        :pubkey,
    }
    return
}

// AES key is encrypted with RSA for async encryption scheme, and rest of data, EncryptedMasterCommand &
// EncryptedSlaveStatus, are encrypted with AES
func ExecKeyExchangeMeta(command *PocketMasterStatusCommand, status *slagent.PocketSlaveStatusAgent, aeskey []byte, aescrypto crypt.AESCryptor, rsacrypto crypt.RsaEncryptor) (meta *PocketMasterAgentMeta, err error) {
    // marshal command
    mc, err := msgpack.Marshal(command)
    if err != nil {
        return nil, err
    }
    // encrypt the marshaled command with AES
    encryptedCommand, err := aescrypto.Encrypt(mc)
    if err != nil {
        return nil, err
    }

    //TODO : since including encrypted status bloats the final meta packet size to 633, we're here to omit it and put encrypted slave name instead
    //TODO : this should later be looked into again
/*
    // marshal status
    ms, err := msgpack.Marshal(status)
    if err != nil {
        return nil, err
    }
    // encrypt the marshaled status with AES
    encryptedStatus, err := aescrypto.Encrypt(ms)
    if err != nil {
        return nil, err
    }
*/
    // encrypted slave name with AES
    encryptedSlaveName, err := aescrypto.Encrypt([]byte(status.SlaveNodeName))
    if err != nil {
        return nil, err
    }
    // encrypt the AES key with RSA
    encryptedAES, AESsignature, err := rsacrypto.EncryptMessage(aeskey)
    if err != nil {
        return nil, err
    }
    meta = &PocketMasterAgentMeta{
        MetaVersion             :MASTER_META_VERSION,
        EncryptedMasterCommand  :encryptedCommand,
        EncryptedSlaveStatus    :encryptedSlaveName, //encryptedStatus,
        EncryptedAESKey         :encryptedAES,
        RsaCryptoSignature      :AESsignature,
    }
    return
}

func SendCryptoCheckMeta(command *PocketMasterStatusCommand, aescrypto crypt.AESCryptor) (meta *PocketMasterAgentMeta, err error) {
    // marshal command
    mc, err := msgpack.Marshal(command)
    if err != nil {
        return nil, err
    }
    // encrypt the marshaled command with AES
    encryptedCommand, err := aescrypto.Encrypt(mc)
    if err != nil {
        return nil, err
    }
    meta = &PocketMasterAgentMeta{
        MetaVersion             :MASTER_META_VERSION,
        EncryptedMasterCommand  :encryptedCommand,
    }
    return
}

func BoundedStatusMeta(command *PocketMasterStatusCommand, aescrypto crypt.AESCryptor) (meta *PocketMasterAgentMeta, err error) {
    // marshal command
    mc, err := msgpack.Marshal(command)
    if err != nil {
        return nil, err
    }
    // encrypt the marshaled command with AES
    encryptedCommand, err := aescrypto.Encrypt(mc)
    if err != nil {
        return nil, err
    }
    meta = &PocketMasterAgentMeta{
        MetaVersion             :MASTER_META_VERSION,
        EncryptedMasterCommand  :encryptedCommand,
    }
    return
}

func BindBrokenMeta() (meta *PocketMasterAgentMeta, err error) {
    return
}
