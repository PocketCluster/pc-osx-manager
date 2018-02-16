package slcontext

import (
    "fmt"
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext/config"
)

type PocketSlaveContext interface {
    // Once sync, all the configuration is saved, and slave node is bounded
    // This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
    // No other place can execute this
    SyncAll() error

    // This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bound.
    // No other place can execute this
    SaveConfiguration() error

    // Discard all data communicated with master (not the one from slave itself such as network info)
    // This should executed on failure from joining states (unbounded, inquired, keyexchange, checkcrypto)
    DiscardAll() error

    // Discard master aes key, ip address, and other session related data
    DiscardMasterSession()

    GetPublicKey() (pubkey []byte)
    GetPrivateKey() (prvkey []byte)
    pcrypto.RsaDecryptor

    SetMasterPublicKey(masterPubkey []byte) error
    GetMasterPublicKey() ([]byte, error)

    SetAESKey(aesKey []byte) error
    // TODO : should this be removed? this is only used in testing. Plus, it does not handle error properlty
    GetAESKey() (aeskey []byte)
    AESCryptor() (pcrypto.AESCryptor, error)
    pcrypto.AESCryptor

    SetClusterID(clusterID string) error
    GetClusterID() (string, error)

    SetMasterIP4Address(ip4Address string) error
    GetMasterIP4Address() (string, error)

    SetSlaveNodeName(nodeName string) error
    GetSlaveNodeName() (string, error)
    GetSlaveNodeNameFQDN() (string, error)

    PocketSlaveSSHInfo
}

// Singleton handling
var (
    singletonContext *slaveContext
    once sync.Once
)

type slaveContext struct {
    config           *config.PocketSlaveConfig

    pocketPublicKey  []byte
    pocketPrivateKey []byte
    pocketDecryptor  pcrypto.RsaDecryptor

    masterPubkey     []byte
    aeskey           []byte
    aesCryptor       pcrypto.AESCryptor
}

// this method should never have an error
func SharedSlaveContext() PocketSlaveContext {
    return getSingletonSlaveContext()
}

func getSingletonSlaveContext() *slaveContext {
    once.Do(func() {
        singletonContext = &slaveContext{}
        cfg := config.LoadPocketSlaveConfig()
        err := initWithConfig(singletonContext, cfg)
        if err != nil {
            // TODO : Trace this log
            log.Panicf("[CRITICAL] %s", err.Error())
        }
    })
    return singletonContext
}

// --- Sync All ---
func initWithConfig(s *slaveContext, cfg *config.PocketSlaveConfig) error {
    var err error
    s.config = cfg

    // pocket public key
    s.pocketPublicKey , err = cfg.SlavePublicKey()
    if err != nil {
        return errors.WithStack(err)
    }

    // pocket private key
    s.pocketPrivateKey, err = cfg.SlavePrivateKey()
    if err != nil {
        return errors.WithStack(err)
    }

    // if master public key exists
    if pcmspubkey, err := cfg.MasterPublicKey(); len(pcmspubkey) != 0 && err == nil {
        s.masterPubkey = pcmspubkey

        if decryptor, err := pcrypto.NewRsaDecryptorFromKeyData(pcmspubkey, s.pocketPrivateKey); decryptor != nil && err == nil {
            s.pocketDecryptor = decryptor
        }
    }

    paddr, err := PrimaryNetworkInterface()
    if err != nil {
        return errors.WithStack(err)
    }

    if paddr.HardwareAddr != cfg.SlaveSection.SlaveMacAddr {
        cfg.SlaveSection.SlaveMacAddr  = paddr.HardwareAddr
    }
    if paddr.PrimaryIP4Addr() != cfg.SlaveSection.SlaveIP4Addr {
        cfg.SlaveSection.SlaveIP4Addr  = paddr.PrimaryIP4Addr()
    }
    if paddr.GatewayAddr != cfg.SlaveSection.SlaveGateway {
        cfg.SlaveSection.SlaveGateway  = paddr.GatewayAddr
    }

    return nil
}

// Once sync, all the configuration is saved, and slave node is bounded
// This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
// No other place can execute this
func (s *slaveContext) SyncAll() error {
    // slave network section
    paddr, err := PrimaryNetworkInterface()
    if err != nil {
        return err
    }
    cfg := s.config
    if paddr.HardwareAddr != cfg.SlaveSection.SlaveMacAddr {
        cfg.SlaveSection.SlaveMacAddr  = paddr.HardwareAddr
    }
    if paddr.PrimaryIP4Addr() != cfg.SlaveSection.SlaveIP4Addr {
        cfg.SlaveSection.SlaveIP4Addr  = paddr.PrimaryIP4Addr()
    }
    if paddr.GatewayAddr != cfg.SlaveSection.SlaveGateway {
        cfg.SlaveSection.SlaveGateway  = paddr.GatewayAddr
    }
    return nil
}

// Discard all data communicated with master (not the one from slave itself such as network info)
// This should executed on failure from joining states (unbounded, inquired, keyexchange, checkcrypto)
func (s *slaveContext) DiscardAll() error {
    // discard aeskey
    s.DiscardMasterSession()

    // remove decryptor
    s.masterPubkey = nil
    s.pocketDecryptor = nil
    // this is to remove master pub key if it exists
    if s.config != nil {
        s.config.ClearMasterPublicKey()
    }

    // master agent name
    s.config.ClusterID                         = ""
    // master ip4 address
    s.config.MasterSection.MasterIP4Address    = ""
    // slave node name
    s.config.SlaveSection.SlaveNodeName        = ""
    // slave auth token
    s.config.SlaveSection.SlaveAuthToken       = ""
    return nil
}

