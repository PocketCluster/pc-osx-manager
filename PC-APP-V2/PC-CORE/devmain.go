package main

import (
    "github.com/stkim1/pc-core/swarm"
)
func main() {
    context := swarm.NewContext("localhost:3275", "1810ffdf37ad898423ada7262f7baf80")
    context.Manage()
}