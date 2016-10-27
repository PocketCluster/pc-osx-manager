package model

import (
    "os"
    "github.com/stkim1/pc-core/context"
)

func DebugModelRepoOpen() (ModelRepo) {
    context.DebugContextPrepared()

    // invalidate singleton instance
    singletonModelRepoInstance()
    repository = &modelRepo{}
    initializeModelRepo(repository)
    return repository
}

func DebugModelRepoClose() {
    CloseModelRepo()
    userDataPath, _ := context.SharedHostContext().ApplicationUserDataDirectory()
    os.Remove(userDataPath + "/core/pc-core.db")
    repository = nil
}

