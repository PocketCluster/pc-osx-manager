package slagent

import (
    "fmt"
    "gopkg.in/vmihailenco/msgpack.v2"
)

func ExampleBoundedStatusAgent() {
    sa, err := BoundedStatusAgent("master-yoda", "slave-jedi", "tz-unknown")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("Version : %s\n", sa.Version)
    fmt.Printf("MasterBoundAgent : %s\n", sa.MasterBoundAgent)
    fmt.Printf("SlaveNodeName : %s\n", sa.SlaveNodeName)
    fmt.Printf("SlaveAddress : %s\n", sa.SlaveAddress)
    fmt.Printf("SlaveNodeMacAddr : %s\n", sa.SlaveNodeMacAddr)
    fmt.Printf("SlaveTimeZone : %s\n", sa.SlaveTimeZone)

    // Output:
    // Version : 1.0.1
    // MasterBoundAgent : master-yoda
    // SlaveNodeName : slave-jedi
    // SlaveAddress : 192.168.1.236
    // SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // SlaveTimeZone : tz-unknown
}

func ExampleBoundedStatusAgentMsgPack() {
    sa, err := BoundedStatusAgent("master-yoda", "slave-jedi", "tz-unknown")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    mp, err := msgpack.Marshal(sa)
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("MsgPack %v len %d", mp, len(mp))
    // Output:
    // MsgPack [134 168 112 99 95 115 108 95 115 116 165 49 46 48 46 49 168 112 99 95 109 97 95 98 97 171 109 97 115 116 101 114 45 121 111 100 97 168 112 99 95 115 108 95 110 109 170 115 108 97 118 101 45 106 101 100 105 168 112 99 95 115 108 95 105 52 173 49 57 50 46 49 54 56 46 49 46 50 51 54 168 112 99 95 115 108 95 109 97 177 97 99 58 98 99 58 51 50 58 57 97 58 56 100 58 54 57 168 112 99 95 115 108 95 116 122 170 116 122 45 117 110 107 110 111 119 110] len 127
}
