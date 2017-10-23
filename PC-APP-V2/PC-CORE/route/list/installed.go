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
        )

        // update package doesn't return error when there is packages to update.
        recs, err := model.AllRecords()
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access installed packages"))
        }
        for i, _ := range recs {
            r := recs[i]

            if pkgs, err := model.FindPackage("pkg_id = ?", r.PkgID); err == nil && len(pkgs) != 0 {
                p := pkgs[0]

                pkgList = append(pkgList, map[string]interface{} {
                    "package-id" : p.PkgID,
                    "description": p.Description,
                    "installed":   true,
                    "menu-name":   p.MenuName,
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
