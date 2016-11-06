package main

import (
    "log"
    "fmt"
    "time"

    "github.com/stkim1/udpnet/ucast"
)

func ucastLocatorTest() {
    channel, err := ucast.NewPocketLocatorChannel()
    if err != nil {
        log.Fatal(err.Error())
        return
    }

    log.Print("[INFO] let's try to listen")
    recvMsg := make(chan []byte, ucast.BEACON_RECVD_BUFFER_SIZE)
    err = channel.Connect(&ucast.ConnParam{
        RecvMessage : recvMsg,
        Timeout     : time.Second,
    })
    if err != nil {
        log.Printf("[ERR] cannot initate unicast client")
    }

    for i := 0; i < 3; i++ {
        log.Print("[INFO] send HELLO! to 192.168.1.152")
        channel.Send("192.168.1.152", []byte("HELLO!"))
    }
    for {
        for entry := range recvMsg {
            fmt.Printf("Got new entry: %s\n", string(entry))
        }
        time.Sleep(time.Second)
    }
}


func main() {
    ucastLocatorTest()
}