package slcontext

import (
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-node-agent/slcontext/config"
)

func DockerEnvironemtPostProcess() error {
    var (
        err error = nil
    )
    err = config.SetupDockerEnvironement("")
    if err != nil {
        log.Debugf(err.Error())
        return errors.WithStack(err)
    }
    err = config.SetupDockerAuthorityCert("")
    if err != nil {
        log.Debugf(err.Error())
        return errors.WithStack(err)
    }
    err = config.AppendAuthCertFowardSystemCertAuthority("")
    if err != nil {
        log.Debugf(err.Error())
        return errors.WithStack(err)
    }
    return nil
}

