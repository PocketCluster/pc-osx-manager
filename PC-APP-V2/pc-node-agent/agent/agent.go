package agent

// ------ VERSION ------
const (
    SLAVE_META_PROTO    = "pc_sl_mt"
    SLAVE_META_VERSION  = "1.0.1"
)

const (
    SLAVE_DISCOVER_PROTO   = "pc_sl_dc"
    SLAVE_DISCOVER_VERSION = "1.0.1"
)

const (
    SLAVE_STATUS_PROTO  = "pc_sl_st"
    SLAVE_STATUS_VERSION= "1.0.1"
)

// ------ PROTOCOL DEFINITIONS ------
const (
    MASTER_COMMAND_TYPE = "pc_ma_ct"
    COMMAND_FIX_BOUND   = "ct_fix_bound"
)

// ------ MASTER SECTION ------
const (
    MASTER_SECTION      = "master"

    // bound-id
    MASTER_BOUND_AGENT  = "pc_ma_ba"
    // master ip4 / ip6
    MASTER_IP4_ADDRESS  = "pc_ma_i4"
    MASTER_IP6_ADDRESS  = "pc_ma_i6"
    // master hostname
    MASTER_HOSTNAME     = "pc_ma_hn"
    // master datetime
    MASTER_DATETIME     = "pc_ma_dt"
    MASTER_TIMEZONE     = "pc_ma_tz"
)

// ------ SLAVE SECTION ------
const (
    SLAVE_SECTION       = "slave"

    // node looks for agent
    SLAVE_LOOKUP_AGENT  = "pc_sl_la"
    SLAVE_NODE_MACADDR  = "pc_sl_ma"
    SLAVE_NODE_NAME     = "pc_sl_nm"
    SLAVE_TIMEZONE      = "pc_sl_tz"
    SLAVE_IP4_ADDRESS   = "pc_sl_i4"
    SLAVE_IP6_ADDRESS   = "pc_sl_i6"
    SLAVE_NAMESERVER    = "pc_sl_ns"

    //TODO check if this is really necessary. If we're to manage SSH sessions with a centralized server, this is not needed
    //SLAVE_CLUSTER_MEMBERS = "pc_sl_cl"
)