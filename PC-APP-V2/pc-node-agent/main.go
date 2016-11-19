package main

import (
    "log"
    "net"

    "github.com/stkim1/netifaces"
)

func main() {
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
}