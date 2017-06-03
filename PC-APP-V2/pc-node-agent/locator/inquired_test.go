package locator

import (
    "testing"
    "time"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/msagent"
)

/* ------------------------------------------------ POSITIVE TEST --------------------------------------------------- */
func TestInquired_KeyExchangeTransition(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm *DebugCommChannel = &DebugCommChannel{}
        debugEvent *DebugEventReceiver = &DebugEventReceiver{}
    )

    meta, endTime, err := msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterWeakPublicKey(), initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm, debugComm, debugEvent)
    if err != nil {
        t.Error(err.Error())
        return
    }
    sd.(*slaveLocator).state = newInquiredState(debugComm, debugComm, debugEvent)

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, endTime.Add(time.Second)); err != nil {
        t.Errorf(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
}

/* ------------------------------------------------ NEGATIVE TEST --------------------------------------------------- */
// inquired -> keyexchange
func Test_Inquired_Keyexchange_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm *DebugCommChannel = &DebugCommChannel{}
        debugEvent *DebugEventReceiver = &DebugEventReceiver{}
        // unbounded -> inquired
        context slcontext.PocketSlaveContext = slcontext.SharedSlaveContext()
        masterTS, slaveTS time.Time = time.Now(), time.Now()
    )

    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm, debugComm, debugEvent)
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
    for i := 0; i < TransitionFailureLimit; i++ {
        // inquired -> keyexchange
        masterTS = slaveTS.Add(time.Second)
        meta, masterTS, err = msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterWeakPublicKey(), masterTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        // void master pubkey to fail transition
        meta.MasterPubkey = nil

        slaveTS = masterTS.Add(time.Second)
        err = sd.TranstionWithMasterMeta(meta, slaveTS)
        if i < (TransitionFailureLimit - 1) {
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

    var (
        debugComm *DebugCommChannel = &DebugCommChannel{}
        debugEvent *DebugEventReceiver = &DebugEventReceiver{}
        // unbounded -> inquired
        context slcontext.PocketSlaveContext = slcontext.SharedSlaveContext()
        slaveTS time.Time = time.Now()
    )

    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm, debugComm, debugEvent)
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

