package discovery

import (
    "time"
    "fmt"
)

func stateTransition(currState SDState, nextCondition func() (SDTranstion)) (nextState SDState, err error) {
    var transtion SDTranstion = nextCondition()

    // Succeed to transition to the next
    if  transtion == SlaveTransitionOk {
        switch currState {
        case SlaveUnbounded:
            nextState = SlaveInquired
        case SlaveInquired:
            nextState = SlaveKeyExchange
        case SlaveKeyExchange:
            nextState = SlaveCryptoCheck
        case SlaveCryptoCheck:
            nextState = SlaveBounded
        case SlaveBindBroken:
            nextState = SlaveBounded
        case SlaveBounded:
            nextState = currState
        default:
            err = fmt.Errorf("[PANIC] 'nextCondition is true and hit default' cannot happen")
        }
        // Fail to transition to the next
    } else if transtion == SlaveTransitionFail {
        switch currState {
        case SlaveUnbounded:
            nextState = SlaveUnbounded
        case SlaveInquired:
            nextState = SlaveUnbounded
        case SlaveKeyExchange:
            nextState = SlaveUnbounded
        case SlaveCryptoCheck:
            nextState = SlaveUnbounded
        case SlaveBindBroken:
            nextState = SlaveBindBroken
        case SlaveBounded:
            nextState = SlaveBindBroken
        default:
            err = fmt.Errorf("[PANIC] 'nextCondition is false and hit default' cannot happen")
        }
        // Idle
    } else {
        nextState = currState
    }
    return
}

func (sd *slaveDiscovery) Unbounded(timestamp *time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) Inquired(timestamp *time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) KeyExchange(timestamp *time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) CryptoCheck(timestamp *time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) Bounded(timestamp *time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}

func (sd *slaveDiscovery) BindBroken(timestamp *time.Time) (err error) {
    state, err := stateTransition(sd.discoveryState, func() SDTranstion {
        return SlaveTransitionOk
    })
    if err != nil {
        return err
    }
    sd.discoveryState = state
    return
}
