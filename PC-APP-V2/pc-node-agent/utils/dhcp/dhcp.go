package dhcp

const (
    PocketDHCPEventSocketPath string = "/var/run/pocketd.sock"
    PocketDHCPModeAgent string       = "pocket.dhcp.agent"
)

type PocketDhcpMeta struct {
    Reason                          string    `json:"reason, omitempty"                             msgpack:"reason, omitempty"`
    Interface                       string    `json:"interface, omitempty"                          msgpack:"interface, omitempty"`
    Medium                          string    `json:"medium, omitempty"                             msgpack:"medium, omitempty"`
    AliasIpAddress                  string    `json:"alias_ip_address, omitempty"                   msgpack:"alias_ip_address, omitempty"`
    IpAddress                       string    `json:"ip_address, omitempty"                         msgpack:"ip_address, omitempty"`
    HostName                        string    `json:"host_name, omitempty"                          msgpack:"host_name, omitempty"`
    NetworkNumber                   string    `json:"network_number, omitempty"                     msgpack:"network_number, omitempty"`
    SubnetMask                      string    `json:"subnet_mask, omitempty"                        msgpack:"subnet_mask, omitempty"`
    BroadcastAddress                string    `json:"broadcast_address, omitempty"                  msgpack:"broadcast_address, omitempty"`
    Routers                         string    `json:"routers, omitempty"                            msgpack:"routers, omitempty"`
    StaticRoutes                    string    `json:"static_routes, omitempty"                      msgpack:"static_routes, omitempty"`
    Rfc3442ClasslessStaticRoutes    string    `json:"rfc3442_classless_static_routes, omitempty"    msgpack:"rfc3442_classless_static_routes, omitempty"`
    DomainName                      string    `json:"domain_name, omitempty"                        msgpack:"domain_name, omitempty"`
    DomainSearch                    string    `json:"domain_search, omitempty"                      msgpack:"domain_search, omitempty"`
    DomainNameServers               string    `json:"domain_name_servers, omitempty"                msgpack:"domain_name_servers, omitempty"`
    NetbiosNameServers              string    `json:"netbios_name_servers, omitempty"               msgpack:"netbios_name_servers, omitempty"`
    NetbiosScope                    string    `json:"netbios_scope, omitempty"                      msgpack:"netbios_scope, omitempty"`
    NtpServers                      string    `json:"ntp_servers, omitempty"                        msgpack:"ntp_servers, omitempty"`
    Ip6Address                      string    `json:"ip6_address, omitempty"                        msgpack:"ip6_address, omitempty"`
    Ip6Prefix                       string    `json:"ip6_prefix, omitempty"                         msgpack:"ip6_prefix, omitempty"`
    Ip6Prefixlen                    string    `json:"ip6_prefixlen, omitempty"                      msgpack:"ip6_prefixlen, omitempty"`
    Dhcp6DomainSearch               string    `json:"dhcp6_domain_search, omitempty"                msgpack:"dhcp6_domain_search, omitempty"`
    Dhcp6NameServers                string    `json:"dhcp6_name_servers, omitempty"                 msgpack:"dhcp6_name_servers, omitempty"`
}

type PocketDhcpEvent struct {
    // the env variables without prefix
    Timestamp    string          `json:"timestamp"             msgpack:"timestamp"`
    Reason       string          `json:"reason"                msgpack:"reason"`
    Interface    string          `json:"interface"             msgpack:"interface"`
    Medium       string          `json:"medium, omitempty"     msgpack:"medium, omitempty"`

    // env meta variables with prefix 'old_', 'cur_', 'new_'
    Old        PocketDhcpMeta    `json:"old, omitempty"        msgpack:"old, inline, omitempty"`
    Current    PocketDhcpMeta    `json:"current, omitempty"    msgpack:"current, inline, omitempty"`
    New        PocketDhcpMeta    `json:"new, omitempty"        msgpack:"new, inline, omitempty"`

    // dhcp client requested checker
    Requested PocketDhcpMeta     `json:"requested, omitempty"  msgpack:"requested, inline, omitempty"`
}