package model

import (
    "os"
    "github.com/stkim1/pc-core/context"
)

func DebugModelRepoPrepare() (ModelRepo) {
    context.DebugContextPrepare()

    // invalidate singleton instance
    singletonModelRepoInstance()
    repository = &modelRepo{}
    initializeModelRepo(repository)
    return repository
}

func DebugModelRepoDestroy() {
    CloseModelRepo()
    userDataPath, _ := context.SharedHostContext().ApplicationUserDataDirectory()
    os.Remove(userDataPath + "/core/pc-core.db")
    repository = nil
}

