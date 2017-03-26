package pcrypto

import (
    "crypto"
    "crypto/rand"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "encoding/asn1"
    "fmt"
    "io"
    "math/big"
    "net"
    "strings"
    "time"

    "github.com/cloudflare/cfssl/config"
    "github.com/stkim1/pcrypto/helpers"
)

type basicConstraints struct {
    IsCA       bool `asn1:"optional"`
    MaxPathLen int  `asn1:"optional,default:-1"`
}

// subName subName contains the SubjectInfo fields.
type subName struct {
    C            string // Country
    ST           string // State
    L            string // Locality
    O            string // OrganisationName
    OU           string // OrganisationalUnitName
    SerialNumber string
}

// requestSubject requestSubject contains the information that should be used to override the
// subject information when signing a certificate.
type requestSubject struct {
    CN           string
    Names        []subName
    SerialNumber string
}

// Name returns the PKIX name for the subject.
func (s *requestSubject) name() pkix.Name {
    var name pkix.Name
    name.CommonName = s.CN

    // appendIf appends to a if s is not an empty string.
    appendIf := func (s string, a *[]string) {
        if s != "" {
            *a = append(*a, s)
        }
    }
    for _, n := range s.Names {
        appendIf(n.C, &name.Country)
        appendIf(n.ST, &name.Province)
        appendIf(n.L, &name.Locality)
        appendIf(n.O, &name.Organization)
        appendIf(n.OU, &name.OrganizationalUnit)
    }
    name.SerialNumber = s.SerialNumber
    return name
}

// extension represents a raw extension to be included in the certificate.  The
// "value" field must be hex encoded.
type extension struct {
    ID       config.OID
    Critical bool
    Value    string
}

// signRequest stores a signature request, which contains the hostname,
// the CSR, optional subject information, and the signature profile.
//
// Extensions provided in the signRequest are copied into the certificate, as
// long as they are in the ExtensionWhitelist for the signer's policy.
// Extensions requested in the CSR are ignored, except for those processed by
// ParseCertificateRequest (mainly subjectAltName).
type signRequest struct {
    dnsName     []string
    ipAddress   []string
    Request     string
    Subject     *requestSubject
    Profile     string
    CRLOverride string
    Label       string
    Serial      *big.Int
    Extensions  []extension
}

type subjectPublicKeyInfo struct {
    Algorithm        pkix.AlgorithmIdentifier
    SubjectPublicKey asn1.BitString
}

func makeNodeCertificateRequest(privKeyPEM []byte, subject *requestSubject) ([]byte, error) {
    privKey, err := helpers.ParsePrivateKeyPEM(privKeyPEM)
    if err != nil {
        return nil, err
    }
    // step: generate the csr request
    val, err := asn1.Marshal(basicConstraints{IsCA:false, MaxPathLen:0})
    if err != nil {
        return nil, err
    }
    // step: generate a csr template
    var csrTemplate = x509.CertificateRequest{
        Subject:            subject.name(),
        SignatureAlgorithm: x509.SHA256WithRSA,
        ExtraExtensions: []pkix.Extension{
            {
                Id:       asn1.ObjectIdentifier{2, 5, 29, 19},
                Value:    val,
                Critical: true,
            },
        },
    }
    csrCertificate, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privKey)
    if err != nil {
        return nil, err
    }
    csrReq := pem.EncodeToMemory(&pem.Block{
        Type: "CERTIFICATE REQUEST", Bytes: csrCertificate,
    })
    return csrReq, nil
}

// ParseCertificateRequest takes an incoming certificate request and
// builds a certificate template from it.
func parseCertificateRequest(sigAlgo x509.SignatureAlgorithm, csrBytes []byte) (template *x509.Certificate, err error) {
    csrv, err := x509.ParseCertificateRequest(csrBytes)
    if err != nil {
        return nil, err
    }

    err = helpers.CheckSignature(csrv, csrv.SignatureAlgorithm, csrv.RawTBSCertificateRequest, csrv.Signature)
    if err != nil {
        return nil, err
    }

    template = &x509.Certificate{
        Subject:            csrv.Subject,
        PublicKeyAlgorithm: csrv.PublicKeyAlgorithm,
        PublicKey:          csrv.PublicKey,
        SignatureAlgorithm: sigAlgo,
        KeyUsage:           x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
        DNSNames:           csrv.DNSNames,
        IPAddresses:        csrv.IPAddresses,
        EmailAddresses:     csrv.EmailAddresses,
    }

    for _, val := range csrv.Extensions {
        // Check the CSR for the X.509 BasicConstraints (RFC 5280, 4.2.1.9)
        // extension and append to template if necessary
        if val.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 19}) {
            var constraints basicConstraints
            var rest []byte

            if rest, err = asn1.Unmarshal(val.Value, &constraints); err != nil {
                return nil, err
            } else if len(rest) != 0 {
                return nil, err
            }

            template.BasicConstraintsValid = true
            template.IsCA = constraints.IsCA
            template.MaxPathLen = constraints.MaxPathLen
            template.MaxPathLenZero = template.MaxPathLen == 0
        }
    }

    return
}

