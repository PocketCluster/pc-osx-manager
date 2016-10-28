package slcontext

import "testing"

func setUp() {
    DebugSlcontextPrepare()
}

func tearDown() {
    DebugSlcontextDestroy()
}

func TestGetDefaultIP4Gateway(t *testing.T) {
    setUp()
    defer tearDown()




}
