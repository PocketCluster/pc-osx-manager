package vmagent

import (
    "time"

    "github.com/stkim1/pc-vbox-core/vcagent"
    "github.com/pkg/errors"
)

type bounded struct {
}

func stateBounded () vboxController {
    return &bounded{
    }
}

func (b *bounded) currentState() VBoxMasterState {
    return VBoxMasterBounded
}

func (b *bounded) transitionWithCoreMeta(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error) {
    var (
        //meta *vcagent.VBoxCoreAgentMeta = nil
        err error = nil
    )

    _, err = vcagent.CoreDecryptBounded(metaPackage, master.rsaDecryptor)
    if err != nil {
        return VBoxMasterTransitionFail, errors.WithStack(err)
    }

    return VBoxMasterTransitionOk, nil
}

func (b *bounded) transitionWithTimeStamp(master *masterControl, ts time.Time) error {
    return nil
}

func (b *bounded) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    return nil
}

func (b *bounded) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
