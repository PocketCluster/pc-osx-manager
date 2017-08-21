package context

import (
    "sync"

    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
)

type HostContextCertificate interface {
    // cert authority
    UpdateCertAuth(bundle *CertAuthBundle)
    CertAuthSigner() (*pcrypto.CaSigner, error)
    CertAuthPublicKey() ([]byte, error)
    CertAuthCertificate() ([]byte, error)

    // host certificate
    UpdateHostCert(bundle *HostCertBundle)
    MasterHostPublicKey() ([]byte, error)
    MasterHostPrivateKey() ([]byte, error)
    MasterHostCertificate() ([]byte, error)

    // beacon certificate
    UpdateBeaconCert(bundle *BeaconCertBundle)
    MasterBeaconPublicKey() ([]byte, error)
    MasterBeaconPrivateKey() ([]byte, error)

    // vbox certificate
    UpdateVBoxCert(bundle *VBoxCertBundle)
    MasterVBoxCtrlPrivateKey() ([]byte, error)
    MasterVBoxCtrlPublicKey() ([]byte, error)
}

type hostCertificate struct {
    sync.Mutex

    // certificate authority
    caBundle                     *CertAuthBundle
    // host certificate
    hostBundle                   *HostCertBundle
    // beacon certificate
    beaconBundle                 *BeaconCertBundle
    // VBox Control Certificate
    vboxBundle                   *VBoxCertBundle
}

// --- Certificate Authority Handling --- //
type CertAuthBundle struct {
    CASigner *pcrypto.CaSigner
    CAPrvKey []byte
    CAPubKey []byte
    CACrtPem []byte
    CASSHChk []byte
}

func (ctx *hostCertificate) UpdateCertAuth(bundle *CertAuthBundle) {
    ctx.Lock()
    defer ctx.Unlock()

    ctx.caBundle = bundle
}

func (ctx *hostCertificate) CertAuthSigner() (*pcrypto.CaSigner, error) {
    ctx.Lock()
    defer ctx.Unlock()

    if ctx.caBundle == nil || ctx.caBundle.CASigner == nil {
        return nil, errors.Errorf("[ERR] invalid cert authority signer")
    }
    return ctx.caBundle.CASigner, nil
}

func (ctx *hostCertificate) CertAuthPublicKey() ([]byte, error) {
    ctx.Lock()
    defer ctx.Unlock()

    if ctx.caBundle == nil || ctx.caBundle.CAPubKey == nil {
        return nil, errors.Errorf("[ERR] invalid cert authority public key")
    }
    return ctx.caBundle.CAPubKey, nil
}

func (ctx *hostCertificate) CertAuthCertificate() ([]byte, error) {
    ctx.Lock()
    defer ctx.Unlock()

    if ctx.caBundle == nil || ctx.caBundle.CACrtPem == nil {
        return nil, errors.Errorf("[ERR] invalid cert authority certificate")
    }
    return ctx.caBundle.CACrtPem, nil
}

// --- Host Certificate Handling --- //
type HostCertBundle struct {
    PrivateKey     []byte
    PublicKey      []byte
    SshKey         []byte
    Certificate    []byte
}

func (ctx *hostCertificate) UpdateHostCert(bundle *HostCertBundle) {
    ctx.Lock()
    defer ctx.Unlock()

    ctx.hostBundle = bundle
}

func (ctx *hostCertificate) MasterHostPublicKey() ([]byte, error) {
    ctx.Lock()
    defer ctx.Unlock()

    if ctx.hostBundle == nil || ctx.hostBundle.PublicKey == nil {
        return nil, errors.Errorf("[ERR] Invalid master public key")
    }
    return ctx.hostBundle.PublicKey, nil
}

func (ctx *hostCertificate) MasterHostPrivateKey() ([]byte, error) {
    ctx.Lock()
    defer ctx.Unlock()

    if ctx.hostBundle == nil || ctx.hostBundle.PrivateKey == nil {
        return nil, errors.Errorf("[ERR] Invalid master private key")
    }
    return ctx.hostBundle.PrivateKey, nil
}

func (ctx *hostCertificate) MasterHostCertificate() ([]byte, error) {
    ctx.Lock()
    defer ctx.Unlock()

    if ctx.hostBundle == nil || ctx.hostBundle.Certificate == nil {
        return nil, errors.Errorf("[ERR] Invalid master certificate data")
    }
    return ctx.hostBundle.Certificate, nil
}

// --- Beacon Certificate Handling --- //
type BeaconCertBundle struct {
    PrivateKey     []byte
    PublicKey      []byte
}

func (ctx *hostCertificate) UpdateBeaconCert(bundle *BeaconCertBundle) {
    ctx.Lock()
    defer ctx.Unlock()

    ctx.beaconBundle = bundle
}

func (ctx *hostCertificate) MasterBeaconPublicKey() ([]byte, error) {
    ctx.Lock()
    defer ctx.Unlock()

    if ctx.beaconBundle == nil || ctx.beaconBundle.PublicKey == nil {
        return nil, errors.Errorf("[ERR] invalid public beacon key")
    }
    return ctx.beaconBundle.PublicKey, nil
}

func (ctx *hostCertificate) MasterBeaconPrivateKey() ([]byte, error) {
    ctx.Lock()
    defer ctx.Unlock()

    if ctx.beaconBundle == nil || ctx.beaconBundle.PrivateKey == nil {
        return nil, errors.Errorf("[ERR] invalid private beacon key")
    }
    return ctx.beaconBundle.PrivateKey, nil
}

type VBoxCertBundle struct {
    PrivateKey     []byte
    PublicKey      []byte
}

func (ctx *hostCertificate) UpdateVBoxCert(bundle *VBoxCertBundle) {
    ctx.Lock()
    defer ctx.Unlock()

    ctx.vboxBundle = bundle
}

func (ctx *hostCertificate) MasterVBoxCtrlPrivateKey() ([]byte, error) {
    ctx.Lock()
    defer ctx.Unlock()

    if ctx.vboxBundle == nil || ctx.vboxBundle.PrivateKey == nil {
        return nil, errors.Errorf("[ERR] invalid vbox private key")
    }
    return ctx.vboxBundle.PrivateKey, nil
}

func (ctx *hostCertificate) MasterVBoxCtrlPublicKey() ([]byte, error) {
    ctx.Lock()
    defer ctx.Unlock()

    if ctx.vboxBundle == nil || ctx.vboxBundle.PublicKey == nil {
        return nil, errors.Errorf("[ERR] invalid vbox public key")
    }
    return ctx.vboxBundle.PublicKey, nil
}
