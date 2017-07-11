package pkg

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

type VBoxCoreState int
const (
    VBoxCoreUnbounded       VBoxCoreState = iota
    VBoxCoreBounded
    VBoxCoreBindBroken
)

func (s VBoxCoreState) String() string {
    var state string
    switch s {
    case VBoxCoreUnbounded:
        state = "VBoxCoreUnbounded"
    case VBoxCoreBounded:
        state = "VBoxCoreBounded"
    case VBoxCoreBindBroken:
        state = "VBoxCoreBindBroken"
    }
    return state
}

// --- Meta Field ---
const (
    VBC_PROTOCOL_VERSION    string = "c_pv"
    VBC_CORE_STATE          string = "c_st"
    VBC_ENCRYPTED_PKG       string = "c_ep"
    VBC_PUBLIC_KEY          string = "c_pk"
    VBC_CRYPTO_SIGNATURE    string = "c_cs"
)

type VBoxCoreMeta struct {
    ProtocolVersion         string             `msgpack:"c_pv"`
    CoreState               VBoxCoreState      `msgpack:"c_st"`
    EncryptedPackage        []byte             `msgpack:"c_ep, inline, omitempty"`
    PublicKey               []byte             `msgpack:"c_pk, omitempty"`
    CryptoSignature         []byte             `msgpack:"c_cs, omitempty"`
    CoreStatus              *VBoxCoreStatus    `msgpack:"-"`
}

// --- Status Field ---
const (
    VBC_EXT_IP4_ADDR_SMASK  string = "c_i4"
    VBC_EXT_IP4_GATEWAY     string = "c_g4"
    VBC_TIMESTAMP           string = "c_ts"
)

type VBoxCoreStatus struct {
    ExtIP4AddrSmask         string             `msgpack:"c_i4"`
    ExtIP4Gateway           string             `msgpack:"c_g4"`
    TimeStamp               time.Time          `msgpack:"c_ts"`
}

// --- Compositions ---
func CorePackingUnboundedStatus(pubkey []byte) ([]byte, error) {
    return corePackingStatus(VBoxCoreUnbounded, pubkey, "", "", nil)
}

func CorePackingBoundedStatus(extAddr, extGateway string, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    return corePackingStatus(VBoxCoreBounded, nil, extAddr, extGateway, rsaEncryptor)
}

func CorePackingBindBrokenStatus(extAddr, extGateway string, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    return corePackingStatus(VBoxCoreBindBroken, nil, extAddr, extGateway, rsaEncryptor)
}

func corePackingStatus(state VBoxCoreState, pubkey []byte, extAddr, extGateway string, rsaEncryptor pcrypto.RsaEncryptor) ([]byte, error) {
    switch state {
        case VBoxCoreUnbounded: {
            var (
                meta = &VBoxCoreMeta{
                    ProtocolVersion:    VBoxCoreVersion,
                    CoreState:          VBoxCoreUnbounded,
                    PublicKey:          pubkey,
                }
                pkg []byte = nil
                err error = nil
            )
            // error check
            if len(pubkey) == 0 {
                return nil, errors.Errorf("[ERR] core status public key cannot be null")
            }

            // packing
            pkg, err = msgpack.Marshal(meta)
            return pkg, errors.WithStack(err)
        }

        default: {
            var (
                meta *VBoxCoreMeta = nil
                status = &VBoxCoreStatus {
                    ExtIP4AddrSmask:    extAddr,
                    ExtIP4Gateway:      extGateway,
                    TimeStamp:          time.Now(),
                }
                spkg, epkg, mpkg []byte = nil, nil, nil
                sig pcrypto.Signature
                err error = nil
            )

            // error check
            if len(extAddr) == 0 {
                return nil, errors.Errorf("[ERR] core status external address cannot be empty")
            }
            if len(extGateway) == 0 {
                return nil, errors.Errorf("[ERR] core status external gateway cannot be empty")
            }
            if rsaEncryptor == nil {
                return nil, errors.Errorf("[ERR] core status RSA Encryptor cannot be null")
            }

            // packaging status
            spkg, err = msgpack.Marshal(status)
            if err != nil {
                return nil, errors.WithStack(err)
            }

            // encrypt status package
            epkg, sig, err = rsaEncryptor.EncryptByRSA(spkg)
            if err != nil {
                return nil, errors.WithStack(err)
            }

            // meta message packing
            meta = &VBoxCoreMeta{
                ProtocolVersion:     VBoxCoreVersion,
                CoreState:           state,
                EncryptedPackage:    epkg,
                CryptoSignature:     sig,
            }
            mpkg, err = msgpack.Marshal(meta)
            return mpkg, errors.WithStack(err)
        }
    }
}

func CoreUnpackingStatus(metaPackage []byte, rsaDecryptor pcrypto.RsaDecryptor) (*VBoxCoreMeta, error) {
    var (
        meta *VBoxCoreMeta
        err error = nil
    )
    // error check
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

    switch meta.CoreState {
        case VBoxCoreUnbounded: {
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
        default: {
            var (
                status *VBoxCoreStatus = nil
                spkg []byte = nil
            )

            // error check
            if rsaDecryptor == nil {
                return nil, errors.Errorf("[ERR] core status RSA Decryptor cannot be null")
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
            if len(status.ExtIP4AddrSmask) == 0 {
                return nil, errors.Errorf("[ERR] invalid ip & subnet mask")
            }
            if len(status.ExtIP4Gateway) == 0 {
                return nil, errors.Errorf("[ERR] invalid gateway")
            }

            // assign status
            meta.CoreStatus = status

            return meta, nil
        }
    }
}