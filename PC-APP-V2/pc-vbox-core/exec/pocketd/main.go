package main

import (
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/service"
    "github.com/stkim1/pc-vbox-core/crcontext"
)

func main() {
    var (
        err error = nil
        app service.AppSupervisor = nil
    )
    log.SetLevel(log.DebugLevel)

    // TODO check user and reject if not root
    // TODO check if this exceed Date limit set


    crcontext.SharedCoreContext()
    app = service.NewAppSupervisor()

    // dhcp listner
    err = initDhcpListner(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // vbox reporter
    err = initVboxCoreReportService(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // DNS Service
    err = initDNSService(app)
    if err != nil {
        log.Panic(errors.WithStack(err))
    }

    // teleport management
    err = initTeleportNodeService(app)
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