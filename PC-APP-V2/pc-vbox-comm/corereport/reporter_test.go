package corereport

import (
    "testing"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pcrypto"
    . "gopkg.in/check.v1"

    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    "github.com/stkim1/pc-vbox-core/crcontext"
)

const (
    clusterID           string = "ZKYQbwGnKJfFRTcW"
    masterExtIP4Addr    string = "192.168.1.105"
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
    crcontext.DebugPrepareCoreContextWithRoot(c.MkDir())
}

func (r *CoreReportTestSuite) TearDownTest(c *C) {
    log.Debugf("\n\n")
    r.core = nil
    r.master = nil
    crcontext.DebugDestroyCoreContext()
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

    core, err := NewCoreReporter(clusterID, pcrypto.TestSlaveNodePrivateKey(), pcrypto.TestSlaveNodePublicKey(), pcrypto.TestMasterStrongPublicKey())
    if err != nil {
        return err
    }
    r.core = core.(*coreReporter)
    return nil
}

// --- Test Body ---
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
    meta, err := cpkg.CoreUnpackingStatus(clusterID, metaPackage, r.master.decryptor)
    c.Assert(err, IsNil)
    c.Assert(meta.CoreStatus.CoreState, Equals, cpkg.VBoxCoreBindBroken)

    // master make acknowledge
    r.master.timestamp = r.master.timestamp.Add(time.Second)
    metaPackage, err = mpkg.MasterPackingBindBrokenAcknowledge(clusterID, masterExtIP4Addr, r.master.encryptor)
    c.Assert(err, IsNil)
    c.Assert(len(metaPackage), Not(Equals), 0)

    // core read
    r.timestamp = r.master.timestamp.Add(time.Second)
    err = r.core.ReadMasterAcknowledgement(metaPackage, r.master.timestamp)
    c.Assert(err, IsNil)
    c.Assert(r.core.CurrentState(), Equals, cpkg.VBoxCoreBounded)
}