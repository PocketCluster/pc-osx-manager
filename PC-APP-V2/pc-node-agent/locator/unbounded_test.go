package locator

import (
    "testing"
    "time"

    "github.com/stkim1/pc-core/msagent"
)

/* ------------------------------------------------ POSITIVE TEST --------------------------------------------------- */
func TestUnboundedState_InquiredTransition(t *testing.T) {
    setUp()
    defer tearDown()

    meta, err := msagent.TestMasterInquireSlaveRespond()
    if err != nil {
        t.Error(err.Error())
        return
    }

    debugComm := &DebugCommChannel{}
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm, debugComm)
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
}

/* ------------------------------------------------ NEGATIVE TEST --------------------------------------------------- */
func Test_Unbounded_Unbounded_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    debugComm := &DebugCommChannel{}
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm, debugComm)
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

func Test_Unbounded_Inquired_MasterMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    // unbounded state
    debugComm := &DebugCommChannel{}
    slaveTS := time.Now().Add(time.Second)
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm, debugComm)
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
    sd, err := NewSlaveLocator(SlaveUnbounded, debugComm, debugComm)
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
