package context

type SCNIAddress struct {
    Flags               uint
    Family              uint8
    IsPrimary           bool
    Addr                string
    Netmask             string
    Broadcast           string
    Peer                string
}

type PCNetworkInterface struct {
    WifiPowerOff        bool
    IsActive            bool
    IsPrimary           bool
    AddrCount           uint
    Address             []*SCNIAddress
    BsdName             string
    DisplayName         string
    MacAddress          string
    MediaType           string
}

type SCNIGateway struct {
    Family              uint8
    IsDefault           bool
    IfaceName           string
    Address             string
}

