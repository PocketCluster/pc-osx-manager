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
    "golang.org/x/crypto/ssh"
    "os"
)

type Signature []byte


const (
    //------------------------------------------------ RSA PRIVATE KEY -----------------------------------------------------
    // As of now (10/13/2016), 1024-bit keysize is ineffective to defend from malicious attack.
    // But, this is required due to 1) slow slave node processing power (2048-bit key take 19.sec to pass tests on Odroid C2
    // and 2) bloated CryptoKeyExchange packet sent to slave (up to 620 bytes).
    //
    // We can mitigate this with asymmetric key size (i.e. master 2048, slave 1024) in the future.
    //
    // As of now (11/27/2016), this key will only be used to generate joining pub/prv key pair
    rsaWeakKeySize int = 1024

    // As of now (11/27/2016), 2048-bit keysize is to be used for SSH and Container key
    rsaStrongKeySize int = 2048

    rsaKeyFilePerm os.FileMode = os.FileMode(0600)
)

type rsaPrivateKey struct {
    *rsa.PrivateKey
    crypto.Hash
}

func newPrivateKeyFromKey(k interface{}) (*rsaPrivateKey, error) {
    switch t := k.(type) {
        case *rsa.PrivateKey:
            return &rsaPrivateKey{
                PrivateKey:t,
                Hash:crypto.SHA256,
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
    return r.signDataWithHash(data, crypto.SHA256)
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
                Hash:crypto.SHA256,
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
    return r.unsignDataWithHash(message, sig, crypto.SHA256)
}

//------------------------------------------------ RSA KEY GENERATION --------------------------------------------------
// generateKeyPairs private/ public/ ssh keypair
func generateKeyPairs(keysize int) ([]byte, []byte, []byte, error) {
    var (
        privateKey *rsa.PrivateKey = nil
        privDer, pubDer, sshBytes []byte = nil, nil, nil
        privBlock, pubBlock *pem.Block = nil, nil
        privPem, pubPem []byte = nil, nil
        sshPub ssh.PublicKey
        err error = nil
    )

    // check key size
    if keysize != rsaStrongKeySize && keysize != rsaWeakKeySize {
        return nil, nil, nil, fmt.Errorf("[ERR] RSA key size should be either 1024 or 2048. Current %d", keysize)
    }

    // generate private key
    privateKey, err = rsa.GenerateKey(rand.Reader, keysize)
    if err != nil {
        return nil, nil, nil, err
    }
    // check the key generated
    err = privateKey.Validate()
    if err != nil {
        return nil, nil, nil, err
    }
    // build private key
    privDer = x509.MarshalPKCS1PrivateKey(privateKey)
    privBlock = &pem.Block{
        Type:    "RSA PRIVATE KEY",
        Headers: nil,
        Bytes:   privDer,
    }
    privPem = pem.EncodeToMemory(privBlock)

    // generate and public key
    pubDer, err = x509.MarshalPKIXPublicKey(privateKey.Public())
    if err != nil {
        return nil, nil, nil, err
    }
    pubBlock = &pem.Block{
        Type:   "PUBLIC KEY",
        Headers: nil,
        Bytes:   pubDer,
    }
    pubPem = pem.EncodeToMemory(pubBlock)

    // generate ssh key
    sshPub, err = ssh.NewPublicKey(privateKey.Public())
    if err != nil {
        return nil, nil, nil, err
    }
    sshBytes = ssh.MarshalAuthorizedKey(sshPub)
    return privPem, pubPem, sshBytes, err
}

func generateKeyFies(prvKeyPath, pubKeyPath, sshKeyPath string, keysize int) error {
    prv, pub, ssh, err := generateKeyPairs(keysize)
    if err != nil {
        return err
    }

    if len(prvKeyPath) != 0 && len(prv) != 0 {
        err = ioutil.WriteFile(prvKeyPath, prv, rsaKeyFilePerm)
        if err != nil {
            return err
        }
    }

    if len(pubKeyPath) != 0 && len(pub) != 0 {
        err = ioutil.WriteFile(pubKeyPath, pub, rsaKeyFilePerm)
        if err != nil {
            return err
        }
    }

    if len(sshKeyPath) != 0 && len(ssh) != 0 {
        err = ioutil.WriteFile(sshKeyPath, ssh, rsaKeyFilePerm)
        if err != nil {
            return err
        }
    }
    return nil
}

// GenerateKeyPair make a pair of public and private keys encoded in PEM format
func GenerateWeakKeyPairFiles(pubKeyPath, prvKeyPath, sshKeyPath string) error {
    return generateKeyFies(prvKeyPath, pubKeyPath, sshKeyPath, rsaWeakKeySize)
}

// GenerateKeyPair make a pair of public and private keys encoded in PEM format
func GenerateStrongKeyPairFiles(pubKeyPath, prvKeyPath, sshKeyPath string) error {
    return generateKeyFies(prvKeyPath, pubKeyPath, sshKeyPath, rsaStrongKeySize)
}

// GenerateKeyPair make a pair of public and private keys encoded in PEM format
func GenerateWeakKeyPair() (pub []byte, prv []byte, ssh []byte, err error) {
    prv, pub, ssh, err = generateKeyPairs(rsaWeakKeySize)
    return
}

// GenerateKeyPair make a pair of public and private keys encoded in PEM format
func GenerateStrongKeyPair() (pub []byte, prv []byte, ssh []byte, err error) {
    prv, pub, ssh, err = generateKeyPairs(rsaStrongKeySize)
    return
}