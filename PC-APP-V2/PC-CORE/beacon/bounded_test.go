package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pcrypto"
)

func Test_BeaconInit_Bounded_OnePass_Transition(t *testing.T) {
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
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave answer master inquiry
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
    if mb.CurrentState() != MasterInquired {
        t.Error("[ERR] Master state is expected to be " + MasterInquired.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave tries to key exchange
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
    if mb.CurrentState() != MasterKeyExchange {
        t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave checks crypto
    if mb.(*masterBeacon).state.(DebugState).AESCryptor() == nil {
        t.Error("[ERR] AES Cryptor is nil. Should not happen.")
        return
    }
    if len(mb.SlaveNode().NodeName) == 0 {
        t.Errorf("[ERR] Slave node name should not be empty")
        return
    }
    aescryptor := mb.(*masterBeacon).state.(DebugState).AESCryptor()
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, mb.SlaveNode().NodeName, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterCryptoCheck {
        t.Error("[ERR] Master state is expected to be " + MasterCryptoCheck.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave is now bounded
    slaveTS = masterTS.Add(time.Second)
    aescryptor = mb.(*masterBeacon).state.(DebugState).AESCryptor()
    sa, end, err = slagent.TestSlaveBoundedStatus(masterAgentName, slaveNodeName, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterBounded {
        t.Error("[ERR] Master state is expected to be " + MasterBounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_Bounded_BindBroken_TimeoutFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        masterTS, slaveTS time.Time = time.Now(), time.Now()
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

    // --- slave lookup master
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

    // --- slave answer master inquiry
    slaveTS = masterTS.Add(time.Second)
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
    if mb.CurrentState() != MasterInquired {
        t.Error("[ERR] Master state is expected to be " + MasterInquired.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave tries to key exchange
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
    if mb.CurrentState() != MasterKeyExchange {
        t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave checks crypto
    if mb.(*masterBeacon).state.(DebugState).AESCryptor() == nil {
        t.Error("[ERR] AES Cryptor is nil. Should not happen.")
        return
    }
    if len(mb.SlaveNode().NodeName) == 0 {
        t.Errorf("[ERR] Slave node name should not be empty")
        return
    }
    aescryptor := mb.(*masterBeacon).state.(DebugState).AESCryptor()
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, mb.SlaveNode().NodeName, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterCryptoCheck {
        t.Error("[ERR] Master state is expected to be " + MasterCryptoCheck.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave is now bounded
    slaveTS = masterTS.Add(time.Second)
    aescryptor = mb.(*masterBeacon).state.(DebugState).AESCryptor()
    sa, end, err = slagent.TestSlaveBoundedStatus(masterAgentName, slaveNodeName, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterBounded {
        t.Error("[ERR] Master state is expected to be " + MasterBounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_Bounded_BindBroken_TooManyMetaFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        masterTS, slaveTS time.Time = time.Now(), time.Now()
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

    // --- slave lookup master
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

    // --- slave answer master inquiry
    slaveTS = masterTS.Add(time.Second)
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
    if mb.CurrentState() != MasterInquired {
        t.Error("[ERR] Master state is expected to be " + MasterInquired.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave tries to key exchange
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
    if mb.CurrentState() != MasterKeyExchange {
        t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave checks crypto
    if mb.(*masterBeacon).state.(DebugState).AESCryptor() == nil {
        t.Error("[ERR] AES Cryptor is nil. Should not happen.")
        return
    }
    if len(mb.SlaveNode().NodeName) == 0 {
        t.Errorf("[ERR] Slave node name should not be empty")
        return
    }
    aescryptor := mb.(*masterBeacon).state.(DebugState).AESCryptor()
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, mb.SlaveNode().NodeName, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterCryptoCheck {
        t.Error("[ERR] Master state is expected to be " + MasterCryptoCheck.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave is now bounded
    slaveTS = masterTS.Add(time.Second)
    aescryptor = mb.(*masterBeacon).state.(DebugState).AESCryptor()
    sa, end, err = slagent.TestSlaveBoundedStatus(masterAgentName, slaveNodeName, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterBounded {
        t.Error("[ERR] Master state is expected to be " + MasterBounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_Bounded_BindBroken_TxActionFail(t *testing.T) {
    setUp()
    defer tearDown()

    var (
        debugComm CommChannel = &DebugCommChannel{}
        masterTS, slaveTS time.Time = time.Now(), time.Now()
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

    // --- slave lookup master
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

    // --- slave answer master inquiry
    slaveTS = masterTS.Add(time.Second)
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
    if mb.CurrentState() != MasterInquired {
        t.Error("[ERR] Master state is expected to be " + MasterInquired.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave tries to key exchange
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
    if mb.CurrentState() != MasterKeyExchange {
        t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave checks crypto
    if mb.(*masterBeacon).state.(DebugState).AESCryptor() == nil {
        t.Error("[ERR] AES Cryptor is nil. Should not happen.")
        return
    }
    if len(mb.SlaveNode().NodeName) == 0 {
        t.Errorf("[ERR] Slave node name should not be empty")
        return
    }
    aescryptor := mb.(*masterBeacon).state.(DebugState).AESCryptor()
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, mb.SlaveNode().NodeName, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterCryptoCheck {
        t.Error("[ERR] Master state is expected to be " + MasterCryptoCheck.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave is now bounded
    slaveTS = masterTS.Add(time.Second)
    aescryptor = mb.(*masterBeacon).state.(DebugState).AESCryptor()
    sa, end, err = slagent.TestSlaveBoundedStatus(masterAgentName, slaveNodeName, aescryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterBounded {
        t.Error("[ERR] Master state is expected to be " + MasterBounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- TX ACTION FAIL ---
    for i := 0; i <= int(TxActionLimit); i++ {
        masterTS = masterTS.Add(time.Millisecond + BoundedTimeout)
        err = mb.TransitionWithTimestamp(masterTS)
        if err != nil {
            t.Logf("%th %s",i + 1, err.Error())
        }
    }
    if mb.CurrentState() != MasterBindBroken {
        t.Error("[ERR] Master state is expected to be " + MasterBindBroken.String() + ". Current : " + mb.CurrentState().String())
        return
    }

}
