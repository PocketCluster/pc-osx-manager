package main

import (
    "fmt"
    "time"
    "crypto/x509"
    "encoding/pem"
    "bytes"
    "crypto/rand"
    "math/big"
    "encoding/asn1"
    "crypto/x509/pkix"
    "crypto/rsa"
    "log"

    "github.com/gravitational/teleport"
    "github.com/gravitational/teleport/lib/auth/native"
    "github.com/gravitational/teleport/lib/services"
    "github.com/gravitational/teleport/lib/auth"

    "github.com/cloudflare/cfssl/csr"
    "github.com/cloudflare/cfssl/helpers"
    "github.com/cloudflare/cfssl/config"
    "github.com/cloudflare/cfssl/selfsign"
    "github.com/cloudflare/cfssl/cli/genkey"

    "github.com/smira/go-uuid/uuid"
)

// *** this is self sign and certificate request at the same time. ***
func CreateCertificateAuthorityAndRequest() ([]byte, []byte, []byte, error) {
    names := pkix.Name{
        CommonName:   "localhost",
        Country:      []string{"US"},
        Organization: []string{"localhost"},
    }

    // step: generate a keypair
    keys, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("unable to genarate private keys, error: %s", err)
    }
    var privateKey bytes.Buffer
    if err := pem.Encode(&privateKey, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(keys)}); err != nil {
        return nil, nil, nil, err
    }

    // step: create the request template
    // step: generate a serial number
    serial, err := rand.Int(rand.Reader, (&big.Int{}).Exp(big.NewInt(2), big.NewInt(159), nil))
    if err != nil {
        return nil, nil, nil, err
    }
    now := time.Now()
    //// generate certificate
    template := x509.Certificate{
        SerialNumber:          serial,
        Subject:               names,
        NotBefore:             now.Add(-10 * time.Minute).UTC(),
        NotAfter:              now.Add(time.Hour * 24 * 365 * 10).UTC(),
        BasicConstraintsValid: true,
        IsCA:                  true,
        KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
    }
    // step: sign the certificate authority
    certificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &keys.PublicKey, keys)
    if err != nil {
        return nil, nil, nil, fmt.Errorf("failed to generate certificate, error: %s", err)
    }
    var certBuf bytes.Buffer
    if err := pem.Encode(&certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: certificate}); err != nil {
        return nil, nil, nil, err
    }
/*
    return &caCertificate{
        privateKey: privateKey.String(),
        publicKey:  request.String(),
        csr:        string(csr),
    }, nil
*/

    // step: generate the csr request
    val, err := asn1.Marshal(basicConstraints{false, 0})
    if err != nil {
        return nil, nil, nil, err
    }
    // step: generate a csr template
    var csrTemplate = x509.CertificateRequest{
        Subject:            names,
        SignatureAlgorithm: x509.SHA512WithRSA,
        ExtraExtensions: []pkix.Extension{
            {
                Id:       asn1.ObjectIdentifier{2, 5, 29, 19},
                Value:    val,
                Critical: true,
            },
        },
    }
    csrCertificate, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, keys)
    if err != nil {
        return nil, nil, nil, err
    }
    csrReq := pem.EncodeToMemory(&pem.Block{
        Type: "CERTIFICATE REQUEST", Bytes: csrCertificate,
    })

    for _, val := range csrTemplate.ExtraExtensions {
        // Check the CSR for the X.509 BasicConstraints (RFC 5280, 4.2.1.9)
        // extension and append to template if necessary
        if val.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 19}) {
            var constraints csr.BasicConstraints
            var rest []byte

            if rest, err = asn1.Unmarshal(val.Value, &constraints); err != nil {
                return nil, nil, nil, err
            } else if len(rest) != 0 {
                return nil, nil, nil, err
            }
            log.Printf("[INFO] is Generated Template from Certificate Authority? %v", constraints)
        }
    }

    return privateKey.Bytes(), certBuf.Bytes(), csrReq, nil
}

