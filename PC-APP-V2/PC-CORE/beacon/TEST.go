package beacon

import (
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/stkim1/pcrypto"
    "github.com/stkim1/pc-core/model"
    "github.com/stkim1/pc-core/utils"
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

const (
    maxRandomSlaveIdLenth = 30
)

type DebugBeaconNotiReceiver struct {
    LastStateSuccessFrom     MasterBeaconState
    LastStateFailureFrom     MasterBeaconState
    Slave                    *model.SlaveNode
    SlaveNodes               []model.SlaveNode
    TransitionTS             time.Time
    IsShutdown               bool
}

func (d *DebugBeaconNotiReceiver) BeaconEventPrepareJoin(slave *model.SlaveNode) error {
    err := slave.SetAuthToken(utils.NewRandomString(maxRandomSlaveIdLenth))
    if err != nil {
        log.Debugf(err.Error())
    }
    return err
}

func (d *DebugBeaconNotiReceiver) BeaconEventResurrect(slaves []model.SlaveNode) error {
    d.SlaveNodes = slaves
    return nil
}

func (d *DebugBeaconNotiReceiver) BeaconEventTranstion(state MasterBeaconState, slave *model.SlaveNode, ts time.Time, transOk bool) error {
    if transOk {
        d.LastStateSuccessFrom = state
    } else {
        d.LastStateFailureFrom = state
    }
    d.Slave = slave
    d.TransitionTS = ts
    return nil
}

func (d *DebugBeaconNotiReceiver) BeaconEventDiscard(slave *model.SlaveNode) error {
    return nil
}

func (d *DebugBeaconNotiReceiver) BeaconEventShutdown() error {
    d.IsShutdown = true
    return nil
}
