package pcrypto

import (
    "crypto/rsa"
    "crypto/rand"
    "crypto/sha1"
    "crypto/x509/pkix"
    "crypto/x509"
    "encoding/asn1"
    "encoding/pem"
    "fmt"
    "io/ioutil"
    "math/big"
    "time"
    "strings"

    "golang.org/x/crypto/ssh"
)

type certError struct {
    s string
}

func (n *certError) Error() string {
    return n.s
}

const (
    pcCertAuthName string = "pc-cert-auth"
)

// ComputeSKI derives an SKI from the certificate's public key in a
// standard manner. This is done by computing the SHA-1 digest of the
// SubjectPublicKeyInfo component of the certificate.
func computeSKI(template *x509.Certificate) ([]byte, error) {
    encodedPub, err := x509.MarshalPKIXPublicKey(template.PublicKey)
    if err != nil {
        return nil, err
    }

    var subPKI subjectPublicKeyInfo
    _, err = asn1.Unmarshal(encodedPub, &subPKI)
    if err != nil {
        return nil, err
    }

    pubHash := sha1.Sum(subPKI.SubjectPublicKey.Bytes)
    return pubHash[:], nil
}

// makeSelfCertAuth generates a self signed certificate that
// is valid for given domain names and ips, returns PEM-encoded bytes with key and cert
func makeSelfCertAuth(commonName, dnsName, country string) ([]byte, []byte, []byte, []byte, error) {
    var (
        privateKey *rsa.PrivateKey
        privDer, pubDer, certDer []byte
        privBlock, pubBlock, certBlock *pem.Block
        privPem, pubPem, certPem []byte
        notBefore, notAfter time.Time
        serialNumber *big.Int
        ski []byte
        err error = nil
    )

    // check country code
    if len(commonName) == 0 {
        return nil, nil, nil, nil, fmt.Errorf("[ERR] Invalid common name")
    }
    if len(dnsName) == 0 {
        return nil, nil, nil, nil, fmt.Errorf("[ERR] Invalid dns name")
    }
    if len(country) == 0 {
        return nil, nil, nil, nil, fmt.Errorf("[ERR] Invalid country code")
    }

    // generate private key
    privateKey, err = rsa.GenerateKey(rand.Reader, rsaStrongKeySize)
    if err != nil {
        return nil, nil, nil, nil, err
    }
    // check the key generated
    err = privateKey.Validate()
    if err != nil {
        return nil, nil, nil, nil, err
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
        return nil, nil, nil, nil, err
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
    // we can use long serial number, but OpenSSL use 64bit long, which seems to be reasonable
    //serialNumber, err = rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
    //serialNumber, err = rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
    serialNumber, err = rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 64))
    if err != nil {
        return nil, nil, nil, nil, err
    }
    /*
    TODO : Subject Key ID -> Authority key id. There must be a link between. We'll handle this later
    subjectKeyId, err = rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
    if err != nil {
        return nil, nil, nil, err
    }
    */
    certEntity := pkix.Name{
        CommonName:   commonName,
        Country:      []string{country},
    }
    certTemplate := &x509.Certificate{
        SignatureAlgorithm:    x509.SHA256WithRSA,
        PublicKeyAlgorithm:    x509.RSA,
        PublicKey:             privateKey.Public(),
        SerialNumber:          serialNumber,
        Issuer:                certEntity,
        Subject:               certEntity,
        NotBefore:             notBefore,
        NotAfter:              notAfter,
        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
        BasicConstraintsValid: true,
        IsCA:                  true,
        DNSNames:              []string{dnsName, "localhost.local"},
    }

    // TODO : (11/28/2016) we're to empty the addresses here to see if certificate works fine. Also, check if ip address is really necessary
    // TODO : (12/03/2016) it's ok to generate Cert Auth without IP. But we'll be using DomainName for tight security
/*
    // collect IP addresses localhost resolves to and add them to the cert. template:
    ips, _ := net.LookupIP("localhost")
    if ips != nil {
        certTemplate.IPAddresses = ips
    }
    certTemplate.IPAddresses = []net.IP{}
*/
    ski, err = computeSKI(certTemplate)
    if err != nil {
        return nil, nil, nil, nil, err
    }
    certTemplate.SubjectKeyId = ski
    certDer, err = x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, privateKey.Public(), privateKey)
    if err != nil {
        return nil, nil, nil, nil, err
    }
    certBlock = &pem.Block{
        Type: "CERTIFICATE",
        Bytes: certDer,
    }
    certPem = pem.EncodeToMemory(certBlock)

    // ssh checking key. added for teleport
    sshpub, err := ssh.NewPublicKey(privateKey.Public())
    if err != nil {
        return nil, nil, nil, nil, err
    }
    sshpubBytes := ssh.MarshalAuthorizedKey(sshpub)
    return pubPem, privPem, certPem, sshpubBytes, nil
}

// TODO : Add Test
// CreateSelfSignedHTTPSCert generates and self-signs a TLS key+cert pair for https connection to the proxy server.
func GenerateClusterCertificateAuthorityFiles(pubKeyPath, prvKeyPath, certPath, sshPath, domainName, country string) error {
    if len(domainName) == 0 {
        return &certError{"[ERR] bad cluster id"}
    }
    if len(country) == 0 {
        return &certError{"[ERR] bad country id"}
    }
    pub, prv, cert, sshbyte, err := makeSelfCertAuth(pcCertAuthName,
        fmt.Sprintf("%s.%s", pcCertAuthName, domainName),
        strings.ToUpper(country))
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
    if len(sshPath) != 0 && len(sshbyte) != 0 {
        err = ioutil.WriteFile(sshPath, sshbyte, rsaKeyFilePerm)
        if err != nil {
            return err
        }
    }
    return nil
}

// TODO : Add Test
// return : public, private, certificate, sshbyte, error
func GenerateClusterCertificateAuthorityData(domainName, country string) ([]byte, []byte, []byte, []byte, error) {
    if len(domainName) == 0 {
        return nil, nil, nil, nil, &certError{"[ERR] bad cluster id"}
    }
    if len(country) == 0 {
        return nil, nil, nil, nil, &certError{"[ERR] bad country id"}
    }
    return makeSelfCertAuth(pcCertAuthName,
        fmt.Sprintf("%s.%s",pcCertAuthName, domainName),
        strings.ToUpper(country))
}
