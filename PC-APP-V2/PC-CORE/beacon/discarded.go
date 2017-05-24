package beacon

import (
    "net"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
)

func discardedState(oldState *beaconState) BeaconState {
    b := &discarded{}

    b.constState                    = MasterDiscarded

    b.constTransitionFailureLimit   = TransitionFailureLimit
    b.constTransitionTimeout        = BoundedTimeout * time.Duration(TxActionLimit)
    b.constTxActionLimit            = TxActionLimit
    b.constTxTimeWindow             = BoundedTimeout

    b.lastTransitionTS              = time.Now()
    b.txActionCount                 = TxActionLimit + 1

    b.timestampTransition           = b.transitionActionWithTimestamp
    b.slaveMetaTransition           = b.transitionWithSlaveMeta

    b.BeaconOnTransitionEvent       = oldState.BeaconOnTransitionEvent
    b.aesKey                        = oldState.aesKey
    b.aesCryptor                    = oldState.aesCryptor
    b.rsaEncryptor                  = oldState.rsaEncryptor
    b.slaveNode                     = oldState.slaveNode
    b.commChan                      = nil

    return b
}

type discarded struct {
    beaconState
}

func (b *discarded) transitionActionWithTimestamp(masterTimestamp time.Time) error {
    return nil
}

func (b *discarded) transitionWithSlaveMeta(sender *net.UDPAddr, meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    return MasterTransitionOk, nil
}
