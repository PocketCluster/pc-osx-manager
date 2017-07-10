package vcagent

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    mpkg "github.com/stkim1/pc-core/vmagent/pkg"
    cpkg "github.com/stkim1/pc-vbox-core/vcagent/pkg"
)

type unbounded struct {}
func stateUnbounded() vboxReporter { return &unbounded{} }

func (u *unbounded) currentState() VBoxCoreState {
    return VBoxCoreUnbounded
}

func (u *unbounded) transitionWithMasterMeta(core *coreReporter, sender interface{}, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
    var (
        ack *mpkg.VBoxMasterAcknowledge = nil
        encryptor pcrypto.RsaEncryptor
        decryptor pcrypto.RsaDecryptor
        err error = nil
    )

    // get acknowledge, encryptor, & decryptor
    ack, encryptor, decryptor, err = mpkg.MasterDecryptedKeyExchange(metaPackage, core.privateKey)
    if err != nil {
        return VBoxCoreTransitionIdle, errors.WithStack(err)
    }
    core.rsaEncryptor = encryptor
    core.rsaDecryptor = decryptor
    core.authToken = ack.AuthToken

    return VBoxCoreTransitionOk, nil
}

func (u *unbounded) transitionWithTimeStamp(core *coreReporter, ts time.Time) error {
    var (
        meta []byte = nil
        err error = nil
    )

    // send pubkey to master
    meta, err = cpkg.CorePackingUnbounded(core.publicKey)
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
