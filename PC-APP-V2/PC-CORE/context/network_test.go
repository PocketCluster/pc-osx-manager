package context

import (
    "testing"
)

var intefaces0 = []*HostNetworkInterface{
    {
        WifiPowerOff        : false,
        IsActive            : true,
        IsPrimary           : true,
        AddrCount           : 1,
        Address             : []*HostIPAddress{
            {
                Flags       : 0x8863,
                Family      : 2,
                IsPrimary   : true,
                Address     : "192.168.1.248",
                Netmask     : "255.255.255.0",
                Broadcast   : "192.168.1.255",
            },
        },
        BsdName             : "en0",
        DisplayName         : "Ethernet",
        MacAddress          : "74:d4:35:f3:b5:20",
        MediaType           : "Ethernet",
    },
    {
        WifiPowerOff        : false,
        IsActive            : true,
        IsPrimary           : false,
        AddrCount           : 1,
        Address             : []*HostIPAddress{
            {
                Flags       : 0x8863,
                Family      : 2,
                IsPrimary   : true,
                Address     : "192.168.1.247",
                Netmask     : "255.255.255.0",
                Broadcast   : "192.168.1.255",
            },
        },
        BsdName             : "en1",
        DisplayName         : "Wi-Fi",
        MacAddress          : "74:d4:35:f3:b5:20",
        MediaType           : "IEEE80211",
    },
    {
        BsdName             : "lo0",
        Address             : nil,
    },
    {
        BsdName             : "gif0",
        Address             : nil,
    },
    {
        BsdName             : "stf0",
    },
}

var gateways0 = []*HostNetworkGateway{
    {
        Family              : 2,
        IsDefault           : true,
        IfaceName           : "en0",
        Address             : "192.168.1.1",
    },
}


func setup() (*hostContext) {
    ctx := &hostContext{}
    initializeHostContext(ctx)
    _context = ctx
    return ctx
}

func teardown() {
    _context = nil
}


func TestSearchPrimaryIPCandidate(t *testing.T) {
    setup()
    defer teardown()

    singletonContextInstance().monitorNetworkInterfaces(intefaces0)

    ip, err := HostPrimaryIPAddress()
    if err != nil {
        t.Error(err.Error())
    }

    if ip != "192.168.1.248" {
        t.Error("[ERR] wrong ip address has selected! It's supposed to be 192.168.1.248")
    }
}