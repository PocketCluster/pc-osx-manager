package main

import (
    "os"
    "path"

    log "github.com/Sirupsen/logrus"
    tefaults "github.com/gravitational/teleport/lib/defaults"
    tervice "github.com/gravitational/teleport/lib/service"

    "github.com/coreos/etcd/embed"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-core/defaults"
    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/extlib/registry"
    "github.com/stkim1/pc-core/extlib/pcssh/sshcfg"
)

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

    // vbox certificate
    vboxBundle, err := buildVBoxReportCertificate(rec.Certdb(), meta.ClusterUUID)
    if err != nil {
        // this is critical
        return nil, errors.WithStack(err)
    }
    ctx.UpdateVBoxCert(vboxBundle)

    // make teleport core config
    teleCfg := sshcfg.MakeMasterConfig(dataDir, true)
    sshcfg.AssignHostUUID(teleCfg, meta.ClusterUUID)
    sshcfg.AssignDatabaseEngine(teleCfg, rec.DataBase())
    sshcfg.AssignCertStorage(teleCfg, rec.Certdb())
    sshcfg.AssignCASigner(teleCfg, caBundle.CASigner)
    sshcfg.AssignHostCertAuth(teleCfg, caBundle.CAPrvKey, caBundle.CASSHChk, meta.ClusterDomain)
    err = sshcfg.ValidateMasterConfig(teleCfg)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // registry configuration
    var regPath = path.Join(dataDir, defaults.RepositoryPathPostfix)
    if _, err := os.Stat(regPath); os.IsNotExist(err) {
        os.MkdirAll(path.Join(regPath, "docker/registry/v2/repositories"), 0700)
        os.MkdirAll(path.Join(regPath, "docker/registry/v2/blobs"),        0700)
    }
    regCfg, err := registry.NewPocketRegistryConfig(false, regPath, hostBundle.Certificate, hostBundle.PrivateKey)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    err = registry.GarbageCollection(regCfg)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    //etcd configuration
    var etcdPath = path.Join(dataDir, defaults.StoragePathPostfix)
    if _, err := os.Stat(etcdPath); os.IsNotExist(err) {
        os.MkdirAll(etcdPath, 0700)
    }
    // recommended parameter values
    // heartbeat : 500
    // election  : 5000 ( heartbeat * 10 )
    // snapshot  : 1000
    // TODO : these parameters need to be dynamically adjusted according to a cluster condition
    etcdCfg, err := embed.NewPocketConfig(etcdPath, caBundle.CACrtPem, hostBundle.Certificate, hostBundle.PrivateKey, 500, 5000, 1000)
    if err != nil {
        // this is critical
        return nil, errors.WithStack(err)
    }
    return &serviceConfig {
        etcdConfig: etcdCfg,
        teleConfig: teleCfg,
        regConfig: regCfg,
    }, nil
}