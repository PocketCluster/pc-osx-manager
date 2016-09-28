package crypt

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "errors"
    "fmt"
    "io/ioutil"
    "crypto/sha1"
)

/*
def verify_signature(pubkey_path, message, signature):
    '''
    Use Crypto.Signature.PKCS1_v1_5 to verify the signature on a message.
    Returns True for valid signature.
    '''
    log.debug('salt.crypt.verify_signature: Loading public key')
    with salt.utils.fopen(pubkey_path) as f:
        pubkey = RSA.importKey(f.read())
    log.debug('salt.crypt.verify_signature: Verifying signature')
    verifier = PKCS1_v1_5.new(pubkey)
    return verifier.verify(SHA.new(message), signature)

 */

func Cryptotest() {
    signer, err := loadPrivateKeyFile("crypt/test.prv")
    if err != nil {
        fmt.Errorf("signer is damaged: %v", err)
    }

    toSign := "date: Thu, 05 Jan 2012 21:31:40 GMT"

    signed, err := signer.Sign([]byte(toSign))
    if err != nil {
        fmt.Errorf("could not sign request: %v", err)
    }
    sig := base64.StdEncoding.EncodeToString(signed)
    fmt.Printf("Signature: %v\n", sig)



    parser, perr := loadPublicKeyFile("crypt/test.pub")
    if perr != nil {
        fmt.Errorf("could not sign request: %v", err)
    }

    err = parser.Unsign([]byte(toSign), signed)
    if err != nil {
        fmt.Errorf("could not unsign request: %v", err)
    }

    fmt.Printf("Unsign error: %v\n", err)
}


func VerifySignature(pubkeyPath string, message string, signature string) error {

    parser, perr := loadPublicKeyFile(pubkeyPath)
    if perr != nil {
        return fmt.Errorf("[ERR] could load public key: %v", perr)
    }

    err := parser.Unsign([]byte(message), []byte(signature))
    if err != nil {
        return fmt.Errorf("[ERR] could not verify message : %v", err)
    }

    return nil
}

//------------------------------------------------ PRIVATE KEY ---------------------------------------------------------

// A Signer is can create signatures that verify against a public key.
type Signer interface {
    // Sign returns raw signature for the given data. This method
    // will apply the hash specified for the keytype to the data.
    Sign(data []byte) ([]byte, error)
}

type rsaPrivateKey struct {
    *rsa.PrivateKey
}

func newSignerFromKey(k interface{}) (Signer, error) {
    var sshKey Signer
    switch t := k.(type) {
    case *rsa.PrivateKey:
        sshKey = &rsaPrivateKey{t}
    default:
        return nil, fmt.Errorf("ssh: unsupported key type %T", k)
    }
    return sshKey, nil
}

func signDataWithHash(r *rsaPrivateKey, data []byte, hashType crypto.Hash) ([]byte, error) {
    switch hashType {
        case crypto.SHA1:{
            h := sha1.New()
            h.Write(data)
            d := h.Sum(nil)
            return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA1, d)
        }
        case crypto.SHA256:{
            h := sha256.New()
            h.Write(data)
            d := h.Sum(nil)
            return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, crypto.SHA256, d)
        }
        default:
            return nil, errors.New("[ERR] SHA should be SHA1 or SHA256")
    }
}

// Sign signs data with rsa-sha hash
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
    return signDataWithHash(r, data, crypto.SHA1)
}

// parsePublicKey parses a PEM encoded private key.
func parsePrivateKey(pemBytes []byte) (Signer, error) {
    block, _ := pem.Decode(pemBytes)
    if block == nil {
        return nil, errors.New("ssh: no key found")
    }

    var rawkey interface{}
    switch block.Type {
    case "RSA PRIVATE KEY":
        rsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
        if err != nil {
            return nil, err
        }
        rawkey = rsa
    default:
        return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
    }
    return newSignerFromKey(rawkey)
}

