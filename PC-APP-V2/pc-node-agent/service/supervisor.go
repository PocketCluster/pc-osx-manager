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
    name    string
    eventC  chan Event
}

// --- --- --- --- --- --- --- --- --- --- --- --- --- service unit --- --- --- --- --- --- --- --- --- --- --- --- - //

type Service interface {
    // These indicate that the service is named and should be individually managed to start.
    // Plus, only named service can be rebuild
    Name() string
    IsRunning() bool

    IsNamedCycle() bool

    Serve() error

    // onExit allows individual services to register a callback function which will be called when a service is asked to
    // exit. Usually services terminate themselves when the callback is called
    OnExit(callback func(interface{})) error

    SetRunning()
    SetStopped()

    GetWaiters() []*waiter

    // this is only allowed for named services
    //Rebuild() error
    // cleanup is only allowed for unnamed services
    Cleanup() error
}

// --- --- --- --- --- --- --- --- Service internal structure to coalesce functions --- --- --- --- --- --- --- --- - //

type ServeFunc func() error
func (s ServeFunc) Serve() error {
    return s()
}

type OnExitFunc func(callback func(interface{})) error
func (o OnExitFunc) OnExit(callback func(interface{})) error {
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

func (s *srvcFuncs) IsNamedCycle() bool {
    return (len(s.name) != 0)
}

func (s *srvcFuncs) SetRunning() {
    s.Lock()
    defer s.Unlock()

    s.state = serviceRunning
}

func (s *srvcFuncs) SetStopped() {
    s.Lock()
    defer s.Unlock()

    s.state = serviceStopped
}

func (s *srvcFuncs) GetWaiters() []*waiter {
    return s.waiters
}

/*
func (s *srvcFuncs) Rebuild() error {
    if !s.IsNamedCycle() {
        return errors.Errorf("[ERR] only named service is allowed to rebuild")
    }
    return nil
}
*/

func (s *srvcFuncs) Cleanup() error {
    if s.IsNamedCycle() {
        return errors.Errorf("[ERR] only unnamed service is allowed to clean up")
    }

    for i := range s.waiters {
        w := s.waiters[i]
        close(w.eventC)
    }
    return nil
}

// --- --- --- --- --- --- --- --- --- --- Service Builder Options and Functions --- --- --- --- --- --- --- --- ---- //

// ServerOption is a functional option passed to the server
type ServiceOption func(app AppSupervisor, s Service) error

/*
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
*/

func BindEventWithService(eventName string, eventC chan Event) ServiceOption {
    return func(app AppSupervisor, s Service) error {
        srv, ok := s.(*srvcFuncs)
        if !ok {
            return errors.Errorf("[ERR] invalid service type to bind event")
        }
        if srv == nil {
            return errors.Errorf("[ERR] null service instance to bind event")
        }
        if srv.IsNamedCycle() {
            return errors.Errorf("[ERR] named service instance cannot be bound with event")
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
            name:       eventName,
            eventC:     eventC,
        }
        srv.waiters = append(srv.waiters, w)

        if !s.IsNamedCycle() {
            sup.eventWaiters[eventName] = append(sup.eventWaiters[eventName], w)
        }
        return nil
    }
}

// --- --- --- --- --- --- --- --- --- --- AppSupervisor Options and Functions --- --- --- --- --- --- --- --- --- -- //

type AppSupervisor interface {
    BroadcastEvent(event Event)

    Register(srv Service) error
    RegisterServiceWithFuncs(sfn ServeFunc, efn OnExitFunc, options... ServiceOption) error
    RegisterNamedServiceWithFuncs(name string, sfn ServeFunc, efn OnExitFunc) error

    IsStopped() bool
    StopChannel() <- chan struct{}

    Start() error
    RunNamedService(name string) error
    Stop() error

    Wait() error

    ServiceCount() int
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

// --- --- --- --- --- --- --- --- --- --- Event Related Methods --- --- --- --- --- --- --- --- --- --- --- --- ---- //

func (a *appSupervisor) getWaiters(name string) []*waiter {
    a.Lock()
    defer a.Unlock()

    waiters := a.eventWaiters[name]
    out := make([]*waiter, len(waiters))
    for i := range waiters {
        out[i] = waiters[i]
    }
    return out
}

func (a *appSupervisor) notifyWaiter(w *waiter, event Event) {
    select {
        case w.eventC <- event:
    }
}

func (a *appSupervisor) fanOut() {
    defer a.serviceWG.Done()

    for {
        select {
            case <-a.stoppedC:
                return
            case event := <-a.eventsC:
                waiters := a.getWaiters(event.Name)
                for _, waiter := range waiters {
                    a.notifyWaiter(waiter, event)
                }
        }
    }
}

func (a *appSupervisor) BroadcastEvent(event Event) {
    a.Lock()
    defer a.Unlock()

    go func() {
        a.eventsC <- event
    }()
}

// --- --- --- --- --- --- --- --- --- --- Service Related Methods --- --- --- --- --- --- --- --- --- --- --- --- -- //

func (a *appSupervisor) runService(service Service) {
    // this func will be called _after_ a service stops running:
    serviceCleanup := func(as *appSupervisor, srv Service) {
        as.Lock()
        defer as.Unlock()

        // since named cycles don't have events associated, we don't need to clean up
        if srv.IsNamedCycle() {
            return
        }

        // cleanup events first
        swaiters := srv.GetWaiters()
        for o := range swaiters {
            sw := swaiters[o]
            awaiters := as.eventWaiters[sw.name]
            if len(awaiters) != 0 {
                var nwaiters []*waiter = []*waiter{}
                for i := range awaiters {
                    aw := awaiters[i]
                    if aw != sw {
                        nwaiters = append(nwaiters, aw)
                    }
                }
                as.eventWaiters[sw.name] = nwaiters
            }
            log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] waiter [%s] cleaned", srv.Name(), srv, sw.name)
        }

        // cleanup service itself
        for i := range as.services {
            el := as.services[i]
            if el == srv {
                as.services = append(as.services[:i], as.services[i+1:]...)
                break
            }
        }
        err := srv.Cleanup()
        if err != nil {
            log.Debug(errors.WithStack(err))
        }
        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] removed", srv.Name(), srv)
    }

    a.serviceWG.Add(1)
    go func(as *appSupervisor, wg *sync.WaitGroup, srv Service, cleanup func(as *appSupervisor, srv Service)) {
        defer wg.Done()

        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] started", srv.Name(), srv)

        srv.SetRunning()
        err := srv.Serve()
        if err != nil {
            log.Debug(errors.WithStack(err))
        }
        cleanup(as, srv)
        err = srv.OnExit(nil)
        if err != nil {
            log.Debug(errors.WithStack(err))
        }
        srv.SetStopped()

        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] exited", srv.Name(), srv)
    }(a, a.serviceWG, service, serviceCleanup)
}

