package vmagent

import (
    "time"

    "github.com/stkim1/pc-vbox-core/vcagent"
    "github.com/pkg/errors"
)

type keyexchange struct {
}

func stateKeyexchange() vboxController {
    return &keyexchange {
    }
}

func (k *keyexchange) currentState() VBoxMasterState {
    return VBoxMasterKeyExchange
}

func (k *keyexchange) transitionWithCoreMeta(master *masterControl, sender interface{}, metaPackage []byte, ts time.Time) (VBoxMasterTransition, error) {
    var (
        //meta *vcagent.VBoxCoreAgentMeta = nil
        err error = nil
    )

    _, err = vcagent.CoreUnpackingUnbounded(metaPackage)
    if err != nil {
        return VBoxMasterTransitionFail, errors.WithStack(err)
    }

    return VBoxMasterTransitionOk, nil
}

func (k *keyexchange) transitionWithTimeStamp(master *masterControl, ts time.Time) error {
    return nil
}

func (k *keyexchange) onStateTranstionSuccess(master *masterControl, ts time.Time) error {
    return nil
}

func (k *keyexchange) onStateTranstionFailure(master *masterControl, ts time.Time) error {
    return nil
}
