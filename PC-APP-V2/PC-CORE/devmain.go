package main

import (
    "log"

    "github.com/stkim1/pcteleport"
    "github.com/stkim1/pc-core/context"
)

func main() {
    context.DebugContextPrepare()
    err := pcteleport.StartCoreTeleport(true)
    if err != nil {
        log.Printf("[ERR] %s", err.Error())
    }
}