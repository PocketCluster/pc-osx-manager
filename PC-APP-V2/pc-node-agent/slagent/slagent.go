package slagent

// ------ VERSION ------
// meta protocol scheme & version
type MetaProtocol string
const (
    SLAVE_META_PROTO        MetaProtocol = "s_pm"
    SLAVE_META_VERSION      MetaProtocol = "1.0.1"
)

// discovery protocol scheme & version
type DiscoveryProtocol string
const (
    SLAVE_DISCOVER_PROTO    DiscoveryProtocol = "s_pd"
    SLAVE_DISCOVER_VERSION  DiscoveryProtocol = "1.0.1"
)

// status protocol scheme & version
type StatusProtocol string
const (
    SLAVE_STATUS_PROTO      StatusProtocol = "s_ps"
    SLAVE_STATUS_VERSION    StatusProtocol = "1.0.1"
)

// Meta Agent field
const (
    SLAVE_MAC_ADDR          = "s_id"
    SLAVE_DISCOVER_AGENT    = "s_ad"
    SLAVE_STATUS_AGENT      = "s_as"
    SLAVE_ENCRYPTED_STATUS  = "s_es"
    SLAVE_PUBLIC_KEY        = "s_pk"
    SLAVE_SSH_KEY           = "s_sh"
)

// Report types (respond to command)
type ResponseType string
const SLAVE_REPORT_TYPE string = "s_rt"
const (
    // slave node looks for a master agent
    SLAVE_LOOKUP_AGENT      ResponseType = "s_la"
    // slave node "check me" if i can join
    SLAVE_WHO_I_AM          ResponseType = "s_wi"
    // checking slave node public key
    SLAVE_SEND_PUBKEY       ResponseType = "s_sp"
    // slave node ready for binding
    SLAVE_CHECK_CRYPTO      ResponseType = "s_br"
    // slave node report status to the master agent
    SLAVE_REPORT_STATUS     ResponseType = "s_rs"
)

// ------ SLAVE SECTION ------
const (
    //SLAVE_SECTION           = "slave"

    // slave info
    SLAVE_HARDWARE          = "s_hw"

    SLAVE_NODE_MACADDR      = "s_ma"
    SLAVE_NODE_NAME         = "s_nm"
    SLAVE_TIMESTAMP         = "s_ts"
    SLAVE_NAMESERVER        = "s_ns"

    SLAVE_IP4_GATEWAY       = "s_g4"
    SLAVE_IP4_ADDRESS       = "s_i4"
    SLAVE_IP4_NETMASK       = "s_n4"

    //SLAVE_IP6_GATEWAY       = "s_g6"
    //SLAVE_IP6_ADDRESS       = "s_i6"
    //SLAVE_IP6_NETMASK       = "s_n6"

    //TODO check if this is really necessary. If we're to manage SSH sessions with a centralized server, this is not needed
    //SLAVE_CLUSTER_MEMBERS = "s_cl"
)
