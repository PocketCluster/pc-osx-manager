package crypt

import (
    "crypto/rsa"
    "fmt"
    "io/ioutil"
)

// A Signer is can create signatures that verify against a public key.
type Signer interface {
    // Sign returns raw signature for the given data. This method
    // will apply the hash specified for the keytype to the data.
    Sign(data []byte) ([]byte, error)
}

// NewSignerFromPrivateKeyFile loads and parses a PEM encoded private key file
// prvkeyPath : has to be an absolute or a valid path
func NewSignerFromPrivateKeyFile(prvkeyPath string) (Signer, error) {
    data, err := ioutil.ReadFile(prvkeyPath); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot open private key file %s for error %v", prvkeyPath, err)
    }
    rawkey, err := parsePrivateKey(data); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot parse private rawkey %v", err)
    }
    return newSignerFromKey(rawkey)
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
