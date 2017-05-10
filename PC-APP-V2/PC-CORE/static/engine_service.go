package main

import (
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

/*** SERVICE SECTION ***/
const (
    stateStopped    = iota
    stateStarted
)

type Service interface {
    Serve() error
}

type serviceFunc func() error

func (s serviceFunc) Serve() error {
    return s()
}

// Event is a special service event that can be generated by various goroutines in the app
type Event struct {
    Name       string
    Payload    interface{}
}

type Waiter struct {
    eventC     chan Event
    cancelC    chan struct{}
}

type ServiceSupervisor interface {
    RegisterService(srv Service)
    WaitForEvent(name string, eventC chan Event, cancelC chan struct{})
    IsStopped() bool
    StopChannel() <- chan struct{}
    StartServices() error
    StopServices() error
}

type srvSupervisor struct {
    state int
    sync.Mutex

    serviceWG       *sync.WaitGroup
    services        []*Service

    eventsC         chan Event
    events          map[string]Event
    eventWaiters    map[string][]*Waiter

    stoppedC        chan struct{}
}

// ServiceCount returns the number of registered and actively running services
func (s *srvSupervisor) serviceCount() int {
    s.Lock()
    defer s.Unlock()
    return len(s.services)
}

func (s *srvSupervisor) getWaiters(name string) []*Waiter {
    s.Lock()
    defer s.Unlock()

    waiters := s.eventWaiters[name]
    out := make([]*Waiter, len(waiters))
    for i := range waiters {
        out[i] = waiters[i]
    }
    return out
}

func (s *srvSupervisor) notifyWaiter(w *Waiter, evt Event) {
    select {
        case w.eventC <- evt:
        case <-w.cancelC:
    }
}

func (s *srvSupervisor) fanOut() {
    defer s.serviceWG.Done()

    for {
        select {
            case <-s.stoppedC:
                log.Debugf("[SUPERVISOR] fanOut should stop")
                return
            case event := <-s.eventsC:
                waiters := s.getWaiters(event.Name)
                for _, waiter := range waiters {
                    s.notifyWaiter(waiter, event)
                }
        }
    }
}

func (s *srvSupervisor) serve(service *Service) {
    // this func will be called _after_ a service stops running:
    removeService := func(srv *Service) {
        s.Lock()
        defer s.Unlock()
        for i, el := range s.services {
            if el == srv {
                s.services = append(s.services[:i], s.services[i+1:]...)
                break
            }
        }
        log.Debugf("[SUPERVISOR] Service %v is done (%v)", *srv, len(s.services))
    }

    s.serviceWG.Add(1)
    go func(srv *Service) {
        defer s.serviceWG.Done()
        defer removeService(srv)

        log.Debugf("[SUPERVISOR] Service %v started (%v)", *srv, s.serviceCount())
        err := (*srv).Serve()
        if err != nil {
            log.Debug(errors.WithStack(err))
        }
    }(service)
}

func (s *srvSupervisor) RegisterService(srv Service) {
    s.Lock()
    defer s.Unlock()
    s.services = append(s.services, &srv)

    log.Debugf("[SUPERVISOR] Service %v added (%v)", srv, len(s.services))

    if s.state == stateStarted {
        s.serve(&srv)
    }
}

func (s *srvSupervisor) RegisterServiceFunc(fn serviceFunc) {
    s.RegisterService(fn)
}

func (s *srvSupervisor) BroadcastEvent(event Event) {
    s.Lock()
    defer s.Unlock()

    s.events[event.Name] = event
    log.Debugf("[SUPERVISOR] BroadcastEvent: %v", &event)

    go func() {
        s.eventsC <- event
    }()
}

func (s *srvSupervisor) WaitForEvent(name string, eventC chan Event, cancelC chan struct{}) {
    s.Lock()
    defer s.Unlock()

    waiter := &Waiter{eventC: eventC, cancelC: cancelC}
    event, ok := s.events[name]
    if ok {
        go s.notifyWaiter(waiter, event)
        return
    }
    s.eventWaiters[name] = append(s.eventWaiters[name], waiter)
}

func (s *srvSupervisor) IsStopped() bool {
    s.Lock()
    defer s.Unlock()

    if s.state == stateStopped {
        return true
    } else {
        return false
    }
}

func (s *srvSupervisor) StopChannel() <- chan struct{} {
    return s.stoppedC
}

func (s *srvSupervisor) StartServices() error {
    s.Lock()
    defer s.Unlock()

    if s.state == stateStarted {
        return nil
    }
    s.state = stateStarted

    s.serviceWG.Add(1)
    go s.fanOut()

    if len(s.services) == 0 {
        log.Warning("[SUPERVISOR] Start() :: nothing to run")
        return nil
    }

    for _, srv := range s.services {
        s.serve(srv)
    }

    return nil
}

func (s *srvSupervisor) StopServices() error {
    var (
        waiters []*Waiter
        w *Waiter
    )

    // this double locking to to prevent ServiceFunc from deadlocked, but enable other variables to be reset
    s.Lock()
    if s.state == stateStopped {
        defer s.Unlock()
        return nil
    }
    s.state = stateStopped
    s.Unlock()

    log.Warning("[SUPERVISOR] stopping services...")
    // we broadcast stopping and wait for all goroutines closed with event channels intact to give grace period
    close(s.stoppedC)
    s.serviceWG.Wait()

    // this double locking to to prevent ServiceFunc from deadlocked, but enable other variables to be reset
    s.Lock()
    defer s.Unlock()

    // we close event channel after stopping all go routines and fanOut function that we can safely close
    close(s.eventsC)
    for _, waiters = range s.eventWaiters {
        for _, w = range waiters {
            close(w.eventC)
            close(w.cancelC)
        }
    }

    // reset
    s.serviceWG    = &sync.WaitGroup{}
    s.eventsC      = make(chan Event, 100)
    s.events       = map[string]Event{}
    s.eventWaiters = make(map[string][]*Waiter)
    s.stoppedC     = make(chan struct{})

    return nil
}

func newServiceSupervisor() *srvSupervisor {
    return &srvSupervisor{
        state:           stateStopped,
        serviceWG:       &sync.WaitGroup{},
        services:        []*Service{},

        eventsC:         make(chan Event, 100),
        events:          map[string]Event{},
        eventWaiters:    make(map[string][]*Waiter),

        stoppedC:        make(chan struct{}),
    }
}