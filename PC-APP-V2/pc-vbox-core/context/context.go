package context

import (
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-vbox-core/context/config"
)

type PocketCoreContext interface {
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
    
    SetClusterID(clusterID string) error
    GetClusterID() (string, error)

    SetMasterIP4Address(ip4Address string) error
    GetMasterIP4Address() (string, error)

    // Discard master aes key, ip address, and other session related data
    DiscardMasterSession()

    // authtoken
    SetCoreAuthToken(authToken string) error
    GetCoreAuthToken() (string, error)

    CoreKeyAndCertPath() string
    CoreConfigPath() string
}

// Singleton handling
var (
    singletonContext *coreContext
    once sync.Once
)

type coreContext struct {
    config           *config.PocketCoreConfig

    pocketPublicKey  []byte
    pocketPrivateKey []byte
    pocketDecryptor  pcrypto.RsaDecryptor

    masterPubkey     []byte
}

// this method should never have an error
func SharedSlaveContext() PocketCoreContext {
    return getSingletonSlaveContext()
}

func getSingletonSlaveContext() *coreContext {
    once.Do(func() {
        var (
            cfg *config.PocketCoreConfig = nil
            err error = nil
        )
        singletonContext = &coreContext{}
        cfg = config.LoadPocketCoreConfig()
        err = initWithConfig(singletonContext, cfg)
        if err != nil {
            // TODO : Trace this log
            log.Panicf("[CRITICAL] %s", errors.WithStack(err).Error())
        }
    })
    return singletonContext
}

// --- Sync All ---
func initWithConfig(sc *coreContext, cfg *config.PocketCoreConfig) error {
    var err error
    sc.config = cfg

    // pocket private key
    sc.pocketPrivateKey, err = cfg.CorePrivateKey()
    if err != nil {
        return errors.WithStack(err)
    }
    // pocket public key
    sc.pocketPublicKey , err = cfg.CorePublicKey()
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

    return nil
}

// Once sync, all the configuration is saved, and slave node is bounded
// This must be executed on success from Unbounded -> Bound, or BindBroken -> Bind
// No other place can execute this
func (sc *coreContext) SyncAll() error {
    return nil
}

// Discard all data communicated with master (not the one from slave itself such as network info)
// This should executed on failure from joining states (unbounded, inquired, keyexchange, checkcrypto)
func (sc *coreContext) DiscardAll() error {
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
    sc.config.ClusterID = ""
    // slave auth token
    sc.config.CoreSection.CoreAuthToken = ""
    return nil
}

// reload all configuration
func (sc *coreContext) ReloadConfiguration() error {
    return initWithConfig(sc, config.LoadPocketCoreConfig())
}

// This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
// No other place can execute this
func (sc *coreContext) SaveConfiguration() error {
    // master pubkey
    mpubkey, err := sc.GetMasterPublicKey()
    if err != nil {
        return errors.WithStack(err)
    }
    sc.config.SaveMasterPublicKey(mpubkey)

    return sc.config.SaveCoreConfig()
}

// decryptor/encryptor interface
func (sc *coreContext) GetPublicKey() ([]byte) {
    return sc.pocketPublicKey
}

func (sc *coreContext) GetPrivateKey() ([]byte) {
    return sc.pocketPrivateKey
}

func (sc *coreContext) DecryptByRSA(crypted []byte, sendSig pcrypto.Signature) ([]byte, error) {
    if sc.pocketDecryptor == nil {
        return nil, errors.Errorf("[ERR] cannot decrypt with null decryptor")
    }
    return sc.pocketDecryptor.DecryptByRSA(crypted, sendSig)
}

// --- Master Public key ---
func (sc *coreContext) SetMasterPublicKey(masterPubkey []byte) error {
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

func (sc *coreContext) GetMasterPublicKey() ([]byte, error) {
    if sc.masterPubkey == nil {
        return nil, errors.Errorf("[ERR] Empty master public key")
    }
    return sc.masterPubkey, nil
}

// --- Master Agent Name ---
func (sc *coreContext) SetClusterID(clusterID string) error {
    if len(clusterID) == 0 {
        return errors.Errorf("[ERR] Cannot set empty master agent name")
    }
    sc.config.ClusterID = clusterID
    return nil
}

func (sc *coreContext) GetClusterID() (string, error) {
    if len(sc.config.ClusterID) == 0 {
        return "", errors.Errorf("[ERR] Empty master agent name")
    }
    return sc.config.ClusterID, nil
}

// --- Master IP4 Address ---
func (sc *coreContext) SetMasterIP4Address(ip4Address string) error {
    if len(ip4Address) == 0 {
        return errors.Errorf("[ERR] Cannot set empty master ip4 address")
    }
    sc.config.MasterSection.MasterIP4Address = ip4Address
    return nil
}

func (sc *coreContext) GetMasterIP4Address() (string, error) {
    if len(sc.config.MasterSection.MasterIP4Address) == 0 {
        return "", errors.Errorf("[ERR] Empty master ip4 address")
    }
    return sc.config.MasterSection.MasterIP4Address , nil
}

// --- Slave Node UUID ---
func (sc *coreContext) SetCoreAuthToken(authToken string) error {
    if len(authToken) == 0 {
        return errors.Errorf("[ERR] cannot assign invalid slave auth token")
    }
    sc.config.CoreSection.CoreAuthToken = authToken
    return nil
}

func (sc *coreContext) GetCoreAuthToken() (string, error) {
    if len(sc.config.CoreSection.CoreAuthToken) == 0 {
        return "", errors.Errorf("[ERR] invalid slave auth token")
    }
    return sc.config.CoreSection.CoreAuthToken, nil
}

// TODO : add tests
func (s *coreContext) CoreKeyAndCertPath() string {
    return s.config.KeyAndCertDir()
}

// TODO : add tests
func (s *coreContext) CoreConfigPath() string {
    return s.config.ConfigDir()
}
