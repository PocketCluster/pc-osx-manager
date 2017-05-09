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

type PocketApplication struct {
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
func NewPocketApplication() *PocketApplication {
    srv := &PocketApplication{
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
    for {
        select {
        case event := <-p.eventsC:
            waiters := p.getWaiters(event.Name)
            for _, waiter := range waiters {
                go p.notifyWaiter(waiter, event)
            }
        case <-p.closer.C:
            return
        }
    }
}
func (p *PocketApplication) serve(srv *Service) {
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
        log.Debugf("[APPLICATION] Service %v is done (%v)", *srv, len(p.services))
    }

    p.wg.Add(1)
    go func() {
        defer p.wg.Done()
        defer removeService()

        log.Debugf("[APPLICATION] Service %v started (%v)", *srv, p.ServiceCount())
        err := (*srv).Serve()
        if err != nil {
            errors.WithStack(err)
        }
    }()
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
    log.Debugf("PocketApplication.BroadcastEvent: %v", &event)

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

func (p *PocketApplication) Wait() error {
    defer p.closer.Close()
    p.wg.Wait()
    return nil
}

// onExit allows individual services to register a callback function which will be
// called when Teleport Process is asked to exit. Usually services terminate themselves
// when the callback is called
func (p *PocketApplication) OnExit(callback func(interface{})) {
    go func() {
        select {
            case <- p.closer.C:
                callback(nil)
        }
    }()
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
