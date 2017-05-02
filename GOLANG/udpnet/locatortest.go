package main

import (
    "log"
    "strconv"
    "time"

    "github.com/stkim1/udpnet/ucast"
)

func ucastLocatorTest() {
    channel, err := ucast.NewPocketLocatorChannel()
    if err != nil {
        log.Fatal(err.Error())
        return
    }

    for i := 0; i < 10; i++ {
        err := channel.Send("192.168.1.220", []byte("HELLO! - " + strconv.Itoa(i)))
        if err != nil {
            log.Fatal(err.Error())
        }
    }
    time.After(time.Millisecond)
    channel.Close()
}

func main() {
    ucastLocatorTest()
}