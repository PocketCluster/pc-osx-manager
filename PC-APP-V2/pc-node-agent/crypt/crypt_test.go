package crypt

import (
    "testing"
    "os"
    "log"
    "encoding/base64"
)

// loadTestPublicKey loads an parses a PEM encoded public key file.
func loadTestPublicKey(path string) (Unsigner, error) {
    return parsePublicKey([]byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`))
}

func loadTestPrivateKey(path string) (Signer, error) {
    return parsePrivateKey([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDCFENGw33yGihy92pDjZQhl0C36rPJj+CvfSC8+q28hxA161QF
NUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6Z4UMR7EOcpfdUE9Hf3m/hs+F
UR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJwoYi+1hqp1fIekaxsyQIDAQAB
AoGBAJR8ZkCUvx5kzv+utdl7T5MnordT1TvoXXJGXK7ZZ+UuvMNUCdN2QPc4sBiA
QWvLw1cSKt5DsKZ8UETpYPy8pPYnnDEz2dDYiaew9+xEpubyeW2oH4Zx71wqBtOK
kqwrXa/pzdpiucRRjk6vE6YY7EBBs/g7uanVpGibOVAEsqH1AkEA7DkjVH28WDUg
f1nqvfn2Kj6CT7nIcE3jGJsZZ7zlZmBmHFDONMLUrXR/Zm3pR5m0tCmBqa5RK95u
412jt1dPIwJBANJT3v8pnkth48bQo/fKel6uEYyboRtA5/uHuHkZ6FQF7OUkGogc
mSJluOdc5t6hI1VsLn0QZEjQZMEOWr+wKSMCQQCC4kXJEsHAve77oP6HtG/IiEn7
kpyUXRNvFsDE0czpJJBvL/aRFUJxuRK91jhjC68sA7NsKMGg5OXb5I5Jj36xAkEA
gIT7aFOYBFwGgQAQkWNKLvySgKbAZRTeLBacpHMuQdl1DfdntvAyqpAZ0lY0RKmW
G6aFKaqQfOXKCyWoUiVknQJAXrlgySFci/2ueKlIE1QqIiLSZ8V8OlpFLRnb1pzI
7U1yQXnTAEFYM560yJlzUpOb1V4cScGd365tiSMvxLOvTA==
-----END RSA PRIVATE KEY-----`))
}

func TestKeyGeneration(t *testing.T) {
    const toSign string = "date: Thu, 05 Jan 2012 21:31:40 GMT"
    if err := GenerateKeyPair("test.pub", "test.pem"); err != nil {
        t.Errorf("failed to generate a key pair %v", err)
    }

    // sign the mssage
    signer, err := NewSignerFromPrivateKeyFile("test.pem"); if err != nil {
        t.Errorf("signer is damaged: %v", err)
    }
    signed, err := signer.Sign([]byte(toSign)); if err != nil {
        t.Errorf("could not sign request: %v", err)
    }

    // print test key
    sig := base64.StdEncoding.EncodeToString(signed)
    log.Printf("Signature: %v\n", sig)

    // unsigned the message
    parser, perr := NewUnsignerFromPublicKeyFile("test.pub"); if perr != nil {
        t.Errorf("could not sign request: %v", err)
    }
    if err = parser.Unsign([]byte(toSign), signed); err != nil {
        t.Errorf("could not unsign request: %v", err)
    }

    os.Remove("test.pem");os.Remove("test.pub")
}