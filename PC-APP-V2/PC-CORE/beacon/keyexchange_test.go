package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/model"
)

func Test_KeyExchange_CryptoCheck_TimeoutFail(t *testing.T) {
    setUp()
    defer tearDown()

    // --- TIMEOUT FAILURE ---
    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
    )

    // test var preperations
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
    masterTS := time.Now()
    mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    slaveTS := masterTS.Add(time.Second)
    sa, end, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveKeyExchangeStatus(masterAgentName, pcrypto.TestSlavePublicKey(), slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    // --- test
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, "INCORRECT-SLAVE-NAME", mb.SlaveNode().SlaveUUID, mb.(*masterBeacon).state.(DebugState).AESCryptor(), slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    err = mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
    if err == nil {
        t.Errorf("[ERR] Incorrect slave name should fail master beacon to transition")
        return
    } else {
        t.Log(err.Error())
    }
    if mb.CurrentState() != MasterKeyExchange {
        t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).state.(DebugState).TransitionFailed() != 1 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
    // 2nd trial
    masterTS = masterTS.Add(time.Millisecond + UnboundedTimeout * time.Duration(TxActionLimit))
    t.Logf("[INFO] masterTS - MasterBeacon.lastSuccessTimestmap : " + masterTS.Sub(mb.(*masterBeacon).state.(DebugState).TransitionSuccessTS()).String())
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

func Test_KeyExchange_CryptoCheck_TooManyMetaFail(t *testing.T) {
    // --- TOO MANY TIMES FAILURE ---
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
    )

    // test var preperations
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
    masterTS := time.Now()
    mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    slaveTS := masterTS.Add(time.Second)
    sa, end, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveKeyExchangeStatus(masterAgentName, pcrypto.TestSlavePublicKey(), slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }

    // --- test
    for i := 0; i < int(TransitionFailureLimit); i++ {
        if mb.CurrentState() != MasterKeyExchange {
            t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
            return
        }
        slaveTS = masterTS.Add(time.Second)
        t.Logf("[INFO] slaveTS - MasterBeacon.lastSuccessTimestmap : " + slaveTS.Sub(mb.(*masterBeacon).state.(DebugState).TransitionSuccessTS()).String())
        sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, "INCORRECT-SLAVE-NAME", mb.SlaveNode().SlaveUUID, mb.(*masterBeacon).state.(DebugState).AESCryptor(), slaveTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        masterTS = end.Add(time.Second)
        if err := mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS); err == nil {
            t.Errorf("[ERR] Incorrect slave name should fail master beacon to transition")
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

func Test_KeyExchange_CryptoCheck_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        debugEvent BeaconOnTransitionEvent = &DebugTransitionEventReceiver{}
        masterTS, slaveTS time.Time = time.Now(), time.Now()
    )

    // test var preperations
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
    mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS)
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    sa, end, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterInquired {
        t.Errorf("[ERR] Master state is expected to be %s. Current : %s. Trial count %d", MasterInquired.String(), mb.CurrentState().String(), mb.(*masterBeacon).state.(DebugState).TransitionFailed())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveKeyExchangeStatus(masterAgentName, pcrypto.TestSlavePublicKey(), slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(slaveAddr, sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterKeyExchange {
        t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
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