// reload all configuration
func (s *slaveContext) ReloadConfiguration() error {
    return initWithConfig(s, config.LoadPocketSlaveConfig())
}

// This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
// No other place can execute this
func (s *slaveContext) SaveConfiguration() error {
    // master pubkey
    if mpubkey, err := s.GetMasterPublicKey(); err != nil {
        return errors.WithStack(err)
    } else {
        s.config.SaveMasterPublicKey(mpubkey)
    }
    // save slave node name to hostname
    if err := s.config.SaveHostname(); err != nil {
        return errors.WithStack(err)
    }
    // update hosts
    if err := s.config.UpdateHostsFile(); err != nil {
        return errors.WithStack(err)
    }
/*
    // slave network interface
    TODO : (2017-05-15) we'll re-evaluate this option later. For not, this is none critical
    if err = s.config.SaveFixedNetworkInterface(); err != nil {
        return errors.WithStack(err)
    }
*/
    // save slave config into yaml
    return s.config.SaveSlaveConfig()
}

// decryptor/encryptor interface
func (s *slaveContext) GetPublicKey() ([]byte) {
    return s.pocketPublicKey
}

func (s *slaveContext) GetPrivateKey() ([]byte) {
    return s.pocketPrivateKey
}

func (s *slaveContext) DecryptByRSA(crypted []byte, sendSig pcrypto.Signature) ([]byte, error) {
    if s.pocketDecryptor == nil {
        return nil, errors.Errorf("[ERR] cannot decrypt with null decryptor")
    }
    return s.pocketDecryptor.DecryptByRSA(crypted, sendSig)
}

// --- Master Public key ---
func (s *slaveContext) SetMasterPublicKey(masterPubkey []byte) error {
    if len(masterPubkey) == 0 {
        return errors.Errorf("[ERR] Master public key is nil")
    }
    s.masterPubkey = masterPubkey

    decryptor, err := pcrypto.NewRsaDecryptorFromKeyData(masterPubkey, s.pocketPrivateKey)
    if err != nil {
        return errors.WithStack(err)
    }
    s.pocketDecryptor = decryptor
    return nil
}

func (s *slaveContext) GetMasterPublicKey() ([]byte, error) {
    if s.masterPubkey == nil {
        return nil, errors.Errorf("[ERR] Empty master public key")
    }
    return s.masterPubkey, nil
}

// --- AES key ---
func (s *slaveContext) SetAESKey(aesKey []byte) error {
    cryptor, err := pcrypto.NewAESCrypto(aesKey)
    if err != nil {
        return nil
    }
    s.aesCryptor = cryptor
    s.aeskey = aesKey
    return nil
}

func (s *slaveContext) GetAESKey() ([]byte) {
    return s.aeskey
}

func (s *slaveContext) DiscardMasterSession() {
    s.aesCryptor = nil
    s.aeskey = nil
    s.config.MasterSection.MasterIP4Address = ""
    return
}

func (s *slaveContext) AESCryptor() (pcrypto.AESCryptor, error) {
    if s.aesCryptor == nil {
        return nil, errors.Errorf("[ERR] AESKey or AESCryptor is not setup")
    }
    return s.aesCryptor, nil
}

func (s *slaveContext) EncryptByAES(data []byte) ([]byte, error) {
    if s.aesCryptor == nil {
        return nil, errors.Errorf("[ERR] Cannot AES encrypt with null cryptor")
    }
    return s.aesCryptor.EncryptByAES(data)
}

func (s *slaveContext) DecryptByAES(data []byte) ([]byte, error) {
    if s.aesCryptor == nil {
        return nil, errors.Errorf("[ERR] Cannot AES decrypt with null cryptor")
    }
    return s.aesCryptor.DecryptByAES(data)
}

// --- Master Agent Name ---
func (s *slaveContext) SetClusterID(clusterID string) error {
    if len(clusterID) == 0 {
        return errors.Errorf("[ERR] Cannot set empty cluster id")
    }
    s.config.ClusterID = clusterID
    return nil
}

func (s *slaveContext) GetClusterID() (string, error) {
    if len(s.config.ClusterID) == 0 {
        return "", errors.Errorf("[ERR] Empty master agent name")
    }
    return s.config.ClusterID, nil
}

// --- Master IP4 Address ---
func (s *slaveContext) SetMasterIP4Address(ip4Address string) error {
    if len(ip4Address) == 0 {
        return errors.Errorf("[ERR] Cannot set empty master ip4 address")
    }
    s.config.MasterSection.MasterIP4Address = ip4Address
    return nil
}

func (s *slaveContext) GetMasterIP4Address() (string, error) {
    if len(s.config.MasterSection.MasterIP4Address) == 0 {
        return "", errors.Errorf("[ERR] Empty master ip4 address")
    }
    return s.config.MasterSection.MasterIP4Address , nil
}

// --- Slave Node Name ---
func (s *slaveContext) SetSlaveNodeName(nodeName string) error {
    if len(nodeName) == 0 {
        return errors.Errorf("[ERR] Cannot set empty slave nodename")
    }
    s.config.SlaveSection.SlaveNodeName = nodeName
    return nil
}

func (s *slaveContext) GetSlaveNodeName() (string, error) {
    if len(s.config.SlaveSection.SlaveNodeName) == 0 {
        return "", errors.Errorf("[ERR] empty slave node name")
    }
    return s.config.SlaveSection.SlaveNodeName, nil
}
func (s *slaveContext) GetSlaveNodeNameFQDN() (string, error) {
    node, err := s.GetSlaveNodeName()
    if err != nil {
        return "", err
    }
    cid, err := s.GetClusterID()
    if err != nil {
        return "", err
    }

    return fmt.Sprintf(node + "." + pcrypto.FormFQDNClusterID, cid), nil
}