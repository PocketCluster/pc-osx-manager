package config

import (
    "log"
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/pborman/uuid"
    "github.com/stkim1/pcrypto"
)

func (pc *PocketCoreConfig) DebugGetRootPath() string {
    return pc.rootPath
}

func DebugConfigPrepare() (*PocketCoreConfig, error) {
    var (
        // check if the path exists and make it if absent
        rootPath string = filepath.Join(os.TempDir(), uuid.New())
    )
    return DebugConfigPrepareWithRoot(rootPath)
}

func DebugConfigPrepareWithRoot(rootPath string) (*PocketCoreConfig, error) {
    var (
        dirConfig        string = DirPathCoreConfig(rootPath)
        dirCerts         string = DirPathCoreCerts(rootPath)
        pathCorePubKey   string = FilePathCoreVboxPublicKey(rootPath)
        pathCorePrvKey   string = FilePathCoreVboxPrivateKey(rootPath)
        pathMasterPubKey string = FilePathMasterVboxPublicKey(rootPath)
    )

    os.MkdirAll(dirConfig, os.ModeDir|0700)
    os.MkdirAll(dirCerts,  os.ModeDir|0700)

    ioutil.WriteFile(pathCorePubKey,   pcrypto.TestSlaveNodePublicKey(),    0600)
    ioutil.WriteFile(pathCorePrvKey,   pcrypto.TestSlaveNodePrivateKey(),   0600)
    ioutil.WriteFile(pathMasterPubKey, pcrypto.TestMasterStrongPublicKey(), 0600)

    return _loadCoreConfig(rootPath), nil
}

func DebugConfigDestory(cfg *PocketCoreConfig) {
    if cfg == nil {
        log.Panic("[CRITICAL] Configuration cannot be null")
    }

    os.RemoveAll(cfg.rootPath)
}
