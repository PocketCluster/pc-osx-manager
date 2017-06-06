package swarm

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "time"
    "net"
    "net/http"
    "strings"
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/docker/swarm/api"
    "gopkg.in/tylerb/graceful.v1"

    "github.com/docker/swarm/cluster"
    "github.com/docker/swarm/cluster/swarm"
    "github.com/docker/swarm/scheduler"
    "github.com/docker/swarm/scheduler/filter"
    "github.com/docker/swarm/scheduler/strategy"
)

// DEPRECATED CONTEXT METHODS //

func NewContext(host, nodeList string, tlsCa, tlsCert, tlsKey string) *SwarmContext {
    discoveryOpt := make(map[string]string)
    clusterOpt := cluster.DriverOpts{}
    tlsConfig, err := loadTLSConfigFromFiles(tlsCa, tlsCert, tlsKey, true)
    if err != nil {
        log.Fatal(err)
    }

    return &SwarmContext {
        discoveryOpt:       discoveryOpt,

        heartbeat:          time.Duration(10 * time.Second),
        refreshMinInterval: time.Duration(5 * time.Second),
        refreshMaxInterval: time.Duration(30 * time.Second),
        refreshRetry:       time.Duration(3),
        failureRetry:       6,

        managerHost:        []string{host},
        debug:              true,

        clusterOpt:         clusterOpt,
        strategy:           "spread",
        tlsConfig:          tlsConfig,
    }
}

// Load the TLS certificates/keys and, if verify is true, the CA.
func loadTLSConfigFromFiles(ca, cert, key string, verify bool) (*tls.Config, error) {
    c, err := tls.LoadX509KeyPair(cert, key)
    if err != nil {
        return nil, fmt.Errorf("Couldn't load X509 key pair (%s, %s): %s. Key encrypted?",
            cert, key, err)
    }

    config := &tls.Config{
        Certificates: []tls.Certificate{c},
        MinVersion:   tls.VersionTLS10,
    }

    if verify {
        certPool := x509.NewCertPool()
        file, err := ioutil.ReadFile(ca)
        if err != nil {
            return nil, fmt.Errorf("Couldn't read CA certificate: %s", err)
        }
        certPool.AppendCertsFromPEM(file)
        config.RootCAs = certPool
        config.ClientAuth = tls.RequireAndVerifyClientCert
        config.ClientCAs = certPool
    } else {
        // If --tlsverify is not supplied, disable CA validation.
        config.InsecureSkipVerify = true
    }

    return config, nil
}

func (context *SwarmContext) Manage() error {
    refreshMinInterval := context.refreshMinInterval
    refreshMaxInterval := context.refreshMaxInterval
    if refreshMinInterval <= time.Duration(0) * time.Second {
        return errors.Errorf("min refresh interval should be a positive number")
    }
    if refreshMaxInterval < refreshMinInterval {
        return errors.Errorf("max refresh interval cannot be less than min refresh interval")
    }
    // engine-refresh-retry is deprecated
    refreshRetry := context.refreshRetry
    if refreshRetry != 3 {
        return errors.Errorf("--engine-refresh-retry is deprecated. Use --engine-failure-retry")
    }
    failureRetry := context.failureRetry
    if failureRetry <= 0 {
        return errors.Errorf("invalid failure retry count")
    }
    engineOpts := &cluster.EngineOpts {
        RefreshMinInterval: refreshMinInterval,
        RefreshMaxInterval: refreshMaxInterval,
        FailureRetry:       failureRetry,
    }

    discovery := createNodeDiscovery(context)
    s, err := strategy.New(context.strategy)
    if err != nil {
        return errors.WithStack(err)
    }

    // see https://github.com/codegangsta/cli/issues/160
    names := []string{"health", "port", "containerslots", "dependency", "affinity", "constraint"}
    fs, err := filter.New(names)
    if err != nil {
        return errors.WithStack(err)
    }

    sched := scheduler.New(s, fs)
    var cl cluster.Cluster
    cl, err = swarm.NewCluster(sched, context.tlsConfig, discovery, context.clusterOpt, engineOpts)
    if err != nil {
        return errors.WithStack(err)
    }

    hosts := context.managerHost
    server := newService(hosts, context.tlsConfig)
    primary := api.NewPrimary(cl, context.tlsConfig, &statusHandler{cl, nil, nil}, context.debug, context.cors)
    server.SetHandler(primary)
    cluster.NewWatchdog(cl)

    return errors.WithStack(server.ListenAndServeMultiHosts())
}

// DEPRECATED SERVICE METHODS //

// Dispatcher is a meta http.Handler. It acts as an http.Handler and forwards
// requests to another http.Handler that can be changed at runtime.
type dispatcher struct {
    handler http.Handler
}

// SetHandler changes the underlying handler.
func (d *dispatcher) SetHandler(handler http.Handler) {
    d.handler = handler
}

// ServeHTTP forwards requests to the underlying handler.
func (d *dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if d.handler == nil {
        //httpError(w, "No dispatcher defined", http.StatusInternalServerError)
        err, status := "No dispatcher defined", http.StatusInternalServerError
        log.WithField("status", status).Errorf("HTTP error: %v", err)
        http.Error(w, err, status)
        return
    }
    d.handler.ServeHTTP(w, r)
}

func newService(hosts []string, tlsConfig *tls.Config) *SwarmService {
    return &SwarmService{
        hosts:      hosts,
        tlsConfig:  tlsConfig,
        dispatcher: &dispatcher{},
    }
}

