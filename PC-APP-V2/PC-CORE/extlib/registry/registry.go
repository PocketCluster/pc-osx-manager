package registry

import (
    "crypto/tls"
    "fmt"
    "io"
    "net"
    "net/http"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/docker/distribution/configuration"
    "github.com/docker/distribution/context"
    "github.com/docker/distribution/health"
    "github.com/docker/distribution/registry/handlers"
    "github.com/docker/distribution/registry/listener"
    "github.com/docker/distribution/uuid"
    "github.com/docker/distribution/version"

    "gopkg.in/tylerb/graceful.v1"
)

// NewRegistry creates a new registry from a context and configuration struct.
func NewPocketRegistry(config *PocketRegistryConfig) (*PocketRegistry, error) {
    var (
        ctx = context.WithVersion(context.Background(), version.Version)
    )
/*
    if config.regConfig.HTTP.Debug.Addr != "" {
        go func(addr string) {
            log.Infof("debug server listening %v", addr)
            if err := http.ListenAndServe(addr, nil); err != nil {
                log.Fatalf("error listening on debug interface: %v", err)
            }
        }(config.regConfig.HTTP.Debug.Addr)
    }
*/
    // inject a logger into the uuid library. warns us if there is a problem
    // with uuid generation under low entropy.
    uuid.Loggerf = context.GetLogger(ctx).Warnf

    app := handlers.NewApp(ctx, config.regConfig)
    // TODO(aaronl): The global scope of the health checks means NewRegistry
    // can only be called once per process.
    app.RegisterHealthChecks()
    handler := configureReporting(app)
    handler = alive("/", handler)
    handler = health.Handler(handler)
    handler = panicHandler(handler)

/*
	// (04/16/2017) logging is disabled for now
    if !config.Log.AccessLog.Disabled {
        handler = gorhandlers.CombinedLoggingHandler(os.Stdout, handler)
    }
*/

    return &PocketRegistry{
        app:    app,
        config: config,
        server: &graceful.Server{
            Timeout: 10 * time.Second,
            NoSignalHandling: true,
            Server: &http.Server{
                Addr: config.regConfig.HTTP.Addr,
                Handler: handler,
            },
        },
    }, nil
}

// A Registry represents a complete instance of the registry.
type PocketRegistry struct {
    config        *PocketRegistryConfig
    app           *handlers.App
    server 		  *graceful.Server
    listener      io.Closer
}

// ListenAndServe runs the registry's HTTP server.
func (r *PocketRegistry) Start() (error) {
    config := r.config

    ln, err := listener.NewListener(config.regConfig.HTTP.Net, config.regConfig.HTTP.Addr)
    if err != nil {
        return err
    }

    // TODO : No HTTP secret provided - generated random secret. This may cause problems with uploads if multiple registries are behind a load-balancer. To provide a shared secret, fill in http.secret in the configuration file or set the REGISTRY_HTTP_SECRET environment variable.
    ln = tls.NewListener(ln, config.tlsConfig)
    log.Debugf("[REG] listening on %v, tls", ln.Addr())

    // start serving
    go func(srv *graceful.Server, l net.Listener) {
        var err = srv.Serve(l)
        if err != nil {
            log.Println("HTTP Server Error - ", err)
        }
    }(r.server, ln)

    r.listener = ln
    return nil
}

func (r *PocketRegistry) Stop(timeout time.Duration) error {
    if r.listener == nil {
        return nil
    }
    ln := r.listener
    r.listener = nil

    err := ln.Close()
    r.server.Stop(timeout)
    return errors.WithStack(err)
}

func configureReporting(app *handlers.App) http.Handler {
    var handler http.Handler = app
/*
    TODO : (04/16/2017) replace this with sentry.io
    if app.Config.Reporting.Bugsnag.APIKey != "" {
        bugsnagConfig := bugsnag.Configuration{
            APIKey: app.Config.Reporting.Bugsnag.APIKey,
            // TODO(brianbland): provide the registry version here
            // AppVersion: "2.0",
        }
        if app.Config.Reporting.Bugsnag.ReleaseStage != "" {
            bugsnagConfig.ReleaseStage = app.Config.Reporting.Bugsnag.ReleaseStage
        }
        if app.Config.Reporting.Bugsnag.Endpoint != "" {
            bugsnagConfig.Endpoint = app.Config.Reporting.Bugsnag.Endpoint
        }
        bugsnag.Configure(bugsnagConfig)

        handler = bugsnag.Handler(handler)
    }
*/

    return handler
}

// TODO : insert sentry.io here
// panicHandler add an HTTP handler to web app. The handler recover the happening
// panic. logrus.Panic transmits panic message to pre-config log hooks, which is
// defined in config.yml.
func panicHandler(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Panic(fmt.Sprintf("%v", err))
            }
        }()
        handler.ServeHTTP(w, r)
    })
}

// alive simply wraps the handler with a route that always returns an http 200
// response when the path is matched. If the path is not matched, the request
// is passed to the provided handler. There is no guarantee of anything but
// that the server is up. Wrap with other handlers (such as health.Handler)
// for greater affect.
func alive(path string, handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == path {
            w.Header().Set("Cache-Control", "no-cache")
            w.WriteHeader(http.StatusOK)
            return
        }

        handler.ServeHTTP(w, r)
    })
}

func nextProtos(config *configuration.Configuration) []string {
    switch config.HTTP.HTTP2.Disabled {
    case true:
        return []string{"http/1.1"}
    default:
        return []string{"h2", "http/1.1"}
    }
}
