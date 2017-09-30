package tlscfg

import (
    "crypto/tls"
    "crypto/x509"

    "github.com/pkg/errors"
)

// Load the TLS certificates/keys and, if verify is true, the CA.
func BuildTLSConfigWithCAcert(ca, cert, key []byte, verify bool) (*tls.Config, error) {
    c, err := tls.X509KeyPair(cert, key)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    config := &tls.Config{
        Certificates: []tls.Certificate{c},
        MinVersion:   tls.VersionTLS10,
    }

    if verify {
        certPool := x509.NewCertPool()
        certPool.AppendCertsFromPEM(ca)
        config.RootCAs = certPool
        config.ClientAuth = tls.RequireAndVerifyClientCert
        config.ClientCAs = certPool
    } else {
        // If --tlsverify is not supplied, disable CA validation.
        config.InsecureSkipVerify = true
    }

    return config, nil
}
