package crypt

import (
    "fmt"
    "io/ioutil"
)

func SignMessageWithPrivateKeyFile(prvkeyPath string, message string) ([]byte, error) {
    // private key signer
    signer, err := NewSignerFromPrivateKeyFile(prvkeyPath); if err != nil {
        return nil, err
    }
    return signer.Sign([]byte(message))
}

// Use PKCS1 v1.5 to verify the signature on a message. Returns True for valid signature.
func VerifySignature(pubkeyPath string, message string, signature []byte) error {
    unsigner, perr := NewUnsignerFromPublicKeyFile(pubkeyPath); if perr != nil {
        return fmt.Errorf("[ERR] could load public key: %v", perr)
    }
    if err := unsigner.Unsign([]byte(message), signature); err != nil {
        return fmt.Errorf("[ERR] could not verify message : %v", err)
    }
    return nil
}

func GenerateSignature(pubkeyPath string, prvkeyPath string) ([]byte, error) {
    // pubkey data
    toSign, err := ioutil.ReadFile(pubkeyPath); if err != nil {
        return nil, fmt.Errorf("[ERR] cannot open public key file %s for error %v", pubkeyPath, err)
    }
    // private key signer
    signer, err := NewSignerFromPrivateKeyFile(prvkeyPath); if err != nil {
        return nil, err
    }
    // sign the pubkey
    signed, err := signer.Sign([]byte(toSign)); if err != nil {
        return nil, err
    }
    return signed, nil
}