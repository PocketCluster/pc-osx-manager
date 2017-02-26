package registry

import (
    "os"
    "fmt"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/Sirupsen/logrus/formatters/logstash"
    "github.com/docker/distribution/configuration"
    "github.com/docker/distribution/context"
)

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

func logLevel(level configuration.Loglevel) log.Level {
    l, err := log.ParseLevel(string(level))
    if err != nil {
        l = log.InfoLevel
        log.Warnf("error parsing level %q: %v, using %q	", level, err, l)
    }

    return l
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

