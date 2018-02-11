package slcontext

import (
    "io/ioutil"

    log "github.com/Sirupsen/logrus"

    "github.com/gravitational/teleport/lib/auth"
    "github.com/stkim1/pc-node-agent/slcontext/config"
)

func writeDockerKeyAndCert(certPack *auth.PocketResponseAuthKeyCert) error {
    var (
        AuthorityCertFile     string = SharedSlaveContext().SlaveEngineAuthCertFileName()
        NodeEngineKeyFile     string = SharedSlaveContext().SlaveEnginePrivateKeyFileName()
        NodeEngineCertFile    string = SharedSlaveContext().SlaveEngineKeyCertFileName()
    )
    log.Debugf("write slave docker auth to %v, key to %v, cert from %v", AuthorityCertFile, NodeEngineKeyFile, NodeEngineCertFile)
    if err := ioutil.WriteFile(AuthorityCertFile, certPack.Auth, 0600); err != nil {
        return err
    }
    if err := ioutil.WriteFile(NodeEngineKeyFile,  certPack.Key, 0600); err != nil {
        return err
    }
    if err := ioutil.WriteFile(NodeEngineCertFile, certPack.Cert, 0600); err != nil {
        return err
    }
    return nil
}

func DockerEnvironemtPostProcess(certPack *auth.PocketResponseAuthKeyCert) error {
    err := writeDockerKeyAndCert(certPack)
    if err != nil {
        log.Debugf(err.Error())
        return err
    }
    err = config.SetupDockerEnvironement("")
    if err != nil {
        log.Debugf(err.Error())
        return err
    }
    err = config.AppendAuthCertFowardSystemCertAuthority("")
    if err != nil {
        log.Debugf(err.Error())
        return err
    }
    err = config.CopyCertAuthForwardCustomCertStorage("")
    if err != nil {
        log.Debugf(err.Error())
        return err
    }
    return nil
}
