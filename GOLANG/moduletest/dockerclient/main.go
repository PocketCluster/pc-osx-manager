package main

import (
    "archive/tar"
    "bytes"
    "time"
    "strings"
    "os"

    log "github.com/Sirupsen/logrus"
    "github.com/gravitational/trace"
    dkclt "github.com/fsouza/go-dockerclient"
)

func main() {
    client, err := dkclt.NewVersionedClient("tcp://192.168.1.154:2375", "1.22")
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }

    imgOpt := dkclt.ListImagesOptions{All:true}
    images, err := client.ListImages(imgOpt)
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }
    for _, img := range images {
        log.Infof("Image ID %s | %s", img.ID, strings.Join(img.RepoTags, ","))
    }

    // list of existing container
    containers, err := client.ListContainers(
        dkclt.ListContainersOptions{
            All:    true,
    })
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }
    for _, cntnr := range containers {
        log.Infof("Container to be removed ID %s | %s", cntnr.ID, cntnr.Image)

        err = client.RemoveContainer(
            dkclt.RemoveContainerOptions{
                ID:     cntnr.ID,
        })
        if err != nil {
            log.Info(trace.Wrap(err))
        }
    }


    // create a container to run
    cntnr, err := client.CreateContainer(
        dkclt.CreateContainerOptions{
            Name:       "TestRunContainer",
            Config:     &dkclt.Config{
                    Image:  "3fd2eb0e5e67ffedd937cd9fa88e443ca73c6057fbb1340ba1d6dafc51b3a24e",
                },
    })
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }

    // run image
    err = client.StartContainer(cntnr.ID, &dkclt.HostConfig{})
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }
    err = client.Logs(dkclt.LogsOptions{
        Container:      cntnr.ID,
        OutputStream:   os.Stdout,
        Stdout:         true,
    })
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }

    // stop container
    err = client.StopContainer(cntnr.ID, 0300)
    if err != nil {
        log.Fatal(trace.Wrap(err))
    }

    // commit container
    img, err := client.CommitContainer(
        dkclt.CommitContainerOptions{
            Container:      cntnr.ID,
            Repository:     "pocketcluster-hadoop",
            Tag:            "modified",
    })
    if err != nil {
        log.Error(trace.Wrap(err))
    } else {
        log.Infof("Committed Image %s", img.ID)
    }

    // once it runs fine, let's remove container
    err = client.RemoveContainer(dkclt.RemoveContainerOptions{
        ID:     cntnr.ID,
    })
    if err != nil {
        log.Info(trace.Wrap(err))
    }

    return

    t := time.Now()
    inputbuf, outputbuf := bytes.NewBuffer(nil), bytes.NewBuffer(nil)

    tr := tar.NewWriter(inputbuf)
    tr.WriteHeader(&tar.Header{Name: "Dockerfile", Size: 10, ModTime: t, AccessTime: t, ChangeTime: t})
    tr.Write([]byte("FROM base\n"))
    tr.Close()
    opts := dkclt.BuildImageOptions{
        Name:         "test",
        InputStream:  inputbuf,
        OutputStream: outputbuf,
    }
    if err := client.BuildImage(opts); err != nil {
        log.Fatal(err)
    }
}
