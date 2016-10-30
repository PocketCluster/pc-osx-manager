package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-node-agent/crypt"
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

func Test_Init_Unbound_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    sa, err := slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // TODO : how to find out this is discovery inquery?
    mb := NewBeaconForSlaveNode()
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }
    if err := mb.TranstionWithSlaveMeta(sa, initTime); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}


func Test_Unbounded_Inquired_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    mb := NewBeaconForSlaveNode()
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    sa, err := slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS := initTime
    if err := mb.TranstionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
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
    if err := mb.TranstionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterInquired {
        t.Error("[ERR] Master state is expected to be " + MasterInquired.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_Inquired_KeyExchange_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    // new beacon
    mb := NewBeaconForSlaveNode()
    if mb.CurrentState() != MasterInit {
        t.Error("[ERR] Master state is expected to be " + MasterInit.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // slave lookup master
    sa, err := slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS := initTime
    mb.TranstionWithSlaveMeta(sa, masterTS)
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // slave answer master inquiry
    slaveTS := initTime.Add(time.Second)
    sa, end, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TranstionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterInquired {
        t.Error("[ERR] Master state is expected to be " + MasterInquired.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // slave tries to key exchange
    sa, end, err = slagent.TestSlaveKeyExchangeStatus(masterAgentName, crypt.TestSlavePublicKey(), initTime)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TranstionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterKeyExchange {
        t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_KeyExchange_CryptoCheck_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    setUp()
    defer tearDown()

    // new beacon
    mb := NewBeaconForSlaveNode()
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
    mb.TranstionWithSlaveMeta(sa, masterTS)
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave answer master inquiry
    slaveTS := initTime.Add(time.Second)
    sa, end, err := slagent.TestSlaveAnswerMasterInquiry(slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TranstionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterInquired {
        t.Error("[ERR] Master state is expected to be " + MasterInquired.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave tries to key exchange
    sa, end, err = slagent.TestSlaveKeyExchangeStatus(masterAgentName, crypt.TestSlavePublicKey(), initTime)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TranstionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterKeyExchange {
        t.Error("[ERR] Master state is expected to be " + MasterKeyExchange.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave checks crypto
    t.Logf("[INFO] masterAgentName %s ", masterAgentName)
    if mb.(*masterBeacon).aesCryptor == nil {
        t.Error("[ERR] AES Cryptor is nil. Should not happen.")
        return
    }
    sa, end, err = slagent.TestSlaveCheckCryptoStatus(masterAgentName, slaveNodeName, mb.(*masterBeacon).aesCryptor, initTime)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
    if err := mb.TranstionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterCryptoCheck {
        t.Error("[ERR] Master state is expected to be " + MasterCryptoCheck.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_CryptoCheck_Bounded_Transition(t *testing.T) {
    setUp()
    defer tearDown()


    sa, end, err := slagent.TestSlaveBoundedStatus(slaveNodeName, crypt.TestAESCryptor, initTime)
    if err != nil {
        t.Error(err.Error())
        return
    }

    masterTS := end.Add(time.Second)
    mb := NewBeaconForSlaveNode()
    mb.(*masterBeacon).beaconState = MasterCryptoCheck

    if err := mb.TranstionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterBounded {
        t.Error("[ERR] Master state is expected to be " + MasterBounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_Bounded_BindBroken_Transition(t *testing.T) {
    setUp()
    defer tearDown()
}


