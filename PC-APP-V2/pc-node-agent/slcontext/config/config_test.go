package config

import (
    "os"
    "testing"
    "reflect"
    "io/ioutil"

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
    dirConfig := DirPathSlaveConfig(cfg.rootPath)
    if _, err := os.Stat(dirConfig); os.IsNotExist(err) {
        t.Error("[ERR] slave config dir should have existed")
        return
    }
    // check if config secure key dir also exists and creat if DNE
    dirCerts := DirPathSlaveCerts(cfg.rootPath)
    if _, err := os.Stat(dirCerts); os.IsNotExist(err) {
        t.Error("[ERR] slave keys dir should have existed")
        return
    }
    if err := cfg.SaveSlaveConfig(); err != nil {
        t.Error(err.Error())
        return
    }

    loadedCfg := loadSlaveConfig(cfg.rootPath)
    if !reflect.DeepEqual(cfg, loadedCfg) {
        t.Error("[ERR] incorect loaded config should be the same.")
        t.Log(spew.Sdump(loadedCfg))
        return
    }

    pubKey, err := cfg.SlavePublicKey()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if len(pubKey) == 0 {
        t.Error("[ERR] public key cannot be null")
        return
    }

    prvKey, err := cfg.SlavePrivateKey()
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


func TestSaveLoadHostName(t *testing.T) {
    cfg, err := DebugConfigPrepare()
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    defer DebugConfigDestory(cfg)

    var hostnamePath string = FilePathSystemHostname(cfg.rootPath)
    const slaveNodeName string = "jedi-skywalker"

    cfg.SlaveSection.SlaveNodeName = slaveNodeName
    err = cfg.SaveHostname()
    if err != nil {
        t.Error(err.Error())
        return
    }
    hname, err := ioutil.ReadFile(hostnamePath)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if string(hname) != slaveNodeName {
        t.Error("[ERR] incorrect hostname")
   }
}