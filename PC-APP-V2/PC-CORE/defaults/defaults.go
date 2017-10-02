package defaults

import (
    "github.com/stkim1/pcrypto"
)

// PLACE ONLY *CONSTANT* VALUES W/O DEPENDENCY IN THIS PACKAGE, !!!PLEASE!!!

const (
    ApplicationVersion         string = "0.1.4"

    ApplicationExpirationDate  string = "2018/01/01 00:00:00 -0000"
)

const (
    ClusterCertAuthPrivateKey  string = "pc_cert_auth"   + pcrypto.FileExtPrivateKey

    ClusterCertAuthPublicKey   string = "pc_cert_auth"   + pcrypto.FileExtPublicKey

    ClusterCertAuthCertificate string = "pc_cert_auth"   + pcrypto.FileExtCertificate

    ClusterCertAuthSshCheck    string = "pc_cert_auth"   + pcrypto.FileExtSSHCertificate

    MasterHostPrivateKey       string = "pc_master_host" + pcrypto.FileExtPrivateKey

    MasterHostPublicKey        string = "pc_master_host" + pcrypto.FileExtPublicKey

    MasterHostCertificate      string = "pc_master_host" + pcrypto.FileExtCertificate

    MasterHostSshKey           string = "pc_master_host" + pcrypto.FileExtSSHCertificate

    MasterBeaconPrivateKey     string = "pc_master_beacon" + pcrypto.FileExtPrivateKey

    MasterBeaconPublicKey      string = "pc_master_beacon" + pcrypto.FileExtPublicKey

    MasterVBoxCtrlPrivateKey   string = "pc_master_vbox_ctrl" + pcrypto.FileExtPrivateKey

    MasterVBoxCtrlPublicKey    string = "pc_master_vbox_ctrl" + pcrypto.FileExtPublicKey
)

const (
    PathPostfixRepository      string = "repository"

    PathPostfixStorage         string = "storage"

    PathPostfixVirtualMachine  string = "pc-core"

    PathPostfixCoreNodeData    string = ".pocket-core-data"

    PathPostfixCoreDataInput   string = "PocketCluster"
)

const (
    PocketTimeDateFormat       string = "2006/01/02 15:04:05 -0700"
)

const (
    PocketClusterCoreName      string = "pc-core"
)

const (
    VBoxDefualtCoreDiskName    string = "pc-core-hdd.vmdk"
    // 46 GB
    //VBoxDefualtCoreDiskSize    uint   = 128000
    VBoxDefualtCoreDiskSize    uint   = 46000
)