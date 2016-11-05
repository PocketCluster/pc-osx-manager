package slagent

import (
    "github.com/stkim1/pc-node-agent/slcontext"
)

type PocketSlaveDiscovery struct {
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

func (sda *PocketSlaveDiscovery) IsAppropriateSlaveInfo() bool {
    if len(sda.SlaveAddress) == 0 || len(sda.SlaveGateway) == 0 || len(sda.SlaveNetmask) == 0 || len(sda.SlaveNodeMacAddr) == 0 {
        return false
    }
    return true
}

func UnboundedMasterDiscovery() (*PocketSlaveDiscovery, error) {
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveDiscovery {
        Version         : SLAVE_DISCOVER_VERSION,
        SlaveResponse   : SLAVE_LOOKUP_AGENT,
        SlaveAddress    : piface.IP.String(),
        SlaveGateway    : piface.GatewayAddr,
        SlaveNetmask    : piface.IPMask.String(),
        SlaveNodeMacAddr: piface.HardwareAddr.String(),
    }, nil
}

func BrokenBindDiscovery(master string) (*PocketSlaveDiscovery, error) {
    piface, err := slcontext.SharedSlaveContext().PrimaryNetworkInterface()
    if err != nil {
        return nil, err
    }
    return &PocketSlaveDiscovery {
        Version         : SLAVE_DISCOVER_VERSION,
        MasterBoundAgent: master,
        SlaveResponse   : SLAVE_LOOKUP_AGENT,
        SlaveAddress    : piface.IP.String(),
        SlaveGateway    : piface.GatewayAddr,
        SlaveNetmask    : piface.IPMask.String(),
        SlaveNodeMacAddr: piface.HardwareAddr.String(),
    }, nil
}