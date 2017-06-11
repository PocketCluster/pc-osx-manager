package service

import (
    "time"

    . "gopkg.in/check.v1"
)

const (
    testEvent1 = "test_event1"
    testEvent2 = "test_event2"
    testEvent3 = "test_event3"
    testEvent4 = "test_event4"
)

const (
    testValue1 = "test_value1"
    testValue2 = "test_value2"
    testValue3 = "test_value3"
    testValue4 = "test_value4"
)

func (s *SupervisorSuite) Test_UnamedService_Receive_Event(c *C) {
    var(
        exitLatch = make(chan bool)
        exitChecker = ""

        eventLatch = make(chan bool)
        eventChecker = ""
        eventC1 = make(chan Event)
        eventC2 = make(chan Event)
        eventC3 = make(chan Event)
    )
    err := s.app.Start()
    c.Assert(err, IsNil)

    s.app.WaitForEvent(testEvent1, eventC1)
    s.app.WaitForEvent(testEvent2, eventC2)
    s.app.WaitForEvent(testEvent3, eventC3)

    err = s.app.RegisterServiceWithFuncs(
        func() error {
            for {
                select {
                    case e := <-eventC1: {
                        eventChecker = e.Payload.(string)
                        eventLatch <- true
                    }
                    case e := <-eventC1: {
                        eventChecker = e.Payload.(string)
                        eventLatch <- true
                    }
                    case e := <-eventC1: {
                        eventChecker = e.Payload.(string)
                        eventLatch <- true
                    }
                    case <- s.app.StopChannel(): {
                        return nil
                    }
                    default:
                }
            }
        },
        func(_ func(interface{})) error {
            exitChecker = exitValue
            exitLatch <- true
            return nil
        },
        BindEventWithService(eventC1, eventC2, eventC3),
    )
    c.Assert(err, IsNil)
    c.Assert(s.app.serviceCount(), Equals, 1)

    s.app.BroadcastEvent(Event{Name:testEvent1, Payload:testValue1})
    <-eventLatch
    c.Assert(eventChecker, Equals, testValue1)

    s.app.BroadcastEvent(Event{Name:testEvent1, Payload:testValue1})
    <-eventLatch
    c.Assert(eventChecker, Equals, testValue1)

    s.app.BroadcastEvent(Event{Name:testEvent1, Payload:testValue1})
    <-eventLatch
    c.Assert(eventChecker, Equals, testValue1)

    err = s.app.Stop()
    c.Assert(err, IsNil)
    <-exitLatch

    // it takes abit to
    time.Sleep(time.Second)
    c.Check(exitChecker, Equals, exitValue)
    c.Assert(s.app.serviceCount(), Equals, 0)
    close(exitLatch)
}