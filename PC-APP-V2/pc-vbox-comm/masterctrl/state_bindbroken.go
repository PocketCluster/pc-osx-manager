package masterctrl

import (
    "time"

    "github.com/pkg/errors"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
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
    meta, err = cpkg.CoreUnpackingStatus(master.clusterID, metaPackage, master.rsaDecryptor)
    if err != nil {
        return VBoxMasterTransitionIdle, errors.WithStack(err)
    }
    if meta.CoreStatus.CoreState != cpkg.VBoxCoreBindBroken {
        return VBoxMasterTransitionIdle, errors.Errorf("[ERR] core state should be VBoxCoreBindBroken")
    }
    err = master.coreNode.UpdateIPv4WithGW(meta.CoreStatus.ExtIP4AddrSmask, meta.CoreStatus.ExtIP4Gateway)
    if err != nil {
        return VBoxMasterTransitionIdle, errors.WithMessage(err,"[ERR] cannot update core node with address and gateway")
    }

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
