package config

import (
    "os"
    "io/ioutil"

    "gopkg.in/yaml.v2"
    "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pc-node-agent/slcontext"
    "fmt"
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
    SLAVE_ADDRESS_KEY   = "address"
    SLAVE_NETMASK_KEY   = "netmask"
    SLAVE_BROADCS_KEY   = "broadcast"
    SLAVE_GATEWAY_KEY   = "gateway"

    // TODO we might not need this
    //NAMESRV             = "dns-nameservers"
)
// TODO : we might not need this
//var IFACE_KEYS []string = []string{ADDRESS, NETMASK, BROADCS, GATEWAY}

const (
    PAGENT_SEND_PORT    = 10060
    PAGENT_RECV_PORT    = 10061
)

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
    hostaddr_file           = "/etc/hosts"
    host_timezone_file      = "/etc/timezone"
    resolve_conf_file       = "/etc/resolv.conf"
)

// ------ SALT DEFAULT ------
const (
    PC_MASTER           = "pc-master"
)

// ------- POCKET EDITOR MARKER ------
const (
    POCKET_START        = "// --------------- POCKETCLUSTER START ---------------"
    POCKET_END          = "// ---------------  POCKETCLUSTER END  ---------------"
)

type PocketSlaveConfig interface {
    // load context with config values
    LoadupContext(ctx *slcontext.PocketSlaveContext) error

    // save context values to config
    SaveContext(ctx *slcontext.PocketSlaveContext) error

    SlavePublicKey() ([]byte, error)
    SlavePrivateKey() ([]byte, error)
    SlaveSSHKey() ([]byte, error)
    MasterPublicKey() ([]byte, error)
}

// This is default public constructor as it does not accept root file path
func NewPocketSlaveConfig() (PocketSlaveConfig) {
    return loadSlaveConfig("")
}

// --- struct
type configMasterSection struct {
    // Master Agent Specific String
    MasterBoundAgent    string                   `yaml:"master-binder-agent"`
    // Last Known IP4
    MasterIP4Address    string                   `yaml:"master-ip4-addr"`
    //MasterIP6Address    string
    //MasterHostName      string
    MasterTimeZone      string                   `yaml:"master-timezone"`
}

type configSlaveSection struct {
    SlaveMacAddr        string                   `yaml:"slave-mac-addr"`
    SlaveNodeName       string                   `yaml:"slave-node-name"`
    SlaveIP4Addr        string                   `yaml:"slave-ip4-addr"`
    //SlaveIP6Addr        string
    SlaveNetMask        string                   `yaml:"slave-net-mask"`
    //SlaveBroadcast      string
    SlaveGateway        string                   `yaml:"slave-gateway"`
    //SlaveNameServ       string
}

type pocketSlaveConfig struct {
    // this field exists to create files at a specific location for testing so ignore
    rootPath            string                   `yaml:"-"`
    ConfigVersion       string                   `yaml:"config-version"`
    BindingStatus       string                   `yaml:"binding-status"`
    MasterSection       *configMasterSection     `yaml:"master-section",inline,flow`
    SlaveSection        *configSlaveSection      `yaml:"slave-section",inline,flow`
}


// --- func
func _brandNewSlaveConfig(rootPath string) (*pocketSlaveConfig) {
    return &pocketSlaveConfig{
        rootPath        :rootPath,
        ConfigVersion   :SLAVE_CONFIG_VAL,
        BindingStatus   :locator.SlaveUnbounded.String(),
        MasterSection   :&configMasterSection{},
        SlaveSection    :&configSlaveSection{},
    }
}

func loadSlaveConfig(rootPath string) (*pocketSlaveConfig) {

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
        var config pocketSlaveConfig
        if err = yaml.Unmarshal(configData, &config); err != nil {
            return _brandNewSlaveConfig(rootPath)
        } else {
            // as rootpath is ignored, we need to restore it
            config.rootPath = rootPath
            return &config
        }
    }
}

func saveSlaveConfig(cfg *pocketSlaveConfig) error {
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
    return ioutil.WriteFile(configFilePath, configData, 0600)
}

func (pc *pocketSlaveConfig) LoadupContext(ctx *slcontext.PocketSlaveContext) error {

    return nil
}

func (pc *pocketSlaveConfig) SaveContext(ctx *slcontext.PocketSlaveContext) error {
    return nil
}

func (pc *pocketSlaveConfig) SlavePublicKey() ([]byte, error) {
    pubKeyPath := pc.rootPath + slave_public_Key_file
    if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(pubKeyPath)
}

func (pc *pocketSlaveConfig) SlavePrivateKey() ([]byte, error) {
    prvKeyPath := pc.rootPath + slave_prvate_Key_file
    if _, err := os.Stat(prvKeyPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(prvKeyPath)
}

func (pc *pocketSlaveConfig) SlaveSSHKey() ([]byte, error) {
    sshKeyPath := pc.rootPath + slave_ssh_Key_file
    if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(sshKeyPath)
}

func (pc *pocketSlaveConfig) MasterPublicKey() ([]byte, error) {
    masterPubKey := pc.rootPath + master_public_Key_file
    if _, err := os.Stat(masterPubKey); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] Master Publickey might have not been synced yet.")
    }
    return ioutil.ReadFile(masterPubKey)
}
