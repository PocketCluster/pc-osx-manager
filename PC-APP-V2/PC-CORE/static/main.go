package main

import "C"
import (
    "net/http"
    "fmt"
    "time"
    "sync"

    log "github.com/Sirupsen/logrus"
    teledefaults "github.com/gravitational/teleport/lib/defaults"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/process"
    "github.com/gravitational/teleport/lib/utils"

    "github.com/coreos/etcd/embed"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/lifecycle"
    "github.com/stkim1/pc-core/event/network"
    "github.com/stkim1/pc-core/event/crash"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/record"
    telesrv "github.com/stkim1/pc-core/extsrv/teleport"
    regisrv "github.com/stkim1/pc-core/extsrv/registry"
    swarmsrv "github.com/stkim1/pc-core/extsrv/swarm"
)
import (
    "github.com/tylerb/graceful"
    "github.com/davecgh/go-spew/spew"

)

func RunWebServer(wg *sync.WaitGroup) *graceful.Server {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        fmt.Fprintf(w, "Welcome to the home page!")
    })
    srv := &graceful.Server{
        Timeout: 10 * time.Second,
        NoSignalHandling: true,
        Server: &http.Server{
            Addr: ":3001",
            Handler: mux,
        },
    }

    go func() {
        defer wg.Done()
        srv.ListenAndServe()
    }()
    return srv
}

func setLogger(debug bool) {
    // debug setup
    if debug {
        utils.InitLoggerDebug()
        log.Info("DEBUG mode logger output configured")
    } else {
        utils.InitLoggerCLI()
        log.Info("NORMAL mode logger configured")
    }
}

func main_old() {
    setLogger(true)

    var wg sync.WaitGroup
    srv := RunWebServer(&wg)

    go func() {
        wg.Wait()
    }()

    // setup context
    ctx := context.SharedHostContext()
    context.SetupBasePath()

    // open database
    dataDir, err := ctx.ApplicationUserDataDirectory()
    if err != nil {
        log.Info(err)
    }
    rec, err := record.OpenRecordGate(dataDir, teledefaults.CoreKeysSqliteFile)
    if err != nil {
        log.Info(err)
    }

    cfg := service.MakeCoreConfig(dataDir, true)
    cfg.AssignCertStorage(rec.Certdb())

    log.Info(spew.Sdump(ctx))

    // Perhaps the first thing main() function needs to do is initiate OSX main
    //C.osxmain(0, nil)

    srv.Stop(time.Second)
    record.CloseRecordGate()
    fmt.Println("pc-core terminated!")
}

type serviceConfig struct {
    etcdConfig     *embed.PocketConfig
    teleConfig     *service.PocketConfig
    regConfig      *regisrv.PocketRegistryConfig
    swarmConfig    *swarmsrv.SwarmContext
}

