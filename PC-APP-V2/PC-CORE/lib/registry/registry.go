package registry

import (
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/docker/distribution/registry"
    "github.com/docker/distribution/context"
    "github.com/docker/distribution/version"
    "github.com/docker/distribution/configuration"
    "github.com/docker/distribution/registry/storage/driver/factory"
    _ "github.com/docker/distribution/configuration"
    _ "github.com/docker/distribution/registry/storage/driver/filesystem"
    "github.com/docker/libtrust"
    "github.com/docker/distribution/registry/storage"

    log "github.com/Sirupsen/logrus"
    logstash "github.com/bshuster-repo/logrus-logstash-hook"
)

func logLevel(level configuration.Loglevel) log.Level {
    l, err := log.ParseLevel(string(level))
    if err != nil {
        l = log.InfoLevel
        log.Warnf("error parsing level %q: %v, using %q	", level, err, l)
    }

    return l
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
    case "logstash":
        log.SetFormatter(&logstash.LogstashFormatter{
            TimestampFormat: time.RFC3339Nano,
        })
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

func ParseYamlConfig(configPath string) (*configuration.Configuration, error) {
    fp, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }

    defer fp.Close()

    config, err := configuration.Parse(fp)
    if err != nil {
        return nil, fmt.Errorf("error parsing %s: %v", configPath , err)
    }
    return config, err
}

func NewRegistryConfig() *configuration.Configuration {
    var(
        // logging options
        accessLog = struct {
            Disabled        bool                        `yaml:"disabled,omitempty"`
        } {
            Disabled:       false,
        }
        log = struct {
            AccessLog struct {
                Disabled    bool                        `yaml:"disabled,omitempty"`
            }                                           `yaml:"accesslog,omitempty"`
            Level           configuration.Loglevel      `yaml:"level"`
            Formatter       string                      `yaml:"formatter,omitempty"`
            Fields          map[string]interface{}      `yaml:"fields,omitempty"`
            Hooks           []configuration.LogHook     `yaml:"hooks,omitempty"`
        }{
            AccessLog:      accessLog,
            Level:          configuration.Loglevel("debug"),
            Formatter:      "text",
            Fields:         map[string]interface{} {
                "service": "registry",
                "environment": "pc-master",
            },
            Hooks:          nil,
        }

        // storage options
        storage = configuration.Storage {
            "cache": configuration.Parameters {
                "blobdescriptor": "inmemory",
            },
            "filesystem": configuration.Parameters {
                "rootdirectory": "/Users/almightykim/Workspace/DKIMG/REGISTRY/data",
                "maxthreads": 32,
            },
            "maintenance": configuration.Parameters {
                "readonly": map[interface{}]interface{} {
                    "enabled": false,
                },
                "uploadpurging": map[interface{}]interface{} {
                    "enabled": true,
                    "age": "24h",
                    "interval": "3h",
                    "dryrun": false,
                },
            },
        }

        // http connection options
        letsEncrypt = struct {
            CacheFile        string     `yaml:"cachefile,omitempty"`
            Email            string     `yaml:"email,omitempty"`
        } {}
        tls = struct {
            Certificate      string     `yaml:"certificate,omitempty"`
            Key              string     `yaml:"key,omitempty"`
            ClientCAs        []string   `yaml:"clientcas,omitempty"`
            LetsEncrypt      struct {
                CacheFile    string     `yaml:"cachefile,omitempty"`
                Email string            `yaml:"email,omitempty"`
            }                           `yaml:"letsencrypt,omitempty"`
        } {
            Certificate:     "/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.cert",
            Key:             "/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.key",
            ClientCAs:       nil,
            LetsEncrypt:     letsEncrypt,
        }
        debug = struct {
            Addr string                 `yaml:"addr,omitempty"`
        } {
            "",
        }
        http2 = struct {
            Disabled        bool        `yaml:"disabled,omitempty"`
        } {
            Disabled:       false,
        }

        // HTTP contains configuration parameters for the registry's http
        // interface.
        http = struct {
            Addr            string      `yaml:"addr,omitempty"`
            Net             string      `yaml:"net,omitempty"`
            Host            string      `yaml:"host,omitempty"`
            Prefix          string      `yaml:"prefix,omitempty"`
            Secret          string      `yaml:"secret,omitempty"`
            RelativeURLs    bool        `yaml:"relativeurls,omitempty"`
            TLS struct {
                Certificate string      `yaml:"certificate,omitempty"`
                Key         string      `yaml:"key,omitempty"`
                ClientCAs   []string    `yaml:"clientcas,omitempty"`
                LetsEncrypt struct {
                    CacheFile string    `yaml:"cachefile,omitempty"`
                    Email   string      `yaml:"email,omitempty"`
                }                       `yaml:"letsencrypt,omitempty"`
            }                           `yaml:"tls,omitempty"`
            Headers http.Header         `yaml:"headers,omitempty"`
            Debug struct {
                Addr        string      `yaml:"addr,omitempty"`
            }                           `yaml:"debug,omitempty"`
            HTTP2 struct {
                Disabled    bool        `yaml:"disabled,omitempty"`
            }                           `yaml:"http2,omitempty"`
        } {
            Addr:           "0.0.0.0:5000",
            Net:            "tcp",
            Host:           "",
            Prefix:         "",
            Secret:         "",
            RelativeURLs:   false,
            TLS:            tls,
            Headers:        http.Header {
                "X-Content-Type-Options": []string{"nosniff"},
            },
            Debug:          debug,
            HTTP2:          http2,
        }

        // Notifications specifies configuration about various endpoint to which
        // registry events are dispatched.
        notifications = configuration.Notifications {
            Endpoints:  nil,
        }

        // Redis configures the redis pool available to the registry webapp.
        redis = struct {
            Addr           string       `yaml:"addr,omitempty"`
            Password       string       `yaml:"password,omitempty"`
            DB             int          `yaml:"db,omitempty"`
            DialTimeout    time.Duration `yaml:"dialtimeout,omitempty"`
            ReadTimeout    time.Duration `yaml:"readtimeout,omitempty"`
            WriteTimeout   time.Duration `yaml:"writetimeout,omitempty"`
            Pool struct {
                MaxIdle    int          `yaml:"maxidle,omitempty"`
                MaxActive  int          `yaml:"maxactive,omitempty"`
                IdleTimeout time.Duration `yaml:"idletimeout,omitempty"`
            }                           `yaml:"pool,omitempty"`
        } {}

        // health check
        health = configuration.Health {
            FileCheckers:  nil,
            HTTPCheckers:  nil,
            TCPCheckers:   nil,
        }

        // Compatibility is used for configurations of working with older or deprecated features.
        compatibility = struct {
            Schema1 struct {
                TrustKey string `yaml:"signingkeyfile,omitempty"`
            } `yaml:"schema1,omitempty"`
        } {}

        // Validation configures validation options for the registry.
        validation = struct {
            Enabled        bool         `yaml:"enabled,omitempty"`
            Manifests struct {
                URLs struct {
                    Allow  []string     `yaml:"allow,omitempty"`
                    Deny   []string     `yaml:"deny,omitempty"`
                }                       `yaml:"urls,omitempty"`
            }                           `yaml:"manifests,omitempty"`
        } {}

        // Policy configures registry policy options.
        policy = struct {
            Repository struct {
                Classes []string        `yaml:"classes"`
            }                           `yaml:"repository,omitempty"`
        } {}
    )

    return &configuration.Configuration {
        Version:        configuration.MajorMinorVersion(0, 1),
        Log:            log,
        Loglevel:       configuration.Loglevel("info"),
        Storage:        storage,
        Auth:           nil,
        Middleware:     nil,
        Reporting:      configuration.Reporting {},
        HTTP:           http,
        Notifications:  notifications,
        Redis:          redis,
        Health:         health,
        Proxy:          configuration.Proxy{},
        Compatibility:  compatibility,
        Validation:     validation,
        Policy:         policy,
    }
}

