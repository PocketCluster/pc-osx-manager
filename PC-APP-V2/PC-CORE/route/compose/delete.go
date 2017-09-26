package compose

import (
    "encoding/json"

    "golang.org/x/net/context"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/docker/libcompose/docker"
    "github.com/docker/libcompose/docker/ctx"
    "github.com/docker/libcompose/project"
    "github.com/docker/libcompose/project/options"

    "github.com/stkim1/pc-core/route"
    "github.com/stkim1/pc-core/route/routepath"
)

func InitPackageDeleteRoutePath(appLife route.Router, feeder route.ResponseFeeder) {

    // install a package
    appLife.POST(routepath.RpathPackageInstall(), func(_, rpath, payload string) error {
        var (
            pkgID      string = ""
        )
        // 1. parse input package id
        err := json.Unmarshal([]byte(payload), &struct {
            PkgID *string `json:"pkg-id"`
        }{&pkgID})
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackDelete, errors.WithMessage(err, "unable to specify package id"))
        }

        // 2. load template
        cTempl, err := loadComposeTemplate(pkgID)
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackDelete, errors.WithMessage(err, "unable to access package template"))
        }

        // 3. build client
        opts, err := newComposeClient()
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackDelete, errors.WithMessage(err, "unable to build orchestration client"))
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
            return feedError(feeder, rpath, packageFeedbackDelete, errors.WithMessage(err, "unable to create project"))
        }

        // 5. delete package
        err = project.Delete(context.Background(), options.Delete{}, []string{}...)
        if err != nil {
            return feedError(feeder, rpath, packageFeedbackDelete, errors.WithMessage(err, "unable to start package"))
        }

        // 6. return feedback
        data, err := json.Marshal(route.ReponseMessage{
            packageFeedbackDelete: {
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
}
