package pcrypto

import (
    "crypto"
    "crypto/rsa"
    "fmt"
    "errors"
)

type RsaDecryptor interface {
    DecryptByRSA(crypted []byte, sendSig Signature) (plain []byte, err error)
}

func NewDecryptorFromKeyFiles(sendPubkeyPath, recvPrvkeyPath string) (RsaDecryptor, error) {
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

func NewDecryptorFromKeyData(sendPubkeyData, recvPrvkeyData []byte) (RsaDecryptor, error) {
    pubkey, err := newPublicKeyFromData(sendPubkeyData); if err != nil {
        return nil, err
    }
    prvkey, err := newPrivateKeyFromData(recvPrvkeyData); if err != nil {
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
    hType := crypto.SHA256
    hash := hType.New()
    hash.Write(plainText)
    opts := &rsa.PSSOptions {SaltLength:rsa.PSSSaltLengthAuto}

    if err := rsa.VerifyPSS(d.sendPubkey.PublicKey, hType, hash.Sum(nil), sendSig, opts); err != nil {
        return err
    }
    return nil
}

func (d *decryption) DecryptByRSA(crypted []byte, sendSig Signature) (plain []byte, err error) {
    plain, err = d.recvPrvkey.decrypt(crypted); if err != nil {
        return nil, errors.New("[ERR] decryption failed due to " + err.Error())
    }
    if err = d.verifySignature(plain, sendSig); err != nil {
        return nil, fmt.Errorf("[ERR] Cannot verify message with sender signature %v", err)
    }
    return plain, nil
}
