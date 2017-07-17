package crcontext

import (
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-vbox-core/crcontext/config"
)

type PocketCoreContext interface {
    // reload all configuration
    ReloadConfiguration() error

    // Discard all data communicated with master (not the one from core itself such as network info)
    // This should executed on failure from joining states (unbounded, inquired, keyexchange, checkcrypto)
    DiscardAll() error

    // Discard master ip address, and other session related data
    DiscardMasterSession() error

    // This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bound.
    // No other place can execute this
    SaveConfiguration() error

    SetClusterID(clusterID string) error
    GetClusterID() (string, error)

    // authtoken
    SetCoreAuthToken(authToken string) error
    GetCoreAuthToken() (string, error)

    GetPrivateKey() (prvkey []byte)
    GetPublicKey() (pubkey []byte)

    SetMasterPublicKey(masterPubkey []byte) error
    GetMasterPublicKey() ([]byte, error)

    SetMasterIP4ExtAddr(ip4Address string) error
    GetMasterIP4ExtAddr() (string, error)

    CoreKeyAndCertPath() string
    CoreConfigPath() string
}

// Singleton handling
var (
    singletonContext *coreContext
    once sync.Once
)

type coreContext struct {
    sync.Mutex

    config           *config.PocketCoreConfig

    pocketPublicKey  []byte
    pocketPrivateKey []byte
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
    var (
        mpubkey []byte = nil
        err error = nil
    )
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
    mpubkey, err = cfg.MasterPublicKey()
    if len(mpubkey) != 0 && err == nil {
        c.masterPubkey = mpubkey
    }
    return nil
}

// reload all configuration
func (c *coreContext) ReloadConfiguration() error {
    return initWithConfig(c, config.LoadPocketCoreConfig())
}

// Discard all data communicated with master (not the one from slave itself such as network info)
// This should executed on failure from joining states (unbounded, inquired, keyexchange, checkcrypto)
func (c *coreContext) DiscardAll() error {
    // discard aeskey
    c.DiscardMasterSession()

    // remove decryptor
    c.masterPubkey = nil
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

func (c *coreContext) DiscardMasterSession() error {
    c.Lock()
    defer c.Unlock()

    c.config.MasterSection.MasterIP4Address = ""
    return nil
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

//--- decryptor/encryptor interface ---
func (c *coreContext) GetPrivateKey() ([]byte) {
    return c.pocketPrivateKey
}

func (c *coreContext) GetPublicKey() ([]byte) {
    return c.pocketPublicKey
}

// --- Master Public key ---
func (c *coreContext) SetMasterPublicKey(masterPubkey []byte) error {
    if len(masterPubkey) == 0 {
        return errors.Errorf("[ERR] invalid master public key")
    }
    c.masterPubkey = masterPubkey
    return nil
}

func (c *coreContext) GetMasterPublicKey() ([]byte, error) {
    if c.masterPubkey == nil {
        return nil, errors.Errorf("[ERR] empty master public key")
    }
    return c.masterPubkey, nil
}

// --- Master IP4 Address ---
func (c *coreContext) SetMasterIP4ExtAddr(ip4Address string) error {
    c.Lock()
    defer c.Unlock()

    if len(ip4Address) == 0 {
        return errors.Errorf("[ERR] invalid master ip4 address")
    }
    c.config.MasterSection.MasterIP4Address = ip4Address
    return nil
}

func (c *coreContext) GetMasterIP4ExtAddr() (string, error) {
    c.Lock()
    defer c.Unlock()

    if len(c.config.MasterSection.MasterIP4Address) == 0 {
        return "", errors.Errorf("[ERR] empty master ip4 address")
    }
    return c.config.MasterSection.MasterIP4Address , nil
}

// TODO : add tests
func (s *coreContext) CoreKeyAndCertPath() string {
    return s.config.KeyAndCertDir()
}

// TODO : add tests
func (s *coreContext) CoreConfigPath() string {
    return s.config.ConfigDir()
}
