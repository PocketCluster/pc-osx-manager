package pkg

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

type VBoxMasterState int
const (
    VBoxMasterUnbounded      VBoxMasterState = iota
    VBoxMasterKeyExchange
    VBoxMasterBounded
    VBoxMasterBindBroken
)

func (s VBoxMasterState) String() string {
    var state string
    switch s {
        case VBoxMasterUnbounded:
            state = "VBoxMasterUnbounded"
        case VBoxMasterKeyExchange:
            state = "VBoxMasterKeyExchange"
        case VBoxMasterBounded:
            state = "VBoxMasterBounded"
        case VBoxMasterBindBroken:
            state = "VBoxMasterBindBroken"
    }
    return state
}

// --- Meta Field ---
const (
    VBM_PROTOCOL_VERSION     string = "m_pv"
    VBM_MASTER_STATE         string = "m_st"
    VBM_ENCRYPTED_PKG        string = "m_ep"
    VBM_MASTER_PUBKEY        string = "m_pk"
    VBM_CRYPTO_SIGNATURE     string = "m_cs"
)

type VBoxMasterMeta struct {
    ProtocolVersion          string                    `msgpack:"m_pv"`
    MasterState              VBoxMasterState           `msgpack:"m_st"`
    EncryptedPackage         []byte                    `msgpack:"m_ep, inline, omitempty"`
    PublicKey                []byte                    `msgpack:"m_pk, omitempty"`
    CryptoSignature          []byte                    `msgpack:"m_cs, omitempty"`
    MasterAcknowledge        *VBoxMasterAcknowledge    `msgpack:"-"`
    Encryptor                pcrypto.RsaEncryptor      `msgpack:"-"`
    Decryptor                pcrypto.RsaDecryptor      `msgpack:"-"`
}

// --- Acknowledge Field ---
const (
    VBM_AUTH_TOKEN           string = "m_at"
    VBM_TIMESTAMP            string = "m_ts"
)

type VBoxMasterAcknowledge struct {
    AuthToken                string    `msgpack:"m_at, omitempty"`
    TimeStamp                time.Time `msgpack:"m_ts"`
}

// --- Compositions ---
func MasterPackingKeyExchangeAcknowledge(authToken string, pubkey []byte, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    return masterPackingAcknowledge(VBoxMasterKeyExchange, authToken, pubkey, rsaEncryptor)
}

func MasterPackingBoundedAcknowledge(rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    return masterPackingAcknowledge(VBoxMasterBounded, "", nil, rsaEncryptor)
}

func MasterPackingBindBrokenAcknowledge(rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    return masterPackingAcknowledge(VBoxMasterBindBroken, "", nil, rsaEncryptor)
}

func masterPackingAcknowledge(state VBoxMasterState, authToken string, pubkey []byte, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    var (
        meta *VBoxMasterMeta = nil
        apkg, epkg, mpkg []byte = nil, nil, nil
        sig pcrypto.Signature = nil
        err error = nil
    )

    if rsaEncryptor == nil {
        return nil, errors.Errorf("[ERR] master RSA Encryptor cannot be null")
    }

    switch state {
        case VBoxMasterKeyExchange: {
            var (
                ack = &VBoxMasterAcknowledge {
                    AuthToken:    authToken,
                    TimeStamp:    time.Now(),
                }
            )

            // error check
            if len(authToken) == 0 {
                return nil, errors.Errorf("[ERR] invalid auth token assignment")
            }
            if len(pubkey) == 0 {
                return nil, errors.Errorf("[ERR] invalid public key passed")
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
                MasterState:         VBoxMasterKeyExchange,
                EncryptedPackage:    epkg,
                PublicKey:           pubkey,
                CryptoSignature:     sig,
            }
            mpkg, err = msgpack.Marshal(meta)
            return mpkg, errors.WithStack(err)
        }
        default: {
            var (
                ack = &VBoxMasterAcknowledge {
                    TimeStamp:    time.Now(),
                }
            )

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
                MasterState:         state,
                EncryptedPackage:    epkg,
                CryptoSignature:     sig,
            }
            mpkg, err = msgpack.Marshal(meta)
            return mpkg, errors.WithStack(err)
        }
    }
}

func MasterUnpackingAcknowledge(metaPackage, prvkey []byte, rsaDecryptor pcrypto.RsaDecryptor) (*VBoxMasterMeta, error) {
    var (
        meta *VBoxMasterMeta = nil
        ack *VBoxMasterAcknowledge = nil
        apkg []byte = nil
        err error = nil
    )

    // unpack meta & error check
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

    switch meta.MasterState {
        case VBoxMasterKeyExchange: {
            var (
                decryptor pcrypto.RsaDecryptor
                encryptor pcrypto.RsaEncryptor
            )

            // error check
            if len(prvkey) == 0 {
                return nil, errors.Errorf("[ERR] private key cannot be null")
            }
            if len(meta.EncryptedPackage) == 0 {
                return nil, errors.Errorf("[ERR] null encrypted ack")
            }
            if len(meta.PublicKey) == 0 {
                return nil, errors.Errorf("[ERR] null public key")
            }
            if len(meta.CryptoSignature) == 0 {
                return nil, errors.Errorf("[ERR] null crypto signature")
            }

            // build encryptor & decryptor
            encryptor, err = pcrypto.NewRsaEncryptorFromKeyData(meta.PublicKey, prvkey)
            if err != nil {
                return nil, errors.WithStack(err)
            }
            decryptor, err = pcrypto.NewRsaDecryptorFromKeyData(meta.PublicKey, prvkey)
            if err != nil {
                return nil, errors.WithStack(err)
            }

            // decrypt message
            apkg, err = decryptor.DecryptByRSA(meta.EncryptedPackage, meta.CryptoSignature)
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
            if len(ack.AuthToken) == 0 {
                return nil, errors.Errorf("[ERR] invalid auth token assignment")
            }

            // assign fields
            meta.MasterAcknowledge = ack
            meta.Encryptor = encryptor
            meta.Decryptor = decryptor

            return meta, nil
        }
        default: {
            // error check
            if rsaDecryptor == nil {
                return nil, errors.Errorf("[ERR] core RSA Decryptor cannot be null")
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
            if len(ack.AuthToken) != 0 {
                return nil, errors.Errorf("[ERR] invalid ack content w/ auth token")
            }

            // assing acknowledge
            meta.MasterAcknowledge = ack

            return meta, nil
        }
    }
}