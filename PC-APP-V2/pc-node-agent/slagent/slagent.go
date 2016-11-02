package slagent

// ------ VERSION ------
// meta protocol scheme & version
type MetaProtocol string
const (
    SLAVE_META_PROTO        MetaProtocol = "pc_sl_pm"
    SLAVE_META_VERSION      MetaProtocol = "1.0.1"
)

// discovery protocol scheme & version
type DiscoveryProtocol string
const (
    SLAVE_DISCOVER_PROTO    DiscoveryProtocol = "pc_sl_pd"
    SLAVE_DISCOVER_VERSION  DiscoveryProtocol = "1.0.1"
)

// status protocol scheme & version
type StatusProtocol string
const (
    SLAVE_STATUS_PROTO      StatusProtocol = "pc_sl_ps"
    SLAVE_STATUS_VERSION    StatusProtocol = "1.0.1"
)

// Meta Agent field
const (
    SLAVE_MAC_ADDR          = "pc_sl_id"
    SLAVE_DISCOVER_AGENT    = "pc_sl_ad"
    SLAVE_STATUS_AGENT      = "pc_sl_as"
    SLAVE_ENCRYPTED_STATUS  = "pc_sl_es"
    SLAVE_PUBLIC_KEY        = "pc_sl_pk"
)

// Report types (respond to command)
type ResponseType string
const SLAVE_REPORT_TYPE string = "pc_sl_rt"
const (
    // slave node looks for a master agent
    SLAVE_LOOKUP_AGENT      ResponseType = "pc_sl_la"
    // slave node "check me" if i can join
    SLAVE_WHO_I_AM          ResponseType = "pc_sl_wi"
    // checking slave node public key
    SLAVE_SEND_PUBKEY       ResponseType = "pc_sl_sp"
    // slave node ready for binding
    SLAVE_CHECK_CRYPTO      ResponseType = "pc_sl_br"
    // slave node report status to the master agent
    SLAVE_REPORT_STATUS     ResponseType = "pc_sl_rs"
)

// ------ SLAVE SECTION ------
const (
    SLAVE_SECTION           = "slave"

    // slave info
    SLAVE_HARDWARE          = "pc_sl_hw"

    SLAVE_NODE_MACADDR      = "pc_sl_ma"
    SLAVE_NODE_NAME         = "pc_sl_nm"
    SLAVE_TIMESTAMP         = "pc_sl_ts"
    SLAVE_NAMESERVER        = "pc_sl_ns"

    SLAVE_IP4_GATEWAY       = "pc_sl_g4"
    SLAVE_IP4_ADDRESS       = "pc_sl_i4"
    SLAVE_IP4_NETMASK       = "pc_sl_n4"

    SLAVE_IP6_GATEWAY       = "pc_sl_g6"
    SLAVE_IP6_ADDRESS       = "pc_sl_i6"
    SLAVE_IP6_NETMASK       = "pc_sl_n6"

    //TODO check if this is really necessary. If we're to manage SSH sessions with a centralized server, this is not needed
    //SLAVE_CLUSTER_MEMBERS = "pc_sl_cl"
)

type SlaveAgent interface {
    // this is main loop of locating service to master
    MonitorLocatingService() error
}