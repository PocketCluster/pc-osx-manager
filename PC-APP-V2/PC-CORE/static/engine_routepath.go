package main

import (
    "encoding/json"
    "net/http"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/route/routepath"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/vboxglue"
)

type ReponseMessage map[string]map[string]interface{}

func initRoutePathService() {

    // check if this system is suitable to run
    theApp.GET(routepath.RpathSystemReadiness(), func(_, path, _ string) error {
        var (
            syserr, nerr, vlerr, vererr error = nil, nil, nil, nil
            vbox vboxglue.VBoxGlue = nil
            data []byte = nil
            response ReponseMessage = nil
        )

        syserr = context.SharedHostContext().CheckHostSuitability()
        if syserr == nil {

            _, nerr = context.SharedHostContext().HostPrimaryAddress()
            if nerr == nil {

                vbox, vlerr = vboxglue.NewGOVboxGlue()
                defer vbox.Close()
                if vlerr == nil {

                    vererr = vbox.CheckVBoxSuitability()
                    if vererr == nil {

                        response = ReponseMessage{"syscheck": {"status": true}}
                    } else {
                        response = ReponseMessage{
                            "syscheck": {
                                "status": false,
                                "error": vererr.Error(),
                            },
                        }
                    }

                } else {
                    response = ReponseMessage{
                        "syscheck": {
                            "status": false,
                            "error": errors.WithMessage(vlerr, "Loading Virtualbox causes an error. Please install latest VirtualBox"),
                        },
                    }
                }

            } else {
                response = ReponseMessage{
                    "syscheck": {
                        "status": false,
                        "error": errors.WithMessage(nerr, "Unable to detect Wi-Fi network. Please enable Wi-Fi"),
                    },
                }
            }

        } else {
            response = ReponseMessage{
                "syscheck": {
                    "status": false,
                    "error": syserr.Error(),
                },
            }
        }

        // inline json marshalling
        data, err := json.Marshal(response)
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

    // check if app is expired
    theApp.GET(routepath.RpathAppExpired(), func(_, path, _ string) error {
        var (
            expired, warn, err = context.SharedHostContext().CheckIsApplicationExpired()
            response ReponseMessage = nil
        )
        if err != nil {
            response = ReponseMessage {
                "expired" : {
                    "status" : expired,
                    "error"  : err.Error(),
                },
            }
        } else if warn != nil {
            response = ReponseMessage {
                "expired" : {
                    "status" : expired,
                    "warning" : warn.Error(),
                },
            }
        } else {
            response = ReponseMessage {
                "expired" : {
                    "status" : expired,
                },
            }
        }

        data, err := json.Marshal(response)
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

    // check if this is the first time run
    theApp.GET(routepath.RpathSystemIsFirstRun(), func(_, path, _ string) error {
        data, err := json.Marshal(ReponseMessage{
            "firsttime": {
                "status" : context.SharedHostContext().CheckIsFistTimeExecution(),
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

    // check if user is authenticated
    theApp.GET(routepath.RpathUserAuthed(), func(_, path, _ string) error {
        data, err := json.Marshal(ReponseMessage{
            "user-auth": {
                "status": true,
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
                "description": pkgs[i].Description,
                "id" : pkgs[i].PkgID,
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

}