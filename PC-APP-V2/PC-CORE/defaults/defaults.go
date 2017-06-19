package defaults

import (
    "github.com/stkim1/pcrypto"
)

const (
    ClusterCertAuthPrivateKey string  = "pc_cert_auth"   + pcrypto.FileExtPrivateKey

    ClusterCertAuthPublicKey string   = "pc_cert_auth"   + pcrypto.FileExtPublicKey

    ClusterCertAuthCertificate string = "pc_cert_auth"   + pcrypto.FileExtCertificate

    ClusterCertAuthSshCheck string    = "pc_cert_auth"   + pcrypto.FileExtSSHCertificate

    MasterHostPrivateKey string       = "pc_master_host" + pcrypto.FileExtPrivateKey

    MasterHostPublicKey string        = "pc_master_host" + pcrypto.FileExtPublicKey

    MasterHostCertificate string      = "pc_master_host" + pcrypto.FileExtCertificate

    MasterHostSshKey string           = "pc_master_host" + pcrypto.FileExtSSHCertificate

    MasterBeaconPrivateKey string     = "pc_master_beacon" + pcrypto.FileExtPrivateKey

    MasterBeaconPublicKey string      = "pc_master_beacon" + pcrypto.FileExtPublicKey
)

const (
    RepositoryPathPostfix string      = "repository"

    StoragePathPostfix string         = "storage"
)

const (
    PocketTimeDateFormat string       = "2006/01/02 15:04:05 -0700"
)