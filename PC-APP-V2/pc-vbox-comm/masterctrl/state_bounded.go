package masterctrl

import (
    "time"

    "github.com/pkg/errors"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
)

type bounded struct {}
func stateBounded () vboxController { return &bounded{} }

func (b *bounded) currentState() mpkg.VBoxMasterState {
    return mpkg.VBoxMasterBounded
}

func (b *bounded) readCoreReport(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error) {
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

func (b *bounded) makeMasterAck(master *masterControl, ts time.Time) ([]byte, error) {
    var (
        ackpkg []byte = nil
        err error = nil
    )

    // send acknowledge packagepackage
    ackpkg, err = mpkg.MasterPackingBoundedAcknowledge(master.clusterID, master.extIP4Addr, master.rsaEncryptor)
    return ackpkg, errors.WithStack(err)
}

func (b *bounded) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    return nil
}

func (b *bounded) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
