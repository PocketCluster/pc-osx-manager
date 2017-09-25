package main

import (
    "github.com/gravitational/teleport/lib/defaults"
    "github.com/gravitational/teleport/lib/service"
    "github.com/gravitational/trace"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/record"
)

func openConfig() (*service.PocketConfig, error) {
    // setup context
    ctx := context.SharedHostContext()
    context.SetupBasePath()

    // open database
    dataDir, err := ctx.ApplicationUserDataDirectory()
    if err != nil {
        log.Info(err)
        return nil, trace.Wrap(err)
    }
    rec, err := record.OpenRecordGate(dataDir, defaults.CoreKeysSqliteFile)
    if err != nil {
        log.Info(err)
        return nil, trace.Wrap(err)
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
            return nil, trace.Wrap(err)
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
        return nil, trace.Wrap(err)
    }
    context.UpdateCertAuth(caBundle)

    // host certificate
    hostBundle, err := hostCertificate(rec.Certdb(), caBundle.CASigner, defaults.CoreHostName, meta.ClusterUUID)
    if err != nil {
        // this is critical
        log.Debugf(err.Error())
        return nil, trace.Wrap(err)
    }
    context.UpdateHostCert(hostBundle)

    // make teleport core config
    cfg := service.MakeCoreConfig(dataDir, true)
    cfg.AssignHostUUID(meta.ClusterUUID)
    cfg.AssignDatabaseEngine(rec.DataBase())
    cfg.AssignCertStorage(rec.Certdb())
    cfg.AssignCASigner(caBundle.CASigner)
    cfg.AssignHostCertAuth(caBundle.CAPrvKey, caBundle.CASSHChk, meta.ClusterDomain)
    err = service.ValidateMasterConfig(cfg)
    if err != nil {
        log.Debugf(err.Error())
        return nil, trace.Wrap(err)
    }

    //log.Info(spew.Sdump(ctx))
    return cfg, nil
}
