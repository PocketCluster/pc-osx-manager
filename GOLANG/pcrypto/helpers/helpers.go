package helpers

import (
    "bytes"
    "crypto"
    "crypto/ecdsa"
    "crypto/rsa"
    "crypto/x509"
    "crypto/elliptic"
    "encoding/pem"
    "fmt"
    "math/big"
    "encoding/asn1"
    "errors"
    "strings"

)

// TODO : add tests
// ParseCertificatePEM parses and returns a PEM-encoded certificate
func ParseCertificatePEM(certPEM []byte) (*x509.Certificate, error) {
    certPEM = bytes.TrimSpace(certPEM)
    cert, rest, err := ParseOneCertificateFromPEM(certPEM)
    if err != nil {
        // Log the actual parsing error but throw a default parse error message.
        return nil, fmt.Errorf("[ERR] Certificate parsing error: %v", err)
    } else if cert == nil {
        return nil, fmt.Errorf("[ERR] Certificate parsing error. Decode failure 1")
    } else if len(rest) > 0 {
        return nil, fmt.Errorf("[ERR] Certificate parsing error. Decode failure 2")
    } else if len(cert) > 1 {
        return nil, fmt.Errorf("[ERR] Certificate parsing error. the PKCS7 object in the PEM file should contain only one certificate")
    }
    return cert[0], nil
}

// TODO : add tests
// ParseOneCertificateFromPEM attempts to parse one PEM encoded certificate object,
// a raw x509 certificate from the top of certsPEM, which itself may
// contain multiple PEM encoded certificate objects.
func ParseOneCertificateFromPEM(certsPEM []byte) ([]*x509.Certificate, []byte, error) {
    block, rest := pem.Decode(certsPEM)
    if block == nil {
        return nil, rest, nil
    }
    cert, err := x509.ParseCertificate(block.Bytes)
    if err != nil {
        return nil, rest, nil
    }
    return []*x509.Certificate{cert}, rest, nil
}

// TODO : add tests
// CheckSignature verifies a signature made by the key on a CSR, such
// as on the CSR itself.
func CheckSignature(csr *x509.CertificateRequest, algo x509.SignatureAlgorithm, signed, signature []byte) error {
    var hashType crypto.Hash

    switch algo {
    case x509.SHA1WithRSA, x509.ECDSAWithSHA1:
        hashType = crypto.SHA1
    case x509.SHA256WithRSA, x509.ECDSAWithSHA256:
        hashType = crypto.SHA256
    case x509.SHA384WithRSA, x509.ECDSAWithSHA384:
        hashType = crypto.SHA384
    case x509.SHA512WithRSA, x509.ECDSAWithSHA512:
        hashType = crypto.SHA512
    default:
        return x509.ErrUnsupportedAlgorithm
    }

    if !hashType.Available() {
        return x509.ErrUnsupportedAlgorithm
    }
    h := hashType.New()

    h.Write(signed)
    digest := h.Sum(nil)

    switch pub := csr.PublicKey.(type) {
    case *rsa.PublicKey:
        return rsa.VerifyPKCS1v15(pub, hashType, digest, signature)
    case *ecdsa.PublicKey:
        ecdsaSig := new(struct{ R, S *big.Int })
        if _, err := asn1.Unmarshal(signature, ecdsaSig); err != nil {
            return err
        }
        if ecdsaSig.R.Sign() <= 0 || ecdsaSig.S.Sign() <= 0 {
            return errors.New("x509: ECDSA signature contained zero or negative values")
        }
        if !ecdsa.Verify(pub, digest, ecdsaSig.R, ecdsaSig.S) {
            return errors.New("x509: ECDSA verification failure")
        }
        return nil
    }
    return x509.ErrUnsupportedAlgorithm
}

// ParsePrivateKeyPEM parses and returns a PEM-encoded private
// key. The private key may be either an unencrypted PKCS#8, PKCS#1,
// or elliptic private key.
func ParsePrivateKeyPEM(keyPEM []byte) (key crypto.Signer, err error) {
    return ParsePrivateKeyPEMWithPassword(keyPEM, nil)
}


// ParsePrivateKeyDER parses a PKCS #1, PKCS #8, or elliptic curve
// DER-encoded private key. The key must not be in PEM format.
func ParsePrivateKeyDER(keyDER []byte) (key crypto.Signer, err error) {
    generalKey, err := x509.ParsePKCS8PrivateKey(keyDER)
    if err != nil {
        generalKey, err = x509.ParsePKCS1PrivateKey(keyDER)
        if err != nil {
            generalKey, err = x509.ParseECPrivateKey(keyDER)
            if err != nil {
                // We don't include the actual error into
                // the final error. The reason might be
                // we don't want to leak any info about
                // the private key.
                return nil, errors.New("[ERR] PrivateKeyError : Parse failed")
            }
        }
    }

    switch generalKey.(type) {
    case *rsa.PrivateKey:
        return generalKey.(*rsa.PrivateKey), nil
    case *ecdsa.PrivateKey:
        return generalKey.(*ecdsa.PrivateKey), nil
    }

    // should never reach here
    return nil, errors.New("[ERR] PrivateKeyError : Parse failed")
}

// ParsePrivateKeyPEMWithPassword parses and returns a PEM-encoded private
// key. The private key may be a potentially encrypted PKCS#8, PKCS#1,
// or elliptic private key.
func ParsePrivateKeyPEMWithPassword(keyPEM []byte, password []byte) (key crypto.Signer, err error) {
    keyDER, err := GetKeyDERFromPEM(keyPEM, password)
    if err != nil {
        return nil, err
    }

    return ParsePrivateKeyDER(keyDER)
}

// GetKeyDERFromPEM parses a PEM-encoded private key and returns DER-format key bytes.
func GetKeyDERFromPEM(in []byte, password []byte) ([]byte, error) {
    keyDER, _ := pem.Decode(in)
    if keyDER != nil {
        if procType, ok := keyDER.Headers["Proc-Type"]; ok {
            if strings.Contains(procType, "ENCRYPTED") {
                if password != nil {
                    return x509.DecryptPEMBlock(keyDER, password)
                }
                return nil, errors.New("[ERR] PrivateKeyError : Decryption failed")
            }
        }
        return keyDER.Bytes, nil
    }

    return nil, errors.New("[ERR] PrivateKeyError : Decode failed")
}

// DefaultSigAlgo returns an appropriate X.509 signature algorithm given
// the CA's private key.
func DefaultSigAlgo(priv crypto.Signer) x509.SignatureAlgorithm {
    pub := priv.Public()
    switch pub := pub.(type) {
    case *rsa.PublicKey:
        keySize := pub.N.BitLen()
        switch {
        case keySize >= 4096:
            return x509.SHA512WithRSA
        case keySize >= 3072:
            return x509.SHA384WithRSA
        case keySize >= 2048:
            return x509.SHA256WithRSA
        default:
            return x509.SHA1WithRSA
        }
    case *ecdsa.PublicKey:
        switch pub.Curve {
        case elliptic.P256():
            return x509.ECDSAWithSHA256
        case elliptic.P384():
            return x509.ECDSAWithSHA384
        case elliptic.P521():
            return x509.ECDSAWithSHA512
        default:
            return x509.ECDSAWithSHA1
        }
    default:
        return x509.UnknownSignatureAlgorithm
    }
}