package main

import (
    "encoding/json"
    "fmt"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/event/route/routepath"
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
            response ReponseMessage = nil
            nowDate = time.Now()
            expDate, err = time.Parse(defaults.PocketTimeDateFormat, defaults.ApplicationExpirationDate)
        )
        if err != nil {
            log.Debugf("parse error %v", err.Error())
        }

        // check if exp date exceed today
        if expDate.After(nowDate) {
            timeLeft := expDate.Sub(nowDate)
            if timeLeft < time.Duration(time.Hour * 24 * 7) {
                response = ReponseMessage {
                    "expired" : {
                        "status" : false,
                        "warning" : fmt.Sprintf("Version %v will be expired within %d days", defaults.ApplicationVersion, int(timeLeft / time.Duration(time.Hour * 24))),
                    },
                }
            } else {
                response = ReponseMessage {
                    "expired" : {
                        "status" : false,
                    },
                }
            }
        } else {
            eYear, eMonth, eDay := expDate.Date()
            response = ReponseMessage {
                "expired" : {
                    "status" : true,
                    "error"  : fmt.Sprintf("Version %v is expired on %d/%d/%d", defaults.ApplicationVersion, eYear, eMonth, eDay),
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
                "status": false,
                "error": "please check your invitation code",
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

    theApp.GET(routepath.RpathCmdServiceStart(), func(_, _, _ string) error {
        return nil
    })

}