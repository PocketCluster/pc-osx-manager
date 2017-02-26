package swarm

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
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

// Load the TLS certificates/keys and, if verify is true, the CA.
func loadTLSConfig(ca, cert, key string, verify bool) (*tls.Config, error) {
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

/*
swarm manage
--debug
--host =            :3376
--advertise=        pc-master:3376

--tlsverify =       true
--tlscacert =       /Users/almightykim/Workspace/DKIMG/CERT/ca-cert.pub
--tlscert =         /Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.cert
--tlskey =          /Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.key

nodes://192.168.1.150:2375,
192.168.1.151:2375,
192.168.1.152:2375,
192.168.1.153:2375,
192.168.1.161:2375,
192.168.1.162:2375,
192.168.1.163:2375,
192.168.1.164:2375,
192.168.1.165:2375,
192.168.1.166:2375
*/

const (
    DefaultTLSCA    = "/Users/almightykim/Workspace/DKIMG/CERT/ca-cert.pub"
    DefaultTLSCert  = "/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.cert"
    DefaultTLSKey   = "/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.key"
)

func (context *SwarmContext) Manage() {
    var (
        tlsConfig *tls.Config
        err       error
    )

    tlsConfig, err = loadTLSConfig(
        DefaultTLSCA,
        DefaultTLSCert,
        DefaultTLSKey,
        true)
    if err != nil {
        log.Fatal(err)
    }

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
    engineOpts := &cluster.EngineOpts {
        RefreshMinInterval: refreshMinInterval,
        RefreshMaxInterval: refreshMaxInterval,
        FailureRetry:       failureRetry,
    }

    // FIXME : this should check the validity of node list (form, # of items, etc)
    uri := context.nodeList
    if uri == "" {
        log.Fatalf("discovery required to manage a cluster.")
    }
    discovery := context.createNodeDiscovery()
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
    server := NewServer(hosts, tlsConfig)
    primary := api.NewPrimary(cl, tlsConfig, &statusHandler{cl, nil, nil}, context.debug, context.cors)
    server.SetHandler(primary)
    cluster.NewWatchdog(cl)

    log.Fatal(server.ListenAndServe())
}
