package pcrypto

import (
    "crypto/rsa"
    "fmt"
    "io/ioutil"
    "crypto"
)

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
        sshKey = &rsaPublicKey{
            PublicKey:t,
            Hash:crypto.SHA256,
        }
    default:
        return nil, fmt.Errorf("ssh: unsupported key type %T", k)
    }
    return sshKey, nil
}

// NewUnsignerFromPublicKeyFile loads and parses a PEM encoded public key file
// pubkeyPath : has to be an absolute or a valid path
func NewUnsignerFromPublicKeyFile(pubkeyPath string) (Unsigner, error) {
    data, err := ioutil.ReadFile(pubkeyPath); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot open public key file from %s : %v", pubkeyPath, err)
    }
    rawkey, err := parsePublicKey(data); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot parse pulic rawkey %v", err)
    }
    return newUnsignerFromKey(rawkey)
}
