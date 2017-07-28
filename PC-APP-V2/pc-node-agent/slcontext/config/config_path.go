package config

import (
    "path"
    "path/filepath"

    "github.com/stkim1/pcrypto"
)

// ------ CONFIGURATION FILES ------
const (
    // HOST GENERAL CONFIG
    dir_system_config           string = "/etc/"
    system_hostname_file        string = "hostname"
    system_timezone_file        string = "timezone"

    // POCKET SPECIFIC CONFIG
    dir_slave_config            string = "pocket"
    slave_config_file           string = "slave.conf.yaml"

    dir_slave_certs             string = "pki"
    // these files are 1024 RSA crypto files used to join network
    slave_public_Key_file       string = "pc_node_beacon"   + pcrypto.FileExtPublicKey
    slave_prvate_Key_file       string = "pc_node_beacon"   + pcrypto.FileExtPrivateKey
    master_public_Key_file      string = "pc_master_beacon" + pcrypto.FileExtPublicKey

    // these files are 2048 RSA crypto files used for Docker & Registry. This should be acquired from Teleport Auth server
    slave_auth_cert_file        string = "pc_node_engine"   + pcrypto.FileExtAuthCertificate
    slave_engine_cert_file      string = "pc_node_engine"   + pcrypto.FileExtCertificate
    slave_engine_key_file       string = "pc_node_engine"   + pcrypto.FileExtPrivateKey

    // these are files used for teleport certificate
    slave_ssh_key_cert_file     string = "pc_node_ssh"      + pcrypto.FileExtSSHCertificate
    slave_ssh_private_key_file  string = "pc_node_ssh"      + pcrypto.FileExtPrivateKey
)

// --- System Configuration --- //
func DirPathSystemConfig(rootPath string) string {
    return path.Join(rootPath, dir_system_config)
}

func FilePathSystemHostname(rootPath string) string {
    return filepath.Join(DirPathSystemConfig(rootPath), system_hostname_file)
}

func FilePathSystemTimezone(rootPath string) string {
    return filepath.Join(DirPathSystemConfig(rootPath), system_timezone_file)
}

// --- Slave Configuration --- //
func DirPathSlaveConfig(rootPath string) string {
    return path.Join(DirPathSystemConfig(rootPath), dir_slave_config)
}

func FilePathSlaveConfig(rootPath string) string {
    return filepath.Join(DirPathSlaveConfig(rootPath), slave_config_file)
}

// --- Slave Certificates --- //
func DirPathSlaveCerts(rootPath string) string {
    return path.Join(DirPathSlaveConfig(rootPath), dir_slave_certs)
}

func FilePathSlavePublicKey(rootPath string) string {
    return filepath.Join(DirPathSlaveCerts(rootPath), slave_public_Key_file)
}

func FilePathSlavePrivateKey(rootPath string) string {
    return filepath.Join(DirPathSlaveCerts(rootPath), slave_prvate_Key_file)
}

func FilePathMasterPublicKey(rootPath string) string {
    return filepath.Join(DirPathSlaveCerts(rootPath), master_public_Key_file)
}

func FilePathSlaveAuthCert(rootPath string) string {
    return filepath.Join(DirPathSlaveCerts(rootPath), slave_auth_cert_file)
}

func FilePathSlaveEngineKeyCert(rootPath string) string {
    return filepath.Join(DirPathSlaveCerts(rootPath), slave_engine_cert_file)
}

func FilePathSlaveEnginePrivateKey(rootPath string) string {
    return filepath.Join(DirPathSlaveCerts(rootPath), slave_engine_key_file)
}

func FilePathSlaveSSHKeyCert(rootPath string) string {
    return filepath.Join(DirPathSlaveCerts(rootPath), slave_ssh_key_cert_file)
}

func FilePathSlaveSSHPrivateKey(rootPath string) string {
    return filepath.Join(DirPathSlaveCerts(rootPath), slave_ssh_private_key_file)
}
