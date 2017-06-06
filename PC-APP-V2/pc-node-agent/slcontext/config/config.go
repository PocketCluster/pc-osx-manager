package config

import (
    "fmt"
    "io/ioutil"
    "os"
    "strings"

    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
    "github.com/pborman/uuid"
    "github.com/stkim1/pcrypto"
)

// ------ CONFIG VERSION -------
const (
    SLAVE_CONFIG_KEY string            = "config-version"
    SLAVE_CONFIG_VAL string            = "1.0.1"
)

const (
    SLAVE_STATUS_KEY string            = "binding-status"
)

// ------ NETWORK INTERFACES ------
const (
    SLAVE_ADDRESS_KEY string           = "address"
    SLAVE_GATEWAY_KEY string           = "gateway"
    SLAVE_NETMASK_KEY string           = "netmask"
    SLAVE_NAMESRV_KEY string           = "dns-nameserver"
    SLAVE_BROADCS_KEY string           = "broadcast"
)

var SLAVE_NETIFACE_KEYS []string       = []string{SLAVE_ADDRESS_KEY, SLAVE_GATEWAY_KEY, SLAVE_NETMASK_KEY, SLAVE_NAMESRV_KEY, SLAVE_BROADCS_KEY}

// Master name server is fixed for now (v.0.1.4)
const SLAVE_NAMESRV_VALUE string       = "pc-master:53535"

// ------ CONFIGURATION FILES ------
const (
    // POCKET SPECIFIC CONFIG
    slave_config_dir string            = "/etc/pocket/"
    slave_config_file string           = slave_config_dir + "slave-conf.yaml"

    slave_keys_dir string              = slave_config_dir + "pki/"
    // these files are 1024 RSA crypto files used to join network
    slave_public_Key_file string       = slave_keys_dir + "pc_node_beacon"   + pcrypto.FileExtPublicKey
    slave_prvate_Key_file string       = slave_keys_dir + "pc_node_beacon"   + pcrypto.FileExtPrivateKey
    master_public_Key_file string      = slave_keys_dir + "pc_master_beacon" + pcrypto.FileExtPublicKey

    // these files are 2048 RSA crypto files used for Docker & Registry. This should be acquired from Teleport Auth server
    SlaveAuthCertFileName string       = slave_keys_dir + "pc_cert_auth"   + pcrypto.FileExtCertificate
    SlaveEngineKeyFileName string      = slave_keys_dir + "pc_node_engine" + pcrypto.FileExtPrivateKey
    SlaveEngineCertFileName string     = slave_keys_dir + "pc_node_engine" + pcrypto.FileExtCertificate

    // these are files used for teleport certificate
    SlaveSSHCertificateFileName string = slave_keys_dir + "pc_node_ssh" + pcrypto.FileExtSSHCertificate
    SlaveSSHPrivateKeyFileName string  = slave_keys_dir + "pc_node_ssh" + pcrypto.FileExtPrivateKey

    // these files are 2048 RSA crypto files used for SSH.
    // 1) This should be acquired from Teleport Auth server
    // 2) This should be handled by teleport process
    //node_private_Key_file string       = "/etc/pocket/pki/node.key"
    //node_certificate_file string       = "/etc/pocket/pki/node.cert"

    // HOST GENERAL CONFIG
    network_iface_file string          = "/etc/network/interfaces"
    hostname_file string               = "/etc/hostname"
    //hostaddr_file string               = "/etc/hosts"
    host_timezone_file string          = "/etc/timezone"
    //resolve_conf_file string           = "/etc/resolv.conf"
)

// ------ SALT DEFAULT ------
const (
    PC_MASTER string                  = "pc-master"
)

// ------- POCKET EDITOR MARKER ------
const (
    POCKET_START string               = "# --------------- POCKETCLUSTER START ---------------"
    POCKET_END string                 = "# ---------------  POCKETCLUSTER END  ---------------"
)

// --- struct
type ConfigMasterSection struct {
    // Master Agent Specific String
    MasterBoundAgent    string                   `yaml:"master-binder-agent"`
    // Last Known IP4
    MasterIP4Address    string                   `yaml:"-"`
    //MasterIP6Address    string                   `yaml:"-"`
    //MasterHostName      string                   `yaml:"-"`
    MasterTimeZone      string                   `yaml:"master-timezone"`
}

type ConfigSlaveSection struct {
    SlaveNodeName       string                   `yaml:"slave-node-name"`
    SlaveNodeUUID       string                   `yaml:"slave-node-uuid"`
    SlaveAuthToken      string                   `yaml:"slave-auth-token"`
    SlaveMacAddr        string                   `yaml:"slave-mac-addr"`
    //SlaveIP6Addr        string
    SlaveIP4Addr        string                   `yaml:"slave-net-ip4"`
    SlaveGateway        string                   `yaml:"slave-net-gateway"`
    SlaveNameServ       string                   `yaml:"slave-net-nameserv"`
    //SlaveBroadcast      string                   `yaml:"slave-net-broadc"`
}

type PocketSlaveConfig struct {
    // this field exists to create files at a specific location for testing so ignore
    rootPath            string                   `yaml:"-"`
    ConfigVersion       string                   `yaml:"config-version"`
    // TODO : we need to avoid cyclic import but need to fix this
    //BindingStatus       string                   `yaml:"binding-status"`
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
        rootPath:         rootPath,
        ConfigVersion:    SLAVE_CONFIG_VAL,
        // TODO : we need to avoid cyclic import but need to fix this
        //BindingStatus:    "SlaveUnbounded", //locator.SlaveUnbounded.String(),
        MasterSection:    &ConfigMasterSection{},
        SlaveSection:     &ConfigSlaveSection{
            SlaveNodeUUID:    uuid.New(),
        },
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

        // config file path
        configFilePath string   = rootPath + slave_config_file

        makeKeys bool           = false
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
    if _, err := os.Stat(pcPubKeyPath); os.IsNotExist(err) {
        makeKeys = true
    }
    if _, err := os.Stat(pcPrvKeyPath); os.IsNotExist(err) {
        makeKeys = true
    }
    if makeKeys {
        pcrypto.GenerateWeakKeyPairFiles(pcPubKeyPath, pcPrvKeyPath, "")
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
        return nil, errors.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(pubKeyPath)
}

func (pc *PocketSlaveConfig) SlavePrivateKey() ([]byte, error) {
    prvKeyPath := pc.rootPath + slave_prvate_Key_file
    if _, err := os.Stat(prvKeyPath); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(prvKeyPath)
}

func (pc *PocketSlaveConfig) MasterPublicKey() ([]byte, error) {
    masterPubKey := pc.rootPath + master_public_Key_file
    if _, err := os.Stat(masterPubKey); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] Master Publickey might have not been synced yet.")
    }
    return ioutil.ReadFile(masterPubKey)
}

func (pc *PocketSlaveConfig) SaveMasterPublicKey(masterPubKey []byte) error {
    if len(masterPubKey) == 0 {
        return errors.Errorf("[ERR] Cannot save empty master key")
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
