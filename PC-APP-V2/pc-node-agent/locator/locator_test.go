package locator

import (
    "time"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/context"
)

var masterAgentName string
var slaveNodeName string
var initSendTimestmap time.Time

func setUp() {
    masterAgentName, _ = context.DebugContextPrepare().MasterAgentName()
    slaveNodeName = "pc-node1"
    //initSendTimestmap, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    initSendTimestmap = time.Now()
    slcontext.DebugSlcontextPrepare()
}

func tearDown() {
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}
