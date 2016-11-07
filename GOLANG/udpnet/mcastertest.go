package main

import (
    "log"
    "time"
    "strconv"

    "github.com/stkim1/udpnet/mcast"
)

func mcasterTest() {
    caster, err := mcast.NewMultiCaster(nil); if err != nil {
        log.Fatal("[ERR] cannot initate Multi-cast client")
    }
    for i := 0; i < 5; i++ {
        caster.Send([]byte("Hello listner - " + strconv.Itoa(i)))
    }
    time.After(time.Millisecond)
    caster.Close()
}

func main() {
    mcasterTest()
}