package slagent

import "github.com/stkim1/pc-node-agent/status"

type PocketSlaveDiscoveryAgent struct {
    Version             DiscoveryProtocol    `msgpack:"pc_sl_pd"`
    // master
    MasterBoundAgent    string               `msgpack:"pc_ms_ba,omitempty"`
    // slave response
    SlaveResponse       ResponseType         `msgpack:"pc_sl_rt,omitempty`

    // slave
    SlaveAddress        string               `msgpack:"pc_sl_i4,omitempty"`
    SlaveGateway        string               `msgpack:"pc_sl_g4,omitempty"`
    SlaveNetmask        string               `msgpack:"pc_sl_n4,omitempty"`
    SlaveNodeMacAddr    string               `msgpack:"pc_sl_ma"`

    // TODO : check if nameserver & node name is really necessary for discovery
    //SlaveNameServer     string     `bson:"pc_sl_ns,omitempty"
    //SlaveNodeName       string     `bson:"pc_sl_nm,omitempty"
}

func (sda *PocketSlaveDiscoveryAgent) IsAppropriateSlaveInfo() bool {
    if len(sda.SlaveAddress) == 0 || len(sda.SlaveGateway) == 0 || len(sda.SlaveNetmask) == 0 || len(sda.SlaveNodeMacAddr) == 0 {
        return false
    }
    return true
}

func UnboundedBroadcastAgent() (agent *PocketSlaveDiscoveryAgent, err error) {
    gwaddr, gwifname, err := status.GetDefaultIP4Gateway()
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
    agent = &PocketSlaveDiscoveryAgent{
        Version         : SLAVE_DISCOVER_VERSION,
        SlaveResponse   : SLAVE_LOOKUP_AGENT,
        SlaveAddress    : ipaddrs[0].IP.String(),
        SlaveGateway    : gwaddr,
        SlaveNetmask    : ipaddrs[0].IPMask.String(),
        SlaveNodeMacAddr: iface.HardwareAddr.String(),
    }
    err = nil
    return
}

func BindBrokenBroadcastAgent(master string) (agent *PocketSlaveDiscoveryAgent, err error) {
    gwaddr, gwifname, err := status.GetDefaultIP4Gateway()
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
    agent = &PocketSlaveDiscoveryAgent{
        Version         : SLAVE_DISCOVER_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_LOOKUP_AGENT,
        SlaveAddress    : ipaddrs[0].IP.String(),
        SlaveGateway    : gwaddr,
        SlaveNetmask    : ipaddrs[0].IPMask.String(),
        SlaveNodeMacAddr: iface.HardwareAddr.String(),
    }
    err = nil
    return
}