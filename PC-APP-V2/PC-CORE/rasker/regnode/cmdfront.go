package regnode

import (
    "encoding/json"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/rasker"
    "github.com/stkim1/pc-core/service"
)

func InitNodeRegisterStop(appLife rasker.RouteTasker,  feeder route.ResponseFeeder) error {
    return appLife.GET(routepath.RpathNodeRegStop(), func(_, rpath, _ string) error {
        // broadcast stop signal
        appLife.BroadcastEvent(service.Event{Name:iventNodeRegisterStop})
        data, err := json.Marshal(route.ReponseMessage{
            "node-reg-stop": {
                "status": true,
            },
        })
        if err != nil {
            return errors.WithStack(err)
        }
        return feeder.FeedResponseForGet(rpath, string(data))
    })
}

func InitNodeRegisterCanidate(appLife rasker.RouteTasker,  feeder route.ResponseFeeder) error {
    return appLife.POST(routepath.RpathNodeRegCandiate(), func(_, rpath, payload string) error {
        var (
            candid = make([]string, 0)
            err error = nil
        )
        err = json.Unmarshal([]byte(payload), &struct {
            Candidates *[]string `json:"candidates"`
        }{&candid})
        if err != nil || len(candid) == 0 {
            return feedPostError(feeder, rpath, "node-reg-candidate", errors.WithMessage(err, "unable to identify candidate nodes"))
        }

        // broadcast stop signal
        appLife.BroadcastEvent(service.Event{
            Name:iventNodeRegisterCandid,
            Payload:candid})
        return nil
    })
}
