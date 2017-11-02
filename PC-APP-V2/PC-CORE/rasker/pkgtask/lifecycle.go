package pkgtask

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

    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
    "github.com/stkim1/pc-core/rasker"
    "github.com/stkim1/pc-core/service"
    "github.com/stkim1/pc-core/service/ivent"
)

// initiate a package
func InitPackageLifeCycle(appLife rasker.RouteTasker, feeder route.ResponseFeeder) error {
    return appLife.POST(routepath.RpathPackageStartup(), func(_, rpath, payload string) error {
        // 1. parse input package id
        var (
            reportC = make(chan service.Event)
            pkgID string = ""
            pkg   *model.Package = nil
        )
        err := json.Unmarshal([]byte(payload), &struct {
            PkgID *string `json:"pkg-id"`
        }{&pkgID})
        if err != nil {
            return feedError(feeder, rpath, fbPackageStartup, errors.WithMessage(err, "unable to specify package id"))
        }

        // TODO check if package has started

        pkgs, err := model.FindPackage("pkg_id = ?", pkgID)
        if err != nil {
            return feedError(feeder, rpath, fbPackageStartup, errors.WithMessage(err, "unable to specify package id"))
        }
        pkg = pkgs[0]
        log.Infof("Package Meta Found %v", pkg)

        // 2. get the node list report
        err = appLife.BindDiscreteEvent(ivent.IventReportLiveNodesResult, reportC)
        if err != nil {
            return feedError(feeder, rpath, fbPackageStartup, errors.WithMessage(err, "unable to access node list report"))
        }
        // ask node list
        appLife.BroadcastEvent(service.Event{Name:ivent.IventReportLiveNodesRequest})
        nr := <- reportC
        appLife.UntieDiscreteEvent(ivent.IventReportLiveNodesResult)
        nlist, ok := nr.Payload.([]string)
        if !ok {
            return feedError(feeder, rpath, fbPackageStartup, errors.WithMessage(err, "unable to access proper node list"))
        }

        // 3. load template
        tmpl, err := model.FindTemplateWithPackageID(pkgID)
        if err != nil {
            return feedError(feeder, rpath, fbPackageStartup, errors.WithMessage(err, "unable to access package template"))
        }
        cTempl, err := buildComposeTemplateWithNodeList(tmpl.Body, nlist)
        log.Info(string(cTempl))
        if err != nil {
            return feedError(feeder, rpath, fbPackageStartup, errors.WithMessage(err, "unable to build package exec plan"))
        }

        // 4. build client
        opts, err := newComposeClient()
        if err != nil {
            return feedError(feeder, rpath, fbPackageStartup, errors.WithMessage(err, "unable to build orchestration client"))
        }

        // 5. build package
        iventKillTag   := fmt.Sprintf("%s%s", iventPackageKillPrefix,     pkgID)
        taskProcessTag := fmt.Sprintf("%s%s", raskerPackageProcessPrefix, pkgID)
        taskKillTag    := fmt.Sprintf("%s%s", raskerPackageKillPrefix,    pkgID)
        project, err := docker.NewPocketProject(&docker.PocketContext{
            Context: &ctx.Context{
                Context: project.Context{
                    ProjectName: pkg.MenuName,
                },
            },
            ClientOptions: opts,
            Manifest: cTempl,
        }, nil)
        if err != nil {
            return feedError(feeder, rpath, fbPackageStartup, errors.WithMessage(err, "unable to create project"))
        }

        // --- --- --- --- --- --- --- --- --- --- --- --- package process list --- --- --- --- --- --- --- --- --- //
        killPsC := make(chan service.Event)
        appLife.RegisterServiceWithFuncs(
            taskProcessTag,
            func() error {
                // cluster process list
                var (
                    rptPath = routepath.RpathPackageProcess()
                    columns = []string{"Id", "Name", "Command", "State", "Ports"}
                    timer   = time.NewTicker(time.Second * 5)
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
                                feedError(feeder, rptPath, fbPackageProcess, errors.WithMessage(err, "unable to list cluster process"))
                            }
                            pslist := allInfo.String(columns, false)

                            // return feedback
                            data, err := json.Marshal(route.ReponseMessage{
                                fbPackageProcess: {
                                    "status":  true,
                                    "pkg-id":  pkgID,
                                    "process": pslist,
                                },
                            })
                            if err != nil {
                                log.Error(err.Error())
                            }
                            err = feeder.FeedResponseForPost(rptPath, string(data))
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
                var (
                    killPath = routepath.RpathPackageKill()
                )
                // since two signals do the same thing (wait until signal kicks)
                select {
                    case <- appLife.StopChannel():
                    case <- killSigC:
                }

                // kill package
                err = project.Kill(context.TODO(), "SIGINT", []string{}...)
                if err != nil {
                    log.Error(err.Error())
                    //return feedError(feeder, killPath, fbPackageKill, errors.WithMessage(err, "unable to stop package"))
                }

                // delete package
                err = project.Delete(context.Background(), options.Delete{}, []string{}...)
                if err != nil {
                    return feedError(feeder, killPath, fbPackageKill, errors.WithMessage(err, "unable to remove package residue"))
                }

                // 7. return feedback
                data, err := json.Marshal(route.ReponseMessage{
                    fbPackageKill: {
                        "status": true,
                        "pkg-id" : pkgID,
                    },
                })
                // this should never happen
                if err != nil {
                    log.Error(err.Error())
                }
                err = feeder.FeedResponseForPost(killPath, string(data))
                return errors.WithStack(err)
            },
            service.BindEventWithService(iventKillTag, killSigC))

        // --- --- --- --- --- --- --- --- --- --- --- --- package startup --- --- --- --- --- --- --- --- --- --- //
        // startup package
        err = project.Up(context.TODO(), options.Up{})
        if err != nil {
            // kill monitor signal and delete
            appLife.BroadcastEvent(service.Event{Name:fmt.Sprintf("%s%s", iventPackageKillPrefix, pkgID)})
            return feedError(feeder, rpath, fbPackageStartup, errors.WithMessage(err, "unable to start package"))
        }
        // return feedback
        data, err := json.Marshal(route.ReponseMessage{
            fbPackageStartup: {
                "status": true,
                "pkg-id": pkgID,
            },
        })
        // this should never happen
        if err != nil {
            log.Error(err.Error())
        }
        err = feeder.FeedResponseForPost(rpath, string(data))
        return errors.WithStack(err)
    })
}
