package crcontext

import (
    "fmt"

    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-vbox-core/crcontext/config"
)

const (
    coreNodeName = "pc-core"
)

// for teleport to build configuration
type PocketCoreProperty interface {
    CoreConfigPath() string
    CoreKeyAndCertPath() string

    CoreSSHCertificateFileName() string
    CoreSSHPrivateKeyFileName() string

    CoreNodeUUID() string
    CoreAuthToken() string
    CoreNodeName() string
    CoreNodeNameFQDN() string
    CoreClusterID() string

    CoreAuthServerAddr() (string, error)
    SetMasterIP4ExtAddr(ip4Address string) error
    GetMasterIP4ExtAddr() (string, error)

}

// TODO : add tests
func (c *coreContext) CoreConfigPath() string {
    return config.DirPathCoreConfig(c.RootPath())
}

// TODO : add tests
func (c *coreContext) CoreKeyAndCertPath() string {
    return config.DirPathCoreCerts(c.RootPath())
}

// TODO : add tests
func (c *coreContext) CoreSSHCertificateFileName() string {
    return config.FilePathCoreSSHKeyCert(c.RootPath())
}

// TODO : add tests
func (c *coreContext) CoreSSHPrivateKeyFileName() string {
    return config.FilePathCoreSSHPrivateKey(c.RootPath())
}

func (c *coreContext) CoreNodeUUID() string {
    return c.CoreSection.CoreNodeUUID
}

func (c *coreContext) CoreAuthToken() string {
    return c.CoreSection.CoreAuthToken
}

func (c *coreContext) CoreNodeName() string {
    return coreNodeName
}

// TODO : add tests
func (c *coreContext) CoreNodeNameFQDN() string {
    return fmt.Sprintf(coreNodeName + "." + pcrypto.FormFQDNClusterID, c.ClusterID)
}

// --- Cluster ID ---
func (c *coreContext) CoreClusterID() string {
    return c.ClusterID
}

func (c *coreContext) CoreAuthServerAddr() (string, error) {
    inif, err := InternalNetworkInterface()
    if err != nil {
        return "", err
    }
    return inif.GatewayAddr, nil
}

func CoreSSHAdvertiseAddr() string {
    return "127.0.0.1"
}

// --- Master IP4 Address ---
func (c *coreContext) SetMasterIP4ExtAddr(ip4Address string) error {
    c.Lock()
    defer c.Unlock()

    if len(ip4Address) == 0 {
        return errors.Errorf("[ERR] invalid master ip4 address")
    }
    c.MasterSection.MasterIP4Address = ip4Address
    return nil
}

func (c *coreContext) GetMasterIP4ExtAddr() (string, error) {
    c.Lock()
    defer c.Unlock()

    if len(c.MasterSection.MasterIP4Address) == 0 {
        return "", errors.Errorf("[ERR] empty master ip4 address")
    }
    return c.MasterSection.MasterIP4Address , nil
}
