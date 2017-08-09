package pkg

import (
    "time"

    "github.com/pkg/errors"
    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/stkim1/pcrypto"
)

// --- Version ---
const (
    VBoxMasterVersion        string = "1.0.4"
)

type VBoxMasterState int
const (
    VBoxMasterBindBroken     VBoxMasterState = iota
    VBoxMasterBounded
)

func (s VBoxMasterState) String() string {
    var state string
    switch s {
        case VBoxMasterBindBroken:
            state = "VBoxMasterBindBroken"
        case VBoxMasterBounded:
            state = "VBoxMasterBounded"
    }
    return state
}

// --- Meta Field ---
const (
    VBM_PROTOCOL_VERSION     string = "m_pv"
    VBM_ENCRYPTED_PKG        string = "m_ep"
    VBM_MASTER_PUBKEY        string = "m_pk"
    VBM_CRYPTO_SIGNATURE     string = "m_cs"
)

type VBoxMasterMeta struct {
    ProtocolVersion          string                    `msgpack:"m_pv"`
    EncryptedPackage         []byte                    `msgpack:"m_ep, inline, omitempty"`
    CryptoSignature          []byte                    `msgpack:"m_cs, omitempty"`
    MasterAcknowledge        *VBoxMasterAcknowledge    `msgpack:"-"`
}

// --- Acknowledge Field ---
const (
    VBM_MASTER_STATE         string = "m_st"
    VBM_CLUSTER_ID           string = "m_ci"
    VBM_AUTH_TOKEN           string = "m_at"
    VBM_EXT_IP4_ADDR         string = "m_e4"
    VBM_TIMESTAMP            string = "m_ts"
)

type VBoxMasterAcknowledge struct {
    MasterState              VBoxMasterState           `msgpack:"m_st"`
    ClusterID                string                    `msgpack:"m_ci, omitempty"`
    ExtIP4Addr               string                    `msgpack:"m_e4, omitempty"`
    TimeStamp                time.Time                 `msgpack:"m_ts"`
}

// --- Compositions ---
func MasterPackingBoundedAcknowledge(clusterID, extIP4Addr string, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    return masterPackingAcknowledge(VBoxMasterBounded, clusterID, extIP4Addr, rsaEncryptor)
}

func MasterPackingBindBrokenAcknowledge(clusterID, extIP4Addr string, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    return masterPackingAcknowledge(VBoxMasterBindBroken, clusterID, extIP4Addr, rsaEncryptor)
}

func masterPackingAcknowledge(state VBoxMasterState, clusterID, extIP4Addr string, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    var (
        meta *VBoxMasterMeta = nil
        apkg, epkg, mpkg []byte = nil, nil, nil
        sig pcrypto.Signature = nil
        err error = nil

        ack = &VBoxMasterAcknowledge {
            MasterState:    state,
            ClusterID:      clusterID,
            ExtIP4Addr:     extIP4Addr,
            TimeStamp:      time.Now(),
        }
    )
    if len(clusterID) == 0 {
        return nil, errors.Errorf("[ERR] invalid cluster id assignment")
    }
    if len(extIP4Addr) == 0 {
        return nil, errors.Errorf("[ERR] invalid external master ip4 address")
    }
    if rsaEncryptor == nil {
        return nil, errors.Errorf("[ERR] master RSA Encryptor cannot be null")
    }

    // package acknowledge
    apkg, err = msgpack.Marshal(ack)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // encrypt ack package
    epkg, sig, err = rsaEncryptor.EncryptByRSA(apkg)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // meta message packing
    meta = &VBoxMasterMeta{
        ProtocolVersion:     VBoxMasterVersion,
        EncryptedPackage:    epkg,
        CryptoSignature:     sig,
    }
    mpkg, err = msgpack.Marshal(meta)
    return mpkg, errors.WithStack(err)
}

func MasterUnpackingAcknowledge(clusterID string, metaPackage []byte, rsaDecryptor pcrypto.RsaDecryptor) (*VBoxMasterMeta, error) {
    var (
        meta *VBoxMasterMeta = nil
        ack *VBoxMasterAcknowledge = nil
        apkg []byte = nil
        err error = nil
    )

    // unpack meta & error check
    if len(clusterID) == 0 {
        return nil, errors.Errorf("[ERR] invalid master cluster id assignment")
    }
    if len(metaPackage) == 0 {
        return nil, errors.Errorf("[ERR] meta package cannot be null")
    }
    err = msgpack.Unmarshal(metaPackage, &meta)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    if meta == nil {
        return nil, errors.Errorf("[ERR] null unpacked meta")
    }
    if meta.ProtocolVersion != VBoxMasterVersion {
        return nil, errors.Errorf("[ERR] incorrect protocol version")
    }
    if len(meta.EncryptedPackage) == 0 {
        return nil, errors.Errorf("[ERR] null encrypted acknowledge")
    }
    if len(meta.CryptoSignature) == 0 {
        return nil, errors.Errorf("[ERR] null crypto signature")
    }

    if rsaDecryptor == nil {
        return nil, errors.Errorf("[ERR] core RSA Decryptor cannot be null")
    }
    // decrypt ack package
    apkg, err = rsaDecryptor.DecryptByRSA(meta.EncryptedPackage, meta.CryptoSignature)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // unpack acknowledge
    err = msgpack.Unmarshal(apkg, &ack)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    if ack == nil {
        return nil, errors.Errorf("[ERR] null unpacked acknowledge")
    }
    if len(ack.ClusterID) == 0 {
        return nil, errors.Errorf("[ERR] invalid cluster id assignment")
    }
    if ack.ClusterID != clusterID {
        return nil, errors.Errorf("[ERR] invalid cluster id from unknown master")
    }
    if len(ack.ExtIP4Addr) == 0 {
        return nil, errors.Errorf("[ERR] invalid external master ip4 address")
    }

    // assing acknowledge
    meta.MasterAcknowledge = ack

    return meta, nil
}