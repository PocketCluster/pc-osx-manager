package crypt

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
)

type RsaEncryptor interface {
    EncryptMessage(plain []byte) (crypted []byte, sig Signature, err error)
}

func NewEncryptorFromKeyFiles(recvPubkeyPath, sendPrvkeyPath string) (RsaEncryptor, error) {
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

func NewEncryptorFromKeyData(recvPubkeyData, sendPrvkeyData []byte) (RsaEncryptor, error) {
    pubkey, err := newPublicKeyFromData(recvPubkeyData); if err != nil {
        return nil, err
    }
    prvkey, err := newPrivateKeyFromData(sendPrvkeyData); if err != nil {
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

func (e *encryptor) EncryptMessage(plain []byte) (crypted []byte, sig Signature, err error) {
    crypted, err = e.recvPubkey.encrypt(plain); if err != nil {
        return nil, nil, err
    }
    sig, err = e.generateSignature(plain); if err != nil {
        return nil, nil, err
    }
    return crypted, sig, nil
}

