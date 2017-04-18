package etcd

import (
    "net/url"

    "github.com/coreos/etcd/embed"
    "github.com/coreos/etcd/pkg/cors"
    "github.com/coreos/etcd/etcdserver"
    "github.com/coreos/etcd/pkg/transport"
)

const (
    DefaultName                     = "pc-master"
    DefaultInitialClusterMember     = "pc-master=http://127.0.0.1:2380"
    DefaultAdvertiseClientURLs      = "https://pc-master:2379"
    DefaultListenClientURLs         = "https://0.0.0.0:2379"
    DefaultInitialAdvertisePeerURLs = "http://127.0.0.1:2380"
    DefaultListenPeerURLs           = "http://127.0.0.1:2380"
    DefaultInitialClusterToken      = "pocketcluster-kvstorage"
)

func NewEtcdConfig(dataDir string, tlsCa, tlsCert, tlsKey []byte) (*embed.Config, error) {
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
        // --data-dir
        Dir:                    dataDir,
        Name:                   DefaultName,                       // --name
        TickMs:                 5000,                              // --heartbeat-interval
        ElectionMs:             50000,                             // --election-timeout
        LPUrls:                 []url.URL{*lpurl},                 // --listen-peer-urls
        LCUrls:                 []url.URL{*lcurl},                 // --listen-client-urls
        APUrls:                 []url.URL{*apurl},                 // --initial-advertise-peer-urls
        ACUrls:                 []url.URL{*acurl},                 // --advertise-client-urls
        ClusterState:           embed.ClusterStateFlagNew,
        InitialCluster:         DefaultInitialClusterMember,       // --initial-cluster
        InitialClusterToken:    DefaultInitialClusterToken,
        StrictReconfigCheck:    true,
        Metrics:                "basic",

        // Do not auto generate any certificate
        ClientAutoTLS:          false,
        PeerAutoTLS:            false,
        Debug:                  false,

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
    return cfg, nil
}