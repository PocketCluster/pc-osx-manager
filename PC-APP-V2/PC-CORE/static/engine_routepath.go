package main

import (
    //log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/event/route/routepath"
)

func initRoutePathService() {

    theApp.GET(routepath.RpathSystemReadiness, func(_ string) error {
        FeedSend("[ROUTE] system is ready to run")
        return nil
    })

    theApp.GET(routepath.RpathAppExpired, func(_ string) error {
        FeedSend("[ROUTE] app is not expired")
        return nil
    })

    theApp.GET(routepath.RpathSystemIsFirstRun, func(_ string) error {
        FeedSend("[ROUTE] this is not the first run")
        return nil
    })

    theApp.GET(routepath.RpathCmdServiceStart, func(_ string) error {

        return nil
    })

}