package main

import (
    "log"
    "net"
    "net/rpc"
    "sync"
)

func runServer(wg *sync.WaitGroup) {
    defer wg.Done()

    rpc.Register(NewRPC())
    l, e := net.Listen("tcp", dsn)
    if e != nil {
        log.Fatal("listen error:", e)
    }
    rpc.Accept(l)
}
