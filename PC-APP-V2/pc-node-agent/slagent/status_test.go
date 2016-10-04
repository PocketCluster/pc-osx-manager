package slagent

import (
    "fmt"
    "time"

    "gopkg.in/vmihailenco/msgpack.v2"
)


func ExampleUnboundedStatusAgent() {
    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    sa, err := UnboundedStatusAgent(&timestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("Version : %s\n", sa.Version)
    fmt.Printf("MasterBoundAgent : %s\n", sa.MasterBoundAgent)
    fmt.Printf("SlaveNodeName : %s\n", sa.SlaveNodeName)
    fmt.Printf("SlaveAddress : %s\n", sa.SlaveAddress)
    fmt.Printf("SlaveNodeMacAddr : %s\n", sa.SlaveNodeMacAddr)
    fmt.Printf("SlaveHardware : %s\n", sa.SlaveHardware)
    fmt.Printf("SlaveTimestamp : %s\n", sa.SlaveTimestamp.String())

    mp, err := msgpack.Marshal(sa)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MsgPack %v len %d", mp, len(mp))

    // Output:
    // Version : 1.0.1
    // MasterBoundAgent : pc_sl_la
    // SlaveNodeName : MacBook-Pro-4.local
    // SlaveAddress : 192.168.1.236
    // SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // SlaveHardware : amd64
    // SlaveTimestamp : 2012-11-01 22:08:41 +0000 +0000
    // MsgPack [135 168 112 99 95 115 108 95 112 115 165 49 46 48 46 49 168 112 99 95 109 97 95 98 97 168 112 99 95 115 108 95 108 97 168 112 99 95 115 108 95 110 109 179 77 97 99 66 111 111 107 45 80 114 111 45 52 46 108 111 99 97 108 168 112 99 95 115 108 95 105 52 173 49 57 50 46 49 54 56 46 49 46 50 51 54 168 112 99 95 115 108 95 109 97 177 97 99 58 98 99 58 51 50 58 57 97 58 56 100 58 54 57 168 112 99 95 115 108 95 104 100 165 97 109 100 54 52 168 112 99 95 115 108 95 116 115 146 206 80 146 242 233 0] len 144
}

func ExampleBoundedStatusAgent() {
    timestmap, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    sa, err := BoundedStatusAgent("master-yoda", &timestmap)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("Version : %s\n", sa.Version)
    fmt.Printf("MasterBoundAgent : %s\n", sa.MasterBoundAgent)
    fmt.Printf("SlaveNodeName : %s\n", sa.SlaveNodeName)
    fmt.Printf("SlaveAddress : %s\n", sa.SlaveAddress)
    fmt.Printf("SlaveNodeMacAddr : %s\n", sa.SlaveNodeMacAddr)
    fmt.Printf("SlaveHardware : %s\n", sa.SlaveHardware)
    fmt.Printf("SlaveTimestamp : %s\n", sa.SlaveTimestamp.String())

    mp, err := msgpack.Marshal(sa)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MsgPack %v len %d", mp, len(mp))

    // Output:
    // Version : 1.0.1
    // MasterBoundAgent : master-yoda
    // SlaveNodeName : MacBook-Pro-4.local
    // SlaveAddress : 192.168.1.236
    // SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // SlaveHardware : amd64
    // SlaveTimestamp : 2012-11-01 22:08:41 +0000 +0000
    // MsgPack [135 168 112 99 95 115 108 95 112 115 165 49 46 48 46 49 168 112 99 95 109 97 95 98 97 171 109 97 115 116 101 114 45 121 111 100 97 168 112 99 95 115 108 95 110 109 179 77 97 99 66 111 111 107 45 80 114 111 45 52 46 108 111 99 97 108 168 112 99 95 115 108 95 105 52 173 49 57 50 46 49 54 56 46 49 46 50 51 54 168 112 99 95 115 108 95 109 97 177 97 99 58 98 99 58 51 50 58 57 97 58 56 100 58 54 57 168 112 99 95 115 108 95 104 100 165 97 109 100 54 52 168 112 99 95 115 108 95 116 115 146 206 80 146 242 233 0] len 147
}