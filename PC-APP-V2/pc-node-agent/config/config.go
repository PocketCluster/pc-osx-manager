package config

import (
    "time"
)

// ------ CONFIG VERSION -------
const (
    SLAVE_CONFIG_KEY    = "config-version"
    SLAVE_CONFIG_VAL    = "1.0.1"
)

const (
    SLAVE_STATUS_KEY    = "binding-status"
    SLAVE_UNBOUNDED_VAL = "unbounded"
    SLAVE_BOUNDED_VAL   = "bounded"
)

// ------ NETWORK INTERFACES ------
const (
    SLAVE_ADDRESS_KEY   = "address"
    SLAVE_NETMASK_KEY   = "netmask"
    SLAVE_BROADCS_KEY   = "broadcast"
    SLAVE_GATEWAY_KEY   = "gateway"

    // TODO we might not need this
    //NAMESRV             = "dns-nameservers"
)
// TODO : we might not need this
//var IFACE_KEYS []string = []string{ADDRESS, NETMASK, BROADCS, GATEWAY}

const (
    PAGENT_SEND_PORT    = 10060
    PAGENT_RECV_PORT    = 10061
)

// ------ CONFIGURATION FILES ------
const (
    // POCKET SPECIFIC CONFIG
    SLAVE_CONFIG_PATH   = "/etc/pocket/node-conf.yaml"
    SLAVE_PRVATE_PATH   = "/etc/pocket/pki/slave.pem"
    SLAVE_PUBLIC_PATH   = "/etc/pocket/pki/slave.pub"
    SLAVE_SSH_PATH      = "/etc/pocket/pki/slave.ssh"
    MASTER_PUBLIC_PATH  = "/etc/pocket/pki/master.pub"

    // HOST GENERAL CONFIG
    NETWORK_IFACE_PATH  = "/etc/network/interfaces"
    HOSTNAME_PATH       = "/etc/hostname"
    HOSTADDR_PATH       = "/etc/hosts"
    HOST_TIMEZONE_PATH  = "/etc/timezone"
    RESOLVE_CONF_PATH   = "/etc/resolv.conf"
)

// ------ SALT DEFAULT ------
const (
    PC_MASTER           = "pc-master"
)

// ------ DEFAULT TIMEOUTS ------
const (
    UNBOUNDED_TIMEOUT   = 3 * time.Second
    BOUNDED_TIMEOUT     = 10 * time.Second
)

// ------- POCKET EDITOR MARKER ------
const (
    POCKET_START        = "// --------------- POCKETCLUSTER START ---------------"
    POCKET_END          = "// ---------------  POCKETCLUSTER END  ---------------"
)

type ConfigMasterSection struct {
    // Master Agent Specific String
    MasterBoundAgent    string                   `yaml:"master-binder-agent"`
    // Last Known IP4
    MasterIP4Address    string                   `yaml:"master-ip4-addr"`
    //MasterIP6Address    string
    //MasterHostName      string
    MasterTimeZone      string                   `yaml:"master-timezone"`
}

type ConfigSlaveSection struct {
    SlaveMacAddr        string                   `yaml:"slave-mac-addr"`
    SlaveNodeName       string                   `yaml:"slave-node-name"`
    SlaveIP4Addr        string                   `yaml:"slave-ip4-addr"`
    //SlaveIP6Addr        string
    SlaveNetMask        string                   `yaml:"slave-net-mask"`
    //SlaveBroadcast      string
    SlaveGateway        string                   `yaml:"slave-gateway"`
    //SlaveNameServ       string
}

type PocketSlaveConfig struct {
    ConfigVersion       string                   `yaml:"config-version"`
    BindingStatus       string                   `yaml:"binding-status"`
    MasterSection       *ConfigMasterSection      `yaml:"master-section",inline,flow`
    SlaveSection        *ConfigSlaveSection       `yaml:"slave-section",inline,flow`
}

func buildInitConfig() (*PocketSlaveConfig) {
    return &PocketSlaveConfig{
        ConfigVersion:SLAVE_CONFIG_VAL,
        BindingStatus:SLAVE_UNBOUNDED_VAL,
        MasterSection:&ConfigMasterSection{
        },
        SlaveSection:&ConfigSlaveSection{
        },
    }
}

