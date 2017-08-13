package main

import (
    //log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/event/route/routepath"
)

func initRoutePathService() {

    theApp.GET(routepath.RpathSystemReadiness(), func(_, path, _ string) error {
        FeedResponseForGet(path,"[ROUTE] system is ready to run")
        /* inline json marshalling
        data, err := json.Marshal(map[string]string{
            FeedType:   "api-feed",
            FeedResult: "api-success",
            FeedMessage: feed,
        })
        */
        return nil
    })

    theApp.GET(routepath.RpathAppExpired(), func(_, path, _ string) error {
        FeedResponseForGet(path,"[ROUTE] app is not expired")
        return nil
    })

    theApp.GET(routepath.RpathSystemIsFirstRun(), func(_, path, _ string) error {
        FeedResponseForGet(path,"[ROUTE] this is not the first run")
        return nil
    })

    theApp.GET(routepath.RpathCmdServiceStart(), func(_, _, _ string) error {

        return nil
    })

}