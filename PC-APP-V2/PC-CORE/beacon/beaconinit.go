package beacon

import (
    "fmt"
    "time"

    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/model"
)

func beaconinitState(comm CommChannel) BeaconState {
    b := &beaconinit{}

    b.constState                    = MasterInit

    b.constTransitionFailureLimit   = TransitionFailureLimit
    b.constTransitionTimeout        = UnboundedTimeout * time.Duration(TxActionLimit)
    b.constTxActionLimit            = TxActionLimit
    b.constTxTimeWindow             = UnboundedTimeout

    b.lastTransitionTS              = time.Now()

    b.timestampTransition           = b.transitionActionWithTimestamp
    b.slaveMetaTransition           = b.beaconInit
    b.onTransitionSuccess           = b.onStateTranstionSuccess
    b.onTransitionFailure           = b.onStateTranstionFailure

    b.slaveNode                     = &model.SlaveNode{}
    b.commChan                      = comm

    b.slaveLocation                 = nil
    b.slaveStatus                   = nil

    return b
}

type beaconinit struct {
    beaconState
}

func (b *beaconinit) transitionActionWithTimestamp(masterTimestamp time.Time) error {
    return nil
}

func (b *beaconinit) beaconInit(meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if meta.DiscoveryAgent == nil || meta.DiscoveryAgent.Version != slagent.SLAVE_DISCOVER_VERSION {
        return MasterTransitionFail, fmt.Errorf("[ERR] Null or incorrect version of slave discovery")
    }
    // if slave isn't looking for agent, then just return. this is not for this state.
    if meta.DiscoveryAgent.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return MasterTransitionIdle, nil
    }
    if len(meta.DiscoveryAgent.MasterBoundAgent) != 0 {
        return MasterTransitionIdle, fmt.Errorf("[ERR] Incorrect slave bind. Slave should not be bound to a master when it looks for joining")
    }
    // slave ip address
    if len(meta.DiscoveryAgent.SlaveAddress) == 0 {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave node address")
    }
    b.slaveNode.IP4Address = meta.DiscoveryAgent.SlaveAddress

    // slave ip gateway
    if len(meta.DiscoveryAgent.SlaveGateway) == 0 {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave node gateway")
    }
    b.slaveNode.IP4Gateway = meta.DiscoveryAgent.SlaveGateway

    // slave ip netmask
    if len(meta.DiscoveryAgent.SlaveNetmask) == 0 {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave node netmask")
    }
    b.slaveNode.IP4Netmask = meta.DiscoveryAgent.SlaveNetmask

    // slave mac address
    if meta.SlaveID != meta.DiscoveryAgent.SlaveNodeMacAddr {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave ID")
    }
    if len(meta.DiscoveryAgent.SlaveNodeMacAddr) == 0 {
        return MasterTransitionFail, fmt.Errorf("[ERR] Inappropriate slave MAC address")
    }
    b.slaveNode.MacAddress = meta.DiscoveryAgent.SlaveNodeMacAddr

    // save slave discovery to send responsed
    b.slaveLocation  = meta.DiscoveryAgent

    return MasterTransitionOk, nil
}

func (b *beaconinit) onStateTranstionSuccess(masterTimestamp time.Time) error {
    return nil
}

func (b *beaconinit) onStateTranstionFailure(masterTimestamp time.Time) error {
    return nil
}
