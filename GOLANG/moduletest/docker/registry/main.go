package main

import (
    "github.com/docker/distribution/configuration"
    "github.com/stkim1/pc-core/extlib/registry"
    log "github.com/Sirupsen/logrus"
    "github.com/davecgh/go-spew/spew"
)

/*
version: 0.1
log:
    level: debug
    formatter: text
    fields:
        service: registry
        environment: pc-master
storage:
    cache:
        blobdescriptor: inmemory
    filesystem:
        rootdirectory: /Users/almightykim/Workspace/DKIMG/REGISTRY/data
        maxthreads: 32
    maintenance:
        uploadpurging:
            enabled: true
            age: 168h
            interval: 24h
            dryrun: false
        readonly:
            enabled: false
http:
    addr: 0.0.0.0:5000
    net: tcp
    secret: mytokensecret
    relativeurls: false
    tls:
        certificate: /Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.cert
        key: /Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.key
    debug:
        addr: 0.0.0.0:5001
    headers:
        X-Content-Type-Options: [nosniff]
 */

func main() {
    var config *configuration.Configuration
    if false {
        config, _ = registry.ParseYamlConfig("/Users/almightykim/Workspace/DKIMG/REGISTRY/config/config.yml")
    } else {
        config = registry.NewRegistrySampleConfig()
    }
    log.Info(spew.Sdump(config))
    registry.Serve(config)
}


