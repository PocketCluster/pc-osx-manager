package service

import (
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
)


const (
    broadcastChannelSize = 64
)

const (
    stateStopped    = iota
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

// Event is a special service event that can be generated by various goroutines in the app
type Event struct {
    Name       string
    Payload    interface{}
}

type waiter struct {
    name       string
    eventC     chan Event
}

// --- --- --- --- --- --- --- --- --- --- --- --- --- service unit --- --- --- --- --- --- --- --- --- --- --- --- - //

type Service interface {
    // These indicate that the service is named and should be individually managed to start.
    // Plus, only named service can be rebuild
    Name() string

    IsRunning() bool

    // Serve() allows services to run
    serve() (interface{}, error)

    // OnExit allows individual services to register a callback function which will be called when a service is asked to
    // exit. Usually services terminate themselves when the callback is called
    onExit(residue interface{}, runerr error) error

    // --- internal methods --- //

    isNamedCycle() bool

    setRunning()
    setStopped()

    getWaiters() []*waiter

    // cleanup is only allowed for unnamed services
    cleanup() error
}

// --- --- --- --- --- --- --- --- Service internal structure to coalesce functions --- --- --- --- --- --- --- --- - //

type serveFunc func() (interface{}, error)
func (s serveFunc) serve() (interface{}, error) {
    return s()
}

type onExitFunc func(residue interface{}, err error) error
func (o onExitFunc) onExit(residue interface{}, err error) error {
    return o(residue, err)
}

type srvcFuncs struct {
    sync.Mutex
    state      int

    name       string
    waiters    []*waiter
    serveFunc
    onExitFunc
}

func (s *srvcFuncs) Name() string {
    return s.name
}

func (s *srvcFuncs) IsRunning() bool {
    s.Lock()
    defer s.Unlock()

    return (s.state == serviceRunning)
}

func (s *srvcFuncs) isNamedCycle() bool {
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

func (s *srvcFuncs) getWaiters() []*waiter {
    return s.waiters
}

func (s *srvcFuncs) cleanup() error {
    if s.isNamedCycle() {
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
type ServiceOption func(app ServiceSupervisor, s Service) error

func BindEventWithService(eventName string, eventC chan Event) ServiceOption {
    return func(app ServiceSupervisor, s Service) error {
        srv, ok := s.(*srvcFuncs)
        if !ok {
            return errors.Errorf("[ERR] invalid service type to bind event")
        }
        if srv == nil {
            return errors.Errorf("[ERR] null service instance to bind event")
        }
        if srv.isNamedCycle() {
            return errors.Errorf("[ERR] named service instance cannot be bound with event")
        }
        sup, ok := app.(*srvcSupervisor)
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

        if !s.isNamedCycle() {
            sup.eventWaiters[eventName] = append(sup.eventWaiters[eventName], w)
        }
        return nil
    }
}

// --- --- --- --- --- --- --- --- --- --- ServiceSupervisor Options and Functions --- --- --- --- --- --- --- --- -- //

type ServiceSupervisor interface {
    BroadcastEvent(event Event)

    RegisterServiceWithFuncs(sfn serveFunc, efn onExitFunc, options... ServiceOption) error
    RegisterNamedServiceWithFuncs(name string, sfn serveFunc, efn onExitFunc) error

    IsStopped() bool
    StopChannel() <- chan struct{}

    StartServices() error
    RunNamedService(name string) error
    StopServices() error
    Refresh() error

    // --- internal service --- //

    serviceCount() int
}

func NewServiceSupervisor() ServiceSupervisor {
    return &srvcSupervisor{
        state:           stateStopped,
        serviceWG:       &sync.WaitGroup{},
        services:        []Service{},

        eventsC:         make(chan Event, broadcastChannelSize),
        eventWaiters:    make(map[string][]*waiter),

        stoppedC:        make(chan struct{}),
    }
}

type srvcSupervisor struct {
    sync.Mutex
    state int

    serviceWG       *sync.WaitGroup
    services        []Service

    eventWaiters    map[string][]*waiter
    eventsC         chan Event
    stoppedC        chan struct{}
}

// --- --- --- --- --- --- --- --- --- --- Event Related Methods --- --- --- --- --- --- --- --- --- --- --- --- ---- //

func (s *srvcSupervisor) getWaiters(name string) []*waiter {
    s.Lock()
    defer s.Unlock()

    waiters := s.eventWaiters[name]
    out := make([]*waiter, len(waiters))
    for i := range waiters {
        out[i] = waiters[i]
    }
    return out
}

func (s *srvcSupervisor) notifyWaiter(w *waiter, evt Event) {
    select {
        case w.eventC <- evt:
    }
}

func (s *srvcSupervisor) fanOut() {
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

func (s *srvcSupervisor) BroadcastEvent(event Event) {
    s.Lock()
    defer s.Unlock()

    go func() {
        s.eventsC <- event
    }()
}

// --- --- --- --- --- --- --- --- --- --- Service Related Methods --- --- --- --- --- --- --- --- --- --- --- --- -- //

func (s *srvcSupervisor) runService(service Service) {
    // this func will be called _after_ a service stops running:
    serviceCleanup := func(ss *srvcSupervisor, srv Service) {
        ss.Lock()
        defer ss.Unlock()

        // since named cycles don't have events associated, we don't need to clean up
        if srv.isNamedCycle() {
            return
        }

        // cleanup events first
        swaiters := srv.getWaiters()
        for o := range swaiters {
            sw := swaiters[o]
            awaiters := ss.eventWaiters[sw.name]
            if len(awaiters) != 0 {
                var nwaiters []*waiter = []*waiter{}
                for i := range awaiters {
                    aw := awaiters[i]
                    if aw != sw {
                        nwaiters = append(nwaiters, aw)
                    }
                }
                ss.eventWaiters[sw.name] = nwaiters
            }
            log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] waiter [%s] cleaned", srv.Name(), srv, sw.name)
        }

        // cleanup service itself
        for i := range ss.services {
            el := ss.services[i]
            if el == srv {
                ss.services = append(ss.services[:i], ss.services[i+1:]...)
                break
            }
        }
        err := srv.cleanup()
        if err != nil {
            log.Debug(errors.WithStack(err))
        }
        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] removed", srv.Name(), srv)
    }

    s.serviceWG.Add(1)
    go func(ss *srvcSupervisor, wg *sync.WaitGroup, srv Service, cleanup func(as *srvcSupervisor, srv Service)) {
        defer wg.Done()

        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] started", srv.Name(), srv)

        srv.setRunning()
        residue, rerr := srv.serve()
        if rerr != nil {
            log.Debug(errors.WithStack(rerr))
        }
        cleanup(ss, srv)
        oerr := srv.onExit(residue, rerr)
        if oerr != nil {
            log.Debug(errors.WithStack(oerr))
        }
        srv.setStopped()

        log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] exited", srv.Name(), srv)
    }(s, s.serviceWG, service, serviceCleanup)
}

func (s *srvcSupervisor) registerService(srv Service) error {
    s.Lock()
    defer s.Unlock()

    // when a service is named, check if there is a service with the same name
    if srv.isNamedCycle() {
        for _, es := range s.services {
            if srv.Name() == es.Name() {
                return errors.Errorf("[ERR] a service with the same name '%s' exists", srv.Name())
            }
        }
    }
    s.services = append(s.services, srv)

    log.Debugf("[SUPERVISOR-SERVICE] ['%s' | %v] added", srv.Name(), srv)

    if s.state == stateStarted && !srv.isNamedCycle() {
        s.runService(srv)
    }
    return nil
}

func (s *srvcSupervisor) RegisterServiceWithFuncs(sfn serveFunc, efn onExitFunc, options... ServiceOption) error {
    srv := &srvcFuncs {
        state:         serviceStopped,
        waiters:       []*waiter{},
        serveFunc:     sfn,
        onExitFunc:    efn,
    }

    for _, opt := range options {
        if err := opt(s, srv); err != nil {
            return errors.WithStack(err)
        }
    }

    return s.registerService(srv)
}

func (s *srvcSupervisor) RegisterNamedServiceWithFuncs(name string, sfn serveFunc, efn onExitFunc) error {
    return s.registerService(&srvcFuncs {
        name:          name,
        state:         serviceStopped,
//        waiters:       []*waiter{},
        serveFunc:     sfn,
        onExitFunc:    efn,
    })
}

func (s *srvcSupervisor) IsStopped() bool {
    s.Lock()
    defer s.Unlock()

    if s.state == stateStopped {
        return true
    } else {
        return false
    }
}

func (s *srvcSupervisor) StopChannel() <- chan struct{} {
    return s.stoppedC
}

func (s *srvcSupervisor) StartServices() error {
    s.Lock()
    defer s.Unlock()

    if s.state == stateStarted {
        return nil
    }
    s.state = stateStarted

    s.serviceWG.Add(1)
    go s.fanOut()

    if len(s.services) == 0 {
        log.Debugf("[SUPERVISOR] Start() :: nothing to run")
        return nil
    }

    for _, srv := range s.services {
        if !srv.isNamedCycle() {
            s.runService(srv)
        }
    }

    return nil
}

func (s *srvcSupervisor) RunNamedService(name string) error {
    s.Lock()
    defer s.Unlock()
    if s.state != stateStarted {
        return nil
    }

    if len(name) == 0 {
        return errors.Errorf("[ERR] cannot find a service without name")
    }

    for _, srv := range s.services {
        if srv.isNamedCycle() && srv.Name() == name {
            if !srv.IsRunning() {
                s.runService(srv)
                return nil
            } else {
                return ServiceRunning
            }
        }
    }

    return errors.Errorf("[ERR] cannot find a service named %s", name)
}

func (s *srvcSupervisor) StopServices() error {
    // this double locking to to prevent ServiceFunc from deadlocked, but enable other variables to be reset
    s.Lock()
    defer s.Unlock()

    if s.state == stateStopped {
        return nil
    }
    s.state = stateStopped

    log.Debugf("[SUPERVISOR] stopping services...")

    // we broadcast stopping and wait for all goroutines closed with event channels intact to give grace period
    // while services are stopping, there are many races to grap supervisor lock so it is unlocked before call gets here
    close(s.stoppedC)
    return nil
}

func (p *srvcSupervisor) Refresh() error {
    p.serviceWG.Wait()

    // we close event channel after stopping all go routines and fanOut function that we can safely close
    close(p.eventsC)
    // reset
    p.serviceWG    = &sync.WaitGroup{}
    p.eventsC      = make(chan Event, broadcastChannelSize)
    p.eventWaiters = make(map[string][]*waiter)
    p.stoppedC     = make(chan struct{})

    return nil
}

// ServiceCount returns the number of registered and actively running services
func (s *srvcSupervisor) serviceCount() int {
    s.Lock()
    defer s.Unlock()

    return len(s.services)
}
