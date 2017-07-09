package vmagent

import (
    "time"

    "github.com/stkim1/pc-vbox-core/vcagent"
    "github.com/pkg/errors"
)

type keyexchange struct {}
func stateKeyexchange() vboxController { return &keyexchange {} }

func (k *keyexchange) currentState() VBoxMasterState {
    return VBoxMasterKeyExchange
}

func (k *keyexchange) transitionWithCoreMeta(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error) {
    var (
        status *vcagent.VBoxCoreStatus
        err error = nil
    )

    // decrypt status package
    status, err = vcagent.CoreDecryptBounded(metaPackage, master.rsaDecryptor)
    if err != nil {
        return VBoxMasterTransitionIdle, errors.WithStack(err)
    }

    // TODO assign core node ip and share
    master.coreNode = status

    return VBoxMasterTransitionOk, nil
}

func (k *keyexchange) transitionWithTimeStamp(master *masterControl, ts time.Time) error {
    var (
        keypkg []byte = nil
        err error = nil
    )

    // TODO : master corenode UUID
    keypkg, err = MasterEncryptedKeyExchange(master.coreNode, master.publicKey, master.rsaEncryptor)
    if err != nil {
        return errors.WithStack(err)
    }

    // send key exchange package

    return nil
}

func (k *keyexchange) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    // save core node public key

    // share core node external ip address

    return nil
}

func (k *keyexchange) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
