package config

import (
    "os"
    "os/user"
)

func DebugConfigPrepare() (*PocketSlaveConfig, error) {
    usr, err := user.Current()
    if err != nil {
        return nil, err
    }
    // check if the path exists and make it if absent
    root := usr.HomeDir + "/temp"
    if _, err := os.Stat(root); os.IsNotExist(err) {
        os.MkdirAll(root, 0700);
    }

    return _loadSlaveConfig(root), nil
}

func DebugConfigDestory(cfg *PocketSlaveConfig) {
    os.Remove(cfg.rootPath + slave_config_file)
    os.Remove(cfg.rootPath + slave_public_Key_file)
    os.Remove(cfg.rootPath + slave_prvate_Key_file)
    os.Remove(cfg.rootPath + slave_ssh_Key_file)
    cfg = nil
}
