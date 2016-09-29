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
    "golang.org/x/crypto/ssh"
)

type Signature []byte

//------------------------------------------------ PRIVATE KEY ---------------------------------------------------------

type rsaPrivateKey struct {
    *rsa.PrivateKey
}

func newPrivateKeyFromKey(k interface{}) (*rsaPrivateKey, error) {
    switch t := k.(type) {
        case *rsa.PrivateKey:
            return &rsaPrivateKey{t}, nil
        default:
            return nil, fmt.Errorf("ssh: unsupported key type %T", k)
    }
}

// parsePublicKey parses a PEM encoded private key.
func parsePrivateKey(pemBytes []byte) (interface{}, error) {
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
    return rawkey, nil
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

// Decrypt returns encrypted payload for the given data.
func (r *rsaPrivateKey) decrypt(data []byte) ([]byte, error) {
    decrypted, err := rsa.DecryptOAEP(sha1.New(), rand.Reader, r.PrivateKey, data, []byte("~pc*crypt^pkg!")); if err != nil {
        return nil, err
    }
    return decrypted, nil
}

func newPrivateKeyFromFile(prvkeyPath string) (*rsaPrivateKey, error) {
    data, err := ioutil.ReadFile(prvkeyPath); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot open private key file %s for error %v", prvkeyPath, err)
    }
    rawkey, err := parsePrivateKey(data); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot parse private rawkey %v", err)
    }
    return newPrivateKeyFromKey(rawkey)
}

// Sign signs data with rsa-sha hash
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
    // TODO : when overal cluster nodes are powerful enough to handle SHA256, please change to that
    return signDataWithHash(r, data, crypto.SHA1)
}

//------------------------------------------------ PUBLIC KEY ----------------------------------------------------------


type rsaPublicKey struct {
    *rsa.PublicKey
}

func newPublicKeyFromKey(k interface{}) (*rsaPublicKey, error) {
    switch t := k.(type) {
        case *rsa.PublicKey:
            return &rsaPublicKey{t}, nil
        default:
            return nil, fmt.Errorf("ssh: unsupported key type %T", k)
    }
}


// parsePublicKey parses a PEM encoded private key.
func parsePublicKey(pemBytes []byte) (interface{}, error) {
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

    return rawkey, nil
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

// Encrypt returns encrypted payload for the given data.
func (r *rsaPublicKey) encrypt(data []byte) ([]byte, error) {
    // The label parameter must be the same for decrypt function
    encrypted, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, r.PublicKey, data, []byte("~pc*crypt^pkg!")); if err != nil {
        return nil, err
    }
    return encrypted, nil
}

func newPublicKeyFromFile(pubkeyPath string) (*rsaPublicKey, error) {
    data, err := ioutil.ReadFile(pubkeyPath); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot open public key file from %s : %v", pubkeyPath, err)
    }
    rawkey, err := parsePublicKey(data); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot parse pulic rawkey %v", err)
    }
    return newPublicKeyFromKey(rawkey)
}

// Unsign verifies the message using a rsa-sha signature
func (r *rsaPublicKey) Unsign(message []byte, sig []byte) error {
    // TODO : when overal cluster nodes are powerful enough to handle SHA256, please change to that
    return unsignDataWithHash(r, message, sig, crypto.SHA1)
}

//------------------------------------------------ KEY GENERATION ------------------------------------------------------

// GenerateKeyPair make a pair of public and private keys encoded in PEM format
func GenerateKeyPair(pubKeyPath, prvkeyPath, sshPubkeyPath string) error {
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

    // generate and write public key for ssh
    pub, err := ssh.NewPublicKey(&privateKey.PublicKey); if err != nil {
        return err
    }
    return ioutil.WriteFile(sshPubkeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
}
