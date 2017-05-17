package beacon

import (
    "testing"

    log "github.com/Sirupsen/logrus"
    . "gopkg.in/check.v1"
    "github.com/stkim1/pc-core/model"
)

func TestRecord(t *testing.T) { TestingT(t) }

type ManagerSuite struct {
    dataDir     string
}

var _ = Suite(&ManagerSuite{})

func (s *ManagerSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
}

func (s *ManagerSuite) TearDownSuite(c *C) {
}

func (s *ManagerSuite) SetUpTest(c *C) {
    var err error

    s.dataDir = c.MkDir()
    _, err = model.DebugRecordGatePrepare(s.dataDir)
    c.Assert(err, IsNil)
}

func (s *ManagerSuite) TearDownTest(c *C) {
    c.Assert(model.DebugRecordGateDestroy(s.dataDir), IsNil)
}
