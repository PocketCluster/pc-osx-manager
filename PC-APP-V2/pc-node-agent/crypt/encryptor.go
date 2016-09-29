package crypt

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
)

type Encryptor interface {
    EncryptMessage(plain string) ([]byte, Signature, error)
}

func NewEncryptorFromKeyFiles(recvPubkeyPath, sendPrvkeyPath string) (Encryptor, error) {
    pubkey, err := newPublicKeyFromFile(recvPubkeyPath); if err != nil {
        return nil, err
    }
    prvkey, err := newPrivateKeyFromFile(sendPrvkeyPath); if err != nil {
        return nil, err
    }
    return &encryptor{
        recvPubkey:pubkey,
        sendPrvkey:prvkey,
    }, nil
}

type encryptor struct {
    recvPubkey      *rsaPublicKey
    sendPrvkey      *rsaPrivateKey
}

func (e *encryptor) generateSignature(plain []byte) (Signature, error) {
    hType := crypto.SHA1
    hash := hType.New()
    hash.Write(plain)
    opts := &rsa.PSSOptions{SaltLength:rsa.PSSSaltLengthAuto}

    return rsa.SignPSS(rand.Reader, e.sendPrvkey.PrivateKey, hType, hash.Sum(nil), opts)
}

func (e *encryptor) EncryptMessage(plain string) ([]byte, Signature, error) {
    binMsg := []byte(plain)
    crypted, err := e.recvPubkey.encrypt(binMsg); if err != nil {
        return nil, nil, err
    }
    sig, err := e.generateSignature(binMsg); if err != nil {
        return nil, nil, err
    }
    return crypted, sig, nil
}

