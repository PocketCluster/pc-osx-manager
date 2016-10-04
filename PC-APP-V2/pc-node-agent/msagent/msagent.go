package msagent

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

