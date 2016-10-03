package agent

import "github.com/stkim1/pc-node-agent/status"

// ------ VERSION ------
const (
    PC_PROTO            = "pc_ver"
    VERSION             = "1.0.1"
)

// ------ PROTOCOL DEFINITIONS ------
const (
    MASTER_COMMAND_TYPE = "pc_ma_ct"
    COMMAND_FIX_BOUND   = "ct_fix_bound"
)

// ------ MASTER SECTION ------
const (
    MASTER_SECTION      = "master"

    // bound-id
    MASTER_BOUND_AGENT  = "pc_ma_ba"
    // master ip4 / ip6
    MASTER_IP4_ADDRESS  = "pc_ma_i4"
    MASTER_IP6_ADDRESS  = "pc_ma_i6"
    // master hostname
    MASTER_HOSTNAME     = "pc_ma_hn"
    // master datetime
    MASTER_DATETIME     = "pc_ma_dt"
    MASTER_TIMEZONE     = "pc_ma_tz"
)

// ------ SLAVE SECTION ------
const (
    SLAVE_SECTION       = "slave"

    // node looks for agent
    SLAVE_LOOKUP_AGENT  = "pc_sl_la"
    SLAVE_NODE_MACADDR  = "pc_sl_ma"
    SLAVE_NODE_NAME     = "pc_sl_nm"
    SLAVE_TIMEZONE      = "pc_sl_tz"
    SLAVE_IP4_ADDRESS   = "pc_sl_i4"
    SLAVE_IP6_ADDRESS   = "pc_sl_i6"
    SLAVE_NAMESERVER    = "pc_sl_ns"

    //TODO check if this is really necessary. If we're to manage SSH sessions with a centralized server, this is not needed
    //SLAVE_CLUSTER_MEMBERS = "pc_sl_cl"
)

type PocketSlaveAgent struct {
    Version             string      `bson:"pc_ver"      json:"pc_ver"`

    // master
    MasterBoundAgent    string      `bson:"pc_ma_ba"    json:"pc_ma_ba"`

    // slave
    SlaveNodeName       string      `bson:"pc_sl_nm,omitempty"    json:"pc_sl_nm"`
    // current interface status
    SlaveAddress        string      `bson:"pc_sl_i4,omitempty"    json:"pc_sl_i4"`
    SlaveNodeMacAddr    string      `bson:"pc_sl_ma,omitempty"    json:"pc_sl_ma"`
    SlaveTimeZone       string      `bson:"pc_sl_tz,omitempty"    json:"pc_sl_tz"`

    // TODO check if nameserver setting is really needed.
    //SlaveNameServer     string      `bson:"pc_sl_ns,omitempty"    json:"pc_sl_ns"`
}

func UnboundBroadcastAgent() (agent *PocketSlaveAgent, err error) {
    _, ifname, err := status.GetDefaultIP4Gateway(); if err != nil {
        return nil, err
    }
    // TODO : should this be fixed to have "eth0"?
    iface, err := status.InterfaceByName(ifname); if err != nil {
        return nil, err
    }
    ipaddrs, err := status.IP4Addrs(iface); if err != nil || len(ipaddrs) == 0 {
        return nil, err
    }
    agent = &PocketSlaveAgent{
        Version:VERSION,
        MasterBoundAgent:SLAVE_LOOKUP_AGENT,
        SlaveNodeMacAddr: iface.HardwareAddr.String(),
        SlaveAddress: ipaddrs[0].IP.String(),
    }
    err = nil
    return
}

func (pa *PocketSlaveAgent) getBoundBroadcast() {
}