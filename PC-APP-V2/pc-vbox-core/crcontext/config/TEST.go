package config

import (
    "os"
    "os/user"
    "log"

    "github.com/pborman/uuid"
)

func (pc *PocketCoreConfig) DebugGetRootPath() string {
    return pc.rootPath
}

func DebugConfigPrepare() (*PocketCoreConfig, error) {

    usr, err := user.Current()
    if err != nil {
        return nil, err
    }

    // TODO : ubuntu does not support TMPDIR env. CI host should have it's TMP as memfs in the future
    //root := os.Getenv("TMPDIR")

    var (
        tuid string = uuid.New()
        // check if the path exists and make it if absent
        root string     = usr.HomeDir + "/temp/" + tuid
    )

    return _loadCoreConfig(root), nil
}

func DebugConfigPrepareWithRoot(rootPath string) (*PocketCoreConfig, error) {
    return _loadCoreConfig(rootPath), nil
}

func DebugConfigDestory(cfg *PocketCoreConfig) {
    if cfg == nil {
        log.Panic("[CRITICAL] Configuration cannot be null")
    }

    os.RemoveAll(cfg.rootPath)

    // TODO : safely remove file. It seems file creation in test case create race condition *following does not work*
/*
    if _, err := os.Stat(cfg.rootPath); os.IsExist(err) {
        os.RemoveAll(cfg.rootPath)
    }
 */
}
