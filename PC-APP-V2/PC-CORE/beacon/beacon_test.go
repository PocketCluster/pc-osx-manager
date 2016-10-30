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


func Test_Init_Bounded_OnePass_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    // test var preperations
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

    // --- slave tries to key exchange
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveKeyExchangeStatus(masterAgentName, crypt.TestSlavePublicKey(), slaveTS)
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
    if err := mb.TranstionWithSlaveMeta(sa, masterTS); err != nil {
        t.Error(err.Error())
        return
    }
    if mb.CurrentState() != MasterCryptoCheck {
        t.Error("[ERR] Master state is expected to be " + MasterCryptoCheck.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    // --- slave is now bounded
    slaveTS = masterTS.Add(time.Second)
    sa, end, err = slagent.TestSlaveBoundedStatus(slaveNodeName, mb.(*masterBeacon).aesCryptor, slaveTS)
    if err != nil {
        t.Error(err.Error())
        return
    }
    masterTS = end.Add(time.Second)
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


