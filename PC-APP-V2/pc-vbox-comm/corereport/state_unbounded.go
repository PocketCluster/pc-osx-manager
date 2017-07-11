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

func (u *unbounded) transitionWithMasterMeta(core *coreReporter, sender interface{}, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
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

func (u *unbounded) transitionWithTimeStamp(core *coreReporter, ts time.Time) error {
    var (
        meta []byte = nil
        err error = nil
    )

    // send pubkey to master
    meta, err = cpkg.CorePackingStatus(cpkg.VBoxCoreUnbounded, core.publicKey, "", "", nil)
    if err != nil {
        return errors.WithStack(err)
    }
    return core.UcastSend("127.0.0.1", meta)
}

func (u *unbounded) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    return nil
}

func (u *unbounded) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return nil
}
