package model

import "os"

func DebugModelRepoOpen() (ModelRepo) {
    singletonModelRepoInstance()

    // invalidate singleton instance
    repository = &modelRepo{}
    initializeModelRepo(repository)
    return repository
}

func DebugModelRepoClose() {
    CloseModelRepo()
    repository = nil
    os.Remove("pc-core.db")
}

