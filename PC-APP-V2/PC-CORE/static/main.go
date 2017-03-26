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
    pcdefaults "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/event/lifecycle"
    "github.com/stkim1/pc-core/event/network"
    "github.com/stkim1/pc-core/hostapi"
    "github.com/stkim1/pc-core/record"
    "github.com/stkim1/pcrypto"
)
import (
    "github.com/tylerb/graceful"
    "github.com/davecgh/go-spew/spew"
    "github.com/cloudflare/cfssl/certdb"
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

func prepEnviornment() {
    setLogger(true)

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

    // new cluster id
    var meta *record.ClusterMeta = nil
    cluster, err := record.FindClusterMeta()
    if err != nil {
        if err == record.NoItemFound {
            meta = record.NewClusterMeta()
            record.UpsertClusterMeta(meta)
        } else {
            // This is critical error. report it to UI and ask them to clean & re-install
            log.Info(err)
        }
    } else {
        meta = cluster[0]
    }
    log.Debugf("Cluster ID %v | UUID %v", meta.ClusterID, meta.ClusterUUID)

    country, err := ctx.CurrentCountryCode()
    if err != nil {
        // (03/26/2017) skip coutry code error and defaults it to US
        log.Debugf(err.Error())
        country = "US"
    }

    //certificate authority generation
    var (
        prvKey, pubKey, certPem []byte = nil, nil, nil
        caPrvKey, rerr = rec.Certdb().GetCertificate(pcdefaults.MasterCertAuthPrivateKey, meta.ClusterUUID)
        caPubKey, uerr = rec.Certdb().GetCertificate(pcdefaults.MasterCertAuthPublicKey, meta.ClusterUUID)
        caCert, cerr   = rec.Certdb().GetCertificate(pcdefaults.MasterCertAuthCertificate, meta.ClusterUUID)
    )
    if (rerr != nil || uerr != nil || cerr != nil) || (len(caPrvKey) == 0 || len(caPubKey) == 0 || len(caCert) == 0) {
        pubKey, prvKey, certPem, err = pcrypto.GenerateClusterCertificateAuthorityData(meta.ClusterID, country)
        if err != nil {
            // this is critical
            log.Debugf(errors.WithStack(err).Error())
        }
        // save private key
        err = rec.Certdb().InsertCertificate(certdb.CertificateRecord{
            PEM:        string(prvKey),
            Serial:     pcdefaults.MasterCertAuthPrivateKey,
            AKI:        meta.ClusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            // this is critical
            log.Debugf(errors.WithStack(err).Error())
        }
        // save public key
        err = rec.Certdb().InsertCertificate(certdb.CertificateRecord{
            PEM:        string(pubKey),
            Serial:     pcdefaults.MasterCertAuthPublicKey,
            AKI:        meta.ClusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            // this is critical
            log.Debugf(errors.WithStack(err).Error())
        }
        // save certificate
        err = rec.Certdb().InsertCertificate(certdb.CertificateRecord{
            PEM:        string(certPem),
            Serial:     pcdefaults.MasterCertAuthCertificate,
            AKI:        meta.ClusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            // this is critical
            log.Debugf(errors.WithStack(err).Error())
        }
    } else {
        prvKey  = []byte(caPrvKey[0].PEM)
        pubKey  = []byte(caPubKey[0].PEM)
        certPem = []byte(caCert[0].PEM)
    }
    caSigner, err := pcrypto.NewCertAuthoritySigner(prvKey, certPem, meta.ClusterID, country)
    if err != nil {
        // this is critical
        log.Debugf(errors.WithStack(err).Error())
    }
    context.SetupCertAuthSigner(caSigner)



    // make teleport core config
    cfg := service.MakeCoreConfig(dataDir, true)
    cfg.AssignHostUUID(meta.ClusterUUID)
    cfg.AssignDatabaseEngine(rec.DataBase())
    cfg.AssignCertStorage(rec.Certdb())
    cfg.AssignCASigner(caSigner)
    err = service.ValidateCoreConfig(cfg)
    if err != nil {
        log.Debugf(err.Error())
    }

    log.Info(spew.Sdump(ctx))
}

func main() {

    prepEnviornment()

    mainLifeCycle(func(a App) {
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

                            hostapi.SendFeedBack("successfully initiated engine ...")
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
                            log.Debugf(spew.Sdump(e.HostInterfaces))
                            context.MonitorNetworkInterfaces(e.HostInterfaces)
                        }
                        case network.NetworkChangeGateway: {
                            log.Debugf(spew.Sdump(e.HostGateways))
                            context.MonitorNetworkGateways(e.HostGateways)
                        }
                    }
                }
            }
        }
    })
}
