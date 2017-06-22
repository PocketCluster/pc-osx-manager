package swarm

import (
    "crypto/tls"
    "net"
    "net/http"
    "strings"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/docker/swarm/api"
    "gopkg.in/tylerb/graceful.v1"
)

func newListener(proto, addr string, tlsConfig *tls.Config) (net.Listener, error) {
    l, err := net.Listen(proto, addr)
    if err != nil {
        if strings.Contains(err.Error(), "address already in use") && strings.Contains(addr, api.DefaultDockerPort) {
            return nil, errors.Errorf("%s: is Docker already running on this machine? Try using a different port", err)
        }
        return nil, err
    }
    if tlsConfig != nil {
        tlsConfig.NextProtos = []string{"http/1.1"}
        l = tls.NewListener(l, tlsConfig)
    }
    return l, nil
}

// NewServer creates an api.Server.
func newStoppableServiceForSingleHost(hosts []string, handler http.Handler, tlsConfig *tls.Config) (*SwarmService, error) {
    var (
        listener net.Listener = nil
        err error = nil
    )
    if len(hosts) == 0 {
        return nil, errors.Errorf("[ERR] serving hosts should be specified")
    }
    if handler == nil {
        return nil, errors.Errorf("[ERR] http handler cannot be nil")
    }
    if tlsConfig == nil {
        return nil, errors.Errorf("[ERR] tls should be configured before setup service")
    }

    protoAddrParts := strings.SplitN(hosts[0], "://", 2)
    if len(protoAddrParts) == 1 {
        protoAddrParts = append([]string{"tcp"}, protoAddrParts...)
    }

    log.WithFields(log.Fields{"proto": protoAddrParts[0], "addr": protoAddrParts[1]}).Info("Listening for HTTP")

    switch protoAddrParts[0] {
        //case "unix":
        //    l, err = newUnixListener(protoAddrParts[1], tlsConfig)
        case "tcp":
            listener, err = newListener("tcp", protoAddrParts[1], tlsConfig)
        default:
            err = errors.Errorf("unsupported protocol: %q", protoAddrParts[0])
    }
    if err != nil {
        return nil, err
    }

    return &SwarmService{
        listener:      listener,
        server:        &graceful.Server{
            NoSignalHandling:   true,
            Server:             &http.Server{
                Addr:           protoAddrParts[1],
                Handler:        handler,
            },
        },
    }, nil
}

// Server is a Docker API server.
type SwarmService struct {
    hosts         []string
    tlsConfig     *tls.Config
    listener      net.Listener
    server        *graceful.Server

    // this field is to be deprecated
    dispatcher    *dispatcher
}

func (s *SwarmService) ListenAndServeSingleHost() error {
    return errors.WithStack(s.server.Serve(s.listener))
}

func (s *SwarmService) Close() error {
    s.server.Stop(10 * time.Second)
    return nil
}
