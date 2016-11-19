package main

import (
    "github.com/stkim1/pc-core/teleport"
    "github.com/stkim1/pc-core/context"
    "time"
)

func main() {
    context.DebugContextPrepare()
    teleport.StartTeleport()

    for {
        time.Sleep(time.Second)
    }
}