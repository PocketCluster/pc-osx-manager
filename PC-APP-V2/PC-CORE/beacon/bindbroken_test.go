package beacon

import (
    "testing"
    "time"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-node-agent/slagent"
)

// ------------------------------------------------- POSITIVE TESTING -----------------------------------------------
func Test_BindBroken_BindRecovery_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        masterTS time.Time = time.Now()
    )

    slave := model.DebugTestSlaveNode()
    mb, err := NewMasterBeacon(MasterBindBroken, slave, debugComm)
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
    err = mb.TransitionWithSlaveMeta(meta, masterTS)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindRecovery {
        t.Error("[ERR] Master state is expected to be " + MasterBindRecovery.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

// ------------------------------------------------- NEGATIVE TESTING -----------------------------------------------
func Test_BindBroken_BindRecovery_TimeoutFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        masterTS time.Time = time.Now()
    )

    slave := model.DebugTestSlaveNode()
    mb, err := NewMasterBeacon(MasterBindBroken, slave, debugComm)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    meta, err := slagent.TestSlaveBindBroken("WRONG_MASTER_NAME")
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    // 1st trial
    masterTS = masterTS.Add(time.Second)
    err = mb.TransitionWithSlaveMeta(meta, masterTS)
    if err != nil {
        t.Log(err.Error())
    }
    // 2nd trial
    masterTS = masterTS.Add(time.Millisecond + UnboundedTimeout * time.Duration(TxActionLimit))
    err = mb.TransitionWithSlaveMeta(meta, masterTS)
    if err != nil {
        t.Log(err.Error())
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterBindBroken.String() + ". Current : " + mb.CurrentState().String())
    }
}

func Test_BindBroken_BindRecovery_TooManyMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        masterTS time.Time = time.Now()
    )

    slave := model.DebugTestSlaveNode()
    mb, err := NewMasterBeacon(MasterBindBroken, slave, debugComm)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    meta, err := slagent.TestSlaveBindBroken("WRONG_MASTER_NAME")
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    for i := 0; i <= int(TransitionFailureLimit); i++ {
        masterTS = masterTS.Add(time.Second)
        err = mb.TransitionWithSlaveMeta(meta, masterTS)
        if err != nil {
            t.Log(err.Error())
        }
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterBindBroken.String() + ". Current : " + mb.CurrentState().String())
    }
}

func Test_BindBroken_BindRecovery_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        masterTS time.Time = time.Now()
    )

    slave := model.DebugTestSlaveNode()
    mb, err := NewMasterBeacon(MasterBindBroken, slave, debugComm)
    if err != nil {
        t.Errorf(err.Error())
        return
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
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
