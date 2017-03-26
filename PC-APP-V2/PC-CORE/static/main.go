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
    "github.com/cloudflare/cfssl/certdb"
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

//certificate authority generation
func certAuthSigner(certRec certdb.Accessor, meta *record.ClusterMeta, country string) (*context.CertAuthBundle, error) {
    var (
        signer *pcrypto.CaSigner = nil
        prvKey []byte  = nil
        pubKey []byte  = nil
        crtPem []byte  = nil
        err error      = nil
        caPrvRec, rerr = certRec.GetCertificate(pcdefaults.ClusterCertAuthPrivateKey, meta.ClusterUUID)
        caPubRec, uerr = certRec.GetCertificate(pcdefaults.ClusterCertAuthPublicKey, meta.ClusterUUID)
        caCrtRec, cerr = certRec.GetCertificate(pcdefaults.ClusterCertAuthCertificate, meta.ClusterUUID)
    )
    if (rerr != nil || uerr != nil || cerr != nil) || (len(caPrvRec) == 0 || len(caPubRec) == 0 || len(caCrtRec) == 0) {
        pubKey, prvKey, crtPem, err = pcrypto.GenerateClusterCertificateAuthorityData(meta.ClusterID, country)
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save private key
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(prvKey),
            Serial:     pcdefaults.ClusterCertAuthPrivateKey,
            AKI:        meta.ClusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save public key
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(pubKey),
            Serial:     pcdefaults.ClusterCertAuthPublicKey,
            AKI:        meta.ClusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save certificate
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(crtPem),
            Serial:     pcdefaults.ClusterCertAuthCertificate,
            AKI:        meta.ClusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
    } else {
        prvKey = []byte(caPrvRec[0].PEM)
        pubKey = []byte(caPubRec[0].PEM)
        crtPem = []byte(caCrtRec[0].PEM)
    }
    signer, err = pcrypto.NewCertAuthoritySigner(prvKey, crtPem, meta.ClusterID, country)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &context.CertAuthBundle{
        CASigner:    signer,
        CAPrvKey:    prvKey,
        CAPubKey:    pubKey,
        CACrtPem:    crtPem,
    }, nil
}

func hostCertificate(certRec certdb.Accessor, caSigner *pcrypto.CaSigner, hostname, clusterUUID string) (*context.HostCertBundle, error) {
    var (
        prvKey []byte  = nil
        pubKey []byte  = nil
        crtPem []byte  = nil
        sshPem []byte  = nil
        err error      = nil

        prvRec, rerr = certRec.GetCertificate(pcdefaults.MasterHostPrivateKey, clusterUUID)
        pubRec, uerr = certRec.GetCertificate(pcdefaults.MasterHostPublicKey, clusterUUID)
        crtRec, cerr = certRec.GetCertificate(pcdefaults.MasterHostCertificate, clusterUUID)
        sshRec, serr = certRec.GetCertificate(pcdefaults.MasterHostSshKey, clusterUUID)
    )

    if (rerr != nil || uerr != nil || cerr != nil || serr != nil) || (len(prvRec) == 0 || len(pubRec) == 0 || len(crtRec) == 0 || len(sshRec) == 0) {
        pubKey, prvKey, sshPem, err = pcrypto.GenerateStrongKeyPair()
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // we're not going to proide ip address for now
        crtPem, err = caSigner.GenerateSignedCertificate(hostname, "", prvKey)
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save private key
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(prvKey),
            Serial:     pcdefaults.MasterHostPrivateKey,
            AKI:        clusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save public key
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(pubKey),
            Serial:     pcdefaults.MasterHostPublicKey,
            AKI:        clusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save cert pem
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(crtPem),
            Serial:     pcdefaults.MasterHostCertificate,
            AKI:        clusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
        // save ssh pem
        err = certRec.InsertCertificate(certdb.CertificateRecord{
            PEM:        string(sshPem),
            Serial:     pcdefaults.MasterHostSshKey,
            AKI:        clusterUUID,
            Status:     "good",
            Reason:     0,
        })
        if err != nil {
            return nil, errors.WithStack(err)
        }
    } else {
        prvKey = []byte(prvRec[0].PEM)
        pubKey = []byte(pubRec[0].PEM)
        crtPem = []byte(crtRec[0].PEM)
        sshPem = []byte(sshRec[0].PEM)
    }
    return &context.HostCertBundle{
        PrivateKey:     prvKey,
        PublicKey:      pubKey,
        SshKey:         sshPem,
        Certificate:    crtPem,
    }, nil
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
        country = "US"
    }

    // certificate authority
    caBundle, err := certAuthSigner(rec.Certdb(), meta, country)
    if err != nil {
        // this is critical
        log.Debugf(err.Error())
    }
    context.UpdateCertAuth(caBundle)

    // host certificate
    hostBundle, err := hostCertificate(rec.Certdb(), caBundle.CASigner, teledefaults.CoreHostName, meta.ClusterUUID)
    if err != nil {
        // this is critical
        log.Debugf(err.Error())
    }
    context.UpdateHostCert(hostBundle)


    // make teleport core config
    cfg := service.MakeCoreConfig(dataDir, true)
    cfg.AssignHostUUID(meta.ClusterUUID)
    cfg.AssignDatabaseEngine(rec.DataBase())
    cfg.AssignCertStorage(rec.Certdb())
    cfg.AssignCASigner(caBundle.CASigner)
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
