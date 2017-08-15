package defaults

import (
    "github.com/stkim1/pcrypto"
)

const (
    ApplicationVersion         string = "0.1.4"

    ApplicationExpirationDate  string = "2017/12/31 23:55:59 -0000"
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
    UserDataPath               string = ".pocket"

    RepositoryPathPostfix      string = "repository"

    StoragePathPostfix         string = "storage"

    VirtualMachinePath         string = "pc-core"
)

const (
    PocketTimeDateFormat       string = "2006/01/02 15:04:05 -0700"
)

const (
    PocketClusterCoreName      string = "pc-core"
)

const (
    VBoxDefualtCoreDiskName    string = "pc-core-hdd.vmdk"
    // 128 GB
    VBoxDefualtCoreDiskSize    uint   = 128000
)