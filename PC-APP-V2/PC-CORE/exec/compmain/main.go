package main

import (
    "os"

    "github.com/docker/libcompose/docker"
    "github.com/docker/libcompose/docker/ctx"
    "github.com/docker/libcompose/project"
    "github.com/docker/libcompose/project/options"

    "golang.org/x/net/context"
    log "github.com/Sirupsen/logrus"
    "github.com/davecgh/go-spew/spew"
)

func main() {
    log.SetOutput(os.Stdout)
    project, err := docker.NewProject(&ctx.Context{
        Context: project.Context{
            ComposeFiles: []string{"pocket-deploy-original.yml"},
            ProjectName:  "pocket-hadoop",
        },
    }, nil)

    if err != nil {
        log.Fatal(err)
    }

    log.Info(spew.Sdump(project))
    return

    err = project.Up(context.Background(), options.Up{})
    if err != nil {
        log.Fatal(err)
    }
}
