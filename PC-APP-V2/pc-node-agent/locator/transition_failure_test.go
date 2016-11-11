package locator

import (
    "testing"
    "time"
    "log"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/msagent"
    "github.com/davecgh/go-spew/spew"
)

// unbounded -> inquired
func Test_Unbounded_Inquired_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    // unbounded state
    debugComm := &DebugCommChannel{}
    slaveTS := time.Now().Add(time.Second)
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveUnbounded {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    TransitionLimit := int(TransitionFailureLimit * TransitionFailureLimit)
    for i := 0 ; i < TransitionLimit; i++ {
        slaveTS = slaveTS.Add(time.Second)
        meta, err := msagent.TestMasterInquireSlaveRespond()
        if err != nil {
            t.Error(err.Error())
            return
        }
        meta.DiscoveryRespond.MasterAddress = ""
        err = sd.TranstionWithMasterMeta(meta, slaveTS)
        if err != nil {
            t.Log(err.Error())
        }
        state, err := sd.CurrentState()
        if err != nil {
            t.Error(err.Error())
            return
        }
        if state != SlaveUnbounded {
            t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
            return
        }
    }
}

func Test_Unbounded_Inquired_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    // inquired transition
    debugComm := &DebugCommChannel{}
    slaveTS := time.Now()
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveUnbounded {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    TransitionLimit := int(TransitionFailureLimit * TransitionFailureLimit)
    for i := 0 ; i < TransitionLimit; i++ {
        slaveTS = slaveTS.Add(time.Millisecond + UnboundedTimeout)
        err = sd.TranstionWithTimestamp(slaveTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        state, err := sd.CurrentState()
        if err != nil {
            t.Error(err.Error())
            return
        }
        if state != SlaveUnbounded {
            t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
            return
        }
        if len(debugComm.LastMcastMessage) == 0 || 508 < len(debugComm.LastMcastMessage) {
            t.Errorf("[ERR] Multicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastMcastMessage))
        }
    }
    if debugComm.MCommCount != uint(TransitionLimit) {
        t.Errorf("[ERR] MultiComm count does not match %d | expected %d", debugComm.MCommCount, TransitionLimit)
    }
}

// inquired -> keyexchange
func Test_Inquired_Keyexchange_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    // unbounded -> inquired
    context := slcontext.SharedSlaveContext()
    debugComm := &DebugCommChannel{}
    slaveTS := time.Now()
    masterTS := time.Now()
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    meta, err := msagent.TestMasterInquireSlaveRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = slaveTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveInquired {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    for i := 0; i < int(TransitionFailureLimit); i++ {
        // inquired -> keyexchange
        masterTS = slaveTS.Add(time.Second)
        meta, masterTS, err = msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), masterTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        // void master pubkey to fail transition
        meta.MasterPubkey = nil

        slaveTS = masterTS.Add(time.Second)
        err = sd.TranstionWithMasterMeta(meta, slaveTS)
        if i < int(TransitionFailureLimit) - 1 {
            if err != nil {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveInquired {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
        } else {
            if err == nil {
                t.Errorf("[ERR] Master meta transition count more than TransitionFailureLimit should generate error")
            } else {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveUnbounded {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
            _, err = context.GetMasterAgent()
            if err == nil {
                t.Errorf("[ERR] Master Agent Name should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterPublicKey()
            if err == nil {
                t.Errorf("[ERR] Master public key should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterIP4Address()
            if err == nil {
                t.Errorf("[ERR] Master address should be empty")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.AESCryptor()
            if err == nil {
                t.Errorf("[ERR] AESCryptor should be nil")
                return
            } else {
                t.Log(err.Error())
            }
        }
    }
}

func Test_Inquired_Keyexchange_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    // unbounded -> inquired
    context := slcontext.SharedSlaveContext()
    debugComm := &DebugCommChannel{}
    slaveTS := time.Now()
    //masterTS := time.Now()
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    meta, err := msagent.TestMasterInquireSlaveRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = slaveTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveInquired {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    for i := 0; i <= int(TxActionLimit); i++ {
        slaveTS = slaveTS.Add(time.Millisecond + UnboundedTimeout)
        err = sd.TranstionWithTimestamp(slaveTS)
        if i < int(TxActionLimit) {
            if err != nil {
                t.Error(err.Error())
                return
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveInquired {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
        } else {
            if err == nil {
                t.Error("[ERR] Tx after TxActionLimit should generate error")
                return
            } else {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveUnbounded {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
            _, err = context.GetMasterAgent()
            if err == nil {
                t.Errorf("[ERR] Master Agent Name should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterPublicKey()
            if err == nil {
                t.Errorf("[ERR] Master public key should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterIP4Address()
            if err == nil {
                t.Errorf("[ERR] Master address should be empty")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.AESCryptor()
            if err == nil {
                t.Errorf("[ERR] AESCryptor should be nil")
                return
            } else {
                t.Log(err.Error())
            }
        }
        if len(debugComm.LastUcastMessage) == 0 || 508 < len(debugComm.LastUcastMessage) {
            t.Errorf("[ERR] Multicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastUcastMessage))
        }
    }
    if debugComm.UCommCount != TxActionLimit {
        t.Errorf("[ERR] MultiComm count does not match %d | expected %d", debugComm.UCommCount, TxActionLimit)
    }
}

// keyexchange -> cryptocheck
func Test_Keyexchange_Cryptocheck_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    context := slcontext.SharedSlaveContext()
    debugComm := &DebugCommChannel{}
    slaveTS := time.Now()
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // unbounded -> inquired
    meta, err := msagent.TestMasterInquireSlaveRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = slaveTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveInquired {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
    // inquired -> keyexchange
    masterTS := slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    slaveTS = masterTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    state, err = sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    for i := 0; i <= int(TransitionFailureLimit); i++ {
        // keyexchange -> cryptocheck
        masterTS = slaveTS.Add(time.Millisecond * 100)
        meta, masterTS, err = msagent.TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName, pcrypto.TestSlavePublicKey(), pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor, masterTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        // make transition fail
        meta.RsaCryptoSignature = nil

        slaveTS = masterTS.Add(time.Millisecond * 100)
        err = sd.TranstionWithMasterMeta(meta, slaveTS)
        if i < int(TransitionFailureLimit - 1) {
            if err != nil {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveKeyExchange {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
        } else {
            if err == nil {
                t.Errorf("[ERR] Should not generate error after failure transition")
                return
            } else {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveUnbounded {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
            _, err = context.GetMasterAgent()
            if err == nil {
                t.Errorf("[ERR] Master Agent Name should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterPublicKey()
            if err == nil {
                t.Errorf("[ERR] Master public key should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterIP4Address()
            if err == nil {
                t.Errorf("[ERR] Master address should be empty")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.AESCryptor()
            if err == nil {
                t.Errorf("[ERR] AESCryptor should be nil")
                return
            } else {
                t.Log(err.Error())
            }
        }
    }
}

func Test_keyexchange_Cryptocheck_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    context := slcontext.SharedSlaveContext()
    debugComm := &DebugCommChannel{}
    slaveTS := time.Now()
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // unbounded -> inquired
    meta, err := msagent.TestMasterInquireSlaveRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = slaveTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveInquired {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
    // inquired -> keyexchange
    masterTS := slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    slaveTS = masterTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    state, err = sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    for i := 0; i <= int(TxActionLimit); i++ {
        slaveTS = slaveTS.Add(time.Millisecond + UnboundedTimeout)
        err = sd.TranstionWithTimestamp(slaveTS)
        if i < int(TxActionLimit) {
            if err != nil {
                t.Error(err.Error())
                return
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveKeyExchange {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
            if len(debugComm.LastUcastMessage) == 0 || 508 < len(debugComm.LastUcastMessage) {
                t.Errorf("[ERR] Multicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastUcastMessage))
                log.Println(spew.Sdump(debugComm.LastUcastMessage))
            }
        } else {
            if err == nil {
                t.Error("[ERR] Tx after TxActionLimit should generate error")
                return
            } else {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveUnbounded {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
            _, err = context.GetMasterAgent()
            if err == nil {
                t.Errorf("[ERR] Master Agent Name should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterPublicKey()
            if err == nil {
                t.Errorf("[ERR] Master public key should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterIP4Address()
            if err == nil {
                t.Errorf("[ERR] Master address should be empty")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.AESCryptor()
            if err == nil {
                t.Errorf("[ERR] AESCryptor should be nil")
                return
            } else {
                t.Log(err.Error())
            }
        }
        if len(debugComm.LastUcastMessage) == 0 || 508 < len(debugComm.LastUcastMessage) {
            t.Errorf("[ERR] Multicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastUcastMessage))
        }
    }
    if debugComm.UCommCount != TxActionLimit {
        t.Errorf("[ERR] MultiComm count does not match %d | expected %d", debugComm.UCommCount, TxActionLimit)
    }
}

// cryptocheck -> bounded
func Test_Cryptocheck_Bounded_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    context := slcontext.SharedSlaveContext()
    debugComm := &DebugCommChannel{}
    slaveTS := time.Now()
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // unbounded -> inquired
    meta, err := msagent.TestMasterInquireSlaveRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = slaveTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveInquired {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
    // inquired -> keyexchange
    masterTS := slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    slaveTS = masterTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    state, err = sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
    // keyexchange -> cryptocheck
    masterTS = slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName, pcrypto.TestSlavePublicKey(), pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor, masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Error(err.Error())
        return
    }
    state, err = sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveCryptoCheck {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    for i := 0; i <= int(TransitionFailureLimit); i++ {
        // cryptocheck -> bounded
        masterTS = slaveTS.Add(time.Millisecond * 100)
        meta, masterTS, err = msagent.TestMasterCheckCryptoCommand(masterAgentName, slaveNodeName, pcrypto.TestAESCryptor, masterTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        // make transition fail
        meta.MetaVersion = ""

        // FIXME : fix transition timeout window
        slaveTS = masterTS.Add(time.Millisecond * 100)
        err = sd.TranstionWithMasterMeta(meta, slaveTS)
        if i < int(TransitionFailureLimit  - 1) {
            if err != nil {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()

            if err != nil {
                t.Error(err.Error())
                return
            }

            if state != SlaveCryptoCheck {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
        } else {
            if err == nil {
                t.Errorf("[ERR] Should not generate error after failure transition")
                return
            } else {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveUnbounded {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
            _, err = context.GetMasterAgent()
            if err == nil {
                t.Errorf("[ERR] Master Agent Name should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterPublicKey()
            if err == nil {
                t.Errorf("[ERR] Master public key should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterIP4Address()
            if err == nil {
                t.Errorf("[ERR] Master address should be empty")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.AESCryptor()
            if err == nil {
                t.Errorf("[ERR] AESCryptor should be nil")
                return
            } else {
                t.Log(err.Error())
            }
        }
    }
}

func Test_Cryptocheck_Bounded_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    context := slcontext.SharedSlaveContext()
    debugComm := &DebugCommChannel{}
    slaveTS := time.Now()
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // unbounded -> inquired
    meta, err := msagent.TestMasterInquireSlaveRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = slaveTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveInquired {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
    // inquired -> keyexchange
    masterTS := slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    slaveTS = masterTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    state, err = sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
    // keyexchange -> cryptocheck
    masterTS = slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName, pcrypto.TestSlavePublicKey(), pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor, masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Error(err.Error())
        return
    }
    state, err = sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveCryptoCheck {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    for i := 0; i <= int(TxActionLimit); i++ {
        slaveTS = slaveTS.Add(time.Millisecond + UnboundedTimeout)
        err = sd.TranstionWithTimestamp(slaveTS)
        if i < int(TxActionLimit) {
            if err != nil {
                t.Error(err.Error())
                return
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveCryptoCheck {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
            if len(debugComm.LastUcastMessage) == 0 || 508 < len(debugComm.LastUcastMessage) {
                t.Errorf("[ERR] Multicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastUcastMessage))
                log.Println(spew.Sdump(debugComm.LastUcastMessage))
            }
        } else {
            if err == nil {
                t.Error("[ERR] Tx after TxActionLimit should generate error")
                return
            } else {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveUnbounded {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
            _, err = context.GetMasterAgent()
            if err == nil {
                t.Errorf("[ERR] Master Agent Name should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterPublicKey()
            if err == nil {
                t.Errorf("[ERR] Master public key should be nil")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.GetMasterIP4Address()
            if err == nil {
                t.Errorf("[ERR] Master address should be empty")
                return
            } else {
                t.Log(err.Error())
            }
            _, err = context.AESCryptor()
            if err == nil {
                t.Errorf("[ERR] AESCryptor should be nil")
                return
            } else {
                t.Log(err.Error())
            }
        }
    }
    if debugComm.UCommCount != TxActionLimit {
        t.Errorf("[ERR] MultiComm count does not match %d | expected %d", debugComm.UCommCount, TxActionLimit)
    }
}

// bounded -> bindbroken
func Test_Bounded_BindBroken_MasterMeta_Fail(t *testing.T) {
    setUp()
    defer tearDown()

    context := slcontext.SharedSlaveContext()
    debugComm := &DebugCommChannel{}
    slaveTS := time.Now()
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // unbounded -> inquired
    meta, err := msagent.TestMasterInquireSlaveRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = slaveTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveInquired {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
    // inquired -> keyexchange
    masterTS := slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    slaveTS = masterTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    state, err = sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
    // keyexchange -> cryptocheck
    masterTS = slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName, pcrypto.TestSlavePublicKey(), pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor, masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Error(err.Error())
        return
    }
    state, err = sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveCryptoCheck {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
    // cryptocheck -> bounded
    masterTS = slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterCheckCryptoCommand(masterAgentName, slaveNodeName, pcrypto.TestAESCryptor, masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Error(err.Error())
        return
    }
    state, err = sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveBounded {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    for i := 0; i <= int(TransitionFailureLimit); i++ {
        // cryptocheck -> bounded
        masterTS = slaveTS.Add(time.Second)
        meta, masterTS, err = msagent.TestMasterBoundedStatusCommand(masterAgentName, slaveNodeName, pcrypto.TestAESCryptor, masterTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        // make transition fail
        meta.MetaVersion = ""

        slaveTS = masterTS.Add(time.Second)
        err = sd.TranstionWithMasterMeta(meta, slaveTS)
        if i < int(TransitionFailureLimit - 1) {
            if err != nil {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveBounded {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
        } else {
            if err == nil {
                t.Errorf("[ERR] Exceeding # transition trial should cause an error")
            } else {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveBindBroken {
                t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
                return
            }
            aeskey := context.GetAESKey()
            if len(aeskey) != 0 {
                t.Errorf("[ERR] AES KEY should be nil")
                return
            }
            _, err = context.AESCryptor()
            if err == nil {
                t.Errorf("[ERR] AESCryptor should be nil")
                return
            } else {
                t.Log(err.Error())
            }
        }
    }
}

func Test_Bounded_BindBroken_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    // Let's have a bounded state
    debugComm := &DebugCommChannel{}
    context := slcontext.SharedSlaveContext()
    context.SetMasterPublicKey(pcrypto.TestMasterPublicKey())
    context.SetMasterAgent(masterAgentName)
    context.SetSlaveNodeName(slaveNodeName)

    // have a slave locator
    sd, err := NewSlaveLocator(SlaveBindBroken, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS := time.Now()
    meta, err := msagent.TestMasterBrokenBindRecoveryCommand(masterAgentName, pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor)
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS := masterTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS);
    if err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveBounded {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    //FIXME : get the exact count. We are now running 7
    slaveTS = time.Now()
    var i uint = 0
    for ; i <= TxActionLimit + 1; i++ {
        slaveTS = slaveTS.Add(time.Second + BoundedTimeout)
        err = sd.TranstionWithTimestamp(slaveTS)
        if i < TxActionLimit {
            if err != nil {
                t.Error(err.Error())
                return
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveBounded {
                t.Errorf("[ERR] Slave state should not change properly | Current : %s\n", state.String())
                return
            }
            if len(debugComm.LastUcastMessage) == 0 || 508 < len(debugComm.LastUcastMessage) {
                t.Errorf("[ERR] Unicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastUcastMessage))
            }
        } else {
            if err != nil {
                t.Log(err.Error())
            }
            state, err = sd.CurrentState()
            if err != nil {
                t.Error(err.Error())
                return
            }
            if state != SlaveBindBroken {
                t.Errorf("[ERR] Slave state should not change properly | Current : %s\n", state.String())
                return
            }
            // FIXME : HOW?
            if len(debugComm.LastMcastMessage) == 0 || 508 < len(debugComm.LastMcastMessage) {
                //t.Errorf("[ERR] Multicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastMcastMessage))
            }
            aesKey := context.GetAESKey()
            if len(aesKey) != 0 {
                t.Error("[ERR] AES key should be nil")
                return
            }
        }
    }
    if debugComm.UCommCount != TxActionLimit {
        t.Errorf("[ERR] comm count does not match %d | expected %d", debugComm.UCommCount, TxActionLimit)
    }
    if len(debugComm.LastUcastMessage) == 0 || 508 < len(debugComm.LastUcastMessage) {
        t.Errorf("[ERR] Unicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastUcastMessage))
    }
    if debugComm.MCommCount != 1 {
        t.Errorf("[ERR] MultiComm count does not match %d / expected %d ", debugComm.MCommCount, 1)
    }
    // FIXME : WHY?
    if len(debugComm.LastMcastMessage) == 0 || 508 < len(debugComm.LastMcastMessage) {
        t.Errorf("[ERR] Multicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastMcastMessage))
    }
}

// bindbroken -> bindbroken
func Test_BindBroken_BindBroken_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    // by the time bind broken state is revived, previous master public key should have been available.
    debugComm := &DebugCommChannel{}
    context := slcontext.SharedSlaveContext()
    context.SetMasterPublicKey(pcrypto.TestMasterPublicKey())
    context.SetMasterAgent(masterAgentName)
    context.SetSlaveNodeName(slaveNodeName)

    sd, err := NewSlaveLocator(SlaveBindBroken, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    slaveTS := time.Now()
    TxCountTarget := TxActionLimit * TxActionLimit
    var i uint = 0
    for ;i < TxCountTarget; i++ {
        slaveTS = slaveTS.Add(time.Millisecond + UnboundedTimeout)
        err = sd.TranstionWithTimestamp(slaveTS);
        if err != nil {
            t.Error(err.Error())
            return
        }
        state, err := sd.CurrentState()
        if err != nil {
            t.Error(err.Error())
            return
        }
        // now broken bind is recovered
        if state != SlaveBindBroken {
            t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
            return
        }
        if len(debugComm.LastMcastMessage) == 0 || 508 < len(debugComm.LastMcastMessage) {
            t.Errorf("[ERR] Multicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastMcastMessage))
        }
    }
    if debugComm.MCommCount != TxCountTarget {
        t.Errorf("[ERR] MultiComm count does not match %d / expected %d ", debugComm.MCommCount, TxCountTarget)
    }
}

// unbounded -> unbounded
func Test_Unbounded_Unbounded_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    debugComm := &DebugCommChannel{}
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm)
    if err != nil {
        t.Error(err.Error())
        return
    }

    /* ---------------------------------------------- make transition failed ---------------------------------------- */
    slaveTS := time.Now()
    TxCountTarget := TxActionLimit * TxActionLimit
    var i uint = 0
    for ;i < TxCountTarget; i++ {
        slaveTS = slaveTS.Add(time.Millisecond + UnboundedTimeout)
        err = sd.TranstionWithTimestamp(slaveTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        state, err := sd.CurrentState()
        if err != nil {
            t.Error(err.Error())
            return
        }
        if state != SlaveUnbounded {
            t.Errorf("[ERR] Slave state should not change properly | Current : %s\n", state.String())
            return
        }
        if len(debugComm.LastMcastMessage) == 0 || 508 < len(debugComm.LastMcastMessage) {
            t.Errorf("[ERR] Multicast message cannot exceed 508 bytes. Current %d", len(debugComm.LastMcastMessage))
        }
    }
    if debugComm.MCommCount != TxCountTarget {
        t.Errorf("[ERR] MultiComm count does not match %d | expected %d", debugComm.MCommCount, TxCountTarget)
    }
}