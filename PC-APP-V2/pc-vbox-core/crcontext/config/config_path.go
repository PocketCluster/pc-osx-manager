package config

import (
    "path/filepath"

    "github.com/stkim1/pcrypto"
)

// ------ CONFIGURATION FILES ------
const (
    // config directory
    dir_core_config             string = "/etc/pocket/"

    // core config file
    core_config_file            string = "core.conf.yaml"
    core_cluster_id_file        string = "cluster.id"
    core_ssh_auth_token_file    string = "ssh.auth.token"
    core_user_name_file         string = "core.user.name"
    core_user_uid_file          string = "core.user.uid"
    core_user_gid_file          string = "core.user.gid"

    // cert directory
    dir_core_certs              string = "pki"

    // these files are 2048 RSA crypto files used to join network
    core_vbox_public_Key_file   string = "pc_core_vbox" + pcrypto.FileExtPublicKey
    core_vbox_prvate_Key_file   string = "pc_core_vbox" + pcrypto.FileExtPrivateKey
    master_vbox_public_Key_file string = "pc_master_vbox" + pcrypto.FileExtPublicKey

    // these files are 2048 RSA crypto files used for Docker & Registry
    core_engine_auth_cert_file  string = "pc_core_engine" + pcrypto.FileExtAuthCertificate
    core_engine_key_cert_file   string = "pc_core_engine" + pcrypto.FileExtCertificate
    core_engine_prvate_key_file string = "pc_core_engine" + pcrypto.FileExtPrivateKey

    // these are files used for teleport certificate
    core_ssh_key_cert_file      string = "pc_core_ssh" + pcrypto.FileExtSSHCertificate
    core_ssh_private_key_file   string = "pc_core_ssh" + pcrypto.FileExtPrivateKey
)

// --- to read config ---
func DirPathCoreConfig(rootPath string) string {
    return filepath.Join(rootPath, dir_core_config)
}

func FilePathCoreConfig(rootPath string) string {
    return filepath.Join(DirPathCoreConfig(rootPath), core_config_file)
}

func FilePathClusterID(rootPath string) string {
    return filepath.Join(DirPathCoreConfig(rootPath), core_cluster_id_file)
}

func FilePathAuthToken(rootPath string) string {
    return filepath.Join(DirPathCoreConfig(rootPath), core_ssh_auth_token_file)
}

// --- To read config ---
func DirPathCoreCerts(rootPath string) string {
    return filepath.Join(DirPathCoreConfig(rootPath), dir_core_certs)
}

func FilePathCoreVboxPublicKey(rootPath string) string {
    return filepath.Join(DirPathCoreCerts(rootPath), core_vbox_public_Key_file)
}

func FilePathCoreVboxPrivateKey(rootPath string) string {
    return filepath.Join(DirPathCoreCerts(rootPath), core_vbox_prvate_Key_file)
}

func FilePathMasterVboxPublicKey(rootPath string) string {
    return filepath.Join(DirPathCoreCerts(rootPath), master_vbox_public_Key_file)
}

func FilePathCoreSSHKeyCert(rootPath string) string {
    return filepath.Join(DirPathCoreCerts(rootPath), core_ssh_key_cert_file)
}

func FilePathCoreSSHPrivateKey(rootPath string) string {
    return filepath.Join(DirPathCoreCerts(rootPath), core_ssh_private_key_file)
}

// --- to build tar archive file ---
func ArchivePathClusterID() string {
    return core_cluster_id_file
}

func ArchivePathAuthToken() string {
    return core_ssh_auth_token_file
}

func ArchivePathUserName() string {
    return core_user_name_file
}

func ArchivePathUserUID() string {
    return core_user_uid_file
}

func ArchivePathUserGID() string {
    return core_user_gid_file
}

func ArchivePathCertsDir() string {
    return dir_core_certs
}

func ArchivePathCoreVboxPublicKey() string {
    return filepath.Join(ArchivePathCertsDir(), core_vbox_public_Key_file)
}

func ArchivePathCoreVboxPrivateKey() string {
    return filepath.Join(ArchivePathCertsDir(), core_vbox_prvate_Key_file)
}

func ArchivePathMasterVboxPublicKey() string {
    return filepath.Join(ArchivePathCertsDir(), master_vbox_public_Key_file)
}

func ArchivePathCoreEngineAuthCert() string {
    return filepath.Join(ArchivePathCertsDir(), core_engine_auth_cert_file)
}

func ArchivePathCoreEngineKeyCert() string {
    return filepath.Join(ArchivePathCertsDir(), core_engine_key_cert_file)
}

func ArchivePathCoreEnginePrivateKey() string {
    return filepath.Join(ArchivePathCertsDir(), core_engine_prvate_key_file)
}
