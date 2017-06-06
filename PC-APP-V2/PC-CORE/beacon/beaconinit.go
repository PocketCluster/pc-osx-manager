package beacon

import (
    "net"
    "time"

    "github.com/pkg/errors"
    "github.com/stkim1/pc-node-agent/slagent"
    "github.com/stkim1/pc-core/model"
)

func beaconinitState(slaveNode *model.SlaveNode, comm CommChannel, event BeaconOnTransitionEvent) BeaconState {
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

    b.BeaconOnTransitionEvent       = event
    b.slaveNode                     = slaveNode
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

func (b *beaconinit) beaconInit(sender *net.UDPAddr, meta *slagent.PocketSlaveAgentMeta, timestamp time.Time) (MasterBeaconTransition, error) {
    if sender == nil {
        return MasterTransitionIdle, errors.Errorf("[ERR] incorrect slave input. slave address should not be nil when receiving multicast while initializing bind.")
    }
    if meta.DiscoveryAgent == nil || meta.DiscoveryAgent.Version != slagent.SLAVE_DISCOVER_VERSION {
        return MasterTransitionFail, errors.Errorf("[ERR] Null or incorrect version of slave discovery")
    }
    // if slave isn't looking for agent, then just return. this is not for this state.
    if meta.DiscoveryAgent.SlaveResponse != slagent.SLAVE_LOOKUP_AGENT {
        return MasterTransitionIdle, nil
    }
    if len(meta.MasterBoundAgent) != 0 {
        return MasterTransitionIdle, errors.Errorf("[ERR] Incorrect slave bind. Slave should not be bound to a master when it looks for joining")
    }

    // TODO : (2015-05-16) we're not checking ip + subnet eligivility for now
    // slave ip address
    addr, err := model.IP4AddrToString(meta.DiscoveryAgent.SlaveAddress)
    if err != nil {
        return MasterTransitionFail, errors.WithStack(err)
    }
    if addr != sender.IP.String() {
        return MasterTransitionFail, errors.Errorf("[ERR] malicious slave ip address.")
    }
    b.slaveNode.IP4Address = meta.DiscoveryAgent.SlaveAddress
    // slave ip gateway
    if len(meta.DiscoveryAgent.SlaveGateway) == 0 {
        return MasterTransitionFail, errors.Errorf("[ERR] Inappropriate slave node gateway")
    }
    b.slaveNode.IP4Gateway = meta.DiscoveryAgent.SlaveGateway
    // slave mac address
    if len(meta.SlaveID) == 0 {
        return MasterTransitionFail, errors.Errorf("[ERR] Inappropriate slave MAC address")
    }
    b.slaveNode.SlaveID = meta.SlaveID

    // save slave discovery to send responsed
    b.slaveLocation = meta.DiscoveryAgent

    return MasterTransitionOk, nil
}

func (b *beaconinit) onStateTranstionSuccess(masterTimestamp time.Time) error {
    return nil
}

func (b *beaconinit) onStateTranstionFailure(masterTimestamp time.Time) error {
    return nil
}
