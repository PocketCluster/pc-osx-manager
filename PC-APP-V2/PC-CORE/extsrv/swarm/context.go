package swarm

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/docker/swarm/cluster"
    "github.com/docker/docker/pkg/discovery"
    "github.com/pkg/errors"
)

const (
    swarmServingHost string = "0.0.0.0:3376"
)

// Context is the global cluster setup context
type SwarmContext struct {
    // discovery option
    discoveryOpt       map[string]string
    // custom discovery backend
    discoveryBack      discovery.Backend

    heartbeat          time.Duration
    refreshMinInterval time.Duration
    refreshMaxInterval time.Duration
    refreshRetry       time.Duration
    failureRetry       int

    managerHost        []string
    debug              bool
    cors               bool

    // cluster strategy
    clusterOpt         cluster.DriverOpts
    strategy           string
    tlsConfig          *tls.Config
}

func NewContextWithCertAndKey(tlsCa, tlsCert, tlsKey []byte, discoveryBackend discovery.Backend) (*SwarmContext, error) {
    var (
        discoveryOpt = make(map[string]string)
        clusterOpt   = cluster.DriverOpts{}
    )

    // TODO : (04/17/2017) We should check if verifying clients with CA would results in errors for clients to connect. It appears to be ok with file version
    tlsConfig, err := buildTLSConfig(tlsCa, tlsCert, tlsKey, true)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    return &SwarmContext {
        discoveryOpt:       discoveryOpt,
        discoveryBack:      discoveryBackend,

        // FIXME : reuse beacon manager constants
        heartbeat:          time.Duration(10 * time.Second),
        refreshMinInterval: time.Duration(5 * time.Second),
        refreshMaxInterval: time.Duration(30 * time.Second),
        failureRetry:       6,
        // TODO : this will be removed.
        refreshRetry:       time.Duration(3),

        managerHost:        []string{swarmServingHost},
        debug:              true,

        clusterOpt:         clusterOpt,
        strategy:           "spread",
        tlsConfig:          tlsConfig,
    }, nil
}

// Load the TLS certificates/keys and, if verify is true, the CA.
func buildTLSConfig(ca, cert, key []byte, verify bool) (*tls.Config, error) {
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

// createNodesDiscovery replaces $GOPATH/src/github.com/docker/swarm/cli/manage/createDiscovery
// Instead of going through dokcker/pkg/discovery interface for the compatiblity with consul,
// this function creates node based backend directly.
func createNodeDiscovery(c *SwarmContext) discovery.Backend {
    // we're to go through BeaconManger directly
    err := c.discoveryBack.Initialize("pc-beacon", c.heartbeat, 0, getDiscoveryOpt(c))
    if err != nil {
        log.Fatal(err)
    }
    return c.discoveryBack
}

func getDiscoveryOpt(c *SwarmContext) map[string]string {
    // Process the store options
    options := map[string]string{}
    for key, value := range c.discoveryOpt {
        options[key] = value
    }
    if _, ok := options["kv.path"]; !ok {
        options["kv.path"] = "docker/swarm/nodes"
    }
    return options
}

// following methods will be deprecated //

func NewContext(host, nodeList string, tlsCa, tlsCert, tlsKey string) *SwarmContext {
    discoveryOpt := make(map[string]string)
    clusterOpt := cluster.DriverOpts{}
    tlsConfig, err := loadTLSConfigFromFiles(tlsCa, tlsCert, tlsKey, true)
    if err != nil {
        log.Fatal(err)
    }

    return &SwarmContext {
        discoveryOpt:       discoveryOpt,

        heartbeat:          time.Duration(10 * time.Second),
        refreshMinInterval: time.Duration(5 * time.Second),
        refreshMaxInterval: time.Duration(30 * time.Second),
        refreshRetry:       time.Duration(3),
        failureRetry:       6,

        managerHost:        []string{host},
        debug:              true,

        clusterOpt:         clusterOpt,
        strategy:           "spread",
        tlsConfig:          tlsConfig,
    }
}

// Load the TLS certificates/keys and, if verify is true, the CA.
func loadTLSConfigFromFiles(ca, cert, key string, verify bool) (*tls.Config, error) {
    c, err := tls.LoadX509KeyPair(cert, key)
    if err != nil {
        return nil, fmt.Errorf("Couldn't load X509 key pair (%s, %s): %s. Key encrypted?",
            cert, key, err)
    }

    config := &tls.Config{
        Certificates: []tls.Certificate{c},
        MinVersion:   tls.VersionTLS10,
    }

    if verify {
        certPool := x509.NewCertPool()
        file, err := ioutil.ReadFile(ca)
        if err != nil {
            return nil, fmt.Errorf("Couldn't read CA certificate: %s", err)
        }
        certPool.AppendCertsFromPEM(file)
        config.RootCAs = certPool
        config.ClientAuth = tls.RequireAndVerifyClientCert
        config.ClientCAs = certPool
    } else {
        // If --tlsverify is not supplied, disable CA validation.
        config.InsecureSkipVerify = true
    }

    return config, nil
}
