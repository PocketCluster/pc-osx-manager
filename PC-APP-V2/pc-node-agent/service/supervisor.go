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
    serviceStopped = iota
    serviceRunning
)

var (
    // Service Running error
    ServiceRunning = &supError{"[ERR] service already running"}
)

type supError struct {
    s string
}

func (n *supError) Error() string {
    return n.s
}

// Event is a special service event that can be generated
// by various goroutines in the application
type Event struct {
    Name    string
    Payload interface{}
}

type waiter struct {
    eventC  chan Event
    service Service
}

type Service interface {
    // These indicate that the service is named and should be individually managed to start.
    // Plus, only named service can be rebuild
    Name() string
    IsRunning() bool

    // --- internal methods ---
    isNamedService() bool

    serve() error
    onExit(callback func(interface{})) error

    setRunning()
    setStopped()

    rebuild() error

    getWaiters() []*waiter
}

// ServerOption is a functional option passed to the server
type ServiceOption func(app AppSupervisor, s Service) error

type ServeFunc func() error
func (s ServeFunc) serve() error {
    return s()
}

type OnExitFunc func(callback func(interface{})) error
func (o OnExitFunc) onExit(callback func(interface{})) error {
    return o(callback)
}

// --- private unexpose methods --- //

type srvcFuncs struct {
    sync.Mutex
    state      int

    name       string
    waiters    []*waiter
    ServeFunc
    OnExitFunc
}

func (s *srvcFuncs) Name() string {
    return s.name
}

func (s *srvcFuncs) IsRunning() bool {
    s.Lock()
    defer s.Unlock()

    return (s.state == serviceRunning)
}

func (s *srvcFuncs) isNamedService() bool {
    return (len(s.name) != 0)
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

func (s *srvcFuncs) rebuild() error {
    return nil
}

func (s *srvcFuncs) getWaiters() []*waiter {
    return s.waiters
}

func MakeServiceNamed(name string) ServiceOption {
    return func(_ AppSupervisor, s Service) error {
        srv, ok := s.(*srvcFuncs)
        if ok {
            srv.name = name
            return nil
        }
        if srv == nil {
            return errors.Errorf("[ERR] null service instance to bind event")
        }
        return errors.Errorf("[ERR] invalid type to make cycle async")
    }
}

func BindEventWithService(eventName string, eventC chan Event) ServiceOption {
    return func(app AppSupervisor, s Service) error {
        srv, ok := s.(*srvcFuncs)
        if !ok {
            return errors.Errorf("[ERR] invalid service type to bind event")
        }
        if srv == nil {
            return errors.Errorf("[ERR] null service instance to bind event")
        }
        sup, ok := app.(*appSupervisor)
        if !ok {
            return errors.Errorf("[ERR] invalid supervisor type to bind event")
        }
        if sup == nil {
            return errors.Errorf("[ERR] null supervisor instance to bind event")
        }
        sup.Lock()
        defer sup.Unlock()

        w := &waiter{
            eventC:     eventC,
            service:    srv,
        }

        sup.eventWaiters[eventName] = append(sup.eventWaiters[eventName], w)
        srv.waiters = append(srv.waiters, w)
        return nil
    }
}

// --- AppSupervisor --- //

type AppSupervisor interface {
    Register(srv Service) error
    RegisterServiceWithFuncs(sfn ServeFunc, efn OnExitFunc, options... ServiceOption) error
    BroadcastEvent(event Event)

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

// NewSupervisor returns new instance of initialized supervisor
func NewAppSupervisor() AppSupervisor {
    srv := &appSupervisor{
        state:           stateStopped,
        services:        []Service{},
        serviceWG:       &sync.WaitGroup{},

        eventWaiters:    make(map[string][]*waiter),
        eventsC:         make(chan Event, 100),
        stoppedC:        make(chan struct{}),
    }

    // for fanOut function
    srv.serviceWG.Add(1)
    go srv.fanOut()
    return srv
}

type appSupervisor struct {
    sync.Mutex
    state           int

    serviceWG       *sync.WaitGroup
    services        []Service

    eventWaiters    map[string][]*waiter
    eventsC         chan Event
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

func (p *appSupervisor) runService(service Service) {
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
    removeEvents := func(srv Service) {
        p.Lock()
        defer p.Unlock()
    }

    p.serviceWG.Add(1)
    go func(srv Service, delSrv func(srv Service), delEvent func(srv Service), wg *sync.WaitGroup) {
        defer wg.Done()
        defer srv.setStopped()
        srv.setRunning()

        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] started", srv.Name(), srv)
        err := srv.serve()
        if err != nil {
            log.Debug(errors.WithStack(err))
        }
        delEvent(srv)
        err = srv.onExit(nil)
        if err != nil {
            log.Debug(errors.WithStack(err))
        }
        if !srv.isNamedService() {
            delSrv(srv)
        }
        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] exited", srv.Name(), srv)
    }(service, removeService, removeEvents, p.serviceWG)
}

func (p *appSupervisor) Register(srv Service) error {
    p.Lock()
    defer p.Unlock()

    // when a service is named, check if there is a service with the same name
    if srv.isNamedService() {
        for _, es := range p.services {
            if srv.Name() == es.Name() {
                return errors.Errorf("[ERR] a service with the same name '%s' exists", srv.Name())
            }
        }
    }
    p.services = append(p.services, srv)

    log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] added", srv.Name(), srv)

    if p.state == stateStarted && !srv.isNamedService() {
        p.runService(srv)
    }
    return nil
}

func (p *appSupervisor) RegisterServiceWithFuncs(sfn ServeFunc, efn OnExitFunc, options... ServiceOption) error {
    srv := &srvcFuncs {
        state:         serviceStopped,
        waiters:       []*waiter{},
        ServeFunc:     sfn,
        OnExitFunc:    efn,
    }

    for _, opt := range options {
        if err := opt(p, srv); err != nil {
            return errors.WithStack(err)
        }
    }

    return p.Register(srv)
}

func (p *appSupervisor) BroadcastEvent(event Event) {
    p.Lock()
    defer p.Unlock()

    go func() {
        p.eventsC <- event
    }()
}

func (p *appSupervisor) WaitForEvent(name string, eventC chan Event) {
    p.Lock()
    defer p.Unlock()

    p.eventWaiters[name] = append(p.eventWaiters[name], &waiter{eventC: eventC})
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
        if !srv.isNamedService() {
            p.runService(srv)
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

    if len(name) == 0 {
        return errors.Errorf("[ERR] cannot find a service without name")
    }

    for _, srv := range p.services {
        if srv.isNamedService() && srv.Name() == name {
            if !srv.IsRunning() {
                p.runService(srv)
                return nil
            } else {
                return ServiceRunning
            }
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
