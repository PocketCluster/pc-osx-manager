package config

import (
    "testing"
    "github.com/davecgh/go-spew/spew"
    "reflect"
    "os"
)

func TestConfigLoadAndSave(t *testing.T) {
    cfg, err := DebugConfigPrepare()
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    defer DebugConfigDestory(cfg)
    t.Log(spew.Sdump(cfg))

    // check if config dir exists, and creat if DNE
    configDirPath := cfg.rootPath + slave_config_dir
    if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
        t.Error("[ERR] slave config dir should have existed")
        return
    }
    // check if config secure key dir also exists and creat if DNE
    keysDirPath := cfg.rootPath + slave_keys_dir
    if _, err := os.Stat(keysDirPath); os.IsNotExist(err) {
        t.Error("[ERR] slave keys dir should have existed")
        return
    }
    if err := cfg.Save(); err != nil {
        t.Error(err.Error())
        return
    }

    loadedCfg := _loadSlaveConfig(cfg.rootPath)
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

    sshkey, err := cfg.SlaveSSHKey()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if len(sshkey) == 0 {
        t.Error("[ERR] ssh key cannot be null")
        return
    }

}