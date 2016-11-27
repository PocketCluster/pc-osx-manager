package slcontext

import (
    "log"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext/config"
)

func DebugSlcontextPrepare() PocketSlaveContext {
    // instead of running singleton creation, we'll invalidate sync.once to disengage in singleton production
    // getSingletonSlaveContext()
    once.Do(func(){})

    // load config and generate
    singletonContext = &slaveContext{}
    cfg, err := config.DebugConfigPrepare()
    if err != nil {
        log.Panic(err.Error())
    }
    err = singletonContext.initWithConfig(cfg)
    if err != nil {
        log.Panic(err.Error())
    }

    // pub/priv keys are generated
    singletonContext.pocketPublicKey    = pcrypto.TestSlavePublicKey()
    singletonContext.pocketPrivateKey   = pcrypto.TestSlavePrivateKey()
    singletonContext.nodePublicKey      = pcrypto.TestSlaveNodePublicKey()
    singletonContext.nodePrivateKey     = pcrypto.TestSlaveNodePrivateKey()

    return singletonContext
}

func DebugSlcontextPrepareWithRoot(rootPath string) PocketSlaveContext {
    // instead of running singleton creation, we'll invalidate sync.once to disengage in singleton production
    // getSingletonSlaveContext()
    once.Do(func(){})

    // load config and generate
    singletonContext = &slaveContext{}
    cfg, err := config.DebugConfigPrepareWithRoot(rootPath)
    if err != nil {
        log.Panic(err.Error())
    }
    err = singletonContext.initWithConfig(cfg)
    if err != nil {
        log.Panic(err.Error())
    }

    // pub/priv keys are generated
    singletonContext.pocketPublicKey    = pcrypto.TestSlavePublicKey()
    singletonContext.pocketPrivateKey   = pcrypto.TestSlavePrivateKey()
    singletonContext.nodePublicKey      = pcrypto.TestSlaveNodePublicKey()
    singletonContext.nodePrivateKey     = pcrypto.TestSlaveNodePrivateKey()

    return singletonContext
}

func DebugSlcontextDestroy() {
    singletonContext.DiscardAll()
    config.DebugConfigDestory(singletonContext.config)
    singletonContext = nil
}