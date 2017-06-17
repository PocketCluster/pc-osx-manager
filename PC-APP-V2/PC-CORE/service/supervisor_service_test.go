package service

import (
    "testing"
    "time"

    . "gopkg.in/check.v1"
    log "github.com/Sirupsen/logrus"
    "fmt"
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

func (s *SupervisorSuite) Test_Service_Run_After_Start(c *C) {
    var(
        exitChecker string = ""
    )
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.RegisterServiceWithFuncs(
        testService1,
        func() error {
            for {
                select {
                    case <- s.app.StopChannel():
                        exitChecker = exitValue
                        log.Debugf("LET THIS SERVICE STOP")
                        return nil
                    default:
                }
            }
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.StopServices()
    c.Assert(err, IsNil)
    c.Check(exitChecker, Equals, exitValue)

    // check cleanup
    c.Assert(s.app.serviceCount(), Equals, 0)
}

func (s *SupervisorSuite) Test_Service_Register_Before_Start(c *C) {
    var(
        exitChecker string = ""
    )
    err := s.app.RegisterServiceWithFuncs(
        testService1,
        func() error {
            for {
                select {
                    case <- s.app.StopChannel():
                        exitChecker = exitValue
                        return nil
                    default:
                }
            }
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.StopServices()
    c.Assert(err, IsNil)
    c.Assert(exitChecker, Equals, exitValue)

    // check cleanup
    c.Assert(s.app.serviceCount(), Equals, 0)
}

func (s *SupervisorSuite) Test_Services_Sycned_Stop(c *C) {
    var(
        exitChecker1 string = ""
        exitChecker2 string = ""
    )
    // start services
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    //register named service1
    err = s.app.RegisterServiceWithFuncs(
        testService1,
        func() error {
            for {
                select {
                case <- s.app.StopChannel():
                    exitChecker1 = exitValue
                    return nil
                default:
                }
            }
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    //register named service2
    err = s.app.RegisterServiceWithFuncs(
        testService2,
        func() error {
            for {
                select {
                case <- s.app.StopChannel():
                    exitChecker2 = exitValue
                    return nil
                default:
                }
            }
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 2)

    // stop service
    err = s.app.StopServices()
    c.Assert(err, IsNil)

    // stop services
    c.Assert(exitChecker1, Equals, exitValue)
    c.Assert(exitChecker2, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 0)
}

func (s *SupervisorSuite) Test_Services_Iteration(c *C) {
    var(
        eventLatch  = make(chan string)
        exitLatch   = make(chan string)
    )
    // start services
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    for i := 0; i < 5; i++ {
        //register named service1
        err = s.app.RegisterServiceWithFuncs(
            testService1,
            func() error {
                for {
                    select {
                        case e := <- eventLatch:
                            exitLatch <- e
                            return nil
                        default:
                    }
                }
            })
        c.Assert(err, IsNil)
        c.Assert(s.app.serviceCount(), Equals, 1)

        eventLatch <- testValue1
        c.Check(<-exitLatch, Equals, testValue1)

        time.Sleep(time.Second)
        c.Assert(s.app.serviceCount(), Equals, 0)
    }

    // stop service
    c.Assert(s.app.serviceCount(), Equals, 0)
    err = s.app.StopServices()
    c.Assert(err, IsNil)
}

func (s *SupervisorSuite) Test_Supervisor_Iteration(c *C) {
    var(
        eventLatch  = make(chan string)
        err error   = nil
    )
    for i := 0; i < 5; i++ {
        // start services
        err = s.app.StartServices()
        c.Assert(err, IsNil)

        //register named service1
        var (
            eventC = make(chan Event)
            serviceTag = fmt.Sprintf("test_service%d",i)
        )
        err = s.app.RegisterServiceWithFuncs(
            serviceTag,
            func() error {
                for {
                    select {
                        case e := <- eventC:
                            eventLatch <- e.Payload.(string)
                        case <- s.app.StopChannel():
                            return nil
                        default:
                    }
                }
            },
            BindEventWithService(testEvent1, eventC))
        c.Assert(err, IsNil)
        c.Assert(s.app.serviceCount(), Equals, 1)

        // event check
        s.app.BroadcastEvent(Event{Name:testEvent1, Payload:testValue1})
        c.Check(<-eventLatch, Equals, testValue1)

        // stop service
        err = s.app.StopServices()
        c.Assert(err, IsNil)

        // cleanup check
        time.Sleep(time.Second)
        c.Assert(s.app.serviceCount(), Equals, 0)
    }
    close(eventLatch)
}