func openContext() (*serviceConfig, error) {
    // setup context
    ctx := context.SharedHostContext()
    context.SetupBasePath()

    // open database
    dataDir, err := ctx.ApplicationUserDataDirectory()
    if err != nil {
        log.Info(err)
        return nil, errors.WithStack(err)
    }
    rec, err := record.OpenRecordGate(dataDir, teledefaults.CoreKeysSqliteFile)
    if err != nil {
        log.Info(err)
        return nil, errors.WithStack(err)
    }

    // new cluster id
    var meta *record.ClusterMeta = nil
    cluster, err := record.FindClusterMeta()
    if err != nil {
        if err == record.NoItemFound {
            meta = record.NewClusterMeta()
            record.UpsertClusterMeta(meta)
        } else {
            // This is critical error. report it to UI and ask them to clean & re-install
            return nil, errors.WithStack(err)
        }
    } else {
        meta = cluster[0]
    }
    log.Debugf("Cluster ID %v | UUID %v", meta.ClusterID, meta.ClusterUUID)

    country, err := ctx.CurrentCountryCode()
    if err != nil {
        // (03/26/2017) skip coutry code error and defaults it to US
        country = "US"
    }

    // certificate authority
    caBundle, err := certAuthSigner(rec.Certdb(), meta, country)
    if err != nil {
        // this is critical
        log.Debugf(err.Error())
        return nil, errors.WithStack(err)
    }
    context.UpdateCertAuth(caBundle)

    // host certificate
    hostBundle, err := hostCertificate(rec.Certdb(), caBundle.CASigner, teledefaults.CoreHostName, meta.ClusterUUID)
    if err != nil {
        // this is critical
        log.Debugf(err.Error())
        return nil, errors.WithStack(err)
    }
    context.UpdateHostCert(hostBundle)

    // make teleport core config
    teleCfg := service.MakeCoreConfig(dataDir, true)
    teleCfg.AssignHostUUID(meta.ClusterUUID)
    teleCfg.AssignDatabaseEngine(rec.DataBase())
    teleCfg.AssignCertStorage(rec.Certdb())
    teleCfg.AssignCASigner(caBundle.CASigner)
    teleCfg.AssignHostCertAuth(caBundle.CAPrvKey, caBundle.CASSHChk, meta.ClusterDomain)
    err = service.ValidateCoreConfig(teleCfg)
    if err != nil {
        log.Debugf(err.Error())
        return nil, errors.WithStack(err)
    }

    // registry configuration
    // TODO : fix datadir. Plus, is it ok not to pass CA pub key? we need to unify TLS configuration
    regCfg, err := regisrv.NewPocketRegistryConfig(false, dataDir, hostBundle.Certificate, hostBundle.PrivateKey)
    if err != nil {
        log.Debugf(err.Error())
        return nil, errors.WithStack(err)
    }

    // swarm configuration
    swarmCfg, err := swarmsrv.NewContextWithCertAndKey(
        "0.0.0.0:3376",
        "192.168.1.150:2375,192.168.1.151:2375,192.168.1.152:2375,192.168.1.153:2375,192.168.1.161:2375,192.168.1.162:2375,192.168.1.163:2375,192.168.1.164:2375,192.168.1.165:2375,192.168.1.166:2375",
        caBundle.CACrtPem,
        hostBundle.Certificate,
        hostBundle.PrivateKey,
    )
    if err != nil {
        // this is critical
        log.Debugf(err.Error())
        return nil, errors.WithStack(err)
    }

    //etcd configuration
    // TODO fix datadir
    etcdCfg, err := embed.NewPocketConfig(dataDir, caBundle.CACrtPem, hostBundle.Certificate, hostBundle.PrivateKey)
    if err != nil {
        // this is critical
        log.Debugf(err.Error())
        return nil, errors.WithStack(err)
    }
    //log.Info(spew.Sdump(ctx))
    return &serviceConfig {
        etcdConfig: etcdCfg,
        teleConfig: teleCfg,
        regConfig: regCfg,
        swarmConfig: swarmCfg,
    }, nil
}

