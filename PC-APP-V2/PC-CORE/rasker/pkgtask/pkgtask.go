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
    "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/route"
)

const (
    fbPackageStartup           string = "package-start"
    fbPackageProcess           string = "package-proc"
    fbPackageKill              string = "package-kill"

    iventPackageKillPrefix     string = "ivent.rasker.package.kill."
    raskerPackageStartupPrefix string = "rasker.pacakge.startup."
    raskerPackageProcessPrefix string = "rasker.pacakge.process."
    raskerPackageKillPrefix    string = "rasker.pacakge.kill."
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

func buildComposeTemplateWithNodeList(template []byte, nodeList []string) ([]byte, error) {
    const (
        nodeNamePrefix = "pc-node"
    )
    if len(nodeList) == 0 {
        return nil, errors.Errorf("unable to generate template with empty node list")
    }
    if len(template) == 0 {
        return nil, errors.Errorf("invalid template to reinstate")
    }

    // build node data
    var snodes = []pongo2.Context{}
    for _, nodeName := range nodeList {

        // we skip pc-core for now
        if !strings.HasPrefix(nodeName, nodeNamePrefix) {
            continue
        }

        nidx := strings.Replace(nodeName, nodeNamePrefix, "", -1)
        nadr, err := strconv.Atoi(nidx)
        if err != nil {
            continue
        }

        snodes = append(
            snodes,
            pongo2.Context{
                "index":   nidx,
                "address": nadr + 1,
                "name":    nodeName,
            })
    }
    if len(snodes) == 0 {
        return nil, errors.Errorf("unable to generate proper node list to initiate cluster")
    }
    data := pongo2.Context{
        "corenode": pongo2.Context{
            "name": defaults.PocketClusterCoreName,
            "address": 1,
        },
        "slavenodes": snodes,
    }

    // bring template into life
    tpl, err := pongo2.FromString(string(template))
    if err != nil {
        log.Error(errors.WithStack(err).Error())
    }

    // reinstate template with data
    return tpl.ExecuteBytes(data)
}