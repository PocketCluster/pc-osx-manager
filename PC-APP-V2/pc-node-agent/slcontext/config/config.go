package config

import (
    "os"
    "fmt"
    "io/ioutil"
    "strings"

    "gopkg.in/yaml.v2"
    "github.com/stkim1/pcrypto"
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

// Master name server is fixed for now (v.0.1.4)
const SLAVE_NAMESRV_VALUE = "pc-master:53535"

// ------ CONFIGURATION FILES ------
const (
    // POCKET SPECIFIC CONFIG
    slave_config_dir        = "/etc/pocket/"
    slave_config_file       = "/etc/pocket/slave-conf.yaml"

    slave_keys_dir          = "/etc/pocket/pki/"
    // these files are 1024 RSA crypto files used to join network
    slave_public_Key_file   = "/etc/pocket/pki/pcslave.pub"
    slave_prvate_Key_file   = "/etc/pocket/pki/pcslave.pem"
    master_public_Key_file  = "/etc/pocket/pki/pcmaster.pub"

    // these files are 2048 RSA crypto files used for SSH, DOCKER
    node_public_Key_file    = "/etc/pocket/pki/node.pub"
    node_private_Key_file   = "/etc/pocket/pki/node.pem"
    node_certificate_file   = "/etc/pocket/pki/node.csr"

    // HOST GENERAL CONFIG
    network_iface_file      = "/etc/network/interfaces"
    hostname_file           = "/etc/hostname"
    //hostaddr_file           = "/etc/hosts"
    host_timezone_file      = "/etc/timezone"
    //resolve_conf_file       = "/etc/resolv.conf"

    SlaveDockerAuthFileName = "docker.auth"
    SlaveDockerKeyFileName  = "docker.key"
    SlaveDockerCertFileName = "docker.cert"

    slave_docker_auth_file  = slave_keys_dir + SlaveDockerAuthFileName
    slave_docker_key_file   = slave_keys_dir + SlaveDockerKeyFileName
    slave_docker_cert_file  = slave_keys_dir + SlaveDockerCertFileName

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

    var (
        // config and key directories
        configDirPath string    = rootPath + slave_config_dir
        keysDirPath string      = rootPath + slave_keys_dir

        // pocket cluster join keys
        pcPubKeyPath string     = rootPath + slave_public_Key_file
        pcPrvKeyPath string     = rootPath + slave_prvate_Key_file

        // node SSH/DOCKER keys
        nodePubKeyPath string   = rootPath + node_public_Key_file
        nodePrvKeyPath string   = rootPath + node_private_Key_file

        // config file path
        configFilePath string   = rootPath + slave_config_file
    )

    // check if config dir exists, and creat if DNE
    if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
        os.MkdirAll(configDirPath, 0700);
    }

    // check if config secure key dir also exists and creat if DNE
    if _, err := os.Stat(keysDirPath); os.IsNotExist(err) {
        os.MkdirAll(keysDirPath, 0700);
    }

    // create pocketcluster join key sets
    var makePcJoinKeys bool = false
    if _, err := os.Stat(pcPubKeyPath); os.IsNotExist(err) {
        makePcJoinKeys = true
    }
    if _, err := os.Stat(pcPrvKeyPath); os.IsNotExist(err) {
        makePcJoinKeys = true
    }
    if makePcJoinKeys {
        pcrypto.GenerateWeakKeyPairFiles(pcPubKeyPath, pcPrvKeyPath, "")
    }

    // create node ssh key sets
    var makeNodeKeys bool = false
    if _, err := os.Stat(nodePubKeyPath); os.IsNotExist(err) {
        makeNodeKeys = true
    }
    if _, err := os.Stat(nodePrvKeyPath); os.IsNotExist(err) {
        makeNodeKeys = true
    }
    if makeNodeKeys {
        pcrypto.GenerateStrongKeyPairFiles(nodePubKeyPath, nodePrvKeyPath, "")
    }

    // check if config file exists in path.
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
        os.MkdirAll(configDirPath, os.ModeDir|0700);
    }

    configFilePath := cfg.rootPath + slave_config_file
    configData, err := yaml.Marshal(cfg)
    if err != nil {
        return err
    }
    if err = ioutil.WriteFile(configFilePath, configData, 0600); err != nil {
        return err
    }
    return nil
}

func (cfg *PocketSlaveConfig) SaveHostname() error {
    // save host file
    hostnamePath := cfg.rootPath + hostname_file
    if len(cfg.SlaveSection.SlaveNodeName) != 0 {
        return ioutil.WriteFile(hostnamePath, []byte(cfg.SlaveSection.SlaveNodeName), 0644);
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

func (pc *PocketSlaveConfig) NodePublicKey() ([]byte, error) {
    pubKeyPath := pc.rootPath + node_public_Key_file
    if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(pubKeyPath)
}

func (pc *PocketSlaveConfig) NodePrivateKey() ([]byte, error) {
    prvKeyPath := pc.rootPath + node_private_Key_file
    if _, err := os.Stat(prvKeyPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(prvKeyPath)
}

func (pc *PocketSlaveConfig) NodeCertificate() ([]byte, error) {
    certPath := pc.rootPath + node_certificate_file
    if _, err := os.Stat(certPath); os.IsNotExist(err) {
        return nil, fmt.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(certPath)
}

func (pc *PocketSlaveConfig) SaveNodeCertificate(certificate []byte) error {
    if len(certificate) == 0 {
        return fmt.Errorf("[ERR] Cannot save empty node certificate")
    }
    filepath := pc.rootPath + node_certificate_file
    return ioutil.WriteFile(filepath, certificate, 0600)
}

func (pc *PocketSlaveConfig) ClearNodeCertificate() error {
    filepath := pc.rootPath + node_certificate_file
    return os.Remove(filepath)
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

func (c *PocketSlaveConfig) KeyAndCertDir() string {
    return c.rootPath + slave_keys_dir
}

func (c *PocketSlaveConfig) ConfigDir() string {
    return c.rootPath + slave_config_dir
}
