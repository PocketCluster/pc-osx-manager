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
    "github.com/gravitational/teleport/lib/utils"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/event/lifecycle"
    "github.com/stkim1/pc-core/event/network"
    "github.com/stkim1/pc-core/event/crash"
    "github.com/stkim1/pc-core/event/operation"
    "github.com/stkim1/pc-core/record"
)
import (
    "github.com/tylerb/graceful"
    "github.com/davecgh/go-spew/spew"
    "github.com/gravitational/teleport/lib/process"
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

func openContext() (*service.PocketConfig, error) {
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
    cfg := service.MakeCoreConfig(dataDir, true)
    cfg.AssignHostUUID(meta.ClusterUUID)
    cfg.AssignDatabaseEngine(rec.DataBase())
    cfg.AssignCertStorage(rec.Certdb())
    cfg.AssignCASigner(caBundle.CASigner)
    cfg.AssignHostCertAuth(caBundle.CAPrvKey, caBundle.CASSHChk, meta.ClusterDomain)
    err = service.ValidateCoreConfig(cfg)
    if err != nil {
        log.Debugf(err.Error())
        return nil, errors.WithStack(err)
    }

    //log.Info(spew.Sdump(ctx))
    return cfg, nil
}

func main() {

    mainLifeCycle(func(a App) {

        var (
            teleConfig *service.PocketConfig = nil
            teleProc *process.PocketCoreProcess = nil
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
                            teleConfig, err = openContext()
                            if err != nil {
                                // TODO send error report
                            }
                            FeedSend("successfully initiated engine ..." + teleConfig.HostUUID)
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
                        teleProc, err = startTeleportCore(teleConfig)
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                    }
                    case operation.CmdTeleportStop: {
                        log.Debugf("[OP] %v", e.String())
                        err = teleProc.Close()
                        if err != nil {
                            log.Debugf("[ERR] " + err.Error())
                        }
                    }
                    default:
                        log.Print("[OP] %v", e.String())
                    }
                }
            }
        }
    })
}
