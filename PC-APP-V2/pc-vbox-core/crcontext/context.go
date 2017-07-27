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

    // Discard master ip address, and other session related data
    DiscardMasterSession() error

    // This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bound.
    // No other place can execute this
    SaveConfiguration() error

    GetClusterID() (string, error)
    GetPrivateKey() ([]byte)
    GetPublicKey() ([]byte)
    GetMasterPublicKey() ([]byte)

    SetMasterIP4ExtAddr(ip4Address string) error
    GetMasterIP4ExtAddr() (string, error)

    PocketCoreSSHInfo
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
    // master public key
    c.masterPubkey, err = cfg.MasterPublicKey()
    if err != nil {
        return errors.WithStack(err)
    }
    return nil
}

// reload all configuration
func (c *coreContext) ReloadConfiguration() error {
    return initWithConfig(c, config.LoadPocketCoreConfig())
}

// Discard master ip address, and other session related data
func (c *coreContext) DiscardMasterSession() error {
    c.Lock()
    defer c.Unlock()

    c.config.MasterSection.MasterIP4Address = ""
    return nil
}

// This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
// No other place can execute this
func (c *coreContext) SaveConfiguration() error {
    return c.config.SaveCoreConfig()
}

// --- Cluster ID ---
func (c *coreContext) GetClusterID() (string, error) {
    if len(c.config.ClusterID) == 0 {
        return "", errors.Errorf("[ERR] cluster id name")
    }
    return c.config.ClusterID, nil
}

//--- decryptor/encryptor interface ---
func (c *coreContext) GetPrivateKey() ([]byte) {
    return c.pocketPrivateKey
}

func (c *coreContext) GetPublicKey() ([]byte) {
    return c.pocketPublicKey
}

// --- Master Public key ---
func (c *coreContext) GetMasterPublicKey() ([]byte) {
    return c.masterPubkey
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
