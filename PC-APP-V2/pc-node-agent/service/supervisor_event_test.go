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
        exitLatch = make(chan string)
        eventLatch = make(chan string)

        eventC1 = make(chan Event)
        eventC2 = make(chan Event)
        eventC3 = make(chan Event)
    )
    err := s.app.Start()
    c.Assert(err, IsNil)

    err = s.app.RegisterServiceWithFuncs(
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
                        return nil
                    }
                    default:
                }
            }
        },
        func(_ func(interface{})) error {
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

    err = s.app.Stop()
    c.Assert(err, IsNil)
    c.Check(<-exitLatch, Equals, exitValue)

    // it takes abit to
    time.Sleep(time.Second)
    c.Assert(s.app.serviceCount(), Equals, 0)
    close(exitLatch)
    close(eventLatch)
}