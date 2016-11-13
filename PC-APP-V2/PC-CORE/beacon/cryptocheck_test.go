package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pcrypto"
)

func Test_CryptoCheck_Bounded_TimeoutFail(t *testing.T) {
    setUp()
    defer tearDown()

    // --- TIMEOUT FAILURE ---
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
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveKeyExchangeStatus(masterAgentName, pcrypto.TestSlavePublicKey(), slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, mb.SlaveNode().NodeName, mb.(*masterBeacon).state.(*beaconState).aesCryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }

    // --- test
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveBoundedStatus(masterAgentName, "INCORRECT-SLAVE-NAME", mb.(*masterBeacon).state.(*beaconState).aesCryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err == nil {
        t.Errorf("[ERR] Incorrect slave name should fail master beacon to transition")
        return
    } else {
        t.Log(err.Error())
    }
    if mb.CurrentState() != MasterCryptoCheck {
        t.Error("[ERR] Master state is expected to be " + MasterCryptoCheck.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).state.(*beaconState).transitionFailureCount != 1 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
    // fail with timestamp
    masterTS = masterTS.Add(time.Second * 11)
    t.Logf("[INFO] masterTS - MasterBeacon.lastSuccessTimestmap : " + mb.(*masterBeacon).state.(*beaconState).lastTransitionTS.String())
    if err := mb.TransitionWithTimestamp(masterTS); err != nil {
        t.Log(err.Error())
    }
    if mb.(*masterBeacon).state.(*beaconState).transitionFailureCount != 1 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).state.(*beaconState).transitionFailureCount)
        return
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_CryptoCheck_Bounded_TooManyMetaFail(t *testing.T) {
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
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveKeyExchangeStatus(masterAgentName, pcrypto.TestSlavePublicKey(), slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, mb.SlaveNode().NodeName, mb.(*masterBeacon).state.(*beaconState).aesCryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    // --- test
    for i := 0; i < int(TransitionFailureLimit); i++ {
        if mb.CurrentState() != MasterCryptoCheck {
            t.Error("[ERR] Master state is expected to be " + MasterCryptoCheck.String() + ". Current : " + mb.CurrentState().String())
            return
        }
        slaveTS = masterTS.Add(time.Second)
        sa, end, err = slagent.TestSlaveBoundedStatus(masterAgentName, "INCORRECT-SLAVE-NAME", mb.(*masterBeacon).state.(*beaconState).aesCryptor, slaveTS)
        if err != nil {
            t.Error(err.Error())
            return
        }
        masterTS = end.Add(time.Second)
        if err := mb.TransitionWithSlaveMeta(sa, masterTS); err == nil {
            t.Errorf("[ERR] Incorrect slave name should fail master beacon to transition")
            return
        } else {
            t.Log(err.Error())
        }
    }
    if mb.(*masterBeacon).state.(*beaconState).transitionFailureCount != 5 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).state.(*beaconState).transitionFailureCount)
        return
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}
