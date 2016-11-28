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

/*
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
*/

// GenerateSelfSignedCert generates a self signed certificate that
// is valid for given domain names and ips, returns PEM-encoded bytes with key and cert
func generateSelfSignedCert(country string, hostNames []string) ([]byte, []byte, []byte, error) {
    var (
        privateKey *rsa.PrivateKey
        privDer, pubDer, certDer []byte
        privBlock, pubBlock, certBlock *pem.Block
        privPem, pubPem, certPem []byte
        notBefore, notAfter time.Time
        serialNumberLimit, serialNumber *big.Int
        err error = nil
    )

    // generate private key
    privateKey, err = rsa.GenerateKey(rand.Reader, rsaStrongKeySize)
    if err != nil {
        return nil, nil, nil, err
    }
    // check the key generated
    err = privateKey.Validate()
    if err != nil {
        return nil, nil, nil, err
    }
    // build private key
    privDer = x509.MarshalPKCS1PrivateKey(privateKey)
    privBlock = &pem.Block{
        Type:    "RSA PRIVATE KEY",
        Headers: nil,
        Bytes:   privDer,
    }
    privPem = pem.EncodeToMemory(privBlock)

    //// generate and public key
    pubDer, err = x509.MarshalPKIXPublicKey(privateKey.Public())
    if err != nil {
        return nil, nil, nil, err
    }
    pubBlock = &pem.Block{
        Type:   "PUBLIC KEY",
        Headers: nil,
        Bytes:   pubDer,
    }
    pubPem = pem.EncodeToMemory(pubBlock)

    //// generate certificate
    notBefore = time.Now()
    notAfter = notBefore.Add(time.Hour * 24 * 365 * 10) // 10 years
    serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)
    serialNumber, err = rand.Int(rand.Reader, serialNumberLimit)
    if err != nil {
        return nil, nil, nil, err
    }
    certEntity := pkix.Name{
        CommonName:   "localhost",
        Country:      []string{country},
        Organization: []string{"localhost"},
    }
    certTemplate := &x509.Certificate{
        SerialNumber:          serialNumber,
        Issuer:                certEntity,
        Subject:               certEntity,
        NotBefore:             notBefore,
        NotAfter:              notAfter,
        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
        BasicConstraintsValid: true,
        IsCA:                  true,
    }
    // collect IP addresses localhost resolves to and add them to the cert. template:
    certTemplate.DNSNames = append(hostNames, "localhost.local")

    // TODO :  (11/28/2016) we're to empty the addresses here to see if certificate works fine. Also, check if ip address is really necessary
    /*
    ips, _ := net.LookupIP("localhost")
    if ips != nil {
        certTemplate.IPAddresses = ips
    }
    */
    certTemplate.IPAddresses = []net.IP{}

    certDer, err = x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, privateKey.Public(), privateKey)
    if err != nil {
        return nil, nil, nil, err
    }
    certBlock = &pem.Block{
        Type: "CERTIFICATE",
        Bytes: certDer,
    }
    certPem = pem.EncodeToMemory(certBlock)

    return privPem, pubPem, certPem, nil
}

// CreateSelfSignedHTTPSCert generates and self-signs a TLS key+cert pair for https connection to the proxy server.
func GenerateSelfSignedCertificateFiles(pubKeyPath, prvKeyPath, certPath, country string) error {
    if len(country) != 2 {
        return fmt.Errorf("[ERR] Invalid country code")
    }

    prv, pub, cert, err := generateSelfSignedCert(country, []string{"localhost", "localhost"})
    if err != nil {
        return err
    }

    if len(pubKeyPath) != 0 && len(pub) != 0 {
        err = ioutil.WriteFile(pubKeyPath, pub, rsaKeyFilePerm)
        if err != nil {
            return err
        }
    }
    if len(prvKeyPath) != 0 && len(prv) != 0 {
        err = ioutil.WriteFile(prvKeyPath, prv, rsaKeyFilePerm)
        if err != nil {
            return err
        }
    }
    if len(certPath) != 0 && len(cert) != 0 {
        err = ioutil.WriteFile(certPath, cert, rsaKeyFilePerm)
        if err != nil {
            return err
        }
    }
    return nil
}