func (a *appSupervisor) Register(srv Service) error {
    a.Lock()
    defer a.Unlock()

    // when a service is named, check if there is a service with the same name
    if srv.IsNamedCycle() {
        for _, es := range a.services {
            if srv.Name() == es.Name() {
                return errors.Errorf("[ERR] a service with the same name '%s' exists", srv.Name())
            }
        }
    }
    a.services = append(a.services, srv)

    log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] added", srv.Name(), srv)

    if a.state == stateStarted && !srv.IsNamedCycle() {
        a.runService(srv)
    }
    return nil
}

func (a *appSupervisor) RegisterServiceWithFuncs(sfn ServeFunc, efn OnExitFunc, options... ServiceOption) error {
    srv := &srvcFuncs {
        state:         serviceStopped,
        waiters:       []*waiter{},
        ServeFunc:     sfn,
        OnExitFunc:    efn,
    }

    for _, opt := range options {
        if err := opt(a, srv); err != nil {
            return errors.WithStack(err)
        }
    }

    return a.Register(srv)
}

func (a *appSupervisor) RegisterNamedServiceWithFuncs(name string, sfn ServeFunc, efn OnExitFunc) error {
    return a.Register(&srvcFuncs {
        name:          name,
        state:         serviceStopped,
//        waiters:       []*waiter{},
        ServeFunc:     sfn,
        OnExitFunc:    efn,
    })
}

func (a *appSupervisor) IsStopped() bool {
    a.Lock()
    defer a.Unlock()

    if a.state == stateStarted {
        return false
    } else {
        return true
    }
}

func (a *appSupervisor) StopChannel() <- chan struct{} {
    return a.stoppedC
}

func (a *appSupervisor) Start() error {
    a.Lock()
    defer a.Unlock()

    if a.state == stateStarted {
        return nil
    }
    a.state = stateStarted

    if len(a.services) == 0 {
        log.Debugf("\n\n[SUPERVISOR-SERVICE] starting services...")
        return nil
    }

    for _, srv := range a.services {
        if !srv.IsNamedCycle() {
            a.runService(srv)
        }
    }

    return nil
}

func (a *appSupervisor) RunNamedService(name string) error {
    a.Lock()
    defer a.Unlock()

    if a.state != stateStarted {
        return nil
    }

    if len(name) == 0 {
        return errors.Errorf("[ERR] cannot find a service without name")
    }

    for _, srv := range a.services {
        if srv.IsNamedCycle() && srv.Name() == name {
            if !srv.IsRunning() {
/*
                swaiters := srv.GetWaiters()
                for i := range swaiters {
                    w := swaiters[i]
                    p.eventWaiters[w.name] = append(p.eventWaiters[w.name], w)
                }
*/
                a.runService(srv)

                return nil
            } else {
                return ServiceRunning
            }
        }
    }

    return errors.Errorf("[ERR] cannot find a service named %s", name)
}

func (a *appSupervisor) Stop() error {
    // this double locking to to prevent ServiceFunc from deadlocked, but enable other variables to be reset
    a.Lock()
    defer a.Unlock()

    if a.state == stateStopped {
        return nil
    }
    a.state = stateStopped

    // we broadcast stopping and wait for all goroutines closed with event channels intact to give grace period
    log.Debugf("[SUPERVISOR-SERVICE] stopping all services...\n")
    close(a.stoppedC)
    return nil
}

func (a *appSupervisor) Wait() error {
    var (
        waiters []*waiter
        w *waiter
    )
    a.serviceWG.Wait()

    // we close event channel after stopping all go routines and fanOut function that we can safely close
    close(a.eventsC)
    return nil

    for _, waiters = range a.eventWaiters {
        for _, w = range waiters {
            close(w.eventC)
        }
    }

    return nil
}

// ServiceCount returns the number of registered and actively running services
func (a *appSupervisor) ServiceCount() int {
    a.Lock()
    defer a.Unlock()

    return len(a.services)
}
