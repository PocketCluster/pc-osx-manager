package main

import (
    "log"
    "time"
    "io/ioutil"

    "github.com/coreos/etcd/embed"
)

func main() {
    cert, _ := ioutil.ReadFile("/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.cert")
    key, _  := ioutil.ReadFile("/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.key")
    ca, _   := ioutil.ReadFile("/Users/almightykim/Workspace/DKIMG/CERT/ca-cert.pub")
    cfg, err := embed.NewPocketConfig("/Users/almightykim/Workspace/DKIMG/ETCD/data", ca, cert, key)
    if err != nil {
        log.Fatal(err)
    }
    e, err := embed.StartPocketEtcd(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer e.Close()
    select {
    case <-e.Server.ReadyNotify():
        log.Printf("Server is ready!")
    case <-time.After(60 * time.Second):
        e.Server.Stop() // trigger a shutdown
        log.Printf("Server took too long to start!")
    }
    log.Fatal(<-e.Err())
}