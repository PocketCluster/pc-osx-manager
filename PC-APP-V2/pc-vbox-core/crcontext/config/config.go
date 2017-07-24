package config

import (
    "io/ioutil"
    "os"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
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
    CORE_CONFIG_DIR  string            = "/etc/pocket/"
    CORE_CLUSTER_ID_FILE string        = CORE_CONFIG_DIR + "cluster.id"
    CORE_SSH_AUTH_TOKEN_FILE string    = CORE_CONFIG_DIR + "ssh.auth.token"
    CORE_USER_NAME_FILE string         = CORE_CONFIG_DIR + "core.user.name"
    core_config_file string            = CORE_CONFIG_DIR + "core-conf.yaml"

    CORE_CERTS_DIR string              = CORE_CONFIG_DIR + "pki/"
    // these files are 2048 RSA crypto files used to join network
    core_public_Key_file string        = CORE_CERTS_DIR + "pc_core_vbox_report" + pcrypto.FileExtPublicKey
    core_prvate_Key_file string        = CORE_CERTS_DIR + "pc_core_vbox_report" + pcrypto.FileExtPrivateKey
    master_public_Key_file string      = CORE_CERTS_DIR + "pc_master_vbox_ctrl" + pcrypto.FileExtPublicKey

    CORE_TLS_AUTH_CERT_FILE string     = CORE_CERTS_DIR + "pc_core_tls"         + pcrypto.FileExtAuthCertificate
    CORE_TLS_PRVATE_KEY_FILE string    = CORE_CERTS_DIR + "pc_core_tls"         + pcrypto.FileExtPrivateKey
    CERT_TLS_CERTIFICATE_FILE string   = CORE_CERTS_DIR + "pc_core_tls"         + pcrypto.FileExtCertificate

    // these files are 2048 RSA crypto files used for Docker & Registry. This should be acquired from Teleport Auth server
    CoreAuthCertFileName string        = CORE_CERTS_DIR + "pc_cert_auth"        + pcrypto.FileExtCertificate
    CoreEngineKeyFileName string       = CORE_CERTS_DIR + "pc_core_engine"      + pcrypto.FileExtPrivateKey
    CoreEngineCertFileName string      = CORE_CERTS_DIR + "pc_core_engine"      + pcrypto.FileExtCertificate

    // these are files used for teleport certificate
    CoreSSHCertificateFileName string  = CORE_CERTS_DIR + "pc_core_ssh"         + pcrypto.FileExtSSHCertificate
    CoreSSHPrivateKeyFileName string   = CORE_CERTS_DIR + "pc_core_ssh"         + pcrypto.FileExtPrivateKey

    // HOST GENERAL CONFIG
    host_timezone_file string          = "/etc/timezone"
)

// ------ SALT DEFAULT ------
const (
    PC_MASTER string                   = "pc-master"
)

// --- struct
type ConfigMasterSection struct {
    MasterIP4Address    string                   `yaml:"-"`
    MasterIP6Address    string                   `yaml:"-"`
    MasterTimeZone      string                   `yaml:"master-timezone"`
}

type ConfigCoreSection struct {
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
        CoreSection:      &ConfigCoreSection{},
    }
}

func _loadCoreConfig(rootPath string) (*PocketCoreConfig) {

    var (
        // config and key directories
        configDirPath string    = rootPath + CORE_CONFIG_DIR
        keysDirPath string      = rootPath + CORE_CERTS_DIR

        // pocket cluster join keys
        pcPubKeyPath string     = rootPath + core_public_Key_file
        pcPrvKeyPath string     = rootPath + core_prvate_Key_file

        // config file path
        configFilePath string   = rootPath + core_config_file

        makeKeys bool           = false

        err error               = nil
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
        err = pcrypto.GenerateStrongKeyPairFiles(pcPubKeyPath, pcPrvKeyPath, "")
        if err != nil {
            log.Panic(errors.WithStack(err).Error())
        }
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
    configDirPath := cfg.rootPath + CORE_CONFIG_DIR
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
        return nil, errors.Errorf("[ERR] public key has not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(pubKeyPath)
}

func (pc *PocketCoreConfig) CorePrivateKey() ([]byte, error) {
    prvKeyPath := pc.rootPath + core_prvate_Key_file
    if _, err := os.Stat(prvKeyPath); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] private key has not been generated properly. This is a critical error")
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
    return c.rootPath + CORE_CERTS_DIR
}

func (c *PocketCoreConfig) ConfigDir() string {
    return c.rootPath + CORE_CONFIG_DIR
}
