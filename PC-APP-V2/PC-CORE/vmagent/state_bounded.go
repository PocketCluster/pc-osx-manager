package vmagent

import (
    "time"

    "github.com/stkim1/pc-vbox-core/vcagent"
    "github.com/pkg/errors"
)

type bounded struct {}
func stateBounded () vboxController { return &bounded{} }

func (b *bounded) currentState() VBoxMasterState {
    return VBoxMasterBounded
}

func (b *bounded) transitionWithCoreMeta(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error) {
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

func (b *bounded) transitionWithTimeStamp(master *masterControl, ts time.Time) error {
    var (
        ackpkg []byte = nil
        err error = nil
    )

    // send acknowledge packagepackage
    ackpkg, err = MasterEncryptedBounded(master.rsaEncryptor)
    if err != nil {
        return errors.WithStack(err)
    }
    return master.UcastSend("127.0.0.1", ackpkg)
}

func (b *bounded) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    return nil
}

func (b *bounded) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
