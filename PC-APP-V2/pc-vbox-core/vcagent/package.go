package vcagent

import (
    "time"

    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
)

// --- Version ---
const (
    VBoxCoreVersion         string = "1.0.0"
)

// --- Meta Field ---
const (
    VBC_PROTOCOL_VERSION    string = "c_pv"
    VBC_ENCRYPTED_PKG       string = "c_ep"
    VBC_PUBLIC_KEY          string = "c_pk"
    VBC_CRYPTO_SIGNATURE    string = "c_cs"
)

type VBoxCoreAgentMeta struct {
    ProtocolVersion         string    `msgpack:"c_pv"`
    EncryptedPackage        []byte    `msgpack:"c_ep, inline, omitempty"`
    PublicKey               []byte    `msgpack:"c_pk, omitempty"`
    CryptoSignature         []byte    `msgpack:"c_cs, omitempty"`
}

// --- Status Field ---
const (
    VBC_EXT_IP4_ADDR_SMASK  string = "c_i4"
    VBC_EXT_IP4_GATEWAY     string = "c_g4"
    VBC_TIMESTAMP           string = "c_ts"
)

type VBoxCoreStatus struct {
    ExtIP4AddrSmask         string    `msgpack:"c_i4"`
    ExtIP4Gateway           string    `msgpack:"c_g4"`
    TimeStamp               time.Time `msgpack:"c_ts"`
}

// --- Compositions ---
func CorePackingUnbounded(pubkey []byte) ([]byte, error) {
    var (
        meta = &VBoxCoreAgentMeta {
            ProtocolVersion:    VBoxCoreVersion,
            PublicKey:          pubkey,
        }
        err error = nil
    )
    pkg, err := msgpack.Marshal(meta)
    return pkg, errors.WithStack(err)
}

func CoreUnpackingUnbounded(metaPackage []byte) (*VBoxCoreAgentMeta, error) {
    var (
        meta *VBoxCoreAgentMeta
        err error = nil
    )
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
    if len(meta.EncryptedPackage) != 0 {
        return nil, errors.Errorf("[ERR] invalid meta package content w/ encrypted status")
    }
    if len(meta.PublicKey) == 0 {
        return nil, errors.Errorf("[ERR] null pubkey in meta package")
    }
    if len(meta.CryptoSignature) != 0 {
        return nil, errors.Errorf("[ERR] invalid meta package content w/ crypto signature")
    }

    return meta, nil
}

func CoreEncryptedBounded(extAddr, extGateway string, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    var (
        status = &VBoxCoreStatus {
            ExtIP4AddrSmask:    extAddr,
            ExtIP4Gateway:      extGateway,
            TimeStamp:          time.Now(),
        }
        err error = nil
    )
    // packaging status
    spkg, err := msgpack.Marshal(status)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // encrypt status package
    epkg, sig, err := rsaEncryptor.EncryptByRSA(spkg)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // meta message packing
    meta := &VBoxCoreAgentMeta {
        ProtocolVersion:     VBoxCoreVersion,
        EncryptedPackage:    epkg,
        CryptoSignature:     sig,
    }
    mpkg, err := msgpack.Marshal(meta)
    return mpkg, errors.WithStack(err)
}

func CoreDecryptBounded(metaPackage []byte, rsaDecryptor pcrypto.RsaDecryptor) (*VBoxCoreStatus, error) {
    var (
        meta *VBoxCoreAgentMeta = nil
        status *VBoxCoreStatus = nil
        err error = nil
    )

    // unpack meta
    err = msgpack.Unmarshal(metaPackage, &meta)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    if meta == nil {
        return nil, errors.Errorf("[ERR] null unpacked meta")
    }
    if meta.ProtocolVersion != VBoxCoreVersion {
        return nil, errors.Errorf("[ERR] incorrect protocol version")
    }
    if len(meta.EncryptedPackage) == 0 {
        return nil, errors.Errorf("[ERR] null encrypted status")
    }
    if len(meta.PublicKey) != 0 {
        return nil, errors.Errorf("[ERR] invalid meta package content w/ pubkey")
    }
    if len(meta.CryptoSignature) == 0 {
        return nil, errors.Errorf("[ERR] null crypto signature")
    }

    // decrypt status package
    spkg, err := rsaDecryptor.DecryptByRSA(meta.EncryptedPackage, meta.CryptoSignature)
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
    if len(status.ExtIP4AddrSmask) == 0 {
        return nil, errors.Errorf("[ERR] invalid ip & subnet mask")
    }
    if len(status.ExtIP4Gateway) == 0 {
        return nil, errors.Errorf("[ERR] invalid gateway")
    }

    return status, nil
}