package beacon

import (
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-node-agent/slcontext"
    "testing"
)

func setUp() {
    context.DebugContextPrepare()
    slcontext.DebugSlcontextPrepare()
}

func tearDown() {
    slcontext.DebugSlcontextDestroy()
    context.DebugContextDestroy()
}

func Test_Unbounded_Inquired_Transition(t *testing.T) {
    setUp()
    defer tearDown()





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


