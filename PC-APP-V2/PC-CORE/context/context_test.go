package context

import "testing"

var intefaces0 = []*PCNetworkInterface {
    {
        WifiPowerOff        : false,
        IsActive            : true,
        IsPrimary           : true,
        AddrCount           : 3,
        Address             : []*SCNIAddress {
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

var gateways0 = []*SCNIGateway {
    {
        Family              : 2,
        IsDefault           : true,
        IfaceName           : "en0",
        Address             : "192.168.1.1",
    },
}

func TestSearchPrimaryIPCandidate(t *testing.T) {

}