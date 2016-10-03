package main

import (
    "log"
    "fmt"
    "time"

    "github.com/stkim1/pc-node-agent/mcast"
    "github.com/stkim1/pc-node-agent/agent"
)

func mcastTest() {
    client, err := mcast.NewClient(); if err != nil {
        log.Fatal("[ERR] cannot initate Multi-cast client")
    }

    recvMsg := make(chan []byte, 32)
    param := mcast.DefaultParams([]byte("this is message"), recvMsg)
    client.Query(param)

    go func() {
        for entry := range recvMsg {
            fmt.Printf("Got new entry: %s\n", string(entry))
        }
    }()

    // Start the lookup
    client.Query(param)

    for {
        time.Sleep(time.Second)
    }
}

func main() {

/*
    ifaces, err := status.Interfaces()
    if err != nil {
        log.Printf("Cannot acquire interface info %v", err)
    }

    for _, ifs := range ifaces {
        fmt.Println(ifs.Name)
        fmt.Println(ifs.HardwareAddr)
        fmt.Println(ifs.Flags)
        addrs, _ := status.IP4Addrs(ifs)
        // handle err
        for _, addr := range addrs {
            fmt.Println("\t" + addr.IPString())
            fmt.Println("\t" + addr.IPMaskString())
        }
        fmt.Println("--------\n")
    }
*/
    pa, err := agent.UnboundedBroadcastAgent(); if err != nil {
        return
    }
    fmt.Printf("addr %s mac %s", pa.SlaveAddress, pa.SlaveNodeMacAddr)
}