package locator

import (
    "time"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/context"
)

var (
    masterAgentName string = ""
    slaveNodeName string = "pc-node1"
    initSendTimestmap time.Time
)

func setUp() {
    masterAgentName, _ = context.DebugContextPrepare().MasterAgentName()
    initSendTimestmap = time.Now()
    slcontext.DebugSlcontextPrepare()
}

func tearDown() {
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}
