package main

import (
    "io/ioutil"
    "net/http"

    log "github.com/Sirupsen/logrus"
    "golang.org/x/net/context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "github.com/stkim1/pc-core/utils/tlscfg"
)

func main() {

    caCert, err := ioutil.ReadFile("ca-cert.pem")
    if err != nil {
        log.Error(err.Error())
        return
    }
    hostCrt, err := ioutil.ReadFile("host-cert.pem")
    if err != nil {
        log.Error(err.Error())
        return
    }
    hostPrv, err := ioutil.ReadFile("host-key.pem")
    if err != nil {
        log.Error(err.Error())
        return
    }
    tlsc, err := tlscfg.BuildTLSConfigWithCAcert(caCert, hostCrt, hostPrv, true)
    if err != nil {
        log.Error(err.Error())
        return
    }

    httpcli := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: tlsc,
        },
    }
    cli, err := client.NewClient("tcp://pc-node1:2376", "1.12", httpcli, nil)
    if err != nil {
        log.Error(err.Error())
        return
    }
    _, err = cli.ImagePull(context.TODO(), "pc-master:5000/arm64v8-ubuntu", types.ImagePullOptions{})
    if err != nil {
        log.Error(err.Error())
        return
    }
}
