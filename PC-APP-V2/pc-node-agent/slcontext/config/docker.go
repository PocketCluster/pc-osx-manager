package config

import (
    "io"
    "io/ioutil"
    "os"
    "path"
    "path/filepath"

    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
)

const (
    cert_path_permission os.FileMode    = os.FileMode(0755)
    cert_file_permission os.FileMode    = os.FileMode(0644)

    DOCKER_ENV_PATH string              = "/etc/default/"
    DOCKER_ENV_FILE string              = DOCKER_ENV_PATH + "docker"

    SYSTEM_AUTH_CERT_NATIVE_FILE string = "/etc/ssl/certs/ca-certificates.crt"
    SYSTEM_AUTH_CERT_BACKUP_PATH string = "/etc/pocket/backup/"
    SYSTEM_AUTH_CERT_BACKUP_FILE string = SYSTEM_AUTH_CERT_BACKUP_PATH + "ca-certificates" + pcrypto.FileExtCertificate

    // CUSTOM CERTIFATE CONFIG
    CUSTOM_CERT_AUTH_PATH        string = "/usr/local/share/ca­certificates/"
    // the key file name should be "pc_node_engine_auth.acr". Due to accepted extention ".crt", we'll use custom name
    CUSTOM_CERT_AUTH_FILE        string = "pc_node_engine_auth.crt"
)

func dockerEnvContent() []byte {
    return []byte(`# PocketCluster Docker Upstart and SysVinit configuration file

DOCKER_OPTS="-H tcp://0.0.0.0:2376 --dns 127.0.0.1 --tlsverify --tlscacert=/etc/pocket/pki/pc_node_engine.acr --tlscert=/etc/pocket/pki/pc_node_engine.crt --tlskey=/etc/pocket/pki/pc_node_engine.pem --cluster-advertise=eth0:2376 --cluster-store=etcd://pc-master:2379 --cluster-store-opt kv.cacertfile=/etc/pocket/pki/pc_node_engine.acr --cluster-store-opt kv.certfile=/etc/pocket/pki/pc_node_engine.crt --cluster-store-opt kv.keyfile=/etc/pocket/pki/pc_node_engine.pem"`)
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

// Setup system cert for docker to connect registry
// This function assumes that docker auth certificate exists, and we just need to append it to system certs collection
func AppendAuthCertFowardSystemCertAuthority(rootPath string) error {
    var (
        slaveAuthCertFile string        = FilePathSlaveEngineAuthCert(rootPath)
        systemAuthCertNativeFile string = path.Join(rootPath, SYSTEM_AUTH_CERT_NATIVE_FILE)
        systemAuthCertBackupPath string = path.Join(rootPath, SYSTEM_AUTH_CERT_BACKUP_PATH)
        systemAuthCertBackupFile string = path.Join(rootPath, SYSTEM_AUTH_CERT_BACKUP_FILE)

        systemCert, dockerAuthCert, updatedCert []byte
        err error = nil
    )
    if !path.IsAbs(systemAuthCertNativeFile) {
        return errors.Errorf("[ERR] invalid root path")
    }
    // following two files should exist
    _, err = os.Stat(slaveAuthCertFile)
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
    dockerAuthCert, err = ioutil.ReadFile(slaveAuthCertFile)
    if err != nil {
        return errors.WithStack(err)
    }
    // concatenate the two
    updatedCert = append(systemCert, dockerAuthCert...)
    // write it to system certificate
    err = ioutil.WriteFile(systemAuthCertNativeFile, updatedCert, os.FileMode(0644))
    return errors.WithStack(err)
}

// --- Custom Cert Path and File --- //
func dirPathCustomCertAuth(rootPath string) string {
    return filepath.Join(rootPath, CUSTOM_CERT_AUTH_PATH)
}

func filePathCustomCertAuth(rootPath string) string {
    return filepath.Join(dirPathCustomCertAuth(rootPath), CUSTOM_CERT_AUTH_FILE)
}

func CopyCertAuthForwardCustomCertStorage(rootPath string) error {
    var (
        slaveAuthCertFile       string = FilePathSlaveEngineAuthCert(rootPath)
        // /usr/local/share/ca­certificates/
        slaveCustomCertAuthPath string = "\x2f\x75\x73\x72\x2f\x6c\x6f\x63\x61\x6c\x2f\x73\x68\x61\x72\x65\x2f\x63\x61\x2d\x63\x65\x72\x74\x69\x66\x69\x63\x61\x74\x65\x73\x2f"
        // /usr/local/share/ca­certificates/pc_node_engine_auth.crt
        slaveCustomCertAuthFile string = "\x2f\x75\x73\x72\x2f\x6c\x6f\x63\x61\x6c\x2f\x73\x68\x61\x72\x65\x2f\x63\x61\x2d\x63\x65\x72\x74\x69\x66\x69\x63\x61\x74\x65\x73\x2f\x70\x63\x5f\x6e\x6f\x64\x65\x5f\x65\x6e\x67\x69\x6e\x65\x5f\x61\x75\x74\x68\x2e\x63\x72\x74"
    )
    // original cert auth should exist
    if _, err := os.Stat(slaveAuthCertFile); os.IsNotExist(err) {
        return errors.WithStack(err)
    }
    // check if custom cert auth storage path exists. if DNE, build the full-path
    if _, err := os.Stat(slaveCustomCertAuthPath); os.IsNotExist(err) {
        if err := os.MkdirAll(slaveCustomCertAuthPath, cert_path_permission); err != nil {
            return errors.WithStack(err)
        }
    }
    // copy cert auth file to custom location so next upgrade will not mess up the cert system
    return copyFile(slaveAuthCertFile, slaveCustomCertAuthFile)
}
