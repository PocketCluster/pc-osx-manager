package slcontext

import (
    "github.com/stkim1/pc-node-agent/crypt"
    "net"
)

type NetworkInterface struct {
    *net.Interface
    *net.IP
    *net.IPMask
    *net.HardwareAddr
    GatewayAddr             string
}

type PocketSlaveContext interface {
    // Once sync, all the configuration is saved, and slave node is bounded
    // This must be executed on success from CheckCrypto -> Bound
    SyncAll() error
    // Discard all data communicated.
    // This should executed on failure from joining states (unbounded, inquired, keyexchange, checkcrypto)
    DiscardAll() error
    // reload all configuration
    ReloadConfiguration() error
    // Save all configuration
    SaveConfiguration() error

    GetPublicKey() (pubkey []byte)
    GetPrivateKey() (prvkey []byte)
    crypt.RsaDecryptor

    SetMasterPublicKey(masterPubkey []byte) error
    GetMasterPublicKey() ([]byte, error)

    SetAESKey(aesKey []byte) error
    GetAESKey() (aeskey []byte)
    DiscardAESKey()
    AESCryptor() (crypt.AESCryptor, error)
    crypt.AESCryptor

    SetMasterAgent(agentName string) error
    GetMasterAgent() (string, error)

    SetMasterIP4Address(ip4Address string) error
    GetMasterIP4Address() (string, error)

    SetSlaveNodeName(nodeName string) error
    GetSlaveNodeName() (string, error)

    PrimaryNetworkInterface() (*NetworkInterface, error)
}
