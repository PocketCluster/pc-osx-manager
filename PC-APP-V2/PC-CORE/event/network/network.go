package network

// Package network defines an event for a system network status change

import (
    "github.com/stkim1/pc-core/context"
)

type Event struct {
    NetworkEvent
    HostInterfaces    []*context.HostNetworkInterface
    HostGateways      []*context.HostNetworkGateway
}

// Direction is the direction of the key event.
type NetworkEvent uint8

const (
    NetworkChangeInterface    NetworkEvent = iota
    NetworkChangeGateway
)
