package slagent

import (
    "runtime"
    "time"

    "gopkg.in/vmihailenco/msgpack.v2"
)

type PocketSlaveStatus struct {
    Version             StatusProtocol  `msgpack:"s_ps"`
    // slave response
    SlaveResponse       ResponseType    `msgpack:"s_rt,omitempty`
    // slave nodename
    SlaveNodeName       string          `msgpack:"s_nm,omitempty"`
    // slave UUID
    SlaveUUID           string          `msgpack:"s_uu,omitempty"`
    SlaveHardware       string          `msgpack:"s_hw"`
    SlaveTimestamp      time.Time       `msgpack:"s_ts"`
}

func (ssa *PocketSlaveStatus) IsAppropriateSlaveInfo() bool {
    if len(ssa.SlaveHardware) == 0 {
        return false
    }
    return true
}

func PackedSlaveStatus(status *PocketSlaveStatus) ([]byte, error) {
    return msgpack.Marshal(status)
}

func UnpackedSlaveStatus(message []byte) (status *PocketSlaveStatus, err error) {
    err = msgpack.Unmarshal(message, &status)
    return
}

// Unbounded
func AnswerMasterInquiryStatus(timestamp time.Time) (*PocketSlaveStatus, error) {
    return &PocketSlaveStatus {
        Version:           SLAVE_STATUS_VERSION,
        SlaveResponse:     SLAVE_WHO_I_AM,
        SlaveHardware:     runtime.GOARCH,
        SlaveTimestamp:    timestamp,
    }, nil
}

func KeyExchangeStatus(timestamp time.Time) (*PocketSlaveStatus, error) {
    return &PocketSlaveStatus {
        Version:           SLAVE_STATUS_VERSION,
        SlaveResponse:     SLAVE_SEND_PUBKEY,
        SlaveHardware:     runtime.GOARCH,
        SlaveTimestamp:    timestamp,
    }, nil
}

func CheckSlaveCryptoStatus(nodename, uuid string, timestamp time.Time) (*PocketSlaveStatus, error) {
    return &PocketSlaveStatus {
        Version:           SLAVE_STATUS_VERSION,
        SlaveResponse:     SLAVE_CHECK_CRYPTO,
        SlaveNodeName:     nodename,
        SlaveUUID:         uuid,
        SlaveHardware:     runtime.GOARCH,
        SlaveTimestamp:    timestamp,
    }, nil
}

func SlaveBoundedStatus(nodename, uuid string, timestamp time.Time) (*PocketSlaveStatus, error) {
    return &PocketSlaveStatus {
        Version:           SLAVE_STATUS_VERSION,
        SlaveResponse:     SLAVE_REPORT_STATUS,
        SlaveNodeName:     nodename,
        SlaveUUID:         uuid,
        SlaveHardware:     runtime.GOARCH,
        SlaveTimestamp:    timestamp,
    }, nil
}
