package beacon

import (
    "time"

    "github.com/stkim1/pcrypto"
)

type DebugCommChannel struct {
    LastUcastMessage []byte
    LastUcastHost    string
    UCommCount       uint
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
    TransitionFailed() uint
    TxActionTS() time.Time
    TxActionFailed() uint
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

func (b *beaconState) TransitionFailed() uint {
    return b.transitionFailureCount
}

func (b *beaconState) TxActionTS() time.Time {
    return b.lastTransmissionTS
}

func (b *beaconState) TxActionFailed() uint {
    return b.txActionCount
}
