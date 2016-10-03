package slagent

import "github.com/stkim1/pc-node-agent/status"

type PocketSlaveStatusAgent struct {
    Version             string      `msgpack:"pc_sl_ps"`

    // master
    MasterBoundAgent    string      `msgpack:"pc_ma_ba"`

    // slave
    SlaveNodeName       string      `msgpack:"pc_sl_nm"`

    // current interface status
    SlaveAddress        string      `msgpack:"pc_sl_i4"`
    SlaveNodeMacAddr    string      `msgpack:"pc_sl_ma"`
    SlaveTimeZone       string      `msgpack:"pc_sl_tz"`
}

func BoundedStatusAgent(master, slave, timezone string) (agent *PocketSlaveStatusAgent, err error) {
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
        SlaveNodeName   : slave,
        SlaveAddress    : ipaddrs[0].IP.String(),
        SlaveNodeMacAddr: iface.HardwareAddr.String(),
        SlaveTimeZone   : timezone,
    }
    err = nil
    return
}