package randstr

import (
    "encoding/binary"
    crand "crypto/rand"
)

func cryptoRandomInt64() int64 {
    rb := make([]byte, 8)
    for {
        if _, err := crand.Read(rb[:]); err == nil {
            return int64(binary.LittleEndian.Uint64(rb[:]))
        }
    }
}

/*
 * (03/25/2017)
 * this function is here to provide random string. It could be changed to even more random
 * string generator with one in crypto
 *
 * (03/04/2017)
 * math random doesn't provice enough randomness, we switched to crypto random
 */
func NewRandomString(length int) string {
    if length == 0 {
        return ""
    }
    const (
        letterBytes   = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ" // 62 characters
        letterIdxBits = 6                        // 6 bits to represent a letter index (2^6 = 64 which is greater than 62)
        letterIdxMask = 1  << letterIdxBits - 1  // All 1-bits, as many as letterIdxBits
        letterIdxMax  = 63 / letterIdxBits       // # of letter indices fitting in 6 bits
    )
    b := make([]byte, length)

    // A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
    for i, cache, remain := length - 1, cryptoRandomInt64(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = cryptoRandomInt64(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }
    return string(b)
}
