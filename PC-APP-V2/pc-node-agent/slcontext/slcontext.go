package slcontext

import (
    "net"

    "github.com/stkim1/pcrypto"
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
    // This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
    // No other place can execute this
    SyncAll() error
    // Discard all data communicated with master (not the one from slave itself such as network info)
    // This should executed on failure from joining states (unbounded, inquired, keyexchange, checkcrypto)
    DiscardAll() error
    // reload all configuration

    // TODO : how to test this?
    // ReloadConfiguration() error

    // This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bound.
    // No other place can execute this
    SaveConfiguration() error
    GetPublicKey() (pubkey []byte)
    GetPrivateKey() (prvkey []byte)
    pcrypto.RsaDecryptor

    SetMasterPublicKey(masterPubkey []byte) error
    GetMasterPublicKey() ([]byte, error)

    SetAESKey(aesKey []byte) error
    // TODO : should this be removed? this is only used in testing. Plus, it does not handle error properlty
    GetAESKey() (aeskey []byte)
    DiscardAESKey()
    AESCryptor() (pcrypto.AESCryptor, error)
    pcrypto.AESCryptor

    SetMasterAgent(agentName string) error
    GetMasterAgent() (string, error)

    SetMasterIP4Address(ip4Address string) error
    GetMasterIP4Address() (string, error)

    SetSlaveNodeName(nodeName string) error
    GetSlaveNodeName() (string, error)

    PrimaryNetworkInterface() (*NetworkInterface, error)

    SlaveKeyAndCertPath() string
    SlaveConfigPath() string
}
