package dhcp

import (
    "net"
    "os"
    "time"

    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "gopkg.in/vmihailenco/msgpack.v2"
    ps "github.com/mitchellh/go-ps"
)

func PocketDHCPEventAgent() {
    if os.Getuid() != 0 {
        log.Error(errors.WithStack(errors.New("invalid permission")))
        return
    }
    // dhclient-script pid
    sps, err := ps.FindProcess(os.Getppid())
    if err != nil {
        log.Error(errors.WithStack(err))
        return
    }
    if sps.Executable() != "dhclient-script" {
        log.Error(errors.WithStack(errors.New("incorrect preliminary executable")))
        return
    }
    // real dhclient pid
    rps, err := ps.FindProcess(sps.PPid())
    if err != nil {
        log.Error(errors.WithStack(err))
        return
    }
    if rps.Executable() != "dhclient" {
        log.Error(errors.WithStack(errors.New("invalid postliminary executable")))
        return
    }

    event := &PocketDhcpEvent{}

    event.Timestamp                               = time.Now().Format(time.RFC3339)
    event.Reason                                  = os.Getenv("reason")
    event.Interface                               = os.Getenv("interface")
    event.Medium                                  = os.Getenv("medium")

    event.Old.Reason                              = os.Getenv("old_reason")
    event.Old.Interface                           = os.Getenv("old_interface")
    event.Old.Medium                              = os.Getenv("old_medium")
    event.Old.AliasIpAddress                      = os.Getenv("old_alias_ip_address")
    event.Old.IpAddress                           = os.Getenv("old_ip_address")
    event.Old.HostName                            = os.Getenv("old_host_name")
    event.Old.NetworkNumber                       = os.Getenv("old_network_number")
    event.Old.SubnetMask                          = os.Getenv("old_subnet_mask")
    event.Old.BroadcastAddress                    = os.Getenv("old_broadcast_address")
    event.Old.Routers                             = os.Getenv("old_routers")
    event.Old.StaticRoutes                        = os.Getenv("old_static_routes")
    event.Old.Rfc3442ClasslessStaticRoutes        = os.Getenv("old_rfc3442_classless_static_routes")
    event.Old.DomainName                          = os.Getenv("old_domain_name")
    event.Old.DomainSearch                        = os.Getenv("old_domain_search")
    event.Old.DomainNameServers                   = os.Getenv("old_domain_name_servers")
    event.Old.NetbiosNameServers                  = os.Getenv("old_netbios_name_servers")
    event.Old.NetbiosScope                        = os.Getenv("old_netbios_scope")
    event.Old.NtpServers                          = os.Getenv("old_ntp_servers")
    event.Old.Ip6Address                          = os.Getenv("old_ip6_address")
    event.Old.Ip6Prefix                           = os.Getenv("old_ip6_prefix")
    event.Old.Ip6Prefixlen                        = os.Getenv("old_ip6_prefixlen")
    event.Old.Dhcp6DomainSearch                   = os.Getenv("old_dhcp6_domain_search")
    event.Old.Dhcp6NameServers                    = os.Getenv("old_dhcp6_name_servers")

    event.Current.Reason                          = os.Getenv("cur_reason")
    event.Current.Interface                       = os.Getenv("cur_interface")
    event.Current.Medium                          = os.Getenv("cur_medium")
    event.Current.AliasIpAddress                  = os.Getenv("cur_alias_ip_address")
    event.Current.IpAddress                       = os.Getenv("cur_ip_address")
    event.Current.HostName                        = os.Getenv("cur_host_name")
    event.Current.NetworkNumber                   = os.Getenv("cur_network_number")
    event.Current.SubnetMask                      = os.Getenv("cur_subnet_mask")
    event.Current.BroadcastAddress                = os.Getenv("cur_broadcast_address")
    event.Current.Routers                         = os.Getenv("cur_routers")
    event.Current.StaticRoutes                    = os.Getenv("cur_static_routes")
    event.Current.Rfc3442ClasslessStaticRoutes    = os.Getenv("cur_rfc3442_classless_static_routes")
    event.Current.DomainName                      = os.Getenv("cur_domain_name")
    event.Current.DomainSearch                    = os.Getenv("cur_domain_search")
    event.Current.DomainNameServers               = os.Getenv("cur_domain_name_servers")
    event.Current.NetbiosNameServers              = os.Getenv("cur_netbios_name_servers")
    event.Current.NetbiosScope                    = os.Getenv("cur_netbios_scope")
    event.Current.NtpServers                      = os.Getenv("cur_ntp_servers")
    event.Current.Ip6Address                      = os.Getenv("cur_ip6_address")
    event.Current.Ip6Prefix                       = os.Getenv("cur_ip6_prefix")
    event.Current.Ip6Prefixlen                    = os.Getenv("cur_ip6_prefixlen")
    event.Current.Dhcp6DomainSearch               = os.Getenv("cur_dhcp6_domain_search")
    event.Current.Dhcp6NameServers                = os.Getenv("cur_dhcp6_name_servers")

    event.New.Reason                              = os.Getenv("new_reason")
    event.New.Interface                           = os.Getenv("new_interface")
    event.New.Medium                              = os.Getenv("new_medium")
    event.New.AliasIpAddress                      = os.Getenv("new_alias_ip_address")
    event.New.IpAddress                           = os.Getenv("new_ip_address")
    event.New.HostName                            = os.Getenv("new_host_name")
    event.New.NetworkNumber                       = os.Getenv("new_network_number")
    event.New.SubnetMask                          = os.Getenv("new_subnet_mask")
    event.New.BroadcastAddress                    = os.Getenv("new_broadcast_address")
    event.New.Routers                             = os.Getenv("new_routers")
    event.New.StaticRoutes                        = os.Getenv("new_static_routes")
    event.New.Rfc3442ClasslessStaticRoutes        = os.Getenv("new_rfc3442_classless_static_routes")
    event.New.DomainName                          = os.Getenv("new_domain_name")
    event.New.DomainSearch                        = os.Getenv("new_domain_search")
    event.New.DomainNameServers                   = os.Getenv("new_domain_name_servers")
    event.New.NetbiosNameServers                  = os.Getenv("new_netbios_name_servers")
    event.New.NetbiosScope                        = os.Getenv("new_netbios_scope")
    event.New.NtpServers                          = os.Getenv("new_ntp_servers")
    event.New.Ip6Address                          = os.Getenv("new_ip6_address")
    event.New.Ip6Prefix                           = os.Getenv("new_ip6_prefix")
    event.New.Ip6Prefixlen                        = os.Getenv("new_ip6_prefixlen")
    event.New.Dhcp6DomainSearch                   = os.Getenv("new_dhcp6_domain_search")
    event.New.Dhcp6NameServers                    = os.Getenv("new_dhcp6_name_servers")

    event.Requested.Reason                        = os.Getenv("requested_reason")
    event.Requested.Interface                     = os.Getenv("requested_interface")
    event.Requested.Medium                        = os.Getenv("requested_medium")
    event.Requested.AliasIpAddress                = os.Getenv("requested_alias_ip_address")
    event.Requested.IpAddress                     = os.Getenv("requested_ip_address")
    event.Requested.HostName                      = os.Getenv("requested_host_name")
    event.Requested.NetworkNumber                 = os.Getenv("requested_network_number")
    event.Requested.SubnetMask                    = os.Getenv("requested_subnet_mask")
    event.Requested.BroadcastAddress              = os.Getenv("requested_broadcast_address")
    event.Requested.Routers                       = os.Getenv("requested_routers")
    event.Requested.StaticRoutes                  = os.Getenv("requested_static_routes")
    event.Requested.Rfc3442ClasslessStaticRoutes  = os.Getenv("requested_rfc3442_classless_static_routes")
    event.Requested.DomainName                    = os.Getenv("requested_domain_name")
    event.Requested.DomainSearch                  = os.Getenv("requested_domain_search")
    event.Requested.DomainNameServers             = os.Getenv("requested_domain_name_servers")
    event.Requested.NetbiosNameServers            = os.Getenv("requested_netbios_name_servers")
    event.Requested.NetbiosScope                  = os.Getenv("requested_netbios_scope")
    event.Requested.NtpServers                    = os.Getenv("requested_ntp_servers")
    event.Requested.Ip6Address                    = os.Getenv("requested_ip6_address")
    event.Requested.Ip6Prefix                     = os.Getenv("requested_ip6_prefix")
    event.Requested.Ip6Prefixlen                  = os.Getenv("requested_ip6_prefixlen")
    event.Requested.Dhcp6DomainSearch             = os.Getenv("requested_dhcp6_domain_search")
    event.Requested.Dhcp6NameServers              = os.Getenv("requested_dhcp6_name_servers")

    conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: PocketDHCPEventSocketPath, Net: "unix"})
    if err != nil {
        log.Error(errors.WithStack(err))
        return
    }
    defer conn.Close()

    msg, err := msgpack.Marshal(event)
    if err != nil {
        log.Error(errors.WithStack(err))
        return
    }

    _, err = conn.Write(msg)
    if err != nil {
        log.Error(errors.WithStack(err))
    }
}
