package discovery

import (
    "testing"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/crypt"
    "github.com/stkim1/pc-node-agent/slcontext"
    "bytes"
    msconfig "github.com/stkim1/pc-core/config"
)

var masterBoundAgentName, _ = msconfig.MasterHostSerial()
const slaveNodeName string = "pc-node1"

var aeskey []byte = []byte("longer means more possible keys ")
var aesenc, _ = crypt.NewAESCrypto(aeskey)

var initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")


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


func testSlavePublicKey() []byte {
    return []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwTd+iJMPFWGUENpOxJsw
jvZMgW0O/pQXJpN5miRdJcjk7ajBEDoo/NJPf5El60sp3+F+VT82ROab3nQknU5b
XXfinny0yvC2JNaORzPh7P8UzPGljXgxOfb+++tgEgSFI5WnBeMcima/Ce7M2AxS
WoAHxmu9AroBc33OMgg7TCpqqyWqbSIlnkizKi7IDopp1F0q92xPQFhFJne9IVDH
Opxdk8aFik3aWunVOla2olc/Vn+rs+J0i9+Kn8e4bHe0M5kGNx5+P/0OD37XsfVy
zNGJRIE1O1DIJ2ZHVtole4mtAt3C8d9lUrI7BacLBjDS1iZ/6kksO9Jhr4Lj9u4+
DQIDAQAB
-----END PUBLIC KEY-----`)
}


func testSlavePrivateKey() []byte {
    return []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAwTd+iJMPFWGUENpOxJswjvZMgW0O/pQXJpN5miRdJcjk7ajB
EDoo/NJPf5El60sp3+F+VT82ROab3nQknU5bXXfinny0yvC2JNaORzPh7P8UzPGl
jXgxOfb+++tgEgSFI5WnBeMcima/Ce7M2AxSWoAHxmu9AroBc33OMgg7TCpqqyWq
bSIlnkizKi7IDopp1F0q92xPQFhFJne9IVDHOpxdk8aFik3aWunVOla2olc/Vn+r
s+J0i9+Kn8e4bHe0M5kGNx5+P/0OD37XsfVyzNGJRIE1O1DIJ2ZHVtole4mtAt3C
8d9lUrI7BacLBjDS1iZ/6kksO9Jhr4Lj9u4+DQIDAQABAoIBAFwX3EqydV0GjnFd
7G9PXNy3To3d8mirI0GyxyIONQuebmdMqQDYB9NBVr0B7OXyhHn+W528LFy44hAs
oYsM3wV07+IEpJOaGecDEPulIgk5J6vrfbIpWKU9MhnW/Yp49xCX8u0ea+sXv/S3
CpHrhZE3Nv1/Oq7DA5ANpas5OzI4rHN1n1PUrwUbqFd8EazKUD4n+TD7aUPG5mMj
k4H2BEcNpfr4aXbjRqoFLzjr3RXQdiaWKfmRe3ZpBIF/iaCr1re2vs/rhnIw2HX7
3qk03Eyj5vu/LZiko/jxeYOfaTnackyZ6vkzjD7rVxMKBAwPqDfgvXK2P/U5HEAF
KDaVhc0CgYEAygvyh2ccYRcEJxkyNaGY0yhK5OclOwjeRqUHkcFloZ59ymkVLKnB
5xtrgsFksWU+DYOyQURbGofnLOvNd+LKrgWomnWh4DHHvVMXui8zavGdgB3P1FdR
YXA+kda+DCHST2OmVSxidneqKFKNLt9lozA1amba3X/V4lD/Llzyuu8CgYEA9M/t
1pKyVvMD9Jblb13jipaY/sOONHS2FvoO+YrNAqwXx9VmYVUUTx/T3z5+ujZbeJ+J
nEx5I/nSAQZuY2IJ/3RmRD+cgszGEuDeocBTZY73yUM6XKexo11pZk6xKqP+c5Gg
csDWQmM30c4lJwx7DJNfDCA+jCN+aEPoqq3bxsMCgYBJPDllwQc1Xg1gSq67Z96o
M0OqYupI0rcW7jynJW28PmGkG6DUNpgVOAgpNgZUkrkCVwkmxSssm7Q8wSAR43/J
wj1R9298fy7CPjssfm1pxzhqtuOdOSVDZ1cWr7rlVOERa7Jfzx3FiSyBPyLzqYAC
vbeu4KdWgD67sNY+LOzCuwKBgF2UzjnjwcBzDOQGepXjsgNcJgfdARMUOjb2R5sk
b9HBryV4cbZrK2RDql4AKblM5hJqCdRxdy1FZf12U+Qxqdi4yg70sgNd+6ljxDbY
qgh8akPJKxoYEFN+dbfiBN9j6PSMimTTShP+kWvl/VW7852PCBo+iSpQtxVsQBhe
dVC5AoGAOqrY/248onrzhEsENCBUl26CIuV0k+oe5De4vzpU+kUM/J4kSYgMcZRj
yHP1nQ1E0mWOV5sWcNJnU8ZhJt3M2pxhGp7bRGl2c8FxhAvWyYSRXtOILremUAhN
3nsy4axG/khKFn3jHJ1WQxxy7aqTrFJKNaDZB+YQ0rsLL6PsTtA=
-----END RSA PRIVATE KEY-----`)
}


