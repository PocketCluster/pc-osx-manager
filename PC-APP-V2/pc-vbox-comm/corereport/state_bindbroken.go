package corereport

import (
    "time"

    "github.com/pkg/errors"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
)

type bindbroken struct {}
func stateBindbroken() vboxReporter { return &bindbroken{} }

func (n *bindbroken) currentState() cpkg.VBoxCoreState {
    return cpkg.VBoxCoreBindBroken
}

func (n *bindbroken) makeCoreReport(core *coreReporter, ts time.Time) ([]byte, error) {
    var (
        meta []byte = nil
        err error = nil
    )

    // send status to master
    // TODO get ip address and gateway
    meta, err = cpkg.CorePackingBindBrokenStatus("127.0.0.1", "192.168.1.1", core.rsaEncryptor)
    return meta, errors.WithStack(err)
}

func (n *bindbroken) readMasterAck(core *coreReporter, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
    var (
        err error = nil
    )

    _, err = mpkg.MasterUnpackingAcknowledge(metaPackage, nil, core.rsaDecryptor)
    if err != nil {
        return VBoxCoreTransitionIdle, errors.WithStack(err)
    }

    return VBoxCoreTransitionOk, nil
}

func (n *bindbroken) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    return nil
}

func (n *bindbroken) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return nil
}
