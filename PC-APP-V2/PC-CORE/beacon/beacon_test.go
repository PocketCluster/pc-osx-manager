package beacon
/*
import (
    "testing"
    "time"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pcrypto"
)

var masterAgentName string
var slaveNodeName string
var initTime time.Time

func setUp() {
    mctx := context.DebugContextPrepare()
    slcontext.DebugSlcontextPrepare()
    model.DebugModelRepoPrepare()

    masterAgentName, _ = mctx.MasterAgentName()
    slaveNodeName = "pc-node1"
    initTime, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
}

func tearDown() {
    model.DebugModelRepoDestroy()
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}


func Test_Init_Bounded_OnePass_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    // test var preperations
    mb := NewMasterBeacon(MasterInit, nil, nil)
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
    masterTS := initTime
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
    if mb.(*masterBeacon).aesCryptor == nil {
        t.Error("[ERR] AES Cryptor is nil. Should not happen.")
        return
    }
    if len(mb.SlaveNode().NodeName) == 0 {
        t.Errorf("[ERR] Slave node name should not be empty")
        return
    }
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, mb.SlaveNode().NodeName, mb.(*masterBeacon).aesCryptor, slaveTS)
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
    sa, end, err = slagent.TestSlaveBoundedStatus(masterAgentName, slaveNodeName, mb.(*masterBeacon).aesCryptor, slaveTS)
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

func Test_Init_Unbounded_Transition_Fail(t *testing.T) {
    setUp()
    defer tearDown()

    // test var preperations
    mb := NewBeaconForSlaveNode()
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- TIMEOUT FAILURE ---
    sa, err := slagent.TestSlaveBindBroken(masterAgentName)
    if err != nil {
        t.Skip(err.Error())
    }

    // 1st trial with incorrect slave meta
    masterTS := time.Now()
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err == nil {
        t.Errorf("[ERR] incorrect slave state should generate error when fed to freshly spwaned beacon")
        return
    } else {
        t.Logf(err.Error())
    }
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).trialFailCount != 1 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
    // 2nd trial with TS +10 sec
    masterTS = masterTS.Add(time.Second * 10)
    if err := mb.TransitionWithTimestamp(masterTS); err != nil {
        t.Log(err.Error())
    }
    if mb.(*masterBeacon).trialFailCount != 1 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).trialFailCount)
        return
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- TOO MANY TIMES FAILURE ---
    mb = NewBeaconForSlaveNode()
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave lookup master
    sa, err = slagent.TestSlaveBindBroken(masterAgentName)
    if err != nil {
        t.Skip(err.Error())
    }

    // 1st trial with incorrect slave meta
    masterTS = time.Now()
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err == nil {
        t.Errorf("[ERR] incorrect slave state should generate error when fed to freshly spwaned beacon")
        return
    }
    // four more times of failure with incorrect slave meta
    for i := 0; i < 4; i++ {
        masterTS = masterTS.Add(time.Second)
        if err := mb.TransitionWithSlaveMeta(sa, masterTS); err == nil {
            t.Errorf("[ERR] incorrect slave state should generate error when fed to freshly spwaned beacon")
            return
        }
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).trialFailCount != 5 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
}

func Test_Unbounded_Inquired_Transition_Fail(t *testing.T) {
    setUp()
    defer tearDown()

    // --- TIMEOUT FAILURE ---
    masterTS := time.Now()
    mb := NewBeaconForSlaveNode()
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    mb.(*masterBeacon).beaconState = MasterUnbounded
    // --- test
    slaveTS := masterTS.Add(time.Second)
    sa, _, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
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
    if mb.(*masterBeacon).trialFailCount != 1 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
    // update with timestamp
    slaveTS = masterTS.Add(time.Second * 11)
    t.Logf("[INFO] slaveTS - MasterBeacon.lastSuccessTimestmap : " + slaveTS.Sub(mb.(*masterBeacon).lastSuccess).String())
    if err := mb.TransitionWithTimestamp(slaveTS); err != nil {
        t.Log(err.Error())
    }
    if mb.(*masterBeacon).trialFailCount != 1 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).trialFailCount)
        return
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }


    // --- TOO MANY TIMES FAILURE ---
    masterTS = time.Now()
    mb = NewBeaconForSlaveNode()
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    mb.(*masterBeacon).beaconState = MasterUnbounded
    // --- test
    for i := 0; i < TransitionFailureLimit; i ++ {
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
    if mb.(*masterBeacon).trialFailCount != 5 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
}

func Test_Inquired_KeyExchange_Fail(t *testing.T) {
    setUp()
    defer tearDown()

    // --- TIMEOUT FAILURE ---
    mb := NewBeaconForSlaveNode()
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
    if mb.(*masterBeacon).trialFailCount != 1 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
    // fail with timestamp
    masterTS = masterTS.Add(time.Second * 11)
    t.Logf("[INFO] slaveTS - MasterBeacon.lastSuccessTimestmap : " + slaveTS.Sub(mb.(*masterBeacon).lastSuccess).String())
    if err := mb.TransitionWithTimestamp(masterTS); err != nil {
        t.Log(err.Error())
    }
    if mb.(*masterBeacon).trialFailCount != 1 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).trialFailCount)
        return
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }


    // --- TOO MANY TIMES FAILURE ---
    mb = NewBeaconForSlaveNode()
    sa, err = slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = time.Now()
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveAnswerMasterInquiry(slaveTS)
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
    for i := 0; i < TransitionFailureLimit; i++ {
        t.Logf("[INFO] Master state : %s. Trial count %d", mb.CurrentState().String(), mb.(*masterBeacon).trialFailCount)

        if mb.(*masterBeacon).trialFailCount != i {
            t.Error("[ERR] Master fail count [%d] should match with trial count [%d]", mb.(*masterBeacon).trialFailCount, i)
            return
        }
        if mb.CurrentState() != MasterInquired {
            t.Errorf("[ERR] Master state is expected to be %s. Current : %s. Trial count %d", MasterInquired.String(), mb.CurrentState().String(), mb.(*masterBeacon).trialFailCount)
            //return
        }
        slaveTS = masterTS.Add(time.Second)
        t.Logf("[INFO] slaveTS - MasterBeacon.lastSuccessTimestmap : " + slaveTS.Sub(mb.(*masterBeacon).lastSuccess).String())
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
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).trialFailCount != 5 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
}

func Test_KeyExchange_CryptoCheck_Fail(t *testing.T) {
    setUp()
    defer tearDown()

    // --- TIMEOUT FAILURE ---
    mb := NewBeaconForSlaveNode()
    sa, err := slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS := time.Now()
    mb.TransitionWithSlaveMeta(sa, masterTS)
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
    // --- test
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, "INCORRECT-SLAVE-NAME", mb.(*masterBeacon).aesCryptor, slaveTS)
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
    if mb.CurrentState() != MasterKeyExchange {
        t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).trialFailCount != 1 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
    // fail with timestamp
    masterTS = masterTS.Add(time.Second * 11)
    t.Logf("[INFO] masterTS - MasterBeacon.lastSuccessTimestmap : " + masterTS.Sub(mb.(*masterBeacon).lastSuccess).String())
    if err := mb.TransitionWithTimestamp(masterTS); err != nil {
        t.Log(err.Error())
    }
    if mb.(*masterBeacon).trialFailCount != 1 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).trialFailCount)
        return
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }


    // --- TOO MANY TIMES FAILURE ---
    mb = NewBeaconForSlaveNode()
    sa, err = slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = time.Now()
    mb.TransitionWithSlaveMeta(sa, masterTS)
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveAnswerMasterInquiry(slaveTS)
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
    // --- test
    for i := 0; i < TransitionFailureLimit; i++ {
        t.Logf("[INFO] Master state : %s. Trial count %d", mb.CurrentState().String(), mb.(*masterBeacon).trialFailCount)

        if mb.CurrentState() != MasterKeyExchange {
            t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
            return
        }
        slaveTS = masterTS.Add(time.Second)
        t.Logf("[INFO] slaveTS - MasterBeacon.lastSuccessTimestmap : " + slaveTS.Sub(mb.(*masterBeacon).lastSuccess).String())
        sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, "INCORRECT-SLAVE-NAME", mb.(*masterBeacon).aesCryptor, slaveTS)
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
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if mb.(*masterBeacon).trialFailCount != 5 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
}

func Test_KeyExchange_Bounded_Fail(t *testing.T) {
    setUp()
    defer tearDown()

    // --- TIMEOUT FAILURE ---
    mb := NewBeaconForSlaveNode()
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
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, mb.SlaveNode().NodeName, mb.(*masterBeacon).aesCryptor, slaveTS)
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
    sa, end, err = slagent.TestSlaveBoundedStatus(masterAgentName, "INCORRECT-SLAVE-NAME", mb.(*masterBeacon).aesCryptor, slaveTS)
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
    if mb.(*masterBeacon).trialFailCount != 1 {
        t.Error("[ERR] Master fail count should have increased")
        return
    }
    // fail with timestamp
    masterTS = masterTS.Add(time.Second * 11)
    t.Logf("[INFO] masterTS - MasterBeacon.lastSuccessTimestmap : " + masterTS.Sub(mb.(*masterBeacon).lastSuccess).String())
    if err := mb.TransitionWithTimestamp(masterTS); err != nil {
        t.Log(err.Error())
    }
    if mb.(*masterBeacon).trialFailCount != 1 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).trialFailCount)
        return
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }


    // --- TOO MANY TIMES FAILURE ---
    mb = NewBeaconForSlaveNode()
    sa, err = slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = time.Now()
    if err := mb.TransitionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveAnswerMasterInquiry(slaveTS)
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
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, mb.SlaveNode().NodeName, mb.(*masterBeacon).aesCryptor, slaveTS)
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
    for i := 0; i < TransitionFailureLimit; i++ {
        if mb.CurrentState() != MasterCryptoCheck {
            t.Error("[ERR] Master state is expected to be " + MasterCryptoCheck.String() + ". Current : " + mb.CurrentState().String())
            return
        }
        slaveTS = masterTS.Add(time.Second)
        sa, end, err = slagent.TestSlaveBoundedStatus(masterAgentName, "INCORRECT-SLAVE-NAME", mb.(*masterBeacon).aesCryptor, slaveTS)
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
    if mb.(*masterBeacon).trialFailCount != 5 {
        t.Errorf("[ERR] Master fail count should have increased. Current count %d", mb.(*masterBeacon).trialFailCount)
        return
    }
    if mb.CurrentState() != MasterDiscarded {
        t.Error("[ERR] Master state is expected to be " + MasterDiscarded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_Bounded_BindBroken_Transition(t *testing.T) {
    setUp()
    defer tearDown()
}

*/