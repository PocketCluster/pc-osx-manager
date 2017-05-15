package slagent

import (
    "runtime"
    "time"

    "github.com/pkg/errors"
    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/stkim1/pc-node-agent/slcontext"
)

type PocketSlaveStatus struct {
    Version             StatusProtocol  `msgpack:"s_ps"`
    // master
    MasterBoundAgent    string          `msgpack:"m_ba,omitempty"`
    // slave response
    SlaveResponse       ResponseType    `msgpack:"s_rt,omitempty`
    // slave
    SlaveNodeName       string          `msgpack:"s_nm,omitempty"`
    // current interface status
    SlaveAddress        string          `msgpack:"s_i4"`
    SlaveNodeMacAddr    string          `msgpack:"s_ma"`
    SlaveHardware       string          `msgpack:"s_hw"`
    SlaveTimestamp      time.Time       `msgpack:"s_ts"`
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
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveStatus {
        Version:             SLAVE_STATUS_VERSION,
        SlaveResponse:       SLAVE_WHO_I_AM,
        SlaveAddress:        piface.PrimaryIP4Addr(),
        SlaveNodeMacAddr:    piface.HardwareAddr,
        SlaveHardware:       runtime.GOARCH,
        SlaveTimestamp:      timestamp,
    }, nil
}

func KeyExchangeStatus(master string, timestamp time.Time) (*PocketSlaveStatus, error) {
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveStatus {
        Version:             SLAVE_STATUS_VERSION,
        MasterBoundAgent:    master,
        SlaveResponse:       SLAVE_SEND_PUBKEY,
        SlaveAddress:        piface.PrimaryIP4Addr(),
        SlaveNodeMacAddr:    piface.HardwareAddr,
        SlaveHardware:       runtime.GOARCH,
        SlaveTimestamp:      timestamp,
    }, nil
}

func CheckSlaveCryptoStatus(master, nodename string, timestamp time.Time) (*PocketSlaveStatus, error) {
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveStatus {
        Version:             SLAVE_STATUS_VERSION,
        MasterBoundAgent:    master,
        SlaveResponse:       SLAVE_CHECK_CRYPTO,
        SlaveNodeName:       nodename,
        SlaveAddress:        piface.PrimaryIP4Addr(),
        SlaveNodeMacAddr:    piface.HardwareAddr,
        SlaveHardware:       runtime.GOARCH,
        SlaveTimestamp:      timestamp,
    }, nil
}

func SlaveBoundedStatus(master, nodename string, timestamp time.Time) (*PocketSlaveStatus, error) {
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveStatus {
        Version:             SLAVE_STATUS_VERSION,
        MasterBoundAgent:    master,
        SlaveResponse:       SLAVE_REPORT_STATUS,
        SlaveNodeName:       nodename,
        SlaveAddress:        piface.PrimaryIP4Addr(),
        SlaveNodeMacAddr:    piface.HardwareAddr,
        SlaveHardware:       runtime.GOARCH,
        SlaveTimestamp:      timestamp,
    }, nil
}

// this is for master bindbroken state. Since majority of sanity check is done by beacon.bindbroken module, we'll just check simple things.
func ConvertBindAttemptDiscoveryAgent(discovery *PocketSlaveDiscovery, slaveNode, slaveHardware string) (*PocketSlaveStatus, error) {
    if len(slaveNode) == 0 {
        return nil, errors.Errorf("[ERR] incorrect slave name")
    }
    if len(slaveHardware) == 0 {
        return nil, errors.Errorf("[ERR] incrrect slave hardware architecture")
    }
    if discovery.Version != SLAVE_DISCOVER_VERSION {
        return nil, errors.Errorf("[ERR] Incorrect SlaveDiscoveryAgent version")
    }
    if len(discovery.MasterBoundAgent) == 0 {
        return nil, errors.Errorf("[ERR] Incorrect master agent name")
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
    if len(discovery.SlaveNodeMacAddr) == 0 {
        return nil, errors.Errorf("[ERR] incorrect slave macaddress")
    }
    return &PocketSlaveStatus{
        Version:             SLAVE_STATUS_VERSION,
        MasterBoundAgent:    discovery.MasterBoundAgent,
        SlaveResponse:       SLAVE_REPORT_STATUS,
        SlaveNodeName:       slaveNode,
        SlaveAddress:        discovery.SlaveAddress,
        SlaveNodeMacAddr:    discovery.SlaveNodeMacAddr,
        SlaveHardware:       slaveHardware,
        // TODO : since discovery agent does not have timestamp, we'll use master timstamp.
        SlaveTimestamp:      time.Now(),
    }, nil
}