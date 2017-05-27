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
    regisrv "github.com/stkim1/pc-core/extsrv/registry"
    swarmsrv "github.com/stkim1/pc-core/extsrv/swarm"
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
}

type serviceConfig struct {
    etcdConfig     *embed.PocketConfig
    teleConfig     *tervice.PocketConfig
    regConfig      *regisrv.PocketRegistryConfig
    swarmConfig    *swarmsrv.SwarmContext
}

func setupServiceConfig() (*serviceConfig, error) {
    // setup context
    ctx := context.SharedHostContext()
    context.SetupBasePath()

    // open database
    dataDir, err := ctx.ApplicationUserDataDirectory()
    if err != nil {
        log.Info(err)
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
    caBundle, err := certAuthSigner(rec.Certdb(), meta, country)
    if err != nil {
        // this is critical
        log.Debugf(err.Error())
        return nil, errors.WithStack(err)
    }
    context.UpdateCertAuth(caBundle)

    // host certificate
    hostBundle, err := hostCertificate(rec.Certdb(), caBundle.CASigner, tefaults.CoreHostName, meta.ClusterUUID)
    if err != nil {
        // this is critical
        log.Debugf(err.Error())
        return nil, errors.WithStack(err)
    }
    context.UpdateHostCert(hostBundle)

    // make teleport core config
    teleCfg := tervice.MakeCoreConfig(dataDir, true)
    teleCfg.AssignHostUUID(meta.ClusterUUID)
    teleCfg.AssignDatabaseEngine(rec.DataBase())
    teleCfg.AssignCertStorage(rec.Certdb())
    teleCfg.AssignCASigner(caBundle.CASigner)
    teleCfg.AssignHostCertAuth(caBundle.CAPrvKey, caBundle.CASSHChk, meta.ClusterDomain)
    err = tervice.ValidateCoreConfig(teleCfg)
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
