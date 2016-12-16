package pcrypto

import (
    "testing"
    "os"
    "io/ioutil"
    "encoding/base64"
    "reflect"
    "fmt"
    "crypto/tls"
)

func TestKeyGeneration(t *testing.T) {
    var toSign []byte = []byte("date: Thu, 05 Jan 2012 21:31:40 GMT")
    if err := GenerateWeakKeyPairFiles("test.pub", "test.pem", "test.ssh"); err != nil {
        t.Errorf("failed to generate a key pair %v", err)
    }

    // sign the mssage
    signer, err := NewSignerFromPrivateKeyFile("test.pem"); if err != nil {
        t.Errorf("signer is damaged: %v", err)
    }
    signed, err := signer.Sign(toSign); if err != nil {
        t.Errorf("could not sign request: %v", err)
    }

    // unsigned the message
    parser, perr := NewUnsignerFromPublicKeyFile("test.pub"); if perr != nil {
        t.Errorf("could not sign request: %v", err)
    }
    if err = parser.Unsign(toSign, signed); err != nil {
        t.Errorf("could not unsign request: %v", err)
    }

    os.Remove("test.pem");os.Remove("test.pub");os.Remove("test.ssh")
}

func TestSignatureGeneration(t *testing.T) {
    if err := ioutil.WriteFile("test.pub", TestMasterPublicKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write public key %v", err)
    }
    if err := ioutil.WriteFile("test.pem", TestMasterPrivateKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write private key %v", err)
    }

    // generate signature
    sig, err := GenerateSignature("test.pub", "test.pem"); if err != nil {
        t.Error(err.Error())
    }
    if len(sig) == 0 {
        t.Error("Empty Signature generated")
    }
    if base64.StdEncoding.EncodeToString(sig) != TestKeySignature {
        t.Error("Wrong Signature generated")
    }

    os.Remove("test.pem");os.Remove("test.pub")
}

