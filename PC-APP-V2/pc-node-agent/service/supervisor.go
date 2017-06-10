package service

import (
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)

const (
    stateStopped = iota
    stateStarted
)

const (
    serviceRunning = iota
    serviceStopped
)

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
    // These indicate that the service is named and should be individually managed to start.
    // Plus, only named service can be rebuild
    Name() string
    Rebuild() error
    IsNamedCycle() bool

    Serve() error
    OnExit(callback func(interface{})) error

    IsRunning() bool
    // this is for internal handling
    setRunning()
    setStopped()
}

// ServerOption is a functional option passed to the server
type ServiceOption func(s Service) error

type ServeFunc func() error
func (s ServeFunc) Serve() error {
    return s()
}

type OnExitFunc func(callback func(interface{})) error
func (o OnExitFunc) OnExit(callback func(interface{})) error {
    return o(callback)
}

type AppSupervisor interface {
    Register(srv Service) error
    RegisterServiceWithFuncs(sfn ServeFunc, efn OnExitFunc, options...ServiceOption) error

    BroadcastEvent(event Event)
    WaitForEvent(name string, eventC chan Event, cancelC chan struct{})

    IsStopped() bool
    StopChannel() <- chan struct{}

    Start() error
    RunNamedService(name string) error
    Stop() error
    Wait() error
    OnExit(callback func(interface{}))

    // internal
    serviceCount() int
}

// --- private unexpose methods --- //

type srvcFuncs struct {
    sync.Mutex
    state      int

    name       string
    ServeFunc
    OnExitFunc
}

func (s *srvcFuncs) Name() string {
    return s.name
}

func (s *srvcFuncs) IsNamedCycle() bool {
    return (len(s.name) != 0)
}

func (s *srvcFuncs) Rebuild() error {
    return nil
}

func (s *srvcFuncs) IsRunning() bool {
    s.Lock()
    defer s.Unlock()

    return (s.state == serviceRunning)
}

func (s *srvcFuncs) setRunning() {
    s.Lock()
    defer s.Unlock()

    s.state = serviceRunning
}

func (s *srvcFuncs) setStopped() {
    s.Lock()
    defer s.Unlock()

    s.state = serviceStopped
}

func MakeServiceNamed(name string) ServiceOption {
    return func(s Service) error {
        srv, ok := s.(*srvcFuncs)
        if ok {
            srv.name = name
            return nil
        }
        return errors.Errorf("[ERR] invalid type to make cycle async")
    }
}

// NewSupervisor returns new instance of initialized supervisor
func NewAppSupervisor() AppSupervisor {
    srv := &appSupervisor{
        state:           stateStopped,
        services:        []Service{},
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

type appSupervisor struct {
    state           int
    sync.Mutex

    serviceWG       *sync.WaitGroup
    services        []Service

    events          map[string]Event
    eventsC         chan Event
    eventWaiters    map[string][]*waiter

    stoppedC        chan struct{}
}

func (p *appSupervisor) getWaiters(name string) []*waiter {
    p.Lock()
    defer p.Unlock()

    waiters := p.eventWaiters[name]
    out := make([]*waiter, len(waiters))
    for i := range waiters {
        out[i] = waiters[i]
    }
    return out
}

func (p *appSupervisor) notifyWaiter(w *waiter, event Event) {
    select {
        case w.eventC <- event:
        case <-w.cancelC:
    }
}

func (p *appSupervisor) fanOut() {
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

func (p *appSupervisor) serve(service Service) {
    // this func will be called _after_ a service stops running:
    removeService := func(srv Service) {
        p.Lock()
        defer p.Unlock()
        for i, el := range p.services {
            // TODO : MAKE 100% SURE THIS COMPARISON WORKS PROPERLY
            if el == srv {
                p.services = append(p.services[:i], p.services[i+1:]...)
                break
            }
        }
        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] removed", srv.Name(), srv)
    }

    p.serviceWG.Add(1)
    go func(srv Service, delSrv func(srv Service), wg *sync.WaitGroup) {
        defer wg.Done()

        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] started", srv.Name(), srv)
        srv.setRunning()
        err := srv.Serve()
        if err != nil {
            log.Debug(errors.WithStack(err))
        }
        err = srv.OnExit(nil)
        if err != nil {
            log.Debug(errors.WithStack(err))
        }
        srv.setStopped()
        if !srv.IsNamedCycle() {
            delSrv(srv)
        }
        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] exited", srv.Name(), srv)
    }(service, removeService, p.serviceWG)
}

