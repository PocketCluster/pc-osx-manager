package slcontext

import (
    "github.com/stkim1/pc-node-agent/crypt"
)

func DebugSlcontextPrepare() PocketSlaveContext {
    getSingletonSlaveContext()
    singletonContext = &slaveContext{}
    initializeSlaveContext(singletonContext)

    // pub/priv keys are generated
    singletonContext.publicKey = crypt.TestSlavePrivateKey()
    singletonContext.privateKey = crypt.TestSlavePrivateKey()

    return singletonContext
}

func DebugSlcontextDestroy() {
    singletonContext.DiscardAll()
    singletonContext = nil
}