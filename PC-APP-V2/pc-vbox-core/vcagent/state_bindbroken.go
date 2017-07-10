package vcagent

import (
    "time"

    "github.com/pkg/errors"
    mpkg "github.com/stkim1/pc-core/vmagent/pkg"
    cpkg "github.com/stkim1/pc-vbox-core/vcagent/pkg"
)

type bindbroken struct {}
func stateBindbroken() vboxReporter { return &bindbroken{} }

func (n *bindbroken) currentState() VBoxCoreState {
    return VBoxCoreBindBroken
}

func (n *bindbroken) transitionWithMasterMeta(core *coreReporter, sender interface{}, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
    var (
        err error = nil
    )

    _, err = mpkg.MasterDecryptedBounded(metaPackage, core.rsaDecryptor)
    if err != nil {
        return VBoxCoreTransitionIdle, errors.WithStack(err)
    }

    return VBoxCoreTransitionOk, nil
}

func (n *bindbroken) transitionWithTimeStamp(core *coreReporter, ts time.Time) error {
    var (
        meta []byte = nil
        err error = nil
    )

    // send status to master
    // TODO get ip address and gateway
    meta, err = cpkg.CoreEncryptedBounded("127.0.0.1", "192.168.1.1", core.rsaEncryptor)
    if err != nil {
        return errors.WithStack(err)
    }
    return core.UcastSend("127.0.0.1", meta)
}

func (n *bindbroken) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    return nil
}

func (n *bindbroken) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return nil
}
