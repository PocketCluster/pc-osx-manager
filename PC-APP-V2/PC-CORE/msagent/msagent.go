package msagent

// ------ VERSION ------
// meta protocol scheme & version
type MetaProtocol string
const (
    MASTER_META_PROTO           MetaProtocol = "m_pm"
    MASTER_META_VERSION         MetaProtocol = "1.0.1"
)

// discovery protocol scheme & version
type RespondProtocol string
const (
    MASTER_RESPOND_PROTO        RespondProtocol = "m_pr"
    MASTER_RESPOND_VERSION      RespondProtocol = "1.0.1"
)

// status protocol scheme & version
type CommandProtocol string
const (
    MASTER_COMMAND_PROTO        CommandProtocol = "m_pc"
    MASTER_COMMAND_VERSION      CommandProtocol = "1.0.1"
)

// Discovery & Status tag for meta
const (
    MASTER_ENCRYPTED_COMMAND    = "m_ec"
    MASTER_DISCOVER_RESPOND     = "m_dr"
    MASTER_STATUS_COMMANDER     = "m_sc"
    MASTER_PUBLIC_KEY           = "m_pk"
    MASTER_ENCRYPTED_AESKEY     = "m_ak"
    MASTER_RSA_SIGNATURE        = "m_sg"
    MASTER_ENCRYPTED_SLAVE      = "m_es"
    MASTER_ENCRYPTED_RESPOND    = "m_er"
)

// ------ COMMAND DEFINITIONS ------
type CommandType string
const MASTER_COMMAND_TYPE  string = "m_ct"
const (
    COMMAND_SLAVE_IDINQUERY     CommandType = "m_si"
    COMMAND_MASTER_DECLARE      CommandType = "m_md"
    COMMAND_EXCHANGE_CRPTKEY    CommandType = "m_ck"
    COMMAND_MASTER_BIND_READY   CommandType = "m_mr"
    COMMAND_SLAVE_ACK           CommandType = "m_ak"
    COMMAND_RECOVER_BIND        CommandType = "m_rb"
    COMMAND_MASTER_AWAY         CommandType = "m_ma"
)

// ------ MASTER SECTION ------
const (
    //MASTER_SECTION              = "master"

    // bound-id
    MASTER_BOUND_AGENT          = "m_ba"
    // master ip4 / ip6
    MASTER_IP4_ADDRESS          = "m_i4"
    MASTER_IP6_ADDRESS          = "m_i6"
    // master datetime
    MASTER_TIMESTAMP            = "m_ts"

    // TODO : Do we need this?
    //MASTER_HOSTNAME        = "m_hn"
)
