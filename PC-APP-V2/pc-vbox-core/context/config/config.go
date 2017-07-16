package config

import (
    "io/ioutil"
    "os"

    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
    "github.com/pborman/uuid"
    "github.com/stkim1/pcrypto"
)

// ------ CONFIG VERSION -------
const (
    CORE_CONFIG_KEY string            = "config-version"
    CORE_CONFIG_VAL string            = "1.0.1"
)

const (
    CORE_STATUS_KEY string            = "binding-status"
)

// ------ CONFIGURATION FILES ------
const (
    // POCKET SPECIFIC CONFIG
    core_config_dir string             = "/etc/pocket/"
    core_config_file string            = core_config_dir + "core-conf.yaml"

    core_keys_dir string               = core_config_dir + "pki/"
    // these files are 2048 RSA crypto files used to join network
    core_public_Key_file string        = core_keys_dir + "pc_core_vbox_report" + pcrypto.FileExtPublicKey
    core_prvate_Key_file string        = core_keys_dir + "pc_core_vbox_report" + pcrypto.FileExtPrivateKey
    master_public_Key_file string      = core_keys_dir + "pc_master_vbox_ctrl" + pcrypto.FileExtPublicKey

    // these files are 2048 RSA crypto files used for Docker & Registry. This should be acquired from Teleport Auth server
    CoreAuthCertFileName string        = core_keys_dir + "pc_cert_auth"        + pcrypto.FileExtCertificate
    CoreEngineKeyFileName string       = core_keys_dir + "pc_core_engine"      + pcrypto.FileExtPrivateKey
    CoreEngineCertFileName string      = core_keys_dir + "pc_core_engine"      + pcrypto.FileExtCertificate

    // these are files used for teleport certificate
    CoreSSHCertificateFileName string  = core_keys_dir + "pc_core_ssh"         + pcrypto.FileExtSSHCertificate
    CoreSSHPrivateKeyFileName string   = core_keys_dir + "pc_core_ssh"         + pcrypto.FileExtPrivateKey

    // HOST GENERAL CONFIG
    host_timezone_file string          = "/etc/timezone"
)

// ------ SALT DEFAULT ------
const (
    PC_MASTER string                   = "pc-master"
)

// --- struct
type ConfigMasterSection struct {
    MasterTimeZone      string                   `yaml:"master-timezone"`
}

type ConfigCoreSection struct {
    CoreNodeName        string                   `yaml:"core-node-name"`
    CoreNodeUUID        string                   `yaml:"core-node-uuid"`
    CoreAuthToken       string                   `yaml:"core-auth-token"`
    CoreMacAddr         string                   `yaml:"core-mac-addr"`
}

type PocketCoreConfig struct {
    // this field exists to create files at a specific location for testing so ignore
    rootPath            string                   `yaml:"-"`
    ConfigVersion       string                   `yaml:"config-version"`

    // Cluster Identity
    ClusterID           string                   `yaml:"cluster-id"`
    MasterSection       *ConfigMasterSection     `yaml:"master-section",inline,flow`
    CoreSection         *ConfigCoreSection       `yaml:"core-section",inline,flow`
}

// This is default public constructor as it does not accept root file path
func LoadPocketCoreConfig() *PocketCoreConfig {
    return _loadCoreConfig("")
}

// --- func
func _brandNewSlaveConfig(rootPath string) (*PocketCoreConfig) {
    return &PocketCoreConfig{
        rootPath:         rootPath,
        ConfigVersion:    CORE_CONFIG_VAL,
        MasterSection:    &ConfigMasterSection{},
        CoreSection:      &ConfigCoreSection{
            CoreNodeUUID:    uuid.New(),
        },
    }
}

func _loadCoreConfig(rootPath string) (*PocketCoreConfig) {

    var (
        // config and key directories
        configDirPath string    = rootPath + core_config_dir
        keysDirPath string      = rootPath + core_keys_dir

        // pocket cluster join keys
        pcPubKeyPath string     = rootPath + core_public_Key_file
        pcPrvKeyPath string     = rootPath + core_prvate_Key_file

        // config file path
        configFilePath string   = rootPath + core_config_file

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
        pcrypto.GenerateStrongKeyPairFiles(pcPubKeyPath, pcPrvKeyPath, "")
    }

    // check if config file exists in path.
    if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
        return _brandNewSlaveConfig(rootPath)
    }

    // if does, unmarshal and load them.
    if configData, err := ioutil.ReadFile(configFilePath); err != nil {
        return _brandNewSlaveConfig(rootPath)
    } else {
        var config PocketCoreConfig
        if err = yaml.Unmarshal(configData, &config); err != nil {
            return _brandNewSlaveConfig(rootPath)
        } else {
            // as rootpath is ignored, we need to restore it
            config.rootPath = rootPath
            return &config
        }
    }
}

func (cfg *PocketCoreConfig) SaveCoreConfig() error {
    // check if config dir exists, and creat if DNE
    configDirPath := cfg.rootPath + core_config_dir
    if _, err := os.Stat(configDirPath); os.IsNotExist(err) {
        os.MkdirAll(configDirPath, os.ModeDir|0700);
    }

    configFilePath := cfg.rootPath + core_config_file
    configData, err := yaml.Marshal(cfg)
    if err != nil {
        return err
    }
    if err = ioutil.WriteFile(configFilePath, configData, 0600); err != nil {
        return err
    }
    return nil
}

func (pc *PocketCoreConfig) CorePublicKey() ([]byte, error) {
    pubKeyPath := pc.rootPath + core_public_Key_file
    if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(pubKeyPath)
}

func (pc *PocketCoreConfig) CorePrivateKey() ([]byte, error) {
    prvKeyPath := pc.rootPath + core_prvate_Key_file
    if _, err := os.Stat(prvKeyPath); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(prvKeyPath)
}

func (pc *PocketCoreConfig) MasterPublicKey() ([]byte, error) {
    masterPubKey := pc.rootPath + master_public_Key_file
    if _, err := os.Stat(masterPubKey); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] Master Publickey might have not been synced yet.")
    }
    return ioutil.ReadFile(masterPubKey)
}

func (pc *PocketCoreConfig) SaveMasterPublicKey(masterPubKey []byte) error {
    if len(masterPubKey) == 0 {
        return errors.Errorf("[ERR] Cannot save empty master key")
    }
    keyPath := pc.rootPath + master_public_Key_file
    return ioutil.WriteFile(keyPath, masterPubKey, 0600)
}

func (pc *PocketCoreConfig) ClearMasterPublicKey() error {
    keyPath := pc.rootPath + master_public_Key_file
    return os.Remove(keyPath)
}

func (c *PocketCoreConfig) KeyAndCertDir() string {
    return c.rootPath + core_keys_dir
}

func (c *PocketCoreConfig) ConfigDir() string {
    return c.rootPath + core_config_dir
}
