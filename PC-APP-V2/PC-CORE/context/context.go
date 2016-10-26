//go:binary-only-package
package context

import (
    "sync"
    "fmt"
)

type HostContext interface {
    RefreshStatus() error

    CocoaHomeDirectory() (string, error)
    PosixHomeDirectory() (string, error)
    FullUserName() (string, error)
    LoginUserName() (string, error)
    UserTemporaryDirectory() (string, error)

    ApplicationSupportDirectory() (string, error)
    ApplicationDocumentsDirectoru() (string, error)
    ApplicationTemporaryDirectory() (string, error)
    ApplicationLibraryCacheDirectory() (string, error)
    ApplicationResourceDirectory() (string, error)
    ApplicationExecutableDirectory() (string, error)

    HostDeviceSerial() (string, error)

    HostPrimaryAddress() (string, error)
    HostDefaultGatewayAddress() (string, error)
}

type hostContext struct {
    hostInterfaces              *[]*HostNetworkInterface
    hostGateways                *[]*HostNetworkGateway

    primaryInteface             *HostNetworkInterface
    primaryAddress              *HostIPAddress
    primaryGateway              *HostNetworkGateway

    cocoaHomePath               string
    posixHomePath               string
    fullUserName                string
    loginUserName               string
    userTempPath                string

    applicationSupportPath      string
    applicationDocumentPath     string
    applicationTempPath         string
    applicationLibCachePath     string
    applicationResourcePath     string
    applicationExecutablePath   string

    hostDeviceSerial            string
}

// singleton initialization
var context *hostContext = nil
var once sync.Once

func SharedHostContext() (HostContext) {
    return singletonContextInstance()
}

func singletonContextInstance() (*hostContext) {
    once.Do(func() {
        context = &hostContext{}
        initializeHostContext(context)
    })
    return context
}

func initializeHostContext(ctx *hostContext) {
}

// take network interfaces
func (ctx *hostContext) monitorNetworkInterfaces(interfaces []*HostNetworkInterface) {
    // TODO : we make an assumption that host's primary interface and network addresses are at the same network segment. This could not be the case, we'll look into it v0.1.5
    ctx.hostInterfaces = &interfaces

    for _, iface := range interfaces {
        if iface.IsPrimary {
            ctx.primaryInteface = iface

            for _, addr := range iface.Address {
                if addr.IsPrimary {
                    ctx.primaryAddress = addr
                }
            }
            break
        }
    }

    // this is backup. It selects 1) Wi-Fi interface with 2) an active ip address.
    // Ethernet is going to be obsolete on lots of new macbook. We'll take wifi as default
    if ctx.primaryInteface == nil {
        for _, iface := range interfaces {
            if iface.IsActive && (iface.MediaType == "IEEE80211" || iface.DisplayName == "Wi-Fi") {
                ctx.primaryInteface = iface
                ctx.primaryAddress = iface.Address[0]
                break
            }
        }
    }
}

func (ctx *hostContext) monitorNetworkGateways(gateways []*HostNetworkGateway) {
    for _, gw := range gateways {
        if gw.IsDefault {
            ctx.primaryGateway = gw
        }
    }
    return
}

func (ctx *hostContext) RefreshStatus() error {

    ctx.cocoaHomePath               = findCocoaHomeDirectory()
    ctx.posixHomePath               = findPosixHomeDirectory()
    ctx.fullUserName                = findFullUserName()
    ctx.loginUserName               = findLoginUserName()
    ctx.userTempPath                = findUserTemporaryDirectory()

    ctx.applicationSupportPath      = findApplicationSupportDirectory()
    ctx.applicationDocumentPath     = findApplicationDocumentsDirectoru()
    ctx.applicationTempPath         = findApplicationTemporaryDirectory()
    ctx.applicationLibCachePath     = findApplicationLibraryCacheDirectory()
    ctx.applicationResourcePath     = findApplicationResourceDirectory()
    ctx.applicationExecutablePath   = findApplicationExecutableDirectory()

    ctx.hostDeviceSerial            = findSerialNumber()
    return nil
}

func (ctx *hostContext) CocoaHomeDirectory() (string, error) {
    if len(ctx.cocoaHomePath) == 0 {
        return "", fmt.Errorf("[ERR] Invalid Cocoa Home Directory")
    }
    return ctx.cocoaHomePath, nil
}

func (ctx *hostContext) PosixHomeDirectory() (string, error) {
    if len(ctx.posixHomePath) == 0 {
        return "", fmt.Errorf("[ERR] Invalid Posix Home Directory")
    }
    return ctx.posixHomePath, nil
}

func (ctx *hostContext) FullUserName() (string, error) {
    if len(ctx.fullUserName) == 0 {
        return "", fmt.Errorf("[ERR] Invalid Full Username")
    }
    return ctx.fullUserName, nil
}

func (ctx *hostContext) LoginUserName() (string, error) {
    if len(ctx.loginUserName) == 0 {
        return "", fmt.Errorf("[ERR] Invalid Login user name")
    }
    return ctx.loginUserName, nil
}

func (ctx *hostContext) UserTemporaryDirectory() (string, error) {
    if len(ctx.userTempPath) == 0 {
        return "", fmt.Errorf("[ERR] Invalid user temp directory")
    }
    return ctx.userTempPath, nil
}

func (ctx *hostContext) ApplicationSupportDirectory() (string, error) {
    if len(ctx.applicationSupportPath) == 0 {
        return "", fmt.Errorf("[ERR] Invalid App support directory")
    }
    return ctx.applicationSupportPath, nil
}

func (ctx *hostContext) ApplicationDocumentsDirectoru() (string, error) {
    if len(ctx.applicationDocumentPath) == 0 {
        return "", fmt.Errorf("[ERR] Invalid App doc directory")
    }
    return ctx.applicationDocumentPath, nil
}

func (ctx *hostContext) ApplicationTemporaryDirectory() (string, error) {
    if len(ctx.applicationTempPath) == 0 {
        return "", fmt.Errorf("[ERR] Invalid App temp directory")
    }
    return ctx.applicationTempPath, nil
}

func (ctx *hostContext) ApplicationLibraryCacheDirectory() (string, error) {
    if len(ctx.applicationLibCachePath) == 0 {
        return "", fmt.Errorf("[ERR] Invalid App lib cache directory")
    }
    return ctx.applicationLibCachePath, nil
}

func (ctx *hostContext) ApplicationResourceDirectory() (string, error) {
    if len(ctx.applicationResourcePath) == 0 {
        return "", fmt.Errorf("[ERR] Invalid app resource directory")
    }
    return ctx.applicationResourcePath, nil
}

func (ctx *hostContext) ApplicationExecutableDirectory() (string, error) {
    if len(ctx.applicationExecutablePath) == 0 {
        return "", fmt.Errorf("[ERR] Invalid app exec directory")
    }
    return ctx.applicationExecutablePath, nil
}

func (ctx *hostContext) HostDeviceSerial() (string, error) {
    if len(ctx.hostDeviceSerial) == 0 {
        return "", fmt.Errorf("[ERR] Invalid host device serial")
    }
    return ctx.hostDeviceSerial, nil
}

func (ctx *hostContext) HostPrimaryAddress() (string, error) {
    addr := ctx.primaryAddress
    if addr != nil {
        return addr.Address, nil
    }

    return "", fmt.Errorf("[ERR] No primary address has been found")
}

func (ctx *hostContext) HostDefaultGatewayAddress() (string, error) {
    gateway := ctx.primaryGateway
    if gateway != nil {
        return gateway.Address, nil
    }

    return "", fmt.Errorf("[ERR] No default gateway is found")
}