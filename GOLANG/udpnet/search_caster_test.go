package main

import (
    "log"
    "time"
    "strconv"
    "sync"

    "github.com/stkim1/udpnet/mcast"
)

func main() {
    var wg sync.WaitGroup
    caster, err := mcast.NewSearchCaster(&wg)
    if err != nil {
        log.Fatal("[ERR] cannot initate Multi-cast client")
    }
    for i := 0; i < 5; i++ {
        caster.Send([]byte("Hello listner - " + strconv.Itoa(i)))
    }
    time.After(time.Millisecond)
    caster.Close()
}