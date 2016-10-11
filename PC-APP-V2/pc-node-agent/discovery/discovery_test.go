package discovery

import (
    "testing"
    "time"
    "io/ioutil"
    "os"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/msagent"
    "github.com/stkim1/pc-node-agent/crypt"
)

const masterBoundAgentName string = "master-yoda"
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
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxzYuc22QSst/dS7geYYK
5l5kLxU0tayNdixkEQ17ix+CUcUbKIsnyftZxaCYT46rQtXgCaYRdJcbB3hmyrOa
vkhTpX79xJZnQmfuamMbZBqitvscxW9zRR9tBUL6vdi/0rpoUwPMEh8+Bw7CgYR0
FK0DhWYBNDfe9HKcyZEv3max8Cdq18htxjEsdYO0iwzhtKRXomBWTdhD5ykd/fAC
VTr4+KEY+IeLvubHVmLUhbE5NgWXxrRpGasDqzKhCTmsa2Ysf712rl57SlH0Wz/M
r3F7aM9YpErzeYLrl0GhQr9BVJxOvXcVd4kmY+XkiCcrkyS1cnghnllh+LCwQu1s
YwIDAQAB
-----END PUBLIC KEY-----`)
}


func testSlavePrivateKey() []byte {
    return []byte(`-----BEGIN RSA PRIVATE KEY-----
