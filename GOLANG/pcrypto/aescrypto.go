package pcrypto

import (
    "crypto/aes"
    "crypto/cipher"
    cryptorand "crypto/rand"
    "fmt"
    "io"
    mathrand "math/rand"
    "time"
)

type AESCryptor interface {
    EncryptByAES(data []byte) (crypted []byte, err error)
    DecryptByAES(data []byte) (plain []byte, err error)
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

//------------------------------------------------ AES CRYPTOR ---------------------------------------------------------

type aesCrpytor struct {
    cipher.Block
}

func (ac *aesCrpytor) EncryptByAES(plain []byte) ([]byte, error) {
    var crypted []byte = make([]byte, aes.BlockSize + len(string(plain)))

    // iv : initialization vector
    iv := crypted[:aes.BlockSize]
    if _, err := io.ReadFull(cryptorand.Reader, iv); err != nil {
        return nil, err
    }

    cfb := cipher.NewCFBEncrypter(ac.Block, iv)
    cfb.XORKeyStream(crypted[aes.BlockSize:], plain)
    return crypted, nil
}

func (ac *aesCrpytor) DecryptByAES(crypted []byte) ([]byte,error) {

    var plain []byte
    if len(crypted) < aes.BlockSize {
        return nil, fmt.Errorf("[ERR] ciphertext too short")
    }

    iv := crypted[:aes.BlockSize]
    plain = make([]byte, len(crypted[aes.BlockSize:]))
    copy(plain, crypted[aes.BlockSize:])

    cfb := cipher.NewCFBDecrypter(ac.Block, iv)
    cfb.XORKeyStream(plain, plain)

    return plain, nil
}

// --- Random Key generation ---
// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
// letterBytes is expanded to 94 chars
const letterBytes = "!\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"
const (
    letterIdxBits = 7                        // 7 bits to represent a letter index (2^7 = 128 which is greater than 94)
    letterIdxMask = 1  << letterIdxBits - 1  // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits      // # of letter indices fitting in 7 bits
)

func randBytesWithMask(n int) []byte {
    b := make([]byte, n)
    // A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
    for i, cache, remain := n - 1, mathrand.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = mathrand.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }
    return b
}

// as mathrand.NewSource doesn't seem to thread safe, we'll avoid this
func randBytesWithMaskSrc(n int) []byte {
    b := make([]byte, n)
    var src = mathrand.NewSource(time.Now().UnixNano())
    // A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
    for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }
    return b
}

// this is crypto rand function generator
func randCryptoBytes(n int) ([]byte) {
    b := make([]byte, n)
    for {
        _, err := cryptorand.Read(b)
        // Note that err == nil only if we read len(b) bytes.
        if err == nil {
            break
        }
    }
    return b
}

func NewAESKey32Byte() []byte {
    // when Cryptography is asked, cryptography grade random string should be in place
    // return randBytesWithMask(32)
    return randCryptoBytes(32)
}