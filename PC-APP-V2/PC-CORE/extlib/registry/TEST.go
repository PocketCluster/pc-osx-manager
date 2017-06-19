package registry

import (
    "crypto/tls"
    "fmt"
    "net"
    "net/http"
    "os"
    "sync"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/docker/distribution/configuration"
    "github.com/docker/distribution/context"
    "github.com/docker/distribution/registry"
    "github.com/docker/distribution/registry/listener"
    "github.com/docker/distribution/version"

    "gopkg.in/tylerb/graceful.v1"
)

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

func NewRegistrySampleConfig() *configuration.Configuration {
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
        httpTLS = struct {
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
            TLS:            httpTLS,
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


// ListenAndServe runs the registry's HTTP server.
func (r *PocketRegistry) StartOnWaitGroup(wg *sync.WaitGroup) (error) {
    config := r.config

    ln, err := listener.NewListener(config.regConfig.HTTP.Net, config.regConfig.HTTP.Addr)
    if err != nil {
        return err
    }

    // TODO : No HTTP secret provided - generated random secret. This may cause problems with uploads if multiple registries are behind a load-balancer. To provide a shared secret, fill in http.secret in the configuration file or set the REGISTRY_HTTP_SECRET environment variable.
    ln = tls.NewListener(ln, config.tlsConfig)
    context.GetLogger(r.app).Infof("listening on %v, tls", ln.Addr())
    go func(w *sync.WaitGroup, srv *graceful.Server, l net.Listener) {
        defer w.Done()

        var err = srv.Serve(l)
        if err != nil {
            log.Println("HTTP Server Error - ", err)
        }
    }(wg, r.server, ln)

    r.listener = ln
    return nil
}
