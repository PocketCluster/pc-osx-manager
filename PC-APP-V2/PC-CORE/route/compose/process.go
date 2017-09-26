package compose

import (
    "encoding/json"

    "golang.org/x/net/context"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/docker/libcompose/docker"
    "github.com/docker/libcompose/docker/ctx"
    "github.com/docker/libcompose/project"

    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
)

func InitPackageProcessRoutePath(appLife route.Router, feeder route.ResponseFeeder) {

    // install a package
    appLife.POST(routepath.RpathPackageInstall(), func(_, rpath, payload string) error {
        var (
            columns []string = []string{"Id", "Name", "Command", "State", "Ports"}
            pkgID     string = ""
        )
        // 1. parse input package id
        err := json.Unmarshal([]byte(payload), &struct {
            PkgID *string `json:"pkg-id"`
        }{&pkgID})
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackProcess, errors.WithMessage(err, "unable to specify package id"))
        }

        // 2. load template
        cTempl, err := loadComposeTemplate(pkgID)
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackProcess, errors.WithMessage(err, "unable to access package template"))
        }

        // 3. build client
        opts, err := newComposeClient()
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackProcess, errors.WithMessage(err, "unable to build orchestration client"))
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
            return feedError(feeder, rpath, packageFeedbackProcess, errors.WithMessage(err, "unable to create project"))
        }

        // 5. cluster process list
        allInfo, err := project.Ps(context.Background(), []string{}...)
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackProcess, errors.WithMessage(err, "unable to list cluster process"))
        }
        pslist := allInfo.String(columns, false)

        // 6. return feedback
        data, err := json.Marshal(route.ReponseMessage{
            packageFeedbackProcess: {
                "status": true,
                "pkg-id" : pkgID,
                "process": pslist,
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
