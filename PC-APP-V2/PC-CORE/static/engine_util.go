package main

import (
    log "github.com/Sirupsen/logrus"
    tefaults "github.com/gravitational/teleport/lib/defaults"
    tervice "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/teleport/lib/utils"

    "github.com/coreos/etcd/embed"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/extlib/registry"
)

func setLogger(debug bool) {
    // debug setup
    if debug {
        utils.InitLoggerDebug()
        log.Info("DEBUG mode logger output configured")
    } else {
        utils.InitLoggerCLI()
        log.Info("NORMAL mode logger configured")
    }
    log.SetFormatter(&log.TextFormatter{
        DisableColors:    true,
    })
}

type serviceConfig struct {
    etcdConfig     *embed.PocketConfig
    teleConfig     *tervice.PocketConfig
    regConfig      *registry.PocketRegistryConfig
}

func setupServiceConfig() (*serviceConfig, error) {
    // setup context
    ctx := context.SharedHostContext()
    context.SetupBasePath()

    // open database
    dataDir, err := ctx.ApplicationUserDataDirectory()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    rec, err := model.OpenRecordGate(dataDir, tefaults.CoreKeysSqliteFile)
    if err != nil {
        log.Info(err)
        return nil, errors.WithStack(err)
    }

    // new cluster id
    var meta *model.ClusterMeta = nil
    cluster, err := model.FindClusterMeta()
    if err != nil {
        if err == model.NoItemFound {
            meta = model.NewClusterMeta()
            model.UpsertClusterMeta(meta)
        } else {
            // This is critical error. report it to UI and ask them to clean & re-install
            return nil, errors.WithStack(err)
        }
    } else {
        meta = cluster[0]
    }
    log.Debugf("Cluster ID %v | UUID %v", meta.ClusterID, meta.ClusterUUID)
    ctx.SetMasterAgentName(meta.ClusterID)

    country, err := ctx.CurrentCountryCode()
    if err != nil {
        // (03/26/2017) skip coutry code error and defaults it to US
        country = "US"
    }

    // certificate authority
    caBundle, err := buildCertAuthSigner(rec.Certdb(), meta, country)
    if err != nil {
        // this is critical
        return nil, errors.WithStack(err)
    }
    ctx.UpdateCertAuth(caBundle)

    // host certificate
    hostBundle, err := buildHostCertificate(rec.Certdb(), caBundle.CASigner, tefaults.CoreHostName, meta.ClusterUUID)
    if err != nil {
        // this is critical
        return nil, errors.WithStack(err)
    }
    ctx.UpdateHostCert(hostBundle)

    // beacon certificate
    beaconBundle, err := buildBeaconCertificate(rec.Certdb(), meta.ClusterUUID)
    if err != nil {
        // this is critical
        return nil, errors.WithStack(err)
    }
    ctx.UpdateBeaconCert(beaconBundle)

    // make teleport core config
    teleCfg := tervice.MakeCoreConfig(dataDir, true)
    teleCfg.AssignHostUUID(meta.ClusterUUID)
    teleCfg.AssignDatabaseEngine(rec.DataBase())
    teleCfg.AssignCertStorage(rec.Certdb())
    teleCfg.AssignCASigner(caBundle.CASigner)
    teleCfg.AssignHostCertAuth(caBundle.CAPrvKey, caBundle.CASSHChk, meta.ClusterDomain)
    err = tervice.ValidateCoreConfig(teleCfg)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // registry configuration
    // TODO : fix datadir. Plus, is it ok not to pass CA pub key? we need to unify TLS configuration
    regCfg, err := registry.NewPocketRegistryConfig(false, dataDir, hostBundle.Certificate, hostBundle.PrivateKey)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    //etcd configuration
    // TODO fix datadir
    etcdCfg, err := embed.NewPocketConfig(dataDir, caBundle.CACrtPem, hostBundle.Certificate, hostBundle.PrivateKey)
    if err != nil {
        // this is critical
        return nil, errors.WithStack(err)
    }
    //log.Info(spew.Sdump(ctx))
    return &serviceConfig {
        etcdConfig: etcdCfg,
        teleConfig: teleCfg,
        regConfig: regCfg,
    }, nil
}
