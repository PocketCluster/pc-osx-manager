package list

import (
    "encoding/json"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
)

// read available package list from backend
func InitRouthPathListInstalled(appLife route.Router, feeder route.ResponseFeeder) {
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

            pkgList = []map[string]interface{}{}
            pkgs    = []*model.Package{}
        )

        // update package doesn't return error when there is packages to update.
        pkgs, err := model.FindPackage("", "")
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access package list"))
        }
        for i, _ := range pkgs {
            pkgList = append(pkgList, map[string]interface{} {
                "package-id" : pkgs[i].PkgID,
                "description": pkgs[i].Description,
                "installed": true,
            })
        }
        data, err := json.Marshal(route.ReponseMessage{
            "package-installed": {
                "status": true,
                "list":   pkgList,
            },
        })
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access package list"))
        }
        err = feeder.FeedResponseForGet(rpath, string(data))
        if err != nil {
            return errors.WithStack(err)
        }
        return nil
    })
}
