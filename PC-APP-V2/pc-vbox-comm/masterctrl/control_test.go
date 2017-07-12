package masterctrl

import (
    "os"
    "testing"
    "time"

    . "gopkg.in/check.v1"
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/model"

    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
)

const (
    authToken           string = "bjAbqvJVCy2Yr2suWu5t2ZnD4Z5336oNJ0bBJWFZ4A0="
    coreExtIpAddrSmMask string = "192.168.1.105/24"
    coreExtGateway      string = "192.168.1.1"
)

func TestMasterControl(t *testing.T) { TestingT(t) }

type coreProperties struct {
    // Core node properties
    publicKey        []byte
    privateKey       []byte
    encryptor        pcrypto.RsaEncryptor
    decryptor        pcrypto.RsaDecryptor
    timestamp        time.Time
}

type MasterControlTestSuite struct {
    timestamp        time.Time
    master           *masterControl

    // core node properties
    core             *coreProperties
}

var _ = Suite(&MasterControlTestSuite{})

func (m *MasterControlTestSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
}

func (m *MasterControlTestSuite) TearDownSuite(c *C) {
}

func (m *MasterControlTestSuite) SetUpTest(c *C) {
    context.DebugContextPrepare()
    model.DebugRecordGatePrepare(os.Getenv("TMPDIR"))

    // setup init time
    m.timestamp, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

    // setup core
    core := &coreProperties {
        publicKey:     pcrypto.TestSlaveNodePublicKey(),
        privateKey:    pcrypto.TestSlaveNodePrivateKey(),
        timestamp:     m.timestamp,
    }
    m.core = core

    log.Debugf("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
}

func (m *MasterControlTestSuite) TearDownTest(c *C) {
    log.Debugf("\n\n")
    m.core = nil
    model.DebugRecordGateDestroy(os.Getenv("TMPDIR"))
    context.DebugContextDestroy()
}

// --- Test Body ---

func (m *MasterControlTestSuite) TestCoreNodeJoin(c *C) {
    // setup test specifics
    {
        // setup core node
        coreNode := model.RetrieveCoreNode()
        coreNode.SetAuthToken(authToken)
        coreNode.CreateCore()

        // setup controller
        ctrl, err := NewVBoxMasterControl(pcrypto.TestMasterStrongPrivateKey(), pcrypto.TestMasterStrongPublicKey(), coreNode, nil)
        if err != nil {
            log.Panic(err.Error())
        }
        m.master = ctrl.(*masterControl)
    }

    // core report
    metaPackage, err := cpkg.CorePackingUnboundedStatus(m.core.publicKey)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)

    // master read and ack
    m.timestamp = m.core.timestamp.Add(time.Second)
    metaPackage, err = m.master.ReadCoreMetaAndMakeMasterAck(nil, metaPackage, m.timestamp)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)
    c.Assert(m.master.CurrentState(), Equals, mpkg.VBoxMasterKeyExchange)

    // core read
    m.core.timestamp = m.timestamp.Add(time.Second)
    meta, err := mpkg.MasterUnpackingAcknowledge(metaPackage, m.core.privateKey, nil)
    c.Assert(err, IsNil)
    c.Assert(meta.MasterState, Equals, mpkg.VBoxMasterKeyExchange)
    c.Assert(meta.MasterAcknowledge.AuthToken, Equals, authToken)
    c.Assert(meta.Encryptor, NotNil)
    c.Assert(meta.Decryptor, NotNil)
    m.core.encryptor = meta.Encryptor
    m.core.decryptor = meta.Decryptor

    // core report
    m.core.timestamp = m.core.timestamp.Add(time.Second)
    metaPackage, err = cpkg.CorePackingBoundedStatus(coreExtIpAddrSmMask, coreExtGateway, m.core.encryptor)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)

    // master read and ack
    m.timestamp = m.core.timestamp.Add(time.Second)
    metaPackage, err = m.master.ReadCoreMetaAndMakeMasterAck(nil, metaPackage, m.timestamp)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)
    c.Assert(m.master.CurrentState(), Equals, mpkg.VBoxMasterBounded)

    // core read
    m.core.timestamp = m.timestamp.Add(time.Second)
    meta, err = mpkg.MasterUnpackingAcknowledge(metaPackage, nil, m.core.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.MasterState, Equals, mpkg.VBoxMasterBounded)
    c.Assert(meta.MasterAcknowledge.AuthToken, Equals, "")
    c.Assert(meta.Encryptor, Equals, nil)
    c.Assert(meta.Decryptor, Equals, nil)
}

func (m *MasterControlTestSuite) TestCoreNodeBindRecovery(c *C) {
    // setup test specifics
    {
        // setup core node
        coreNode := model.RetrieveCoreNode()
        coreNode.SetAuthToken(authToken)
        coreNode.CreateCore()
        coreNode.PublicKey = m.core.publicKey
        coreNode.IP4Address = coreExtIpAddrSmMask
        coreNode.IP4Gateway = coreExtGateway
        coreNode.JoinCore()

        // re-setup controller
        ctrl, err := NewVBoxMasterControl(pcrypto.TestMasterStrongPrivateKey(), pcrypto.TestMasterStrongPublicKey(), coreNode, nil)
        if err != nil {
            log.Panic(err.Error())
        }
        m.master = ctrl.(*masterControl)
    }

    // re-setup core encryptor & decryptor
    encryptor, _ := pcrypto.NewRsaEncryptorFromKeyData(pcrypto.TestMasterStrongPublicKey(), m.core.privateKey)
    decryptor, _ := pcrypto.NewRsaDecryptorFromKeyData(pcrypto.TestMasterStrongPublicKey(), m.core.privateKey)
    m.core.encryptor = encryptor
    m.core.decryptor = decryptor


    // core report
    metaPackage, err := cpkg.CorePackingBindBrokenStatus(coreExtIpAddrSmMask, coreExtGateway, m.core.encryptor)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)

    // master read and ack
    m.timestamp = m.core.timestamp.Add(time.Second)
    metaPackage, err = m.master.ReadCoreMetaAndMakeMasterAck(nil, metaPackage, m.timestamp)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)
    c.Assert(m.master.CurrentState(), Equals, mpkg.VBoxMasterBounded)

    // core read
    m.core.timestamp = m.timestamp.Add(time.Second)
    meta, err := mpkg.MasterUnpackingAcknowledge(metaPackage, nil, m.core.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.MasterState, Equals, mpkg.VBoxMasterBounded)
    c.Assert(meta.MasterAcknowledge.AuthToken, Equals, "")
    c.Assert(meta.Encryptor, Equals, nil)
    c.Assert(meta.Decryptor, Equals, nil)

    // core report
    m.core.timestamp = m.core.timestamp.Add(time.Second)
    metaPackage, err = cpkg.CorePackingBoundedStatus(coreExtIpAddrSmMask, coreExtGateway, m.core.encryptor)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)

    // master read and ack
    m.timestamp = m.core.timestamp.Add(time.Second)
    metaPackage, err = m.master.ReadCoreMetaAndMakeMasterAck(nil, metaPackage, m.timestamp)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)
    c.Assert(m.master.CurrentState(), Equals, mpkg.VBoxMasterBounded)

    // core read
    m.core.timestamp = m.timestamp.Add(time.Second)
    meta, err = mpkg.MasterUnpackingAcknowledge(metaPackage, nil, m.core.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.MasterState, Equals, mpkg.VBoxMasterBounded)
    c.Assert(meta.MasterAcknowledge.AuthToken, Equals, "")
    c.Assert(meta.Encryptor, Equals, nil)
    c.Assert(meta.Decryptor, Equals, nil)
}