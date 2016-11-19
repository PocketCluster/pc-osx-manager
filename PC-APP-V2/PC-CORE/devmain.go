package main

import (
    "github.com/stkim1/pc-core/teleport"
    "github.com/stkim1/pc-core/context"
)

func main() {
    context.DebugContextPrepare()
    teleport.StartTeleport()
}