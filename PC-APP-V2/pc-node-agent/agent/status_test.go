package agent

import "fmt"

func ExampleBoundedStatusAgent() {
    ua, err := BoundedStatusAgent("master-yoda", "slave-jedi", "tz-unknown")
    if err != nil {
        fmt.Printf(err.Error())
        return
    }
    fmt.Printf("Version : %s\n", ua.Version)
    fmt.Printf("MasterBoundAgent : %s\n", ua.MasterBoundAgent)
    fmt.Printf("SlaveNodeName : %s\n", ua.SlaveNodeName)
    fmt.Printf("SlaveAddress : %s\n", ua.SlaveAddress)
    fmt.Printf("SlaveNodeMacAddr : %s\n", ua.SlaveNodeMacAddr)
    fmt.Printf("SlaveTimeZone : %s\n", ua.SlaveTimeZone)

    // Output:
    // Version : 1.0.1
    // MasterBoundAgent : master-yoda
    // SlaveNodeName : slave-jedi
    // SlaveAddress : 192.168.1.236
    // SlaveNodeMacAddr : ac:bc:32:9a:8d:69
    // SlaveTimeZone : tz-unknown
}

