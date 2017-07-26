package corereport

import (
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-vbox-core/crcontext"
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
    eni, err := crcontext.ExternalNetworkInterface()
    if err != nil {
        return nil, errors.WithStack(err)
    }

    meta, err = cpkg.CorePackingBoundedStatus(core.clusterID, eni.IP4Address[0], eni.GatewayAddr, core.rsaEncryptor)
    return meta, errors.WithStack(err)
}

func (b *bounded) readMasterAck(core *coreReporter, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
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

func (b *bounded) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    return nil
}

func (b *bounded) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return errors.WithStack(crcontext.SharedCoreContext().DiscardMasterSession())
}
