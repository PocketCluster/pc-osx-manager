package main

import (
    "log"

    "github.com/stkim1/udpnet/mcast"
)

func mlistenerTest() {
    listener, err := mcast.NewMultiListener(nil); if err != nil {
        log.Fatal("[ERR] cannot initate Multi-cast client")
    }

    for v := range listener.ChRead {
        log.Println(string(v.Message) + " : " + v.Address.IP.String())
    }
}

func main() {
    mlistenerTest()
}