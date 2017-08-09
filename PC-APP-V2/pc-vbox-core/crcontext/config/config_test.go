package config

import (
    "os"
    "testing"
    "reflect"
    "github.com/davecgh/go-spew/spew"
    "github.com/stkim1/pcrypto"
)

func TestConfigSaveReload(t *testing.T) {
    cfg, err := DebugConfigPrepare()
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    defer DebugConfigDestory(cfg)
    t.Log(spew.Sdump(cfg))

    // check if config dir exists, and creat if DNE
    dirConfig := DirPathCoreConfig(cfg.rootPath)
    if _, err := os.Stat(dirConfig); os.IsNotExist(err) {
        t.Error("[ERR] slave config dir should have existed")
        return
    }
    // check if config secure key dir also exists and creat if DNE
    dirCerts := DirPathCoreCerts(cfg.rootPath)
    if _, err := os.Stat(dirCerts); os.IsNotExist(err) {
        t.Error("[ERR] slave keys dir should have existed")
        return
    }
    if err := cfg.SaveCoreConfig(); err != nil {
        t.Error(err.Error())
        return
    }

    loadedCfg := loadCoreConfig(cfg.rootPath)
    if !reflect.DeepEqual(cfg, loadedCfg) {
        t.Error("[ERR] incorect loaded config should be the same.")
        t.Log(spew.Sdump(loadedCfg))
        return
    }

    pubKey, err := cfg.CorePublicKey()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if len(pubKey) == 0 {
        t.Error("[ERR] public key cannot be null")
        return
    }

    prvKey, err := cfg.CorePrivateKey()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if len(prvKey) == 0 {
        t.Error("[ERR] private key cannot be null")
        return
    }
    master, err := cfg.MasterPublicKey()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if !reflect.DeepEqual(master, pcrypto.TestMasterStrongPublicKey()) {
        t.Error("[ERR] Master Publickey is different!")
        return
    }

    if cfg.ClusterID != TestClusterID {
        t.Error("[ERR] cluster id is different")
        return
    }

    if cfg.CoreSection.CoreAuthToken != TestAuthToken {
        t.Error("[ERR] auth token is different")
        return
    }
}

func TestConfigMultiReload(t *testing.T) {
    cfg1, err := DebugConfigPrepare()
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    defer DebugConfigDestory(cfg1)

    cfg2 := loadCoreConfig(cfg1.RootPath())
    if !reflect.DeepEqual(cfg1, cfg2) {
        t.Errorf("[ERR] 2nd configuration should be identical without crash")
    }

    cfg3 := loadCoreConfig(cfg1.RootPath())
    if !reflect.DeepEqual(cfg2, cfg3) {
        t.Errorf("[ERR] third configuration should be identical without crash")
    }
}