func loadTestPrivateKey(path string) (Signer, error) {
    return parsePrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDCFENGw33yGihy92pDjZQhl0C36rPJj+CvfSC8+q28hxA161QF
NUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6Z4UMR7EOcpfdUE9Hf3m/hs+F
UR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJwoYi+1hqp1fIekaxsyQIDAQAB
AoGBAJR8ZkCUvx5kzv+utdl7T5MnordT1TvoXXJGXK7ZZ+UuvMNUCdN2QPc4sBiA
QWvLw1cSKt5DsKZ8UETpYPy8pPYnnDEz2dDYiaew9+xEpubyeW2oH4Zx71wqBtOK
kqwrXa/pzdpiucRRjk6vE6YY7EBBs/g7uanVpGibOVAEsqH1AkEA7DkjVH28WDUg
f1nqvfn2Kj6CT7nIcE3jGJsZZ7zlZmBmHFDONMLUrXR/Zm3pR5m0tCmBqa5RK95u
412jt1dPIwJBANJT3v8pnkth48bQo/fKel6uEYyboRtA5/uHuHkZ6FQF7OUkGogc
mSJluOdc5t6hI1VsLn0QZEjQZMEOWr+wKSMCQQCC4kXJEsHAve77oP6HtG/IiEn7
kpyUXRNvFsDE0czpJJBvL/aRFUJxuRK91jhjC68sA7NsKMGg5OXb5I5Jj36xAkEA
gIT7aFOYBFwGgQAQkWNKLvySgKbAZRTeLBacpHMuQdl1DfdntvAyqpAZ0lY0RKmW
G6aFKaqQfOXKCyWoUiVknQJAXrlgySFci/2ueKlIE1QqIiLSZ8V8OlpFLRnb1pzI
7U1yQXnTAEFYM560yJlzUpOb1V4cScGd365tiSMvxLOvTA==
-----END RSA PRIVATE KEY-----`))
}

// loadPrivateKeyFile loads and parses a PEM encoded private key file
// prvkeyPath : has to be an absolute or a valid path
func loadPrivateKeyFile(prvkeyPath string) (Signer, error) {
    data, err := ioutil.ReadFile(prvkeyPath); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot open private key file %s for error %v", prvkeyPath, err)
    }
    return parsePrivateKey(data)
}

//------------------------------------------------ PUBLIC KEY ----------------------------------------------------------

type rsaPublicKey struct {
    *rsa.PublicKey
}

// A Signer is can create signatures that verify against a public key.
type Unsigner interface {
    // Sign returns raw signature for the given data. This method
    // will apply the hash specified for the keytype to the data.
    Unsign(data[]byte, sig []byte) error
}

func newUnsignerFromKey(k interface{}) (Unsigner, error) {
    var sshKey Unsigner
    switch t := k.(type) {
    case *rsa.PublicKey:
        sshKey = &rsaPublicKey{t}
    default:
        return nil, fmt.Errorf("ssh: unsupported key type %T", k)
    }
    return sshKey, nil
}

func unsignDataWithHash(r *rsaPublicKey, message []byte, sig []byte, hashType crypto.Hash) error {
    switch hashType {
        case crypto.SHA1:{
            h := sha1.New()
            h.Write(message)
            d := h.Sum(nil)
            return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA1, d, sig)
        }
        case crypto.SHA256:{
            h := sha256.New()
            h.Write(message)
            d := h.Sum(nil)
            return rsa.VerifyPKCS1v15(r.PublicKey, crypto.SHA256, d, sig)
        }
        default:
            return errors.New("[ERR] SHA should be SHA1 or SHA256")
    }
}

// Unsign verifies the message using a rsa-sha signature
func (r *rsaPublicKey) Unsign(message []byte, sig []byte) error {
    return unsignDataWithHash(r, message, sig, crypto.SHA1)
}

// parsePublicKey parses a PEM encoded private key.
func parsePublicKey(pemBytes []byte) (Unsigner, error) {
    block, _ := pem.Decode(pemBytes)
    if block == nil {
        return nil, errors.New("ssh: no key found")
    }

    var rawkey interface{}
    switch block.Type {
    case "PUBLIC KEY":
        rsa, err := x509.ParsePKIXPublicKey(block.Bytes)
        if err != nil {
            return nil, err
        }
        rawkey = rsa
    default:
        return nil, fmt.Errorf("ssh: unsupported key type %q", block.Type)
    }

    return newUnsignerFromKey(rawkey)
}

// loadPublicKeyFile loads and parses a PEM encoded public key file
// pubkeyPath : has to be an absolute or a valid path
func loadPublicKeyFile(pubkeyPath string) (Unsigner, error) {
    data, err := ioutil.ReadFile(pubkeyPath); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot open public key file %s for error %v", pubkeyPath, err)
    }
    return parsePublicKey(data)
}

// loadTestPublicKey loads an parses a PEM encoded public key file.
func loadTestPublicKey(path string) (Unsigner, error) {
    return parsePublicKey([]byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`))
}
