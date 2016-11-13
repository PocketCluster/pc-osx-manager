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
func Test_BindBroken_Bounded_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    // by the time bind broken state is revived, previous master public key should have been available.
    context := slcontext.SharedSlaveContext()
    context.SetMasterPublicKey(pcrypto.TestMasterPublicKey())
    context.SetMasterAgent(masterAgentName)
    context.SetSlaveNodeName(slaveNodeName)

    sd, err := NewSlaveLocator(SlaveBindBroken, &DebugCommChannel{})
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
    // execute state transition
    slaveTS := masterTS.Add(time.Second)
    if err = sd.TranstionWithMasterMeta(meta, slaveTS); err != nil {
        t.Error(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    // now broken bind is recovered
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


/* ------------------------------------------------ NEGATIVE TEST --------------------------------------------------- */
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
