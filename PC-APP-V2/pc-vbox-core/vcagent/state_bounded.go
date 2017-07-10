package vcagent

import "time"

type bounded struct {}
func stateBounded() vboxReporter { return &bounded{} }

func (b *bounded) currentState() VBoxCoreState {
    return VBoxCoreBounded
}

func (b *bounded) transitionWithMasterMeta(core *coreReporter, sender interface{}, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
    return VBoxCoreTransitionOk, nil
}

func (b *bounded) transitionWithTimeStamp(core *coreReporter, ts time.Time) error {
    return nil
}

func (b *bounded) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    return nil
}

func (b *bounded) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return nil
}
