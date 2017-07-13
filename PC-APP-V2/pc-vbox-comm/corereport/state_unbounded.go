package corereport

import (
    "time"

    "github.com/pkg/errors"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
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
    core.authToken = meta.MasterAcknowledge.AuthToken

    return VBoxCoreTransitionOk, nil
}

func (u *unbounded) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    // TODO : save uuid to disk
    return nil
}

func (u *unbounded) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return nil
}
