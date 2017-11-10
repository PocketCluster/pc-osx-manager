package main

import (
    "encoding/json"
    ctx "context"

    tclient "github.com/gravitational/teleport/lib/client"
    tervice "github.com/gravitational/teleport/lib/service"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/context"
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

// stop baseservice
func stopBaseService(appLife *appMainLife, feeder route.ResponseFeeder) error {
    // binf channel
    resultC := make(chan service.Event)
    err := appLife.BindDiscreteEvent(ivent.IventMonitorStopResult, resultC)
    if err != nil {
        return errors.WithMessage(err,"[LIFE] unable to stop monitoring...")
    }

    // ask node list
    appLife.BroadcastEvent(service.Event{Name:ivent.IventMonitorStopRequest})
    <- resultC
    appLife.UntieDiscreteEvent(ivent.IventMonitorStopResult)
    appLife.StopServices()

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

// shutdown nodes and stop services
func shutdownCluster(appLife *appMainLife, feeder route.ResponseFeeder, pcsshCfg *tervice.PocketConfig) error {
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

    return stopBaseService(appLife, feeder)
}