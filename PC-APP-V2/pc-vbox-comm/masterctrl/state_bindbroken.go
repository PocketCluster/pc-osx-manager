package masterctrl

import (
    "time"

    "github.com/pkg/errors"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
)

type bindbroken struct {}
func stateBindbroken() vboxController { return &bindbroken {} }

func (n *bindbroken) currentState() mpkg.VBoxMasterState {
    return mpkg.VBoxMasterBindBroken
}

func (n *bindbroken) readCoreReport(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error) {
    var (
        meta *cpkg.VBoxCoreMeta
        err error = nil
    )

    // decrypt & update status package
    meta, err = cpkg.CoreUnpackingStatus(metaPackage, master.rsaDecryptor)
    if err != nil {
        return VBoxMasterTransitionIdle, errors.WithStack(err)
    }
    if meta.CoreState != cpkg.VBoxCoreBindBroken {
        return VBoxMasterTransitionIdle, errors.Errorf("[ERR] core state should be VBoxCoreBindBroken")
    }
    // TODO need lock
    master.coreNode.IP4Address = meta.CoreStatus.ExtIP4AddrSmask
    master.coreNode.IP4Gateway = meta.CoreStatus.ExtIP4Gateway

    return VBoxMasterTransitionOk, nil
}

func (n *bindbroken) makeMasterAck(master *masterControl, ts time.Time) ([]byte, error) {
    return nil, errors.Errorf("[ERR] VBoxMasterBindBroken cannot yield output")
}

func (n *bindbroken) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    return nil
}

func (n *bindbroken) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
