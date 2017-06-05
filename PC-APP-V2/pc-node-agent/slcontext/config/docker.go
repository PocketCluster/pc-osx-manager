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
    cert_path_permission os.FileMode    = os.FileMode(0755)
    cert_file_permission os.FileMode    = os.FileMode(0644)

    DOCKER_ENV_PATH string              = "/etc/default/"
    DOCKER_ENV_FILE string              = DOCKER_ENV_PATH + "docker"

    DOCKER_AUTH_CERT_PATH string        = "/etc/docker/certs.d/pc-master/"
    DOCKER_AUTH_CERT_FILE string        = DOCKER_AUTH_CERT_PATH + "ca" + pcrypto.FileExtCertificate

    SYSTEM_AUTH_CERT_NATIVE_FILE string = "/etc/ssl/certs/ca-certificates.crt"
    SYSTEM_AUTH_CERT_BACKUP_PATH string = slave_config_dir + "backup/"
    SYSTEM_AUTH_CERT_BACKUP_FILE string = slave_config_dir + "ca-certificates" + pcrypto.FileExtCertificate
)

func dockerEnvContent() []byte {
    return []byte(`# PocketCluster Docker Upstart and SysVinit configuration file

DOCKER_OPTS="-H tcp://0.0.0.0:2375 --dns 127.0.0.1 --tlsverify --tlscacert=/etc/pocket/pki/pc_cert_auth.crt --tlscert=/etc/pocket/pki/pc_node_engine.crt --tlskey=/etc/pocket/pki/pc_node_engine.pem --cluster-advertise=eth0:2375 --cluster-store=etcd://pc-master:2379 --cluster-store-opt kv.cacertfile=/etc/pocket/pki/pc_cert_auth.crt --cluster-store-opt kv.certfile=/etc/pocket/pki/pc_node_engine.crt --cluster-store-opt kv.keyfile=/etc/pocket/pki/pc_node_engine.pem"`)
}

func copyFile(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return errors.WithStack(err)
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return errors.WithStack(err)
    }
    defer dstFile.Close()

    // check first var for number of bytes copied
    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        return errors.WithStack(err)
    }

    err = dstFile.Sync()
    if err != nil {
        return errors.WithStack(err)
    }

    err = os.Chmod(dst, cert_file_permission)
    if err != nil {
        return errors.WithStack(err)
    }
    return nil
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
/*
    TODO : this isn't necessary unless you want to make sure there only is one problem, DNE.
    _, err = os.Stat(dockerEnvFile)
    if err != nil && !os.IsNotExist(err) {
        return errors.WithStack(err)
    }
*/
    os.Remove(dockerEnvFile)

    err = ioutil.WriteFile(dockerEnvFile, dockerEnvContent(), cert_file_permission)
    return errors.WithStack(err)
}

func SetupDockerAuthorityCert(rootPath string) error {
    var (
        slaveAuthCertFile string  = path.Join(rootPath, SlaveAuthCertFileName)
        dockerAuthCertPath string = path.Join(rootPath, DOCKER_AUTH_CERT_PATH)
        dockerAuthCertFile string = path.Join(rootPath, DOCKER_AUTH_CERT_FILE)
        err error = nil
    )
    if !path.IsAbs(dockerAuthCertPath) {
        return errors.Errorf("[ERR] invalid root path")
    }
    _, err = os.Stat(slaveAuthCertFile)
    if err != nil {
        return errors.WithStack(err)
    }
    _, err = os.Stat(dockerAuthCertPath)
    if err != nil {
        if os.IsNotExist(err) {
            os.MkdirAll(dockerAuthCertPath, cert_path_permission);
        } else {
            return errors.WithStack(err)
        }
    }
/*
    TODO : this isn't necessary unless you want to make sure there only is one problem, DNE.
    _, err = os.Stat(dockerAuthCertFile)
    if err != nil && !os.IsNotExist(err) {
        return errors.WithStack(err)
    }
*/
    os.Remove(dockerAuthCertFile)

    return copyFile(slaveAuthCertFile, dockerAuthCertFile)
}

// Setup system cert for docker to connect registry
// This function assumes that docker auth certificate exists, and we just need to append it to system certs collection
func AppendAuthCertFowardSystemCertAuthority(rootPath string) error {
    var (
        systemAuthCertNativeFile string = path.Join(rootPath, SYSTEM_AUTH_CERT_NATIVE_FILE)
        systemAuthCertBackupPath string = path.Join(rootPath, SYSTEM_AUTH_CERT_BACKUP_PATH)
        systemAuthCertBackupFile string = path.Join(rootPath, SYSTEM_AUTH_CERT_BACKUP_FILE)
        dockerAuthCertFile       string = path.Join(rootPath, DOCKER_AUTH_CERT_FILE)

        systemCert, dockerAuthCert, updatedCert []byte
        err error = nil
    )
    if !path.IsAbs(systemAuthCertNativeFile) {
        return errors.Errorf("[ERR] invalid root path")
    }
    // following two files should exist
    _, err = os.Stat(dockerAuthCertFile)
    if err != nil {
        return errors.WithStack(err)
    }
    _, err = os.Stat(systemAuthCertNativeFile)
    if err != nil {
        return errors.WithStack(err)
    }
    // check if backup has appropriate folder. If DNE make one.
    _, err = os.Stat(systemAuthCertBackupPath)
    if err != nil {
        if os.IsNotExist(err) {
            os.MkdirAll(systemAuthCertBackupPath, cert_path_permission);
        } else {
            return errors.WithStack(err)
        }
    }
    // check if system certificate back up exists. If DNE, copy the original to backup location
    _, err = os.Stat(systemAuthCertBackupFile)
    if err != nil {
        if os.IsNotExist(err) {
            err = copyFile(systemAuthCertNativeFile, systemAuthCertBackupFile)
            if err != nil {
                return err
            }
        } else {
            return errors.WithStack(err)
        }
    }

    // read backed-up native certificate
    systemCert, err = ioutil.ReadFile(systemAuthCertBackupFile)
    if err != nil {
        return errors.WithStack(err)
    }
    // read downloaded docker certificate
    dockerAuthCert, err = ioutil.ReadFile(dockerAuthCertFile)
    if err != nil {
        return errors.WithStack(err)
    }
    // concatenate the two
    updatedCert = append(systemCert, dockerAuthCert...)
    // write it to system certificate
    err = ioutil.WriteFile(systemAuthCertNativeFile, updatedCert, os.FileMode(0644))
    return errors.WithStack(err)
}