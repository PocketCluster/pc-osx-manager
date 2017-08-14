package main

import (
    "encoding/json"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/event/route/routepath"
    "github.com/pkg/errors"
)

func initRoutePathService() {

    theApp.GET(routepath.RpathSystemReadiness(), func(_, path, _ string) error {
        // inline json marshalling
        data, err := json.Marshal(map[string]map[string]interface{}{
            "cpu": {
                "status": true,
            },
            "mem": {
                "status": false,
                "error": "Insufficient memory is present in system",
            },
            "hdd": {
                "status": false,
                "error": "Insufficient hdd space is present in system",
            },
            "net": {
                "status": true,
            },
            "vbox": {
                "status": true,
            },
        })
        if err != nil {
            log.Debugf(err.Error())
            return errors.WithStack(err)
        }
        err = FeedResponseForGet(path, string(data))
        if err != nil {
            return errors.WithStack(err)
        }
        return nil
    })

    theApp.GET(routepath.RpathAppExpired(), func(_, path, _ string) error {
        return nil
    })
    theApp.GET(routepath.RpathSystemIsFirstRun(), func(_, path, _ string) error {
        return nil
    })
    theApp.GET(routepath.RpathCmdServiceStart(), func(_, _, _ string) error {
        return nil
    })

}