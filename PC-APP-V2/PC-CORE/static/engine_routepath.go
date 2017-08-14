package main

import (
    "encoding/json"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/route/routepath"
)

func initRoutePathService() {

    // check if this system is suitable to run
    theApp.GET(routepath.RpathSystemReadiness(), func(_, path, _ string) error {

        cpu := context.SharedHostContext().HostPhysicalCoreCount()
        mem := context.SharedHostContext().HostPhysicalMemorySize()
        _, avail := context.SharedHostContext().HostStorageSpaceStatus()

        // inline json marshalling
        data, err := json.Marshal(map[string]map[string]interface{}{
            "cpu": {
                "status": true,
                "count": cpu,
            },
            "mem": {
                "status": false,
                "size": mem,
                "error": "Insufficient memory is present in system",
            },
            "hdd": {
                "status": false,
                "size": avail,
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

    // check if app is expired
    theApp.GET(routepath.RpathAppExpired(), func(_, path, _ string) error {
        data, err := json.Marshal(map[string]map[string]interface{}{
            "expired": {
                "status": true,
                "error": "expired by 2017/09/11",
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

    // check if user is authenticated
    theApp.GET(routepath.RpathUserAuthed(), func(_, path, _ string) error {
        data, err := json.Marshal(map[string]map[string]interface{}{
            "user-auth": {
                "status": false,
                "error": "please check your invitation code",
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

    // check if this is the first time run
    theApp.GET(routepath.RpathSystemIsFirstRun(), func(_, path, _ string) error {
        data, err := json.Marshal(map[string]map[string]interface{}{
            "user-auth": {
                "status": false,
                "error": "please check your invitation code",
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


    theApp.GET(routepath.RpathCmdServiceStart(), func(_, _, _ string) error {
        return nil
    })

}