func masterIdentityInqueryRespond() (meta *msagent.PocketMasterAgentMeta, err error) {
    // ------------- Let's Suppose you've sent an unbounded inquery from a node over multicast net ---------------------
    ua, err := slagent.UnboundedMasterSearchDiscovery()
    if err != nil {
        return
    }
    psm, err := slagent.PackedSlaveMeta(slagent.UnboundedMasterSearchMeta(ua))
    if err != nil {
        return
    }
    // -------------- over master, it's received the message and need to make an inquiry "Who R U"? --------------------
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        return
    }
    cmd, err := msagent.SlaveIdentityInqueryRespond(usm.DiscoveryAgent)
    if err != nil {
        return
    }
    meta = msagent.SlaveIdentityInquiryMeta(cmd)
    return
}

func masterIdentityFixationRespond() (meta *msagent.PocketMasterAgentMeta, err error) {
    agent, err := slagent.MasterAnswerInquiryStatus(initSendTimestmap)
    if err != nil {
        return
    }
    msa, err := slagent.MasterAnswerInquiryAgent(agent)
    if err != nil {
        return
    }
    cmd, err := msagent.MasterDeclarationCommand(msa.StatusAgent, initSendTimestmap.Add(time.Second))
    if err != nil {
        return
    }
    meta = msagent.MasterDeclarationMeta(cmd, testMasterPublicKey())
    return
}

func masterKeyExchangeCommand() (meta *msagent.PocketMasterAgentMeta, err error) {
    agent, err := slagent.KeyExchangeStatus(masterBoundAgentName, initSendTimestmap)
    if err != nil {
        return
    }
    msa, err := slagent.KeyExchangeMeta(agent, testSlavePublicKey())
    if err != nil {
        return
    }
    psm, err := slagent.PackedSlaveMeta(msa)
    if err != nil {
        return
    }
    //-------------- over master, we've received the message ----------------------
    // suppose we've sort out what this is.
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        return
    }
    // master preperation
    timestmap := initSendTimestmap.Add(time.Second)
    // encryptor
    rsaenc ,err := crypt.NewEncryptorFromKeyData(usm.SlavePubKey, testMasterPrivateKey())
    if err != nil {
        return
    }
    // responding commnad
    cmd, slvstat, err := msagent.ExchangeCryptoKeyAndNameCommand(usm.StatusAgent, slaveNodeName, timestmap)
    if err != nil {
        return
    }
    meta, err = msagent.ExchangeCryptoKeyAndNameMeta(cmd, slvstat, aeskey, aesenc, rsaenc)
    return
}

