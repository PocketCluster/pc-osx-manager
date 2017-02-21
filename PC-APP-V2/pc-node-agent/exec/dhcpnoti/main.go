package main

import (
    "encoding/json"
    "log"
    "os"
    "time"

    dhcp "github.com/stkim1/pc-node-agent/network"
)

func main() {
    var (
        writeToFile = false
    )

    dhcpEvent := new(dhcp.DhcpEvent)

    dhcpEvent.Timestamp                               = time.Now().Format(time.RFC3339)
    dhcpEvent.Reason                                  = os.Getenv("reason")
    dhcpEvent.Interface                               = os.Getenv("interface")
    dhcpEvent.Medium                                  = os.Getenv("medium")

    dhcpEvent.Old.Reason                              = os.Getenv("old_reason")
    dhcpEvent.Old.Interface                           = os.Getenv("old_interface")
    dhcpEvent.Old.Medium                              = os.Getenv("old_medium")
    dhcpEvent.Old.AliasIpAddress                      = os.Getenv("old_alias_ip_address")
    dhcpEvent.Old.IpAddress                           = os.Getenv("old_ip_address")
    dhcpEvent.Old.HostName                            = os.Getenv("old_host_name")
    dhcpEvent.Old.NetworkNumber                       = os.Getenv("old_network_number")
    dhcpEvent.Old.SubnetMask                          = os.Getenv("old_subnet_mask")
    dhcpEvent.Old.BroadcastAddress                    = os.Getenv("old_broadcast_address")
    dhcpEvent.Old.Routers                             = os.Getenv("old_routers")
    dhcpEvent.Old.StaticRoutes                        = os.Getenv("old_static_routes")
    dhcpEvent.Old.Rfc3442ClasslessStaticRoutes        = os.Getenv("old_rfc3442_classless_static_routes")
    dhcpEvent.Old.DomainName                          = os.Getenv("old_domain_name")
    dhcpEvent.Old.DomainSearch                        = os.Getenv("old_domain_search")
    dhcpEvent.Old.DomainNameServers                   = os.Getenv("old_domain_name_servers")
    dhcpEvent.Old.NetbiosNameServers                  = os.Getenv("old_netbios_name_servers")
    dhcpEvent.Old.NetbiosScope                        = os.Getenv("old_netbios_scope")
    dhcpEvent.Old.NtpServers                          = os.Getenv("old_ntp_servers")
    dhcpEvent.Old.Ip6Address                          = os.Getenv("old_ip6_address")
    dhcpEvent.Old.Ip6Prefix                           = os.Getenv("old_ip6_prefix")
    dhcpEvent.Old.Ip6Prefixlen                        = os.Getenv("old_ip6_prefixlen")
    dhcpEvent.Old.Dhcp6DomainSearch                   = os.Getenv("old_dhcp6_domain_search")
    dhcpEvent.Old.Dhcp6NameServers                    = os.Getenv("old_dhcp6_name_servers")

    dhcpEvent.Current.Reason                          = os.Getenv("cur_reason")
    dhcpEvent.Current.Interface                       = os.Getenv("cur_interface")
    dhcpEvent.Current.Medium                          = os.Getenv("cur_medium")
    dhcpEvent.Current.AliasIpAddress                  = os.Getenv("cur_alias_ip_address")
    dhcpEvent.Current.IpAddress                       = os.Getenv("cur_ip_address")
    dhcpEvent.Current.HostName                        = os.Getenv("cur_host_name")
    dhcpEvent.Current.NetworkNumber                   = os.Getenv("cur_network_number")
    dhcpEvent.Current.SubnetMask                      = os.Getenv("cur_subnet_mask")
    dhcpEvent.Current.BroadcastAddress                = os.Getenv("cur_broadcast_address")
    dhcpEvent.Current.Routers                         = os.Getenv("cur_routers")
    dhcpEvent.Current.StaticRoutes                    = os.Getenv("cur_static_routes")
    dhcpEvent.Current.Rfc3442ClasslessStaticRoutes    = os.Getenv("cur_rfc3442_classless_static_routes")
    dhcpEvent.Current.DomainName                      = os.Getenv("cur_domain_name")
    dhcpEvent.Current.DomainSearch                    = os.Getenv("cur_domain_search")
    dhcpEvent.Current.DomainNameServers               = os.Getenv("cur_domain_name_servers")
    dhcpEvent.Current.NetbiosNameServers              = os.Getenv("cur_netbios_name_servers")
    dhcpEvent.Current.NetbiosScope                    = os.Getenv("cur_netbios_scope")
    dhcpEvent.Current.NtpServers                      = os.Getenv("cur_ntp_servers")
    dhcpEvent.Current.Ip6Address                      = os.Getenv("cur_ip6_address")
    dhcpEvent.Current.Ip6Prefix                       = os.Getenv("cur_ip6_prefix")
    dhcpEvent.Current.Ip6Prefixlen                    = os.Getenv("cur_ip6_prefixlen")
    dhcpEvent.Current.Dhcp6DomainSearch               = os.Getenv("cur_dhcp6_domain_search")
    dhcpEvent.Current.Dhcp6NameServers                = os.Getenv("cur_dhcp6_name_servers")

    dhcpEvent.New.Reason                              = os.Getenv("new_reason")
    dhcpEvent.New.Interface                           = os.Getenv("new_interface")
    dhcpEvent.New.Medium                              = os.Getenv("new_medium")
    dhcpEvent.New.AliasIpAddress                      = os.Getenv("new_alias_ip_address")
    dhcpEvent.New.IpAddress                           = os.Getenv("new_ip_address")
    dhcpEvent.New.HostName                            = os.Getenv("new_host_name")
    dhcpEvent.New.NetworkNumber                       = os.Getenv("new_network_number")
    dhcpEvent.New.SubnetMask                          = os.Getenv("new_subnet_mask")
    dhcpEvent.New.BroadcastAddress                    = os.Getenv("new_broadcast_address")
    dhcpEvent.New.Routers                             = os.Getenv("new_routers")
    dhcpEvent.New.StaticRoutes                        = os.Getenv("new_static_routes")
    dhcpEvent.New.Rfc3442ClasslessStaticRoutes        = os.Getenv("new_rfc3442_classless_static_routes")
    dhcpEvent.New.DomainName                          = os.Getenv("new_domain_name")
    dhcpEvent.New.DomainSearch                        = os.Getenv("new_domain_search")
    dhcpEvent.New.DomainNameServers                   = os.Getenv("new_domain_name_servers")
    dhcpEvent.New.NetbiosNameServers                  = os.Getenv("new_netbios_name_servers")
    dhcpEvent.New.NetbiosScope                        = os.Getenv("new_netbios_scope")
    dhcpEvent.New.NtpServers                          = os.Getenv("new_ntp_servers")
    dhcpEvent.New.Ip6Address                          = os.Getenv("new_ip6_address")
    dhcpEvent.New.Ip6Prefix                           = os.Getenv("new_ip6_prefix")
    dhcpEvent.New.Ip6Prefixlen                        = os.Getenv("new_ip6_prefixlen")
    dhcpEvent.New.Dhcp6DomainSearch                   = os.Getenv("new_dhcp6_domain_search")
    dhcpEvent.New.Dhcp6NameServers                    = os.Getenv("new_dhcp6_name_servers")


    if writeToFile {
        jsonOut, err := json.Marshal(dhcpEvent)
        if err != nil {
            log.Print(err.Error())
        }
        out, err := os.OpenFile("/tmp/dh-client-env.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
        if err != nil {
            if !os.IsExist(err) {
                out, err = os.Create("/tmp/dh-client-env.log")
                if err != nil {
                    log.Fatal(err.Error())
                }
            } else {
                log.Fatal(err.Error())
            }
        }
        out.WriteString("\n--------------------- " + time.Now().Format("Mon, 2 Jan 2006 15:04:05 MST") + " ---------------------\n")
        out.WriteString(string(jsonOut))
        err = out.Close()
        if err != nil {
            log.Fatal(err.Error())
        }
    } else {
        json.NewEncoder(os.Stdout).Encode(dhcpEvent)
    }
}