package main

import (
    //log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/event/route/routepath"
)

func initRoutePathService() {

    theApp.GET(routepath.RpathSystemReadiness(), func(_, _, _ string) error {
        EventFeedPost("/v1/error/message","[ROUTE] system is ready to run")
        /* inline json marshalling
        data, err := json.Marshal(map[string]string{
            FeedType:   "api-feed",
            FeedResult: "api-success",
            FeedMessage: feed,
        })
        */
        return nil
    })

    theApp.GET(routepath.RpathAppExpired(), func(_, _, _ string) error {
        EventFeedPost("/v1/error/message","[ROUTE] app is not expired")
        return nil
    })

    theApp.GET(routepath.RpathSystemIsFirstRun(), func(_, _, _ string) error {
        EventFeedPost("/v1/error/message","[ROUTE] this is not the first run")
        return nil
    })

    theApp.GET(routepath.RpathCmdServiceStart(), func(_, _, _ string) error {

        return nil
    })

}