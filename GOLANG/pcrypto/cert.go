package pcrypto

import (
    "crypto/tls"
    "crypto/rsa"
    "crypto/rand"
    "crypto/x509/pkix"
    "crypto/x509"
    "encoding/pem"
    "io/ioutil"
    "os"
    "time"
    "math/big"
    "fmt"
    "net"
    "path/filepath"
)

const (
    // DefaultLRUCapacity is a capacity for LRU session cache
    DefaultLRUCapacity = 1024
    // DefaultCertTTL sets the TTL of the self-signed certificate (1 year)
    DefaultCertTTL = (24 * time.Hour) * 365
)

// CreateTLSConfiguration sets up default TLS configuration
func CreateTLSConfiguration(certFile, keyFile string) (*tls.Config, error) {
    config := &tls.Config{}

    if _, err := os.Stat(certFile); err != nil {
        return nil, fmt.Errorf("[ERR] certificate is not accessible by '%v', %s", certFile, err.Error())
    }
    if _, err := os.Stat(keyFile); err != nil {
        return nil, fmt.Errorf("[ERR] certificate is not accessible by '%v', %s", certFile, err.Error())
    }

    //log.Infof("[PROXY] TLS cert=%v key=%v", certFile, keyFile)
    cert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return nil, err
    }

    config.Certificates = []tls.Certificate{cert}

    config.CipherSuites = []uint16{
        tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

        tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
        tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,

        tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
        tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,

        tls.TLS_RSA_WITH_AES_256_CBC_SHA,
        tls.TLS_RSA_WITH_AES_128_CBC_SHA,
    }

    config.MinVersion = tls.VersionTLS12
    config.SessionTicketsDisabled = false
    config.ClientSessionCache = tls.NewLRUClientSessionCache(
        DefaultLRUCapacity)

    return config, nil
}

// TLSCredentials keeps the typical 3 components of a proper HTTPS configuration
type TLSCredentials struct {
    // PublicKey in PEM format
    PublicKey []byte
    // PrivateKey in PEM format
    PrivateKey []byte
    Cert       []byte
}

// GenerateSelfSignedCert generates a self signed certificate that
// is valid for given domain names and ips, returns PEM-encoded bytes with key and cert
func generateSelfSignedCert(country string, hostNames []string) (*TLSCredentials, error) {
    priv, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, err
    }
    notBefore := time.Now()
    notAfter := notBefore.Add(time.Hour * 24 * 365 * 10) // 10 years

    serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
    serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
    if err != nil {
        return nil, err
    }

    entity := pkix.Name{
        CommonName:   "localhost",
        Country:      []string{country},
        Organization: []string{"localhost"},
    }

    template := x509.Certificate{
        SerialNumber:          serialNumber,
        Issuer:                entity,
        Subject:               entity,
        NotBefore:             notBefore,
        NotAfter:              notAfter,
        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
        BasicConstraintsValid: true,
        IsCA:                  true,
    }

    // collect IP addresses localhost resolves to and add them to the cert. template:
    template.DNSNames = append(hostNames, "localhost.local")
    ips, _ := net.LookupIP("localhost")
    if ips != nil {
        template.IPAddresses = ips
    }
    derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
    if err != nil {
        return nil, err
    }

    publicKeyBytes, err := x509.MarshalPKIXPublicKey(priv.Public())
    if err != nil {
        return nil, err
    }

    return &TLSCredentials{
        PublicKey:  pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyBytes}),
        PrivateKey: pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)}),
        Cert:       pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes}),
    }, nil
}

// CreateSelfSignedHTTPSCert generates and self-signs a TLS key+cert pair for https connection to the proxy server.
func CreateSelfSignedHTTPSCert(selfSignedKeyPath, selfSignedCertPath, certDirPath, country string) error {
    // if cert save path does not exist, return error
    if _, err := os.Stat(certDirPath); err != nil {
        if os.IsNotExist(err) {
            return err
        }
    }
    if len(selfSignedKeyPath) == 0 {
        return fmt.Errorf("[ERR] Invalid SelfSignedKeyPath")
    }
    if len(selfSignedCertPath) == 0 {
        return fmt.Errorf("[ERR] Invalid SelfSignedCertPath")
    }
    if len(country) != 2 {
        return fmt.Errorf("[ERR] Invalid country code")
    }

    keyPath := filepath.Join(certDirPath, selfSignedKeyPath)
    certPath := filepath.Join(certDirPath, selfSignedCertPath)

    // return the existing pair if they ahve already been generated:
    _, err := tls.LoadX509KeyPair(certPath, keyPath)
    if err == nil {
        return nil
    }
    if !os.IsNotExist(err) {
        return err
    }

    // "[CONFIG] Generating self signed key and cert to %v %v", keyPath, certPath)
    creds, err := generateSelfSignedCert(country, []string{"localhost", "localhost"})
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(keyPath, creds.PrivateKey, 0600)
    if err != nil {
        return err
    }
    err = ioutil.WriteFile(certPath, creds.Cert, 0600)
    if err != nil {
        return err
    }
    return nil
}

