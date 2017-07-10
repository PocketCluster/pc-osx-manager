package masterctrl

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
)

type unbounded struct {}
func stateUnbounded() vboxController { return &unbounded{}}

func (u *unbounded) currentState() VBoxMasterState {
    return VBoxMasterUnbounded
}

func (u *unbounded) transitionWithCoreMeta(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error) {
    var (
        meta *cpkg.VBoxCoreAgentMeta = nil
        encryptor pcrypto.RsaEncryptor
        decryptor pcrypto.RsaDecryptor
        err error = nil
    )

    // unpack unbounded
    meta, err = cpkg.CoreUnpackingUnbounded(metaPackage)
    if err != nil {
        return VBoxMasterTransitionIdle, errors.WithStack(err)
    }

    // encryptor, decryptor
    encryptor, err = pcrypto.NewRsaEncryptorFromKeyData(meta.PublicKey, master.privateKey)
    if err != nil {
        return VBoxMasterTransitionFail, errors.WithStack(err)
    }
    decryptor, err = pcrypto.NewRsaDecryptorFromKeyData(meta.PublicKey, master.privateKey)
    if err != nil {
        return VBoxMasterTransitionFail, errors.WithStack(err)
    }
    master.rsaEncryptor = encryptor
    master.rsaDecryptor = decryptor

    // save core node public key
    master.coreNode.PublicKey = meta.PublicKey

    return VBoxMasterTransitionOk, nil
}

func (u *unbounded) transitionWithTimeStamp(master *masterControl, ts time.Time) error {
    return nil
}

func (u *unbounded) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    return nil
}

func (u *unbounded) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
