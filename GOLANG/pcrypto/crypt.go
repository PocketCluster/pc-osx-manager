package pcrypto

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "errors"
    "fmt"
    "io/ioutil"
    "os"
    "golang.org/x/crypto/ssh"
)

type Signature []byte

//------------------------------------------------ RSA PRIVATE KEY -----------------------------------------------------
// TODO : Implement stronger encryption key
// As of now (10/13/2016), 1024-bit keysize is ineffective to defend from malicious attack.
// But, this is required due to 1) slow slave node processing power (2048-bit key take 19.sec to pass tests on Odroid C2
// and 2) bloated CryptoKeyExchange packet sent to slave (up to 620 bytes).
//
// We can mitigate this with asymmetric key size (i.e. master 2048, slave 1024) in the future.

const rsaKeySize int = 1024

type rsaPrivateKey struct {
    *rsa.PrivateKey
    crypto.Hash
}

func newPrivateKeyFromKey(k interface{}) (*rsaPrivateKey, error) {
    switch t := k.(type) {
        case *rsa.PrivateKey:
            return &rsaPrivateKey{
                PrivateKey:t,
                Hash:crypto.SHA1,
            }, nil
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

func newPrivateKeyFromFile(prvkeyPath string) (*rsaPrivateKey, error) {
    data, err := ioutil.ReadFile(prvkeyPath); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot open private key file %s for error %v", prvkeyPath, err)
    }
    rawkey, err := parsePrivateKey(data); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot parse private rawkey %v", err)
    }
    return newPrivateKeyFromKey(rawkey)
}

func newPrivateKeyFromData(prvkeyData []byte) (*rsaPrivateKey, error) {
    if len(prvkeyData) == 0 {
        return nil, fmt.Errorf("[ERR] cannot create private with null data")
    }
    rawkey, err := parsePrivateKey(prvkeyData); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot parse private rawkey %v", err)
    }
    return newPrivateKeyFromKey(rawkey)
}

// Decrypt returns encrypted payload for the given data.
func (r *rsaPrivateKey) decrypt(data []byte) ([]byte, error) {
    decrypted, err := rsa.DecryptOAEP(r.Hash.New(), rand.Reader, r.PrivateKey, data, []byte("~pc*crypt^pkg!")); if err != nil {
        return nil, err
    }
    return decrypted, nil
}

func (r *rsaPrivateKey) signDataWithHash(data []byte, hashType crypto.Hash) ([]byte, error) {
    h := r.Hash.New()
    h.Write(data)
    d := h.Sum(nil)
    return rsa.SignPKCS1v15(rand.Reader, r.PrivateKey, r.Hash, d)
}

// Sign signs data with rsa-sha hash
func (r *rsaPrivateKey) Sign(data []byte) ([]byte, error) {
    // TODO : when overal cluster nodes are powerful enough to handle SHA256, please change to that
    return r.signDataWithHash(data, crypto.SHA1)
}

//------------------------------------------------ RSA PUBLIC KEY ------------------------------------------------------

type rsaPublicKey struct {
    *rsa.PublicKey
    crypto.Hash
}

func newPublicKeyFromKey(k interface{}) (*rsaPublicKey, error) {
    switch t := k.(type) {
        case *rsa.PublicKey:
            return &rsaPublicKey{
                PublicKey:t,
                Hash:crypto.SHA1,
            }, nil
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


func newPublicKeyFromFile(pubkeyPath string) (*rsaPublicKey, error) {
    data, err := ioutil.ReadFile(pubkeyPath); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot open public key file from %s : %v", pubkeyPath, err)
    }
    rawkey, err := parsePublicKey(data); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot parse pulic rawkey %v", err)
    }
    return newPublicKeyFromKey(rawkey)
}

func newPublicKeyFromData(pubkeyData []byte) (*rsaPublicKey, error) {
    if len(pubkeyData) == 0 {
        return nil, fmt.Errorf("[ERR] cannot create public key with null data")
    }
    rawkey, err := parsePublicKey(pubkeyData); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot parse pulic rawkey %v", err)
    }
    return newPublicKeyFromKey(rawkey)
}

// Encrypt returns encrypted payload for the given data.
func (r *rsaPublicKey) encrypt(data []byte) ([]byte, error) {
    // The label parameter must be the same for decrypt function
    encrypted, err := rsa.EncryptOAEP(r.Hash.New(), rand.Reader, r.PublicKey, data, []byte("~pc*crypt^pkg!")); if err != nil {
        return nil, err
    }
    return encrypted, nil
}

func (r *rsaPublicKey) unsignDataWithHash(message []byte, sig []byte, hashType crypto.Hash) error {
    h := r.Hash.New()
    h.Write(message)
    d := h.Sum(nil)
    return rsa.VerifyPKCS1v15(r.PublicKey, r.Hash, d, sig)
}

// Unsign verifies the message using a rsa-sha signature
func (r *rsaPublicKey) Unsign(message []byte, sig []byte) error {
    // TODO : when overal cluster nodes are powerful enough to handle SHA256, please change to that
    return r.unsignDataWithHash(message, sig, crypto.SHA1)
}

//------------------------------------------------ RSA KEY GENERATION --------------------------------------------------

// GenerateKeyPair make a pair of public and private keys encoded in PEM format
func GenerateKeyPair(pubKeyPath, prvkeyPath, sshPubkeyPath string) error {
    privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize); if err != nil {
        return err
    }
    if err = privateKey.Validate(); err != nil {
        return err
    }

    // generate and write private key as PEM
    prvkeyFile, err := os.Create(prvkeyPath)
    // as go defer works in LIFO, we'll chage permissions first
    defer os.Chmod(prvkeyPath, 0600)
    defer prvkeyFile.Close()
    if err != nil {
        return err
    }
    prvkeyPEM := &pem.Block{
        Type: "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    }
    if err := pem.Encode(prvkeyFile, prvkeyPEM); err != nil {
        return err
    }

    // generate and write public key as PEM
    pubkeyFile, err := os.Create(pubKeyPath)
    defer os.Chmod(pubKeyPath, 0600)
    defer pubkeyFile.Close()
    if err != nil {
        return err
    }
    pubkeyM, err := x509.MarshalPKIXPublicKey(privateKey.Public()); if err != nil {
        return err
    }
    pubkeyPEM := &pem.Block{
        Type: "PUBLIC KEY",
        Bytes:pubkeyM,
    }
    if err = pem.Encode(pubkeyFile, pubkeyPEM); err != nil {
        return err
    }

    // generate and write public key for ssh
    pub, err := ssh.NewPublicKey(&privateKey.PublicKey); if err != nil {
        return err
    }
    return ioutil.WriteFile(sshPubkeyPath, ssh.MarshalAuthorizedKey(pub), 0600)
}