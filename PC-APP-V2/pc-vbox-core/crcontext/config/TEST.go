package config

import (
    "log"
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/pborman/uuid"
    "github.com/stkim1/pcrypto"
)

const (
    TestCoreUserName    string = "almightykim"
    TestClusterID       string = "p3l4WI26Bd50hzAo"
    TestAuthToken       string = "22acc6140aa95e69c9bfd6ed778645f9"
)

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
        pathClusterID    string = FilePathClusterID(rootPath)
        pathAuthToken    string = FilePathAuthToken(rootPath)

        dirCerts         string = DirPathCoreCerts(rootPath)
        pathCorePubKey   string = FilePathCoreVboxPublicKey(rootPath)
        pathCorePrvKey   string = FilePathCoreVboxPrivateKey(rootPath)
        pathMasterPubKey string = FilePathMasterVboxPublicKey(rootPath)
    )

    os.MkdirAll(dirConfig,             os.ModeDir|0700)
    ioutil.WriteFile(pathClusterID,    []byte(TestClusterID),               0600)
    ioutil.WriteFile(pathAuthToken,    []byte(TestAuthToken),               0600)

    os.MkdirAll(dirCerts,              os.ModeDir|0700)
    ioutil.WriteFile(pathCorePubKey,   pcrypto.TestSlaveNodePublicKey(),    0600)
    ioutil.WriteFile(pathCorePrvKey,   pcrypto.TestSlaveNodePrivateKey(),   0600)
    ioutil.WriteFile(pathMasterPubKey, pcrypto.TestMasterStrongPublicKey(), 0600)

    return loadCoreConfig(rootPath), nil
}

func DebugConfigDestory(cfg *PocketCoreConfig) {
    if cfg == nil {
        log.Panic("[CRITICAL] Configuration cannot be null")
    }

    os.RemoveAll(cfg.rootPath)
}
