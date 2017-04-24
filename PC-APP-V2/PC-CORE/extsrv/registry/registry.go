package registry

import (
    "crypto/tls"
    "fmt"
    "net/http"
    "os"
    "time"
    "io"

    log "github.com/Sirupsen/logrus"
    //logstash "github.com/bshuster-repo/logrus-logstash-hook"
    //"github.com/bugsnag/bugsnag-go"
    "github.com/docker/distribution/configuration"
    "github.com/docker/distribution/context"
    "github.com/docker/distribution/health"
    "github.com/docker/distribution/registry/handlers"
    "github.com/docker/distribution/registry/listener"
    "github.com/docker/distribution/uuid"
    "github.com/docker/distribution/version"
)

// NewRegistry creates a new registry from a context and configuration struct.
func NewPocketRegistry(config *PocketRegistryConfig) (*PocketRegistry, error) {
    var (
        err error
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
    ctx, err = configureLogging(ctx, config.regConfig)
    if err != nil {
        return nil, fmt.Errorf("error configuring logger: %v", err)
    }

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
    server := &http.Server{
        Handler: handler,
    }

    return &PocketRegistry{
        app:    app,
        config: config,
        server: server,
    }, nil
}

// A Registry represents a complete instance of the registry.
// TODO(aaronl): It might make sense for Registry to become an interface.
type PocketRegistry struct {
    config        *PocketRegistryConfig
    app           *handlers.App
    server 		  *http.Server
    listener      io.Closer
}

// ListenAndServe runs the registry's HTTP server.
func (r *PocketRegistry) ListenAndServe() error {
    config := r.config

    ln, err := listener.NewListener(config.regConfig.HTTP.Net, config.regConfig.HTTP.Addr)
    if err != nil {
        return err
    }

	ln = tls.NewListener(ln, config.tlsConfig)
	context.GetLogger(r.app).Infof("listening on %v, tls", ln.Addr())
    return r.server.Serve(ln)
}

// ListenAndServe runs the registry's HTTP server.
func (r *PocketRegistry) Start() (error) {
    config := r.config

    ln, err := listener.NewListener(config.regConfig.HTTP.Net, config.regConfig.HTTP.Addr)
    if err != nil {
        return err
    }

    ln = tls.NewListener(ln, config.tlsConfig)
    context.GetLogger(r.app).Infof("listening on %v, tls", ln.Addr())
    go func() {
        var err = r.server.Serve(ln)
        if err != nil {
            log.Println("HTTP Server Error - ", err)
        }
    }()

    r.listener = ln
    return nil
}

func (r *PocketRegistry) Close() error {
    if r.listener == nil {
        return nil
    }
    ln := r.listener
    r.listener = nil
    return ln.Close()
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

// configureLogging prepares the context with a logger using the
// configuration.
func configureLogging(ctx context.Context, config *configuration.Configuration) (context.Context, error) {
    if config.Log.Level == "" && config.Log.Formatter == "" {
        // If no config for logging is set, fallback to deprecated "Loglevel".
        log.SetLevel(logLevel(config.Loglevel))
        ctx = context.WithLogger(ctx, context.GetLogger(ctx))
        return ctx, nil
    }

    log.SetLevel(logLevel(config.Log.Level))

    formatter := config.Log.Formatter
    if formatter == "" {
        formatter = "text" // default formatter
    }

    switch formatter {
    case "json":
        log.SetFormatter(&log.JSONFormatter{
            TimestampFormat: time.RFC3339Nano,
        })
    case "text":
        log.SetFormatter(&log.TextFormatter{
            TimestampFormat: time.RFC3339Nano,
        })
/*
    case "logstash":
        log.SetFormatter(&logstash.LogstashFormatter{
            TimestampFormat: time.RFC3339Nano,
        })
*/
    default:
        // just let the library use default on empty string.
        if config.Log.Formatter != "" {
            return ctx, fmt.Errorf("unsupported logging formatter: %q", config.Log.Formatter)
        }
    }

    if config.Log.Formatter != "" {
        log.Debugf("using %q logging formatter", config.Log.Formatter)
    }

    if len(config.Log.Fields) > 0 {
        // build up the static fields, if present.
        var fields []interface{}
        for k := range config.Log.Fields {
            fields = append(fields, k)
        }

        ctx = context.WithValues(ctx, config.Log.Fields)
        ctx = context.WithLogger(ctx, context.GetLogger(ctx, fields...))
    }

    return ctx, nil
}

func logLevel(level configuration.Loglevel) log.Level {
    l, err := log.ParseLevel(string(level))
    if err != nil {
        l = log.InfoLevel
        log.Warnf("error parsing level %q: %v, using %q    ", level, err, l)
    }

    return l
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

func resolveConfiguration(args []string) (*configuration.Configuration, error) {
    var configurationPath string

    if len(args) > 0 {
        configurationPath = args[0]
    } else if os.Getenv("REGISTRY_CONFIGURATION_PATH") != "" {
        configurationPath = os.Getenv("REGISTRY_CONFIGURATION_PATH")
    }

    if configurationPath == "" {
        return nil, fmt.Errorf("configuration path unspecified")
    }

    fp, err := os.Open(configurationPath)
    if err != nil {
        return nil, err
    }

    defer fp.Close()

    config, err := configuration.Parse(fp)
    if err != nil {
        return nil, fmt.Errorf("error parsing %s: %v", configurationPath, err)
    }

    return config, nil
}

func nextProtos(config *configuration.Configuration) []string {
    switch config.HTTP.HTTP2.Disabled {
    case true:
        return []string{"http/1.1"}
    default:
        return []string{"h2", "http/1.1"}
    }
}
