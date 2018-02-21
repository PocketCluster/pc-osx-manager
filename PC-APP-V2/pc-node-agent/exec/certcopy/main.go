package main

import (
    log "github.com/Sirupsen/logrus"

    "github.com/stkim1/pc-node-agent/slcontext/config"
)

func main() {
    err := config.CopyCertAuthForwardCustomCertStorage("")
    if err != nil {
        log.Debugf(err.Error())
    }
}
