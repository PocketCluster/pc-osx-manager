package msagent

import (
    "fmt"
    "time"
    "testing"
    "os"
    "io/ioutil"
    "bytes"

    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-node-agent/slagent"
)

func ExampleUnboundedInqueryMeta() {
    // Let's Suppose you've received an unbounded inquery from a node over multicast net.
    ua, err := slagent.UnboundedBroadcastAgent()
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    psm, err := slagent.PackedSlaveMeta(slagent.DiscoveryMetaAgent(ua))
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    //-------------- over master, we've received the message and need to make an inquiry "Who R U"? --------------------
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // TODO : we need ways to identify if what this package is
    cmd, err := IdentityInqueryRespond(usm.DiscoveryAgent)
    meta := UnboundedInqueryMeta(cmd)
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // msgpack verfication
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MetaVersion : %s\n",                         meta.MetaVersion)
    fmt.Printf("DiscoveryRespond.Version : %s\n",            meta.DiscoveryRespond.Version)
    fmt.Printf("DiscoveryRespond.MasterBoundAgent : %s\n",   meta.DiscoveryRespond.MasterBoundAgent)
    fmt.Printf("DiscoveryRespond.MasterCommandType : %s\n",  meta.DiscoveryRespond.MasterCommandType)
    fmt.Printf("DiscoveryRespond.MasterAddress : %s\n",      meta.DiscoveryRespond.MasterAddress)
    fmt.Print("------------------\n")
    fmt.Printf("MsgPack Length : %d\n", len(mp))
    fmt.Print("------------------\n")
    fmt.Printf("MetaVersion : %s\n",                         umeta.MetaVersion)
    fmt.Printf("DiscoveryRespond.Version : %s\n",            umeta.DiscoveryRespond.Version)
    fmt.Printf("DiscoveryRespond.MasterBoundAgent : %s\n",   umeta.DiscoveryRespond.MasterBoundAgent)
    fmt.Printf("DiscoveryRespond.MasterCommandType : %s\n",  umeta.DiscoveryRespond.MasterCommandType)
    fmt.Printf("DiscoveryRespond.MasterAddress : %s\n",      umeta.DiscoveryRespond.MasterAddress)
    // Output:
    // MetaVersion : 1.0.1
    // DiscoveryRespond.Version : 1.0.1
    // DiscoveryRespond.MasterBoundAgent : C02QF026G8WL
    // DiscoveryRespond.MasterCommandType : pc_ms_wr
    // DiscoveryRespond.MasterAddress : 192.168.1.236
    // ------------------
    // MsgPack Length : 164
    // ------------------
    // MetaVersion : 1.0.1
    // DiscoveryRespond.Version : 1.0.1
    // DiscoveryRespond.MasterBoundAgent : C02QF026G8WL
    // DiscoveryRespond.MasterCommandType : pc_ms_wr
    // DiscoveryRespond.MasterAddress : 192.168.1.236
}

func testMasterPublicKey() []byte {
    return []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`)
}

func testMasterPrivateKey() []byte {
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

func ExampleIdentityInqeuryMeta() {
    // suppose slave agent has answered question who it is
    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    agent, err := slagent.InquiredAgent(timestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    msa, err := slagent.InquiredMetaAgent(agent)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    mpsm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(mpsm)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    timestmap, err = time.Parse(time.RFC3339, "2012-11-01T22:08:42+00:00")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    cmd, err := MasterIdentityRevealCommand(usm.StatusAgent, timestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    meta := IdentityInqueryMeta(cmd, testMasterPublicKey())
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MetaVersion : %s\n",                     meta.MetaVersion)
    fmt.Printf("StatusCommand.Version : %s\n",           meta.StatusCommand.Version)
    fmt.Printf("StatusCommand.MasterBoundAgent : %s\n",  meta.StatusCommand.MasterBoundAgent)
    fmt.Printf("StatusCommand.MasterCommandType : %s\n", meta.StatusCommand.MasterCommandType)
    fmt.Printf("StatusCommand.MasterAddress : %s\n",     meta.StatusCommand.MasterAddress)
    fmt.Printf("StatusCommand.MasterTimestamp : %s\n",   meta.StatusCommand.MasterTimestamp.String())
    fmt.Print("------------------\n")
    fmt.Printf("MsgPack Length : %d / pubkey Length : %d\n", len(mp), len(umeta.MasterPubkey))
    fmt.Print("------------------\n")
    fmt.Printf("MetaVersion : %s\n",                     meta.MetaVersion)
    fmt.Printf("StatusCommand.Version : %s\n",           meta.StatusCommand.Version)
    fmt.Printf("StatusCommand.MasterBoundAgent : %s\n",  meta.StatusCommand.MasterBoundAgent)
    fmt.Printf("StatusCommand.MasterCommandType : %s\n", meta.StatusCommand.MasterCommandType)
    fmt.Printf("StatusCommand.MasterAddress : %s\n",     meta.StatusCommand.MasterAddress)
    fmt.Printf("StatusCommand.MasterTimestamp : %s\n",   meta.StatusCommand.MasterTimestamp.String())
    // Output:
    // MetaVersion : 1.0.1
    // StatusCommand.Version : 1.0.1
    // StatusCommand.MasterBoundAgent : C02QF026G8WL
    // StatusCommand.MasterCommandType : pc_ms_sp
    // StatusCommand.MasterAddress : 192.168.1.236
    // StatusCommand.MasterTimestamp : 2012-11-01 22:08:42 +0000 +0000
    // ------------------
    // MsgPack Length : 453 / pubkey Length : 271
    // ------------------
    // MetaVersion : 1.0.1
    // StatusCommand.Version : 1.0.1
    // StatusCommand.MasterBoundAgent : C02QF026G8WL
    // StatusCommand.MasterCommandType : pc_ms_sp
    // StatusCommand.MasterAddress : 192.168.1.236
    // StatusCommand.MasterTimestamp : 2012-11-01 22:08:42 +0000 +0000
}

func testSlavePublicKey() []byte {
    return []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCqGKukO1De7zhZj6+H0qtjTkVx
wTCpvKe4eCZ0FPqri0cb2JZfXJ/DgYSF6vUpwmJG8wVQZKjeGcjDOL5UlsuusFnc
CzWBQ7RKNUSesmQRMSGkVb1/3j+skZ6UtW+5u09lHNsj6tQ51s1SPrCBkedbNf0T
p0GbMJDyR4e9T04ZZwIDAQAB
-----END PUBLIC KEY-----`)
}

