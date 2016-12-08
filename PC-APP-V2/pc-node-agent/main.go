package main

import (
    "log"
    "net"

    "github.com/stkim1/netifaces"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pcteleport"
)

func main() {
    slcontext.DebugSlcontextPrepare()
    defer slcontext.DebugSlcontextDestroy()
    gateway, err := netifaces.FindSystemGateways()
    if err != nil {
        log.Print(err.Error())
    }
    defer gateway.Release()

    gwaddr, gwiface, err := gateway.DefaultIP4Gateway()
    log.Printf("GW ADDR %s | GW IFACE %s", gwaddr, gwiface)
    _, err = net.InterfaceByName(gwiface)
    if err != nil {
        log.Print(err.Error())
    }

    err = pcteleport.StartNodeTeleport("192.168.1.150", "c9s93fd9-3333-91d3-9999-c9s93fd98f43", true)
    if err != nil {
        log.Print(err.Error())
    }
}