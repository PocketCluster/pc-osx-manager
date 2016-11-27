package config

import (
    "os"
    "os/user"
    "io/ioutil"
    "fmt"
    "log"

    "github.com/pborman/uuid"
)

func (pc *PocketSlaveConfig) DebugGetRootPath() string {
    return pc.rootPath
}

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
    var (
        tuid string = uuid.New()
        // check if the path exists and make it if absent
        root string     = usr.HomeDir + "/temp/" + tuid
        // check if network directory is ready
        network string  = root + "/etc/network"
        // network interface file
        netiface string = root + "/etc/network/interfaces"
    )

    if _, err := os.Stat(network); os.IsNotExist(err) {
        if err := os.MkdirAll(network, 0755); err != nil {
            return nil, err
        }
        fmt.Println(root + "/etc/network/ creation success.")
    }

    if _, err := os.Stat(netiface); os.IsNotExist(err) {
        if err = ioutil.WriteFile(netiface, ifaces, 0644); err != nil {
            return nil, err
        }
        fmt.Println(root + "/etc/network/interfaces creation success.")
    }
    return _loadSlaveConfig(root), nil
}

func DebugConfigPrepareWithRoot(rootPath string) (*PocketSlaveConfig, error) {

    ifaces := []byte(`# interfaces(5) file used by ifup(8) and ifdown(8)
# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d

# The loopback network interface
auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp`)

    var (
        // check if the path exists and make it if absent
        root string     = rootPath
        // check if network directory is ready
        network string  = rootPath + "/etc/network"
        // network interface file
        netiface string = rootPath + "/etc/network/interfaces"
    )

    if _, err := os.Stat(network); os.IsNotExist(err) {
        if err := os.MkdirAll(network, 0755); err != nil {
            return nil, err
        }
        fmt.Println(root + "/etc/network/ creation success.")
    } else {
        fmt.Println(root + "/etc/network/ creation exists.")
    }

    if _, err := os.Stat(netiface); os.IsNotExist(err) {
        if err = ioutil.WriteFile(netiface, ifaces, 0644); err != nil {
            return nil, err
        }
        fmt.Println(root + "/etc/network/interfaces creation success.")
    } else {
        fmt.Println(root + "/etc/network/interfaces exists.")
    }
    return _loadSlaveConfig(root), nil
}


func DebugConfigDestory(cfg *PocketSlaveConfig) {
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
