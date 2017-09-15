package config

import (
    tefaults "github.com/gravitational/teleport/lib/defaults"
    tervice "github.com/gravitational/teleport/lib/service"
    "github.com/coreos/etcd/embed"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"

    "github.com/stkim1/pc-core/context"
    "github.com/stkim1/pc-core/extlib/registry"
    "github.com/stkim1/pc-core/extlib/pcssh/sshcfg"
    "github.com/stkim1/pc-core/model"
)

type ServiceConfig struct {
    ETCD  *embed.PocketConfig
    PCSSH *tervice.PocketConfig
    REG   *registry.PocketRegistryConfig
}

func InitServiceConfig() (*ServiceConfig, error) {
    // setup context
    ctx := context.SharedHostContext()
    err := context.SetupBasePath()
    if err != nil {
        // this is critical
        return nil, errors.WithStack(err)
    }

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
    ctx.SetClusterMeta(meta)

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
    teleCfg, err := sshcfg.MakeMasterConfig(context.SharedHostContext(), true)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    sshcfg.AssignDatabaseEngine(teleCfg, rec.DataBase())
    sshcfg.AssignCertStorage(teleCfg, rec.Certdb())
    sshcfg.AssignCASigner(teleCfg, caBundle.CASigner)
    sshcfg.AssignHostCertAuth(teleCfg, caBundle.CAPrvKey, caBundle.CASSHChk, meta.ClusterDomain)
    err = sshcfg.ValidateMasterConfig(teleCfg)
    if err != nil {
        return nil, errors.WithStack(err)
    }

    // registry configuration
    regPath, err := context.SharedHostContext().ApplicationRepositoryDirectory()
    if err != nil {
        return nil, errors.WithStack(err)
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
    etcdPath, err := context.SharedHostContext().ApplicationStorageDirectory()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    // recommended parameter values
    // heartbeat : 500
    // election  : 5000 ( heartbeat * 10 )
    // snapshot  : 1000
    // TODO : these parameters need to be dynamically adjusted according to a cluster condition
    etcdCfg, err := embed.NewPocketConfig(etcdPath, caBundle.CACrtPem, hostBundle.Certificate, hostBundle.PrivateKey, 500, 5000, 1000, true)
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &ServiceConfig{
        ETCD:  etcdCfg,
        PCSSH: teleCfg,
        REG:   regCfg,
    }, nil
}