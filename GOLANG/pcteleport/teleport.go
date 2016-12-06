package pcteleport

import (
    "os"
    "fmt"
    "time"
    "io/ioutil"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    "github.com/gravitational/teleport/lib/utils"
    "github.com/gravitational/teleport/lib/config"
    "github.com/gravitational/teleport/lib/services"

    "github.com/stkim1/pcteleport/process"
    "github.com/stkim1/pcteleport/pcconfig"
)

func StartCoreTeleport(debug bool) error {
    cfg := pcconfig.MakeCoreTeleportConfig()
    if debug {
        cfg.Console = ioutil.Discard
        utils.InitLoggerDebug()
        trace.SetDebug(true)
        log.Info("Teleport DEBUG output configured")
    } else {
        utils.InitLoggerCLI()
        log.Info("Teleport NORMAL cli output configured")
    }

    // add static tokens
    for _, token := range []config.StaticToken{"node:d52527f9-b260-41d0-bb5a-e23b0cfe0f8f", "node:c9s93fd9-3333-91d3-9999-c9s93fd98f43"} {
        roles, tokenValue, err := token.Parse()
        if err != nil {
            return trace.Wrap(err)
        }
        cfg.Auth.StaticTokens = append(cfg.Auth.StaticTokens, services.ProvisionToken{Token: tokenValue, Roles: roles, Expires: time.Unix(0, 0)})
    }

    // add temporary token
    //srv, err := service.NewTeleport(cfg)
    srv, err := process.NewCoreTeleport(cfg)
    if err != nil {
        return err
        //return trace.Wrap(err, "initializing teleport")
    }

    if err := srv.Start(); err != nil {
        return err
        //return trace.Wrap(err, "starting teleport")
    }
    // create the pid file
    if cfg.PIDFile != "" {
        f, err := os.OpenFile(cfg.PIDFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
        if err != nil {
            return trace.Wrap(err, "failed to create the PID file")
        }
        fmt.Fprintf(f, "%v", os.Getpid())
        defer f.Close()
    }
    srv.Wait()
    return nil
}
