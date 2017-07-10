package vcagent

import "time"

type bindbroken struct {}
func stateBindbroken() vboxReporter { return &bindbroken{} }

func (n *bindbroken) currentState() VBoxCoreState {
    return VBoxCoreBindBroken
}

func (n *bindbroken) transitionWithMasterMeta(core *coreReporter, sender interface{}, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
    return VBoxCoreTransitionOk, nil
}

func (n *bindbroken) transitionWithTimeStamp(core *coreReporter, ts time.Time) error {
    return nil
}

func (n *bindbroken) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    return nil
}

func (n *bindbroken) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return nil
}
