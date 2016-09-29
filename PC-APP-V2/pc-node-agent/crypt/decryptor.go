package crypt

import (
    "crypto"
    "crypto/rsa"
    "fmt"
    "errors"
)

type Decryptor interface {
    DecryptMessage(crypted string, sendSig Signature) ([]byte, error)
}

func NewDecryptorFromKeyFiles(sendPubkeyPath, recvPrvkeyPath string) (Decryptor, error) {
    pubkey, err := newPublicKeyFromFile(sendPubkeyPath); if err != nil {
        return nil, err
    }
    prvkey, err := newPrivateKeyFromFile(recvPrvkeyPath); if err != nil {
        return nil, err
    }
    return &decryption{
        sendPubkey: pubkey,
        recvPrvkey: prvkey,
    }, nil
}

type decryption struct {
    sendPubkey      *rsaPublicKey
    recvPrvkey      *rsaPrivateKey
}

func (d *decryption) verifySignature(plainText []byte, sendSig Signature) error {
    hType := crypto.SHA1
    hash := hType.New()
    hash.Write(plainText)
    opts := &rsa.PSSOptions {SaltLength:rsa.PSSSaltLengthAuto}

    if err := rsa.VerifyPSS(d.sendPubkey.PublicKey, hType, hash.Sum(nil), sendSig, opts); err != nil {
        return err
    }
    return nil
}

func (d *decryption) DecryptMessage(crypted string, sendSig Signature) ([]byte, error) {
    plain, err := d.recvPrvkey.decrypt([]byte(crypted)); if err != nil {
        return nil, errors.New("[ERR] decryption failed due to " + err.Error())
    }
    if err := d.verifySignature(plain, sendSig); err != nil {
        return nil, fmt.Errorf("[ERR] Cannot verify message with sender signature %v", err)
    }
    return plain, nil
}
