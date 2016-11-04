package config

import (
    "os"
    "fmt"
    "io/ioutil"
    "strings"

    "gopkg.in/yaml.v2"
    "github.com/stkim1/pc-node-agent/crypt"
)

// ------ CONFIG VERSION -------
const (
    SLAVE_CONFIG_KEY    = "config-version"
    SLAVE_CONFIG_VAL    = "1.0.1"
)

const (
    SLAVE_STATUS_KEY    = "binding-status"
)

// ------ NETWORK INTERFACES ------
const (
    SLAVE_ADDRESS_KEY = "address"
    SLAVE_GATEWAY_KEY = "gateway"
    SLAVE_NETMASK_KEY = "netmask"
    SLAVE_NAMESRV_KEY = "dns-nameserver"
    SLAVE_BROADCS_KEY = "broadcast"
)

var SLAVE_NETIFACE_KEYS []string = []string{SLAVE_ADDRESS_KEY, SLAVE_GATEWAY_KEY, SLAVE_NETMASK_KEY, SLAVE_NAMESRV_KEY, SLAVE_BROADCS_KEY}

const (
    PAGENT_SEND_PORT    = 10060
    PAGENT_RECV_PORT    = 10061
)

// Master name server is fixed for now (v.0.1.4)
const SLAVE_NAMESRV_VALUE = "pc-master:53535"

// ------ CONFIGURATION FILES ------
const (
    // POCKET SPECIFIC CONFIG
    slave_config_dir        = "/etc/pocket/"
    slave_config_file       = "/etc/pocket/slave-conf.yaml"

    slave_keys_dir          = "/etc/pocket/pki"
    slave_prvate_Key_file   = "/etc/pocket/pki/slave.pem"
    slave_public_Key_file   = "/etc/pocket/pki/slave.pub"
    slave_ssh_Key_file      = "/etc/pocket/pki/slave.ssh"
    master_public_Key_file  = "/etc/pocket/pki/master.pub"

    // HOST GENERAL CONFIG
    network_iface_file      = "/etc/network/interfaces"
    hostname_file           = "/etc/hostname"
    //hostaddr_file           = "/etc/hosts"
    host_timezone_file      = "/etc/timezone"
    //resolve_conf_file       = "/etc/resolv.conf"
)

// ------ SALT DEFAULT ------
const (
    PC_MASTER           = "pc-master"
)

// ------- POCKET EDITOR MARKER ------
const (
    POCKET_START        = "# --------------- POCKETCLUSTER START ---------------"
    POCKET_END          = "# ---------------  POCKETCLUSTER END  ---------------"
)

// --- struct
type ConfigMasterSection struct {
    // Master Agent Specific String
    MasterBoundAgent    string                   `yaml:"master-binder-agent"`
    // Last Known IP4
    MasterIP4Address    string                   `yaml:"master-ip4-addr"`
    //MasterIP6Address    string
    //MasterHostName      string
    MasterTimeZone      string                   `yaml:"master-timezone"`
}

type ConfigSlaveSection struct {
    SlaveNodeName       string                   `yaml:"slave-node-name"`
    SlaveMacAddr        string                   `yaml:"slave-mac-addr"`
    //SlaveIP6Addr        string
    SlaveIP4Addr        string                   `yaml:"slave-net-ip4"`
    SlaveGateway        string                   `yaml:"slave-net-gateway"`
    SlaveNetMask        string                   `yaml:"slave-net-mask"`
    SlaveNameServ       string                   `yaml:"slave-net-nameserv"`
    //SlaveBroadcast      string                   `yaml:"slave-net-broadc"`
}

type PocketSlaveConfig struct {
    // this field exists to create files at a specific location for testing so ignore
    rootPath            string                   `yaml:"-"`
    ConfigVersion       string                   `yaml:"config-version"`
    BindingStatus       string                   `yaml:"binding-status"`
    MasterSection       *ConfigMasterSection     `yaml:"master-section",inline,flow`
    SlaveSection        *ConfigSlaveSection      `yaml:"slave-section",inline,flow`
}

// This is default public constructor as it does not accept root file path
func LoadPocketSlaveConfig() *PocketSlaveConfig {
    return _loadSlaveConfig("")
}

// --- func
func _brandNewSlaveConfig(rootPath string) (*PocketSlaveConfig) {
    return &PocketSlaveConfig{
        rootPath        :rootPath,
        ConfigVersion   :SLAVE_CONFIG_VAL,
        // TODO : we need to avoid cyclic import but need to fix this
        BindingStatus   : "SlaveUnbounded", //locator.SlaveUnbounded.String(),
        MasterSection   :&ConfigMasterSection{},
        SlaveSection    :&ConfigSlaveSection{},
    }
}

