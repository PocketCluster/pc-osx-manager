package swarm

import (
    "crypto/tls"
/*
    "crypto/x509"
    "io/ioutil"
*/
    //"path"
    //"strings"
    "time"
    log "github.com/Sirupsen/logrus"

    // implicit loading and initialization
    _ "github.com/docker/docker/pkg/discovery/file"
    _ "github.com/docker/docker/pkg/discovery/nodes"
    _ "github.com/docker/swarm/discovery/token"

    "github.com/docker/leadership"
    "github.com/docker/swarm/api"
    "github.com/docker/swarm/cluster"
    "github.com/docker/swarm/cluster/swarm"
    "github.com/docker/swarm/scheduler"
    "github.com/docker/swarm/scheduler/filter"
    "github.com/docker/swarm/scheduler/strategy"
)

type logHandler struct {
}

func (h *logHandler) Handle(e *cluster.Event) error {
    id := e.ID
    // Trim IDs to 12 chars.
    if len(id) > 12 {
        id = id[:12]
    }
    log.WithFields(log.Fields{"node": e.Engine.Name, "id": id, "from": e.From, "status": e.Status}).Debug("Event received")
    return nil
}

type statusHandler struct {
    cluster   cluster.Cluster
    candidate *leadership.Candidate
    follower  *leadership.Follower
}

func (h *statusHandler) Status() [][2]string {
    var status [][2]string

    // for now, (09/21/2016) we won't have replication. Whatever comes in, just return primary
/*
    if h.candidate != nil && !h.candidate.IsLeader() {
        status = [][2]string{
            {"Role", "replica"},
            {"Primary", h.follower.Leader()},
        }
    } else {
        status = [][2]string{
            {"Role", "primary"},
        }
    }
*/
    status = [][2]string{{"Role", "primary"},}
    status = append(status, h.cluster.Info()...)
    return status
}

func (context *SwarmContext) Manage() {
    var (
        tlsConfig *tls.Config
        err       error
    )

    // we'll look into TLS certificate later
/*
    // If either --tls or --tlsverify are specified, load the certificates.
    if c.Bool("tls") || c.Bool("tlsverify") {
        if !c.IsSet("tlscert") || !c.IsSet("tlskey") {
            log.Fatal("--tlscert and --tlskey must be provided when using --tls")
        }
        if c.Bool("tlsverify") && !c.IsSet("tlscacert") {
            log.Fatal("--tlscacert must be provided when using --tlsverify")
        }
        tlsConfig, err = loadTLSConfig(
            c.String("tlscacert"),
            c.String("tlscert"),
            c.String("tlskey"),
            c.Bool("tlsverify"))
        if err != nil {
            log.Fatal(err)
        }
    } else {
        // Otherwise, if neither --tls nor --tlsverify are specified, abort if
        // the other flags are passed as they will be ignored.
        if c.IsSet("tlscert") || c.IsSet("tlskey") || c.IsSet("tlscacert") {
            log.Fatal("--tlscert, --tlskey and --tlscacert require the use of either --tls or --tlsverify")
        }
    }
*/
    refreshMinInterval := context.refreshMinInterval
    refreshMaxInterval := context.refreshMaxInterval
    if refreshMinInterval <= time.Duration(0)*time.Second {
        log.Fatal("min refresh interval should be a positive number")
    }
    if refreshMaxInterval < refreshMinInterval {
        log.Fatal("max refresh interval cannot be less than min refresh interval")
    }
    // engine-refresh-retry is deprecated
    refreshRetry := context.refreshRetry
    if refreshRetry != 3 {
        log.Fatal("--engine-refresh-retry is deprecated. Use --engine-failure-retry")
    }
    failureRetry := context.failureRetry
    if failureRetry <= 0 {
        log.Fatal("invalid failure retry count")
    }
    engineOpts := &cluster.EngineOpts{
        RefreshMinInterval: refreshMinInterval,
        RefreshMaxInterval: refreshMaxInterval,
        FailureRetry:       failureRetry,
    }

    uri := context.discoveryURI
    if uri == "" {
        log.Fatalf("discovery required to manage a cluster.")
    }
    //discovery := context.createDiscovery(uri)
    discovery := context.createTokenDiscovery()
    s, err := strategy.New(context.strategy)
    if err != nil {
        log.Fatal(err)
    }

    // see https://github.com/codegangsta/cli/issues/160
    names := []string{"health", "port", "containerslots", "dependency", "affinity", "constraint"}
    fs, err := filter.New(names)
    if err != nil {
        log.Fatal(err)
    }

    sched := scheduler.New(s, fs)
    var cl cluster.Cluster
    cl, err = swarm.NewCluster(sched, tlsConfig, discovery, context.clusterOpt, engineOpts)
    if err != nil {
        log.Fatal(err)
    }

    hosts := context.managerHost
    //server := api.NewServer(hosts, tlsConfig)
    server := NewServer(hosts, tlsConfig)

    primary := api.NewPrimary(cl, tlsConfig, &statusHandler{cl, nil, nil}, context.debug, context.cors)
    server.SetHandler(primary)
    cluster.NewWatchdog(cl)

    log.Fatal(server.ListenAndServe())
}