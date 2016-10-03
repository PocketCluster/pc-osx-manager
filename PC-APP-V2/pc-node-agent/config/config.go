package config

import (
    "time"
)

// ------ NETWORK INTERFACES ------
const (
    ADDRESS             = "address"
    NETMASK             = "netmask"
    BROADCS             = "broadcast"
    GATEWAY             = "gateway"
    NAMESRV             = "dns-nameservers"
)
var IFACE_KEYS []string = []string{ADDRESS, NETMASK, BROADCS, GATEWAY, NAMESRV}

const (
    PAGENT_SEND_PORT    = 10060
    PAGENT_RECV_PORT    = 10061
)

// ------ CONFIGURATION FILES ------
const (
    // POCKET SPECIFIC CONFIG
    CONFIG_PATH         = "/etc/pocket/conf.ini"
    SLAVE_PRVATE_KEY    = "/etc/pocket/pki/slave.pem"
    SLAVE_PUBLIC_KEY    = "/etc/pocket/pki/slave.pub"
    SLAVE_PUBLIC_SSH    = "/etc/pocket/pki/slave.ssh"
    MASTER_PUBLIC_KEY   = "/etc/pocket/pki/master.pub"

    // HOST GENERAL CONFIG
    NET_IFACE           = "/etc/network/interfaces"
    HOSTNAME_FILE       = "/etc/hostname"
    HOSTADDR_FILE       = "/etc/hosts"
    HOST_TIMEZONE       = "/etc/timezone"
    RESOLVE_CONF        = "/etc/resolv.conf"
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
