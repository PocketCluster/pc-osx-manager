package main

import (
    "fmt"

    "github.com/stkim1/findgate"
)

func main() {
    gw, _ := findgate.DefaultIPv4Gateway()
    fmt.Printf("Interface %s | Gateway %s | Mask %s | Usable %v\n", gw.Interface, gw.Address, gw.Mask, gw.IsUsable())

    gwlist, _ := findgate.AllIPv4Gateways()
    for iface, list := range gwlist {
        for _, gw := range list {
            fmt.Printf("[%s] Interface %s | Gateway %s | Mask %s | Usable %v\n", iface, gw.Interface, gw.Address, gw.Mask, gw.IsUsable())
        }
    }
}