package pkg

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    "gopkg.in/vmihailenco/msgpack.v2"
)

// --- Version ---
const (
    VBoxCoreVersion         string = "1.0.4"
)

type VBoxCoreState int
const (
    VBoxCoreBindBroken      VBoxCoreState = iota
    VBoxCoreBounded
)

func (s VBoxCoreState) String() string {
    var state string
    switch s {
        case VBoxCoreBounded:
            return "VBoxCoreBounded"
        case VBoxCoreBindBroken:
            return "VBoxCoreBindBroken"
    }
    return state
}

// --- Meta Field ---
const (
    VBC_PROTOCOL_VERSION    string = "c_pv"
    VBC_ENCRYPTED_PKG       string = "c_ep"
    VBC_PUBLIC_KEY          string = "c_pk"
    VBC_CRYPTO_SIGNATURE    string = "c_cs"
)

type VBoxCoreMeta struct {
    ProtocolVersion         string             `msgpack:"c_pv"`
    EncryptedPackage        []byte             `msgpack:"c_ep, inline, omitempty"`
    CryptoSignature         []byte             `msgpack:"c_cs, omitempty"`
    CoreStatus              *VBoxCoreStatus    `msgpack:"-"`
}

// --- Status Field ---
const (
    VBC_CORE_STATE          string = "c_st"
    VBC_CLUSTER_ID          string = "c_ci"
    VBC_EXT_IP4_ADDR_SMASK  string = "c_i4"
    VBC_EXT_IP4_GATEWAY     string = "c_g4"
    VBC_TIMESTAMP           string = "c_ts"
)

type VBoxCoreStatus struct {
    CoreState               VBoxCoreState      `msgpack:"c_st"`
    ClusterID               string             `msgpack:"c_ci"`
    ExtIP4AddrSmask         string             `msgpack:"c_i4"`
    ExtIP4Gateway           string             `msgpack:"c_g4"`
    TimeStamp               time.Time          `msgpack:"c_ts"`
}

// --- Compositions ---
func CorePackingBoundedStatus(clusterID, extAddr, extGateway string, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    return corePackingStatus(VBoxCoreBounded, clusterID, extAddr, extGateway, rsaEncryptor)
}

func CorePackingBindBrokenStatus(clusterID, extAddr, extGateway string, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    return corePackingStatus(VBoxCoreBindBroken, clusterID, extAddr, extGateway, rsaEncryptor)
}

func corePackingStatus(state VBoxCoreState, clusterID, extAddr, extGateway string, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    var (
        meta *VBoxCoreMeta = nil
        spkg, epkg, mpkg []byte = nil, nil, nil
        sig pcrypto.Signature
        err error = nil

        status = &VBoxCoreStatus {
            CoreState:          state,
            ClusterID:          clusterID,
            ExtIP4AddrSmask:    extAddr,
            ExtIP4Gateway:      extGateway,
            TimeStamp:          time.Now(),
        }
    )

    // error check
    if len(clusterID) == 0 {
        return nil, errors.Errorf("[ERR] core cluster id cannot be empty")
    }
    if len(extAddr) == 0 {
        return nil, errors.Errorf("[ERR] core status external address cannot be empty")
    }
    // TODO : this can really be zero. find out more cases
    if len(extGateway) == 0 {
        return nil, errors.Errorf("[ERR] core status external gateway cannot be empty")
    }

    // packaging status
    spkg, err = msgpack.Marshal(status)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // encrypt status package
    if rsaEncryptor == nil {
        return nil, errors.Errorf("[ERR] core status RSA Encryptor cannot be null")
    }
    epkg, sig, err = rsaEncryptor.EncryptByRSA(spkg)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // meta message packing
    meta = &VBoxCoreMeta{
        ProtocolVersion:     VBoxCoreVersion,
        EncryptedPackage:    epkg,
        CryptoSignature:     sig,
    }
    mpkg, err = msgpack.Marshal(meta)
    return mpkg, errors.WithStack(err)
}

func CoreUnpackingStatus(clusterID string, metaPackage []byte, rsaDecryptor pcrypto.RsaDecryptor) (*VBoxCoreMeta, error) {
    var (
        meta *VBoxCoreMeta = nil
        status *VBoxCoreStatus = nil
        spkg []byte = nil
        err error = nil
    )
    // error check
    if len(clusterID) == 0 {
        return nil, errors.Errorf("[ERR] core cluster id cannot be empty")
    }
    if len(metaPackage) == 0 {
        return nil, errors.Errorf("[ERR] meta package cannot be null")
    }

    // unpack meta package
    err = msgpack.Unmarshal(metaPackage, &meta)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    if meta == nil {
        return nil, errors.Errorf("[ERR] null unpacked meta")
    }
    if meta.ProtocolVersion != VBoxCoreVersion {
        return nil, errors.Errorf("[ERR] invalid protocol version")
    }

    // error check
    if rsaDecryptor == nil {
        return nil, errors.Errorf("[ERR] core status RSA Decryptor cannot be null")
    }
    if len(meta.EncryptedPackage) == 0 {
        return nil, errors.Errorf("[ERR] null encrypted status")
    }
    if len(meta.CryptoSignature) == 0 {
        return nil, errors.Errorf("[ERR] null crypto signature")
    }

    // decrypt status package
    spkg, err = rsaDecryptor.DecryptByRSA(meta.EncryptedPackage, meta.CryptoSignature)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // unpack status
    err = msgpack.Unmarshal(spkg, &status)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    if status == nil {
        return nil, errors.Errorf("[ERR] null unpacked status")
    }
    if len(status.ClusterID) == 0 {
        return nil, errors.Errorf("[ERR] invalid cluster id")
    }
    if status.ClusterID != clusterID {
        return nil, errors.Errorf("[ERR] invalid status from unkown core node")
    }
    if len(status.ExtIP4AddrSmask) == 0 {
        return nil, errors.Errorf("[ERR] invalid ip & subnet mask")
    }
    // TODO : this can really be zero. find out more cases
    if len(status.ExtIP4Gateway) == 0 {
        return nil, errors.Errorf("[ERR] invalid gateway")
    }

    // assign status
    meta.CoreStatus = status

    return meta, nil
}