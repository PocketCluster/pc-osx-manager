package corereport

import (
    "testing"
    "time"

    . "gopkg.in/check.v1"
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pcrypto"
)

func TestCoreReport(t *testing.T) { TestingT(t) }

type CoreReportTestSuite struct {
    CoreReporter    *coreReporter
    InitTime         time.Time
}

var _ = Suite(&CoreReportTestSuite{})

func (r *CoreReportTestSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
    r.CoreReporter.publicKey  = pcrypto.TestSlaveNodePublicKey()
    r.CoreReporter.privateKey = pcrypto.TestSlaveNodePrivateKey()
}

func (r *CoreReportTestSuite) TearDownSuite(c *C) {
}

func (r *CoreReportTestSuite) SetUpTest(c *C) {
    log.Debugf("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
}

func (r *CoreReportTestSuite) TearDownTest(c *C) {
    log.Debugf("\n\n")
}

// --- Test Body ---

