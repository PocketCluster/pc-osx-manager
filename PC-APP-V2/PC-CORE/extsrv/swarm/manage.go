package swarm

import (
    "time"

    log "github.com/Sirupsen/logrus"
    // implicit loading and initialization
    _ "github.com/docker/docker/pkg/discovery/nodes"
    "github.com/docker/leadership"
    "github.com/docker/swarm/api"
    "github.com/docker/swarm/cluster"
    "github.com/docker/swarm/cluster/swarm"
    "github.com/docker/swarm/scheduler"
    "github.com/docker/swarm/scheduler/filter"
    "github.com/docker/swarm/scheduler/strategy"

    "github.com/pkg/errors"
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

func (context *SwarmContext) Manage() error {
    refreshMinInterval := context.refreshMinInterval
    refreshMaxInterval := context.refreshMaxInterval
    if refreshMinInterval <= time.Duration(0) * time.Second {
        return errors.Errorf("min refresh interval should be a positive number")
    }
    if refreshMaxInterval < refreshMinInterval {
        return errors.Errorf("max refresh interval cannot be less than min refresh interval")
    }
    // engine-refresh-retry is deprecated
    refreshRetry := context.refreshRetry
    if refreshRetry != 3 {
        return errors.Errorf("--engine-refresh-retry is deprecated. Use --engine-failure-retry")
    }
    failureRetry := context.failureRetry
    if failureRetry <= 0 {
        return errors.Errorf("invalid failure retry count")
    }
    engineOpts := &cluster.EngineOpts {
        RefreshMinInterval: refreshMinInterval,
        RefreshMaxInterval: refreshMaxInterval,
        FailureRetry:       failureRetry,
    }

    discovery := createNodeDiscovery(context)
    s, err := strategy.New(context.strategy)
    if err != nil {
        return errors.WithStack(err)
    }

    // see https://github.com/codegangsta/cli/issues/160
    names := []string{"health", "port", "containerslots", "dependency", "affinity", "constraint"}
    fs, err := filter.New(names)
    if err != nil {
        return errors.WithStack(err)
    }

    sched := scheduler.New(s, fs)
    var cl cluster.Cluster
    cl, err = swarm.NewCluster(sched, context.tlsConfig, discovery, context.clusterOpt, engineOpts)
    if err != nil {
        return errors.WithStack(err)
    }

    hosts := context.managerHost
    server := newService(hosts, context.tlsConfig)
    primary := api.NewPrimary(cl, context.tlsConfig, &statusHandler{cl, nil, nil}, context.debug, context.cors)
    server.SetHandler(primary)
    cluster.NewWatchdog(cl)

    return errors.WithStack(server.ListenAndServeMultiHosts())
}

func NewSwarmServer(context *SwarmContext) (*Service, error) {
    refreshMinInterval := context.refreshMinInterval
    refreshMaxInterval := context.refreshMaxInterval
    if refreshMinInterval <= time.Duration(0) * time.Second {
        return nil, errors.Errorf("min refresh interval should be a positive number")
    }
    if refreshMaxInterval < refreshMinInterval {
        return nil, errors.Errorf("max refresh interval cannot be less than min refresh interval")
    }
    // engine-refresh-retry is deprecated
    refreshRetry := context.refreshRetry
    if refreshRetry != 3 {
        return nil, errors.Errorf("--engine-refresh-retry is deprecated. Use --engine-failure-retry")
    }
    failureRetry := context.failureRetry
    if failureRetry <= 0 {
        return nil, errors.Errorf("invalid failure retry count")
    }
    engineOpts := &cluster.EngineOpts {
        RefreshMinInterval: refreshMinInterval,
        RefreshMaxInterval: refreshMaxInterval,
        FailureRetry:       failureRetry,
    }

    discovery := createNodeDiscovery(context)
    s, err := strategy.New(context.strategy)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // see https://github.com/codegangsta/cli/issues/160
    names := []string{"health", "port", "containerslots", "dependency", "affinity", "constraint"}
    fs, err := filter.New(names)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    sched := scheduler.New(s, fs)
    var cl cluster.Cluster
    cl, err = swarm.NewCluster(sched, context.tlsConfig, discovery, context.clusterOpt, engineOpts)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    hosts := context.managerHost
    server, err := newStoppableServiceForSingleHost(hosts, context.tlsConfig)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    primary := api.NewPrimary(cl, context.tlsConfig, &statusHandler{cl, nil, nil}, context.debug, context.cors)
    server.SetHandler(primary)
    cluster.NewWatchdog(cl)

    return server, nil
}