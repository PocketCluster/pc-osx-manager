package swarm

import (
    "crypto/tls"
    "fmt"
    "net"
    "net/http"
    "strings"
    "time"
    "sync"

    log "github.com/Sirupsen/logrus"
    "github.com/docker/swarm/api"
    "gopkg.in/tylerb/graceful.v1"
)

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

// Server is a Docker API server.
type Server struct {
    hosts      []string
    tlsConfig  *tls.Config
    dispatcher *dispatcher
}

// NewServer creates an api.Server.
func NewServer(hosts []string, tlsConfig *tls.Config) *Server {
    return &Server{
        hosts:      hosts,
        tlsConfig:  tlsConfig,
        dispatcher: &dispatcher{},
    }
}

// SetHandler is used to overwrite the HTTP handler for the API.
// This can be the api router or a reverse proxy.
func (s *Server) SetHandler(handler http.Handler) {
    s.dispatcher.SetHandler(handler)
}

func newListener(proto, addr string, tlsConfig *tls.Config) (net.Listener, error) {
    l, err := net.Listen(proto, addr)
    if err != nil {
        if strings.Contains(err.Error(), "address already in use") && strings.Contains(addr, api.DefaultDockerPort) {
            return nil, fmt.Errorf("%s: is Docker already running on this machine? Try using a different port", err)
        }
        return nil, err
    }
    if tlsConfig != nil {
        tlsConfig.NextProtos = []string{"http/1.1"}
        l = tls.NewListener(l, tlsConfig)
    }
    return l, nil
}

// ListenAndServe starts an HTTP server on each host to listen on its
// TCP or Unix network address and calls Serve on each host's server
// to handle requests on incoming connections.
//
// The expected format for a host string is [protocol://]address. The protocol
// must be either "tcp" or "unix", with "tcp" used by default if not specified.
func (s *Server) ListenAndServe() error {
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
                err = fmt.Errorf("unsupported protocol: %q", protoAddrParts[0])
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

func (s *Server) ListenAndServeMultiHostsOnWaitGroup(wg *sync.WaitGroup) ([]*graceful.Server, []error) {
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
                err = fmt.Errorf("unsupported protocol: %q", protoAddrParts[0])
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
func (s *Server) ListenAndServeOnWaitGroup(wg *sync.WaitGroup) (*graceful.Server, error) {
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
            err = fmt.Errorf("unsupported protocol: %q", protoAddrParts[0])
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