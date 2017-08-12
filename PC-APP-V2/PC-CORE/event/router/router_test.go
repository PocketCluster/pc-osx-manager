package router

import (
    "testing"

    . "gopkg.in/check.v1"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

func TestRouter(t *testing.T) { TestingT(t) }

type RouterSuite struct {
    router    *Router
}

var _ = Suite(&RouterSuite{})

func (s *RouterSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
}

func (s *RouterSuite) TearDownSuite(c *C) {
}

func (s *RouterSuite) SetUpTest(c *C) {
    log.Debugf("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
    s.router = NewRouter(func() error {
        return errors.Errorf("/ path should not be accessed")
    })
}

func (s *RouterSuite) TearDownTest(c *C) {
    log.Debugf("\n\n")
}

/// ---


func (s *RouterSuite) Test_GetBasicTest(c *C) {
    var (
        handleVar = ""
    )
    s.router.GET("/v1/system/monitor", func() error {
        handleVar = "test"
        return nil
    })
}