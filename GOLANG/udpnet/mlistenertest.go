package main

import (
    "log"
    "net"

    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/netifaces"
)

func mlistenerTest() {
    gateway, err := netifaces.FindSystemGateways()
    if err != nil {
        log.Panic(err.Error())
    }
    defer gateway.Release()

    _, gwiface, _ := gateway.DefaultIP4Gateway()
    iface, _ := net.InterfaceByName(gwiface)
    log.Print("[INFO] we'll start listening from " + gwiface)

    listener, err := mcast.NewMultiListener(iface, nil); if err != nil {
        log.Fatal("[ERR] cannot initate Multi-cast client")
    }

    for v := range listener.ChRead {
        log.Println(string(v.Message) + " : " + v.Address.IP.String())
    }
}

func main() {
    mlistenerTest()
}