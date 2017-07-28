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

/* ------------------------------------------------ POSITIVE TEST --------------------------------------------------- */
// cryptocheck -> bounded
func Test_Cryptocheck_Bounded_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm *DebugCommChannel = &DebugCommChannel{}
        debugEvent *DebugEventReceiver = &DebugEventReceiver{}
        context slcontext.PocketSlaveContext = slcontext.SharedSlaveContext()
        slaveTS time.Time = time.Now()
    )

    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm, debugComm, debugEvent)
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
    meta, masterTS, err = msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterWeakPublicKey(), masterTS)
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
    meta, masterTS, err = msagent.TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName, authToken, pcrypto.TestSlavePublicKey(), pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor, masterTS)
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
    for i := 0; i <= TransitionFailureLimit; i++ {
        // cryptocheck -> bounded
        masterTS = slaveTS.Add(time.Millisecond * 100)
        meta, masterTS, err = msagent.TestMasterCheckCryptoCommand(masterAgentName, slaveNodeName, authToken, pcrypto.TestAESCryptor, masterTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        // make transition fail
        meta.MetaVersion = ""

        // FIXME : fix transition timeout window
        slaveTS = masterTS.Add(time.Millisecond * 100)
        err = sd.TranstionWithMasterMeta(meta, slaveTS)
        if i < (TransitionFailureLimit  - 1) {
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
            _, err = context.GetClusterID()
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

/* ------------------------------------------------ NEGATIVE TEST --------------------------------------------------- */
func Test_Cryptocheck_Bounded_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm *DebugCommChannel = &DebugCommChannel{}
        debugEvent *DebugEventReceiver = &DebugEventReceiver{}
        context slcontext.PocketSlaveContext = slcontext.SharedSlaveContext()
        slaveTS time.Time = time.Now()
    )

    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm, debugComm, debugEvent)
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
    meta, masterTS, err = msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterWeakPublicKey(), masterTS)
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
    meta, masterTS, err = msagent.TestMasterKeyExchangeCommand(masterAgentName, slaveNodeName, authToken, pcrypto.TestSlavePublicKey(), pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor, masterTS)
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
    for i := 0; i <= TxActionLimit; i++ {
        slaveTS = slaveTS.Add(time.Millisecond + UnboundedTimeout)
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
            _, err = context.GetClusterID()
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

