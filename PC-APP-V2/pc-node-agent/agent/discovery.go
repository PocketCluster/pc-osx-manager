package agent

import "github.com/stkim1/pc-node-agent/status"

type PocketSlaveDiscoveryAgent struct {
    Version             string      `bson:"pc_sl_dc"    json:"pc_sl_dc"`
    // master
    MasterBoundAgent    string      `bson:"pc_ma_ba"    json:"pc_ma_ba"`

    // slave
    SlaveAddress        string      `bson:"pc_sl_i4,omitempty"    json:"pc_sl_i4"`
    SlaveGateway        string      `bson:"pc_sl_ma,omitempty"    json:"pc_sl_ma"`
    SlaveNetmask        string      `bson:"pc_sl_ma,omitempty"    json:"pc_sl_ma"`
    SlaveNodeMacAddr    string      `bson:"pc_sl_ma,omitempty"    json:"pc_sl_ma"`

    // TODO : check if nameserver & node name is really necessary for discovery
    //SlaveNameServer     string      `bson:"pc_sl_ns,omitempty"    json:"pc_sl_ns"`
    //SlaveNodeName       string      `bson:"pc_sl_nm,omitempty"    json:"pc_sl_nm"`
}

func UnboundedBroadcastAgent() (agent *PocketSlaveDiscoveryAgent, err error) {
    gwaddr, gwifname, err := status.GetDefaultIP4Gateway(); if err != nil {
        return nil, err
    }
    // TODO : should this be fixed to have "eth0"?
    iface, err := status.InterfaceByName(gwifname); if err != nil {
        return nil, err
    }
    ipaddrs, err := iface.IP4Addrs(); if err != nil {
        return nil, err
    }
    agent = &PocketSlaveDiscoveryAgent{
        Version         : SLAVE_DISCOVER_VERSION,
        MasterBoundAgent: SLAVE_LOOKUP_AGENT,
        SlaveAddress    : ipaddrs[0].IP.String(),
        SlaveGateway    : gwaddr,
        SlaveNetmask    : ipaddrs[0].IPMask.String(),
        SlaveNodeMacAddr: iface.HardwareAddr.String(),
    }
    err = nil
    return
}

func BoundedBroadcastAgent(master string) (agent *PocketSlaveDiscoveryAgent, err error) {
    gwaddr, gwifname, err := status.GetDefaultIP4Gateway(); if err != nil {
        return nil, err
    }
    // TODO : should this be fixed to have "eth0"?
    iface, err := status.InterfaceByName(gwifname); if err != nil {
        return nil, err
    }
    ipaddrs, err := iface.IP4Addrs(); if err != nil {
        return nil, err
    }
    agent = &PocketSlaveDiscoveryAgent{
        Version         : SLAVE_DISCOVER_VERSION,
        MasterBoundAgent: master,
        SlaveAddress    : ipaddrs[0].IP.String(),
        SlaveGateway    : gwaddr,
        SlaveNetmask    : ipaddrs[0].IPMask.String(),
        SlaveNodeMacAddr: iface.HardwareAddr.String(),
    }
    err = nil
    return
}