package network

// Package network defines an event for a system network status change

import (
    "fmt"

    "github.com/stkim1/pc-core/context"
)

// Direction is the direction of the key event.
type NetworkEvent uint8

const (
    NetworkChangeInterface    NetworkEvent = iota
    NetworkChangeGateway
)

func (n NetworkEvent) String() string {
    switch n {
    case NetworkChangeInterface:
        return "NetworkChangeInterface"
    case NetworkChangeGateway:
        return "NetworkChangeGateway"
    default:
        return fmt.Sprintf("lifecycle.Stage(%d)", n)
    }
}

type Event struct {
    NetworkEvent
    HostInterfaces    []*context.HostNetworkInterface
    HostGateways      []*context.HostNetworkGateway
}

func (e *Event) String() string {
    return e.NetworkEvent.String()
}
