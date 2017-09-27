package compose

import (
    "encoding/json"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/docker/libcompose/docker/client"

    pcctx "github.com/stkim1/pc-core/context"
    _ "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/route"
)

const (
    packageFeedbackStartup string = "package-startup"
    packageFeedbackProcess string = "package-process"
    packageFeedbackKill    string = "package-kill"
)

func feedError(feeder route.ResponseFeeder, rpath, fpath string, irr error) error {
    log.Error(irr.Error())
    data, frr := json.Marshal(route.ReponseMessage{
        fpath: {
            "status": false,
            "error" : irr.Error(),
        },
    })
    // this should never happen
    if frr != nil {
        log.Error(frr.Error())
    }
    frr = feeder.FeedResponseForPost(rpath, string(data))
    if frr != nil {
        log.Error(frr.Error())
    }
    return irr
}

func newComposeClient() (*client.PocketClientOption, error) {
    caCert, err := pcctx.SharedHostContext().CertAuthCertificate()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    hostCrt, err := pcctx.SharedHostContext().MasterHostCertificate()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    hostKey, err := pcctx.SharedHostContext().MasterHostPrivateKey()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return client.NewPocketCientOption(caCert, hostCrt, hostKey, "tcp://pc-master:3376")
}


const (
    iventPackageKillPrefix   string = "ivent.package.kill."
    taskPackageStartupPrefix string = "task.pacakge.startup."
    taskPackageProcessPrefix string = "task.pacakge.process."
    taskPackageKillPrefix    string = "task.pacakge.kill."
)

func loadComposeTemplate(pkgID string) ([]byte, error) {
    return nil, errors.Errorf("no package for %v found", pkgID)
}
