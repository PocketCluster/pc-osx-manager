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
        app service.AppSupervisor = service.NewAppSupervisor()
    )
    log.SetLevel(log.DebugLevel)
    crcontext.SharedCoreContext()

    err = initVboxCoreReportService(app)
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