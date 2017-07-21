package main

import (
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/service"
)

func main() {
    var (
        err error = nil
        app service.AppSupervisor = nil
    )
    log.SetLevel(log.DebugLevel)

    // TODO check user and reject if not root
    // TODO check if this exceed Date limit set for raspberry pi 3

    // initialize slave context
    slcontext.SharedSlaveContext()
    app = service.NewAppSupervisor()

    // dhcp listner
    err = initDhcpListner(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // search service
    _, err = mcast.NewSearchCaster(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // beacon service
    _, err = ucast.NewBeaconAgent(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // agent service
    err = initAgentService(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // teleport management
    err = initTeleportNodeService(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // DNS service
    err = initDNSService(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // application
    err = app.Start()
    if err != nil {
        log.Panic(errors.WithStack(err))
    }
    app.Wait()
}

