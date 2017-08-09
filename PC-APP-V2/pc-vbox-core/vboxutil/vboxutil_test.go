package vboxutil

import (
    "os"
    "testing"

    . "gopkg.in/check.v1"
    log "github.com/Sirupsen/logrus"
)

func TestVboxUtil(t *testing.T) { TestingT(t) }

var _ = Suite(&VboxUtilSuite{})

type VboxUtilSuite struct {
    dataDir     string
}

func (s *VboxUtilSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
}

func (s *VboxUtilSuite) TearDownSuite(c *C) {}

func (s *VboxUtilSuite) SetUpTest(c *C) {
    s.dataDir = c.MkDir()
}

func (s *VboxUtilSuite) TearDownTest(c *C) {
    err := os.RemoveAll(s.dataDir)
    if err != nil {
        log.Panic(err.Error())
    }
}
