package crypt

import (
    "testing"
    "os"
    "io/ioutil"
    "encoding/base64"
)

// loadTestPublicKey loads an parses a PEM encoded public key file.
func testPublicKey() []byte {
    return []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`)
}

func testPrivateKey() []byte {
    return []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`)
}

func TestKeyGeneration(t *testing.T) {
    const toSign string = "date: Thu, 05 Jan 2012 21:31:40 GMT"
    if err := GenerateKeyPair("test.pub", "test.pem", "test.ssh", ); err != nil {
        t.Errorf("failed to generate a key pair %v", err)
    }

    // sign the mssage
    signer, err := NewSignerFromPrivateKeyFile("test.pem"); if err != nil {
        t.Errorf("signer is damaged: %v", err)
    }
    signed, err := signer.Sign([]byte(toSign)); if err != nil {
        t.Errorf("could not sign request: %v", err)
    }

    // unsigned the message
    parser, perr := NewUnsignerFromPublicKeyFile("test.pub"); if perr != nil {
        t.Errorf("could not sign request: %v", err)
    }
    if err = parser.Unsign([]byte(toSign), signed); err != nil {
        t.Errorf("could not unsign request: %v", err)
    }

    os.Remove("test.pem");os.Remove("test.pub");os.Remove("test.ssh")
}

func TestSignatureGeneration(t *testing.T) {
    const testSig string = "K9crXmaFVpvoXAEB/QUguOENDIJ2AhWTgbj8JAPAHbatQqaes19ycSaZgGCrg5NhGvgP13Wf/zR0ny6PR0V8FTUjzxaK2fGDElythqwW7QISyRPKFayRSNjGOC9/d74d31JB2/05Tuk4hksb5u90bc1y+t5RYMArDn8aJjx2GA8="
    if err := ioutil.WriteFile("test.pub", testPublicKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write public key %v", err)
    }
    if err := ioutil.WriteFile("test.pem", testPrivateKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write private key %v", err)
    }

    // generate signature
    sig, err := GenerateSignature("test.pub", "test.pem"); if err != nil {
        t.Error(err.Error())
    }
    if len(sig) == 0 {
        t.Error("Empty Signature generated")
    }
    if base64.StdEncoding.EncodeToString(sig) != testSig {
        t.Error("Wrong Signature generated")
    }

    os.Remove("test.pem");os.Remove("test.pub")
}

func TestSignatureVerification(t *testing.T) {
    var orgMsg string = string(testPublicKey())
    if err := ioutil.WriteFile("test.pub", testPublicKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write public key %v", err)
    }
    if err := ioutil.WriteFile("test.pem", testPrivateKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write private key %v", err)
    }

    // generate signature
    signature, err := GenerateSignature("test.pub", "test.pem"); if err != nil {
        t.Error(err.Error())
    }
    // verify message with signature
    if err = VerifySignature("test.pub", orgMsg, signature); err != nil {
        t.Errorf(err.Error())
    }

    os.Remove("test.pem");os.Remove("test.pub")
}

func TestMessageSigning(t *testing.T) {
    const orgMsg string = "date: Thu, 05 Jan 2012 21:31:40 GMT"
    if err := ioutil.WriteFile("test.pub", testPublicKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write public key %v", err)
    }
    if err := ioutil.WriteFile("test.pem", testPrivateKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write private key %v", err)
    }

    // message signing
    signed, err := SignMessageWithPrivateKeyFile("test.pem", orgMsg); if err != nil {
        t.Error(err.Error())
    }
    // verify message sign
    if err = VerifySignature("test.pub", orgMsg, signed); err != nil {
        t.Errorf(err.Error())
    }

    os.Remove("test.pem");os.Remove("test.pub")
}

func TestEncDecMessage(t *testing.T) {
    const orgMsg string = "date: Thu, 05 Jan 2012 21:31:40 GMT"
    if err := ioutil.WriteFile("sendtest.pub", testPublicKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write public key %v", err)
    }
    if err := ioutil.WriteFile("sendtest.pem", testPrivateKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write private key %v", err)
    }
    if err := GenerateKeyPair("recvtest.pub", "recvtest.pem", "recvtest.ssh", ); err != nil {
        t.Errorf("failed to generate a key pair %v", err)
    }

    // encryptor
    encr ,err := NewEncryptorFromKeyFiles("recvtest.pub", "sendtest.pem"); if  err != nil {
        t.Errorf(err.Error())
    }
    crypted, sig, err := encr.EncryptMessage(orgMsg); if err != nil {
        t.Errorf(err.Error())
    }
    // decryptor
    decr, err := NewDecryptorFromKeyFiles("sendtest.pub", "recvtest.pem"); if err != nil {
        t.Errorf(err.Error())
    }
    plain, err := decr.DecryptMessage(string(crypted), sig); if err != nil {
        t.Errorf(err.Error())
    }
    // comp
    if orgMsg != string(plain) {
        t.Error("Original Message and Decrypted message are different" + string(plain))
    }

    os.Remove("sendtest.pem");os.Remove("sendtest.pub")
    os.Remove("recvtest.pem");os.Remove("recvtest.pub");os.Remove("recvtest.ssh")
}