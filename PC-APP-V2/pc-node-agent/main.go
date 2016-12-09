package main

import (
    "log"
    "net"
    "time"

    "github.com/stkim1/netifaces"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/stkim1/pcteleport"
)

func main() {
    gateway, err := netifaces.FindSystemGateways()
    if err != nil {
        log.Print(err.Error())
    }
    gwaddr, gwiface, err := gateway.DefaultIP4Gateway()
    log.Printf("GW ADDR %s | GW IFACE %s", gwaddr, gwiface)
    _, err = net.InterfaceByName(gwiface)
    if err != nil {
        log.Print(err.Error())
    }
    gateway.Release()

    slcontext.DebugSlcontextPrepare()
    slcontext.SharedSlaveContext().SetSlaveNodeName("pc-node1")
    err = pcteleport.StartNodeTeleport("192.168.1.150", "c9s93fd9-3333-91d3-9999-c9s93fd98f43", true)
    if err != nil {
        log.Print(err.Error())
    }
    for {
        time.Sleep(time.Second)
    }
    slcontext.DebugSlcontextDestroy()
}