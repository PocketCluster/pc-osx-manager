package slagent

import (
    "runtime"
    "os"
    "time"

    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/stkim1/pc-node-agent/slcontext"
)

type PocketSlaveStatus struct {
    Version             StatusProtocol  `msgpack:"pc_sl_ps"`
    // master
    MasterBoundAgent    string          `msgpack:"pc_ms_ba,omitempty"`
    // slave response
    SlaveResponse       ResponseType    `msgpack:"pc_sl_rt,omitempty`
    // slave
    SlaveNodeName       string          `msgpack:"pc_sl_nm,omitempty"`
    // current interface status
    SlaveAddress        string          `msgpack:"pc_sl_i4"`
    SlaveNodeMacAddr    string          `msgpack:"pc_sl_ma"`
    SlaveHardware       string          `msgpack:"pc_sl_hw"`
    SlaveTimestamp      time.Time       `msgpack:"pc_sl_ts"`
}

func (ssa *PocketSlaveStatus) IsAppropriateSlaveInfo() bool {
    if len(ssa.SlaveAddress) == 0 || len(ssa.SlaveNodeMacAddr) == 0 || len(ssa.SlaveHardware) == 0 {
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
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveStatus{
        Version         : SLAVE_STATUS_VERSION,
        SlaveResponse   : SLAVE_WHO_I_AM,
        SlaveAddress    : piface.IP.String(),
        SlaveNodeMacAddr: piface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }, nil
}

func KeyExchangeStatus(master string, timestamp time.Time) (*PocketSlaveStatus, error) {
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveStatus{
        Version         : SLAVE_STATUS_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_SEND_PUBKEY,
        SlaveAddress    : piface.IP.String(),
        SlaveNodeMacAddr: piface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }, nil
}

func SlaveBindReadyStatus(master, nodename string, timestamp time.Time) (*PocketSlaveStatus, error) {
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveStatus{
        Version         : SLAVE_STATUS_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_BIND_READY,
        SlaveNodeName   : nodename,
        SlaveAddress    : piface.IP.String(),
        SlaveNodeMacAddr: piface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }, nil
}

func SlaveBoundedStatus(master string, timestamp time.Time) (*PocketSlaveStatus, error) {
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    hostname, err := os.Hostname()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveStatus{
        Version         : SLAVE_STATUS_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_REPORT_STATUS,
        SlaveNodeName   : hostname,
        SlaveAddress    : piface.IP.String(),
        SlaveNodeMacAddr: piface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }, nil
}