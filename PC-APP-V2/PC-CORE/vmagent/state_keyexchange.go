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

    // decrypt & update status package
    status, err = vcagent.CoreDecryptBounded(metaPackage, master.rsaDecryptor)
    if err != nil {
        return VBoxMasterTransitionIdle, errors.WithStack(err)
    }
    master.coreNode.IP4Address = status.ExtIP4AddrSmask
    master.coreNode.IP4Gateway = status.ExtIP4Gateway

    return VBoxMasterTransitionOk, nil
}

func (k *keyexchange) transitionWithTimeStamp(master *masterControl, ts time.Time) error {
    var (
        authToken string = ""
        keypkg []byte = nil
        err error = nil
    )

    // send key exchange package
    authToken, err = master.coreNode.GetAuthToken()
    if err != nil {
        return errors.WithStack(err)
    }
    keypkg, err = MasterEncryptedKeyExchange(authToken, master.publicKey, master.rsaEncryptor)
    if err != nil {
        return errors.WithStack(err)
    }
    return master.UcastSend("127.0.0.1", keypkg)
}

func (k *keyexchange) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    return master.coreNode.JoinCore()
}

func (k *keyexchange) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
