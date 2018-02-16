package config

import (
    "fmt"
    "io/ioutil"
    "path/filepath"
    "strings"

    "github.com/pkg/errors"
)

// ------ NETWORK INTERFACES ------
const (
    SLAVE_ADDRESS_KEY      string = "address"
    SLAVE_GATEWAY_KEY      string = "gateway"
    SLAVE_NETMASK_KEY      string = "netmask"
    SLAVE_NAMESRV_KEY      string = "dns-nameserver"
    SLAVE_BROADCS_KEY      string = "broadcast"
)

var SLAVE_NETIFACE_KEYS  []string = []string{SLAVE_ADDRESS_KEY, SLAVE_GATEWAY_KEY, SLAVE_NETMASK_KEY, SLAVE_NAMESRV_KEY, SLAVE_BROADCS_KEY}

// ------- POCKET EDITOR MARKER ------
const (
    POCKET_START           string = "# --------------- POCKETCLUSTER START ---------------"
    POCKET_END             string = "# ---------------  POCKETCLUSTER END  ---------------"
)

// ------ SALT DEFAULT ------
const (
    PC_MASTER              string = "pc-master"
)

const (
    network_iface_file     string = "/etc/network/interfaces"
)

// --- network interface redefinition
func _network_iface_redefined(slaveConfig *SlaveConfigSection) []string {
    return []string {
        POCKET_START,
        "iface eth0 inet static",
        fmt.Sprintf("%s %s", SLAVE_ADDRESS_KEY, slaveConfig.SlaveIP4Addr),
        fmt.Sprintf("%s %s", SLAVE_GATEWAY_KEY, slaveConfig.SlaveGateway),
        fmt.Sprintf("%s %s", SLAVE_NAMESRV_KEY, slaveConfig.SlaveNameServ),
        //fmt.Sprintf("%s %s", SLAVE_BROADCS_KEY, slaveConfig.SlaveBroadcast),
        POCKET_END,
    }
}

func _network_keyword_prefixed(line string) bool {
    for _, exl := range SLAVE_NETIFACE_KEYS {
        if strings.HasPrefix(line, exl) {
            return true
        }
    }
    return false
}

func _fixateNetworkInterfaces(slaveConfig *SlaveConfigSection, ifaceData []string) []string {
    var ifacelines []string
    is_pocket_defiend := false
    is_pocket_editing := false
    // first scan
    for _, l := range ifaceData {
        line := strings.TrimSpace(l)
        if strings.HasPrefix(l, POCKET_START) {
            is_pocket_defiend = true
            is_pocket_editing = true
            ifacelines = append(ifacelines, _network_iface_redefined(slaveConfig)...)
            continue
        }
        if strings.HasPrefix(l, POCKET_END) {
            is_pocket_editing = false
            continue
        }
        if !is_pocket_editing {
            ifacelines = append(ifacelines, line)
        }
    }
    // second scan in case there is no pocket section
    if !is_pocket_defiend {
        ifacelines = nil
        for _, l := range ifaceData {
            line := strings.TrimSpace(l)
            if strings.HasPrefix(line, "iface eth0 inet") {
                ifacelines = append(ifacelines, _network_iface_redefined(slaveConfig)...)
                continue
            }
            if _network_keyword_prefixed(line) {
                continue
            }
            ifacelines = append(ifacelines, line)
        }
    }
    return ifacelines
}

func (c *PocketSlaveConfig) SaveFixedNetworkInterface() error {
    var (
        ifaceFilePath = filepath.Join(c.rootPath, network_iface_file)
        ifaceFileContent, err = ioutil.ReadFile(ifaceFilePath)
    )
    if err != nil {
        return errors.WithStack(err)
    }
    var (
        ifaceData = strings.Split(string(ifaceFileContent),"\n")
        fixedIfaceData = _fixateNetworkInterfaces(c.SlaveSection, ifaceData)
        fixedIfaceContent = []byte(strings.Join(fixedIfaceData, "\n"))
    )
    return ioutil.WriteFile(ifaceFilePath, fixedIfaceContent, 0644)
}