func _loadSlaveConfig(rootPath string) (*PocketSlaveConfig) {

    // check if config dir exists, and creat if DNE
    configDirPath := rootPath + slave_config_dir
    if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
        os.MkdirAll(configDirPath, 0700);
    }
    // check if config secure key dir also exists and creat if DNE
    keysDirPath := rootPath + slave_keys_dir
    if _, err := os.Stat(keysDirPath); os.IsNotExist(err) {
        os.MkdirAll(keysDirPath, 0700);
    }

    var shouldGenerateKeys bool = false
    pubKeyPath := rootPath + slave_public_Key_file
    prvKeyPath := rootPath + slave_prvate_Key_file
    sshKeyPath := rootPath + slave_ssh_Key_file
    if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
        shouldGenerateKeys = true
    }
    if _, err := os.Stat(prvKeyPath); os.IsNotExist(err) {
        shouldGenerateKeys = true
    }
    if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
        shouldGenerateKeys = true
    }
    if shouldGenerateKeys {
        crypt.GenerateKeyPair(pubKeyPath, prvKeyPath, sshKeyPath)
    }

    // check if config file exists in path.
    configFilePath := rootPath + slave_config_file
    if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
        return _brandNewSlaveConfig(rootPath)
    }

    // if does, unmarshal and load them.
    if configData, err := ioutil.ReadFile(configFilePath); err != nil {
        return _brandNewSlaveConfig(rootPath)
    } else {
        var config PocketSlaveConfig
        if err = yaml.Unmarshal(configData, &config); err != nil {
            return _brandNewSlaveConfig(rootPath)
        } else {
            // as rootpath is ignored, we need to restore it
            config.rootPath = rootPath
            return &config
        }
    }
}

func (cfg *PocketSlaveConfig) SaveSlaveConfig() error {
    // check if config dir exists, and creat if DNE
    configDirPath := cfg.rootPath + slave_config_dir
    if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
        os.MkdirAll(configDirPath, 0700);
    }

    configFilePath := cfg.rootPath + slave_config_file
    configData, err := yaml.Marshal(cfg)
    if err != nil {
        return err
    }
    if err = ioutil.WriteFile(configFilePath, configData, 0600); err != nil {
        return err
    }

    // save host file
    hostnamePath := cfg.rootPath + hostname_file
    if len(cfg.SlaveSection.SlaveNodeName) != 0 {
        err = ioutil.WriteFile(hostnamePath, []byte(cfg.SlaveSection.SlaveNodeName), 0644);
        if err != nil {
            return err
        }
    }
    return nil
}

func (pc *PocketSlaveConfig) SlavePublicKey() ([]byte, error) {
    pubKeyPath := pc.rootPath + slave_public_Key_file
    if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(pubKeyPath)
}

func (pc *PocketSlaveConfig) SlavePrivateKey() ([]byte, error) {
    prvKeyPath := pc.rootPath + slave_prvate_Key_file
    if _, err := os.Stat(prvKeyPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(prvKeyPath)
}

func (pc *PocketSlaveConfig) SlaveSSHKey() ([]byte, error) {
    sshKeyPath := pc.rootPath + slave_ssh_Key_file
    if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(sshKeyPath)
}

func (pc *PocketSlaveConfig) MasterPublicKey() ([]byte, error) {
    masterPubKey := pc.rootPath + master_public_Key_file
    if _, err := os.Stat(masterPubKey); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] Master Publickey might have not been synced yet.")
    }
    return ioutil.ReadFile(masterPubKey)
}

func (pc *PocketSlaveConfig) SaveMasterPublicKey(masterPubKey []byte) error {
    if len(masterPubKey) == 0 {
        return fmt.Errorf("[ERR] Cannot save empty master key")
    }
    keyPath := pc.rootPath + master_public_Key_file
    return ioutil.WriteFile(keyPath, masterPubKey, 0600)
}

func (pc *PocketSlaveConfig) ClearMasterPublicKey() error {
    keyPath := pc.rootPath + master_public_Key_file
    return os.Remove(keyPath)
}

// --- network interface redefinition
func _network_iface_redefined(slaveConfig *ConfigSlaveSection) []string {
    return []string {
        POCKET_START,
        "iface eth0 inet static",
        fmt.Sprintf("%s %s", SLAVE_ADDRESS_KEY, slaveConfig.SlaveIP4Addr),
        fmt.Sprintf("%s %s", SLAVE_GATEWAY_KEY, slaveConfig.SlaveGateway),
        fmt.Sprintf("%s %s", SLAVE_NETMASK_KEY, slaveConfig.SlaveNetMask),
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

func _fixateNetworkInterfaces(slaveConfig *ConfigSlaveSection, ifaceData []string) []string {
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

func (pc *PocketSlaveConfig) SaveFixedNetworkInterface() error {
    var ifaceFilePath string = pc.rootPath + network_iface_file
    ifaceFileContent, err := ioutil.ReadFile(ifaceFilePath)
    var ifaceData []string = strings.Split(string(ifaceFileContent),"\n")
    if err != nil {
        return err
    }

    var fixedIfaceData []string = _fixateNetworkInterfaces(pc.SlaveSection, ifaceData)
    var fixedIfaceContent []byte = []byte(strings.Join(fixedIfaceData, "\n"))
    return ioutil.WriteFile(ifaceFilePath, fixedIfaceContent, 0644)
}