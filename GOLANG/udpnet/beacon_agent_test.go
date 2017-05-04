package main

import (
    "fmt"
    "log"
    "sync"

    "github.com/stkim1/udpnet/ucast"
)

func main() {
    var wg sync.WaitGroup
    channel, err := ucast.NewBeaconAgent(&wg)
    if err != nil {
        log.Fatal(err.Error())
        return
    }
    for entry := range channel.ChRead {
        fmt.Printf("Got new entry: %s\n", string(entry.Message))
    }
}