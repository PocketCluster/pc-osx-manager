package config

import (
    "os"
    "os/user"
    "io/ioutil"
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
        if err = ioutil.WriteFile(netiface, ifaces, 0644); err != nil {
            return nil, err
        }
    }
    return _loadSlaveConfig(root), nil
}

func DebugConfigDestory(cfg *PocketSlaveConfig) {

    os.Remove(cfg.rootPath + slave_config_file)
    os.Remove(cfg.rootPath + slave_public_Key_file)
    os.Remove(cfg.rootPath + slave_prvate_Key_file)
    os.Remove(cfg.rootPath + slave_ssh_Key_file)
    os.Remove(cfg.rootPath + master_public_Key_file)

    os.Remove(cfg.rootPath + hostname_file)
    os.Remove(cfg.rootPath + network_iface_file)
    cfg = nil
}
