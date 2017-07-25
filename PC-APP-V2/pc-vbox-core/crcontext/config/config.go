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
    CORE_CONFIG_VAL             string = "1.0.4"
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
type MasterConfigSection struct {
    MasterIP4Address    string                   `yaml:"-"`
    MasterIP6Address    string                   `yaml:"-"`
    MasterTimeZone      string                   `yaml:"master-timezone"`
}

type CoreConfigSection struct {
    CoreAuthToken       string                   `yaml:"core-auth-token"`
    CoreMacAddr         string                   `yaml:"core-mac-addr"`
}

type PocketCoreConfig struct {
    // this field exists to create files at a specific location for testing so ignore
    rootPath            string                   `yaml:"-"`
    ConfigVersion       string                   `yaml:"config-version"`

    // Cluster Identity
    ClusterID           string                   `yaml:"cluster-id"`
    MasterSection       *MasterConfigSection     `yaml:"master-section",inline,flow`
    CoreSection         *CoreConfigSection       `yaml:"core-section",inline,flow`
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
        MasterSection:    &MasterConfigSection{},
        CoreSection:      &CoreConfigSection{},
    }
}

func _loadCoreConfig(rootPath string) (*PocketCoreConfig) {
    var (
        // directories
        dirConfig      string = DirPathCoreConfig(rootPath)
        dirCerts       string = DirPathCoreCerts(rootPath)
        // file path
        pathCoreConfig string = FilePathCoreConfig(rootPath)

        cfgData        []byte            = nil
        config         *PocketCoreConfig = nil
        err            error             = nil
    )

    // check if config dir exists, and creat if DNE
    _, err = os.Stat(dirConfig)
    if err != nil {
        log.Panic("[CRITICAL] config directory should have existed")
    }
    // check if config secure key dir also exists and creat if DNE
    _, err = os.Stat(dirCerts)
    if err != nil {
        log.Panic("[CRITICAL] certs directory should have existed")
    }

    // if does, unmarshal and load them.
    cfgData, err = ioutil.ReadFile(pathCoreConfig)
    if err != nil {
        return _brandNewSlaveConfig(rootPath)
    } else {
        err = yaml.Unmarshal(cfgData, config)
        if err != nil {
            return _brandNewSlaveConfig(rootPath)
        } else {
            // as rootpath is ignored, we need to restore it
            config.rootPath = rootPath
            return config
        }
    }
}

func (c *PocketCoreConfig) SaveCoreConfig() error {
    // check if config dir exists, and creat if DNE
    dirConfig := DirPathCoreConfig(c.rootPath)
    if _, err := os.Stat(dirConfig); os.IsNotExist(err) {
        os.MkdirAll(dirConfig, os.ModeDir|0700);
    }

    pathCoreConfig := FilePathCoreConfig(c.rootPath)
    configData, err := yaml.Marshal(c)
    if err != nil {
        return err
    }
    if err = ioutil.WriteFile(pathCoreConfig, configData, 0600); err != nil {
        return err
    }
    return nil
}

func (c *PocketCoreConfig) CorePublicKey() ([]byte, error) {
    path := FilePathCoreVboxPublicKey(c.rootPath)
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] public key has not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(path)
}

func (c *PocketCoreConfig) CorePrivateKey() ([]byte, error) {
    path := FilePathCoreVboxPrivateKey(c.rootPath)
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] private key has not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(path)
}

func (c *PocketCoreConfig) MasterPublicKey() ([]byte, error) {
    path := FilePathMasterVboxPublicKey(c.rootPath)
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] Master Publickey might have not been synced yet.")
    }
    return ioutil.ReadFile(path)
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
