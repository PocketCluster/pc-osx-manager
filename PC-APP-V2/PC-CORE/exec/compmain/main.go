package main

import (
    "os"
    "io/ioutil"

    "github.com/docker/libcompose/docker"
    "github.com/docker/libcompose/docker/ctx"
    "github.com/docker/libcompose/docker/client"
    "github.com/docker/libcompose/project"
    "github.com/docker/libcompose/project/options"

    "golang.org/x/net/context"
    log "github.com/Sirupsen/logrus"
    //"github.com/davecgh/go-spew/spew"
)

func main() {
    log.SetOutput(os.Stdout)
    composeBytes, err := ioutil.ReadFile("pocket-deploy.json")
    if err != nil && !os.IsNotExist(err) {
        log.Fatal("Failed to open the compose file: pocket-deploy.json")
    }
    caCert, err  := ioutil.ReadFile("/Users/almightykim/Workspace/DKIMG/CERT/ca-cert.pub")
    if err != nil {
        log.Fatal(err.Error())
    }
    tlsCert, err := ioutil.ReadFile("/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.cert")
    if err != nil {
        log.Fatal(err.Error())
    }
    tlsKey, err  := ioutil.ReadFile("/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.key")
    if err != nil {
        log.Fatal(err.Error())
    }

    opts, err := client.NewPocketCientOption(caCert, tlsCert, tlsKey, "tcp://192.168.1.150:3376")
    if err != nil {
        log.Fatal(err.Error())
    }
    project, err := docker.NewPocketProject(&docker.PocketContext{
        Context: &ctx.Context{
            Context: project.Context{
                ProjectName:  "pocket-hadoop",
            },
        },
        ClientOptions: opts,
        Manifest: composeBytes,
    }, nil)
    if err != nil {
        log.Fatal(err)
    }

    //log.Info(spew.Sdump(project))
    allInfo, err := project.Ps(context.Background(), []string{}...)
    if err != nil {
        log.Fatal(err)
    }
    columns := []string{"Id", "Name", "Command", "State", "Ports"}
    os.Stdout.WriteString(allInfo.String(columns, false))
    return

    err = project.Up(context.Background(), options.Up{})
    if err != nil {
        log.Fatal(err)
    }
}
