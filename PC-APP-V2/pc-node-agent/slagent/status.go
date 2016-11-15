package slagent

import (
    "runtime"
    "time"

    "gopkg.in/vmihailenco/msgpack.v2"
    "github.com/stkim1/pc-node-agent/slcontext"
    "fmt"
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
    return &PocketSlaveStatus {
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
    return &PocketSlaveStatus {
        Version         : SLAVE_STATUS_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_SEND_PUBKEY,
        SlaveAddress    : piface.IP.String(),
        SlaveNodeMacAddr: piface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }, nil
}

func CheckSlaveCryptoStatus(master, nodename string, timestamp time.Time) (*PocketSlaveStatus, error) {
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveStatus {
        Version         : SLAVE_STATUS_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_CHECK_CRYPTO,
        SlaveNodeName   : nodename,
        SlaveAddress    : piface.IP.String(),
        SlaveNodeMacAddr: piface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }, nil
}

func SlaveBoundedStatus(master, nodename string, timestamp time.Time) (*PocketSlaveStatus, error) {
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveStatus {
        Version         : SLAVE_STATUS_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_REPORT_STATUS,
        SlaveNodeName   : nodename,
        SlaveAddress    : piface.IP.String(),
        SlaveNodeMacAddr: piface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }, nil
}

// this is for master bindbroken state. Since majority of sanity check is done by beacon.bindbroken module, we'll just check simple things.
func ConvertBindAttemptDiscoveryAgent(discovery *PocketSlaveDiscovery, slaveNode, slaveHardware string) (*PocketSlaveStatus, error) {
    if len(slaveNode) == 0 {
        return nil, fmt.Errorf("[ERR] incorrect slave name")
    }
    if len(slaveHardware) == 0 {
        return nil, fmt.Errorf("[ERR] incrrect slave hardware architecture")
    }
    if discovery.Version != SLAVE_DISCOVER_VERSION {
        return nil, fmt.Errorf("[ERR] Incorrect SlaveDiscoveryAgent version")
    }
    if len(discovery.MasterBoundAgent) == 0 {
        return nil, fmt.Errorf("[ERR] Incorrect master agent name")
    }
    if discovery.SlaveResponse != SLAVE_LOOKUP_AGENT {
        return nil, fmt.Errorf("[ERR] incorrect slave discovery response")
    }
    if len(discovery.SlaveAddress) == 0 {
        return nil, fmt.Errorf("[ERR] incorrect slave address")
    }
    if len(discovery.SlaveGateway) == 0 {
        return nil, fmt.Errorf("[ERR] incorrect slave gateway")
    }
    if len(discovery.SlaveNetmask) == 0 {
        return nil, fmt.Errorf("[ERR] incorrect slave netmask")
    }
    if len(discovery.SlaveNodeMacAddr) == 0 {
        return nil, fmt.Errorf("[ERR] incorrect slave macaddress")
    }
    return &PocketSlaveStatus{
        Version             : SLAVE_STATUS_VERSION,
        MasterBoundAgent    : discovery.MasterBoundAgent,
        SlaveResponse       : SLAVE_REPORT_STATUS,
        SlaveNodeName       : slaveNode,
        SlaveAddress        : discovery.SlaveAddress,
        SlaveNodeMacAddr    : discovery.SlaveNodeMacAddr,
        SlaveHardware       : slaveHardware,
        // TODO : since discovery agent does not have timestamp, we'll use master timstamp.
        SlaveTimestamp      : time.Now(),
    }, nil
}