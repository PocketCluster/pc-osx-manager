package msagent

import "github.com/stkim1/pc-node-agent/slagent"

// ------ VERSION ------
// meta protocol scheme & version
type MetaProtocol string
const (
    MASTER_META_PROTO      MetaProtocol = "pc_ms_pm"
    MASTER_META_VERSION    MetaProtocol = "1.0.1"
)

// discovery protocol scheme & version
type DiscoveryProtocol string
const (
    MASTER_DISCOVERY_PROTO    DiscoveryProtocol = "pc_ms_pd"
    MASTER_DISCOVERY_VERSION  DiscoveryProtocol = "1.0.1"
)

// status protocol scheme & version
type StatusProtocol string
const (
    MASTER_STATUS_PROTO    StatusProtocol = "pc_ms_ps"
    MASTER_STATUS_VERSION  StatusProtocol = "1.0.1"
)

// Discovery & Status tag for meta
const (
    MASTER_DISCOVER_RESP   = "pc_ms_dr"
    MASTER_STATUS_COLLECT  = "pc_ms_sc"
)

// ------ COMMAND DEFINITIONS ------
type CommandType string
const MASTER_COMMAND_TYPE  string = "pc_ms_ct"
const (
    COMMAND_WHO_R_U             CommandType = "pc_ms_wr"
    COMMAND_ASK_PUBKEY          CommandType = "pc_ms_ap"
    COMMAND_SEND_AES            CommandType = "pc_ms_sa"
    COMMAND_MASTER_BIND_READY   CommandType = "pc_ms_mr"
    COMMAND_SLAVE_ACK           CommandType = "pc_ms_ak"
)

// ------ MASTER SECTION ------
const (
    MASTER_SECTION         = "master"

    // bound-id
    MASTER_BOUND_AGENT     = "pc_ms_ba"
    // master ip4 / ip6
    MASTER_IP4_ADDRESS     = "pc_ms_i4"
    MASTER_IP6_ADDRESS     = "pc_ms_i6"
    // master datetime
    MASTER_TIMESTAMP       = "pc_ms_ts"

    // TODO : Do we need this?
    //MASTER_HOSTNAME        = "pc_ms_hn"
)

func MakeWhoruInquery(unboundedDiscovery *slagent.PocketSlaveDiscoveryAgent, masterAgent, ip4Address string) (responder *PocketMasterDiscoveryResponder, err error) {
    // TODO : sanity checker here
    return &PocketMasterDiscoveryResponder{
        Version          :MASTER_DISCOVERY_VERSION,
        MasterBoundAgent :masterAgent,
        MasterCommandType:COMMAND_WHO_R_U,
        MasterAddress    :ip4Address,
    }, nil
}

func MakeMasterPubkeyDelivery(unboundedStatus *slagent.PocketSlaveStatusAgent) (collector *PocketMasterDiscoveryResponder, err error) {
    return nil, nil
}

func ExchangeMasterSlaveKeys(unboundedStatus *slagent.PocketSlaveStatusAgent) (collector *PocketMasterStatusCollector, err error) {
    return nil, nil
}

func CheckMasterSlaveCrypto(unboundedStatus *slagent.PocketSlaveStatusAgent) (collector *PocketMasterStatusCollector, err error) {
    return nil, nil
}

func SendMasterBindReady(unboundedStatus *slagent.PocketSlaveStatusAgent) (collector *PocketMasterStatusCollector, err error) {
    return nil, nil
}
