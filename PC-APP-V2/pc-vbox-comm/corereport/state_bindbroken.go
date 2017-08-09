package corereport

import (
    "time"

    "github.com/pkg/errors"
    cpkg "github.com/stkim1/pc-vbox-comm/corereport/pkg"
    mpkg "github.com/stkim1/pc-vbox-comm/masterctrl/pkg"
    "github.com/stkim1/pc-vbox-core/crcontext"
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
    eni, err := crcontext.ExternalNetworkInterface()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    meta, err = cpkg.CorePackingBindBrokenStatus(core.clusterID, eni.IP4Address[0], eni.GatewayAddr, core.rsaEncryptor)
    return meta, errors.WithStack(err)
}

func (n *bindbroken) readMasterAck(core *coreReporter, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
    var (
        meta *mpkg.VBoxMasterMeta = nil
        err error = nil
    )

    meta, err = mpkg.MasterUnpackingAcknowledge(core.clusterID, metaPackage, core.rsaDecryptor)
    if err != nil {
        return VBoxCoreTransitionIdle, errors.WithStack(err)
    }
    err = crcontext.SharedCoreContext().SetMasterIP4ExtAddr(meta.MasterAcknowledge.ExtIP4Addr)
    if err != nil {
        return VBoxCoreTransitionIdle, errors.WithStack(err)
    }

    return VBoxCoreTransitionOk, nil
}

func (n *bindbroken) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    return nil
}

func (n *bindbroken) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return errors.WithStack(crcontext.SharedCoreContext().DiscardMasterSession())
}