Proc-Type: 4,ENCRYPTED
DEK-Info: DES-EDE3-CBC,32495A90F3FF199D
lrMAsSjjkKiRxGdgR8p5kZJj0AFgdWYa3OT2snIXnN5+/p7j13PSkseUcrAFyokc
V9pgeDfitAhb9lpdjxjjuxRcuQjBfmNVLPF9MFyNOvhrprGNukUh/12oSKO9dFEt
s39F/2h6Ld5IQrGt3gZaBB1aGO+tw3ill1VBy2zGPIDeuSz6DS3GG/oQ2gLSSMP4
OVfQ32Oajo496iHRkdIh/7Hho7BNzMYr1GxrYTcE9/Znr6xgeSdNT37CCeCH8cmP
aEAUgSMTeIMVSpILwkKeNvBURic1EWaqXRgPRIWK0vNyOCs/+jNoFISnV4pu1ROF
92vayHDNSVw9wHcdSQ75XSE4Msawqv5U1iI7e2lD64uo1qhmJdrPcXDJQCiDbh+F
hQhF+wAoLRvMNwwhg+LttL8vXqMDQl3olsWSvWPs6b/MZpB0qwd1bklzA6P+PeAU
sfOvTqi9edIOfKqvXqTXEhBP8qC7ZtOKLGnryZb7W04SSVrNtuJUFRcLiqu+w/F/
MSxGSGalYpzIZ1B5HLQqISgWMXdbt39uMeeooeZjkuI3VIllFjtybecjPR9ZYQPt
FFEP1XqNXjLFmGh84TXtvGLWretWM1OZmN8UKKUeATqrr7zuh5AYGAIbXd8BvweL
Pigl9ei0hTculPqohvkoc5x1srPBvzHrirGlxOYjW3fc4kDgZpy+6ik5k5g7JWQD
lbXCRz3HGazgUPeiwUr06a52vhgT7QuNIUZqdHb4IfCYs2pQTLHzQjAqvVk1mm2D
kh4myIcTtf69BFcu/Wuptm3NaKd1nwk1squR6psvcTXOWII81pstnxNYkrokx4r2
7YVllNruOD+cMDNZbIG2CwT6V9ukIS8tl9EJp8eyb0a1uAEc22BNOjYHPF50beWF
ukf3uc0SA+G3zhmXCM5sMf5OxVjKr5jgcir7kySY5KbmG71omYhczgr4H0qgxYo9
Zyj2wMKrTHLfFOpd4OOEun9Gi3srqlKZep7Hj7gNyUwZu1qiBvElmBVmp0HJxT0N
mktuaVbaFgBsTS0/us1EqWvCA4REh1Ut/NoA9oG3JFt0lGDstTw1j+orDmIHOmSu
7FKYzr0uCz14AkLMSOixdPD1F0YyED1NMVnRVXw77HiAFGmb0CDi2KEg70pEKpn3
ksa8oe0MQi6oEwlMsAxVTXOB1wblTBuSBeaECzTzWE+/DHF+QQfQi8kAjjSdmmMJ
yN+shdBWHYRGYnxRkTatONhcDBIY7sZV7wolYHz/rf7dpYUZf37vdQnYV8FpO1um
Ya0GslyRJ5GqMBfDS1cQKne+FvVHxEE2YqEGBcOYhx/JI2soE8aA8W4XffN+DoEy
ZkinJ/+BOwJ/zUI9GZtwB4JXqbNEE+j7r7/fJO9KxfPp4MPK4YWu0H0EUWONpVwe
TWtbRhQUCOe4PVSC/Vv1pstvMD/D+E/0L4GQNHxr+xyFxuvILty5lvFTxoAVYpqD
u8gNhk3NWefTrlSkhY4N+tPP6o7E4t3y40nOA/d9qaqiid+lYcIDB0cJTpZvgeeQ
ijohxY3PHruU4vVZa37ITQnco9az6lsy18vbU0bOyK2fEZ2R9XVO8fH11jiV8oGH
-----END RSA PRIVATE KEY-----`)
}



func masterIdentityInqueryRespond() (meta *msagent.PocketMasterAgentMeta, err error) {
    // ------------- Let's Suppose you've sent an unbounded inquery from a node over multicast net ---------------------
    ua, err := slagent.UnboundedBroadcastAgent()
    if err != nil {
        return
    }
    psm, err := slagent.PackedSlaveMeta(slagent.DiscoveryMetaAgent(ua))
    if err != nil {
        return
    }
    // -------------- over master, it's received the message and need to make an inquiry "Who R U"? --------------------
    usm, err := slagent.UnpackedSlaveMeta(psm)
    if err != nil {
        return
    }
    cmd, err := msagent.IdentityInqueryRespond(usm.DiscoveryAgent)
    if err != nil {
        return
    }
    meta = msagent.UnboundedInqueryMeta(cmd)
    return
}

func masterIdentityRespond() (meta *msagent.PocketMasterAgentMeta, err error) {
    agent, err := slagent.InquiredAgent(initSendTimestmap)
    if err != nil {
        return
    }
    msa, err := slagent.InquiredMetaAgent(agent)
    if err != nil {
        return
    }
    cmd, err := msagent.MasterIdentityRevealCommand(msa.StatusAgent, initSendTimestmap.Add(time.Second))
    if err != nil {
        return
    }
    meta = msagent.IdentityInqueryMeta(cmd, testMasterPublicKey())
    return
}

func masterKeyExchangeRespond() (meta *msagent.PocketMasterAgentMeta, err error) {
    agent, err := slagent.KeyExchangeAgent(masterBoundAgentName, initSendTimestmap)
    if err != nil {
        return
    }
    msa, err := slagent.KeyExchangeMetaAgent(agent, testSlavePublicKey())
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
    err = ioutil.WriteFile("recvtest.pub", testSlavePublicKey(), os.ModePerm)
    defer os.Remove("recvtest.pub")
    if err != nil {
        return
    }
    err = ioutil.WriteFile("sendtest.pem", testMasterPrivateKey(), os.ModePerm)
    defer os.Remove("sendtest.pem")
    if err != nil {
        return
    }
    rsaenc ,err := crypt.NewEncryptorFromKeyFiles("recvtest.pub", "sendtest.pem")
    if err != nil {
        return
    }
    // responding commnad
    cmd, slvstat, err := msagent.CryptoKeyAndNameSetCommand(usm.StatusAgent, slaveNodeName, timestmap)
    if err != nil {
        return
    }
    meta, err = msagent.ExecKeyExchangeMeta(cmd, slvstat, aeskey, aesenc, rsaenc)
    return
}

//--

func TestUnboundedState_InquiredTransition(t *testing.T) {
    meta, err := masterIdentityInqueryRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }

    ssd := NewSlaveDiscovery()
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
    meta, err := masterIdentityRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    sd := NewSlaveDiscovery()
    sd.(*slaveDiscovery).discoveryState = SlaveKeyExchange

    // execute state transition
    sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2))
    if sd.CurrentState() != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }
}

func TestKeyExchange_CryptoCheckTransition(t *testing.T) {
    meta, err := masterKeyExchangeRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    sd := NewSlaveDiscovery()
    sd.(*slaveDiscovery).discoveryState = SlaveKeyExchange

    // execute state transition
    sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2))
    if sd.CurrentState() != SlaveCryptoCheck {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", sd.CurrentState().String())
        return
    }
}