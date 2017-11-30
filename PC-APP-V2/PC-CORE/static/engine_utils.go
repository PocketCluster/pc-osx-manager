package main

import (
    "encoding/json"
    ctx "context"

    tclient "github.com/gravitational/teleport/lib/client"
    tervice "github.com/gravitational/teleport/lib/service"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/extlib/pcssh/sshadmin"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
)

// report context initialization prep status
func reportContextInit(appLife *appMainLife, feeder route.ResponseFeeder) error {
    data, err := json.Marshal(route.ReponseMessage{
        "sys-context-init": {
            "status": true,
        },
    })
    if err != nil {
        return errors.WithStack(err)
    }
    err = feeder.FeedResponseForGet(routepath.RpathSystemContextInit(), string(data))
    return errors.WithStack(err)
}

// report network initialization prep status
func reportNetworkInit(appLife *appMainLife, feeder route.ResponseFeeder) error {
    var (
        message route.ReponseMessage
    )
    _, err := context.SharedHostContext().HostPrimaryInterfaceShortName()
    if err != nil {
        message = route.ReponseMessage{
            "sys-network-init": {
                "status": false,
                "error": err.Error(),
            },
        }
    } else {
        message = route.ReponseMessage{
            "sys-network-init": {
                "status": true,
            },
        }
    }

    data, err := json.Marshal(message)
    if err != nil {
        return errors.WithStack(err)
    }
    err = feeder.FeedResponseForGet(routepath.RpathSystemNetworkInit(), string(data))
    return errors.WithStack(err)
}

// setup base users
func setupBaseUsers(pcsshCfg *tervice.PocketConfig) error {
    cli, err := sshadmin.OpenAdminClientWithAuthService(pcsshCfg)
    if err != nil {
        return errors.WithStack(err)
    }
    roots, err := model.FindUserMetaWithLogin("root")
    if err != nil {
        return errors.WithStack(err)
    }
    err = sshadmin.CreateTeleportUser(cli, "root", roots[0].Password)
    if err != nil {
        return errors.WithStack(err)
    }
    uname, err := context.SharedHostContext().LoginUserName()
    if err != nil {
        return errors.WithStack(err)
    }
    lusers, err := model.FindUserMetaWithLogin(uname)
    if err != nil {
        return errors.WithStack(err)
    }
    err = sshadmin.CreateTeleportUser(cli, uname, lusers[0].Password)
    return errors.WithStack(err)
}

// shutdown nodes and stop services
func shutdownSlaveNodes(appLife *appMainLife, pcsshCfg *tervice.PocketConfig) error {
    // find root user
    roots, err := model.FindUserMetaWithLogin("root")
    if err != nil {
        return errors.WithStack(err)
    }

    // get the node list report
    reportC := make(chan service.Event)
    err = appLife.BindDiscreteEvent(ivent.IventReportLiveNodesResult, reportC)
    if err != nil {
        return errors.WithStack(err)
    }
    // ask node list
    appLife.BroadcastEvent(service.Event{Name:ivent.IventReportLiveNodesRequest})
    nr := <- reportC
    appLife.UntieDiscreteEvent(ivent.IventReportLiveNodesResult)
    nlist, ok := nr.Payload.([]string)
    if !ok {
        return errors.Errorf("unable to access proper node list")
    }

    // execute shutdown command
    for _, node := range nlist {
        if node == "pc-core" {
            continue
        }
        clt, err := tclient.MakeNewClient(pcsshCfg, "root", node)
        if err != nil {
            log.Error(err.Error())
            continue
        }
        err = clt.APISSH(ctx.TODO(), []string{"/sbin/shutdown", "-h", "now"}, roots[0].Password, false)
        if err != nil {
            log.Error(err.Error())
            continue
        }
        clt.Logout()
    }

    return nil
}

// stop baseservice
func stopHealthMonitor(appLife *appMainLife) error {
    // stop health and wait
    resultC := make(chan service.Event)
    err := appLife.BindDiscreteEvent(ivent.IventMonitorStopResult, resultC)
    if err != nil {
        return errors.WithMessage(err,"[LIFE] unable to stop monitoring...")
    }
    appLife.BroadcastEvent(service.Event{Name:ivent.IventMonitorStopRequest})
    <- resultC
    return appLife.UntieDiscreteEvent(ivent.IventMonitorStopResult)
}

func feedShutdownReadySignal(feeder route.ResponseFeeder) error {
    // we send it's ok to quit signal to frontend
    data, err := json.Marshal(route.ReponseMessage{
        "app-shutdown-ready": {
            "status": true,
        },
    })
    if err != nil {
        return errors.WithStack(err)
    }
    err = feeder.FeedResponseForGet(routepath.RpathAppPrepShutdown(), string(data))
    return errors.WithStack(err)
}
