package slcontext

import "github.com/stkim1/pc-node-agent/crypt"

type PocketSlaveContext interface {
    GetPublicKey() (pubkey []byte)
    GetPrivateKey() (prvkey []byte)
    crypt.RsaDecryptor
    crypt.RsaEncryptor

    SetMasterPublicKey(masterPubkey []byte) (err error)
    GetMasterPublicKey() (pubkey []byte, err error)
    DiscardMasterPublicKey() (err error)
    SyncMasterPublicKey() (err error)

    SetAESKey(aesKey []byte) (err error)
    GetAESKey() (aeskey []byte)
    DiscardAESKey()
    crypt.AESCryptor

    SetMasterAgent(agentName string) (err error)
    GetMasterAgent() (agentName string, err error)
    SyncMasterAgent() (err error)
    DiscardMasterAgent() (err error)

    SetMasterIP4Address(ip4Address string) (err error)
    GetMasterIP4Address() (ip4Address string, err error)
    SyncMasterIP4Address() (err error)
    DiscardMasterIP4Address() (err error)

    SetSlaveNodeName(nodeName string) (err error)
    GetSlaveNodeName() (nodeName string, err error)
    SyncSlaveNodeName() (err error)
    DiscardSlaveNodeName() (err error)

    SyncAll() (err error)
    DiscardAll() (err error)
}
