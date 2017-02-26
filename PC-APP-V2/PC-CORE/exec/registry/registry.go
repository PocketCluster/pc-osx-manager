package registry

import (
    "fmt"
    "net/http"
    "os"

    log "github.com/Sirupsen/logrus"

    "github.com/docker/libtrust"
    _ "github.com/docker/distribution/configuration"
    "github.com/docker/distribution/registry"
    "github.com/docker/distribution/context"
    "github.com/docker/distribution/version"
    "github.com/docker/distribution/registry/storage/driver/factory"
    "github.com/docker/distribution/registry/storage"
)

func Serve(args []string) {
    // setup context
    ctx := context.WithVersion(context.Background(), version.Version)

    config, err := resolveConfiguration(args)
    if err != nil {
        fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
        os.Exit(1)
    }

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

func GarbageCollection(args []string) {
    config, err := resolveConfiguration(args)
    if err != nil {
        fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
        return
    }

    driver, err := factory.Create(config.Storage.Type(), config.Storage.Parameters())
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to construct %s driver: %v", config.Storage.Type(), err)
        return
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

    registry, err := storage.NewRegistry(ctx, driver, storage.DisableSchema1Signatures, storage.Schema1SigningKey(k))
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to construct registry: %v", err)
        os.Exit(1)
    }

    var dryRun bool = false
    err = storage.MarkAndSweep(ctx, driver, registry, dryRun)
    if err != nil {
        fmt.Fprintf(os.Stderr, "failed to garbage collect: %v", err)
        os.Exit(1)
    }
}
