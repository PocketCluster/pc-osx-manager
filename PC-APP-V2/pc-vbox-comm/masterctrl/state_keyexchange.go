package masterctrl

import (
    "time"

    "github.com/pkg/errors"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
)

type keyexchange struct {}
func stateKeyexchange() vboxController { return &keyexchange {} }

func (k *keyexchange) currentState() mpkg.VBoxMasterState {
    return mpkg.VBoxMasterKeyExchange
}

func (k *keyexchange) readCoreReport(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error) {
    var (
        meta *cpkg.VBoxCoreMeta
        err error = nil
    )

    // decrypt & update status package
    meta, err = cpkg.CoreUnpackingStatus(metaPackage, master.rsaDecryptor)
    if err != nil {
        return VBoxMasterTransitionIdle, errors.WithStack(err)
    }
    if meta.CoreState != cpkg.VBoxCoreBounded {
        return VBoxMasterTransitionIdle, errors.Errorf("[ERR] core state should be VBoxCoreBounded")
    }
    // TODO need lock
    master.coreNode.IP4Address = meta.CoreStatus.ExtIP4AddrSmask
    master.coreNode.IP4Gateway = meta.CoreStatus.ExtIP4Gateway

    return VBoxMasterTransitionOk, nil
}

func (k *keyexchange) makeMasterAck(master *masterControl, ts time.Time) ([]byte, error) {
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

func (k *keyexchange) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    return master.coreNode.JoinCore()
}

func (k *keyexchange) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