// SetHandler is used to overwrite the HTTP handler for the API.
// This can be the api router or a reverse proxy.
func (s *SwarmService) SetHandler(handler http.Handler) {
    s.dispatcher.SetHandler(handler)
}

// ListenAndServe starts an HTTP server on each host to listen on its
// TCP or Unix network address and calls Serve on each host's server
// to handle requests on incoming connections.
//
// The expected format for a host string is [protocol://]address. The protocol
// must be either "tcp" or "unix", with "tcp" used by default if not specified.
func (s *SwarmService) ListenAndServeMultiHosts() error {
    chErrors := make(chan error, len(s.hosts))

    for _, host := range s.hosts {
        protoAddrParts := strings.SplitN(host, "://", 2)
        if len(protoAddrParts) == 1 {
            protoAddrParts = append([]string{"tcp"}, protoAddrParts...)
        }

        go func() {
            log.WithFields(log.Fields{"proto": protoAddrParts[0], "addr": protoAddrParts[1]}).Info("Listening for HTTP")

            var (
                l      net.Listener
                err    error
                server = &http.Server{
                    Addr:    protoAddrParts[1],
                    Handler: s.dispatcher,
                }
            )

            switch protoAddrParts[0] {
            //case "unix":
            //    l, err = newUnixListener(protoAddrParts[1], s.tlsConfig)
            case "tcp":
                l, err = newListener("tcp", protoAddrParts[1], s.tlsConfig)
            default:
                err = errors.Errorf("unsupported protocol: %q", protoAddrParts[0])
            }

            if err != nil {
                chErrors <- err
            } else {
                chErrors <- server.Serve(l)
            }
        }()
    }

    for i := 0; i < len(s.hosts); i++ {
        err := <-chErrors
        if err != nil {
            return err
        }
    }
    return nil
}

func (s *SwarmService) ListenAndServeMultiHostsOnWaitGroup(wg *sync.WaitGroup) ([]*graceful.Server, []error) {
    var (
        chErrors    = make(chan error, len(s.hosts))
        chServers   = make(chan *graceful.Server, len(s.hosts))

        slErrors    = []error{}
        slServers   = []*graceful.Server{}
    )

    for _, host := range s.hosts {
        protoAddrParts := strings.SplitN(host, "://", 2)
        if len(protoAddrParts) == 1 {
            protoAddrParts = append([]string{"tcp"}, protoAddrParts...)
        }

        go func() {
            defer wg.Done()
            log.WithFields(log.Fields{"proto": protoAddrParts[0], "addr": protoAddrParts[1]}).Info("Listening for HTTP")

            var (
                l      net.Listener
                err    error
                server = &graceful.Server{
                    Timeout: 10 * time.Second,
                    NoSignalHandling: true,
                    Server: &http.Server{
                        Addr:    protoAddrParts[1],
                        Handler: s.dispatcher,
                    },
                }
            )

            switch protoAddrParts[0] {
                //case "unix":
                //    l, err = newUnixListener(protoAddrParts[1], s.tlsConfig)
                case "tcp":
                    l, err = newListener("tcp", protoAddrParts[1], s.tlsConfig)
                default:
                    err = errors.Errorf("unsupported protocol: %q", protoAddrParts[0])
            }

            if err != nil {
                chErrors <- err
                chServers <- nil
            } else {
                chErrors <- server.Serve(l)
                chServers <- server
            }
        }()
    }

    for i := 0; i < len(s.hosts); i++ {
        err := <-chErrors
        if err != nil {
            slErrors = append(slErrors, err)
        }

        srv := <-chServers
        if srv != nil {
            slServers = append(slServers, srv)
        }
    }
    return slServers, slErrors
}

// ListenAndServeOnWaitGroup starts an HTTP server on the first host in the list to listen on its
// TCP or Unix network address and calls Serve on each host's server
// to handle requests on incoming connections.
//
// The expected format for a host string is [protocol://]address. The protocol
// must be either "tcp" or "unix", with "tcp" used by default if not specified.
func (s *SwarmService) ListenAndServeOnWaitGroup(wg *sync.WaitGroup) (*graceful.Server, error) {
    var (
        chErrors    = make(chan error)
        chServers   = make(chan *graceful.Server)
    )

    host := s.hosts[0]
    protoAddrParts := strings.SplitN(host, "://", 2)
    if len(protoAddrParts) == 1 {
        protoAddrParts = append([]string{"tcp"}, protoAddrParts...)
    }

    go func(wg *sync.WaitGroup) {
        defer wg.Done()
        log.WithFields(log.Fields{"proto": protoAddrParts[0], "addr": protoAddrParts[1]}).Info("Listening for HTTP")

        var (
            l      net.Listener
            err    error
            server = &graceful.Server{
                Timeout: 10 * time.Second,
                NoSignalHandling: true,
                Server: &http.Server{
                    Addr:    protoAddrParts[1],
                    Handler: s.dispatcher,
                },
            }
        )

        switch protoAddrParts[0] {
            //case "unix":
            //    l, err = newUnixListener(protoAddrParts[1], s.tlsConfig)
            case "tcp":
                l, err = newListener("tcp", protoAddrParts[1], s.tlsConfig)
            default:
                err = errors.Errorf("unsupported protocol: %q", protoAddrParts[0])
        }

        chServers <- server
        chErrors <- err
        server.Serve(l)
/*
        // TODO : this error message has to be routed toward UI layer
        if err != nil {
            chErrors <- err
        } else {
            chErrors <- server.Serve(l)
        }
*/
    }(wg)

    srv := <-chServers
    err := <-chErrors
    log.Print("all the values retrieved")
    return srv, err
}