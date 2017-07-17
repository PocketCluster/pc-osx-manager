package crcontext

import (
    "log"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-vbox-core/crcontext/config"
)

func DebugPrepareCoreContext() PocketCoreContext {
    // instead of running singleton creation, we'll invalidate sync.once to disengage in singleton production
    // getSingletonSlaveContext()
    once.Do(func(){})

    // load config and generate
    singletonContext = &coreContext{}
    cfg, err := config.DebugConfigPrepare()
    if err != nil {
        log.Panic(err.Error())
    }
    err = initWithConfig(singletonContext, cfg)
    if err != nil {
        log.Panic(err.Error())
    }

    // pub/priv keys are generated
    singletonContext.pocketPublicKey    = pcrypto.TestSlavePublicKey()
    singletonContext.pocketPrivateKey   = pcrypto.TestSlavePrivateKey()

    return singletonContext
}

func DebugPrepareCoreContextWithRoot(rootPath string) PocketCoreContext {
    // instead of running singleton creation, we'll invalidate sync.once to disengage in singleton production
    // getSingletonSlaveContext()
    once.Do(func(){})

    // load config and generate
    singletonContext = &coreContext{}
    cfg, err := config.DebugConfigPrepareWithRoot(rootPath)
    if err != nil {
        log.Panic(err.Error())
    }
    err = initWithConfig(singletonContext, cfg)
    if err != nil {
        log.Panic(err.Error())
    }

    // pub/priv keys are generated
    singletonContext.pocketPublicKey    = pcrypto.TestSlavePublicKey()
    singletonContext.pocketPrivateKey   = pcrypto.TestSlavePrivateKey()

    return singletonContext
}

func DebugDestroyCoreContext() {
    singletonContext.DiscardAll()
    config.DebugConfigDestory(singletonContext.config)
    singletonContext = nil
}