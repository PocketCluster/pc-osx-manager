package pkg

import (
    "testing"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pcrypto"
    . "gopkg.in/check.v1"
)

const (
    extIpAddrSmMask string = "192.168.1.105/24"
    extGateway      string = "192.168.1.1"
    clusterID       string = "ZKYQbwGnKJfFRTcW"
)

func TestCorePackage(t *testing.T) { TestingT(t) }

type PackageTestSuite struct {
    encryptor  pcrypto.RsaEncryptor
    decryptor  pcrypto.RsaDecryptor
}

var _ = Suite(&PackageTestSuite{})

func (p *PackageTestSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
    p.encryptor, _ = pcrypto.NewRsaEncryptorFromKeyData(pcrypto.TestMasterStrongPublicKey(), pcrypto.TestSlaveNodePrivateKey())
    p.decryptor, _ = pcrypto.NewRsaDecryptorFromKeyData(pcrypto.TestSlaveNodePublicKey(),    pcrypto.TestMasterStrongPrivateKey())
}

func (p *PackageTestSuite) TearDownSuite(c *C) {
}

func (p *PackageTestSuite) SetUpTest(c *C) {
    log.Debugf("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
}

func (p *PackageTestSuite) TearDownTest(c *C) {
    log.Debugf("\n\n")
}

// ---

func (p *PackageTestSuite) Test_BoundedStatus_Package(c *C) {
    metaPackage, err := CorePackingBoundedStatus(clusterID, extIpAddrSmMask, extGateway, p.encryptor)
    c.Assert(err, IsNil)

    meta, err := CoreUnpackingStatus(clusterID, metaPackage, p.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.CoreStatus.CoreState, Equals, VBoxCoreBounded)
}

func (p *PackageTestSuite) Test_BindBrokenStatus_Package(c *C) {
    metaPackage, err := CorePackingBindBrokenStatus(clusterID, extIpAddrSmMask, extGateway, p.encryptor)
    c.Assert(err, IsNil)

    meta, err := CoreUnpackingStatus(clusterID, metaPackage, p.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.CoreStatus.CoreState, Equals, VBoxCoreBindBroken)
}
