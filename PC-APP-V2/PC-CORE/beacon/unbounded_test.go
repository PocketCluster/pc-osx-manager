package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
)

func Test_Unbounded_Inquired_Transition_TimeoutFail(t *testing.T) {
    setUp()
    defer tearDown()

    // --- VARIABLE PREP ---
    debugComm := &DebugCommChannel{}
    masterTS := time.Now()
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
    masterTS = time.Now()
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- test ---
    slaveTS := masterTS.Add(time.Second)
    sa, _, err = slagent.TestSlaveAnswerMasterInquiry(slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    // this is an error injection
    sa.StatusAgent.MasterBoundAgent = "MASTER-YODA"
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err == nil {
        t.Errorf("[ERR] incorrect slave state. Should generate error with wrong master name")
        return
    } else {
        t.Logf(err.Error())
    }
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != 1 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
    // update with timestamp
    slaveTS = masterTS.Add(time.Second * 11)
    t.Logf("[INFO] slaveTS - MasterBeacon.lastSuccessTimestmap : " + slaveTS.Sub(mb.(*masterBeacon).state.(DebugState).TransitionSuccessTS()).String())
    if err := mb.TransitionWithTimestamp(slaveTS); err != nil {
        t.Log(err.Error())
    }
    if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != 1 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).state.(DebugState).TransitionFailed())
        return
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_Unbounded_Inquired_Transition_TooManyMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    // --- VARIABLE PREP ---
    debugComm := &DebugCommChannel{}
    masterTS := time.Now()
    slaveTS := time.Now()
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
    masterTS = time.Now()
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- test
    for i := 0; i < int(TransitionFailureLimit); i ++ {
        slaveTS = masterTS.Add(time.Second)
        sa, end, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        // this is an error injection
        sa.StatusAgent.MasterBoundAgent = "MASTER-YODA"
        masterTS = end.Add(time.Second)
        if err := mb.TransitionWithSlaveMeta(sa, masterTS); err == nil {
            t.Errorf("[ERR] incorrect slave state. Should generate error with wrong master name")
            return
        } else {
            t.Logf(err.Error())
        }
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != 5 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
}