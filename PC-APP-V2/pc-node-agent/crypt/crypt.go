package crypt

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "fmt"
    "io/ioutil"
    "crypto/sha1"
    "os"
)

// Use PKCS1 v1.5 to verify the signature on a message. Returns True for valid signature.
func VerifySignature(pubkeyPath string, message string, signature string) error {

    parser, perr := NewUnsignerFromPublicKeyFile(pubkeyPath)
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
    // TODO : when overal cluster nodes are powerful enough to handle SHA256, please change to that
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

// NewSignerFromPrivateKeyFile loads and parses a PEM encoded private key file
// prvkeyPath : has to be an absolute or a valid path
func NewSignerFromPrivateKeyFile(prvkeyPath string) (Signer, error) {
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
    // TODO : when overal cluster nodes are powerful enough to handle SHA256, please change to that
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

// NewUnsignerFromPublicKeyFile loads and parses a PEM encoded public key file
// pubkeyPath : has to be an absolute or a valid path
func NewUnsignerFromPublicKeyFile(pubkeyPath string) (Unsigner, error) {
    data, err := ioutil.ReadFile(pubkeyPath); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot open public key file %s for error %v", pubkeyPath, err)
    }
    return parsePublicKey(data)
}


//------------------------------------------------ KEY GENERATION ------------------------------------------------------

// GenerateKeyPair make a pair of public and private keys encoded in PEM format
func GenerateKeyPair(pubKeyPath string, prvkeyPath string) error {
    privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return err
    }

    // generate and write private key as PEM
    prvkeyFile, err := os.Create(prvkeyPath)
    defer prvkeyFile.Close()
    if err != nil {
        return err
    }
    prvkeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
    if err := pem.Encode(prvkeyFile, prvkeyPEM); err != nil {
        return err
    }

    // generate and write public key as PEM
    pubkeyFile, err := os.Create(pubKeyPath)
    defer pubkeyFile.Close()
    if err != nil {
        return err
    }
    pubkeyM, err := x509.MarshalPKIXPublicKey(privateKey.Public()); if err != nil {
        return err
    }
    pubkeyPEM := &pem.Block{Type: "PUBLIC KEY", Bytes:pubkeyM}
    if err = pem.Encode(pubkeyFile, pubkeyPEM); err != nil {
        return err
    }

    return nil
    // we can also write in SSH pubkey file format, but that's not the focus of this function
    /*
        // generate and write public key
        pub, err := ssh.NewPublicKey(&privateKey.PublicKey); if err != nil {
            return err
        }
        return ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
    */
}