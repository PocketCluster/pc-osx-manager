package main

import (
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

const (
    stateStopped = iota
    stateStarted
)

type PocketApplication struct {
    state           int
    sync.Mutex

    serviceWG       *sync.WaitGroup
    services        []*Service

    events          map[string]Event
    eventsC         chan Event
    eventWaiters    map[string][]*waiter

    stoppedC        chan struct{}
}

// NewSupervisor returns new instance of initialized supervisor
func NewPocketApplication() *PocketApplication {
    srv := &PocketApplication{
        state:           stateStopped,
        services:        []*Service{},
        serviceWG:       &sync.WaitGroup{},

        events:          map[string]Event{},
        eventsC:         make(chan Event, 100),
        eventWaiters:    make(map[string][]*waiter),

        stoppedC:        make(chan struct{}),
    }

    // for fanOut function
    srv.serviceWG.Add(1)
    go srv.fanOut()
    return srv
}

func (p *PocketApplication) getWaiters(name string) []*waiter {
    p.Lock()
    defer p.Unlock()

    waiters := p.eventWaiters[name]
    out := make([]*waiter, len(waiters))
    for i := range waiters {
        out[i] = waiters[i]
    }
    return out
}

func (p *PocketApplication) notifyWaiter(w *waiter, event Event) {
    select {
        case w.eventC <- event:
        case <-w.cancelC:
    }
}

func (p *PocketApplication) fanOut() {
    defer p.serviceWG.Done()

    for {
        select {
            case <-p.stoppedC:
                return
            case event := <-p.eventsC:
                waiters := p.getWaiters(event.Name)
                for _, waiter := range waiters {
                    p.notifyWaiter(waiter, event)
                }
        }
    }
}

func (p *PocketApplication) serve(service *Service) {
    // this func will be called _after_ a service stops running:
    removeService := func(srv *Service) {
        p.Lock()
        defer p.Unlock()
        for i, el := range p.services {
            if el == srv {
                p.services = append(p.services[:i], p.services[i+1:]...)
                break
            }
        }
        log.Debugf("[APPLICATION] Service %v is done (%v)", *srv, len(p.services))
    }

    p.serviceWG.Add(1)
    go func(srv *Service) {
        defer p.serviceWG.Done()
        defer removeService(srv)

        log.Debugf("[APPLICATION] Service %v started (%v)", *srv, p.ServiceCount())
        err := (*srv).Serve()
        if err != nil {
            log.Debug(errors.WithStack(err))
        }
    }(service)
}


func (p *PocketApplication) Register(srv Service) {
    p.Lock()
    defer p.Unlock()
    p.services = append(p.services, &srv)

    log.Debugf("[APPLICATION] Service %v added (%v)", srv, len(p.services))

    if p.state == stateStarted {
        p.serve(&srv)
    }
}

func (p *PocketApplication) RegisterFunc(fn ServiceFunc) {
    p.Register(fn)
}

func (p *PocketApplication) BroadcastEvent(event Event) {
    p.Lock()
    defer p.Unlock()

    p.events[event.Name] = event
//    log.Debugf("[APPLICATION] BroadcastEvent: %v", &event)

    go func() {
        p.eventsC <- event
    }()
}

func (p *PocketApplication) WaitForEvent(name string, eventC chan Event, cancelC chan struct{}) {
    p.Lock()
    defer p.Unlock()

    waiter := &waiter{eventC: eventC, cancelC: cancelC}
    event, ok := p.events[name]
    if ok {
        go p.notifyWaiter(waiter, event)
        return
    }
    p.eventWaiters[name] = append(p.eventWaiters[name], waiter)
}

// ServiceCount returns the number of registered and actively running services
func (p *PocketApplication) ServiceCount() int {
    p.Lock()
    defer p.Unlock()
    return len(p.services)
}

func (p *PocketApplication) IsStopped() bool {
    p.Lock()
    defer p.Unlock()

    if p.state == stateStarted {
        return false
    } else {
        return true
    }
}

func (p *PocketApplication) Start() error {
    p.Lock()
    defer p.Unlock()
    p.state = stateStarted

    if len(p.services) == 0 {
        log.Warning("PocketApplication.Start(): nothing to run")
        return nil
    }

    for _, srv := range p.services {
        p.serve(srv)
    }

    return nil
}

func (s *PocketApplication) Stop() error {
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
    return nil
}

func (p *PocketApplication) Wait() error {
    var (
        waiters []*waiter
        w *waiter
    )
    p.serviceWG.Wait()

    // we close event channel after stopping all go routines and fanOut function that we can safely close
    close(p.eventsC)
    for _, waiters = range p.eventWaiters {
        for _, w = range waiters {
            close(w.eventC)
            close(w.cancelC)
        }
    }

    return nil
}

// onExit allows individual services to register a callback function which will be
// called when Teleport Process is asked to exit. Usually services terminate themselves
// when the callback is called
func (p *PocketApplication) OnExit(callback func(interface{})) {
    go func() {
        select {
            case <- p.stoppedC:
                callback(nil)
        }
    }()
}

// Event is a special service event that can be generated
// by various goroutines in the application
type Event struct {
    Name    string
    Payload interface{}
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
