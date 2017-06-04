package config

import (
    "io"
    "io/ioutil"
    "os"
    "path"

    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
)

const (
    cert_path_permission = os.FileMode(0755)
    cert_file_permission = os.FileMode(0644)

    DOCKER_ENV_PATH string = "/etc/default/"
    DOCKER_ENV_FILE string = DOCKER_ENV_PATH + "docker"

    DOCKER_AUTH_CERT_PATH string = "/etc/docker/certs.d/pc-master/"
    DOCKER_AUTH_CERT_FILE string = DOCKER_AUTH_CERT_PATH + "pc_cert_auth" + pcrypto.FileExtCertificate
)

func dockerEnvContent() []byte {
    return []byte(`# PocketCluster Docker Upstart and SysVinit configuration file

DOCKER_OPTS="-H tcp://0.0.0.0:2375 --dns 127.0.0.1 --tlsverify --tlscacert=/etc/pocket/pki/pc_cert_auth.crt --tlscert=/etc/pocket/pki/pc_node_engine.crt --tlskey=/etc/pocket/pki/pc_node_engine.pem --cluster-advertise=eth0:2375 --cluster-store=etcd://pc-master:2379 --cluster-store-opt kv.cacertfile=/etc/pocket/pki/pc_cert_auth.crt --cluster-store-opt kv.certfile=/etc/pocket/pki/pc_node_engine.crt --cluster-store-opt kv.keyfile=/etc/pocket/pki/pc_node_engine.pem"`)
}

func SetupDockerEnvironement(rootPath string) error {
    var (
        dockerEnvPath string = path.Join(rootPath, DOCKER_ENV_PATH)
        dockerEnvFile string = path.Join(rootPath, DOCKER_ENV_FILE)
        err error = nil
    )
    if !path.IsAbs(dockerEnvPath) {
        return errors.Errorf("[ERR] invalid root path")
    }
    _, err = os.Stat(dockerEnvPath)
    if err != nil {
        if os.IsNotExist(err) {
            os.MkdirAll(dockerEnvPath, cert_path_permission);
        } else {
            return errors.WithStack(err)
        }
    }
    _, err = os.Stat(dockerEnvFile)
    if err != nil {
        if os.IsExist(err) {
            os.Remove(dockerEnvFile)
        } else {
            return errors.WithStack(err)
        }
    }

    err = ioutil.WriteFile(dockerEnvFile, dockerEnvContent(), cert_file_permission)
    return errors.WithStack(err)
}

func SetupDockerAuthorityCert(rootPath string) error {
    var (
        dockerAuthCertPath string = path.Join(rootPath, DOCKER_AUTH_CERT_PATH)
        dockerAuthCertFile string = path.Join(rootPath, DOCKER_AUTH_CERT_FILE)
        err error = nil
    )
    if !path.IsAbs(dockerAuthCertPath) {
        return errors.Errorf("[ERR] invalid root path")
    }
    _, err = os.Stat(dockerAuthCertPath)
    if err != nil {
        if os.IsNotExist(err) {
            os.MkdirAll(dockerAuthCertPath, cert_path_permission);
        } else {
            return errors.WithStack(err)
        }
    }
    _, err = os.Stat(dockerAuthCertPath)
    if err != nil {
        if os.IsExist(err) {
            os.Remove(dockerAuthCertPath)
        } else {
            return errors.WithStack(err)
        }
    }

    srcFile, err := os.Open(SlaveAuthCertFileName)
    if err != nil {
        return errors.WithStack(err)
    }
    defer srcFile.Close()

    destFile, err := os.Create(dockerAuthCertFile) // creates if file doesn't exist
    if err != nil {
        return errors.WithStack(err)
    }
    defer destFile.Close()

    _, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
    if err != nil {
        return errors.WithStack(err)
    }

    err = destFile.Sync()
    if err != nil {
        return errors.WithStack(err)
    }

    err = os.Chmod(dockerAuthCertFile, cert_file_permission)
    if err != nil {
        return errors.WithStack(err)
    }
    return nil
}