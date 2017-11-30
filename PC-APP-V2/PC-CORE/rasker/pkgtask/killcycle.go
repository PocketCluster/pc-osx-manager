package pkgtask

import (
    "encoding/json"
    "fmt"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/rasker"
    "github.com/stkim1/pc-core/service"
)

// kill a package
func InitPackageKillCycle(appLife rasker.RouteTasker, feeder route.ResponseFeeder) error {
    return appLife.POST(routepath.RpathPackageKill(), func(_, rpath, payload string) error {
        // 1. parse input package id
        var (
            pkgID string = ""
            err error = nil
        )
        err = json.Unmarshal([]byte(payload), &struct {
            PkgID *string `json:"pkg-id"`
        }{&pkgID})
        if err != nil || len(pkgID) == 0 {
            return feedError(feeder, rpath, fbPackageKill, pkgID, errors.WithMessage(err, "unable to specify package id"))
        }

        // broadcast kill signal
        appLife.BroadcastEvent(service.Event{Name:fmt.Sprintf("%s%s", iventPackageKillPrefix, pkgID)})
        data, err := json.Marshal(route.ReponseMessage{
            "package-kill": {
                "status": true,
                "pkg-id": pkgID,
            },
        })
        if err != nil {
            return errors.WithStack(err)
        }
        return feeder.FeedResponseForPost(rpath, string(data))
    })
}