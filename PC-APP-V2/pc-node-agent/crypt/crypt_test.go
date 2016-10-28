package crypt

import (
    "testing"
    "os"
    "io/ioutil"
    "encoding/base64"
    "reflect"
)

func TestKeyGeneration(t *testing.T) {
    var toSign []byte = []byte("date: Thu, 05 Jan 2012 21:31:40 GMT")
    if err := GenerateKeyPair("test.pub", "test.pem", "test.ssh", ); err != nil {
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
    if err := GenerateKeyPair("recvtest.pub", "recvtest.pem", "recvtest.ssh"); err != nil {
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
    plain, err := decr.DecryptMessage(crypted, sig); if err != nil {
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
    if err := GenerateKeyPair("recvtest.pub", "recvtest.pem", "recvtest.ssh"); err != nil {
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
    encr ,err := NewEncryptorFromKeyData(recvTestPubKey, sendTestPrvKey); if  err != nil {
        t.Errorf(err.Error())
    }
    crypted, sig, err := encr.EncryptMessage(orgMsg); if err != nil {
        t.Errorf(err.Error())
    }
    // decryptor
    decr, err := NewDecryptorFromKeyData(sendTestPubKey, recvTestPrvKey); if err != nil {
        t.Errorf(err.Error())
    }
    plain, err := decr.DecryptMessage(crypted, sig); if err != nil {
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
    crypted, err := ac.Encrypt(text); if err != nil {
        t.Errorf("Cannot encrypt message with AES %v", err)
    }
    plain, err := ac.Decrypt(crypted); if err != nil {
        t.Errorf("Cannot decrypt message with AES %v", err)
    }
    if string(plain) != string(text) {
        t.Errorf("Orinal and decrypted are different")
    }
}
