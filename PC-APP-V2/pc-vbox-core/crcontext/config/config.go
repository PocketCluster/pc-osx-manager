package config

import (
    "io/ioutil"
    "os"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "gopkg.in/yaml.v2"
)

// ------ CONFIG VERSION -------
const (
    PC_MASTER                   string = "pc-master"
    CORE_STATUS_KEY             string = "binding-status"
    CORE_CONFIG_KEY             string = "config-version"
    CORE_CONFIG_VAL             string = "1.0.4"
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

// --- functions ---
func _brandNewSlaveConfig(rootPath string) (*PocketCoreConfig) {
    return &PocketCoreConfig {
        rootPath:         rootPath,
        ConfigVersion:    CORE_CONFIG_VAL,
        MasterSection:    &MasterConfigSection{},
        CoreSection:      &CoreConfigSection{},
    }
}

func _loadCoreConfig(rootPath string) (*PocketCoreConfig) {
    var (
        dirConfig      string = DirPathCoreConfig(rootPath)
        dirCerts       string = DirPathCoreCerts(rootPath)
        pathCoreConfig string = FilePathCoreConfig(rootPath)

        config         *PocketCoreConfig = &PocketCoreConfig{}
        cfgData        []byte            = nil
        err            error             = nil
    )

    // Check if config dir exists, and crash if doesn't
    // Core /etc/pocket should have existed before run pocketd
    _, err = os.Stat(dirConfig)
    if err != nil {
        log.Panic("[CRITICAL] config directory should have existed")
    }
    // Check if config secure key dir also exists
    // Core /etc/pocket should have existed before run pocketd
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

// This is default public constructor as it does not accept root file path
func LoadPocketCoreConfig() *PocketCoreConfig {
    return _loadCoreConfig("")
}

func (c *PocketCoreConfig) SaveCoreConfig() error {
    // check if config dir exists, and creat if DNE
    dir := DirPathCoreConfig(c.rootPath)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        os.MkdirAll(dir, os.ModeDir|0700);
    }

    path := FilePathCoreConfig(c.rootPath)
    config, err := yaml.Marshal(c)
    if err != nil {
        return err
    }
    if err = ioutil.WriteFile(path, config, 0600); err != nil {
        return err
    }
    return nil
}

func (c *PocketCoreConfig) CorePublicKey() ([]byte, error) {
    path := FilePathCoreVboxPublicKey(c.rootPath)
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] core public key has not been generated properly")
    }
    return ioutil.ReadFile(path)
}

func (c *PocketCoreConfig) CorePrivateKey() ([]byte, error) {
    path := FilePathCoreVboxPrivateKey(c.rootPath)
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] core private key has not been generated properly")
    }
    return ioutil.ReadFile(path)
}

func (c *PocketCoreConfig) MasterPublicKey() ([]byte, error) {
    path := FilePathMasterVboxPublicKey(c.rootPath)
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] master public key has not been generated properly")
    }
    return ioutil.ReadFile(path)
}
