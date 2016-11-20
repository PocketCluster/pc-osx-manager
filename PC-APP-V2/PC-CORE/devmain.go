package main

import (
    "time"
    "log"

    "github.com/stkim1/pc-core/teleport"
    "github.com/stkim1/pc-core/context"
)

func main() {
    context.DebugContextPrepare()
    err := teleport.StartTeleport(false)
    if err != nil {
        log.Printf("[ERR] %s", err.Error())
    }

    for {
        time.Sleep(time.Second)
    }
}