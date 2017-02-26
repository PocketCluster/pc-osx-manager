package main

import (
    "log"
    "time"
    "net/url"

    "github.com/coreos/etcd/embed"
    "github.com/coreos/etcd/pkg/cors"
    "github.com/coreos/etcd/etcdserver"
    "github.com/coreos/etcd/pkg/transport"
)

/*
# leaner config. listen to local only for peer
bin/etcd \

--data-dir="/Users/almightykim/Workspace/DKIMG/ETCD/data" \
--name="pc-master" \

--heartbeat-interval=5000 \
--election-timeout=50000 \

--listen-peer-urls="http://127.0.0.1:2380" \
--listen-client-urls="https://0.0.0.0:2379" \
--initial-advertise-peer-urls="http://127.0.0.1:2380" \
--advertise-client-urls="https://pc-master:2379" \

--initial-cluster="pc-master=http://127.0.0.1:2380" \

--cert-file="/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.cert" \
--key-file="/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.key" \
--trusted-ca-file="/Users/almightykim/Workspace/DKIMG/CERT/ca-cert.pub" \
--client-cert-auth=true \

--debug
*/

const (
    DefaultAdvertiseClientURLs      = "https://pc-master:2379"
    DefaultListenClientURLs         = "https://0.0.0.0:2379"
    DefaultInitialAdvertisePeerURLs = "http://127.0.0.1:2380"
    DefaultListenPeerURLs           = "http://127.0.0.1:2380"
    DefaultName                     = "pc-master"
    DefaultInitialClusterMember     = "pc-master=http://127.0.0.1:2380"
)

/*
 For full Config options, take a look at "github.com/coreos/etcd/etcdmain/config.go"
 */
func NewEtcdConfig(dataDir string) *embed.Config {
    // --listen-peer-urls
    lpurl, _ := url.Parse(DefaultListenPeerURLs)
    // --initial-advertise-peer-urls
    apurl, _ := url.Parse(DefaultInitialAdvertisePeerURLs)
    // --listen-client-urls
    lcurl, _ := url.Parse(DefaultListenClientURLs)
    // --advertise-client-urls
    acurl, _ := url.Parse(DefaultAdvertiseClientURLs)

    cfg := &embed.Config {
        CorsInfo:               &cors.CORSInfo{},
        MaxSnapFiles:           embed.DefaultMaxSnapshots,
        MaxWalFiles:            embed.DefaultMaxWALs,
        SnapCount:              etcdserver.DefaultSnapCount,
        Dir:                    dataDir,                           // --data-dir
        Name:                   DefaultName,                       // --name
        TickMs:                 5000,                              // --heartbeat-interval
        ElectionMs:             50000,                             // --election-timeout
        LPUrls:                 []url.URL{*lpurl},                 // --listen-peer-urls
        LCUrls:                 []url.URL{*lcurl},                 // --listen-client-urls
        APUrls:                 []url.URL{*apurl},                 // --initial-advertise-peer-urls
        ACUrls:                 []url.URL{*acurl},                 // --advertise-client-urls
        ClusterState:           embed.ClusterStateFlagNew,
        InitialCluster:         DefaultInitialClusterMember,       // --initial-cluster
        InitialClusterToken:    "etcd-cluster",
        StrictReconfigCheck:    true,
        Metrics:                "basic",

        // Do not auto generate any certificate
        ClientAutoTLS:          false,
        PeerAutoTLS:            false,
        Debug:                  true,

        // client certificate options
        ClientTLSInfo:          transport.TLSInfo {
            // --cert-file
            CertFile:           "/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.cert",
            // --key-file
            KeyFile:            "/Users/almightykim/Workspace/DKIMG/PC-MASTER/pc-master.key",
            // --trusted-ca-file
            TrustedCAFile:      "/Users/almightykim/Workspace/DKIMG/CERT/ca-cert.pub",
            // --client-cert-auth
            ClientCertAuth:     true,
        },
    }
    return cfg
}

func main() {
    cfg := NewEtcdConfig("/Users/almightykim/Workspace/DKIMG/ETCD/data")
    e, err := embed.StartEtcd(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer e.Close()
    select {
    case <-e.Server.ReadyNotify():
        log.Printf("Server is ready!")
    case <-time.After(60 * time.Second):
        e.Server.Stop() // trigger a shutdown
        log.Printf("Server took too long to start!")
    }
    log.Fatal(<-e.Err())
}