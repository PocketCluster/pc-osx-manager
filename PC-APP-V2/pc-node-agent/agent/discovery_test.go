package agent

import (
    "fmt"
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
    // Output:
    // Version : 1.0.1
    // MasterBoundAgent : pc_sl_la
    // SlaveAddress : 192.168.1.236
    // SlaveGateway : 192.168.1.1
    // SlaveNetmask : ffffff00
    // SlaveNodeMacAddr : ac:bc:32:9a:8d:69
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
    // Output:
    // Version : 1.0.1
    // MasterBoundAgent : master-yoda
    // SlaveAddress : 192.168.1.236
    // SlaveGateway : 192.168.1.1
    // SlaveNetmask : ffffff00
    // SlaveNodeMacAddr : ac:bc:32:9a:8d:69
}
