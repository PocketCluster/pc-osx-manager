package main

import (
    "fmt"

    "github.com/stkim1/findgate"
)

func main() {
    gw, _ := findgate.DefaultGatewayWithInterface()
    fmt.Printf("Interface %s | Gateway %s | Mask %s\n", gw.Interface, gw.Address, gw.Mask)
}
