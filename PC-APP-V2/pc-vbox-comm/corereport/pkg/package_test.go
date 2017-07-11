package pkg

import (
    "testing"

    . "gopkg.in/check.v1"
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pcrypto"
)

const (
    extIpAddrSmMask string = "192.168.1.105/24"
    extGateway      string = "192.168.1.1"
)

func TestCorePackage(t *testing.T) { TestingT(t) }

type PackageTestSuite struct {
    privateKey []byte
    publicKey  []byte
    encryptor  pcrypto.RsaEncryptor
    decryptor  pcrypto.RsaDecryptor
}

var _ = Suite(&PackageTestSuite{})

func (p *PackageTestSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
    p.publicKey  = pcrypto.TestSlavePublicKey()
    p.encryptor  = pcrypto.TestSlaveRSAEncryptor
    p.decryptor  = pcrypto.TestMasterRSADecryptor
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

func (p *PackageTestSuite) Test_Unbounded_Status_Package(c *C) {
    metaPackage, err := CorePackingUnboundedStatus(p.publicKey)
    c.Assert(err, IsNil)

    meta, err := CoreUnpackingStatus(metaPackage, nil)
    c.Assert(err, IsNil)
    c.Assert(meta.CoreState, Equals, VBoxCoreUnbounded)
}

func (p *PackageTestSuite) Test_BoundedStatus_Package(c *C) {
    metaPackage, err := CorePackingBoundedStatus(extIpAddrSmMask, extGateway, p.encryptor)
    c.Assert(err, IsNil)

    meta, err := CoreUnpackingStatus(metaPackage, p.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.CoreState, Equals, VBoxCoreBounded)
}

func (p *PackageTestSuite) Test_BindBrokenStatus_Package(c *C) {
    metaPackage, err := CorePackingBindBrokenStatus(extIpAddrSmMask, extGateway, p.encryptor)
    c.Assert(err, IsNil)

    meta, err := CoreUnpackingStatus(metaPackage, p.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.CoreState, Equals, VBoxCoreBindBroken)
}
