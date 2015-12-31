__author__ = 'stkim1'

# ------ VERSION ------
PC_PROTO            = 'pc_ver'
VERSION             = '1.0.0'

# ------ network interfaces ------
ADDRESS             = 'address'
NETMASK             = 'netmask'
BROADCS             = 'broadcast'
GATEWAY             = 'gateway'
NAMESRV             = 'dns-nameservers'
IFACE_KEYS          = [ADDRESS, NETMASK, BROADCS, GATEWAY, NAMESRV]

# ------ protocol definitions ------
MASTER_COMMAND_TYPE = 'pc_ma_ct'
COMMAND_FIX_BOUND   = 'ct_fix_bound'

# ------ MASTER SECTION ------
MASTER_SECTION      = 'master'

# bound-id
MASTER_BOUND_AGENT  = 'pc_ma_ba'
# master ip4 / ip6
MASTER_IP4_ADDRESS  = 'pc_ma_i4'
MASTER_IP6_ADDRESS  = 'pc_ma_i6'
# master hostname
MASTER_HOSTNAME     = 'pc_ma_hn'
# master datetime
MASTER_DATETIME     = 'pc_ma_dt'
MASTER_TIMEZONE     = 'pc_ma_tz'

# ------ SLAVE SECTION ------
SLAVE_SECTION       = 'slave'

# node looks for agent
SLAVE_LOOKUP_AGENT  = 'pc_sl_la'
SLAVE_NODE_MACADDR  = 'pc_sl_ma'
SLAVE_NODE_NAME     = 'pc_sl_nm'
SLAVE_TIMEZONE      = 'pc_sl_tz'
SLAVE_CLUSTER_MEMBERS = 'pc_sl_cl'

# ------ network configuration ------
POCKETCAST_GROUP    = '239.193.127.127'
PAGENT_SEND_PORT    = 10060
PAGENT_RECV_PORT    = 10061

POCKETCAST_SEND     = (POCKETCAST_GROUP, PAGENT_SEND_PORT)
POCKETCAST_RECV     = (POCKETCAST_GROUP, PAGENT_RECV_PORT)

# ------ configuration files ------
CONFIG_PATH         = '/etc/pocket/conf.ini'
NET_IFACE           = '/etc/network/interfaces'
HOSTNAME_FILE       = '/etc/hostname'
HOSTADDR_FILE       = '/etc/hosts'
HOST_TIMEZONE       = '/etc/timezone'
SALT_MINION_FILE    = '/etc/salt/minion_id'

# ------ SALT DEFAULT ------
SALT_MASTER         = 'salt'
PC_MASTER           = 'pc-master'

# ------ DEFAULT TIMEOUTS ------
UNBOUNDED_TIMEOUT   = 3
BOUNDED_TIMEOUT     = 10

# ------- POCKET EDITOR MARKER ------
POCKET_START = '# --------------- POCKETCLUSTER START ---------------'
POCKET_END   = '# ---------------  POCKETCLUSTER END  ---------------'