package pkg

import (
    "testing"

    . "gopkg.in/check.v1"
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pcrypto"
)

const (
    authToken string = "bjAbqvJVCy2Yr2suWu5t2ZnD4Z5336oNJ0bBJWFZ4A0="
    clusterID string = "ZKYQbwGnKJfFRTcW"
    extIP4Adr string = "192.168.1.105"
)

func TestMasterPackage(t *testing.T) { TestingT(t) }

type PackageTestSuite struct {
    publicKey  []byte
    privateKey []byte
    encryptor  pcrypto.RsaEncryptor
    decryptor  pcrypto.RsaDecryptor
}

var _ = Suite(&PackageTestSuite{})

func (p *PackageTestSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
    p.publicKey    = pcrypto.TestMasterStrongPublicKey()
    p.privateKey   = pcrypto.TestMasterStrongPrivateKey()
    p.encryptor, _ = pcrypto.NewRsaEncryptorFromKeyData(pcrypto.TestSlaveNodePublicKey(), pcrypto.TestMasterStrongPrivateKey())
    p.decryptor, _ = pcrypto.NewRsaDecryptorFromKeyData(pcrypto.TestMasterStrongPublicKey(), pcrypto.TestSlaveNodePrivateKey())
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

func (p *PackageTestSuite) Test_KeyExchangeStatus_Package(c *C) {
    // master side
    metaPackage, err := MasterPackingKeyExchangeAcknowledge(clusterID, authToken, extIP4Adr, p.publicKey, p.encryptor)
    c.Assert(err, IsNil)

    // core side
    meta, err := MasterUnpackingAcknowledge(metaPackage, pcrypto.TestSlaveNodePrivateKey(), nil)
    c.Assert(err, IsNil)
    c.Assert(meta.MasterState, Equals, VBoxMasterKeyExchange)
}

func (p *PackageTestSuite) Test_BoundedStatus_Package(c *C) {
    // master side
    metaPackage, err := MasterPackingBoundedAcknowledge(clusterID, extIP4Adr, p.encryptor)
    c.Assert(err, IsNil)

    // core side
    meta, err := MasterUnpackingAcknowledge(metaPackage, nil, p.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.MasterState, Equals, VBoxMasterBounded)
}

func (p *PackageTestSuite) Test_BindBrokenStatus_Package(c *C) {
    // master side
    metaPackage, err := MasterPackingBindBrokenAcknowledge(clusterID, extIP4Adr, p.encryptor)
    c.Assert(err, IsNil)

    // core side
    meta, err := MasterUnpackingAcknowledge(metaPackage, nil, p.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.MasterState, Equals, VBoxMasterBindBroken)
}
