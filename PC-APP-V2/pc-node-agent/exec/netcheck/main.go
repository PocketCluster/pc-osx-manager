package main

import (
    "fmt"
    "net"

    "github.com/stkim1/netifaces"
    "github.com/stkim1/pc-node-agent/slcontext"
    "github.com/davecgh/go-spew/spew"
)

func main_uwrapped() {

    // 2017-05-15 network address and interface were checked on RPI3 and Odroid with Golang 1.7.5
    // Both hardware works fine without any error! ;)

    gateway, err := netifaces.FindSystemGateways()
    if err != nil {
        fmt.Printf(err.Error())
    }
    gwaddr, gwiface, err := gateway.DefaultIP4Gateway()
    if err != nil {
        fmt.Printf(err.Error())
    }
    fmt.Printf("gateway address %v | gateway interface %v\n", gwaddr, gwiface)
    iface, err := net.InterfaceByName(gwiface)
    if err != nil {
        fmt.Printf(err.Error())
    }
    ifAddrs, err := iface.Addrs()
    if err != nil {
        fmt.Printf(err.Error())
    }
    for _, a := range ifAddrs {
        fmt.Printf("address %v | network %v\n", a.String(), a.Network())
    }
    gateway.Release()
}

func main() {
    niface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Print(spew.Sdump(niface))
}