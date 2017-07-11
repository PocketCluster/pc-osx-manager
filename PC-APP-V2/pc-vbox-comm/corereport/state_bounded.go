package corereport

import (
    "time"

    "github.com/pkg/errors"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
)

type bounded struct {}
func stateBounded() vboxReporter { return &bounded{} }

func (b *bounded) currentState() cpkg.VBoxCoreState {
    return cpkg.VBoxCoreBounded
}

func (b *bounded) makeCoreReport(core *coreReporter, ts time.Time) ([]byte, error) {
    var (
        meta []byte = nil
        err error = nil
    )

    // send status to master
    // TODO get ip address and gateway
    meta, err = cpkg.CorePackingStatus(cpkg.VBoxCoreBounded, nil, "127.0.0.1", "192.168.1.1", core.rsaEncryptor)
    return meta, errors.WithStack(err)
}

func (b *bounded) readMasterAck(core *coreReporter, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
    var (
        err error = nil
    )

    _, err = mpkg.MasterUnpackingAcknowledge(metaPackage, nil, core.rsaDecryptor)
    if err != nil {
        return VBoxCoreTransitionIdle, errors.WithStack(err)
    }

    return VBoxCoreTransitionOk, nil
}

func (b *bounded) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    return nil
}

func (b *bounded) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return nil
}
