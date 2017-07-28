package config

import (
    "io/ioutil"
    "os"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/pborman/uuid"
    "gopkg.in/yaml.v2"
)

// ------ CONFIG VERSION -------
const (
    PC_MASTER           string = "pc-master"
    CORE_STATUS_KEY     string = "binding-status"
    CORE_CONFIG_KEY     string = "config-version"
    CORE_CONFIG_VER     string = "1.0.1"
)

// --- structs ---
type MasterConfigSection struct {
    MasterIP4Address    string                   `yaml:"-"`
    MasterIP6Address    string                   `yaml:"-"`
    MasterTimeZone      string                   `yaml:"master-timezone"`
}

type CoreConfigSection struct {
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
    MasterSection       *MasterConfigSection     `yaml:"master-section"`
    CoreSection         *CoreConfigSection       `yaml:"core-section"`
}

// --- functions ---
func brandNewSlaveConfig(rootPath string) (*PocketCoreConfig) {
    var (
        pathCoreConfig   string = FilePathCoreConfig(rootPath)
        pathClusterID    string = FilePathClusterID(rootPath)
        pathAuthToken    string = FilePathAuthToken(rootPath)
        cfg              *PocketCoreConfig = nil
        cfgData          []byte = nil
        clusterID        []byte = nil
        authToken        []byte = nil
        err              error  = nil
    )

    // read & delete cluster id
    clusterID, err = ioutil.ReadFile(pathClusterID)
    if err != nil {
        log.Panic(errors.WithStack(err).Error())
    } else {
        os.Remove(pathClusterID)
    }

    // read & delete auth token
    authToken, err = ioutil.ReadFile(pathAuthToken)
    if err != nil {
        log.Panic(errors.WithStack(err).Error())
    } else {
        os.Remove(pathAuthToken)
    }

    // core config
    cfg = &PocketCoreConfig {
        rootPath:         rootPath,
        ConfigVersion:    CORE_CONFIG_VER,
        ClusterID:        string(clusterID),
        MasterSection:    &MasterConfigSection{},
        CoreSection:      &CoreConfigSection{
            CoreNodeUUID:     uuid.New(),
            CoreAuthToken:    string(authToken),
        },
    }
    cfgData, err = yaml.Marshal(cfg)
    if err != nil {
        log.Panic(errors.WithStack(err).Error())
    }
    err = ioutil.WriteFile(pathCoreConfig, cfgData, 0600)
    if err != nil {
        log.Panic(errors.WithStack(err).Error())
    }
    return cfg
}

func loadCoreConfig(rootPath string) (*PocketCoreConfig) {
    var (
        dirConfig      string = DirPathCoreConfig(rootPath)
        dirCerts       string = DirPathCoreCerts(rootPath)
        pathCoreConfig string = FilePathCoreConfig(rootPath)

        cfg            *PocketCoreConfig = &PocketCoreConfig{}
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
        return brandNewSlaveConfig(rootPath)
    }
    err = yaml.Unmarshal(cfgData, cfg)
    if err != nil {
        log.Panic("[CRITICAL] config file cannot be properly read %v", err.Error())
    }
    // as rootpath is ignored, we need to restore it
    cfg.rootPath = rootPath
    return cfg
}

// This is default public constructor as it does not accept root file path
func LoadPocketCoreConfig() *PocketCoreConfig {
    return loadCoreConfig("")
}

func (c *PocketCoreConfig) RootPath() string {
    return c.rootPath
}

func (c *PocketCoreConfig) SaveCoreConfig() error {
    cfg, err := yaml.Marshal(c)
    if err != nil {
        return errors.WithStack(err)
    }
    err = ioutil.WriteFile(FilePathCoreConfig(c.rootPath), cfg, 0600)
    if err != nil {
        return errors.WithStack(err)
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
