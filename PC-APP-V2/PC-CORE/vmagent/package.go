package vmagent

import (
    "time"

    "github.com/pkg/errors"
    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/stkim1/pcrypto"
)

// --- Version ---
const (
    VBoxMasterVersion        string = "1.0.0"
)

// --- Meta Field ---
const (
    VBM_PROTOCOL_VERSION     string = "m_pv"
    VBM_ENCRYPTED_PKG        string = "m_ep"
    VBM_MASTER_PUBKEY        string = "m_pk"
    VBM_CRYPTO_SIGNATURE     string = "m_cs"
)

type VBoxMasterAgentMeta struct {
    ProtocolVersion          string    `msgpack:"m_pv"`
    EncryptedPackage         []byte    `msgpack:"m_ep, inline, omitempty"`
    PublicKey                []byte    `msgpack:"m_pk, omitempty"`
    CryptoSignature          []byte    `msgpack:"m_cs, omitempty"`
}

// --- Acknowledge Field ---
const (
    VBM_CORE_UUID            string = "m_cu"
    VBM_TIMESTAMP            string = "m_ts"
)

type VBoxMasterAcknowledge struct {
    CoreUUID                 string    `msgpack:"m_cu, omitempty"`
    TimeStamp                time.Time `msgpack:"m_ts"`
}

// --- Compositions ---
func MasterEncryptedKeyExchange(coreUUID string, pubkey []byte, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    var (
        ack = &VBoxMasterAcknowledge {
            CoreUUID:     coreUUID,
            TimeStamp:    time.Now(),
        }
        err error = nil
    )
    if len(coreUUID) == 0 {
        return nil, errors.Errorf("[ERR] invalid core uuid assignment")
    }
    if len(pubkey) == 0 {
        return nil, errors.Errorf("[ERR] invalid public key passed")
    }
    if rsaEncryptor == nil {
        return nil, errors.Errorf("[ERR] invalid encryptor passed")
    }

    // package acknowledge
    apkg, err := msgpack.Marshal(ack)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // encrypt ack package
    epkg, sig, err := rsaEncryptor.EncryptByRSA(apkg)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // meta message packing
    meta := &VBoxMasterAgentMeta {
        ProtocolVersion:     VBoxMasterVersion,
        EncryptedPackage:    epkg,
        PublicKey:           pubkey,
        CryptoSignature:     sig,
    }
    mpkg, err := msgpack.Marshal(meta)
    return mpkg, errors.WithStack(err)
}

func MasterDecryptedKeyExchange(metaPackage, prvkey []byte) (*VBoxMasterAcknowledge, pcrypto.RsaDecryptor, error) {
    var (
        meta *VBoxMasterAgentMeta = nil
        ack *VBoxMasterAcknowledge = nil
        err error = nil
    )

    // unpack meta
    err = msgpack.Unmarshal(metaPackage, &meta)
    if err != nil {
        return nil, nil, errors.WithStack(err)
    }
    if meta == nil {
        return nil, nil, errors.Errorf("[ERR] null unpacked meta")
    }
    if meta.ProtocolVersion != VBoxMasterVersion {
        return nil, nil, errors.Errorf("[ERR] incorrect protocol version")
    }
    if len(meta.EncryptedPackage) == 0 {
        return nil, nil, errors.Errorf("[ERR] null encrypted ack")
    }
    if len(meta.PublicKey) == 0 {
        return nil, nil, errors.Errorf("[ERR] null public key")
    }
    if len(meta.CryptoSignature) == 0 {
        return nil, nil, errors.Errorf("[ERR] null crypto signature")
    }

    // build decryptor
    rsaDecrypto, err := pcrypto.NewRsaDecryptorFromKeyData(meta.PublicKey, prvkey)
    if err != nil {
        return nil, nil, errors.Errorf("[ERR] cannot build decryptor")
    }

    // decrypt message
    apkg, err := rsaDecrypto.DecryptByRSA(meta.EncryptedPackage, meta.CryptoSignature)
    if err != nil {
        return nil, nil, errors.Errorf("[ERR] cannot build decryptor")
    }

    // unpack acknowledge
    err = msgpack.Unmarshal(apkg, &ack)
    if err != nil {
        return nil, nil, errors.WithStack(err)
    }
    if ack == nil {
        return nil, nil, errors.Errorf("[ERR] null unpacked acknowledge")
    }
    if len(ack.CoreUUID) == 0 {
        return nil, nil, errors.Errorf("[ERR] invalid core uuid assignment")
    }

    return ack, rsaDecrypto, nil
}

func MasterEncryptedBounded(rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    var (
        ack = &VBoxMasterAcknowledge {
            TimeStamp:    time.Now(),
        }
        err error = nil
    )

    // package acknowledge
    apkg, err := msgpack.Marshal(ack)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // encrypt ack package
    epkg, sig, err := rsaEncryptor.EncryptByRSA(apkg)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // meta message packing
    meta := &VBoxMasterAgentMeta {
        ProtocolVersion:     VBoxMasterVersion,
        EncryptedPackage:    epkg,
        CryptoSignature:     sig,
    }
    mpkg, err := msgpack.Marshal(meta)
    return mpkg, errors.WithStack(err)
}

func MasterDecryptedBounded(metaPackage []byte, rsaDecryptor pcrypto.RsaDecryptor) (*VBoxMasterAcknowledge, error) {
    var (
        meta *VBoxMasterAgentMeta = nil
        ack *VBoxMasterAcknowledge = nil
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
    if meta.ProtocolVersion != VBoxMasterVersion {
        return nil, errors.Errorf("[ERR] incorrect protocol version")
    }
    if len(meta.EncryptedPackage) == 0 {
        return nil, errors.Errorf("[ERR] null encrypted acknowledge")
    }
    if len(meta.PublicKey) != 0 {
        return nil, errors.Errorf("[ERR] invalid meta package content w/ pubkey")
    }
    if len(meta.CryptoSignature) == 0 {
        return nil, errors.Errorf("[ERR] null crypto signature")
    }

    // decrypt ack package
    apkg, err := rsaDecryptor.DecryptByRSA(meta.EncryptedPackage, meta.CryptoSignature)
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
    if len(ack.CoreUUID) != 0 {
        return nil, errors.Errorf("[ERR] invalid ack content w/ core uuid")
    }

    return ack, nil
}