package initcheck

import (
    "encoding/json"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/vboxglue"
)

type ReponseMessage map[string]map[string]interface{}

func InitRoutePathServices(appLife route.Router) {

    // check if this system is suitable to run
    appLife.GET(routepath.RpathSystemReadiness(), func(_, path, _ string) error {
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
    appLife.GET(routepath.RpathAppExpired(), func(_, path, _ string) error {
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
    appLife.GET(routepath.RpathSystemIsFirstRun(), func(_, path, _ string) error {
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
    appLife.GET(routepath.RpathUserAuthed(), func(_, path, _ string) error {
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
}