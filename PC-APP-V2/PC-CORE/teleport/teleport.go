package teleport

import (
    "os"
    "fmt"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/trace"
)

func onStart(config *service.Config) error {
    srv, err := service.NewTeleport(config)
    if err != nil {
        return trace.Wrap(err, "initializing teleport")
    }
    if err := srv.Start(); err != nil {
        return trace.Wrap(err, "starting teleport")
    }

    // create the pid file
    if config.PIDFile != "" {
        f, err := os.OpenFile(config.PIDFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
        if err != nil {
            return trace.Wrap(err, "failed to create the PID file")
        }
        fmt.Fprintf(f, "%v", os.Getpid())
        defer f.Close()
    }
    srv.Wait()
    return nil
}

func StartTeleport() error {
    return onStart(makePocketTeleportConfig())
}