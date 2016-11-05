package main

import (
    "log"
    "fmt"
    "time"

    "github.com/stkim1/pc-node-agent/net/mcast"
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
}