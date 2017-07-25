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

    loadedCfg := _loadCoreConfig(cfg.rootPath)
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
}

func TestConfigSaveReloadPublicMasterKey(t *testing.T) {
    cfg, err := DebugConfigPrepare()
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    defer DebugConfigDestory(cfg)

    err = cfg.SaveMasterPublicKey(pcrypto.TestMasterWeakPublicKey())
    if err != nil {
        t.Error(err.Error())
        return
    }

    master, err := cfg.MasterPublicKey()
    if err != nil {
        t.Error(err.Error())
        return
    }

    if !reflect.DeepEqual(master, pcrypto.TestMasterWeakPublicKey()) {
        t.Error("[ERR] Master Publickey is different!")
        return
    }
}
