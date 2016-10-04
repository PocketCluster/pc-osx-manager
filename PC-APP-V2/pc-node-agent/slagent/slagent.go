package slagent

// ------ VERSION ------
// meta protocol scheme & version
const (
    SLAVE_META_PROTO    = "pc_sl_pm"
    SLAVE_META_VERSION  = "1.0.1"
)

// discovery protocol scheme & version
const (
    SLAVE_DISCOVER_PROTO   = "pc_sl_pd"
    SLAVE_DISCOVER_VERSION = "1.0.1"
)

// status protocol scheme & version
const (
    SLAVE_STATUS_PROTO  = "pc_sl_ps"
    SLAVE_STATUS_VERSION= "1.0.1"
)

// Discovery & Status tag for meta
const (
    SLAVE_DISCOVER_AGENT= "pc_sl_ad"
    SLAVE_STATUS_AGENT  = "pc_sl_as"
)

// ------ SLAVE SECTION ------
const (
    SLAVE_SECTION       = "slave"

    // node looks for agent
    SLAVE_LOOKUP_AGENT  = "pc_sl_la"

    // slave info
    SLAVE_HARDWARE      = "pc_sl_hw"
    SLAVE_NODE_MACADDR  = "pc_sl_ma"
    SLAVE_NODE_NAME     = "pc_sl_nm"
    SLAVE_TIMESTAMP     = "pc_sl_ts"
    SLAVE_IP4_ADDRESS   = "pc_sl_i4"
    SLAVE_IP6_ADDRESS   = "pc_sl_i6"
    SLAVE_NAMESERVER    = "pc_sl_ns"

    //TODO check if this is really necessary. If we're to manage SSH sessions with a centralized server, this is not needed
    //SLAVE_CLUSTER_MEMBERS = "pc_sl_cl"
)