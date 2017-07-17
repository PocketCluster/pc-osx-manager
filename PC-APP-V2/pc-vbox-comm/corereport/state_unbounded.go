package corereport

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"

    "github.com/stkim1/pc-vbox-core/crcontext"
)

type unbounded struct {}
func stateUnbounded() vboxReporter { return &unbounded{} }

func (u *unbounded) currentState() cpkg.VBoxCoreState {
    return cpkg.VBoxCoreUnbounded
}

func (u *unbounded) makeCoreReport(core *coreReporter, ts time.Time) ([]byte, error) {
    var (
        meta []byte = nil
        err error = nil
    )

    // send pubkey to master
    meta, err = cpkg.CorePackingUnboundedStatus(core.publicKey)
    return meta, errors.WithStack(err)
}

func (u *unbounded) readMasterAck(core *coreReporter, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
    var (
        meta *mpkg.VBoxMasterMeta = nil
        err error = nil
    )

    // get acknowledge, encryptor, & decryptor
    meta, err = mpkg.MasterUnpackingAcknowledge(metaPackage, core.privateKey, nil)
    if err != nil {
        return VBoxCoreTransitionIdle, errors.WithStack(err)
    }
    core.rsaEncryptor = meta.Encryptor
    core.rsaDecryptor = meta.Decryptor
    core.masterPubKey = meta.PublicKey
    core.clusterID = meta.MasterAcknowledge.ClusterID
    core.authToken = meta.MasterAcknowledge.AuthToken
    core.masterExtIp4Addr = meta.MasterAcknowledge.ExtIP4Addr

    return VBoxCoreTransitionOk, nil
}

func (u *unbounded) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    var err error = nil

    err = crcontext.SharedCoreContext().SetMasterPublicKey(core.masterPubKey)
    if err != nil {
        log.Debugf("[ERR] master public key saving error %v", errors.WithStack(err).Error())
    }
    err = crcontext.SharedCoreContext().SetClusterID(core.clusterID)
    if err != nil {
        log.Debugf("[ERR] cluster id saving error %v", errors.WithStack(err).Error())
    }
    err = crcontext.SharedCoreContext().SetCoreAuthToken(core.authToken)
    if err != nil {
        log.Debugf("[ERR] auth token saving error %v", errors.WithStack(err).Error())
    }
    err = crcontext.SharedCoreContext().SetMasterIP4ExtAddr(core.masterExtIp4Addr)
    if err != nil {
        log.Debugf("[ERR] master external ip4 addr saving error %v", errors.WithStack(err).Error())
    }
    err = crcontext.SharedCoreContext().SaveConfiguration()
    if err != nil {
        log.Debugf("[ERR] sync shared context error %v", errors.WithStack(err).Error())
    }
    return nil
}

func (u *unbounded) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return errors.WithStack(crcontext.SharedCoreContext().DiscardAll())
}
