package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pcrypto"
)

func Test_Inquired_KeyExchange_TimeoutFail(t *testing.T) {
    setUp()
    defer tearDown()

    debugComm := &DebugCommChannel{}

    // test var preperations
    mb, err := NewMasterBeacon(MasterInit, nil, debugComm)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    sa, err := slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS := time.Now()
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS := masterTS.Add(time.Second)
    sa, end, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }

    // --- TEST
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveKeyExchangeStatus("MASTER-YODA", pcrypto.TestSlavePublicKey(), slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // first trial
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err == nil {
        t.Error("[ERR] Master beacon should have failed if wrong master name is fed!")
        return
    } else {
        t.Log(err.Error())
    }
    if mb.CurrentState() != MasterInquired {
        t.Error("[ERR] Master state is expected to be " + MasterInquired.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != 1 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }

    // 2nd try
    masterTS = masterTS.Add(time.Millisecond + UnboundedTimeout * time.Duration(TxActionLimit))
    err = mb.TransitionWithSlaveMeta(sa, masterTS)
    if err != nil {
        t.Log(err.Error())
    }
    if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != 0 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).state.(DebugState).TransitionFailed())
        return
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_Inquired_KeyExchange_TooManyMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    debugComm := &DebugCommChannel{}

    // test var preperations
    mb, err := NewMasterBeacon(MasterInit, nil, debugComm)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    sa, err := slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS := time.Now()
    err = mb.TransitionWithSlaveMeta(sa, masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS := masterTS.Add(time.Second)
    sa, end, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }

    // --- test ---
    for i := 1; i <= int(TransitionFailureLimit); i++ {
        if mb.CurrentState() != MasterInquired {
            t.Errorf("[ERR] Master state is expected to be %s. Current : %s. Trial count %d", MasterInquired.String(), mb.CurrentState().String(), mb.(*masterBeacon).state.(DebugState).TransitionFailed())
            return
        }
        slaveTS = masterTS.Add(time.Second)
        t.Logf("[INFO] slaveTS - MasterBeacon.lastSuccessTimestmap : " + slaveTS.Sub(mb.(*masterBeacon).state.(DebugState).TransitionSuccessTS()).String())
        sa, end, err = slagent.TestSlaveKeyExchangeStatus("MASTER-YODA", pcrypto.TestSlavePublicKey(), slaveTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        masterTS = end.Add(time.Second)
        if err := mb.TransitionWithSlaveMeta(sa, masterTS); err == nil {
            t.Error("[ERR] Master beacon should have failed if wrong master name is fed!")
            return
        } else {
            t.Log(err.Error())
        }
        if i < int(TransitionFailureLimit) {
            if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != uint(i) {
                t.Errorf("[ERR] Master fail count [%d] should have increased appropriatetly : iteration %d", mb.(*masterBeacon).state.(DebugState).TransitionFailed(), i)
                return
            }
        }
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != 0 {
        t.Errorf("[ERR] Wrong master discarded count %d", mb.(*masterBeacon).state.(DebugState).TransitionFailed())
        return
    }
}