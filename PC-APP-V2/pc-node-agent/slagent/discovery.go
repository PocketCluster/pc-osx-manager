package slagent

import (
    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slcontext"
)

type PocketSlaveDiscovery struct {
    Version             DiscoveryProtocol    `msgpack:"s_pd"`
    // master
    MasterBoundAgent    string               `msgpack:"m_ba,omitempty"`
    // slave response
    SlaveResponse       ResponseType         `msgpack:"s_rt,omitempty`

    // slave
    SlaveAddress        string               `msgpack:"s_i4,omitempty"`
    SlaveGateway        string               `msgpack:"s_g4,omitempty"`
    SlaveNodeMacAddr    string               `msgpack:"s_ma"`

    // TODO : check if nameserver & node name is really necessary for discovery
    //SlaveNameServer     string     `bson:"pc_sl_ns,omitempty"
    //SlaveNodeName       string     `bson:"pc_sl_nm,omitempty"
}

func (sda *PocketSlaveDiscovery) IsAppropriateSlaveInfo() bool {
    if len(sda.SlaveAddress) == 0 || len(sda.SlaveGateway) == 0 || len(sda.SlaveNodeMacAddr) == 0 {
        return false
    }
    return true
}

func UnboundedMasterDiscovery() (*PocketSlaveDiscovery, error) {
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketSlaveDiscovery {
        Version:             SLAVE_DISCOVER_VERSION,
        SlaveResponse:       SLAVE_LOOKUP_AGENT,
        SlaveAddress:        piface.PrimaryIP4Addr(),
        SlaveGateway:        piface.GatewayAddr,
        SlaveNodeMacAddr:    piface.HardwareAddr,
    }, nil
}

func BrokenBindDiscovery(master string) (*PocketSlaveDiscovery, error) {
    piface, err := slcontext.PrimaryNetworkInterface()
    if err != nil {
        return nil, errors.WithStack(err)
    }
    return &PocketSlaveDiscovery {
        Version:             SLAVE_DISCOVER_VERSION,
        MasterBoundAgent:    master,
        SlaveResponse:       SLAVE_LOOKUP_AGENT,
        SlaveAddress:        piface.PrimaryIP4Addr(),
        SlaveGateway:        piface.GatewayAddr,
        SlaveNodeMacAddr:    piface.HardwareAddr,
    }, nil
}