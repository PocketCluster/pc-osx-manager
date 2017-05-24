package locator

import (
    "time"

    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/context"
    "github.com/pborman/uuid"
)

var (
    masterAgentName, slaveNodeName, slaveUUID string
    initSendTimestmap time.Time
)

func setUp() {
    masterAgentName, _ = context.DebugContextPrepare().MasterAgentName()
    slaveNodeName = "pc-node1"
    slaveUUID = uuid.New()
    initSendTimestmap = time.Now()
    slcontext.DebugSlcontextPrepare()
}

func tearDown() {
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}
