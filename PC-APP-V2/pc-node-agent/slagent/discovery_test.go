package slagent

import (
    "fmt"
    "gopkg.in/vmihailenco/msgpack.v2"
)

func ExampleUnboundedBroadcastAgent() {
    ua, err := UnboundedBroadcastAgent()
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("Version : %s\n", ua.Version)
    fmt.Printf("MasterBoundAgent : %s\n", ua.MasterBoundAgent)
    fmt.Printf("SlaveAddress : %s\n", ua.SlaveAddress)
    fmt.Printf("SlaveGateway : %s\n", ua.SlaveGateway)
    fmt.Printf("SlaveNetmask : %s\n", ua.SlaveNetmask)
    fmt.Printf("SlaveNodeMacAddr : %s\n", ua.SlaveNodeMacAddr)

    mp, err := msgpack.Marshal(ua)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MsgPack : %v, Length : %d", mp, len(mp))
    // Output:
    // Version : 1.0.1
    // MasterBoundAgent : pc_sl_la
    // SlaveAddress : 192.168.1.236
    // SlaveGateway : 192.168.1.1
    // SlaveNetmask : ffffff00
    // SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // MsgPack : [134 168 112 99 95 115 108 95 112 100 165 49 46 48 46 49 168 112 99 95 109 97 95 98 97 168 112 99 95 115 108 95 108 97 168 112 99 95 115 108 95 105 52 173 49 57 50 46 49 54 56 46 49 46 50 51 54 168 112 99 95 115 108 95 109 97 171 49 57 50 46 49 54 56 46 49 46 49 168 112 99 95 115 108 95 109 97 168 102 102 102 102 102 102 48 48 168 112 99 95 115 108 95 109 97 177 97 99 58 98 99 58 51 50 58 57 97 58 56 100 58 54 57], Length : 123
}

func ExampleBoundedBroadcastAgent() {
    ba, err := BoundedBroadcastAgent("master-yoda")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("Version : %s\n", ba.Version)
    fmt.Printf("MasterBoundAgent : %s\n", ba.MasterBoundAgent)
    fmt.Printf("SlaveAddress : %s\n", ba.SlaveAddress)
    fmt.Printf("SlaveGateway : %s\n", ba.SlaveGateway)
    fmt.Printf("SlaveNetmask : %s\n", ba.SlaveNetmask)
    fmt.Printf("SlaveNodeMacAddr : %s\n", ba.SlaveNodeMacAddr)

    mp, err := msgpack.Marshal(ba)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MsgPack : %v, Length : %d", mp, len(mp))

    // Output:
    // Version : 1.0.1
    // MasterBoundAgent : master-yoda
    // SlaveAddress : 192.168.1.236
    // SlaveGateway : 192.168.1.1
    // SlaveNetmask : ffffff00
    // SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // MsgPack : [134 168 112 99 95 115 108 95 112 100 165 49 46 48 46 49 168 112 99 95 109 97 95 98 97 171 109 97 115 116 101 114 45 121 111 100 97 168 112 99 95 115 108 95 105 52 173 49 57 50 46 49 54 56 46 49 46 50 51 54 168 112 99 95 115 108 95 109 97 171 49 57 50 46 49 54 56 46 49 46 49 168 112 99 95 115 108 95 109 97 168 102 102 102 102 102 102 48 48 168 112 99 95 115 108 95 109 97 177 97 99 58 98 99 58 51 50 58 57 97 58 56 100 58 54 57], Length : 126
}
