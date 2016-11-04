package slcontext

import (
    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-node-agent/slcontext/config"
    "fmt"
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
    singletonContext.publicKey = crypt.TestSlavePrivateKey()
    singletonContext.privateKey = crypt.TestSlavePrivateKey()

    return singletonContext
}

func DebugSlcontextDestroy() {
    singletonContext.DiscardAll()
    singletonContext = nil
}