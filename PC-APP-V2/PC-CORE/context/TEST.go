//+build !build
package context

var test_intefaces = []*HostNetworkInterface{
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

var test_gateways = []*HostNetworkGateway{
    {
        Family              : 2,
        IsDefault           : true,
        IfaceName           : "en0",
        Address             : "192.168.1.1",
    },
}

func debugContextSetup() (*hostContext) {
    context = &hostContext{}
    initializeHostContext(context)
    return context
}

func debugContextTeardown() {
    context = nil
}

func DebugContextPrepared() (HostContext) {
    hostContext := singletonContextInstance()

    hostContext.cocoaHomePath               = "/Users/almightykim"
    hostContext.posixHomePath               = "/Users/almightykim"
    hostContext.fullUserName                = "Almighty Kim"
    hostContext.loginUserName               = "almightykim"
    hostContext.userTempPath                = "/var/folders/1s/nn_7b2vd75g6lfs5_vxcgt_c0000gn/T/"

    hostContext.applicationSupportPath      = "/Users/almightykim/Library/Application Support/SysUtil"
    hostContext.applicationDocumentPath     = "/Users/almightykim/Documents"
    hostContext.applicationTempPath         = "/var/folders/1s/nn_7b2vd75g6lfs5_vxcgt_c0000gn/T/"
    hostContext.applicationLibCachePath     = "/Users/almightykim/Library/Caches"
    hostContext.applicationResourcePath     = "/Users/almightykim/Library/Developer/Xcode/DerivedData/SysUtil-dsrzjqwmorphavfrktsexyevvird/Build/Products/Debug/SysUtil.app/Contents/Resources"
    hostContext.applicationExecutablePath   = "/Users/almightykim/Library/Developer/Xcode/DerivedData/SysUtil-dsrzjqwmorphavfrktsexyevvird/Build/Products/Debug/SysUtil.app/Contents/MacOS/SysUtil"

    hostContext.hostDeviceSerial            = "G8815052XYL"

    hostContext.monitorNetworkGateways(test_gateways)
    hostContext.monitorNetworkInterfaces(test_intefaces)

    return hostContext
}

func DebugContextDestroyed() {
}