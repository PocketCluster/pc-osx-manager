package config

import (
    "os"
    "io/ioutil"

    "github.com/pkg/errors"
)

const (
    DOCKER_ENV_PATH string = "/etc/default/"
    DOCKER_ENV_FILE string = DOCKER_ENV_PATH + "docker"
)

func dockerEnvContent() []byte {
    return []byte(`# PocketCluster Docker Upstart and SysVinit configuration file

DOCKER_OPTS="-H tcp://0.0.0.0:2375 --dns 127.0.0.1 --tlsverify --tlscacert=/etc/pocket/pki/pc_cert_auth.crt --tlscert=/etc/pocket/pki/pc_node_engine.crt --tlskey=/etc/pocket/pki/pc_node_engine.pem --cluster-advertise=eth0:2375 --cluster-store=etcd://pc-master:2379 --cluster-store-opt kv.cacertfile=/etc/pocket/pki/pc_cert_auth.crt --cluster-store-opt kv.certfile=/etc/pocket/pki/pc_node_engine.crt --cluster-store-opt kv.keyfile=/etc/pocket/pki/pc_node_engine.pem"`)
}

func SetupDockerEnvironement(rootPath string) error {
    var (
        dockerEnvPath string = rootPath + DOCKER_ENV_PATH
        dockerEnvFile string = rootPath + DOCKER_ENV_FILE
        err error = nil
    )
    if _, err = os.Stat(dockerEnvPath); os.IsNotExist(err) {
        os.MkdirAll(dockerEnvPath, 0755);
    }

    if _, err = os.Stat(dockerEnvFile); os.IsExist(err) {
        os.Remove(dockerEnvFile)
    }

    err = ioutil.WriteFile(dockerEnvFile, dockerEnvContent(), os.FileMode(0644))
    return errors.WithStack(err)
}