package service

import (
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

func (s *SupervisorSuite) Test_Service_Receive_MultiEvent(c *C) {
    var(
        exitChecker string = ""
        eventLatch = make(chan string)

        eventC1    = make(chan Event)
        eventC2    = make(chan Event)
        eventC3    = make(chan Event)
    )
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.RegisterServiceWithFuncs(
        testService1,
        func() error {
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
                        exitChecker = exitValue
                        return nil
                    }
                    default:
                }
            }
        },
        BindEventWithService(testEvent1, eventC1),
        BindEventWithService(testEvent2, eventC2),
        BindEventWithService(testEvent3, eventC3))
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
    c.Assert(exitChecker, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 0)

    // check if waiter queue is empty
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters), Equals, 3)
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters[testEvent1]), Equals, 0)
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters[testEvent2]), Equals, 0)
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters[testEvent3]), Equals, 0)
    close(eventLatch)
}

func (s *SupervisorSuite) Test_Multiple_Service_Receive_Event(c *C) {
    var(
        exitChecker1 string = ""
        exitChecker2 string = ""
        exitChecker3 string = ""

        eventLatch = make(chan string)

        eventC1    = make(chan Event)
        eventC2    = make(chan Event)
        eventC3    = make(chan Event)
    )
    err := s.app.StartServices()
    c.Assert(err, IsNil)

    err = s.app.RegisterServiceWithFuncs(
        testService1,
        func() error {
            for {
                select {
                    case e := <-eventC1: {
                        eventLatch <- e.Payload.(string)
                    }
                    case <- s.app.StopChannel(): {
                        exitChecker1 = exitValue
                        return nil
                    }
                    default:
                }
            }
        },
        BindEventWithService(testEvent1, eventC1))
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    err = s.app.RegisterServiceWithFuncs(
        testService2,
        func() error {
            for {
                select {
                    case e := <-eventC2: {
                        eventLatch <- e.Payload.(string)
                    }
                    case <- s.app.StopChannel(): {
                        exitChecker2 = exitValue
                        return nil
                    }
                    default:
                }
            }
        },
        BindEventWithService(testEvent1, eventC2))
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 2)

    err = s.app.RegisterServiceWithFuncs(
        testService3,
        func() error {
            for {
                select {
                    case e := <-eventC3: {
                        eventLatch <- e.Payload.(string)
                    }
                    case <- s.app.StopChannel(): {
                        exitChecker3 = exitValue
                        return nil
                    }
                    default:
                }
            }
        },
        BindEventWithService(testEvent1, eventC3))
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 3)

    s.app.BroadcastEvent(Event{Name:testEvent1, Payload:testValue1})
    c.Assert(<-eventLatch, Equals, testValue1)
    c.Assert(<-eventLatch, Equals, testValue1)
    c.Assert(<-eventLatch, Equals, testValue1)

    err = s.app.StopServices()
    c.Assert(err, IsNil)
    c.Check(exitChecker1, Equals, exitValue)
    c.Check(exitChecker2, Equals, exitValue)
    c.Check(exitChecker3, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 0)

    // check if water queue is empty
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters), Equals, 1)
    c.Assert(len(s.app.(*srvcSupervisor).eventWaiters[testEvent1]), Equals, 0)

    // close everything
    close(eventLatch)
}
