package service

import (
    "time"

    . "gopkg.in/check.v1"
)

const (
    testEvent1 = "test_event1"
    testEvent2 = "test_event2"
    testEvent3 = "test_event3"
)

const (
    testValue1 = "test_value1"
    testValue2 = "test_value2"
    testValue3 = "test_value3"
)

func (s *SupervisorSuite) Test_UnnamedService_Receive_MultiEvent(c *C) {
    var(
        exitLatch  = make(chan string)
        eventLatch = make(chan string)

        eventC1    = make(chan Event)
        eventC2    = make(chan Event)
        eventC3    = make(chan Event)
    )
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.RegisterServiceWithFuncs(
        func() (interface{}, error) {
            for {
                select {
                    case e := <-eventC1: {
                        eventLatch <- e.Payload.(string)
                    }
                    case e := <-eventC2: {
                        eventLatch <- e.Payload.(string)
                    }
                    case e := <-eventC3: {
                        eventLatch <- e.Payload.(string)
                    }
                    case <- s.app.StopChannel(): {
                        return nil, nil
                    }
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch <- exitValue
            return nil
        },
        BindEventWithService(testEvent1, eventC1),
        BindEventWithService(testEvent2, eventC2),
        BindEventWithService(testEvent3, eventC3),
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    s.app.BroadcastEvent(Event{Name:testEvent1, Payload:testValue1})
    c.Assert(<-eventLatch, Equals, testValue1)

    s.app.BroadcastEvent(Event{Name:testEvent2, Payload:testValue2})
    c.Assert(<-eventLatch, Equals, testValue2)

    s.app.BroadcastEvent(Event{Name:testEvent3, Payload:testValue3})
    c.Assert(<-eventLatch, Equals, testValue3)

    err = s.app.StopServices()
    c.Assert(err, IsNil)
    c.Check(<-exitLatch, Equals, exitValue)

    // it takes abit to
    time.Sleep(time.Second)
    c.Assert(s.app.serviceCount(), Equals, 0)

    // check if water queue is empty
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters), Equals, 3)
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters[testEvent1]), Equals, 0)
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters[testEvent2]), Equals, 0)
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters[testEvent3]), Equals, 0)

    // close everything
    close(exitLatch)
    close(eventLatch)
}

func (s *SupervisorSuite) Test_Multiple_UnnamedService_Receive_Event(c *C) {
    var(
        exitLatch  = make(chan string)
        eventLatch = make(chan string)

        eventC1    = make(chan Event)
        eventC2    = make(chan Event)
        eventC3    = make(chan Event)
    )
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.RegisterServiceWithFuncs(
        func() (interface{}, error) {
            for {
                select {
                    case e := <-eventC1: {
                        eventLatch <- e.Payload.(string)
                    }
                    case <- s.app.StopChannel(): {
                        return nil, nil
                    }
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch <- exitValue
            return nil
        },
        BindEventWithService(testEvent1, eventC1),
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.RegisterServiceWithFuncs(
        func() (interface{}, error) {
            for {
                select {
                    case e := <-eventC2: {
                        eventLatch <- e.Payload.(string)
                    }
                    case <- s.app.StopChannel(): {
                        return nil, nil
                    }
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch <- exitValue
            return nil
        },
        BindEventWithService(testEvent1, eventC2),
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 2)

    err = s.app.RegisterServiceWithFuncs(
        func() (interface{}, error) {
            for {
                select {
                    case e := <-eventC3: {
                        eventLatch <- e.Payload.(string)
                    }
                    case <- s.app.StopChannel(): {
                        return nil, nil
                    }
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch <- exitValue
            return nil
        },
        BindEventWithService(testEvent1, eventC3),
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 3)

    s.app.BroadcastEvent(Event{Name:testEvent1, Payload:testValue1})
    c.Assert(<-eventLatch, Equals, testValue1)
    c.Assert(<-eventLatch, Equals, testValue1)
    c.Assert(<-eventLatch, Equals, testValue1)

    err = s.app.StopServices()
    c.Assert(err, IsNil)
    c.Check(<-exitLatch, Equals, exitValue)
    c.Check(<-exitLatch, Equals, exitValue)
    c.Check(<-exitLatch, Equals, exitValue)

    // it takes abit to
    time.Sleep(time.Second)
    c.Assert(s.app.serviceCount(), Equals, 0)

    // check if water queue is empty
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters), Equals, 1)
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters[testEvent1]), Equals, 0)

    // close everything
    close(exitLatch)
    close(eventLatch)
}

func (s *SupervisorSuite) Test_Restart_Multiple_NamedService_Receive_Event(c *C) {
    var(
        eventLatch = make(chan string)
        exitLatch1 = make(chan string)
        exitLatch2 = make(chan string)
        exitLatch3 = make(chan string)
        controlExitLatch = make(chan string)

        eventC1    = make(chan string)
        eventC2    = make(chan string)
        eventC3    = make(chan string)
    )
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.RegisterNamedServiceWithFuncs(
        testService1,
        func() (interface{}, error) {
            for {
                select {
                    case e := <-eventC1: {
                        eventLatch <- e
                    }
                    case <- exitLatch1: {
                        return nil, nil
                    }
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

    err = s.app.RegisterNamedServiceWithFuncs(
        testService2,
        func() (interface{}, error) {
            for {
                select {
                    case e := <-eventC2: {
                        eventLatch <- e
                    }
                    case <- exitLatch2: {
                        return nil, nil
                    }
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

    err = s.app.RegisterNamedServiceWithFuncs(
        testService3,
        func() (interface{}, error) {
            for {
                select {
                    case e := <-eventC3: {
                        eventLatch <- e
                    }
                    case <- exitLatch3: {
                        return nil, nil
                    }
                    default:
                }
            }
        },
        func(_ interface{}, _ error) error {
            exitLatch3 <- exitValue
            return nil
        })
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 3)

    err = s.app.RegisterServiceWithFuncs(
        func() (interface{}, error) {
            for i := 0; i < 5; i++ {
                // start services within a service
                err = s.app.RunNamedService(testService1)
                c.Assert(err, IsNil)
                err = s.app.RunNamedService(testService2)
                c.Assert(err, IsNil)
                err = s.app.RunNamedService(testService3)
                c.Assert(err, IsNil)

                // check if all three receives event correct
                eventC1 <- testValue1
                eventC2 <- testValue1
                eventC3 <- testValue1
                c.Assert(<-eventLatch, Equals, testValue1)
                c.Assert(<-eventLatch, Equals, testValue1)
                c.Assert(<-eventLatch, Equals, testValue1)

                // check
                exitLatch1 <- ""
                c.Check(<-exitLatch1, Equals, exitValue)
                exitLatch2 <- ""
                c.Check(<-exitLatch2, Equals, exitValue)
                exitLatch3 <- ""
                c.Check(<-exitLatch3, Equals, exitValue)

                // it takes abit to restart
                time.Sleep(time.Second)
            }
            return nil, nil
        },
        func(_ interface{}, _ error) error {
            controlExitLatch <- exitValue
            return nil
        },
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 4)

    c.Assert(<-controlExitLatch, Equals, exitValue)
    err = s.app.StopServices()
    c.Assert(err, IsNil)

    c.Assert(s.app.serviceCount(), Equals, 3)
    // check if water queue is empty
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters), Equals, 0)

    // close everything
    close(eventLatch)
    close(exitLatch1)
    close(exitLatch2)
    close(exitLatch3)
    close(eventC1)
    close(eventC2)
    close(eventC3)
    close(controlExitLatch)
}