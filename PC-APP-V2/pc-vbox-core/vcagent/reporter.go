package vcagent

import "time"

type VBoxCoreState int
const (
    VBoxCoreUnbounded     VBoxCoreState = iota
    VBoxCoreBounded
    VBoxCoreBindBroken
)

type VBoxCoreTransition int
const (
    VBoxCoreTransitionFail    VBoxCoreTransition = iota
    VBoxCoreTransitionOk
    VBoxCoreTransitionIdle
)

func (s VBoxCoreState) String() string {
    var state string
    switch s {
        case VBoxCoreUnbounded:
            state = "VBoxCoreUnbounded"
        case VBoxCoreBounded:
            state = "VBoxCoreBounded"
        case VBoxCoreBindBroken:
            state = "VBoxCoreBindBroken"
    }
    return state
}

type CommChannel interface {
    //McastSend(data []byte) error
    UcastSend(target string, data []byte) error
}

type CommChannelFunc func(target string, data []byte) error
func (c CommChannelFunc) UcastSend(target string, data []byte) error {
    return c(target, data)
}

// MasterBeacon is assigned individually for each slave node.
type VBoxCoreReporter interface {
    CurrentState() VBoxCoreState
    TransitionWithTimestamp(timestamp time.Time) error
    TransitionWithCoreMeta(sender interface{}, metaPackage []byte, timestamp time.Time) error
}

type coreReporter struct {
    state    VBoxCoreState
}

func (c *coreReporter) CurrentState() VBoxCoreState {
    return c.state
}

func (c *coreReporter) TransitionWithTimestamp(timestamp time.Time) error {
    return nil
}

func (c *coreReporter) TransitionWithCoreMeta(sender interface{}, metaPackage []byte, timestamp time.Time) error {
    return nil
}
