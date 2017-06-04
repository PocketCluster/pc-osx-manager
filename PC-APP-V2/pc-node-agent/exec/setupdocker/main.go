package main

import (
    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-node-agent/slcontext/config"
)

func main() {
    err := config.SetupDockerEnvironement("")
    if err != nil {
        log.Info(err.Error())
    }

    err = config.SetupDockerAuthorityCert("")
    if err != nil {
        log.Info(err.Error())
    }
}
