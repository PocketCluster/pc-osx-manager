package main

import (
    "github.com/stkim1/pc-core/swarm"
)
func main() {
    context := swarm.NewContext("localhost:3275", "192.168.1.151:2375,192.168.1.152:2375")
    context.Manage()
}