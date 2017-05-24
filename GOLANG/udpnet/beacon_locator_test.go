package main

import (
    "log"
    "strconv"
    "sync"
    "time"

    "github.com/stkim1/udpnet/ucast"
)

func main() {
    var wg sync.WaitGroup
    channel, err := ucast.NewBeaconLocator(&wg)
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