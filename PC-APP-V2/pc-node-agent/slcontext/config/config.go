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
    SLAVE_CONFIG_VAL          string = "1.0.1"
    SLAVE_STATUS_BOUNDED      string = "slave_bounded"
    SLAVE_STATUS_UNBOUNDED    string = "slave_unbounded"
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
    //BindingStatus       string                   `yaml:"binding-status"`
    MasterSection       *ConfigMasterSection     `yaml:"master-section"`
    SlaveSection        *ConfigSlaveSection      `yaml:"slave-section"`
}

// This is default public constructor as it does not accept root file path
func LoadPocketSlaveConfig() *PocketSlaveConfig {
    return loadSlaveConfig("")
}

// --- func
func brandNewSlaveConfig(rootPath string) (*PocketSlaveConfig) {
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

func loadSlaveConfig(rootPath string) (*PocketSlaveConfig) {

    var (
        // config and key directories
        dirSlaveConfig string = DirPathSlaveConfig(rootPath)
        pathConfigFile string = FilePathSlaveConfig(rootPath)
        dirSlaveCert   string = DirPathSlaveCerts(rootPath)
        pathPubKey     string = FilePathSlavePublicKey(rootPath)
        pathPrvKey     string = FilePathSlavePrivateKey(rootPath)

        cfg            *PocketSlaveConfig = &PocketSlaveConfig{}
        makeKeys       bool   = false
        data           []byte = nil
        err            error  = nil
    )

    // check if config dir exists, and creat if DNE
    if _, err := os.Stat(dirSlaveConfig); os.IsNotExist(err) {
        os.MkdirAll(dirSlaveConfig, 0700)
    }

    // check if config secure key dir also exists and creat if DNE
    if _, err := os.Stat(dirSlaveCert); os.IsNotExist(err) {
        os.MkdirAll(dirSlaveCert, 0700)
    }

    // create pocketcluster join key sets
    if _, err := os.Stat(pathPubKey); os.IsNotExist(err) {
        makeKeys = true
    }
    if _, err := os.Stat(pathPrvKey); os.IsNotExist(err) {
        makeKeys = true
    }
    if makeKeys {
        pcrypto.GenerateWeakKeyPairFiles(pathPubKey, pathPrvKey, "")
    }

    // check if config file exists in path.
    if _, err := os.Stat(pathConfigFile); os.IsNotExist(err) {
        return brandNewSlaveConfig(rootPath)
    }

    // if does, unmarshal and load them.
    data, err = ioutil.ReadFile(pathConfigFile)
    if err != nil {
        return brandNewSlaveConfig(rootPath)
    }
    err = yaml.Unmarshal(data, cfg)
    if err != nil {
        return brandNewSlaveConfig(rootPath)
    }
    // as rootpath is ignored, we need to restore it
    cfg.rootPath = rootPath
    return cfg
}

func (c *PocketSlaveConfig) SaveSlaveConfig() error {
    // check if config dir exists, and creat if DNE
    var (
        dirSlaveConfig string = DirPathSlaveConfig(c.rootPath)
        pathConfigFile string = FilePathSlaveConfig(c.rootPath)
        data           []byte = nil
        err            error  = nil
    )
    if _, err := os.Stat(dirSlaveConfig); os.IsNotExist(err) {
        os.MkdirAll(dirSlaveConfig, os.ModeDir|0700);
    }

    data, err = yaml.Marshal(c)
    if err != nil {
        return errors.WithStack(err)
    }
    err = ioutil.WriteFile(pathConfigFile, data, 0600)
    if err != nil {
        return errors.WithStack(err)
    }
    return nil
}

func (c *PocketSlaveConfig) SaveHostname() error {
    // save host file
    pathHostname := FilePathSystemHostname(c.rootPath)
    if len(c.SlaveSection.SlaveNodeName) != 0 {
        return ioutil.WriteFile(pathHostname, []byte(c.SlaveSection.SlaveNodeName), 0644);
    }
    return nil
}

func (c *PocketSlaveConfig) SlavePublicKey() ([]byte, error) {
    pathPubKey := FilePathSlavePublicKey(c.rootPath)
    if _, err := os.Stat(pathPubKey); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(pathPubKey)
}

func (c *PocketSlaveConfig) SlavePrivateKey() ([]byte, error) {
    pathPrvKey := FilePathSlavePrivateKey(c.rootPath)
    if _, err := os.Stat(pathPrvKey); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] keys have not been generated properly. This is a critical error")
    }
    return ioutil.ReadFile(pathPrvKey)
}

func (c *PocketSlaveConfig) MasterPublicKey() ([]byte, error) {
    pathMasterPubKey := FilePathMasterPublicKey(c.rootPath)
    if _, err := os.Stat(pathMasterPubKey); os.IsNotExist(err) {
        return nil, errors.Errorf("[ERR] master publickey might have not been synced yet.")
    }
    return ioutil.ReadFile(pathMasterPubKey)
}

func (c *PocketSlaveConfig) SaveMasterPublicKey(masterPubKey []byte) error {
    if len(masterPubKey) == 0 {
        return errors.Errorf("[ERR] cannot save empty master key")
    }
    return ioutil.WriteFile(FilePathMasterPublicKey(c.rootPath), masterPubKey, 0600)
}

func (c *PocketSlaveConfig) ClearMasterPublicKey() error {
    return os.Remove(FilePathMasterPublicKey(c.rootPath))
}

func (c *PocketSlaveConfig) ConfigDir() string {
    return DirPathSlaveConfig(c.rootPath)
}

func (c *PocketSlaveConfig) KeyAndCertDir() string {
    return DirPathSlaveCerts(c.rootPath)
}
