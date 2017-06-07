package main

import (
    "encoding/json"
    "flag"
    "net"
    "os"
    "time"

    "github.com/stkim1/pc-node-agent/utils/dhcp"
    "gopkg.in/vmihailenco/msgpack.v2"
    log "github.com/Sirupsen/logrus"
    "github.com/pkg/errors"
    "github.com/davecgh/go-spew/spew"
    process "github.com/mitchellh/go-ps"
)

const (
    modeDhcpAgent       = "dhcpagent"
    devJsonPrint        = "jsonprint"
)

var (
    mode    = flag.String("mode", "", "Execution mode")
    dev     = flag.String("dev",  "", "Developement")
)

func dhcpAgent() {
    if os.Getuid() != 0 {
        log.Error(errors.WithStack(errors.New("Insufficient Permission")))
        return
    }
    // dhclient-script pid
    sps, err := process.FindProcess(os.Getppid())
    if err != nil {
        log.Error(errors.WithStack(err))
        return
    }
    if sps.Executable() != "dhclient-script" {
        log.Error(errors.WithStack(errors.New("Incorrect preliminary executable")))
        return
    }
    // real dhclient pid
    rps, err := process.FindProcess(sps.PPid())
    if err != nil {
        log.Error(errors.WithStack(err))
        return
    }
    if rps.Executable() != "dhclient" {
        log.Error(errors.WithStack(errors.New("Incorrect postliminary executable")))
        return
    }

    dhcpEvent := &dhcp.DhcpEvent{}

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

    dhcpEvent.Requested.Reason                        = os.Getenv("requested_reason")
    dhcpEvent.Requested.Interface                     = os.Getenv("requested_interface")
    dhcpEvent.Requested.Medium                        = os.Getenv("requested_medium")
    dhcpEvent.Requested.AliasIpAddress                = os.Getenv("requested_alias_ip_address")
    dhcpEvent.Requested.IpAddress                     = os.Getenv("requested_ip_address")
    dhcpEvent.Requested.HostName                      = os.Getenv("requested_host_name")
    dhcpEvent.Requested.NetworkNumber                 = os.Getenv("requested_network_number")
    dhcpEvent.Requested.SubnetMask                    = os.Getenv("requested_subnet_mask")
    dhcpEvent.Requested.BroadcastAddress              = os.Getenv("requested_broadcast_address")
    dhcpEvent.Requested.Routers                       = os.Getenv("requested_routers")
    dhcpEvent.Requested.StaticRoutes                  = os.Getenv("requested_static_routes")
    dhcpEvent.Requested.Rfc3442ClasslessStaticRoutes  = os.Getenv("requested_rfc3442_classless_static_routes")
    dhcpEvent.Requested.DomainName                    = os.Getenv("requested_domain_name")
    dhcpEvent.Requested.DomainSearch                  = os.Getenv("requested_domain_search")
    dhcpEvent.Requested.DomainNameServers             = os.Getenv("requested_domain_name_servers")
    dhcpEvent.Requested.NetbiosNameServers            = os.Getenv("requested_netbios_name_servers")
    dhcpEvent.Requested.NetbiosScope                  = os.Getenv("requested_netbios_scope")
    dhcpEvent.Requested.NtpServers                    = os.Getenv("requested_ntp_servers")
    dhcpEvent.Requested.Ip6Address                    = os.Getenv("requested_ip6_address")
    dhcpEvent.Requested.Ip6Prefix                     = os.Getenv("requested_ip6_prefix")
    dhcpEvent.Requested.Ip6Prefixlen                  = os.Getenv("requested_ip6_prefixlen")
    dhcpEvent.Requested.Dhcp6DomainSearch             = os.Getenv("requested_dhcp6_domain_search")
    dhcpEvent.Requested.Dhcp6NameServers              = os.Getenv("requested_dhcp6_name_servers")

    conn, err := net.DialUnix("unix", nil, &net.UnixAddr{dhcp.DHCPEventSocketPath, "unix"})
    if err != nil {
        log.Error(errors.WithStack(err))
        return
    }
    defer conn.Close()

    msg, err := msgpack.Marshal(dhcpEvent)
    if err != nil {
        log.Error(errors.WithStack(err))
        return
    }

    _, err = conn.Write(msg)
    if err != nil {
        log.Error(errors.WithStack(err))
    }

    if len(*dev) != 0 && *dev == devJsonPrint {
        json.NewEncoder(os.Stdout).Encode(struct {
            Event         *dhcp.DhcpEvent    `json:"dhcp_event, omitempty"`
            Pid           int                `json:"dhcp_pid, omitempty"`
            Executable    string             `json:"dhcp_executable, omitempty"`
        }{
            Event:        dhcpEvent,
            Pid:          os.Getpid(),
            Executable:   rps.Executable(),
        })
    }
}

func pocketDaemon() {
    log.Info("Pocket Daemon Started...")

    buf := make([]byte, 20480)
    dhcpEvent := &dhcp.DhcpEvent{}

    // firstly clear off previous socket
    os.Remove(dhcp.DHCPEventSocketPath)
    listen, err := net.ListenUnix("unix", &net.UnixAddr{dhcp.DHCPEventSocketPath, "unix"})
    if err != nil {
        log.Error(errors.WithStack(err))
        return
    }
    defer os.Remove(dhcp.DHCPEventSocketPath)
    defer listen.Close()

    for {
        conn, err := listen.AcceptUnix()
        if err != nil {
            log.Error(errors.WithStack(err))
            continue
        }
        count, err := conn.Read(buf)
        if err != nil {
            log.Error(errors.WithStack(err))
            continue
        }
        err = msgpack.Unmarshal(buf[0:count], dhcpEvent)
        if err != nil {
            log.Error(errors.WithStack(err))
            continue
        }

        log.Info(spew.Sdump(dhcpEvent))

        err = conn.Close()
        if err != nil {
            log.Error(errors.WithStack(err))
            continue
        }
    }
}

func main() {
    flag.Parse()

    switch *mode {
    case modeDhcpAgent:
        dhcpAgent()
        break
    default:
        pocketDaemon()
    }
}