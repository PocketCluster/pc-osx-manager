package masterctrl

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pcrypto"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
)

type unbounded struct {}
func stateUnbounded() vboxController { return &unbounded{}}

func (u *unbounded) currentState() mpkg.VBoxMasterState {
    return mpkg.VBoxMasterUnbounded
}

func (u *unbounded) readCoreReport(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error) {
    var (
        meta *cpkg.VBoxCoreMeta = nil
        encryptor pcrypto.RsaEncryptor
        decryptor pcrypto.RsaDecryptor
        err error = nil
    )

    // unpack unbounded
    meta, err = cpkg.CoreUnpackingStatus(metaPackage, nil)
    if err != nil {
        return VBoxMasterTransitionIdle, errors.WithStack(err)
    }
    if meta.CoreState != cpkg.VBoxCoreUnbounded {
        return VBoxMasterTransitionIdle, errors.Errorf("[ERR] core state should be VBoxCoreUnbounded")
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

func (u *unbounded) makeMasterAck(master *masterControl, ts time.Time) ([]byte, error) {
    var (
        authToken string = ""
        keypkg []byte = nil
        err error = nil
    )

    // send key exchange package
    authToken, err = master.coreNode.GetAuthToken()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    keypkg, err = mpkg.MasterPackingKeyExchangeAcknowledge(authToken, master.publicKey, master.rsaEncryptor)
    return keypkg, errors.WithStack(err)
}

func (u *unbounded) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    return nil
}

func (u *unbounded) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
