package corereport

import (
    "testing"
    "time"

    . "gopkg.in/check.v1"
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pcrypto"

    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
)

const (
    authToken           string = "bjAbqvJVCy2Yr2suWu5t2ZnD4Z5336oNJ0bBJWFZ4A0="
    coreExtIpAddrSmMask string = "192.168.1.105/24"
    coreExtGateway      string = "192.168.1.1"
)

func TestCoreReport(t *testing.T) { TestingT(t) }

type masterProperty struct {
    // Core node properties
    publicKey     []byte
    privateKey    []byte
    encryptor     pcrypto.RsaEncryptor
    decryptor     pcrypto.RsaDecryptor
    timestamp     time.Time
}

type CoreReportTestSuite struct {
    core          *coreReporter
    timestamp     time.Time

    // master properties
    master        *masterProperty
}

var _ = Suite(&CoreReportTestSuite{})

func (r *CoreReportTestSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)

    // setup init time
    r.timestamp, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
}

func (r *CoreReportTestSuite) TearDownSuite(c *C) {}

func (r *CoreReportTestSuite) SetUpTest(c *C) {
    log.Debugf("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
}

func (r *CoreReportTestSuite) TearDownTest(c *C) {
    log.Debugf("\n\n")
    r.core = nil
    r.master = nil
}

func (r *CoreReportTestSuite) prepareUnboundedCoreMaster() error {
    m := &masterProperty {
        privateKey:    pcrypto.TestMasterStrongPrivateKey(),
        publicKey:     pcrypto.TestMasterStrongPublicKey(),
        timestamp:     r.timestamp,
    }
    r.master = m

    core, err := NewCoreReporter(cpkg.VBoxCoreUnbounded, pcrypto.TestSlaveNodePrivateKey(), pcrypto.TestSlaveNodePublicKey(), nil)
    if err != nil {
        return err
    }
    r.core = core.(*coreReporter)
    return nil
}

func (r *CoreReportTestSuite) prepareBindBrokenCoreMaster() error {
    encryptor, err := pcrypto.NewRsaEncryptorFromKeyData(pcrypto.TestSlaveNodePublicKey(), pcrypto.TestMasterStrongPrivateKey())
    if err != nil {
        return err
    }
    decryptor, err := pcrypto.NewRsaDecryptorFromKeyData(pcrypto.TestSlaveNodePublicKey(), pcrypto.TestMasterStrongPrivateKey())
    if err != nil {
        return err
    }
    m := &masterProperty {
        privateKey:    pcrypto.TestMasterStrongPrivateKey(),
        publicKey:     pcrypto.TestMasterStrongPublicKey(),
        encryptor:     encryptor,
        decryptor:     decryptor,
        timestamp:     r.timestamp,
    }
    r.master = m

    core, err := NewCoreReporter(cpkg.VBoxCoreBindBroken, pcrypto.TestSlaveNodePrivateKey(), pcrypto.TestSlaveNodePublicKey(), pcrypto.TestMasterStrongPublicKey())
    if err != nil {
        return err
    }
    r.core = core.(*coreReporter)
    return nil
}

// --- Test Body ---

func (r *CoreReportTestSuite) Test_Unbounded_Core_Joining_To_Master(c *C) {
    // setup test specifics
    err := r.prepareUnboundedCoreMaster()
    if err != nil {
        log.Panic(err.Error())
    }

    // check core state
    c.Assert(r.core.CurrentState(), Equals, cpkg.VBoxCoreUnbounded)

    // core makes report
    metaPackage, err := r.core.MakeCoreReporter(r.timestamp)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)

    // master read
    r.master.timestamp = r.timestamp.Add(time.Second)
    meta, err := cpkg.CoreUnpackingStatus(metaPackage, nil)
    c.Assert(err, IsNil)
    c.Assert(meta.CoreState, Equals, cpkg.VBoxCoreUnbounded)

    // master build encryptor & decryptor
    r.master.timestamp = r.master.timestamp.Add(time.Second)
    encryptor, err := pcrypto.NewRsaEncryptorFromKeyData(meta.PublicKey, pcrypto.TestMasterStrongPrivateKey())
    c.Assert(err, IsNil)
    r.master.encryptor = encryptor
    decryptor, err := pcrypto.NewRsaDecryptorFromKeyData(meta.PublicKey, pcrypto.TestMasterStrongPrivateKey())
    c.Assert(err, IsNil)
    r.master.decryptor = decryptor

    // master make acknowledge
    r.master.timestamp = r.master.timestamp.Add(time.Second)
    metaPackage, err = mpkg.MasterPackingKeyExchangeAcknowledge(authToken, r.master.publicKey, r.master.encryptor)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)

    // core read acknowledge
    r.timestamp = r.master.timestamp.Add(time.Second)
    err = r.core.ReadMasterAcknowledgement(metaPackage, r.master.timestamp)
    c.Assert(err, IsNil)
    c.Assert(r.core.CurrentState(), Equals, cpkg.VBoxCoreBounded)
    c.Assert(r.core.authToken, Equals, authToken)
}

func (r *CoreReportTestSuite) Test_BindBroken_Core_Joining_To_Master(c *C) {
    // setup test specifics
    err := r.prepareBindBrokenCoreMaster()
    if err != nil {
        log.Panic(err.Error())
    }

    // check core state
    c.Assert(r.core.CurrentState(), Equals, cpkg.VBoxCoreBindBroken)

    // core makes report
    metaPackage, err := r.core.MakeCoreReporter(r.timestamp)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)

    // master read
    r.master.timestamp = r.timestamp.Add(time.Second)
    meta, err := cpkg.CoreUnpackingStatus(metaPackage, r.master.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.CoreState, Equals, cpkg.VBoxCoreBindBroken)

    // master make acknowledge
    r.master.timestamp = r.master.timestamp.Add(time.Second)
    metaPackage, err = mpkg.MasterPackingBindBrokenAcknowledge(r.master.encryptor)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)

    // core read
    r.timestamp = r.master.timestamp.Add(time.Second)
    err = r.core.ReadMasterAcknowledgement(metaPackage, r.master.timestamp)
    c.Assert(err, IsNil)
    c.Assert(r.core.CurrentState(), Equals, cpkg.VBoxCoreBounded)
    c.Assert(r.core.authToken, Equals, "")
}