package context

import "github.com/stkim1/pcrypto"

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

func DebugContextPrepare() (HostContext) {
    // once singleton is assigned, it will not assign again. This is how we invalidate singleton ops
    singletonContextInstance()
    context = &hostContext{
        cocoaHomePath:               "/Users/almightykim",
        posixHomePath:               "/Users/almightykim",
        fullUserName:                "Almighty Kim",
        loginUserName:               "almightykim",
        userTempPath:                "/var/folders/1s/nn_7b2vd75g6lfs5_vxcgt_c0000gn/T/",

        applicationSupportPath:      "/Users/almightykim/Library/Application Support/SysUtil",
        applicationDocumentPath:     "/Users/almightykim/Documents",
        applicationTempPath:         "/var/folders/1s/nn_7b2vd75g6lfs5_vxcgt_c0000gn/T/",
        applicationLibCachePath:     "/Users/almightykim/Library/Caches",
        applicationResourcePath:     "/Users/almightykim/Library/Developer/Xcode/DerivedData/SysUtil-dsrzjqwmorphavfrktsexyevvird/Build/Products/Debug/SysUtil.app/Contents/Resources",
        applicationExecutablePath:   "/Users/almightykim/Library/Developer/Xcode/DerivedData/SysUtil-dsrzjqwmorphavfrktsexyevvird/Build/Products/Debug/SysUtil.app/Contents/MacOS/SysUtil",

        hostDeviceSerial:            "G8815052XYL",

        processorCount:              8,
        activeProcessorCount:        8,
        physicalMemorySize:          34359738368,

        publicKeyData:               pcrypto.TestMasterPublicKey(),
        privateKeyData:              pcrypto.TestMasterPrivateKey(),

        currentLanguageCode:         "EN",
        currentCountryCode:          "KR",
    }

    context.monitorNetworkGateways(test_gateways)
    context.monitorNetworkInterfaces(test_intefaces)

    return context
}

func DebugContextDestroy() {
    debugContextTeardown()
}