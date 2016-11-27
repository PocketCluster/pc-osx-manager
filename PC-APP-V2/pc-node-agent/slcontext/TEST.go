package slcontext

import (
    "fmt"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext/config"
)

func DebugSlcontextPrepare() PocketSlaveContext {
    getSingletonSlaveContext()
    singletonContext = &slaveContext{}
    cfg, err := config.DebugConfigPrepare()
    if err != nil {
        fmt.Print(err.Error())
    }
    initializeSlaveContext(singletonContext, cfg)
    // pub/priv keys are generated
    singletonContext.pocketPublicKey = pcrypto.TestSlavePublicKey()
    singletonContext.pocketPrivateKey = pcrypto.TestSlavePrivateKey()

    return singletonContext
}

func DebugSlcontextDestroy() {
    singletonContext.DiscardAll()
    config.DebugConfigDestory(singletonContext.config)
    singletonContext = nil
}