package pcteleport

import (
    "os"
    "fmt"
    "time"
    "log"

    "github.com/gravitational/trace"
    "github.com/gravitational/teleport/lib/config"
    "github.com/gravitational/teleport/lib/services"

    "github.com/stkim1/pcteleport/process"
    "github.com/stkim1/pcteleport/pcconfig"
)

func StartCoreTeleport(debug bool) error {
    cfg := pcconfig.MakeCoreTeleportConfig(debug)

    // add static tokens
    for _, token := range []config.StaticToken{"node:d52527f9-b260-41d0-bb5a-e23b0cfe0f8f", "node:c9s93fd9-3333-91d3-9999-c9s93fd98f43"} {
        roles, tokenValue, err := token.Parse()
        if err != nil {
            return trace.Wrap(err)
        }
        cfg.Auth.StaticTokens = append(cfg.Auth.StaticTokens, services.ProvisionToken{Token: tokenValue, Roles: roles, Expires: time.Unix(0, 0)})
    }

    // add temporary token
    srv, err := process.NewCoreTeleport(cfg)
    if err != nil {
        return trace.Wrap(err, "initializing teleport")
    }

    if err := srv.Start(); err != nil {
        return trace.Wrap(err, "starting teleport")
    }
    srv.Wait()
    return nil
}

func StartNodeTeleport(authServerAddr, authToken string, debug bool) error {
    cfg, err := pcconfig.MakeNodeTeleportConfig(authServerAddr, authToken, debug)
    if err != nil {
        log.Print(err.Error())
        return trace.Wrap(err, "error in initializing teleport")
    }

    log.Println("Node config created")

    // add temporary token
    srv, err := process.NewNodeTeleport(cfg)
    if err != nil {
        log.Print(err.Error())
        return trace.Wrap(err, "error in initializing teleport")
    }

    if err := srv.Start(); err != nil {
        log.Print(err.Error())
        return trace.Wrap(err, "starting teleport")
    }
    // create the pid file
    if cfg.PIDFile != "" {
        f, err := os.OpenFile(cfg.PIDFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
        if err != nil {
            log.Print(err.Error())
            return trace.Wrap(err, "failed to create the PID file")
        }
        fmt.Fprintf(f, "%v", os.Getpid())
        defer f.Close()
    }
    srv.Wait()
    return nil
}

