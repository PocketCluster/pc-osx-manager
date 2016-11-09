package main

import (
    "testing"
    "time"
    //"bytes"

    //"github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/msagent"
    . "github.com/stkim1/pc-node-agent/locator"
    "github.com/stkim1/pcrypto"
)

var masterAgentName string
var slaveNodeName string
var initSendTimestmap time.Time

func setUp() {
    masterAgentName, _ = context.DebugContextPrepare().MasterAgentName()
    slaveNodeName = "pc-node1"
    initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    slcontext.DebugSlcontextPrepare()
}

func tearDown() {
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}

func TestInquired_KeyExchangeTransition(t *testing.T) {
    setUp()
    defer tearDown()

    meta, endTime, err := msagent.TestMasterAgentDeclarationCommand(pcrypto.TestMasterPublicKey(), initSendTimestmap)
    if err != nil {
        t.Error(err.Error())
        return
    }

    // set to slave discovery state to "Inquired"
    sd, err := NewSlaveLocator(SlaveUnbounded)
    if err != nil {
        t.Error(err.Error())
        return
    }
    //sd.(*slaveLocator).state = newInquiredState()

    // execute state transition
    if err = sd.TranstionWithMasterMeta(meta, endTime.Add(time.Second)); err != nil {
        t.Errorf(err.Error())
        return
    }
    state, err := sd.CurrentState()
    if err != nil {
        t.Error(err.Error())
        return
    }
    if state != SlaveKeyExchange {
        t.Errorf("[ERR] Slave state does not change properly | Current : %s\n", state.String())
        return
    }
}

func main() {
    t := &testing.T{}
    TestInquired_KeyExchangeTransition(t)
}