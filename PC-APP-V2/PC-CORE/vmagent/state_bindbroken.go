package vmagent

import (
    "time"

    "github.com/stkim1/pc-vbox-core/vcagent"
    "github.com/pkg/errors"
)

type bindbroken struct {}
func stateBindbroken() vboxController { return &bindbroken {} }

func (n *bindbroken) currentState() VBoxMasterState {
    return VBoxMasterBindBroken
}

func (n *bindbroken) transitionWithCoreMeta(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error) {
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

func (n *bindbroken) transitionWithTimeStamp(master *masterControl, ts time.Time) error {
    return nil
}

func (n *bindbroken) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    return nil
}

func (n *bindbroken) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
