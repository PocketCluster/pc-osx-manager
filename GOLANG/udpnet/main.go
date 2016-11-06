package main

import (
    "log"
    "fmt"
    "time"

    "github.com/stkim1/udpnet/mcast"
    "github.com/stkim1/udpnet/ucast"
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

func ucastLocatorTest() {
    channel, err := ucast.NewPocketLocatorChannel()
    if err != nil {
        log.Fatal(err.Error())
        return
    }

    recvMsg := make(chan []byte, 32)
    go func() {
        for entry := range recvMsg {
            fmt.Printf("Got new entry: %s\n", string(entry))
        }
    }()

    err = channel.Connect(&ucast.ConnParam{
        RecvMessage : recvMsg,
        Timeout     : time.Second,
    })
    if err != nil {
        log.Printf("[ERR] cannot initate Multi-cast client")
    }

    for i := 3; i < 3; i++ {
        channel.Send("192.168.1.152", []byte("HELLO!"))
    }
    for {
        time.Sleep(time.Second)
    }
}


func main() {
    ucastLocatorTest()
}