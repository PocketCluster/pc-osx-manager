package crypt

import (
    "crypto/aes"
    "crypto/cipher"
    "fmt"
)

type AESCryptor interface {
    Encrypt(data []byte) (crypted []byte, err error)
    Decrypt(data []byte) (plain []byte, err error)
}

func NewAESCrypto(key []byte) (crypto AESCryptor, err error) {
    // TODO : check if key's length is 2's complement
    if len(key) != 32 {
        err = fmt.Errorf("[ERR] key length is too short")
        return nil, err
    }

    var block cipher.Block
    if block, err = aes.NewCipher(key); err != nil {
        return nil, err
    }
    crypto = &aesCrpytor{Block:block}
    return crypto, err
}