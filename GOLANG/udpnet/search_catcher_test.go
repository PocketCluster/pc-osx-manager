package main

import (
    "log"
    "sync"

    "github.com/stkim1/udpnet/mcast"
)

func main() {
    var wg sync.WaitGroup
    listener, err := mcast.NewSearchCatcher("en0", &wg)
    if err != nil {
        log.Fatal("[ERR] cannot initate Multi-cast client")
    }
    for v := range listener.ChRead {
        log.Println(string(v.Message) + " : " + v.Address.IP.String())
    }
}