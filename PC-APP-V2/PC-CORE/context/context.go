//go:binary-only-package
package context

import (
    "os"
    "sync"

    "github.com/ricochet2200/go-disk-usage/du"
    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
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
    ApplicationUserDataDirectory() (string, error)

    HostPrimaryAddress() (string, error)
    HostDefaultGatewayAddress() (string, error)

    HostProcessorCount() uint
    HostActiveProcessorCount() uint
    HostPhysicalMemorySize() uint64
    HostStorageSpaceStatus() (total uint64, available uint64)

    CurrentCountryCode() (string, error)
    CurrentLanguageCode() (string, error)
    MasterAgentName() (string, error)
    SetMasterAgentName(man string)

    // cert authority
    CertAuthSigner() (*pcrypto.CaSigner, error)
    CertAuthPublicKey() ([]byte, error)
    // host certificate
    MasterHostPublicKey() ([]byte, error)
    MasterHostPrivateKey() ([]byte, error)
    MasterHostCertificate() ([]byte, error)
}

type hostContext struct {
    sync.Mutex

    hostInterfaces               *[]*HostNetworkInterface
    hostGateways                 *[]*HostNetworkGateway

    primaryInteface              *HostNetworkInterface
    primaryAddress               *HostIPAddress
    primaryGateway               *HostNetworkGateway

    cocoaHomePath                string
    posixHomePath                string
    fullUserName                 string
    loginUserName                string
    userTempPath                 string

    applicationSupportPath       string
    applicationDocumentPath      string
    applicationTempPath          string
    applicationLibCachePath      string
    applicationResourcePath      string
    applicationExecutablePath    string

    processorCount               uint
    activeProcessorCount         uint
    physicalMemorySize           uint64
    clusterPublicName            string

    currentCountryCode           string
    currentLanguageCode          string

    // certificate authority
    *pcrypto.CaSigner
    caPrivateKey                 []byte
    caPublicKey                  []byte
    caCertificate                []byte
    caSSHChecker                 []byte

    // host certificate
    hostPrivateKey               []byte
    hostPublicKey                []byte
    hostSshKey                   []byte
    hostCertifcate               []byte
}

// singleton initialization
var _context *hostContext = nil
var _once sync.Once

func SharedHostContext() (HostContext) {
    return singletonContextInstance()
}

func singletonContextInstance() (*hostContext) {
    _once.Do(func() {
        _context = &hostContext{}
        _context.RefreshStatus()
    })
    return _context
}

// take network interfaces
func MonitorNetworkInterfaces(interfaces []*HostNetworkInterface) {
    singletonContextInstance().refreshNetworkInterfaces(interfaces)
}

