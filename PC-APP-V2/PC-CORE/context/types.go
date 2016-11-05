package context

type HostIPAddress struct {
    Flags               uint
    Family              uint8
    IsPrimary           bool
    Address             string
    Netmask             string
    Broadcast           string
    Peer                string
}

type HostNetworkInterface struct {
    WifiPowerOff        bool
    IsActive            bool
    IsPrimary           bool
    AddrCount           uint
    Address             []*HostIPAddress
    BsdName             string
    DisplayName         string
    MacAddress          string
    MediaType           string
}

type HostNetworkGateway struct {
    Family              uint8
    IsDefault           bool
    IfaceName           string
    Address             string
}