func masterCryptoCheckCommand() (meta *msagent.PocketMasterAgentMeta, err error) {
    agent, err := slagent.SlaveBindReadyStatus(masterBoundAgentName, slaveNodeName, initSendTimestmap)
    if err != nil {
        return
    }
    msa, err := slagent.SlaveBindReadyMeta(agent, aesenc)
    if err != nil {
        return
    }
    //-------------- over master, we've received the message ----------------------
    mdsa, err := aesenc.Decrypt(msa.EncryptedStatus)
    if err != nil {
        return
    }
    // unmarshaled, slave-status
    ussa, err := slagent.UnpackedSlaveStatus(mdsa)
    if err != nil {
        return
    }
    // master preperation
    timestmap := initSendTimestmap.Add(time.Second)
    if err != nil {
        return
    }
    // master crypto check state command
    cmd, err := msagent.MasterBindReadyCommand(ussa, timestmap)
    if err != nil {
        return
    }
    meta, err = msagent.MasterBindReadyMeta(cmd, aesenc)
    return
}

//--

func TestUnboundedState_InquiredTransition(t *testing.T) {
    meta, err := masterIdentityInqueryRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }

    context := slcontext.DebugSlaveContext(testSlavePublicKey(), testSlavePrivateKey())
    ssd := NewSlaveDiscovery(context)
    err = ssd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second))
    if err != nil {
        t.Error(err.Error())
        return
    }

    if ssd.CurrentState() != SlaveInquired {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", ssd.CurrentState().String())
        return
    }
}

func TestInquired_KeyExchangeTransition(t *testing.T) {
    meta, err := masterIdentityFixationRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    context := slcontext.DebugSlaveContext(testSlavePublicKey(), testSlavePrivateKey())
    sd := NewSlaveDiscovery(context)
    sd.(*slaveDiscovery).discoveryState = SlaveInquired

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2)); err != nil {
        t.Errorf(err.Error())
        return
    }
    if sd.CurrentState() != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }
}

func TestKeyExchange_CryptoCheckTransition(t *testing.T) {
    meta, err := masterIdentityFixationRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    context := slcontext.DebugSlaveContext(testSlavePublicKey(), testSlavePrivateKey())
    sd := NewSlaveDiscovery(context)
    sd.(*slaveDiscovery).discoveryState = SlaveInquired

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2)); err != nil {
        t.Errorf(err.Error())
        return
    }
    if sd.CurrentState() != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    // get master meta with aeskey
    meta, err = masterKeyExchangeCommand()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 3)); err != nil {
        t.Error(err.Error())
        return
    }
    if sd.CurrentState() != SlaveCryptoCheck {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }


    // Verification
    if msName, _ := context.GetMasterAgent(); msName != masterBoundAgentName {
        t.Errorf("[ERR] master node name is setup inappropriately | Current : %s\n", msName)
        return
    }
    if snName, _ := context.GetSlaveNodeName(); snName != slaveNodeName {
        t.Errorf("[ERR] slave node name is setup inappropriately | Current : %s\n", snName)
        return
    }
    if bytes.Compare(context.GetAESKey(), aeskey) != 0 {
        t.Errorf("[ERR] slave aes key is setup inappropriately")
        return
    }
}

func TestCryptoCheck_BoundedTransition(t *testing.T) {
    meta, err := masterIdentityFixationRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    context := slcontext.DebugSlaveContext(testSlavePublicKey(), testSlavePrivateKey())
    sd := NewSlaveDiscovery(context)
    sd.(*slaveDiscovery).discoveryState = SlaveInquired

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2)); err != nil {
        t.Errorf(err.Error())
        return
    }

    // get master meta with aeskey
    meta, err = masterKeyExchangeCommand()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 3)); err != nil {
        t.Error(err.Error())
        return
    }
    if sd.CurrentState() != SlaveCryptoCheck {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    // get master bind ready
    meta, err = masterCryptoCheckCommand()
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 4)); err != nil {
        t.Error(err.Error())
        return
    }

    if sd.CurrentState() != SlaveBounded {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }

    // Verification
    if msName, _ := context.GetMasterAgent(); msName != masterBoundAgentName {
        t.Errorf("[ERR] master node name is setup inappropriately | Current : %s\n", msName)
        return
    }
    if snName, _ := context.GetSlaveNodeName(); snName != slaveNodeName {
        t.Errorf("[ERR] slave node name is setup inappropriately | Current : %s\n", snName)
        return
    }
    if bytes.Compare(context.GetAESKey(), aeskey) != 0 {
        t.Errorf("[ERR] slave aes key is setup inappropriately")
        return
    }
}