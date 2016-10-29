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
    mb.TranstionWithSlaveMeta(sm, initSendTimestmap)
}

func Test_Inquired_KeyExchange_Transition(t *testing.T) {
    setUp()
    defer tearDown()

}

func Test_KeyExchange_CryptoCheck_Transition(t *testing.T) {
    setUp()
    defer tearDown()
}


func Test_CryptoCheck_Bounded_Transition(t *testing.T) {
    setUp()
    defer tearDown()
}


