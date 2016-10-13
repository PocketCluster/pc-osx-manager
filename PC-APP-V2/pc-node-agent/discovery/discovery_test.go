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
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDCFENGw33yGihy92pDjZQhl0C3
6rPJj+CvfSC8+q28hxA161QFNUd13wuCTUcq0Qd2qsBe/2hFyc2DCJJg0h1L78+6
Z4UMR7EOcpfdUE9Hf3m/hs+FUR45uBJeDK1HSFHD8bHKD6kv8FPGfJTotc+2xjJw
oYi+1hqp1fIekaxsyQIDAQAB
-----END PUBLIC KEY-----`)
}

func testSlavePrivateKey() []byte {
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
    agent, err := slagent.AnswerMasterInquiryStatus(initSendTimestmap)
    if err != nil {
        return
    }
    msa, err := slagent.AnswerMasterInquiryMeta(agent)
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

func masterBrokenBindRecoveryCommand() (meta *msagent.PocketMasterAgentMeta, err error) {
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