package pkglist

import (
    "encoding/json"
    "fmt"
    "strings"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/rasker"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
)

// read available package list from backend
func InitRouthPathListInstalled(appLife rasker.RouteTasker, feeder route.ResponseFeeder) {
    // get the list of available packages
    appLife.GET(routepath.RpathPackageListInstalled(), func(_, rpath, _ string) error {
        var (
            feedError = func(irr error) error {
                data, frr := json.Marshal(route.ReponseMessage{
                    "package-installed": {
                        "status": false,
                        "error" : irr.Error(),
                    },
                })
                if frr != nil {
                    log.Debugf(frr.Error())
                }
                frr = feeder.FeedResponseForGet(rpath, string(data))
                if frr != nil {
                    log.Debugf(frr.Error())
                }
                return errors.WithStack(irr)
            }

            coreAddrC = make(chan service.Event)
            pkgList []map[string]interface{} = nil
        )

        // update package doesn't return error when there is packages to update.
        recs, err := model.AllRecords()
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access installed packages"))
        }

        // acquire core external address
        if err = appLife.BindDiscreteEvent(ivent.IventReportCoreAddrResult, coreAddrC); err != nil {
            return feedError(errors.WithMessage(err, "Unable to install package due to invalid node status"))
        }
        appLife.BroadcastEvent(service.Event{Name:ivent.IventReportCoreAddrRequest})
        cr := <- coreAddrC
        appLife.UntieDiscreteEvent(ivent.IventReportCoreAddrResult)
        coreAddr, ok := cr.Payload.(string)
        if !ok {
            return feedError(errors.WithMessage(cr.Payload.(error), "Unable to acquire core node address"))
        }

        for i, _ := range recs {
            r := recs[i]

            if pkgs, err := model.FindPackage("pkg_id = ?", r.PkgID); err == nil && len(pkgs) != 0 {
                var (
                    aPkg                         = pkgs[0]
                    consoles []map[string]string = nil
                )
                for _, prt := range strings.Split(aPkg.WebPorts, "|") {
                    pds := strings.Split(prt, "!")
                    consoles = append(
                        consoles,
                        map[string]string {
                            "address": fmt.Sprintf("%s%s:%s", pds[0], coreAddr, pds[1]),
                            "name":    pds[2],
                        })
                }

                pkgList = append(pkgList, map[string]interface{} {
                    "package-id" : aPkg.PkgID,
                    "description": aPkg.Description,
                    "installed":   true,
                    "menu-name":   aPkg.MenuName,
                    "consoles":    consoles,
                })
            }
        }
        data, err := json.Marshal(route.ReponseMessage{
            "package-installed": {
                "status": true,
                "list":   pkgList,
            },
        })
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access installed packages"))
        }
        err = feeder.FeedResponseForGet(rpath, string(data))
        if err != nil {
            return errors.WithStack(err)
        }
        return nil
    })
}
