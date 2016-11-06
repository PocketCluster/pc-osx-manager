package main

import (
    "log"
    "fmt"

    "github.com/stkim1/udpnet/ucast"
)

func ucastBeaonTest() {
    channel, err := ucast.NewPocketBeaconChannel(nil)
    if err != nil {
        log.Fatal(err.Error())
        return
    }
    for entry := range channel.ChRead {
        fmt.Printf("Got new entry: %s\n", string(entry.Pack))
    }
}

func main() {
    ucastBeaonTest()
}