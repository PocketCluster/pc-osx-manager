package slagent

import (
    "runtime"
    "time"

    "github.com/pkg/errors"
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

// this is for master bindbroken state. Since majority of sanity check is done by beacon.bindbroken module, we'll just check simple things.
func ConvertDiscoveryToStatus(discovery *PocketSlaveDiscovery, slaveNode, slaveUUID, slaveHardware string) (*PocketSlaveStatus, error) {
    if len(slaveNode) == 0 {
        return nil, errors.Errorf("[ERR] incorrect slave name")
    }
    if len(slaveUUID) == 0 {
        return nil, errors.Errorf("[ERR] incorrect slave uuid")
    }
    if len(slaveHardware) == 0 {
        return nil, errors.Errorf("[ERR] incorrect slave hardware architecture")
    }
    if discovery.Version != SLAVE_DISCOVER_VERSION {
        return nil, errors.Errorf("[ERR] Incorrect SlaveDiscoveryAgent version")
    }
    if discovery.SlaveResponse != SLAVE_LOOKUP_AGENT {
        return nil, errors.Errorf("[ERR] incorrect slave discovery response")
    }
    if len(discovery.SlaveAddress) == 0 {
        return nil, errors.Errorf("[ERR] incorrect slave address")
    }
    if len(discovery.SlaveGateway) == 0 {
        return nil, errors.Errorf("[ERR] incorrect slave gateway")
    }
    // TODO : since discovery agent does not have timestamp, we'll use master timstamp.
    return &PocketSlaveStatus{
        Version:           SLAVE_STATUS_VERSION,
        SlaveResponse:     SLAVE_REPORT_STATUS,
        SlaveNodeName:     slaveNode,
        SlaveUUID:         slaveUUID,
        SlaveHardware:     slaveHardware,
        SlaveTimestamp:    time.Now(),
    }, nil
}