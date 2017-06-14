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
    app    ServiceSupervisor
}

var _ = Suite(&SupervisorSuite{})

func (s *SupervisorSuite) SetUpSuite(c *C) {
    log.SetLevel(log.DebugLevel)
}

func (s *SupervisorSuite) TearDownSuite(c *C) {
}

func (s *SupervisorSuite) SetUpTest(c *C) {
    log.Debugf("--- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- --- ---")
    s.app = NewServiceSupervisor()
}

func (s *SupervisorSuite) TearDownTest(c *C) {
    s.app = nil
    log.Debugf("\n\n")
}

/// ---

func (s *SupervisorSuite) Test_Start_Stop(c *C) {
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.StopServices()
    c.Assert(err, IsNil)
}

func (s *SupervisorSuite) Test_UnamedService_Run_After_Start(c *C) {
    var(
        exitLatch = make(chan string)
    )
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.RegisterServiceWithFuncs(
        func() (interface{}, error) {
            for {
                select {
                    case <- s.app.StopChannel():
                        log.Debugf("LET THIS SERVICE STOP")
                        return nil, nil
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch <- exitValue
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.StopServices()
    c.Assert(err, IsNil)
    c.Check(<-exitLatch, Equals, exitValue)

    c.Assert(s.app.serviceCount(), Equals, 0)
    close(exitLatch)
}

func (s *SupervisorSuite) Test_UnnamedService_Register_Before_Start(c *C) {
    var(
        exitLatch = make(chan string)
    )
    err := s.app.RegisterServiceWithFuncs(
        func() (interface{}, error) {
            for {
                select {
                    case <- s.app.StopChannel():
                        return nil, nil
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch <- exitValue
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.StopServices()
    c.Assert(err, IsNil)
    c.Check(<-exitLatch, Equals, exitValue)

    c.Assert(s.app.serviceCount(), Equals, 0)
    close(exitLatch)
}

func (s *SupervisorSuite) Test_NamedService_Unsycned_Stop(c *C) {
    var(
        exitSignal = make(chan bool)
        exitLatch = make(chan string)
    )
    err := s.app.RegisterNamedServiceWithFuncs(
        testService1,
        func() (interface{}, error) {
            for {
                select {
                    case <- exitSignal:
                        log.Debug("finishing serve()...")
                        return nil, nil
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch <- exitValue
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.RunNamedService(testService1)
    c.Assert(err, IsNil)

    exitSignal <- true
    c.Check(<-exitLatch, Equals, exitValue)

    c.Assert(s.app.serviceCount(), Equals, 1)
    close(exitLatch)
}

func (s *SupervisorSuite) Test_NamedService_MultiCycle(c *C) {
    var(
        exitSignal = make(chan bool)
        exitLatch = make(chan string)
    )
    // start service
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    // add a service
    err = s.app.RegisterNamedServiceWithFuncs(
        testService1,
        func() (interface{}, error) {
            for {
                select {
                    case <- exitSignal:
                        return nil, nil
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch <- exitValue
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    // run multiple times
    for i := 0; i < 5; i++ {
        // start service
        err = s.app.RunNamedService(testService1)
        c.Assert(err, IsNil)

        // stop service
        exitSignal <- true
        c.Check(<-exitLatch, Equals, exitValue)

        c.Assert(s.app.serviceCount(), Equals, 1)
        time.Sleep(time.Second)
    }
    // close everything
    close(exitSignal)
    close(exitLatch)
}

func (s *SupervisorSuite) Test_NamedService_Sycned_Stop(c *C) {
    var(
        exitLatch = make(chan string)
    )
    err := s.app.RegisterNamedServiceWithFuncs(
        testService2,
        func() (interface{}, error) {
            for {
                select {
                    case <- s.app.StopChannel():
                        log.Debug("finishing serve()...")
                        return nil, nil
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch <- exitValue
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.RunNamedService(testService2)
    c.Assert(err, IsNil)

    err = s.app.StopServices()
    c.Assert(err, IsNil)
    c.Check(<-exitLatch, Equals, exitValue)

    c.Assert(s.app.serviceCount(), Equals, 1)
    close(exitLatch)
}


func (s *SupervisorSuite) Test_NamedServices_Sycned_Stop(c *C) {
    var(
        exitLatch1   = make(chan string)
        exitLatch2   = make(chan string)
    )
    // start services
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    //register named service1
    err = s.app.RegisterNamedServiceWithFuncs(
        testService1,
        func() (interface{}, error) {
            for {
                select {
                    case <- s.app.StopChannel():
                        return nil, nil
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch1 <- exitValue
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    //register named service2
    err = s.app.RegisterNamedServiceWithFuncs(
        testService2,
        func() (interface{}, error) {
            for {
                select {
                    case <- s.app.StopChannel():
                        return nil, nil
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch2 <- exitValue
            return nil
        })
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
    err = s.app.StopServices()
    c.Assert(err, IsNil)

    // stop service 1
    c.Check(<-exitLatch1, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 2)
    close(exitLatch1)

    // stop service 2
    c.Check(<-exitLatch2, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 2)
    close(exitLatch2)
}

func (s *SupervisorSuite) Test_NamedAndUnnamed_Services_Unsycned_Stop(c *C) {
    var(
        exitSignal   = make(chan bool)
        exitLatch1   = make(chan string)
        exitLatch2   = make(chan string)
    )
    //register named service
    err := s.app.RegisterNamedServiceWithFuncs(
        testService3,
        func() (interface{}, error) {
            for {
                select {
                    case <- exitSignal:
                        return nil, nil
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch1 <- exitValue
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    //register unnamed service
    err = s.app.RegisterServiceWithFuncs(
        func() (interface{}, error) {
            for {
                select {
                    case <- s.app.StopChannel():
                        return nil, nil
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch2 <- exitValue
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 2)

    // start services
    err = s.app.StartServices()
    c.Assert(err, IsNil)
    err = s.app.RunNamedService(testService3)
    c.Assert(err, IsNil)

    // stop service 1
    exitSignal <- true
    c.Check(<-exitLatch1, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 2)
    close(exitLatch1)

    // stop service 2
    err = s.app.StopServices()
    c.Assert(err, IsNil)
    c.Check(<-exitLatch2, Equals, exitValue)

    time.Sleep(time.Second)
    c.Assert(s.app.serviceCount(), Equals, 1)
    close(exitLatch2)
    // close exit signal
    close(exitSignal)
}