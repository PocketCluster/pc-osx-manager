package crcontext

import (
    "fmt"

    "github.com/stkim1/pc-vbox-core/crcontext/config"
    "github.com/stkim1/pcrypto"
)

const (
    coreNodeName = "pc-core"
)

// for teleport to build configuration
type PocketCoreSSHInfo interface {
    CoreNodeName() string
    CoreNodeNameFQDN() string

    CoreNodeUUID() string
    CoreAuthToken() string

    CoreConfigPath() string
    CoreKeyAndCertPath() string

    CoreSSHCertificateFileName() string
    CoreSSHPrivateKeyFileName() string
    CoreSSHAdvertiseAddr() string

    CoreAuthServerAddr() (string, error)
}

func (c *coreContext) CoreNodeName() string {
    return coreNodeName
}

// TODO : add tests
func (c *coreContext) CoreNodeNameFQDN() string {
    return fmt.Sprintf(coreNodeName + "." + pcrypto.FormFQDNClusterID, c.config.ClusterID)
}

func (c *coreContext) CoreNodeUUID() string {
    return c.config.CoreSection.CoreNodeUUID
}

func (c *coreContext) CoreAuthToken() string {
    return c.config.CoreSection.CoreAuthToken
}

// TODO : add tests
func (c *coreContext) CoreConfigPath() string {
    return config.DirPathCoreConfig(c.config.RootPath())
}

// TODO : add tests
func (c *coreContext) CoreKeyAndCertPath() string {
    return config.DirPathCoreCerts(c.config.RootPath())
}

// TODO : add tests
func (c *coreContext) CoreSSHCertificateFileName() string {
    return config.FilePathCoreSSHKeyCert(c.config.RootPath())
}

func (c *coreContext) CoreSSHAdvertiseAddr() string {
    return "127.0.0.1"
}

// TODO : add tests
func (c *coreContext) CoreSSHPrivateKeyFileName() string {
    return config.FilePathCoreSSHPrivateKey(c.config.RootPath())
}

func (c *coreContext) CoreAuthServerAddr() (string, error) {
    inif, err := InternalNetworkInterface()
    if err != nil {
        return "", err
    }
    return inif.GatewayAddr, nil
}