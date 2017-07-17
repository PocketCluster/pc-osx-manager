package crcontext

import (
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-vbox-core/crcontext/config"
)

type PocketCoreContext interface {
    // Once sync, all the configuration is saved, and core node is bounded
    // This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
    // No other place can execute this
    SyncAll() error

    // Discard all data communicated with master (not the one from core itself such as network info)
    // This should executed on failure from joining states (unbounded, inquired, keyexchange, checkcrypto)
    DiscardAll() error

    // Discard master ip address, and other session related data
    DiscardMasterSession()

    // reload all configuration
    ReloadConfiguration() error

    // This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bound.
    // No other place can execute this
    SaveConfiguration() error

    SetClusterID(clusterID string) error
    GetClusterID() (string, error)

    GetPublicKey() (pubkey []byte)
    GetPrivateKey() (prvkey []byte)
    pcrypto.RsaDecryptor

    SetMasterPublicKey(masterPubkey []byte) error
    GetMasterPublicKey() ([]byte, error)

    SetMasterIP4Address(ip4Address string) error
    GetMasterIP4Address() (string, error)

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
func SharedCoreContext() PocketCoreContext {
    return getSingletonCoreContext()
}

func getSingletonCoreContext() *coreContext {
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
func initWithConfig(c *coreContext, cfg *config.PocketCoreConfig) error {
    var err error
    c.config = cfg

    // pocket private key
    c.pocketPrivateKey, err = cfg.CorePrivateKey()
    if err != nil {
        return errors.WithStack(err)
    }
    // pocket public key
    c.pocketPublicKey , err = cfg.CorePublicKey()
    if err != nil {
        return errors.WithStack(err)
    }

    // if master public key exists
    if pcmspubkey, err := cfg.MasterPublicKey(); len(pcmspubkey) != 0 && err == nil {
        c.masterPubkey = pcmspubkey

        if decryptor, err := pcrypto.NewRsaDecryptorFromKeyData(pcmspubkey, c.pocketPrivateKey); decryptor != nil && err == nil {
            c.pocketDecryptor = decryptor
        }
    }

    return nil
}

// Once sync, all the configuration is saved, and slave node is bounded
// This must be executed on success from Unbounded -> Bound, or BindBroken -> Bind
// No other place can execute this
func (c *coreContext) SyncAll() error {
    return nil
}

// Discard all data communicated with master (not the one from slave itself such as network info)
// This should executed on failure from joining states (unbounded, inquired, keyexchange, checkcrypto)
func (c *coreContext) DiscardAll() error {
    // discard aeskey
    c.DiscardMasterSession()

    // remove decryptor
    c.masterPubkey = nil
    c.pocketDecryptor = nil
    // this is to remove master pub key if it exists
    if c.config != nil {
        c.config.ClearMasterPublicKey()
    }
    // master agent name
    c.config.ClusterID = ""
    // slave auth token
    c.config.CoreSection.CoreAuthToken = ""
    return nil
}

func (c *coreContext) DiscardMasterSession() {
    c.config.MasterSection.MasterIP4Address = ""
    return
}

// reload all configuration
func (c *coreContext) ReloadConfiguration() error {
    return initWithConfig(c, config.LoadPocketCoreConfig())
}

// This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
// No other place can execute this
func (c *coreContext) SaveConfiguration() error {
    // master pubkey
    mpubkey, err := c.GetMasterPublicKey()
    if err != nil {
        return errors.WithStack(err)
    }
    c.config.SaveMasterPublicKey(mpubkey)

    return c.config.SaveCoreConfig()
}

// --- Master Agent Name ---
func (c *coreContext) SetClusterID(clusterID string) error {
    if len(clusterID) == 0 {
        return errors.Errorf("[ERR] invalid cluster id to set")
    }
    c.config.ClusterID = clusterID
    return nil
}

func (c *coreContext) GetClusterID() (string, error) {
    if len(c.config.ClusterID) == 0 {
        return "", errors.Errorf("[ERR] cluster id name")
    }
    return c.config.ClusterID, nil
}

//--- decryptor/encryptor interface ---
func (c *coreContext) GetPublicKey() ([]byte) {
    return c.pocketPublicKey
}

func (c *coreContext) GetPrivateKey() ([]byte) {
    return c.pocketPrivateKey
}

func (c *coreContext) DecryptByRSA(crypted []byte, sendSig pcrypto.Signature) ([]byte, error) {
    if c.pocketDecryptor == nil {
        return nil, errors.Errorf("[ERR] cannot decrypt with null decryptor")
    }
    return c.pocketDecryptor.DecryptByRSA(crypted, sendSig)
}

// --- Master Public key ---
func (c *coreContext) SetMasterPublicKey(masterPubkey []byte) error {
    if len(masterPubkey) == 0 {
        return errors.Errorf("[ERR] Master public key is nil")
    }
    c.masterPubkey = masterPubkey

    decryptor, err := pcrypto.NewRsaDecryptorFromKeyData(masterPubkey, c.pocketPrivateKey)
    if err != nil {
        return errors.WithStack(err)
    }
    c.pocketDecryptor = decryptor
    return nil
}

func (c *coreContext) GetMasterPublicKey() ([]byte, error) {
    if c.masterPubkey == nil {
        return nil, errors.Errorf("[ERR] Empty master public key")
    }
    return c.masterPubkey, nil
}

// --- Master IP4 Address ---
func (c *coreContext) SetMasterIP4Address(ip4Address string) error {
    if len(ip4Address) == 0 {
        return errors.Errorf("[ERR] Cannot set empty master ip4 address")
    }
    c.config.MasterSection.MasterIP4Address = ip4Address
    return nil
}

func (c *coreContext) GetMasterIP4Address() (string, error) {
    if len(c.config.MasterSection.MasterIP4Address) == 0 {
        return "", errors.Errorf("[ERR] Empty master ip4 address")
    }
    return c.config.MasterSection.MasterIP4Address , nil
}

// --- Auth Token ---
func (c *coreContext) SetCoreAuthToken(authToken string) error {
    if len(authToken) == 0 {
        return errors.Errorf("[ERR] cannot assign invalid core auth token")
    }
    c.config.CoreSection.CoreAuthToken = authToken
    return nil
}

func (c *coreContext) GetCoreAuthToken() (string, error) {
    if len(c.config.CoreSection.CoreAuthToken) == 0 {
        return "", errors.Errorf("[ERR] invalid core auth token")
    }
    return c.config.CoreSection.CoreAuthToken, nil
}

// TODO : add tests
func (s *coreContext) CoreKeyAndCertPath() string {
    return s.config.KeyAndCertDir()
}

// TODO : add tests
func (s *coreContext) CoreConfigPath() string {
    return s.config.ConfigDir()
}
