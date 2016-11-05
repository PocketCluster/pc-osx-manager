// linux setup util
{
    'have_ipv6_socket_ioctls': [
        
    ],
    'have_getnameinfo': True,
    'have_sockaddr_sa_len': False,
    'have_pf_netlink': True,
    'have_getifaddrs': True,
    'have_pf_route': False,
    'have_headers': [
        'netash/ash.h',
        'netatalk/at.h',
        'netax25/ax25.h',
        'neteconet/ec.h',
        'netipx/ipx.h',
        'netpacket/packet.h',
        'linux/irda.h',
        'linux/atm.h',
        'linux/llc.h',
        'linux/tipc.h',
        'linux/dn.h'
    ],
    'have_sysctl_ctl_net': False,
    'have_sockaddrs': [
        'at',
        'ax25',
        'in',
        'in6',
        'ipx',
        'un',
        'ash',
        'ec',
        'll',
        'atmpvc',
        'atmsvc',
        'dn',
        'irda',
        'llc'
    ]
}


#ifdef __linux__

#  ifndef HAVE_GETIFADDRS
#    define HAVE_GETIFADDRS 1
#  endif

#  ifndef HAVE_GETNAMEINFO
#    define HAVE_GETNAMEINFO 1
#  endif

#  ifdef HAVE_IPV6_SOCKET_IOCTLS
#    undef HAVE_IPV6_SOCKET_IOCTLS
#  endif

#  ifdef HAVE_SOCKET_IOCTLS
#    undef HAVE_SOCKET_IOCTLS
#  endif

#  ifndef HAVE_NETASH_ASH_H
#    define HAVE_NETASH_ASH_H 1
#  endif

#  ifndef HAVE_NETATALK_AT_H
#    define HAVE_NETATALK_AT_H 1
#  endif

#  ifndef HAVE_NETAX25_AX25_H
#    define HAVE_NETAX25_AX25_H 1
#  endif

#  ifndef HAVE_NETECONET_EC_H
#    define HAVE_NETECONET_EC_H 1
#  endif

#  ifndef HAVE_NETIPX_IPX_H
#    define HAVE_NETIPX_IPX_H 1
#  endif

#  ifndef HAVE_NETPACKET_PACKET_H
#    define HAVE_NETPACKET_PACKET_H 1
#  endif

#  ifndef HAVE_LINUX_IRDA_H
#    define HAVE_LINUX_IRDA_H 1
#  endif

#  ifndef HAVE_LINUX_ATM_H
#    define HAVE_LINUX_ATM_H 1
#  endif

#  ifndef HAVE_LINUX_LLC_H
#    define HAVE_LINUX_LLC_H 1
#  endif

#  ifndef HAVE_LINUX_TIPC_H
#    define HAVE_LINUX_TIPC_H 1
#  endif

#  ifndef HAVE_LINUX_DN_H
#    define HAVE_LINUX_DN_H 1
#  endif

#  ifdef HAVE_SOCKADDR_SA_LEN
#    undef HAVE_SOCKADDR_SA_LEN
#  endif

#  ifndef HAVE_SOCKADDR_AT
#    define HAVE_SOCKADDR_AT 1
#  endif

#  ifndef HAVE_SOCKADDR_AX25
#    define HAVE_SOCKADDR_AX25 1
#  endif

#  ifndef HAVE_SOCKADDR_IN
#    define HAVE_SOCKADDR_IN 1
#  endif

#  ifndef HAVE_SOCKADDR_IN6
#    define HAVE_SOCKADDR_IN6 1
#  endif

#  ifndef HAVE_SOCKADDR_IPX
#    define HAVE_SOCKADDR_IPX 1
#  endif

#  ifndef HAVE_SOCKADDR_UN
#    define HAVE_SOCKADDR_UN 1
#  endif

#  ifndef HAVE_SOCKADDR_EC
#    define HAVE_SOCKADDR_EC 1
#  endif

#  ifndef HAVE_SOCKADDR_LL
#    define HAVE_SOCKADDR_LL 1
#  endif

#  ifndef HAVE_SOCKADDR_ATMPVC
#    define HAVE_SOCKADDR_ATMPVC 1
#  endif

#  ifndef HAVE_SOCKADDR_ATMSVC
#    define HAVE_SOCKADDR_ATMSVC 1
#  endif

#  ifndef HAVE_SOCKADDR_DN
#    define HAVE_SOCKADDR_DN 1
#  endif

#  ifndef HAVE_SOCKADDR_IRDA
#    define HAVE_SOCKADDR_IRDA 1
#  endif

#  ifndef HAVE_SOCKADDR_LLC
#    define HAVE_SOCKADDR_LLC 1
#  endif

#  ifdef HAVE_PF_ROUTE
#    undef HAVE_PF_ROUTE
#  endif

#  ifdef HAVE_SYSCTL_CTL_NET
#    undef HAVE_SYSCTL_CTL_NET
#  endif

#  ifndef HAVE_PF_NETLINK
#    define HAVE_PF_NETLINK 1
#  endif

#  endif



// OSX
{
    'have_ipv6_socket_ioctls': [
        'SIOCGIFAFLAG_IN6'
    ],
    'have_getnameinfo': True,
    'have_sockaddr_sa_len': True,
    'have_pf_netlink': False,
    'have_getifaddrs': True,
    'have_pf_route': True,
    'have_headers': [
        'net/if_dl.h'
    ],
    'have_sysctl_ctl_net': True
}


#  elif defined(__APPLE__) || defined(__MACH__)

#  ifndef HAVE_GETIFADDRS
#    define HAVE_GETIFADDRS 1
#  endif

#  ifndef HAVE_GETNAMEINFO
#    define HAVE_GETNAMEINFO 1
#  endif

#  ifndef HAVE_IPV6_SOCKET_IOCTLS
#    define HAVE_IPV6_SOCKET_IOCTLS 1
#  endif

#  ifndef HAVE_SIOCGIFAFLAG_IN6
#    define HAVE_SIOCGIFAFLAG_IN6 1
#  endif

#  ifdef HAVE_SOCKET_IOCTLS
#    undef HAVE_SOCKET_IOCTLS
#  endif

#  ifndef HAVE_NET_IF_DL_H
#    define HAVE_NET_IF_DL_H 1
#  endif

#  ifndef HAVE_SOCKADDR_SA_LEN
#    define HAVE_SOCKADDR_SA_LEN 1
#  endif

#  ifndef HAVE_PF_ROUTE
#    define HAVE_PF_ROUTE 1
#  endif

#  ifndef HAVE_SYSCTL_CTL_NET
#    define HAVE_SYSCTL_CTL_NET 1
#  endif

#  ifdef HAVE_PF_NETLINK
#    undef HAVE_PF_NETLINK
#  endif

