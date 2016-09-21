package swarm

import (
    "time"
    "github.com/docker/swarm/cluster"
    "github.com/docker/docker/pkg/discovery"
    log "github.com/Sirupsen/logrus"
    "github.com/docker/swarm/discovery/token"
)

// Context is the global cluster setup context
type SwarmContext struct {
    discoveryOpt       map[string]string
    discoveryURI       string             // e.g. token://1810ffdf37ad898423ada7262f7baf80

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
}

func NewContext(host string, token string) *SwarmContext {
    discoveryOpt := make(map[string]string)
    clusterOpt := cluster.DriverOpts{}

    return &SwarmContext{
        discoveryOpt: discoveryOpt,
        discoveryURI: token,

        heartbeat: time.Duration(1 * time.Second),
        refreshMinInterval: time.Duration(5 * time.Second),
        refreshMaxInterval: time.Duration(10 * time.Second),
        refreshRetry: time.Duration(3),
        failureRetry: 5,

        managerHost: []string{host},
        debug: true,

        clusterOpt:clusterOpt,
        strategy: "spread",
    }
}

// Initialize the discovery service.
func (c *SwarmContext)createDiscovery(uri string) discovery.Backend {
    hb := c.heartbeat
    if hb < 1*time.Second {
        log.Fatal("--heartbeat should be at least one second")
    }
    // Set up discovery.
    discovery, err := discovery.New(uri, hb, 0, c.getDiscoveryOpt())
    if err != nil {
        log.Fatal(err)
    }

    return discovery
}

// createTokenDiscovery replaces $GOPATH/src/github.com/docker/swarm/cli/manage/createDiscovery
// Instead of going through dokcker/pkg/discovery interface for the compatiblity with consul,
// this function creates token backend directly.
func (c *SwarmContext)createTokenDiscovery() discovery.Backend {
    hb := c.heartbeat
    if hb < 1*time.Second {
        log.Fatal("--heartbeat should be at least one second")
    }
    discovery := &token.Discovery{}
    err := discovery.Initialize(c.discoveryURI, hb, 0, c.getDiscoveryOpt())



    if err != nil {
        log.Fatal(err)
    }
    return discovery
}


func (c *SwarmContext)getDiscoveryOpt() map[string]string {
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