func Serve(config *configuration.Configuration) {
    // setup context
    ctx := context.WithVersion(context.Background(), version.Version)

    if config.HTTP.Debug.Addr != "" {
        go func(addr string) {
            log.Infof("debug server listening %v", addr)
            if err := http.ListenAndServe(addr, nil); err != nil {
                log.Fatalf("error listening on debug interface: %v", err)
            }
        }(config.HTTP.Debug.Addr)
    }

    registry, err := registry.NewRegistry(ctx, config)
    if err != nil {
        log.Fatalln(err)
    }

    if err = registry.ListenAndServe(); err != nil {
        log.Fatalln(err)
    }
}

func GarbageCollection(config *configuration.Configuration) {
    driver, err := factory.Create(config.Storage.Type(), config.Storage.Parameters())
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to construct %s driver: %v", config.Storage.Type(), err)
        os.Exit(1)
    }

    ctx := context.Background()
    ctx, err = configureLogging(ctx, config)
    if err != nil {
        fmt.Fprintf(os.Stderr, "unable to configure logging with config: %s", err)
        os.Exit(1)
    }

    k, err := libtrust.GenerateECP256PrivateKey()
    if err != nil {
        fmt.Fprint(os.Stderr, err)
        os.Exit(1)
    }

    registry, err := storage.NewRegistry(ctx, driver, storage.Schema1SigningKey(k))
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to construct registry: %v", err)
        os.Exit(1)
    }

    err = storage.MarkAndSweep(ctx, driver, registry, false)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to garbage collect: %v", err)
        os.Exit(1)
    }}

