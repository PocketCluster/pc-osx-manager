package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
)

func Test_Init_Unbounded_Transition_TimeoutFail(t *testing.T) {
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

    // --- TIMEOUT FAILURE ---
    sa, err := slagent.TestSlaveBindBroken(masterAgentName)
    if err != nil {
        t.Skip(err.Error())
    }

    // 1st fail
    masterTS := time.Now()
    err = mb.TransitionWithSlaveMeta(sa, masterTS)
    if err == nil {
        t.Errorf("[ERR] incorrect slave state should generate error when fed to freshly spwaned beacon")
        return
    } else {
        t.Logf(err.Error())
    }
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != 1 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }

    // 2nd fail
    masterTS = masterTS.Add(time.Millisecond + UnboundedTimeout * time.Duration(TxActionLimit))
    err = mb.TransitionWithSlaveMeta(sa, masterTS)
    if err != nil {
        t.Log(err.Error())
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != 0 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).state.(DebugState).TransitionFailed())
        return
    }
}

func Test_Init_Unbounded_Transition_TooManyMetaFail(t *testing.T) {
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

    // --- slave lookup master
    sa, err := slagent.TestSlaveBindBroken(masterAgentName)
    if err != nil {
        t.Skip(err.Error())
    }
    masterTS := time.Now()
    // four more times of failure with incorrect slave meta
    for i := 0; i < int(TransitionFailureLimit); i++ {
        masterTS = masterTS.Add(time.Second)
        err = mb.TransitionWithSlaveMeta(sa, masterTS)
        if err == nil {
            t.Errorf("[ERR] incorrect slave state should generate error when fed to freshly spwaned beacon")
            return
        } else {
            t.Log(err.Error())
        }
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != 0 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
}

func Test_BeaconInit_Unbounded_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        masterTS time.Time = time.Now()
    )

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

    // --- TX ACTION FAIL ---
    for i := 0; i <= int(TxActionLimit); i++ {
        masterTS = masterTS.Add(time.Millisecond + UnboundedTimeout)
        err = mb.TransitionWithTimestamp(masterTS)
        if err != nil {
            t.Log(err.Error())
        }
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if len(debugComm.(*DebugCommChannel).LastUcastMessage) != 0 {
        t.Error("[ERR] CommChannel Ucast Message should not contain any messages")
        return
    }
    if mb.SlaveNode().IP4Address != debugComm.(*DebugCommChannel).LastUcastHost {
        t.Error("[ERR] CommChannel Ucast Message should match slave node address")
        return
    }
    if debugComm.(*DebugCommChannel).UCommCount != 0 {
        t.Errorf("[ERR] MultiComm count does not match %d | expected %d", debugComm.(*DebugCommChannel).UCommCount, 0)
    }
}
