package vcagent

import "time"

type unbounded struct {}
func stateUnbounded() vboxReporter { return &unbounded{} }

func (u *unbounded) currentState() VBoxCoreState {
    return VBoxCoreUnbounded
}

func (u *unbounded) transitionWithMasterMeta(core *coreReporter, sender interface{}, metaPackage []byte, ts time.Time) (VBoxCoreTransition, error) {
    return VBoxCoreTransitionOk, nil
}

func (u *unbounded) transitionWithTimeStamp(core *coreReporter, ts time.Time) error {
    return nil
}

func (u *unbounded) onStateTranstionSuccess(core *coreReporter, ts time.Time) error {
    return nil
}

func (u *unbounded) onStateTranstionFailure(core *coreReporter, ts time.Time) error {
    return nil
}
