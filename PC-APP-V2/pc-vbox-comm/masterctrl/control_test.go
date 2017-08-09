package masterctrl

import (
    "os"
    "testing"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pcrypto"
    . "gopkg.in/check.v1"

    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
)

const (
    authToken           string = "bjAbqvJVCy2Yr2suWu5t2ZnD4Z5336oNJ0bBJWFZ4A0="
    clusterID           string = "ZKYQbwGnKJfFRTcW"
    masterExtIp4Addr    string = "192.168.1.100"
    coreExtIp4AdrSmMask string = "192.168.1.105/24"
    coreExtIp4Gateway   string = "192.168.1.1"
)

func TestMasterControl(t *testing.T) { TestingT(t) }

type coreProperty struct {
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
    core             *coreProperty
}

var _ = Suite(&MasterControlTestSuite{})

func (m *MasterControlTestSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
}

func (m *MasterControlTestSuite) TearDownSuite(c *C) {}

func (m *MasterControlTestSuite) SetUpTest(c *C) {
    context.DebugContextPrepare()
    model.DebugRecordGatePrepare(os.Getenv("TMPDIR"))

    // setup init time
    m.timestamp, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

    log.Debugf("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
}

func (m *MasterControlTestSuite) TearDownTest(c *C) {
    log.Debugf("\n\n")
    m.core = nil
    m.master = nil
    model.DebugRecordGateDestroy(os.Getenv("TMPDIR"))
    context.DebugContextDestroy()
}

func (m *MasterControlTestSuite) prepareBindBrokenMasterCore() error {
    // setup core node
    coreNode := model.RetrieveCoreNode()
    coreNode.SetAuthToken(authToken)
    coreNode.PublicKey   = pcrypto.TestSlaveNodePublicKey()
    coreNode.PrivateKey  = pcrypto.TestSlaveNodePrivateKey()
    err := coreNode.CreateCore()
    if err != nil {
        return err
    }

    coreNode.IP4Address = coreExtIp4AdrSmMask
    coreNode.IP4Gateway = coreExtIp4Gateway
    err = coreNode.Update()
    if err != nil {
        return err
    }

    // re-setup controller
    ctrl, err := NewVBoxMasterControl(clusterID, masterExtIp4Addr, pcrypto.TestMasterStrongPrivateKey(), pcrypto.TestMasterStrongPublicKey(), coreNode, nil)
    if err != nil {
        return err
    }
    m.master = ctrl.(*masterControl)

    // setup core
    encryptor, err := pcrypto.NewRsaEncryptorFromKeyData(pcrypto.TestMasterStrongPublicKey(), pcrypto.TestSlaveNodePrivateKey())
    if err != nil {
        return err
    }
    decryptor, err := pcrypto.NewRsaDecryptorFromKeyData(pcrypto.TestMasterStrongPublicKey(), pcrypto.TestSlaveNodePrivateKey())
    if err != nil {
        return err
    }
    core := &coreProperty{
        publicKey:     pcrypto.TestSlaveNodePublicKey(),
        privateKey:    pcrypto.TestSlaveNodePrivateKey(),
        timestamp:     m.timestamp,
        encryptor:     encryptor,
        decryptor:     decryptor,
    }
    m.core = core
    return nil
}

// --- Test Body ---

func (m *MasterControlTestSuite) TestCoreNodeBindRecovery(c *C) {
    // setup test specifics
    err := m.prepareBindBrokenMasterCore()
    if err != nil {
        log.Panic(err.Error())
    }

    // check master status
    c.Assert(m.master.CurrentState(), Equals, mpkg.VBoxMasterBindBroken)

    // core report
    metaPackage, err := cpkg.CorePackingBindBrokenStatus(clusterID, coreExtIp4AdrSmMask, coreExtIp4Gateway, m.core.encryptor)
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
    meta, err := mpkg.MasterUnpackingAcknowledge(clusterID, metaPackage, m.core.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.MasterAcknowledge.MasterState, Equals, mpkg.VBoxMasterBounded)

    // core report
    m.core.timestamp = m.core.timestamp.Add(time.Second)
    metaPackage, err = cpkg.CorePackingBoundedStatus(clusterID, coreExtIp4AdrSmMask, coreExtIp4Gateway, m.core.encryptor)
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
    meta, err = mpkg.MasterUnpackingAcknowledge(clusterID, metaPackage, m.core.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.MasterAcknowledge.MasterState, Equals, mpkg.VBoxMasterBounded)
}