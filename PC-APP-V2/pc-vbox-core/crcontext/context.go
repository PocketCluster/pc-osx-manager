package crcontext

import (
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-vbox-core/crcontext/config"
)

// Singleton handling
var (
    singletonContext *coreContext
    once sync.Once
)

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
    var err error = nil

    c.PocketCoreConfig = cfg
    c.coreCertificate = new(coreCertificate)
    c.coreInstallImage = new(coreInstallImage)

    // pocket private key
    c.pocketPrivateKey, err = cfg.LoadCorePrivateKey()
    if err != nil {
        return errors.WithStack(err)
    }
    // pocket public key
    c.pocketPublicKey , err = cfg.LoadCorePublicKey()
    if err != nil {
        return errors.WithStack(err)
    }
    // master public key
    c.masterPubkey, err = cfg.LoadMasterPublicKey()
    if err != nil {
        return errors.WithStack(err)
    }
    return nil
}

type PocketCoreContext interface {
    // reload all configuration
    ReloadConfiguration() error

    // Discard master ip address, and other session related data
    DiscardMasterSession() error

    // This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bound.
    // No other place can execute this
    SaveConfiguration() error

    PocketCoreProperty
    PocketCertificate
    PocketCoreInstallImage
}

type coreContext struct {
    sync.Mutex

    *config.PocketCoreConfig
    *coreCertificate
    *coreInstallImage
}

// this method should never have an error
func SharedCoreContext() PocketCoreContext {
    return getSingletonCoreContext()
}

// reload all configuration
func (c *coreContext) ReloadConfiguration() error {
    return initWithConfig(c, config.LoadPocketCoreConfig())
}

// Discard master ip address, and other session related data
func (c *coreContext) DiscardMasterSession() error {
    c.Lock()
    defer c.Unlock()

    c.MasterSection.MasterIP4Address = ""
    return nil
}

// This must be executed on success from CheckCrypto -> Bound, or BindBroken -> Bind
// No other place can execute this
func (c *coreContext) SaveConfiguration() error {
    return c.SaveCoreConfig()
}