func main() {

    mainLifeCycle(func(a App) {

        var (
            serviceConfig *serviceConfig = nil
            teleProc *process.PocketCoreProcess = nil
            regiProc *regisrv.PocketRegistry = nil
            swarmProc *swarmsrv.Server
//            etcdProc *embed.PocketEtcd
            err error = nil
        )

        for e := range a.Events() {
            switch e := a.Filter(e).(type) {

                case lifecycle.Event: {
                    switch e.Crosses(lifecycle.StageDead) {
                        case lifecycle.CrossOn: {
                            log.Debugf("[LIFE] app is Dead %v", e.String())
                        }
                        case lifecycle.CrossOff: {
                            log.Debugf("[LIFE] app is not dead %v", e.String())
                        }
                    }

                    switch e.Crosses(lifecycle.StageAlive) {
                        case lifecycle.CrossOn: {
                            log.Debugf("[LIFE] app is now alive %v", e.String())
                            log.Debugf("[PREP] PREPARING GOLANG CONTEXT")
                            serviceConfig, err = openContext()
                            if err != nil {
                                // TODO send error report
                            }
                            FeedSend("successfully initiated engine ..." + serviceConfig.teleConfig.HostUUID)
                        }
                        case lifecycle.CrossOff: {
                            log.Debugf("[LIFE] app is inactive %v", e.String())
                        }
                    }

                    switch e.Crosses(lifecycle.StageVisible) {
                        case lifecycle.CrossOn: {
                            log.Debugf("[LIFE] app is visible %v", e.String())
                        }
                        case lifecycle.CrossOff: {
                            log.Debugf("[LIFE] app is invisible %v", e.String())
                        }
                    }

                    switch e.Crosses(lifecycle.StageFocused) {
                        case lifecycle.CrossOn: {
                            log.Debugf("[LIFE] app is focused %v", e.String())
                        }
                        case lifecycle.CrossOff: {
                            log.Debugf("[LIFE] app is not focused %v", e.String())
                        }
                    }
                }

                case network.Event: {
                    switch e.NetworkEvent {
                        case network.NetworkChangeInterface: {
                            //log.Debugf(spew.Sdump(e.HostInterfaces))
                            log.Debugf("[NET] %v", e.String())
                            context.MonitorNetworkInterfaces(e.HostInterfaces)
                        }
                        case network.NetworkChangeGateway: {
                            //log.Debugf(spew.Sdump(e.HostGateways))
                            log.Debugf("[NET] %v", e.String())
                            context.MonitorNetworkGateways(e.HostGateways)
                        }
                    }
                }

                // artificial crash
                case crash.Crash: {
                    switch e.Reason {
                    case crash.CrashEmergentExit: {
                        log.Printf("[CRASH] COCOA SIDE RUNTIME IS DESTORYED. WE NEED TO CLOSE GOLANG SIDE AS WELL. %v", e.String())
                    }
                    default:
                        log.Printf("crash! %v", e.String())
                    }
                }

                // operational Command
                case operation.Operation: {
                    switch e.Command {
                    case operation.CmdTeleportStart: {
                        log.Debugf("[OP] %v", e.String())

                        teleProc, err = telesrv.NewTeleportCore(serviceConfig.teleConfig)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        err = teleProc.Start()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                    }
                    case operation.CmdTeleportStop: {
                        err = teleProc.Close()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        err = teleProc.Wait()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        log.Debugf("[OP] %v", e.String())
                    }

                    case operation.CmdRegistryStart: {
                        log.Debugf("[OP] %v", e.String())
                        regiProc, err = regisrv.NewPocketRegistry(serviceConfig.regConfig)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        err = regiProc.Start()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                    }
                    case operation.CmdRegistryStop: {
                        err = regiProc.Close()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        log.Debugf("[OP] %v", e.String())
                    }

                    case operation.CmdCntrOrchStart: {
                        swarmProc, err = swarmsrv.NewSwarmServer(serviceConfig.swarmConfig)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        err = swarmProc.ListenAndServe()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        log.Debugf("[OP] %v", e.String())
                    }
                    case operation.CmdCntrOrchStop: {
                        log.Debugf("[OP] %v", e.String())
                    }

                    case operation.CmdStorageStart: {
                        log.Debugf("[OP] %v", e.String())
/*
                        etcdProc, err = embed.StartPocketEtcd(serviceConfig.etcdConfig)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                        etcdProc.Server.Start()
*/
                        go func() {
                            e, err := embed.StartPocketEtcd(serviceConfig.etcdConfig)
                            if err != nil {
                                log.Debugf(err.Error())
                                return
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
                        }()
                    }
                    case operation.CmdStorageStop: {
                        log.Debugf("[OP] %v", e.String())
//                        etcdProc.Server.Stop()
                    }
                    default:
                        log.Print("[OP] %v", e.String())
                    }
                }
            }
        }
    })
}
