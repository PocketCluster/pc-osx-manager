package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/model"
)

func Test_BindRecovery_Bounded_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
        masterTS, slaveTS time.Time = time.Now(), time.Now()
    )

    // slave model has to be exclusively in 'joined' state
    slave := model.DebugTestSlaveNode()
    slave.State = model.SNMStateJoined
    mb, err := NewMasterBeacon(MasterBindBroken, slave, debugComm, debugEvent)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    meta, err := slagent.TestSlaveBindBroken(masterAgentName)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    masterTS = time.Now()
    err = mb.TransitionWithSlaveMeta(slaveAddr, meta, masterTS)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindRecovery {
        t.Error("[ERR] Master state is expected to be " + MasterBindRecovery.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    aescryptor := mb.(*masterBeacon).state.(DebugState).AESCryptor()
    sa, end, err := slagent.TestSlaveBoundedStatus(masterAgentName, slaveNodeName, mb.SlaveNode().SlaveUUID, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterBounded {
        t.Error("[ERR] Master state is expected to be " + MasterBounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

// ------------------------------------------------- NEGATIVE TESTING -----------------------------------------------
func Test_BindRecovery_Bounded_TimeoutFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
        masterTS, slaveTS time.Time = time.Now(), time.Now()
    )

    slave := model.DebugTestSlaveNode()
    mb, err := NewMasterBeacon(MasterBindBroken, slave, debugComm, debugEvent)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    meta, err := slagent.TestSlaveBindBroken(masterAgentName)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    masterTS = time.Now()
    err = mb.TransitionWithSlaveMeta(slaveAddr, meta, masterTS)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindRecovery {
        t.Error("[ERR] Master state is expected to be " + MasterBindRecovery.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    aescryptor := mb.(*masterBeacon).state.(DebugState).AESCryptor()

    // first trial with error
    slaveTS = masterTS.Add(time.Second)
    sa, end, err := slagent.TestSlaveBoundedStatus("WRONG_MASTER_NAME", slaveNodeName, slave.SlaveUUID, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS); err != nil {
        t.Log(err.Error())
    }
    if mb.CurrentState() != MasterBindRecovery {
        t.Error("[ERR] Master state is expected to be " + MasterBounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // 2nd trial
    slaveTS = masterTS.Add(time.Millisecond + BoundedTimeout * time.Duration(TxActionLimit))
    sa, end, err = slagent.TestSlaveBoundedStatus("WRONG_MASTER_NAME", slaveNodeName, slave.SlaveUUID, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    err = mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
    if err != nil {
        t.Log(err.Error())
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterBindBroken.String() + ". Current : " + mb.CurrentState().String())
    }
}

func Test_BindRecovery_Bounded_TooManyMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
        masterTS, slaveTS time.Time = time.Now(), time.Now()
    )

    slave := model.DebugTestSlaveNode()
    mb, err := NewMasterBeacon(MasterBindBroken, slave, debugComm, debugEvent)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    meta, err := slagent.TestSlaveBindBroken(masterAgentName)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    masterTS = time.Now()
    err = mb.TransitionWithSlaveMeta(slaveAddr, meta, masterTS)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindRecovery {
        t.Error("[ERR] Master state is expected to be " + MasterBindRecovery.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    aescryptor := mb.(*masterBeacon).state.(DebugState).AESCryptor()

    for i := 0; i <= int(TransitionFailureLimit); i++ {
        // first trial with error
        slaveTS = masterTS.Add(time.Second)
        sa, end, err := slagent.TestSlaveBoundedStatus("WRONG_MASTER_NAME", slaveNodeName, slave.SlaveUUID, aescryptor, slaveTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        masterTS = end.Add(time.Second)
        err = mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
        if err != nil {
            t.Log(err.Error())
        }
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterBindBroken.String() + ". Current : " + mb.CurrentState().String())
    }
}

func Test_BindRecovery_Bounded_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
        masterTS time.Time = time.Now()
    )

    slave := model.DebugTestSlaveNode()
    mb, err := NewMasterBeacon(MasterBindBroken, slave, debugComm, debugEvent)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    meta, err := slagent.TestSlaveBindBroken(masterAgentName)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    masterTS = time.Now()
    err = mb.TransitionWithSlaveMeta(slaveAddr, meta, masterTS)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindRecovery {
        t.Error("[ERR] Master state is expected to be " + MasterBindRecovery.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    for i := 0; i <= int(TxActionLimit); i++ {
        masterTS = masterTS.Add(time.Millisecond + UnboundedTimeout)
        err = mb.TransitionWithTimestamp(masterTS)
        if err != nil {
            t.Log(err.Error())
        }
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterBindBroken.String() + ". Current : " + mb.CurrentState().String())
    }
}
