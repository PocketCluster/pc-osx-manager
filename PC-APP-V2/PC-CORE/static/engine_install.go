package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    //"github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/route/routepath"
    "github.com/stkim1/pc-core/model"
)

func initInstallRoutePath() {

    // get the list of available packages
    theApp.GET(routepath.RpathPackageList(), func(_, path, _ string) error {
        var (
            feedError = func(errMsg string) error {
                data, err := json.Marshal(ReponseMessage{
                    "package-list": {
                        "status": false,
                        "error" : errMsg,
                    },
                })
                if err != nil {
                    log.Debugf(err.Error())
                    return errors.WithStack(err)
                }
                err = FeedResponseForGet(path, string(data))
                return errors.WithStack(err)
            }

            pkgList = []map[string]interface{}{}
        )

        req, err :=  http.NewRequest("GET", "https://api.pocketcluster.io/service/v014/package/list", nil)
        if err != nil {
            log.Debugf(errors.WithStack(err).Error())
            return feedError("Unable to access package list. Reason : " + errors.WithStack(err).Error())
        }
        //req.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
        req.Header.Add("User-Agent", "PocketCluster/0.1.4 (OSX)")
        req.Header.Set("Content-Type", " application/json; charset=utf-8")
        client := &http.Client{
            Timeout: 10 * time.Second,
        }
        resp, err := client.Do(req)
        if err != nil {
            log.Debugf(errors.WithStack(err).Error())
            return feedError("Unable to access package list. Reason : " + errors.WithStack(err).Error())
        }
        defer resp.Body.Close()

        if resp.StatusCode != 200 {
            return feedError(errors.Errorf("Service unavailable. Status : %d", resp.StatusCode).Error())
        }
        var pkgs = []*model.Package{}
        err = json.NewDecoder(resp.Body).Decode(&pkgs)
        if err != nil {
            log.Debugf(errors.WithStack(err).Error())
            return feedError("Unable to access package list. Reason : " + errors.WithStack(err).Error())
        }
        if len(pkgs) == 0 {
            return feedError("No package avaiable. Contact us at Slack channel.")
        } else {
            // update package doesn't return error when there is packages to update.
            model.UpdatePackages(pkgs)
        }

        for i, _ := range pkgs {
            pkgList = append(pkgList, map[string]interface{} {
                "package-id" : pkgs[i].PkgID,
                "description": pkgs[i].Description,
                "installed": false,
            })
        }
        data, err := json.Marshal(ReponseMessage{
            "package-list": {
                "status": true,
                "list":   pkgList,
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

    // install a package
    theApp.POST(routepath.RpathPackageInstall(), func(_, path, payload string) error {
        var (
            feedError = func(errMsg string) error {
                data, err := json.Marshal(ReponseMessage{
                    "package-install": {
                        "status": false,
                        "error" : errMsg,
                    },
                })
                if err != nil {
                    log.Debugf(err.Error())
                    return errors.WithStack(err)
                }
                err = FeedResponseForPost(path, string(data))
                return errors.WithStack(err)
            }
            pkgID string = ""
        )
        err := json.Unmarshal([]byte(payload), &struct {
            PkgID *string `json:"pkg-id"`
        }{&pkgID})
        if err != nil {
            feedError(errors.WithStack(err).Error())
            return err
        }

        pkgs, _ := model.FindPackage("pkg_id = ?", pkgID)
        if len(pkgs) == 0 {
            errMsg := fmt.Sprintf("selected package %s is not available", pkgID)
            feedError(errMsg)
            return errors.Errorf(errMsg)
        }

        feedError("This is test error")
        return nil
    })
}