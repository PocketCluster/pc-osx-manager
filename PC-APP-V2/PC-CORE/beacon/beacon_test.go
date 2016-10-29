package beacon

import (
    "testing"
    "time"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-node-agent/slagent"
)

var masterAgentName string
var slaveNodeName string
var initSendTimestmap time.Time

func setUp() {
    mctx := context.DebugContextPrepare()
    slcontext.DebugSlcontextPrepare()

    masterAgentName, _ = mctx.MasterAgentName()
    slaveNodeName = "pc-node1"
    initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
}

func tearDown() {
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}

func Test_Unbounded_Inquired_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    sm, err := slagent.TestSlaveUnboundedMasterSearchDiscovery()
    if err != nil {
        t.Error(err.Error())
        return
    }

    // TODO : how to find out this is discovery inquery?
    mb := NewBeaconForSlaveNode()
    if mb.CurrentState() != MasterUnbounded {
        t.Error("[ERR] Master state is expected to be " + MasterUnbounded.String() + ". Current : " + mb.CurrentState().String())
        return
    }

    mb.TranstionWithSlaveMeta(sm, initSendTimestmap)
    if mb.CurrentState() != MasterInquired {
        t.Error("[ERR] Master state is expected to be " + MasterInquired.String() + ". Current : " + mb.CurrentState().String())
        return
    }
}

func Test_Inquired_KeyExchange_Transition(t *testing.T) {
    setUp()
    defer tearDown()

    sa, end, err := slagent.TestSlaveAnswerMasterInquiry(initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    masterTS := end.Add(time.Second)
    mb := NewBeaconForSlaveNode()
    mb.(*masterBeacon).beaconState = MasterInquired
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

    sa, end, err := slagent.TestSlaveKeyExchangeStatus(masterAgentName, initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    masterTS := end.Add(time.Second)
    mb := NewBeaconForSlaveNode()
    mb.(*masterBeacon).beaconState = MasterKeyExchange
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

    sa, end, err := slagent.TestSlaveCheckCryptoStatus(masterAgentName, slaveNodeName, initSendTimestmap)
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