// PopulateSubjectFromCSR has functionality similar to Name, except
// it fills the fields of the resulting pkix.Name with req's if the
// subject's corresponding fields are empty
func populateSubjectFromCSR(subject *requestSubject, req pkix.Name) pkix.Name {
    // if no subject, use req
    if subject == nil {
        return req
    }

    name := subject.name()

    if name.CommonName == "" {
        name.CommonName = req.CommonName
    }
    // replaceSliceIfEmpty replaces the contents of replaced with newContents if
    // the slice referenced by replaced is empty
    replaceSliceIfEmpty := func (replaced, newContents *[]string) {
        if len(*replaced) == 0 {
            *replaced = *newContents
        }
    }
    replaceSliceIfEmpty(&name.Country, &req.Country)
    replaceSliceIfEmpty(&name.Province, &req.Province)
    replaceSliceIfEmpty(&name.Locality, &req.Locality)
    replaceSliceIfEmpty(&name.Organization, &req.Organization)
    replaceSliceIfEmpty(&name.OrganizationalUnit, &req.OrganizationalUnit)
    if name.SerialNumber == "" {
        name.SerialNumber = req.SerialNumber
    }
    return name
}

// FillTemplate is a utility function that tries to load as much of
// the certificate template as possible from the profiles and current
// template. It fills in the key uses, expiration, revocation URLs
// and SKI.
func fillTemplate(template *x509.Certificate, defaultProfile, profile *config.SigningProfile) error {
    var (
        ski             []byte
        eku             []x509.ExtKeyUsage
        ku              x509.KeyUsage
        backdate        time.Duration
        expiry          time.Duration
        notBefore       time.Time
        notAfter        time.Time
        crlURL, ocspURL string
        err error
    )
    ski, err = computeSKI(template)
    if err != nil {
        return err
    }
    // The third value returned from Usages is a list of unknown key usages.
    // This should be used when validating the profile at load, and isn't used
    // here.
    ku, eku, _ = profile.Usages()
    if ku == 0 && len(eku) == 0 {
        return fmt.Errorf("[ERR] PolicyError NoKeyUsages presents")
    }
    if expiry = profile.Expiry; expiry == 0 {
        expiry = defaultProfile.Expiry
    }
    if crlURL = profile.CRL; crlURL == "" {
        crlURL = defaultProfile.CRL
    }
    if ocspURL = profile.OCSP; ocspURL == "" {
        ocspURL = defaultProfile.OCSP
    }
    if backdate = profile.Backdate; backdate == 0 {
        backdate = -5 * time.Minute
    } else {
        backdate = -1 * profile.Backdate
    }
    if !profile.NotBefore.IsZero() {
        notBefore = profile.NotBefore.UTC()
    } else {
        notBefore = time.Now().Round(time.Minute).Add(backdate).UTC()
    }
    if !profile.NotAfter.IsZero() {
        notAfter = profile.NotAfter.UTC()
    } else {
        notAfter = notBefore.Add(expiry).UTC()
    }
    template.NotBefore = notBefore
    template.NotAfter = notAfter
    template.KeyUsage = ku
    template.ExtKeyUsage = eku
    template.BasicConstraintsValid = true
    template.IsCA = profile.CAConstraint.IsCA
    template.SubjectKeyId = ski

    return nil
}

func signRequestTemplate(s *CaSigner, template *x509.Certificate, profile *config.SigningProfile) (cert []byte, err error) {
    var distPoints = template.CRLDistributionPoints
    err = fillTemplate(template, s.policy.Default, profile)
    if distPoints != nil && len(distPoints) > 0 {
        template.CRLDistributionPoints = distPoints
    }
    if err != nil {
        return
    }

    derBytes, err := x509.CreateCertificate(rand.Reader, template, s.caCertPEM, template.PublicKey, s.caPrvKeyPEM)
    if err != nil {
        return nil, err
    }

    cert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
    //log.Printf("[INFO] signed certificate with serial number %v", template.SerialNumber)
    return
}

