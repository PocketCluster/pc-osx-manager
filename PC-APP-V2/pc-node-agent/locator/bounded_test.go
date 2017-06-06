package locator

import (
    "testing"
    "time"
    "bytes"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/msagent"
)

/* ------------------------------------------------ POSITIVE TEST --------------------------------------------------- */
func Test_Unbounded_Bounded_Onepass(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm *DebugCommChannel = &DebugCommChannel{}
        debugEvent *DebugEventReceiver = &DebugEventReceiver{}
        context slcontext.PocketSlaveContext = slcontext.SharedSlaveContext()
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
    err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2))
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
    meta, masterTS, err := msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterWeakPublicKey(), initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    slaveTS := masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
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

    // cryptocheck -> bounded
    masterTS = slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterCheckCryptoCommand(masterAgentName, authToken, slaveNodeName, pcrypto.TestAESCryptor, masterTS)
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
    // Verification
    if msName, _ := context.GetMasterAgent(); msName != masterAgentName {
        t.Errorf("[ERR] master node name is setup inappropriately | Current : %s\n", msName)
        return
    }
    if snName, _ := context.GetSlaveNodeName(); snName != slaveNodeName {
        t.Errorf("[ERR] slave node name is setup inappropriately | Current : %s\n", snName)
        return
    }
    if bytes.Compare(context.GetAESKey(), pcrypto.TestAESKey) != 0 {
        t.Errorf("[ERR] slave aes key is setup inappropriately")
        return
    }
}

func Test_Bounded_Unbroken_Loop(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm *DebugCommChannel = &DebugCommChannel{}
        debugEvent *DebugEventReceiver = &DebugEventReceiver{}
        context slcontext.PocketSlaveContext = slcontext.SharedSlaveContext()
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
    err = sd.TranstionWithMasterMeta(meta, initSendTimestmap.Add(time.Second * 2))
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
    meta, masterTS, err := msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterWeakPublicKey(), initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // execute state transition
    slaveTS := masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
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
    // MASTER META TRANSITION ACTION
    masterTS = slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterCheckCryptoCommand(masterAgentName, slaveNodeName, authToken, pcrypto.TestAESCryptor, masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    err = sd.TranstionWithMasterMeta(meta, slaveTS)
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
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }

    // bounded loop
    for i := 0 ; i < 100; i++ {
        // MASTER META TRANSITION ACTION
        masterTS = slaveTS.Add(time.Second)
        meta, masterTS, err = msagent.TestMasterBoundedStatusCommand(masterAgentName, slaveNodeName, authToken, pcrypto.TestAESCryptor, masterTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        slaveTS = masterTS.Add(time.Second)
        err = sd.TranstionWithMasterMeta(meta, slaveTS)
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
            t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
            return
        }
        // TRANSMISSION ACTION
        slaveTS = slaveTS.Add(time.Second)
        err = sd.TranstionWithTimestamp(slaveTS)
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
            t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
            return
        }
        // Verification
        if msName, _ := context.GetMasterAgent(); msName != masterAgentName {
            t.Errorf("[ERR] master node name is setup inappropriately | Current : %s\n", msName)
            return
        }
        if snName, _ := context.GetSlaveNodeName(); snName != slaveNodeName {
            t.Errorf("[ERR] slave node name is setup inappropriately | Current : %s\n", snName)
            return
        }
        if bytes.Compare(context.GetAESKey(), pcrypto.TestAESKey) != 0 {
            t.Errorf("[ERR] slave aes key is setup inappropriately")
            return
        }
    }
}

/* ------------------------------------------------ NEGATIVE TEST --------------------------------------------------- */
// bounded -> bindbroken
func Test_Bounded_BindBroken_MasterMeta_Fail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm *DebugCommChannel = &DebugCommChannel{}
        debugEvent *DebugEventReceiver = &DebugEventReceiver{}
        context slcontext.PocketSlaveContext = slcontext.SharedSlaveContext()
        slaveTS = time.Now()
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
    // cryptocheck -> bounded
    masterTS = slaveTS.Add(time.Second)
    meta, masterTS, err = msagent.TestMasterCheckCryptoCommand(masterAgentName, slaveNodeName, authToken, pcrypto.TestAESCryptor, masterTS)
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
    for i := 0; i <= TransitionFailureLimit; i++ {
        // cryptocheck -> bounded
        masterTS = slaveTS.Add(time.Second)
        meta, masterTS, err = msagent.TestMasterBoundedStatusCommand(masterAgentName, slaveNodeName, authToken, pcrypto.TestAESCryptor, masterTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        // make transition fail
        meta.MetaVersion = ""

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

    var (
        debugComm *DebugCommChannel = &DebugCommChannel{}
        debugEvent *DebugEventReceiver = &DebugEventReceiver{}
        context slcontext.PocketSlaveContext = slcontext.SharedSlaveContext()
        masterTS, slaveTS = time.Now(), time.Now()
    )

    context.SetMasterPublicKey(pcrypto.TestMasterWeakPublicKey())
    context.SetMasterAgent(masterAgentName)
    context.SetSlaveNodeName(slaveNodeName)
    context.SetSlaveAuthToken(authToken)

    // have a slave locator
    sd, err := NewSlaveLocator(SlaveBindBroken, debugComm, debugComm, debugEvent)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = time.Now()
    meta, err := msagent.TestMasterBrokenBindRecoveryCommand(masterAgentName, pcrypto.TestAESKey, pcrypto.TestAESCryptor, pcrypto.TestMasterRSAEncryptor)
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
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
    for i := 0; i <= TxActionLimit + 1; i++ {
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

