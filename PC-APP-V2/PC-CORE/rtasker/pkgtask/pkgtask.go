package pkgtask

import (
    "encoding/json"
    "strings"
    "strconv"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/flosch/pongo2"
    "github.com/docker/libcompose/docker/client"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/model"
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
    caCert, err := context.SharedHostContext().CertAuthCertificate()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    hostCrt, err := context.SharedHostContext().MasterHostCertificate()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    hostKey, err := context.SharedHostContext().MasterHostPrivateKey()
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

func loadComposeTemplate(pkgID string, nodeList []string) ([]byte, error) {
    if len(nodeList) == 0 {
        return nil, errors.Errorf("unable to generate template with empty node list")
    }
    // retrieve template
    tmpl, err := model.FindTemplateWithPackageID(pkgID)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // build node data
    var snodes = []pongo2.Context{}
    for _, node := range nodeList {

        if node == "pc-core" {
            continue
        }

        nidx := strings.Replace(node,"pc-node", "", -1)
        nadr, err := strconv.Atoi(nidx)
        if err != nil {
            continue
        }

        snodes = append(
            snodes,
            pongo2.Context{
                "index":   nidx,
                "address": nadr + 1,
                "name":    node,
            })
    }
    if len(snodes) == 0 {
        return nil, errors.Errorf("unable to generate proper node list to initiate cluster")
    }
    data := pongo2.Context{
        "corenode":   pongo2.Context{
            "name": "pc-core",
            "address": 1,
        },
        "slavenodes": snodes,
    }

    // bring template into life
    tpl, err := pongo2.FromString(string(tmpl.Body))
    if err != nil {
        log.Error(errors.WithStack(err).Error())
    }

    // reinstate template with data
    return tpl.ExecuteBytes(data)
}