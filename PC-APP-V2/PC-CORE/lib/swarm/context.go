package swarm

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/docker/swarm/cluster"
    "github.com/docker/docker/pkg/discovery"
    "github.com/docker/docker/pkg/discovery/nodes"
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
}

func NewContext(host string, nodeList string) *SwarmContext {
    discoveryOpt := make(map[string]string)
    clusterOpt := cluster.DriverOpts{}

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
    }
}

// createNodesDiscovery replaces $GOPATH/src/github.com/docker/swarm/cli/manage/createDiscovery
// Instead of going through dokcker/pkg/discovery interface for the compatiblity with consul,
// this function creates node based backend directly.
func (c *SwarmContext) CreateNodeDiscovery() discovery.Backend {
    hb := c.heartbeat
    if hb < 1*time.Second {
        log.Fatal("--heartbeat should be at least one second")
    }
    discovery := &nodes.Discovery{}
    err := discovery.Initialize(c.nodeList, hb, 0, c.DiscoveryOpt())
    if err != nil {
        log.Fatal(err)
    }
    return discovery
}

func (c *SwarmContext) DiscoveryOpt() map[string]string {
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
