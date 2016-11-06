package main

import (
    "log"
    "fmt"
    "time"

    "github.com/stkim1/udpnet/ucast"
)

func ucastBeaconTest() {
    channel, err := ucast.NewPocketBeaconChannel()
    if err != nil {
        log.Fatal(err.Error())
        return
    }

    recvMsg := make(chan []byte, ucast.BEACON_RECVD_BUFFER_SIZE)
    go func() {
        err = channel.Connect(&ucast.ConnParam{
            RecvMessage : recvMsg,
            Timeout     : time.Second,
        })
        if err != nil {
            log.Printf("[ERR] cannot initate unicast client")
        }
    }()

    for {
        for entry := range recvMsg {
            fmt.Printf("Got new entry: %s\n", string(entry))
        }
        //time.Sleep(time.Second)
    }
}

func main() {
    ucastBeaconTest()
}