func (ctx *hostContext) refreshNetworkInterfaces(interfaces []*HostNetworkInterface) {
    ctx.Lock()
    defer ctx.Unlock()

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

func MonitorNetworkGateways(gateways []*HostNetworkGateway) {
    singletonContextInstance().refreshNetworkGateways(gateways)
}

func (ctx *hostContext) refreshNetworkGateways(gateways []*HostNetworkGateway) {
    ctx.Lock()
    defer ctx.Unlock()

    for _, gw := range gateways {
        if gw.IsDefault {
            ctx.primaryGateway = gw
        }
    }
    return
}

type CertAuthBundle struct {
    CASigner *pcrypto.CaSigner
    CAPrvKey []byte
    CAPubKey []byte
    CACrtPem []byte
    CASSHChk []byte
}

func UpdateCertAuth(bundle *CertAuthBundle) {
    ctx := singletonContextInstance()
    ctx.Lock()
    defer ctx.Unlock()

    ctx.CaSigner        = bundle.CASigner
    ctx.caPrivateKey    = bundle.CAPrvKey
    ctx.caPublicKey     = bundle.CAPubKey
    ctx.caCertificate   = bundle.CACrtPem
    ctx.caSSHChecker    = bundle.CASSHChk
}

type HostCertBundle struct {
    PrivateKey     []byte
    PublicKey      []byte
    SshKey         []byte
    Certificate    []byte
}

func UpdateHostCert(bundle *HostCertBundle) {
    ctx := singletonContextInstance()
    ctx.Lock()
    defer ctx.Unlock()

    ctx.hostPrivateKey  = bundle.PrivateKey
    ctx.hostPublicKey   = bundle.PublicKey
    ctx.hostSshKey      = bundle.SshKey
    ctx.hostCertifcate  = bundle.Certificate
}

func (ctx *hostContext) CocoaHomeDirectory() (string, error) {
    if len(ctx.cocoaHomePath) == 0 {
        return "", errors.Errorf("[ERR] Invalid Cocoa Home Directory")
    }
    return ctx.cocoaHomePath, nil
}

func (ctx *hostContext) PosixHomeDirectory() (string, error) {
    if len(ctx.posixHomePath) == 0 {
        return "", errors.Errorf("[ERR] Invalid Posix Home Directory")
    }
    return ctx.posixHomePath, nil
}

func (ctx *hostContext) FullUserName() (string, error) {
    if len(ctx.fullUserName) == 0 {
        return "", errors.Errorf("[ERR] Invalid Full Username")
    }
    return ctx.fullUserName, nil
}

func (ctx *hostContext) LoginUserName() (string, error) {
    if len(ctx.loginUserName) == 0 {
        return "", errors.Errorf("[ERR] Invalid Login user name")
    }
    return ctx.loginUserName, nil
}

func (ctx *hostContext) UserTemporaryDirectory() (string, error) {
    if len(ctx.userTempPath) == 0 {
        return "", errors.Errorf("[ERR] Invalid user temp directory")
    }
    return ctx.userTempPath, nil
}

func (ctx *hostContext) ApplicationSupportDirectory() (string, error) {
    if len(ctx.applicationSupportPath) == 0 {
        return "", errors.Errorf("[ERR] Invalid App support directory")
    }
    return ctx.applicationSupportPath, nil
}

func (ctx *hostContext) ApplicationDocumentsDirectoru() (string, error) {
    if len(ctx.applicationDocumentPath) == 0 {
        return "", errors.Errorf("[ERR] Invalid App doc directory")
    }
    return ctx.applicationDocumentPath, nil
}

func (ctx *hostContext) ApplicationTemporaryDirectory() (string, error) {
    if len(ctx.applicationTempPath) == 0 {
        return "", errors.Errorf("[ERR] Invalid App temp directory")
    }
    return ctx.applicationTempPath, nil
}

func (ctx *hostContext) ApplicationLibraryCacheDirectory() (string, error) {
    if len(ctx.applicationLibCachePath) == 0 {
        return "", errors.Errorf("[ERR] Invalid App lib cache directory")
    }
    return ctx.applicationLibCachePath, nil
}

func (ctx *hostContext) ApplicationResourceDirectory() (string, error) {
    if len(ctx.applicationResourcePath) == 0 {
        return "", errors.Errorf("[ERR] Invalid app resource directory")
    }
    return ctx.applicationResourcePath, nil
}

func (ctx *hostContext) ApplicationExecutableDirectory() (string, error) {
    if len(ctx.applicationExecutablePath) == 0 {
        return "", errors.Errorf("[ERR] Invalid app exec directory")
    }
    return ctx.applicationExecutablePath, nil
}

func (ctx *hostContext) ApplicationUserDataDirectory() (string, error) {
    pHome, err := ctx.PosixHomeDirectory()
    if err != nil {
        return "", err
    }
    dataPath := pHome + "/.pocket"

    // create the data directory if it's missing
    _, err = os.Stat(dataPath)
    if os.IsNotExist(err) {
        err := os.MkdirAll(dataPath, os.ModeDir|0700)
        if err != nil {
            return "", err
        }
    }

    return dataPath, nil
}

func (ctx *hostContext) HostPrimaryAddress() (string, error) {
    addr := ctx.primaryAddress
    if addr != nil {
        return addr.Address, nil
    }

    return "", errors.Errorf("[ERR] No primary address has been found")
}

func (ctx *hostContext) HostDefaultGatewayAddress() (string, error) {
    gateway := ctx.primaryGateway
    if gateway != nil {
        return gateway.Address, nil
    }

    return "", errors.Errorf("[ERR] No default gateway is found")
}

func (ctx *hostContext) HostProcessorCount() uint {
    return ctx.processorCount
}

func (ctx *hostContext) HostActiveProcessorCount() uint {
    return ctx.activeProcessorCount
}

func (ctx *hostContext) HostPhysicalMemorySize() uint64 {
    var MB = uint64(1024 * 1024)
    return uint64(ctx.physicalMemorySize / MB)
}

func (ctx *hostContext) HostStorageSpaceStatus() (total uint64, available uint64) {
    var MB = uint64(1024 * 1024)
    usage := du.NewDiskUsage("/")

/*
    fmt.Println("Free:", usage.Free()/(MB))
    fmt.Println("Available:", usage.Available()/(MB))
    fmt.Println("Size:", usage.Size()/(MB))
    fmt.Println("Used:", usage.Used()/(MB))
    fmt.Println("Usage:", usage.Usage()*100, "%")
*/

    total = uint64(usage.Size()/(MB))
    available = uint64(usage.Available()/(MB))
    return
}

//TODO : master specific identifier is necessary
func (ctx *hostContext) MasterAgentName() (string, error) {
    if len(ctx.clusterPublicName) == 0 {
        return "", errors.Errorf("[ERR] Invalid host device serial")
    }
    return ctx.clusterPublicName, nil
}

func (ctx *hostContext) SetMasterAgentName(man string) {
    ctx.clusterPublicName = man
}

func (ctx *hostContext) CurrentCountryCode() (string, error) {
    if len(ctx.currentCountryCode) == 0 {
        return "", errors.Errorf("[ERR] Invalid country code")
    }
    return ctx.currentCountryCode, nil
}

func (ctx *hostContext) CurrentLanguageCode() (string, error) {
    if len(ctx.currentLanguageCode) == 0 {
        return "", errors.Errorf("[ERR] Invalid language code")
    }
    return ctx.currentLanguageCode, nil
}

func (ctx *hostContext) CertAuthSigner() (*pcrypto.CaSigner, error) {
    if ctx.CaSigner == nil {
        return nil, errors.Errorf("[ERR] invalid cert authority signer")
    }
    return ctx.CaSigner, nil
}

func (ctx *hostContext) CertAuthPublicKey() ([]byte, error) {
    if ctx.caPublicKey == nil {
        return nil, errors.Errorf("[ERR] invalid cert public key")
    }
    return ctx.caPublicKey, nil
}

func (ctx *hostContext) MasterHostPublicKey() ([]byte, error) {
    if len(ctx.hostPublicKey) == 0 {
        return nil, errors.Errorf("[ERR] Invalid master public key")
    }
    return ctx.hostPublicKey, nil
}

func (ctx *hostContext) MasterHostPrivateKey() ([]byte, error) {
    if len(ctx.hostPrivateKey) == 0 {
        return nil, errors.Errorf("[ERR] Invalid master private key")
    }
    return ctx.hostPrivateKey, nil
}

func (ctx *hostContext) MasterHostCertificate() ([]byte, error) {
    if len(ctx.hostCertifcate) == 0 {
        return nil, errors.Errorf("[ERR] Invalid master certificate data")
    }
    return ctx.hostCertifcate, nil
}