// ---------------------------------------------------------------------------------------------------------------------
func GetHostCA(authority auth.Authority, DomainName string) (*services.CertAuthority, error) {
    priv, pub, err := authority.GenerateKeyPair("")
    if err != nil {
        return nil, err
    }
    return &services.CertAuthority{
        DomainName:   DomainName,
        Type:         services.HostCA,
        SigningKeys:  [][]byte{priv},
        CheckingKeys: [][]byte{pub},
    }, nil
}

func GenerateHostCert(authority auth.Authority, key []byte, hostID, authDomain string, roles teleport.Roles, ttl time.Duration) ([]byte, error) {
    ca, err := GetHostCA(authority, authDomain)
    if err != nil {
        return nil, err
    }
    privateKey, err := ca.FirstSigningKey()
    if err != nil {
        return nil, err
    }
    return authority.GenerateHostCert(privateKey, key, hostID, authDomain, roles, ttl)
}

func TestRegisterHostUsingToken() {
    DomainName := "localhost"
    hostID := uuid.New()
    authority := native.New()
    roles := []teleport.Role{teleport.RoleNode}

    k, pub, err := authority.GenerateKeyPair("")
    if err != nil {
        fmt.Print(err.Error())
        return
    }
    // we always append authority's domain to resulting node name,
    // that's how we make sure that nodes are uniquely identified/found
    // in cases when we have multiple environments/organizations
    fqdn := fmt.Sprintf("%s.%s", hostID, DomainName)
    c, err := GenerateHostCert(authority, pub, fqdn, DomainName, roles, 0)
    if err != nil {
        fmt.Print(err.Error())
    }

    fmt.Printf("\t *** PRIVATE KEY ***\n%s\n\t *** PUBLIC KEY ***\n%s\n\t *** CERTIFICATE ***\n%s", string(k), string(pub), string(c))
}

// ---------------------------------------------------------------------------------------------------------------------
func TestSelfSingingCertificate() {
    var req = &csr.CertificateRequest{
        CN:    "pc-master",
        Names: []csr.Name {
            {
                C: "KR",
            },
        },
        Hosts: []string{"pc-master"},
        KeyRequest: &csr.BasicKeyRequest{"rsa", 2048},
    }

    var key, csrPEM []byte
    g := &csr.Generator{Validator: genkey.Validator}
    csrPEM, key, err := g.ProcessRequest(req)
    if err != nil {
        key = nil
        return
    }

    priv, err := helpers.ParsePrivateKeyPEM(key)
    if err != nil {
        key = nil
        return
    }

    var profile *config.SigningProfile
    if profile == nil {
        profile = config.DefaultConfig()
        profile.Expiry = 2190 * time.Hour
    }

    signed, err := selfsign.Sign(priv, csrPEM, profile)
    if err != nil {
        key = nil
        priv = nil
        return
    }

    fmt.Printf("[SIGNED CSR]\n%s\n",string(signed))
}

// ---------------------------------------------------------------------------------------------------------------------
func TestCreateCertificateAuthority() {
    pk, ca, req, err := CreateCertificateAuthorityAndRequest()
    if err != nil {
        fmt.Printf(err.Error())
    }

    fmt.Printf("*** PRIVATE KEY ***\n%s\n\n*** CERTIFICATE ***\n%s\n\n*** REQUEST ***\n%s\n", string(pk), string(ca), string(req))

    p, rest := pem.Decode(pk)
    fmt.Printf("\n\n Block type %s \n\n Block Bytes %s \n\n Rest %s", p.Type, string(p.Bytes), string(rest))
}

func test() {
    TestCreateCertificateAuthority()
    fmt.Print("\n\n\n\t -*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-\n\n\n")

    TestRegisterHostUsingToken()
    fmt.Printf("\n\n\n\t -*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-\n\n\n")

    TestSelfSingingCertificate()
}