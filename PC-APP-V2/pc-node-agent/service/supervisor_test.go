package service

import (
    "testing"
    "time"

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
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.Stop()
    c.Assert(err, IsNil)
    <-exitLatch

    c.Check(exitChecker, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 0)
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
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.Start()
    c.Assert(err, IsNil)

    err = s.app.Stop()
    c.Assert(err, IsNil)
    <-exitLatch

    c.Check(exitChecker, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 0)
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
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.Start()
    c.Assert(err, IsNil)

    err = s.app.RunNamedService(testService1)
    c.Assert(err, IsNil)

    exitSignal <- true
    <-exitLatch

    c.Check(exitChecker, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 1)
    close(exitLatch)
}

func (s *SupervisorSuite) Test_NamedService_MultiCycle(c *C) {
    var(
        exitSignal = make(chan bool)
        exitLatch = make(chan bool)
        exitChecker = ""
    )
    // start service
    err := s.app.Start()
    c.Assert(err, IsNil)

    // add a service
    err = s.app.RegisterServiceWithFuncs(
        func() error {
            for {
                select {
                case <- exitSignal:
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
    c.Assert(s.app.serviceCount(), Equals, 1)

    // run multiple times
    for i := 0; i < 5; i++ {
        // start service
        err = s.app.RunNamedService(testService1)
        c.Assert(err, IsNil)

        // stop service
        exitSignal <- true
        <-exitLatch

        c.Check(exitChecker, Equals, exitValue)
        c.Assert(s.app.serviceCount(), Equals, 1)
    }
    // close everything
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
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.Start()
    c.Assert(err, IsNil)

    err = s.app.RunNamedService(testService2)
    c.Assert(err, IsNil)

    err = s.app.Stop()
    c.Assert(err, IsNil)
    <-exitLatch

    c.Check(exitChecker, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 1)
    close(exitLatch)
}


func (s *SupervisorSuite) Test_NamedServices_Sycned_Stop(c *C) {
    var(
        exitLatch1   = make(chan bool)
        exitChecker1 = ""
        exitLatch2   = make(chan bool)
        exitChecker2 = ""
    )
    // start services
    err := s.app.Start()
    c.Assert(err, IsNil)

    //register named service1
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
            exitChecker1 = exitValue
            exitLatch1 <- true
            return nil
        },
        MakeServiceNamed(testService1),
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    //register named service2
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
            exitChecker2 = exitValue
            exitLatch2 <- true
            return nil
        },
        MakeServiceNamed(testService2),
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 2)

    // run services
    err = s.app.RunNamedService(testService1)
    c.Assert(err, IsNil)
    err = s.app.RunNamedService(testService2)
    c.Assert(err, IsNil)

    // wait...
    time.Sleep(time.Second)

    // stop service
    err = s.app.Stop()
    c.Assert(err, IsNil)

    // stop service 1
    <-exitLatch1
    c.Check(exitChecker1, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 2)
    close(exitLatch1)

    // stop service 2
    <-exitLatch2
    c.Check(exitChecker2, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 2)
    close(exitLatch2)
}

func (s *SupervisorSuite) Test_NamedAndUnnamed_Services_Unsycned_Stop(c *C) {
    var(
        exitSignal   = make(chan bool)
        exitLatch1   = make(chan bool)
        exitChecker1 = ""
        exitLatch2   = make(chan bool)
        exitChecker2 = ""
    )
    //register named service
    err := s.app.RegisterServiceWithFuncs(
        func() error {
            for {
                select {
                    case <- exitSignal:
                        return nil
                    default:
                }
            }
        },
        func(_ func(interface{})) error {
            exitChecker1 = exitValue
            exitLatch1 <- true
            return nil
        },
        MakeServiceNamed(testService3),
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    //register unnamed service
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
            exitChecker2 = exitValue
            exitLatch2 <- true
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 2)

    // start services
    err = s.app.Start()
    c.Assert(err, IsNil)
    err = s.app.RunNamedService(testService3)
    c.Assert(err, IsNil)

    // stop service 1
    exitSignal <- true
    <-exitLatch1
    c.Check(exitChecker1, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 2)
    close(exitLatch1)

    // stop service 2
    err = s.app.Stop()
    c.Assert(err, IsNil)
    <-exitLatch2
    c.Check(exitChecker2, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 1)
    close(exitLatch2)
}