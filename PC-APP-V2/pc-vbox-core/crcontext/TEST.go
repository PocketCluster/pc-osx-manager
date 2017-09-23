package crcontext

import (
    "log"
    "os"
    "path/filepath"

    "github.com/pborman/uuid"
    "github.com/stkim1/pc-vbox-core/crcontext/config"
)

func DebugPrepareCoreContext() PocketCoreContext {
    var (
        // check if the path exists and make it if absent
        rootPath string = filepath.Join(os.TempDir(), uuid.New())
    )
    return DebugPrepareCoreContextWithRoot(rootPath)
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

    return singletonContext
}

func DebugDestroyCoreContext() {
    singletonContext.DiscardMasterSession()
    config.DebugConfigDestory(singletonContext.PocketCoreConfig)
    singletonContext = nil
}