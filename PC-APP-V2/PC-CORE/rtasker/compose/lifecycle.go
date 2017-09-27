package compose

import (
    "encoding/json"
    "fmt"
    "time"

    "golang.org/x/net/context"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/docker/libcompose/docker"
    "github.com/docker/libcompose/docker/ctx"
    "github.com/docker/libcompose/project"
    "github.com/docker/libcompose/project/options"

    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/rtasker"
    "github.com/stkim1/pc-core/service"
)

func InitPackageProcess(appLife rtasker.RouteTasker, feeder route.ResponseFeeder) error {

    // install a package
    appLife.POST(routepath.RpathPackageStartup(), func(_, rpath, payload string) error {
        // 1. parse input package id
        var (
            pkgID string = ""
        )
        err := json.Unmarshal([]byte(payload), &struct {
            PkgID *string `json:"pkg-id"`
        }{&pkgID})
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackStartup, errors.WithMessage(err, "unable to specify package id"))
        }

        // 2. load template
        cTempl, err := loadComposeTemplate(pkgID)
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackStartup, errors.WithMessage(err, "unable to access package template"))
        }

        // 3. build client
        opts, err := newComposeClient()
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackStartup, errors.WithMessage(err, "unable to build orchestration client"))
        }

        // 4. build package
        project, err := docker.NewPocketProject(&docker.PocketContext{
            Context: &ctx.Context{
                Context: project.Context{
                    ProjectName:  "pocket-hadoop",
                },
            },
            ClientOptions: opts,
            Manifest: cTempl,
        }, nil)
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackStartup, errors.WithMessage(err, "unable to create project"))
        }

        var (
            iventKillTag   string = fmt.Sprintf("%s%s", iventPackageKillPrefix, pkgID)
            taskStartTag   string = fmt.Sprintf("%s%s", taskPackageStartupPrefix, pkgID)
            taskProcessTag string = fmt.Sprintf("%s%s", taskPackageProcessPrefix, pkgID)
            taskKillTag    string = fmt.Sprintf("%s%s", taskPackageKillPrefix, pkgID)
        )

        // --- --- --- --- --- --- --- --- --- --- --- --- package startup --- --- --- --- --- --- --- --- --- --- //
        appLife.RegisterServiceWithFuncs(
            taskStartTag,
            func() error {
                // startup package
                err = project.Up(context.TODO(), options.Up{})
                if err != nil {
                    return feedError(feeder, rpath, packageFeedbackStartup, errors.WithMessage(err, "unable to start package"))
                }
                // return feedback
                data, err := json.Marshal(route.ReponseMessage{
                    packageFeedbackStartup: {
                        "status": true,
                        "pkg-id" : pkgID,
                    },
                })
                // this should never happen
                if err != nil {
                    log.Error(err.Error())
                }
                err = feeder.FeedResponseForPost(rpath, string(data))
                return errors.WithStack(err)
            })

        // --- --- --- --- --- --- --- --- --- --- --- --- package process list --- --- --- --- --- --- --- --- --- //
        killPsC := make(chan service.Event)
        appLife.RegisterServiceWithFuncs(
            taskProcessTag,
            func() error {
                // cluster process list
                var (
                    columns []string = []string{"Id", "Name", "Command", "State", "Ports"}
                    timer = time.NewTicker(time.Second * 3)
                )

                // process list doesn't quit until signals comes in
                for {
                    select {
                        case <- appLife.StopChannel(): {
                            timer.Stop()
                            return nil
                        }
                        case <- killPsC: {
                            timer.Stop()
                            return nil
                        }
                        case <- timer.C: {
                            allInfo, err := project.Ps(context.Background(), []string{}...)
                            if err != nil {
                                feedError(feeder, rpath, packageFeedbackProcess, errors.WithMessage(err, "unable to list cluster process"))
                            }
                            pslist := allInfo.String(columns, false)

                            // return feedback
                            data, err := json.Marshal(route.ReponseMessage{
                                packageFeedbackProcess: {
                                    "status": true,
                                    "pkg-id" : pkgID,
                                    "process": pslist,
                                },
                            })
                            if err != nil {
                                log.Error(err.Error())
                            }
                            err = feeder.FeedResponseForPost(rpath, string(data))
                            if err != nil {
                                log.Error(err.Error())
                            }
                        }
                    }
                }
            },
            service.BindEventWithService(iventKillTag, killPsC))

        // --- --- --- --- --- --- --- --- --- --- --- --- package kill & delete --- --- --- --- --- --- --- --- --- //
        killSigC := make(chan service.Event)
        appLife.RegisterServiceWithFuncs(
            taskKillTag,
            func() error {
                // since two signals do the same thing (wait until signal kicks)
                select {
                    case <- appLife.StopChannel():
                    case <- killSigC:
                }

                // kill package
                err = project.Kill(context.TODO(), "SIGINT", []string{}...)
                if err != nil {
                    return feedError(feeder, rpath, packageFeedbackKill, errors.WithMessage(err, "unable to stop package"))
                }

                // 6. kill package
                err = project.Delete(context.Background(), options.Delete{}, []string{}...)
                if err != nil {
                    return feedError(feeder, rpath, packageFeedbackKill, errors.WithMessage(err, "unable to remove package residue"))
                }

                // 7. return feedback
                data, err := json.Marshal(route.ReponseMessage{
                    packageFeedbackKill: {
                        "status": true,
                        "pkg-id" : pkgID,
                    },
                })
                // this should never happen
                if err != nil {
                    log.Error(err.Error())
                }
                err = feeder.FeedResponseForPost(rpath, string(data))
                return errors.WithStack(err)
            },
            service.BindEventWithService(iventKillTag, killSigC))

        // we should not return anything since this is unnecessary as package startup will report success/failure of
        // package startup status
        return nil
    })

    return nil
}
