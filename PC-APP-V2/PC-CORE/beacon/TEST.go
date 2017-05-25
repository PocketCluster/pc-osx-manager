package beacon

import (
    "time"

    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/model"
)

type DebugCommChannel struct {
    LastUcastMessage []byte
    LastUcastHost    string
    UCommCount       int
}

func (dc *DebugCommChannel) UcastSend(target string, data []byte) error {
    dc.LastUcastMessage = data
    dc.LastUcastHost = target
    dc.UCommCount++
    return nil
}

type DebugState interface {
    AESKey() []byte
    AESCryptor() pcrypto.AESCryptor
    TransitionSuccessTS() time.Time
    TransitionFailed() int
    TxActionTS() time.Time
    TxActionFailed() int
}

func (b *beaconState) AESKey() []byte {
    return b.aesKey
}

func (b *beaconState) AESCryptor() pcrypto.AESCryptor {
    return b.aesCryptor
}

func (b *beaconState) TransitionSuccessTS() time.Time {
    return b.lastTransitionTS
}

func (b *beaconState) TransitionFailed() int {
    return b.transitionFailureCount
}

func (b *beaconState) TxActionTS() time.Time {
    return b.lastTransmissionTS
}

func (b *beaconState) TxActionFailed() int {
    return b.txActionCount
}

type DebugTransitionEventReceiver struct {
    LastStateSuccessFrom     MasterBeaconState
    LastStateFailureFrom     MasterBeaconState
    Slave                    *model.SlaveNode
    TransitionTS             time.Time
}

func (d *DebugTransitionEventReceiver) OnStateTranstionSuccess(state MasterBeaconState, slave *model.SlaveNode, ts time.Time) error {
    d.LastStateSuccessFrom = state
    d.Slave = slave
    d.TransitionTS = ts
    return nil
}

func (d *DebugTransitionEventReceiver) OnStateTranstionFailure(state MasterBeaconState, slave *model.SlaveNode, ts time.Time) error {
    d.LastStateFailureFrom = state
    d.Slave = slave
    d.TransitionTS = ts
    return nil
}
