package main

import (
    "fmt"
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

const (
    stateCreated    = iota
    stateStarted
)

type PocketSupervisor struct {
    state int
    sync.Mutex
    wg              *sync.WaitGroup

    services        []*Service
    errors          []error
    events          map[string]Event
    eventsC         chan Event
    eventWaiters    map[string][]*waiter
    closer          *CloseBroadcaster
}

// NewSupervisor returns new instance of initialized supervisor
func NewPocketSupervisor() *PocketSupervisor {
    srv := &PocketSupervisor{
        services:        []*Service{},
        wg:              &sync.WaitGroup{},

        events:          map[string]Event{},
        eventsC:         make(chan Event, 100),
        eventWaiters:    make(map[string][]*waiter),
        closer:          &CloseBroadcaster{
            C: make(chan struct{}),
        },
    }
    go srv.fanOut()
    return srv
}

func (s *PocketSupervisor) getWaiters(name string) []*waiter {
    s.Lock()
    defer s.Unlock()

    waiters := s.eventWaiters[name]
    out := make([]*waiter, len(waiters))
    for i := range waiters {
        out[i] = waiters[i]
    }
    return out
}

func (s *PocketSupervisor) notifyWaiter(w *waiter, event Event) {
    select {
    case w.eventC <- event:
    case <-w.cancelC:
    }
}

func (s *PocketSupervisor) fanOut() {
    for {
        select {
        case event := <-s.eventsC:
            waiters := s.getWaiters(event.Name)
            for _, waiter := range waiters {
                go s.notifyWaiter(waiter, event)
            }
        case <-s.closer.C:
            return
        }
    }
}
func (p *PocketSupervisor) serve(srv *Service) {
    // this func will be called _after_ a service stops running:
    removeService := func() {
        p.Lock()
        defer p.Unlock()
        for i, el := range p.services {
            if el == srv {
                p.services = append(p.services[:i], p.services[i+1:]...)
                break
            }
        }
        log.Debugf("[SUPERVISOR] Service %v is done (%v)", *srv, len(p.services))
    }

    p.wg.Add(1)
    go func() {
        defer p.wg.Done()
        defer removeService()

        log.Debugf("[SUPERVISOR] Service %v started (%v)", *srv, p.ServiceCount())
        err := (*srv).Serve()
        if err != nil {
            errors.WithStack(err)
        }
    }()
}

func (p *PocketSupervisor) Register(srv Service) {
    p.Lock()
    defer p.Unlock()
    p.services = append(p.services, &srv)

    log.Debugf("[SUPERVISOR] Service %v added (%v)", srv, len(p.services))

    if p.state == stateStarted {
        p.serve(&srv)
    }
}

// ServiceCount returns the number of registered and actively running services
func (p *PocketSupervisor) ServiceCount() int {
    p.Lock()
    defer p.Unlock()
    return len(p.services)
}

func (p *PocketSupervisor) RegisterFunc(fn ServiceFunc) {
    p.Register(fn)
}

func (s *PocketSupervisor) Start() error {
    s.Lock()
    defer s.Unlock()
    s.state = stateStarted

    if len(s.services) == 0 {
        log.Warning("supervisor.Start(): nothing to run")
        return nil
    }

    for _, srv := range s.services {
        s.serve(srv)
    }

    return nil
}

func (s *PocketSupervisor) Wait() error {
    defer s.closer.Close()
    s.wg.Wait()
    return nil
}

func (s *PocketSupervisor) Run() error {
    if err := s.Start(); err != nil {
        return errors.WithStack(err)
    }
    return s.Wait()
}

func (s *PocketSupervisor) BroadcastEvent(event Event) {
    s.Lock()
    defer s.Unlock()
    s.events[event.Name] = event
    log.Debugf("BroadcastEvent: %v", &event)

    go func() {
        s.eventsC <- event
    }()
}

func (s *PocketSupervisor) WaitForEvent(name string, eventC chan Event, cancelC chan struct{}) {
    s.Lock()
    defer s.Unlock()

    waiter := &waiter{eventC: eventC, cancelC: cancelC}
    event, ok := s.events[name]
    if ok {
        go s.notifyWaiter(waiter, event)
        return
    }
    s.eventWaiters[name] = append(s.eventWaiters[name], waiter)
}


// Event is a special service event that can be generated
// by various goroutines in the supervisor
type Event struct {
    Name    string
    Payload interface{}
}

func (e *Event) String() string {
    return fmt.Sprintf("event(%v)", e.Name)
}

type waiter struct {
    eventC  chan Event
    cancelC chan struct{}
}

type Service interface {
    Serve() error
}

type ServiceFunc func() error

func (s ServiceFunc) Serve() error {
    return s()
}

// CloseBroadcaster is a helper struct
// that implements io.Closer and uses channel
// to broadcast it's closed state once called
type CloseBroadcaster struct {
    sync.Once
    C chan struct{}
}

// Close closes channel (once) to start broadcasting it's closed state
func (b *CloseBroadcaster) Close() error {
    b.Do(func() {
        close(b.C)
    })
    return nil
}