func TestSignatureVerification(t *testing.T) {
    var orgMsg []byte = TestMasterPublicKey()
    if err := ioutil.WriteFile("test.pub", TestMasterPublicKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write public key %v", err)
    }
    if err := ioutil.WriteFile("test.pem", TestMasterPrivateKey(), os.ModePerm); err != nil {
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
    var orgMsg []byte = []byte("date: Thu, 05 Jan 2012 21:31:40 GMT")
    if err := ioutil.WriteFile("test.pub", TestMasterPublicKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write public key %v", err)
    }
    if err := ioutil.WriteFile("test.pem", TestMasterPrivateKey(), os.ModePerm); err != nil {
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

func TestEncDecMessageWithFile(t *testing.T) {
    var orgMsg []byte = []byte("date: Thu, 05 Jan 2012 21:31:40 GMT")
    if err := ioutil.WriteFile("sendtest.pub", TestMasterPublicKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write public key %v", err)
    }
    if err := ioutil.WriteFile("sendtest.pem", TestMasterPrivateKey(), os.ModePerm); err != nil {
        t.Errorf("Fail to write private key %v", err)
    }
    if err := GenerateWeakKeyPairFiles("recvtest.pub", "recvtest.pem", "recvtest.ssh"); err != nil {
        t.Errorf("failed to generate a key pair %v", err)
    }

    // encryptor
    encr ,err := NewRsaEncryptorFromKeyFiles("recvtest.pub", "sendtest.pem"); if  err != nil {
        t.Errorf(err.Error())
    }
    crypted, sig, err := encr.EncryptByRSA(orgMsg); if err != nil {
        t.Errorf(err.Error())
    }
    // decryptor
    decr, err := NewRsaDecryptorFromKeyFiles("sendtest.pub", "recvtest.pem"); if err != nil {
        t.Errorf(err.Error())
    }
    plain, err := decr.DecryptByRSA(crypted, sig); if err != nil {
        t.Errorf(err.Error())
    }
    // comp
    if !reflect.DeepEqual(orgMsg, plain) {
        t.Error("Original Message and Decrypted message are different" + string(plain))
    }

    os.Remove("sendtest.pem");os.Remove("sendtest.pub")
    os.Remove("recvtest.pem");os.Remove("recvtest.pub");os.Remove("recvtest.ssh")
}

func TestEncDecMessageWithData(t *testing.T) {
    var orgMsg []byte = []byte("date: Thu, 05 Jan 2012 21:31:40 GMT")
    if err := GenerateWeakKeyPairFiles("recvtest.pub", "recvtest.pem", "recvtest.ssh"); err != nil {
        t.Errorf("failed to generate a key pair %v", err)
    }
    sendTestPubKey := TestMasterPublicKey()
    sendTestPrvKey := TestMasterPrivateKey()
    recvTestPubKey, err := ioutil.ReadFile("recvtest.pub")
    if err != nil {
        t.Error(err.Error())
    }
    recvTestPrvKey, err := ioutil.ReadFile("recvtest.pem")
    if err != nil {
        t.Error(err.Error())
    }

    // encryptor
    encr ,err := NewRsaEncryptorFromKeyData(recvTestPubKey, sendTestPrvKey); if  err != nil {
        t.Errorf(err.Error())
    }
    crypted, sig, err := encr.EncryptByRSA(orgMsg); if err != nil {
        t.Errorf(err.Error())
    }
    // decryptor
    decr, err := NewRsaDecryptorFromKeyData(sendTestPubKey, recvTestPrvKey); if err != nil {
        t.Errorf(err.Error())
    }
    plain, err := decr.DecryptByRSA(crypted, sig); if err != nil {
        t.Errorf(err.Error())
    }
    // comp
    if !reflect.DeepEqual(orgMsg, plain) {
        t.Error("Original Message and Decrypted message are different" + string(plain))
    }
    os.Remove("recvtest.pem");os.Remove("recvtest.pub");os.Remove("recvtest.ssh")
}

func TestAESEncyptionDecryption(t *testing.T) {
    key := []byte("longer means more possible keys ")
    text := []byte("This is the unecrypted data. Referring to it as plain text.")

    ac, err := NewAESCrypto(key); if err != nil {
        t.Errorf("Cannot create AES cryptor %v", err)
    }
    crypted, err := ac.EncryptByAES(text); if err != nil {
        t.Errorf("Cannot encrypt message with AES %v", err)
    }
    plain, err := ac.DecryptByAES(crypted); if err != nil {
        t.Errorf("Cannot decrypt message with AES %v", err)
    }
    if string(plain) != string(text) {
        t.Errorf("Orinal and decrypted are different")
    }
}

func TestAESKeyGeneration(t *testing.T) {
    var key1 []byte = NewAESKey32Byte()
    var key2 []byte = NewAESKey32Byte()

    // `go test -v ./…` to view log
    t.Log("Key 1 - randBytesWithMask : " + string(key1))
    t.Log("Key 2 - randBytesWithMask : " + string(key2))
    t.Log("Key 3 - randBytesWithMaskSrc : " + string(randBytesWithMaskSrc(32)))
    t.Log("Key 4 - randBytesWithMaskSrc: " + string(randBytesWithMaskSrc(32)))
    key, _ := randCryptoBytes(32)
    t.Log("Key 5 - randCryptoBytes: " + string(key))
    key, _ = randCryptoBytes(32)
    t.Log("Key 6 - randCryptoBytes: " + string(key))

    if reflect.DeepEqual(key1, key2) {
        t.Error("Randome AES Keys are not different enough")
    }
}

func BenchmarkRandBytesWithMask(b *testing.B) {
    for i := 0; i < b.N; i++ {
        randBytesWithMask(32)
    }
}

// this is 10 times slower than BenchmarkRandBytesWithMask
func BenchmarkRandCryptoByte(b *testing.B) {
    for i := 0; i < b.N; i++ {
        randCryptoBytes(32)
    }
}

func ExampleWeakRsaKeyEncryption() {
    // master key pair
    mpub, mprv, _, merr := GenerateWeakKeyPair()
    if merr != nil {
        fmt.Printf(merr.Error())
        return
    }
    //slave key pair
    spub, sprv, _, serr := GenerateWeakKeyPair()
    if serr != nil {
        fmt.Printf(serr.Error())
        return
    }

    // encryptor
    encr ,err := NewRsaEncryptorFromKeyData(spub, mprv)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // decryptor
    decr, err := NewRsaDecryptorFromKeyData(mpub, sprv)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    // descryption
    crypted, sig, err := encr.EncryptByRSA(TestAESKey)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    plain, err := decr.DecryptByRSA(crypted, sig);
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // comp
    if !reflect.DeepEqual(TestAESKey, plain) {
        fmt.Printf("[ERR] Unidentical original Message and Decrypted message" + string(plain))
        return
    }
    fmt.Printf("Original Message Size %d | Encrypted Message Size %d | Signature Size %d", len(TestAESKey), len(crypted), len(sig))
    // Output:
    // Original Message Size 32 | Encrypted Message Size 128 | Signature Size 128
}

func ExampleStrongRsaKeyEncryption() {
    // master key pair
    mpub, mprv, _, merr := GenerateStrongKeyPair()
    if merr != nil {
        fmt.Printf(merr.Error())
        return
    }
    //slave key pair
    spub, sprv, _, serr := GenerateStrongKeyPair()
    if serr != nil {
        fmt.Printf(serr.Error())
        return
    }

    // encryptor
    encr ,err := NewRsaEncryptorFromKeyData(spub, mprv)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // decryptor
    decr, err := NewRsaDecryptorFromKeyData(mpub, sprv)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }

    // descryption
    crypted, sig, err := encr.EncryptByRSA(TestAESKey)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    plain, err := decr.DecryptByRSA(crypted, sig);
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // comp
    if !reflect.DeepEqual(TestAESKey, plain) {
        fmt.Printf("[ERR] Unidentical original Message and Decrypted message" + string(plain))
        return
    }
    fmt.Printf("Original Message Size %d | Encrypted Message Size %d | Signature Size %d", len(TestAESKey), len(crypted), len(sig))
    // Output:
    // Original Message Size 32 | Encrypted Message Size 256 | Signature Size 256
}

func TestLoadStrongX509KeyPair(t *testing.T) {
    if err := GenerateClusterCertificateAuthorityFiles("recvtest.pub", "recvtest.pem", "recvtest.cert", "cluster-id-here", "KR"); err != nil {
        t.Errorf("failed to generate a key pair %v", err)
    }

    _, err := tls.LoadX509KeyPair("recvtest.cert", "recvtest.pem")
    if !os.IsNotExist(err) {
        t.Log("[INFO] File does not exists")
    }
    if err != nil {
        t.Error(err.Error())
    }

    os.Remove("recvtest.pub");os.Remove("recvtest.pem");os.Remove("recvtest.cert");
}
