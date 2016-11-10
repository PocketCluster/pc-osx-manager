package config

import (
    "os"
    "os/user"
    "io/ioutil"
    "fmt"
)

func DebugConfigPrepare() (*PocketSlaveConfig, error) {

    ifaces := []byte(`# interfaces(5) file used by ifup(8) and ifdown(8)
# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d

# The loopback network interface
auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp`)

    usr, err := user.Current()
    if err != nil {
        return nil, err
    }
    // check if the path exists and make it if absent
    root := usr.HomeDir + "/temp"
    if _, err := os.Stat(root); os.IsNotExist(err) {
        os.MkdirAll(root, 0700);
    }
    // check if network directory is ready
    network := root + "/etc/network"
    if _, err := os.Stat(network); os.IsNotExist(err) {
        os.MkdirAll(network, 0755);
    }

    netiface := root + "/etc/network/interfaces"
    if _, err := os.Stat(netiface); os.IsNotExist(err) {
        fmt.Println("/etc/network/interfaces DNE. Let's make one")
        if err = ioutil.WriteFile(netiface, ifaces, 0644); err != nil {
            return nil, err
        }
        fmt.Println("/etc/network/interfaces creation success.")
    }
    return _loadSlaveConfig(root), nil
}

func DebugConfigDestory(cfg *PocketSlaveConfig) {
    if _, err := os.Stat(cfg.rootPath + slave_config_file); os.IsExist(err) {
        os.Remove(cfg.rootPath + slave_config_file)
    }
    if _, err := os.Stat(cfg.rootPath + slave_public_Key_file); os.IsExist(err) {
        os.Remove(cfg.rootPath + slave_public_Key_file)
    }
    if _, err := os.Stat(cfg.rootPath + slave_prvate_Key_file); os.IsExist(err) {
        os.Remove(cfg.rootPath + slave_prvate_Key_file)
    }
    if _, err := os.Stat(cfg.rootPath + slave_ssh_Key_file); os.IsExist(err) {
        os.Remove(cfg.rootPath + slave_ssh_Key_file)
    }
    if _, err := os.Stat(cfg.rootPath + master_public_Key_file); os.IsExist(err) {
        os.Remove(cfg.rootPath + master_public_Key_file)
    }
    if _, err := os.Stat(cfg.rootPath + hostname_file); os.IsExist(err) {
        os.Remove(cfg.rootPath + hostname_file)
    }
    if _, err := os.Stat(cfg.rootPath + network_iface_file); os.IsExist(err) {
        os.Remove(cfg.rootPath + network_iface_file)
    }
    cfg = nil
}
