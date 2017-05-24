package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/model"
)

func Test_Unbounded_Inquired_Transition_TimeoutFail(t *testing.T) {
    setUp()
    defer tearDown()

    // --- VARIABLE PREP ---
    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
        masterTS = time.Now()
    )

    mb, err := NewMasterBeacon(MasterInit, model.NewSlaveNode(slaveSanitizer), debugComm, debugEvent)
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
    if err := mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS); err != nil {
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
    sa.StatusAgent.Version = ""
    err = mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
    if err == nil {
        t.Errorf("[ERR] incorrect slave status version should generate error")
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
    masterTS = masterTS.Add(time.Millisecond + UnboundedTimeout * time.Duration(TxActionLimit))
    t.Logf("[INFO] slaveTS - MasterBeacon.lastSuccessTimestmap : " + masterTS.Sub(mb.(*masterBeacon).state.(DebugState).TransitionSuccessTS()).String())
    err = mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
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

func Test_Unbounded_Inquired_Transition_TooManyMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    // --- VARIABLE PREP ---
    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
        masterTS, slaveTS time.Time = time.Now(), time.Now()
    )

    mb, err := NewMasterBeacon(MasterInit, model.NewSlaveNode(slaveSanitizer), debugComm, debugEvent)
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
    err = mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
    if err != nil {
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
        sa.StatusAgent.Version = ""
        masterTS = end.Add(time.Second)
        err = mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
        if err == nil {
            t.Errorf("[ERR] incorrect slave status version should generate error")
            return
        } else {
            t.Logf(err.Error())
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

func Test_Unbounded_Inquired_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
        masterTS time.Time = time.Now()
    )

    mb, err := NewMasterBeacon(MasterInit, model.NewSlaveNode(slaveSanitizer), debugComm, debugEvent)
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
    err = mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
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
    if len(debugComm.(*DebugCommChannel).LastUcastMessage) == 0 {
        t.Error("[ERR] CommChannel Ucast Message should contain proper messages")
        return
    }
    addr, err := mb.SlaveNode().IP4AddrString()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if addr != debugComm.(*DebugCommChannel).LastUcastHost {
        t.Error("[ERR] CommChannel Ucast Message should match slave node address")
        return
    }
    if debugComm.(*DebugCommChannel).UCommCount != TxActionLimit {
        t.Errorf("[ERR] MultiComm count does not match %d | expected %d", debugComm.(*DebugCommChannel).UCommCount, TxActionLimit)
    }
}
