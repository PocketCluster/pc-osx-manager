package msagent

// ------ VERSION ------
// meta protocol scheme & version
type MetaProtocol string
const (
    MASTER_META_PROTO           MetaProtocol = "pc_ms_pm"
    MASTER_META_VERSION         MetaProtocol = "1.0.1"
)

// discovery protocol scheme & version
type RespondProtocol string
const (
    MASTER_RESPOND_PROTO        RespondProtocol = "pc_ms_pr"
    MASTER_RESPOND_VERSION      RespondProtocol = "1.0.1"
)

// status protocol scheme & version
type CommandProtocol string
const (
    MASTER_COMMAND_PROTO        CommandProtocol = "pc_ms_pc"
    MASTER_COMMAND_VERSION      CommandProtocol = "1.0.1"
)

// Discovery & Status tag for meta
const (
    MASTER_ENCRYPTED_COMMAND    = "pc_ms_ec"
    MASTER_DISCOVER_RESPOND     = "pc_ms_dr"
    MASTER_STATUS_COMMANDER     = "pc_ms_sc"
    MASTER_PUBLIC_KEY           = "pc_ms_pk"
    MASTER_ENCRYPTED_AESKEY     = "pc_ms_ak"
    MASTER_RSA_SIGNATURE        = "pc_ms_sg"
    MASTER_ENCRYPTED_SLAVE      = "pc_ms_es"
    MASTER_ENCRYPTED_RESPOND    = "pc_ms_er"
)

// ------ COMMAND DEFINITIONS ------
type CommandType string
const MASTER_COMMAND_TYPE  string = "pc_ms_ct"
const (
    COMMAND_WHO_R_U             CommandType = "pc_ms_wr"
    COMMAND_SEND_PUBKEY         CommandType = "pc_ms_sp"
    COMMAND_SEND_AES            CommandType = "pc_ms_sa"
    COMMAND_MASTER_BIND_READY   CommandType = "pc_ms_mr"
    COMMAND_SLAVE_ACK           CommandType = "pc_ms_ak"
    COMMAND_RECOVER_BIND        CommandType = "pc_ms_rb"

)

// ------ MASTER SECTION ------
const (
    MASTER_SECTION              = "master"

    // bound-id
    MASTER_BOUND_AGENT          = "pc_ms_ba"
    // master ip4 / ip6
    MASTER_IP4_ADDRESS          = "pc_ms_i4"
    MASTER_IP6_ADDRESS          = "pc_ms_i6"
    // master datetime
    MASTER_TIMESTAMP            = "pc_ms_ts"

    // TODO : Do we need this?
    //MASTER_HOSTNAME        = "pc_ms_hn"
)
