package teleport

import (
    "os"
    "fmt"
    "io/ioutil"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/trace"
)

func onStart(config *service.Config) error {
    srv, err := service.NewTeleport(config)
    if err != nil {
        return err
        //return trace.Wrap(err, "initializing teleport")
    }

    if err := srv.Start(); err != nil {
        return err
        //return trace.Wrap(err, "starting teleport")
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

func StartTeleport(debug bool) error {
    cfg := makePocketTeleportConfig()
    if debug {
        cfg.Console = ioutil.Discard
        utils.InitLoggerDebug()
        trace.SetDebug(true)
        log.Info("Teleport DEBUG output configured")
    } else {
        utils.InitLoggerCLI()
        log.Info("Teleport NORMAL cli output configured")
    }
    return onStart(cfg)
}