// Sign signs a new certificate based on the PEM-encoded client
// certificate or certificate request with the signing profile,
// specified by profileName.
func signCertificateRequest(signer *CaSigner, req signRequest) ([]byte, error) {
    block, _ := pem.Decode([]byte(req.Request))
    if block == nil {
        return nil, fmt.Errorf("[ERR] CSR decoding failed")
    }

    if block.Type != "NEW CERTIFICATE REQUEST" && block.Type != "CERTIFICATE REQUEST" {
        return nil, fmt.Errorf("[ERR] not a certificate or csr")
    }

    // Copy out only the fields from the CSR authorized by policy. If the profile contains no explicit whitelist, assume
    // that all fields should be copied from the CSR
    template, err := parseCertificateRequest(signer.sigAlgo, block.Bytes)
    if err != nil {
        return nil, err
    }

    // fills template's IPAddresses, EmailAddresses, and DNSNames with the
    // content of hosts, if it is not nil.
    if len(req.dnsName) != 0 {
        template.DNSNames = []string{}
        for i := range req.dnsName {
            template.DNSNames = append(template.DNSNames, req.dnsName[i])
        }
    }
    if len(req.ipAddress) != 0 {
        template.IPAddresses = []net.IP{}
        for i := range req.ipAddress {
            if ip := net.ParseIP(req.ipAddress[i]); ip != nil {
                template.IPAddresses = append(template.IPAddresses, ip)
            }
        }
    }
    template.Subject = populateSubjectFromCSR(req.Subject, template.Subject)

    // RFC 5280 4.1.2.2:
    // Certificate users MUST be able to handle serialNumber
    // values up to 20 octets.  Conforming CAs MUST NOT use
    // serialNumber values longer than 20 octets.
    //
    // If CFSSL is providing the serial numbers, it makes
    // sense to use the max supported size.
    serialNumber := make([]byte, 20)
    _, err = io.ReadFull(rand.Reader, serialNumber)
    if err != nil {
        return nil, err
    }
    // SetBytes interprets buf as the bytes of a big-endian
    // unsigned integer. The leading byte should be masked
    // off to ensure it isn't negative.
    serialNumber[0] &= 0x7F
    template.SerialNumber = new(big.Int).SetBytes(serialNumber)

    signedCert, err := signRequestTemplate(signer, template, signer.policy.Default)
    if err != nil {
        return nil, err
    }
    return signedCert, nil
}

type CaSigner struct {
    caPrvKeyPEM crypto.Signer
    caCertPEM   *x509.Certificate
    caCert      []byte
    sigAlgo     x509.SignatureAlgorithm
    policy      *config.Signing
    domainName  string
    country     string
}

func NewCertAuthoritySigner(caKey, caCert []byte, domainName, country string) (*CaSigner, error) {
    caPrvKeyPEM, err := helpers.ParsePrivateKeyPEM(caKey)
    if err != nil {
        return nil, err
    }
    caCertPEM, err := helpers.ParseCertificatePEM(caCert)
    if err != nil {
        return nil, err
    }
    policy := &config.Signing{
        Profiles: map[string]*config.SigningProfile{},
        Default:  config.DefaultConfig(),
    }
    if len(domainName) == 0 {
        return nil, &certError{"[ERR] Invalid cluster id"}
    }
    if len(country) == 0 {
        return nil, &certError{"[ERR] Invalid country code"}
    }
    return &CaSigner{
        caPrvKeyPEM:   caPrvKeyPEM,
        caCertPEM:     caCertPEM,
        caCert:        caCert,
        sigAlgo:       helpers.DefaultSigAlgo(caPrvKeyPEM),
        policy:        policy,
        domainName:    domainName,
        country:       strings.ToUpper(country),
    }, nil
}

// TODO : Add Test + Remove cfssl packages + check authority record to Database
// GenerateSignedCertificate returned a signed certificate
func (s *CaSigner) GenerateSignedCertificate(hostname, ipAddress string, privateKey []byte) ([]byte, error) {
    if len(hostname) == 0 {
        return nil, &certError{"[ERR] Invalid hostname"}
    }
    // TODO : we need to check if CA's cluster id the same as csr's clusteid
    subject := &requestSubject{
        CN:    hostname,
        Names: []subName{
            {
                C: s.country,
            },
        },
    }
    nodereq, err := makeNodeCertificateRequest(privateKey, subject)
    if err != nil {
        return nil, err
    }
    fqdn := fmt.Sprintf("%s.%s", strings.ToLower(hostname), s.domainName)
    creq := signRequest{
        dnsName:    []string{hostname, fqdn},
        Request:    string(nodereq),
        Subject:    subject,
    }
    if len(ipAddress) != 0 {
        creq.ipAddress = []string{ipAddress}
    }
    if !s.policy.Valid() {
        return nil, &certError{"[ERR] Invalid Certificate Authority Policy"}
    }
    return signCertificateRequest(s, creq)
}

func (s *CaSigner) CertificateAuthority() []byte {
    return s.caCert
}