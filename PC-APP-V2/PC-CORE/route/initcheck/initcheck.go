package initcheck

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/url"
    "strings"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "golang.org/x/crypto/ripemd160"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/vboxglue"
    "github.com/stkim1/pc-core/utils/apireq"
)

func InitApplicationCheck(appLife route.Router, feeder route.ResponseFeeder) {

    // check if this system is suitable to run
    appLife.GET(routepath.RpathSystemReadiness(), func(_, path, _ string) error {
        var (
            syserr, nerr, vlerr, vererr error = nil, nil, nil, nil
            vbox vboxglue.VBoxGlue = nil
            data []byte = nil
            response route.ReponseMessage = nil
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

                        response = route.ReponseMessage{"syscheck": {"status": true}}
                    } else {
                        response = route.ReponseMessage{
                            "syscheck": {
                                "status": false,
                                "error": vererr.Error(),
                            },
                        }
                    }

                } else {
                    response = route.ReponseMessage{
                        "syscheck": {
                            "status": false,
                            "error": errors.WithMessage(vlerr, "Loading Virtualbox causes an error. Please install latest VirtualBox"),
                        },
                    }
                }

            } else {
                response = route.ReponseMessage{
                    "syscheck": {
                        "status": false,
                        "error": errors.WithMessage(nerr, "Unable to detect Wi-Fi network. Please enable Wi-Fi"),
                    },
                }
            }

        } else {
            response = route.ReponseMessage{
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
        err = feeder.FeedResponseForGet(path, string(data))
        if err != nil {
            return errors.WithStack(err)
        }
        return nil
    })

    // check if app is expired
    appLife.GET(routepath.RpathAppExpired(), func(_, path, _ string) error {
        var (
            expired, warn, err = context.SharedHostContext().CheckIsApplicationExpired()
            response route.ReponseMessage = nil
        )
        if err != nil {
            response = route.ReponseMessage {
                "expired" : {
                    "status" : expired,
                    "error"  : err.Error(),
                },
            }
        } else if warn != nil {
            response = route.ReponseMessage {
                "expired" : {
                    "status" : expired,
                    "warning" : warn.Error(),
                },
            }
        } else {
            response = route.ReponseMessage {
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
        err = feeder.FeedResponseForGet(path, string(data))
        if err != nil {
            return errors.WithStack(err)
        }
        return nil
    })

    // check if this is the first time run
    appLife.GET(routepath.RpathSystemIsFirstRun(), func(_, path, _ string) error {
        data, err := json.Marshal(route.ReponseMessage{
            "firsttime": {
                "status" : context.SharedHostContext().CheckIsFistTimeExecution(),
            },
        })
        if err != nil {
            log.Debugf(err.Error())
            return errors.WithStack(err)
        }
        err = feeder.FeedResponseForGet(path, string(data))
        if err != nil {
            return errors.WithStack(err)
        }
        return nil
    })

    // check if user is authenticated
    appLife.POST(routepath.RpathUserAuthed(), func(_, path, payload string) error {
        var (
            invitation, autherr string
            status bool = false
        )

        // 1. parse input package id
        err := json.Unmarshal([]byte(payload), &struct {
            PkgID *string `json:"invitation"`
        }{&invitation})
        if err != nil || len(invitation) != 24 {
            data, _ := json.Marshal(route.ReponseMessage{
                "user-auth": {
                    "status": false,
                    "error":  "invalid invitation code. please provide valid one",
                },
            })
            feeder.FeedResponseForPost(path, string(data))
            return err
        }

        // 2. connect to service
        val := url.Values{}
        rh := ripemd160.New()
        rh.Write([]byte(invitation))
        val.Set("invitation", fmt.Sprintf("%x", rh.Sum(nil)))

        rh.Reset()
        rh.Write([]byte(context.SharedHostContext().HostDeviceSerial()))
        val.Set("device", fmt.Sprintf("%x", rh.Sum(nil)))

        req, _ :=  http.NewRequest("POST", fmt.Sprintf("%s/service/v014/auth/check", defaults.PocketClusterAPIHost), strings.NewReader(val.Encode()))
        req.Header.Add("User-Agent", "PocketCluster/0.1.4 (OSX)")
        req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        req.ProtoAtLeast(1, 1)
        client := apireq.NewClient(apireq.ConnTimeout, true)
        resp, err := apireq.ReadRequest(req, client)
        if err != nil {
            data, _ := json.Marshal(route.ReponseMessage{
                "user-auth": {
                    "status": false,
                    "error":  "unable to connect service. please try again",
                },
            })
            feeder.FeedResponseForPost(path, string(data))
            return err
        }

        // 3. read response
        err = json.Unmarshal(resp, &struct {
            AuthError *string `json:"error"`
        }{&autherr})
        if err != nil {
            data, _ := json.Marshal(route.ReponseMessage{
                "user-auth": {
                    "status": false,
                    "error":  "invalid response from service. please try again",
                },
            })
            feeder.FeedResponseForPost(path, string(data))
            return err
        }

        // 4. return auth check value
        if len(autherr) == 0 {
            status = true
        } else {
            status = false
        }
        data, _ := json.Marshal(route.ReponseMessage{
            "user-auth": {
                "status": status,
                "error": autherr,
            },
        })
        feeder.FeedResponseForPost(path, string(data))
        return nil
    })
}