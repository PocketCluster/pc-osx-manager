package slagent

import (
    "runtime"
    "os"
    "time"

    "github.com/stkim1/pc-node-agent/status"
)

type PocketSlaveStatusAgent struct {
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

func (ssa *PocketSlaveStatusAgent) IsAppropriateSlaveInfo() bool {
    if len(ssa.SlaveAddress) == 0 || len(ssa.SlaveNodeMacAddr) == 0 || len(ssa.SlaveHardware) == 0 {
        return false
    }
    return true
}

// Unbounded
func InquiredAgent(timestamp time.Time) (agent *PocketSlaveStatusAgent, err error) {
    _, gwifname, err := status.GetDefaultIP4Gateway(); if err != nil {
        return nil, err
    }
    // TODO : should this be fixed to have "eth0"?
    iface, err := status.InterfaceByName(gwifname); if err != nil {
        return nil, err
    }
    ipaddrs, err := iface.IP4Addrs(); if err != nil {
        return nil, err
    }
    agent = &PocketSlaveStatusAgent{
        Version         : SLAVE_STATUS_VERSION,
        SlaveResponse   : SLAVE_WHO_I_AM,
        SlaveAddress    : ipaddrs[0].IP.String(),
        SlaveNodeMacAddr: iface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }
    err = nil
    return
}

func KeyExchangeAgent(master string, timestamp time.Time) (agent *PocketSlaveStatusAgent, err error) {
    _, gwifname, err := status.GetDefaultIP4Gateway(); if err != nil {
        return nil, err
    }
    // TODO : should this be fixed to have "eth0"?
    iface, err := status.InterfaceByName(gwifname); if err != nil {
        return nil, err
    }
    ipaddrs, err := iface.IP4Addrs(); if err != nil {
        return nil, err
    }
    agent = &PocketSlaveStatusAgent{
        Version         : SLAVE_STATUS_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_SEND_PUBKEY,
        SlaveAddress    : ipaddrs[0].IP.String(),
        SlaveNodeMacAddr: iface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }
    err = nil
    return
}

func SlaveBindReadyAgent(master, nodename string, timestamp time.Time) (agent *PocketSlaveStatusAgent, err error) {
    _, gwifname, err := status.GetDefaultIP4Gateway(); if err != nil {
        return nil, err
    }
    // TODO : should this be fixed to have "eth0"?
    iface, err := status.InterfaceByName(gwifname); if err != nil {
        return nil, err
    }
    ipaddrs, err := iface.IP4Addrs(); if err != nil {
        return nil, err
    }
    agent = &PocketSlaveStatusAgent{
        Version         : SLAVE_STATUS_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_BIND_READY,
        SlaveNodeName   : nodename,
        SlaveAddress    : ipaddrs[0].IP.String(),
        SlaveNodeMacAddr: iface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }
    err = nil
    return
}

func BoundedStatusAgent(master string, timestamp time.Time) (agent *PocketSlaveStatusAgent, err error) {
    _, gwifname, err := status.GetDefaultIP4Gateway()
    if err != nil {
        return nil, err
    }
    // TODO : should this be fixed to have "eth0"?
    iface, err := status.InterfaceByName(gwifname)
    if err != nil {
        return nil, err
    }
    ipaddrs, err := iface.IP4Addrs()
    if err != nil {
        return nil, err
    }
    hostname, err := os.Hostname()
    if err != nil {
        return nil, err
    }
    agent = &PocketSlaveStatusAgent{
        Version         : SLAVE_STATUS_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_REPORT_STATUS,
        SlaveNodeName   : hostname,
        SlaveAddress    : ipaddrs[0].IP.String(),
        SlaveNodeMacAddr: iface.HardwareAddr.String(),
        SlaveHardware   : runtime.GOARCH,
        SlaveTimestamp  : timestamp,
    }
    err = nil
    return
}