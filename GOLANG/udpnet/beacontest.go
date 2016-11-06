package main

import (
    "log"
    "fmt"
    "github.com/stkim1/udpnet/ucast"
)

/*
func ucastOldBeaconTest() {
    channel, err := ucast.NewPocketBeaconChannel()
    if err != nil {
        log.Fatal(err.Error())
        return
    }

    recvMsg := make(chan []byte, ucast.PC_MAX_UDP_BUF_SIZE)
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
*/

func ucastNewBeaonTest() {
    channel, err := ucast.NewPocketBeaconChannel(nil)
    if err != nil {
        log.Fatal(err.Error())
        return
    }

    for entry := range channel.ChRead {
        fmt.Printf("Got new entry: %s\n", string(entry.Pack))
    }

/*
    for {
        select {
        case <- channel.ChRead:
            for entry := range channel.ChRead {
                fmt.Printf("Got new entry: %s\n", string(entry.Pack))
            }
        }
    }
*/
}

func main() {
    ucastNewBeaonTest()
}