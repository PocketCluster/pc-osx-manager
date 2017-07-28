package config

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"

    "github.com/pborman/uuid"
)

func (c *PocketSlaveConfig) DebugGetRootPath() string {
    return c.rootPath
}

func DebugConfigPrepare() (*PocketSlaveConfig, error) {
    var (
        // check if the path exists and make it if absent
        rootPath string = filepath.Join(os.TempDir(), uuid.New())
    )
    return DebugConfigPrepareWithRoot(rootPath)
}

func DebugConfigPrepareWithRoot(rootPath string) (*PocketSlaveConfig, error) {
    var (
        dirNetwork   string = filepath.Join(rootPath, "/etc/network")
        pathNetiface string = filepath.Join(rootPath, "/etc/network/interfaces")
        ifaces       []byte = []byte(`# interfaces(5) file used by ifup(8) and ifdown(8)
# Include files from /etc/network/interfaces.d:
source-directory /etc/network/interfaces.d

# The loopback network interface
auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp`)
    )

    if _, err := os.Stat(dirNetwork); os.IsNotExist(err) {
        if err := os.MkdirAll(dirNetwork, 0755); err != nil {
            return nil, err
        }
        fmt.Println(dirNetwork + " creation success.")
    } else {
        fmt.Println(dirNetwork + " exists.")
    }

    if _, err := os.Stat(pathNetiface); os.IsNotExist(err) {
        if err = ioutil.WriteFile(pathNetiface, ifaces, 0644); err != nil {
            return nil, err
        }
        fmt.Println(pathNetiface + "creation success.")
    } else {
        fmt.Println(pathNetiface + "exists.")
    }
    return loadSlaveConfig(rootPath), nil
}

func DebugConfigDestory(cfg *PocketSlaveConfig) {
    if cfg == nil {
        log.Panic("[CRITICAL] Configuration cannot be null")
    }

    os.RemoveAll(cfg.rootPath)
}
