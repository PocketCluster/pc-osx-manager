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
    "github.com/docker/docker/pkg/discovery/nodes"
    "github.com/pkg/errors"
)

// Context is the global cluster setup context
type SwarmContext struct {
    discoveryOpt       map[string]string

    // https://docs.docker.com/swarm/discovery/#/to-use-a-node-list
    // e.g : <node_ip1:2375>,<node_ip2:2375>
    nodeList           string

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

func NewContext(host, nodeList string, tlsCa, tlsCert, tlsKey string) *SwarmContext {
    discoveryOpt := make(map[string]string)
    clusterOpt := cluster.DriverOpts{}
    tlsConfig, err := loadTLSConfigFromFiles(tlsCa, tlsCert, tlsKey, true)
    if err != nil {
        log.Fatal(err)
    }

    return &SwarmContext {
        discoveryOpt:       discoveryOpt,
        nodeList:           nodeList,

        heartbeat:          time.Duration(1 * time.Second),
        refreshMinInterval: time.Duration(5 * time.Second),
        refreshMaxInterval: time.Duration(10 * time.Second),
        refreshRetry:       time.Duration(3),
        failureRetry:       5,

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

func NewContextWithCertAndKey(host, nodeList string, tlsCa, tlsCert, tlsKey []byte) (*SwarmContext, error) {
    discoveryOpt := make(map[string]string)
    clusterOpt := cluster.DriverOpts{}

    // TODO : (04/17/2017) We should check if verifying clients with CA would results in errors for clients to connect. It appears to be ok with file version
    tlsConfig, err := buildTLSConfig(tlsCa, tlsCert, tlsKey, true)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    return &SwarmContext {
        discoveryOpt:       discoveryOpt,
        nodeList:           nodeList,

        heartbeat:          time.Duration(1 * time.Second),
        refreshMinInterval: time.Duration(5 * time.Second),
        refreshMaxInterval: time.Duration(10 * time.Second),
        refreshRetry:       time.Duration(3),
        failureRetry:       5,

        managerHost:        []string{host},
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
    hb := c.heartbeat
    if hb < 1*time.Second {
        log.Fatal("--heartbeat should be at least one second")
    }
    discovery := &nodes.Discovery{}
    err := discovery.Initialize(c.nodeList, hb, 0, getDiscoveryOpt(c))
    if err != nil {
        log.Fatal(err)
    }
    return discovery
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
