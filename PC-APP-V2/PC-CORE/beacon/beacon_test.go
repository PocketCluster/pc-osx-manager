package beacon

import (
    "time"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pc-core/model"
)

var (
    masterAgentName string
    slaveNodeName string
    initTime time.Time
)

func setUp() {
    mctx := context.DebugContextPrepare()
    slcontext.DebugSlcontextPrepare()
    model.DebugModelRepoPrepare()

    masterAgentName, _ = mctx.MasterAgentName()
    slaveNodeName = model.DEBUG_SLAVE_NODE_NAME
    initTime, _ = time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
}

func tearDown() {
    model.DebugModelRepoDestroy()
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}
