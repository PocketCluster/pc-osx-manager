package slcontext

import (
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
    AESCryptor() (pcrypto.AESCryptor, error)
    pcrypto.AESCryptor

    SetMasterAgent(agentName string) error
    GetMasterAgent() (string, error)

    SetMasterIP4Address(ip4Address string) error
    GetMasterIP4Address() (string, error)

    // Discard master aes key, ip address, and other session related data
    DiscardMasterSession()

    SetSlaveNodeName(nodeName string) error
    GetSlaveNodeName() (string, error)
    GetSlaveNodeUUID() (string, error)

    // authtoken
    SetSlaveAuthToken(authToken string) error
    GetSlaveAuthToken() (string, error)

    SlaveKeyAndCertPath() string
    SlaveConfigPath() string
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
func initWithConfig(sc *slaveContext, cfg *config.PocketSlaveConfig) error {
    var err error
    sc.config = cfg

    // pocket public key
    sc.pocketPublicKey , err = cfg.SlavePublicKey()
    if err != nil {
        return errors.WithStack(err)
    }

    // pocket private key
    sc.pocketPrivateKey, err = cfg.SlavePrivateKey()
    if err != nil {
        return errors.WithStack(err)
    }

    // if master public key exists
    if pcmspubkey, err := cfg.MasterPublicKey(); len(pcmspubkey) != 0 && err == nil {
        sc.masterPubkey = pcmspubkey

        if decryptor, err := pcrypto.NewRsaDecryptorFromKeyData(pcmspubkey, sc.pocketPrivateKey); decryptor != nil && err == nil {
            sc.pocketDecryptor = decryptor
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
    if config.SLAVE_NAMESRV_VALUE != cfg.SlaveSection.SlaveNameServ {
        cfg.SlaveSection.SlaveNameServ = config.SLAVE_NAMESRV_VALUE
    }

    return nil
}

// Once sync, all the configuration is saved, and slave node is bounded
// This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
// No other place can execute this
func (sc *slaveContext) SyncAll() error {
    // slave network section
    paddr, err := PrimaryNetworkInterface()
    if err != nil {
        return err
    }
    cfg := sc.config
    if paddr.HardwareAddr != cfg.SlaveSection.SlaveMacAddr {
        cfg.SlaveSection.SlaveMacAddr  = paddr.HardwareAddr
    }
    if paddr.PrimaryIP4Addr() != cfg.SlaveSection.SlaveIP4Addr {
        cfg.SlaveSection.SlaveIP4Addr  = paddr.PrimaryIP4Addr()
    }
    if paddr.GatewayAddr != cfg.SlaveSection.SlaveGateway {
        cfg.SlaveSection.SlaveGateway  = paddr.GatewayAddr
    }
    if config.SLAVE_NAMESRV_VALUE != cfg.SlaveSection.SlaveNameServ {
        cfg.SlaveSection.SlaveNameServ = config.SLAVE_NAMESRV_VALUE
    }
    return nil
}

// Discard all data communicated with master (not the one from slave itself such as network info)
// This should executed on failure from joining states (unbounded, inquired, keyexchange, checkcrypto)
func (sc *slaveContext) DiscardAll() error {
    // discard aeskey
    sc.DiscardMasterSession()

    // remove decryptor
    sc.masterPubkey = nil
    sc.pocketDecryptor = nil
    // this is to remove master pub key if it exists
    if sc.config != nil {
        sc.config.ClearMasterPublicKey()
    }

    // master agent name
    sc.config.MasterSection.MasterBoundAgent    = ""
    // master ip4 address
    sc.config.MasterSection.MasterIP4Address    = ""
    // slave node name
    sc.config.SlaveSection.SlaveNodeName        = ""
    // slave auth token
    sc.config.SlaveSection.SlaveAuthToken       = ""
    return nil
}

// reload all configuration
func (sc *slaveContext) ReloadConfiguration() error {
    return initWithConfig(sc, config.LoadPocketSlaveConfig())
}

// This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
// No other place can execute this
func (sc *slaveContext) SaveConfiguration() error {
    // master pubkey
    mpubkey, err := sc.GetMasterPublicKey()
    if err != nil {
        return errors.WithStack(err)
    }
    sc.config.SaveMasterPublicKey(mpubkey)

    // TODO : (2017-05-15) we'll re-evaluate this option later. For not, this is none critical
/*
    // slave network interface
    err = sc.config.SaveFixedNetworkInterface()
    if err != nil {
        return errors.WithStack(err)
    }
*/
    // save slave node name to hostname
    err = sc.config.SaveHostname()
    if err != nil {
        return errors.WithStack(err)
    }
    // whole slave config
    return sc.config.SaveSlaveConfig()
}

// decryptor/encryptor interface
func (sc *slaveContext) GetPublicKey() ([]byte) {
    return sc.pocketPublicKey
}

func (sc *slaveContext) GetPrivateKey() ([]byte) {
    return sc.pocketPrivateKey
}

func (sc *slaveContext) DecryptByRSA(crypted []byte, sendSig pcrypto.Signature) ([]byte, error) {
    if sc.pocketDecryptor == nil {
        return nil, errors.Errorf("[ERR] cannot decrypt with null decryptor")
    }
    return sc.pocketDecryptor.DecryptByRSA(crypted, sendSig)
}

// --- Master Public key ---
func (sc *slaveContext) SetMasterPublicKey(masterPubkey []byte) error {
    if len(masterPubkey) == 0 {
        return errors.Errorf("[ERR] Master public key is nil")
    }
    sc.masterPubkey = masterPubkey

    decryptor, err := pcrypto.NewRsaDecryptorFromKeyData(masterPubkey, sc.pocketPrivateKey)
    if err != nil {
        return errors.WithStack(err)
    }
    sc.pocketDecryptor = decryptor
    return nil
}

func (sc *slaveContext) GetMasterPublicKey() ([]byte, error) {
    if sc.masterPubkey == nil {
        return nil, errors.Errorf("[ERR] Empty master public key")
    }
    return sc.masterPubkey, nil
}

// --- AES key ---
func (sc *slaveContext) SetAESKey(aesKey []byte) error {
    cryptor, err := pcrypto.NewAESCrypto(aesKey)
    if err != nil {
        return nil
    }
    sc.aesCryptor = cryptor
    sc.aeskey = aesKey
    return nil
}

func (sc *slaveContext) GetAESKey() ([]byte) {
    return sc.aeskey
}

func (sc *slaveContext) DiscardMasterSession() {
    sc.aesCryptor = nil
    sc.aeskey = nil
    sc.config.MasterSection.MasterIP4Address = ""
    return
}

func (sc *slaveContext) AESCryptor() (pcrypto.AESCryptor, error) {
    if sc.aesCryptor == nil {
        return nil, errors.Errorf("[ERR] AESKey or AESCryptor is not setup")
    }
    return sc.aesCryptor, nil
}

func (sc *slaveContext) EncryptByAES(data []byte) ([]byte, error) {
    if sc.aesCryptor == nil {
        return nil, errors.Errorf("[ERR] Cannot AES encrypt with null cryptor")
    }
    return sc.aesCryptor.EncryptByAES(data)
}

func (sc *slaveContext) DecryptByAES(data []byte) ([]byte, error) {
    if sc.aesCryptor == nil {
        return nil, errors.Errorf("[ERR] Cannot AES decrypt with null cryptor")
    }
    return sc.aesCryptor.DecryptByAES(data)
}

// --- Master Agent Name ---
func (sc *slaveContext) SetMasterAgent(agentName string) error {
    if len(agentName) == 0 {
        return errors.Errorf("[ERR] Cannot set empty master agent name")
    }
    sc.config.MasterSection.MasterBoundAgent = agentName
    return nil
}

func (sc *slaveContext) GetMasterAgent() (string, error) {
    if len(sc.config.MasterSection.MasterBoundAgent) == 0 {
        return "", errors.Errorf("[ERR] Empty master agent name")
    }
    return sc.config.MasterSection.MasterBoundAgent, nil
}

// --- Master IP4 Address ---
func (sc *slaveContext) SetMasterIP4Address(ip4Address string) error {
    if len(ip4Address) == 0 {
        return errors.Errorf("[ERR] Cannot set empty master ip4 address")
    }
    sc.config.MasterSection.MasterIP4Address = ip4Address
    return nil
}

func (sc *slaveContext) GetMasterIP4Address() (string, error) {
    if len(sc.config.MasterSection.MasterIP4Address) == 0 {
        return "", errors.Errorf("[ERR] Empty master ip4 address")
    }
    return sc.config.MasterSection.MasterIP4Address , nil
}

// --- Slave Node Name ---
func (sc *slaveContext) SetSlaveNodeName(nodeName string) error {
    if len(nodeName) == 0 {
        return errors.Errorf("[ERR] Cannot set empty slave nodename")
    }
    sc.config.SlaveSection.SlaveNodeName = nodeName
    return nil
}

func (sc *slaveContext) GetSlaveNodeName() (string, error) {
    if len(sc.config.SlaveSection.SlaveNodeName) == 0 {
        return "", errors.Errorf("[ERR] empty slave node name")
    }
    return sc.config.SlaveSection.SlaveNodeName, nil
}

// --- Slave Node UUID ---
func (sc *slaveContext) GetSlaveNodeUUID() (string, error) {
    if len(sc.config.SlaveSection.SlaveNodeUUID) == 0 {
        return "", errors.Errorf("[ERR] invalid slave node UUID")
    }
    return sc.config.SlaveSection.SlaveNodeUUID, nil
}

func (sc *slaveContext) SetSlaveAuthToken(authToken string) error {
    if len(authToken) == 0 {
        return errors.Errorf("[ERR] cannot assign invalid slave auth token")
    }
    sc.config.SlaveSection.SlaveAuthToken = authToken
    return nil
}

func (sc *slaveContext) GetSlaveAuthToken() (string, error) {
    if len(sc.config.SlaveSection.SlaveAuthToken) == 0 {
        return "", errors.Errorf("[ERR] invalid slave auth token")
    }
    return sc.config.SlaveSection.SlaveAuthToken, nil
}

// TODO : add tests
func (s *slaveContext) SlaveKeyAndCertPath() string {
    return s.config.KeyAndCertDir()
}

// TODO : add tests
func (s *slaveContext) SlaveConfigPath() string {
    return s.config.ConfigDir()
}
