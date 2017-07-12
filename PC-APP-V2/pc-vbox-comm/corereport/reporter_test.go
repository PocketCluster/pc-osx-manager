package corereport

import (
    "testing"
    "time"

    . "gopkg.in/check.v1"
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pcrypto"

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

func (r *CoreReportTestSuite) Test_Unbounded_Joining_To_Master(c *C) {
    // setup test specifics
    err := r.prepareUnboundedCoreMaster()
    if err != nil {
        log.Panic(err.Error())
    }

}

func (r *CoreReportTestSuite) Test_BindBroken_Joining_To_Master(c *C) {
    // setup test specifics
    err := r.prepareBindBrokenCoreMaster()
    if err != nil {
        log.Panic(err.Error())
    }
}