func (p *appSupervisor) Register(srv Service) error {
    p.Lock()
    defer p.Unlock()

    // when a service is named, check if there is a service with the same name
    if len(srv.Name()) != 0 {
        for _, es := range p.services {
            if srv.Name() == es.Name() {
                return errors.Errorf("[ERR] a service with the same name '%s' exists", srv.Name())
            }
        }
    }
    p.services = append(p.services, srv)

    log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] added", srv.Name(), srv)

    if p.state == stateStarted && !srv.IsNamedCycle() {
        p.serve(srv)
    }
    return nil
}

func (p *appSupervisor) RegisterServiceWithFuncs(sfn ServeFunc, efn OnExitFunc, options...ServiceOption) error {
    srv := &srvcFuncs {
        ServeFunc:     sfn,
        OnExitFunc:    efn,
    }

    for _, opt := range options {
        if err := opt(srv); err != nil {
            return errors.WithStack(err)
        }
    }

    return p.Register(srv)
}

func (p *appSupervisor) BroadcastEvent(event Event) {
    p.Lock()
    defer p.Unlock()

    p.events[event.Name] = event

    go func() {
        p.eventsC <- event
    }()
}

func (p *appSupervisor) WaitForEvent(name string, eventC chan Event, cancelC chan struct{}) {
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
func (p *appSupervisor) serviceCount() int {
    p.Lock()
    defer p.Unlock()
    return len(p.services)
}

func (p *appSupervisor) IsStopped() bool {
    p.Lock()
    defer p.Unlock()

    if p.state == stateStarted {
        return false
    } else {
        return true
    }
}

func (p *appSupervisor) StopChannel() <- chan struct{} {
    return p.stoppedC
}

func (p *appSupervisor) Start() error {
    p.Lock()
    if p.state == stateStarted {
        defer p.Unlock()
        return nil
    }
    defer p.Unlock()
    p.state = stateStarted

    if len(p.services) == 0 {
        log.Debugf("\n\n[SUPERVISOR-SERVICE] starting services...")
        return nil
    }

    for _, srv := range p.services {
        if !srv.IsNamedCycle() {
            p.serve(srv)
        }
    }

    return nil
}

func (p *appSupervisor) RunNamedService(name string) error {
    p.Lock()
    if p.state != stateStarted {
        defer p.Unlock()
        return nil
    }
    defer p.Unlock()

    for _, srv := range p.services {
        if srv.IsNamedCycle() && srv.Name() == name {
            p.serve(srv)
            return nil
        }
    }

    return errors.Errorf("[ERR] cannot find a service named %s", name)
}


func (s *appSupervisor) Stop() error {
    // this double locking to to prevent ServiceFunc from deadlocked, but enable other variables to be reset
    s.Lock()
    if s.state == stateStopped {
        defer s.Unlock()
        return nil
    }
    defer s.Unlock()
    s.state = stateStopped

    log.Debugf("[SUPERVISOR-SERVICE] stopping all services...\n")
    // we broadcast stopping and wait for all goroutines closed with event channels intact to give grace period
    close(s.stoppedC)
    return nil
}

func (p *appSupervisor) Wait() error {
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
func (p *appSupervisor) OnExit(callback func(interface{})) {
    go func() {
        select {
            case <- p.stoppedC:
                callback(nil)
        }
    }()
}
