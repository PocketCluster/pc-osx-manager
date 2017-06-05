package config

import (
    "path"
    "testing"
    "reflect"
    "io/ioutil"

    "github.com/stkim1/pcrypto"
)

func TestDockerEnvironment(t *testing.T) {
    cfg, err := DebugConfigPrepare()
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    defer DebugConfigDestory(cfg)

    // there shouldn't be an error
    err = SetupDockerEnvironement(cfg.rootPath)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // when there is no env file
    env, err := ioutil.ReadFile(path.Join(cfg.rootPath, DOCKER_ENV_FILE))
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if !reflect.DeepEqual(env, dockerEnvContent()) {
        t.Errorf(err.Error())
        return
    }

    // when there is file exists, overwrite it first
    err = ioutil.WriteFile(path.Join(cfg.rootPath, DOCKER_ENV_FILE), []byte("nothing"), cert_file_permission)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // there shouldn't be an error
    err = SetupDockerEnvironement(cfg.rootPath)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // when there is no env file
    env, err = ioutil.ReadFile(path.Join(cfg.rootPath, DOCKER_ENV_FILE))
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if !reflect.DeepEqual(env, dockerEnvContent()) {
        t.Errorf(err.Error())
    }
}

func TestDockerAuthorityCert(t *testing.T) {
    cfg, err := DebugConfigPrepare()
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    defer DebugConfigDestory(cfg)

    // first save some file to pki place
    err = ioutil.WriteFile(path.Join(cfg.rootPath, SlaveAuthCertFileName), pcrypto.TestCertPublicAuth(), cert_file_permission)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // setup docker authority
    err = SetupDockerAuthorityCert(cfg.rootPath)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // compare two certs to make sure they are the same
    dca, err := ioutil.ReadFile(path.Join(cfg.rootPath, DOCKER_AUTH_CERT_FILE))
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if !reflect.DeepEqual(dca, pcrypto.TestCertPublicAuth()) {
        t.Errorf(err.Error())
        return
    }

    // when there is file exists, overwrite it first
    err = ioutil.WriteFile(path.Join(cfg.rootPath, DOCKER_AUTH_CERT_FILE), []byte("nothing"), cert_file_permission)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // setup docker authority
    err = SetupDockerAuthorityCert(cfg.rootPath)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // compare two certs to make sure they are the same
    dca, err = ioutil.ReadFile(path.Join(cfg.rootPath, DOCKER_AUTH_CERT_FILE))
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if !reflect.DeepEqual(dca, pcrypto.TestCertPublicAuth()) {
        t.Errorf(err.Error())
        return
    }
}

func TestAppendAuthCertFowardSystemCertAuthority(t *testing.T) {
    cfg, err := DebugConfigPrepare()
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    defer DebugConfigDestory(cfg)

}