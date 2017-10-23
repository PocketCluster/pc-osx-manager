package list

import (
    "encoding/json"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/utils/apireq"
)

// read available package list from backend
func InitRouthPathListAvailable(appLife route.Router, feeder route.ResponseFeeder) {
    // get the list of available packages
    appLife.GET(routepath.RpathPackageListAvailable(), func(_, rpath, _ string) error {
        var (
            feedError = func(irr error) error {
                data, frr := json.Marshal(route.ReponseMessage{
                    "package-available": {
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
            pkgs    = []model.Package{}
        )
        req, err := apireq.NewRequest("https://api.pocketcluster.io/service/v014/package/list", false)
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access package list"))
        }
        client := apireq.NewClient(apireq.ConnTimeout, true)
        resp, err := apireq.ReadRequest(req, client)
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access package list"))
        }
        err = json.Unmarshal(resp, &pkgs)
        if err != nil {
            return feedError(errors.WithMessage(err, "Unable to access package list"))
        }
        if len(pkgs) == 0 {
            return feedError(errors.Errorf("No package avaiable. Contact us at Slack channel."))
        } else {
            // update package doesn't return error when there is packages to update.
            model.UpsertPackages(pkgs)
        }

        for i, _ := range pkgs {
            pkgList = append(pkgList, map[string]interface{} {
                "package-id" : pkgs[i].PkgID,
                "description": pkgs[i].Description,
                "installed": false,
            })
        }
        data, err := json.Marshal(route.ReponseMessage{
            "package-available": {
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
