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
)

const (
    AuthToken string = "bjAbqvJVCy2Yr2suWu5t2ZnD4Z5336oNJ0bBJWFZ4A0="
)

func TestMasterControl(t *testing.T) { TestingT(t) }

type MasterControlTestSuite struct {
    MasterControl    *masterControl
    InitTime         time.Time
}

var _ = Suite(&MasterControlTestSuite{})

func (m *MasterControlTestSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
    context.DebugContextPrepare()
    model.DebugRecordGatePrepare(os.Getenv("TMPDIR"))

    // setup core node
    coreNode := model.NewCoreNode()
    coreNode.SetAuthToken(AuthToken)
    coreNode.JoinCore()

    // setup controller
    ctrl, err := NewVBoxMasterControl(pcrypto.TestMasterStrongPrivateKey(), pcrypto.TestMasterStrongPublicKey(), coreNode, nil)
    if err != nil {
        log.Panic(err.Error())
    }
    m.MasterControl = ctrl.(*masterControl)

    // setup init time
    m.InitTime, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
}

func (m *MasterControlTestSuite) TearDownSuite(c *C) {
    model.DebugRecordGateDestroy(os.Getenv("TMPDIR"))
    context.DebugContextDestroy()
}

func (m *MasterControlTestSuite) SetUpTest(c *C) {
    log.Debugf("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
}

func (m *MasterControlTestSuite) TearDownTest(c *C) {
    log.Debugf("\n\n")
}

// --- Test Body ---