func TestExecKeyExchangeMeta(t *testing.T) {
    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }
    agent, err := slagent.KeyExchangeAgent("master-yoda", timestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    msa, err := slagent.KeyExchangeMetaAgent(agent, testSlavePublicKey())
    if err != nil {
        t.Error(err.Error())
        return
    }
    psm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        t.Error(err.Error())
        return
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // master preperation
    timestmap, err = time.Parse(time.RFC3339, "2012-11-01T22:08:42+00:00")
    if err != nil {
        t.Error(err.Error())
        return
    }
    aeskey := []byte("longer means more possible keys ")
    aesenc, err := crypt.NewAESCrypto(aeskey)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // encryptor
    err = ioutil.WriteFile("sendtest.pub", testMasterPublicKey(), os.ModePerm)
    defer os.Remove("sendtest.pem")
    if err != nil {
        t.Errorf("Fail to write public key %v", err)
        return
    }
    err = ioutil.WriteFile("sendtest.pem", testMasterPrivateKey(), os.ModePerm)
    defer os.Remove("sendtest.pub")
    if err != nil {
        t.Errorf("Fail to write private key %v", err)
        return
    }
    err = ioutil.WriteFile("recvtest.pub", testSlavePublicKey(), os.ModePerm)
    defer os.Remove("recvtest.pub")
    if err != nil {
        t.Errorf("Fail to write private key %v", err)
        return
    }
    rsaenc ,err := crypt.NewEncryptorFromKeyFiles("recvtest.pub", "sendtest.pem")
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // responding commnad
    cmd, slvstat, err := CryptoKeyAndNameSetCommand(usm.StatusAgent, "pc-node1", timestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    meta, err := ExecKeyExchangeMeta(cmd, slvstat, aeskey, aesenc, rsaenc)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // verification step
    mp, err := PackedMasterMeta(meta)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if len(mp) != 481 {
        t.Errorf("[ERR] package meta message size [%d] does not match an expectant", len(mp))
        return
    }
    // verification step
    umeta, err := UnpackedMasterMeta(mp)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    if meta.MetaVersion != umeta.MetaVersion {
        t.Errorf("[ERR] package/unpacked meta version differs")
        return
    }
    if bytes.Compare(meta.EncryptedMasterCommand, umeta.EncryptedMasterCommand) != 0{
        t.Errorf("[ERR] package/unpacked encrypted command differs")
        return
    }
    if bytes.Compare(meta.EncryptedSlaveStatus, umeta.EncryptedSlaveStatus) != 0{
        t.Errorf("[ERR] package/unpacked encrypted slave response differs")
        return
    }
    if bytes.Compare(meta.EncryptedAESKey, umeta.EncryptedAESKey) != 0{
        t.Errorf("[ERR] package/unpacked encrypted aes key differs")
        return
    }
    if bytes.Compare(meta.RsaCryptoSignature, umeta.RsaCryptoSignature) != 0{
        t.Errorf("[ERR] package/unpacked encryption signature differs")
        return
    }
}

func ExampleSendCryptoCheckMeta() {
}

func ExampleBoundedStatusMeta() {
}

func ExampleBindBrokenMeta() {
}