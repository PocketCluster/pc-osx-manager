package config

import (
    "os"
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

    // make slave pki path
    err = os.MkdirAll(path.Join(cfg.rootPath, CORE_CERTS_DIR), cert_path_permission);
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // first save some file to pki place
    err = ioutil.WriteFile(path.Join(cfg.rootPath, CoreAuthCertFileName), pcrypto.TestCertPublicAuth(), cert_file_permission)
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

    // *** This run assumes that we don't have backup *** //
    // prepare final cert
    var updatedCert []byte = append(pcrypto.TestCertPublicAuth(), pcrypto.TestCertPublicAuth()...)
    // make system native cert dir
    err = os.MkdirAll(path.Join(cfg.rootPath, "/etc/ssl/certs/"), cert_path_permission);
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // system native cert
    err = ioutil.WriteFile(path.Join(cfg.rootPath, SYSTEM_AUTH_CERT_NATIVE_FILE), pcrypto.TestCertPublicAuth(), cert_file_permission)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // make slave pki path
    err = os.MkdirAll(path.Join(cfg.rootPath, CORE_CERTS_DIR), cert_path_permission);
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // slave cert
    err = ioutil.WriteFile(path.Join(cfg.rootPath, CoreAuthCertFileName), pcrypto.TestCertPublicAuth(), cert_file_permission)
    if err != nil {
        t.Errorf(err.Error())
        return
    }

    // append new docker cert to system native
    err = AppendAuthCertFowardSystemCertAuthority(cfg.rootPath)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    sca, err := ioutil.ReadFile(path.Join(cfg.rootPath, SYSTEM_AUTH_CERT_NATIVE_FILE))
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if !reflect.DeepEqual(sca, updatedCert) {
        t.Errorf(err.Error())
        return
    }

    // *** This run assumes that we "HAVE" backup *** //
    // firstly make sure backup is what we are looking for
    sbc, err := ioutil.ReadFile(path.Join(cfg.rootPath, SYSTEM_AUTH_CERT_BACKUP_FILE))
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if !reflect.DeepEqual(sbc, pcrypto.TestCertPublicAuth()) {
        t.Errorf(err.Error())
        return
    }
    // then append cert to make sure everything goes as expected
    err = AppendAuthCertFowardSystemCertAuthority(cfg.rootPath)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    sca, err = ioutil.ReadFile(path.Join(cfg.rootPath, SYSTEM_AUTH_CERT_NATIVE_FILE))
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if !reflect.DeepEqual(sca, updatedCert) {
        t.Errorf(err.Error())
        return
    }
}