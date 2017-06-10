package service

import (
    "testing"

    . "gopkg.in/check.v1"
    log "github.com/Sirupsen/logrus"
)

const (
    exitValue       = "no_error"
    testService1    = "test_service1"
    testService2    = "test_service2"
    testService3    = "test_service3"
)

func TestSupervisor(t *testing.T) { TestingT(t) }

type SupervisorSuite struct {
    app    AppSupervisor
}

var _ = Suite(&SupervisorSuite{})

func (s *SupervisorSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
}

func (s *SupervisorSuite) TearDownSuite(c *C) {
}

func (s *SupervisorSuite) SetUpTest(c *C) {
    s.app = NewAppSupervisor()
}

func (s *SupervisorSuite) TearDownTest(c *C) {
    err := s.app.Stop()
    c.Assert(err, IsNil)
    s.app = nil
}

/// ---

func (s *SupervisorSuite) TestStartStop(c *C) {
    err := s.app.Start()
    c.Assert(err, IsNil)

    err = s.app.Stop()
    c.Assert(err, IsNil)
}

func (s *SupervisorSuite) Test_UnamedService_Run_After_Start(c *C) {
    var(
        exitLatch = make(chan bool)
        exitChecker = ""
    )
    err := s.app.Start()
    c.Assert(err, IsNil)

    err = s.app.RegisterServiceWithFuncs(
        func() error {
            for {
                select {
                    case <- s.app.StopChannel():
                        return nil
                    default:
                }
            }
        },
        func(_ func(interface{})) error {
            exitChecker = exitValue
            exitLatch <- true
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.ServiceCount(), Equals, 1)

    err = s.app.Stop()
    c.Assert(err, IsNil)
    <-exitLatch

    c.Check(exitChecker, Equals, exitValue)
    c.Assert(s.app.ServiceCount(), Equals, 0)
    close(exitLatch)
}

func (s *SupervisorSuite) Test_UnnamedService_Register_Before_Start(c *C) {
    var(
        exitLatch = make(chan bool)
        exitChecker = ""
    )
    err := s.app.RegisterServiceWithFuncs(
        func() error {
            for {
                select {
                    case <- s.app.StopChannel():
                        return nil
                    default:
                }
            }
        },
        func(_ func(interface{})) error {
            exitChecker = exitValue
            exitLatch <- true
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.ServiceCount(), Equals, 1)

    err = s.app.Start()
    c.Assert(err, IsNil)

    err = s.app.Stop()
    c.Assert(err, IsNil)
    <-exitLatch

    c.Check(exitChecker, Equals, exitValue)
    c.Assert(s.app.ServiceCount(), Equals, 0)
    close(exitLatch)
}

func (s *SupervisorSuite) Test_NamedService_Unsycned_Stop(c *C) {
    var(
        exitSignal = make(chan bool)
        exitLatch = make(chan bool)
        exitChecker = ""
    )
    err := s.app.RegisterServiceWithFuncs(
        func() error {
            for {
                select {
                    case <- exitSignal:
                        log.Debug("finishing serve()...")
                        return nil
                    default:
                }
            }
        },
        func(_ func(interface{})) error {
            exitChecker = exitValue
            exitLatch <- true
            return nil
        },
        MakeServiceNamed(testService1),
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.ServiceCount(), Equals, 1)

    err = s.app.Start()
    c.Assert(err, IsNil)

    err = s.app.RunNamedService(testService1)
    c.Assert(err, IsNil)

    exitSignal <- true
    <-exitLatch

    c.Check(exitChecker, Equals, exitValue)
    c.Assert(s.app.ServiceCount(), Equals, 1)
    close(exitLatch)
}

func (s *SupervisorSuite) Test_NamedService_Sycned_Stop(c *C) {
    var(
        exitLatch = make(chan bool)
        exitChecker = ""
    )
    err := s.app.RegisterServiceWithFuncs(
        func() error {
            for {
                select {
                    case <- s.app.StopChannel():
                        log.Debug("finishing serve()...")
                        return nil
                    default:
                }
            }
        },
        func(_ func(interface{})) error {
            exitChecker = exitValue
            exitLatch <- true
            return nil
        },
        MakeServiceNamed(testService2),
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.ServiceCount(), Equals, 1)

    err = s.app.Start()
    c.Assert(err, IsNil)

    err = s.app.RunNamedService(testService2)
    c.Assert(err, IsNil)

    err = s.app.Stop()
    c.Assert(err, IsNil)
    <-exitLatch

    c.Check(exitChecker, Equals, exitValue)
    c.Assert(s.app.ServiceCount(), Equals, 1)
    close(exitLatch)
}