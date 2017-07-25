package config

import (
    "io/ioutil"
    "os"
    "path/filepath"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
    "github.com/stkim1/pcrypto"
)

// ------ CONFIG VERSION -------
const (
    PC_MASTER                   string = "pc-master"
    CORE_STATUS_KEY             string = "binding-status"
    CORE_CONFIG_KEY             string = "config-version"
    CORE_CONFIG_VAL             string = "1.0.1"
)

// ------ CONFIGURATION FILES ------
const (
    // config directory
    dir_core_config             string = "/etc/pocket/"

    // core config file
    core_config_file            string = "core.conf.yaml"
    core_cluster_id_file        string = "cluster.id"
    core_ssh_auth_token_file    string = "ssh.auth.token"
    core_user_name_file         string = "core.user.name"

    // cert directory
    dir_core_certs              string = "pki"

    // these files are 2048 RSA crypto files used to join network
    core_vbox_public_Key_file   string = "pc_core_vbox" + pcrypto.FileExtPublicKey
    core_vbox_prvate_Key_file   string = "pc_core_vbox" + pcrypto.FileExtPrivateKey
    master_vbox_public_Key_file string = "pc_master_vbox" + pcrypto.FileExtPublicKey

    // these files are 2048 RSA crypto files used for Docker & Registry
    core_engine_auth_cert_file  string = "pc_core_engine" + pcrypto.FileExtAuthCertificate
    core_engine_key_cert_file   string = "pc_core_engine" + pcrypto.FileExtCertificate
    core_engine_prvate_key_file string = "pc_core_engine" + pcrypto.FileExtPrivateKey

    // these are files used for teleport certificate
    core_ssh_key_cert_file      string = "pc_core_ssh" + pcrypto.FileExtSSHCertificate
    core_ssh_private_key_file   string = "pc_core_ssh" + pcrypto.FileExtPrivateKey
)

// --- structs ---
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

// --- functions ---
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
        dirConfig      string = DirPathCoreConfig(rootPath)
        dirCerts       string = DirPathCoreCerts(rootPath)

        // pocket cluster join keys
        pathCorePubKey string = FilePathCoreVboxPublicKey(rootPath)
        pathCorePrvKey string = FilePathCoreVboxPrivateKey(rootPath)

        // config file path
        pathConfigFile string = FilePathCoreConfig(rootPath)
        makeKeys       bool   = false

        err            error  = nil
    )

    // check if config dir exists, and creat if DNE
    if _, err := os.Stat(dirConfig); os.IsNotExist(err) {
        os.MkdirAll(dirConfig, 0700);
    }

    // check if config secure key dir also exists and creat if DNE
    if _, err := os.Stat(dirCerts); os.IsNotExist(err) {
        os.MkdirAll(dirCerts, 0700);
    }

    // create pocketcluster join key sets
    if _, err := os.Stat(pathCorePubKey); os.IsNotExist(err) {
        makeKeys = true
    }
    if _, err := os.Stat(pathCorePrvKey); os.IsNotExist(err) {
        makeKeys = true
    }
    if makeKeys {
        err = pcrypto.GenerateStrongKeyPairFiles(pathCorePubKey, pathCorePrvKey, "")
        if err != nil {
            log.Panic(errors.WithStack(err).Error())
        }
    }

    // check if config file exists in path.
    if _, err := os.Stat(pathConfigFile); os.IsNotExist(err) {
        return _brandNewSlaveConfig(rootPath)
    }

    // if does, unmarshal and load them.
    if configData, err := ioutil.ReadFile(pathConfigFile); err != nil {
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
    configDirPath := cfg.rootPath + dir_core_config
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
    pubKeyPath := pc.rootPath + core_vbox_public_Key_file
    if _, err := os.Stat(pubKeyPath); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] public key has not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(pubKeyPath)
}

func (pc *PocketCoreConfig) CorePrivateKey() ([]byte, error) {
    prvKeyPath := pc.rootPath + core_vbox_prvate_Key_file
    if _, err := os.Stat(prvKeyPath); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] private key has not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(prvKeyPath)
}

func (pc *PocketCoreConfig) MasterPublicKey() ([]byte, error) {
    masterPubKey := pc.rootPath + master_vbox_public_Key_file
    if _, err := os.Stat(masterPubKey); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] Master Publickey might have not been synced yet.")
    }
    return ioutil.ReadFile(masterPubKey)
}

// --- to read config ---
func DirPathCoreConfig(rootPath string) string {
    return filepath.Join(rootPath, dir_core_config)
}

func FilePathCoreConfig(rootPath string) string {
    return filepath.Join(DirPathCoreConfig(rootPath), core_config_file)
}

func FilePathClusterID(rootPath string) string {
    return filepath.Join(DirPathCoreConfig(rootPath), core_cluster_id_file)
}

func FilePathAuthToken(rootPath string) string {
    return filepath.Join(DirPathCoreConfig(rootPath), core_ssh_auth_token_file)
}

// --- to read certs ---
func DirPathCoreCerts(rootPath string) string {
    return filepath.Join(DirPathCoreConfig(rootPath), dir_core_certs)
}

func FilePathCoreVboxPublicKey(rootPath string) string {
    return filepath.Join(DirPathCoreCerts(rootPath), core_vbox_public_Key_file)
}

func FilePathCoreVboxPrivateKey(rootPath string) string {
    return filepath.Join(DirPathCoreCerts(rootPath), core_vbox_prvate_Key_file)
}

func FilePathMasterVboxPublicKey(rootPath string) string {
    return filepath.Join(DirPathCoreCerts(rootPath), master_vbox_public_Key_file)
}

func FilePathCoreSSHKeyCert(rootPath string) string {
    return filepath.Join(DirPathCoreCerts(rootPath), core_ssh_key_cert_file)
}

func FilePathCoreSSHPrivateKey(rootPath string) string {
    return filepath.Join(DirPathCoreCerts(rootPath), core_ssh_private_key_file)
}

// --- to build tar archive file ---
func ArchivePathClusterID() string {
    return core_cluster_id_file
}

func ArchivePathAuthToken() string {
    return core_ssh_auth_token_file
}

func ArchivePathUserName() string {
    return core_user_name_file
}

func ArchivePathCertsDir() string {
    return dir_core_certs
}

func ArchivePathCoreVboxPublicKey() string {
    return filepath.Join(ArchivePathCertsDir(), core_vbox_public_Key_file)
}

func ArchivePathCoreVboxPrivateKey() string {
    return filepath.Join(ArchivePathCertsDir(), core_vbox_prvate_Key_file)
}

func ArchivePathMasterVboxPublicKey() string {
    return filepath.Join(ArchivePathCertsDir(), master_vbox_public_Key_file)
}

func ArchivePathCoreEngineAuthCert() string {
    return filepath.Join(ArchivePathCertsDir(), core_engine_auth_cert_file)
}

func ArchivePathCoreEngineKeyCert() string {
    return filepath.Join(ArchivePathCertsDir(), core_engine_key_cert_file)
}

func ArchivePathCoreEnginePrivateKey() string {
    return filepath.Join(ArchivePathCertsDir(), core_engine_prvate_key_file)
}
