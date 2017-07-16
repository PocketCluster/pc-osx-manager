/*
 * Copyright (c) 2007-2014 Alastair Houghton

 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:

 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.

 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

#include "netifaces.h"

#ifndef WIN32

#  include <unistd.h>

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

#  ifndef HAVE_SOCKADDR_ASH
#    define HAVE_SOCKADDR_ASH 1
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

#  elif defined(__APPLE__) || defined(__MACH__) //ifdef __linux__

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

#  endif // __APPLE__

#  include <stddef.h>
#  include <stdio.h>
#  include <string.h>
#  include <stdlib.h>
#  include <errno.h>

#  include <sys/types.h>
#  include <sys/socket.h>
#  include <net/if.h>
#  include <netdb.h>

#  if HAVE_PF_ROUTE
#    include <net/route.h>
#  endif

#  if HAVE_SYSCTL_CTL_NET
#    include <sys/sysctl.h>
#    include <net/route.h>
#  endif

#  if HAVE_PF_NETLINK
#    include <asm/types.h>
#    include <linux/netlink.h>
#    include <linux/rtnetlink.h>
#    include <arpa/inet.h>
#  endif

#  if HAVE_GETIFADDRS
#    if HAVE_IPV6_SOCKET_IOCTLS
#      include <sys/ioctl.h>
#      include <netinet/in.h>
#      include <netinet/in_var.h>
#    endif
#  endif

#  if HAVE_SOCKET_IOCTLS
#    include <sys/ioctl.h>
#    include <netinet/in.h>
#    include <arpa/inet.h>
#    if defined(__sun)
#      include <unistd.h>
#      include <stropts.h>
#      include <sys/sockio.h>
#    endif
#  endif /* HAVE_SOCKET_IOCTLS */

/* For logical interfaces support we convert all names to same name prefixed
   with l */
#  if HAVE_SIOCGLIFNUM
#     define CNAME(x) l##x
#  else
#   define CNAME(x) x
#  endif

#  if HAVE_NET_IF_DL_H
#    include <net/if_dl.h>
#  endif

/* For the benefit of stupid platforms (Linux), include all the sockaddr
   definitions we can lay our hands on. It can also be useful for the benefit
   of another stupid platform (FreeBSD, see PR 152036). */
#include <netinet/in.h>
#  if HAVE_NETASH_ASH_H
#    include <netash/ash.h>
#  endif
#  if HAVE_NETATALK_AT_H
#    include <netatalk/at.h>
#  endif
#  if HAVE_NETAX25_AX25_H
#    include <netax25/ax25.h>
#  endif
#  if HAVE_NETECONET_EC_H
#    include <neteconet/ec.h>
#  endif
#  if HAVE_NETIPX_IPX_H
#    include <netipx/ipx.h>
#  endif
#  if HAVE_NETPACKET_PACKET_H
#    include <netpacket/packet.h>
#  endif
#  if HAVE_NETROSE_ROSE_H
#    include <netrose/rose.h>
#  endif
#  if HAVE_LINUX_IRDA_H
#    include <linux/irda.h>
#  endif
#  if HAVE_LINUX_ATM_H
#    include <linux/atm.h>
#  endif
#  if HAVE_LINUX_LLC_H
#    include <linux/llc.h>
#  endif
#  if HAVE_LINUX_TIPC_H
#    include <linux/tipc.h>
#  endif
#  if HAVE_LINUX_DN_H
#    include <linux/dn.h>
#  endif

/* Map address families to sizes of sockaddr structs */
static int af_to_len(int af) 
{
  switch (af) {
  case AF_INET: return sizeof (struct sockaddr_in);
#if defined(AF_INET6) && HAVE_SOCKADDR_IN6
  case AF_INET6: return sizeof (struct sockaddr_in6);
#endif
#if defined(AF_AX25) && HAVE_SOCKADDR_AX25
#  if defined(AF_NETROM)
  case AF_NETROM: /* I'm assuming this is carried over x25 */
#  endif
  case AF_AX25: return sizeof (struct sockaddr_ax25);
#endif
#if defined(AF_IPX) && HAVE_SOCKADDR_IPX
  case AF_IPX: return sizeof (struct sockaddr_ipx);
#endif
#if defined(AF_APPLETALK) && HAVE_SOCKADDR_AT
  case AF_APPLETALK: return sizeof (struct sockaddr_at);
#endif
#if defined(AF_ATMPVC) && HAVE_SOCKADDR_ATMPVC
  case AF_ATMPVC: return sizeof (struct sockaddr_atmpvc);
#endif
#if defined(AF_ATMSVC) && HAVE_SOCKADDR_ATMSVC
  case AF_ATMSVC: return sizeof (struct sockaddr_atmsvc);
#endif
#if defined(AF_X25) && HAVE_SOCKADDR_X25
  case AF_X25: return sizeof (struct sockaddr_x25);
#endif
#if defined(AF_ROSE) && HAVE_SOCKADDR_ROSE
  case AF_ROSE: return sizeof (struct sockaddr_rose);
#endif
#if defined(AF_DECnet) && HAVE_SOCKADDR_DN
  case AF_DECnet: return sizeof (struct sockaddr_dn);
#endif
#if defined(AF_PACKET) && HAVE_SOCKADDR_LL
  case AF_PACKET: return sizeof (struct sockaddr_ll);
#endif
#if defined(AF_ASH) && HAVE_SOCKADDR_ASH
  case AF_ASH: return sizeof (struct sockaddr_ash);
#endif
#if defined(AF_ECONET) && HAVE_SOCKADDR_EC
  case AF_ECONET: return sizeof (struct sockaddr_ec);
#endif
#if defined(AF_IRDA) && HAVE_SOCKADDR_IRDA
  case AF_IRDA: return sizeof (struct sockaddr_irda);
#endif
#if defined(AF_LINK) && HAVE_SOCKADDR_DL
  case AF_LINK: return sizeof (struct sockaddr_dl);
#endif
  }
  return sizeof (struct sockaddr);
}

#if !HAVE_SOCKADDR_SA_LEN
  #define SA_LEN(sa)      af_to_len(sa->sa_family)
  #if HAVE_SIOCGLIFNUM
    #define SS_LEN(sa)      af_to_len(sa->ss_family)
  #else
    #define SS_LEN(sa)      SA_LEN(sa)
  #endif
#else
  #define SA_LEN(sa)      sa->sa_len
#endif /* !HAVE_SOCKADDR_SA_LEN */

#  if HAVE_GETIFADDRS
#    include <ifaddrs.h>
#  endif /* HAVE_GETIFADDRS */

#  if !HAVE_GETIFADDRS && (!HAVE_SOCKET_IOCTLS || !HAVE_SIOCGIFCONF)
/* If the platform doesn't define, what we need, barf.  If you're seeing this,
   it means you need to write suitable code to retrieve interface information
   on your system. */
#    error You need to add code for your platform.
#  endif

#endif

#ifndef TRUE
#define TRUE 1
#endif

#ifndef FALSE
#define FALSE 0
#endif

/* On systems without AF_LINK (Windows, for instance), define it anyway, but
   give it a crazy value.  On Linux, which has AF_PACKET but not AF_LINK,
   define AF_LINK as the latter instead. */
#ifndef AF_LINK
#  ifdef AF_PACKET
#    define AF_LINK  AF_PACKET
#  else
#    define AF_LINK  -1000
#  endif
#  define HAVE_AF_LINK 0
#else
#  define HAVE_AF_LINK 1
#endif

/* -- Utility Functions ----------------------------------------------------- */

#if !defined(WIN32)
#if  !HAVE_GETNAMEINFO
#undef getnameinfo
#undef NI_NUMERICHOST

#define getnameinfo our_getnameinfo
#define NI_NUMERICHOST 1

/* A very simple getnameinfo() for platforms without */
static int
getnameinfo (const struct sockaddr *addr, int addr_len,
             char *buffer, int buflen,
             char *buf2, int buf2len,
             int flags)
{
  switch (addr->sa_family) {
  case AF_INET:
    {
      const struct sockaddr_in *sin = (struct sockaddr_in *)addr;
      const unsigned char *bytes = (unsigned char *)&sin->sin_addr.s_addr;
      char tmpbuf[20];

      sprintf (tmpbuf, "%d.%d.%d.%d",
               bytes[0], bytes[1], bytes[2], bytes[3]);

      strncpy (buffer, tmpbuf, buflen);
    }
    break;
#ifdef AF_INET6
  case AF_INET6:
    {
      const struct sockaddr_in6 *sin = (const struct sockaddr_in6 *)addr;
      const unsigned char *bytes = sin->sin6_addr.s6_addr;
      int n;
      char tmpbuf[80], *ptr = tmpbuf;
      int done_double_colon = FALSE;
      int colon_mode = FALSE;

      for (n = 0; n < 8; ++n) {
        unsigned char b1 = bytes[2 * n];
        unsigned char b2 = bytes[2 * n + 1];

        if (b1) {
          if (colon_mode) {
            colon_mode = FALSE;
            *ptr++ = ':';
          }
          sprintf (ptr, "%x%02x", b1, b2);
          ptr += strlen (ptr);
          *ptr++ = ':';
        } else if (b2) {
          if (colon_mode) {
            colon_mode = FALSE;
            *ptr++ = ':';
          }
          sprintf (ptr, "%x", b2);
          ptr += strlen (ptr);
          *ptr++ = ':';
        } else {
          if (!colon_mode) {
            if (done_double_colon) {
              *ptr++ = '0';
              *ptr++ = ':';
            } else {
              if (n == 0)
                *ptr++ = ':';
              colon_mode = TRUE;
              done_double_colon = TRUE;
            }
          }
        }
      }
      if (colon_mode) {
        colon_mode = FALSE;
        *ptr++ = ':';
        *ptr++ = '\0';
      } else {
        *--ptr = '\0';
      }

      strncpy (buffer, tmpbuf, buflen);
    }
    break;
#endif /* AF_INET6 */
  default:
    return -1;
  }

  return 0;
}
#endif

static int
string_from_sockaddr (struct sockaddr *addr,
                      char *buffer,
                      int buflen)
{
  struct sockaddr* bigaddr = 0;
  int failure;
  struct sockaddr* gniaddr;
  socklen_t gnilen;

  if (!addr || addr->sa_family == AF_UNSPEC)
    return -1;

  if (SA_LEN(addr) < af_to_len(addr->sa_family)) {
    /* Sometimes ifa_netmask can be truncated. So let's detruncate it.  FreeBSD
       PR: kern/152036: getifaddrs(3) returns truncated sockaddrs for netmasks
       -- http://www.freebsd.org/cgi/query-pr.cgi?pr=152036 */
    gnilen = af_to_len(addr->sa_family);
    bigaddr = calloc(1, gnilen);
    if (!bigaddr)
      return -1;
    memcpy(bigaddr, addr, SA_LEN(addr));
#if HAVE_SOCKADDR_SA_LEN
    bigaddr->sa_len = gnilen;
#endif
    gniaddr = bigaddr;
  } else {
    gnilen = SA_LEN(addr);
    gniaddr = addr;
  }

  failure = getnameinfo (gniaddr, gnilen,
                         buffer, buflen,
                         NULL, 0,
                         NI_NUMERICHOST);

  if (bigaddr) {
    free(bigaddr);
    bigaddr = 0;
  }

  if (failure) {
    size_t n, len;
    char *ptr;
    const char *data;
      
    len = SA_LEN(addr);

#if HAVE_AF_LINK
    /* BSD-like systems have AF_LINK */
    if (addr->sa_family == AF_LINK) {
      struct sockaddr_dl *dladdr = (struct sockaddr_dl *)addr;
      len = dladdr->sdl_alen;
      data = LLADDR(dladdr);
    } else {
#endif
#if defined(AF_PACKET)
      /* Linux has AF_PACKET instead */
      if (addr->sa_family == AF_PACKET) {
        struct sockaddr_ll *lladdr = (struct sockaddr_ll *)addr;
        len = lladdr->sll_halen;
        data = (const char *)lladdr->sll_addr;
      } else {
#endif
        /* We don't know anything about this sockaddr, so just display
           the entire data area in binary. */
        len -= (sizeof (struct sockaddr) - sizeof (addr->sa_data));
        data = addr->sa_data;
#if defined(AF_PACKET)
      }
#endif
#if HAVE_AF_LINK
    }
#endif

    if (buflen < 3 * len)
      return -1;

    ptr = buffer;
    buffer[0] = '\0';

    for (n = 0; n < len; ++n) {
      sprintf (ptr, "%02x:", data[n] & 0xff);
      ptr += 3;
    }
    if (len)
      *--ptr = '\0';
  }

  if (!buffer[0])
    return -1;

  return 0;
}

/* Tries to format in CIDR form where possible; falls back to using
   string_from_sockaddr(). */
static int
string_from_netmask (struct sockaddr *addr,
                     char *buffer,
                     int buflen)
{
#ifdef AF_INET6
  if (addr && addr->sa_family == AF_INET6) {
    struct sockaddr_in6 *sin6 = (struct sockaddr_in6 *)addr;
    unsigned n = 16;
    unsigned zeroes = 0;
    unsigned prefix;
    unsigned bytes;
    char *bufptr = buffer;
    char *bufend = buffer + buflen;
    char pfxbuf[16];

    while (n--) {
      unsigned char byte = sin6->sin6_addr.s6_addr[n];

      /* We need to count the rightmost zeroes */
      unsigned char x = byte;
      unsigned zx = 8;

      x &= -x;
      if (x)
        --zx;
      if (x & 0x0f)
        zx -= 4;
      if (x & 0x03)
        zx -= 2;
      if (x & 0x05)
        zx -= 1;

      zeroes += zx;

      if (byte)
        break;
    }

    prefix = 128 - zeroes;
    bytes = 2 * ((prefix + 15) / 16);

    for (n = 0; n < bytes; ++n) {
      unsigned char byte = sin6->sin6_addr.s6_addr[n];
      char ch1, ch2;

      if (n && !(n & 1)) {
        if (bufptr < bufend)
          *bufptr++ = ':';
      }

      ch1 = '0' + (byte >> 4);
      if (ch1 > '9')
        ch1 += 'a' - '0' - 10;
      ch2 = '0' + (byte & 0xf);
      if (ch2 > '9')
        ch2 += 'a' - '0' - 10;

      if (bufptr < bufend)
        *bufptr++ = ch1;
      if (bufptr < bufend)
        *bufptr++ = ch2;
    }

    if (bytes < 16) {
      if (bufend - bufptr > 2) {
        *bufptr++ = ':';
        *bufptr++ = ':';
      }
    }

    sprintf (pfxbuf, "/%u", prefix);

    if (bufend - bufptr > strlen(pfxbuf))
      strcpy (bufptr, pfxbuf);

    if (buflen)
      buffer[buflen - 1] = '\0';

    return 0;
  }
#endif

  return string_from_sockaddr(addr, buffer, buflen);
}
#endif /* !defined(WIN32) */

/* -- ifaddresses() --------------------------------------------------------- */

static bool
append_address(Address** results, Address* address) {
    Address* addr = *results;
    if (address == NULL) {
        return false;
    }
    if (*results == NULL) {
        *results = address;
        return true;
    }
    while(addr->next != NULL) {
        // as we have the information already, we do not need to add it again.
        if (addr == address || (addr->addr != NULL && address->addr != NULL && strcmp(addr->addr, address->addr) == 0) ) {
            return true;
        }
        addr = addr->next;
    }
    addr->next = address;
    return true;
}

static void
release_addresses(Address** results) {
    Address *head = *results, *tail = NULL;
    if (*results == NULL) {
        return;
    }
    // traseverse linked list
    while(head != NULL) {
        tail = head;
        head = head->next;
        
        if (tail->addr != NULL) {
            free(tail->addr);
        }
        if (tail->netmask != NULL) {
            free(tail->netmask);
        }
        if (tail->broadcast != NULL) {
            free(tail->broadcast);
        }
        if (tail->peer != NULL) {
            free(tail->peer);
        }
        tail->next = NULL;
        free(tail);
        tail = NULL;
    }
    *results = NULL;
}

/*!
	@function ifaddrs
	@discussion Returns ipv4/ipv6 addresses (dotted format) linked to the interface.
	@param ifname The network interface name.
            results List where address should go
	@result The list of ipv4 addresses linked the interface;
            NULL if no ipv4 addresses are supported or linked.
 */

static int
ifaddrs (Address **results, const char *ifname)
{
  int found = FALSE;
#if HAVE_GETIFADDRS
  struct ifaddrs *addrs = NULL;
  struct ifaddrs *addr = NULL;
#endif
  
  if (ifname == NULL || strlen(ifname) == 0) {
    return EINVAL;
  }

#if HAVE_GETIFADDRS
  /* .. UNIX, with getifaddrs() ............................................. */

  if (getifaddrs (&addrs) < 0) {
    return ENETUNREACH;
  }

  for (addr = addrs; addr; addr = addr->ifa_next) {
    char buffer[256];
    Address *address = NULL;

    if (strcmp (addr->ifa_name, ifname) != 0)
      continue;
 
    /* We mark the interface as found, even if there are no addresses;
       this results in sensible behaviour for these few cases. */
    found = TRUE;

    /* Sometimes there are records without addresses (e.g. in the case of a
       dial-up connection via ppp, which on Linux can have a link address
       record with no actual address).  We skip these as they aren't useful.
       Thanks to Christian Kauhaus for reporting this issue. */
    if (!addr->ifa_addr)
      continue;
      
    /* As it is ready to pull address inforation, we'll alloc an address */
    address = (Address *) calloc(1, sizeof (Address));
    
#if HAVE_IPV6_SOCKET_IOCTLS
    /* For IPv6 addresses we try to get the flags. */
    if (addr->ifa_addr->sa_family == AF_INET6) {
      struct sockaddr_in6 *sin;
      struct in6_ifreq ifr6;
      
      int sock6 = socket (AF_INET6, SOCK_DGRAM, 0);

      if (sock6 < 0) {
        freeifaddrs (addrs);
        free(address);
        return ENETUNREACH;
      }
      
      sin = (struct sockaddr_in6 *)addr->ifa_addr;
      strncpy (ifr6.ifr_name, addr->ifa_name, IFNAMSIZ);
      ifr6.ifr_addr = *sin;
      
      if (ioctl (sock6, SIOCGIFAFLAG_IN6, &ifr6) >= 0) {
        address->flags = ifr6.ifr_ifru.ifru_flags6;
      }

      close (sock6);
    } else
#endif /* HAVE_IPV6_SOCKET_IOCTLS */
    {
      address->flags = addr->ifa_flags;
    }

    if (string_from_sockaddr (addr->ifa_addr, buffer, sizeof (buffer)) == 0) {
      size_t addr_len = strlen(buffer);
      char* addr_str = (char *) malloc (addr_len * sizeof(char));
      memcpy(addr_str, buffer, addr_len);
      address->addr = addr_str;
    }

    if (string_from_netmask (addr->ifa_netmask, buffer, sizeof (buffer)) == 0) {
      size_t netmask_len = strlen(buffer);
      char* netmask_str = (char *) malloc (netmask_len * sizeof(char));
      memcpy(netmask_str, buffer, netmask_len);
      address->netmask = netmask_str;
    }
      
    if (string_from_sockaddr (addr->ifa_broadaddr, buffer, sizeof (buffer)) == 0) {
      size_t braddr_len = strlen(buffer);
      char* braddr_str = (char *) malloc (braddr_len * sizeof(char));
      memcpy(braddr_str, buffer, braddr_len);
      address->broadcast = braddr_str;
    }

    /* Cygwin's implementation of getaddrinfo() is buggy and returns broadcast
       addresses for 169.254.0.0/16.  Nix them here. */
    if (addr->ifa_addr->sa_family == AF_INET) {
      struct sockaddr_in *sin = (struct sockaddr_in *)addr->ifa_addr;
      if ((ntohl(sin->sin_addr.s_addr) & 0xffff0000) == 0xa9fe0000) {
        if (address->broadcast != NULL && strlen(address->broadcast) != 0) {
          free(address->broadcast);
          address->broadcast = NULL;
        }
      }
    }

    if (address->broadcast != NULL && strlen(address->broadcast) != 0) {
      if (addr->ifa_flags & (IFF_POINTOPOINT | IFF_LOOPBACK)) {
        address->peer = address->broadcast;
        address->broadcast = NULL;
      }
    }
      
    address->family = addr->ifa_addr->sa_family;
    
    append_address(results, address);
  }

  freeifaddrs (addrs);
#elif HAVE_SOCKET_IOCTLS
  /* .. UNIX, with SIOC ioctls() ............................................ */
  
  int sock = socket(AF_INET, SOCK_DGRAM, 0);

  if (sock < 0) {
    return ENETUNREACH;
  }

  struct CNAME(ifreq) ifr;
  int is_p2p = FALSE;
  char buffer[256];
  /* As it is ready to pull address inforation, we'll alloc an address */
  Address *address = (Address *) calloc(1, sizeof (Address));
    
  strncpy (ifr.CNAME(ifr_name), ifname, IFNAMSIZ);

#if HAVE_SIOCGIFHWADDR
  if (ioctl (sock, SIOCGIFHWADDR, &ifr) == 0) {
    found = TRUE;

    if (string_from_sockaddr ((struct sockaddr *)&ifr.CNAME(ifr_addr), buffer, sizeof (buffer)) == 0) {
      size_t addr_len = strlen(buffer);
      char* addr_str = (char *) malloc (addr_len * sizeof(char));
      memcpy(addr_str, buffer, addr_len);
      address->addr = addr_str;
        
      address->family = AF_LINK;

      append_address(results, address);
    }
  }
#endif

#if HAVE_SIOCGIFADDR
#if HAVE_SIOCGLIFNUM
  if (ioctl (sock, SIOCGLIFADDR, &ifr) == 0) {
#else
  if (ioctl (sock, SIOCGIFADDR, &ifr) == 0) {
#endif
    found = TRUE;

    if (string_from_sockaddr ((struct sockaddr *)&ifr.CNAME(ifr_addr), buffer, sizeof (buffer)) == 0) {
      size_t addr_len = strlen(buffer);
      char* addr_str = (char *) malloc (addr_len * sizeof(char));
      memcpy(addr_str, buffer, addr_len);
      address->addr = addr_str;
    }
  }
#endif

#if HAVE_SIOCGIFNETMASK
#if HAVE_SIOCGLIFNUM
  if (ioctl (sock, SIOCGLIFNETMASK, &ifr) == 0) {
#else
  if (ioctl (sock, SIOCGIFNETMASK, &ifr) == 0) {
#endif
    found = TRUE;

    if (string_from_sockaddr ((struct sockaddr *)&ifr.CNAME(ifr_addr), buffer, sizeof (buffer)) == 0) {
      size_t netmask_len = strlen(buffer);
      char* netmask_str = (char *) malloc (netmask_len * sizeof(char));
      memcpy(netmask_str, buffer, netmask_len);
      address->netmask = netmask_str;
    }
  }
#endif

#if HAVE_SIOCGIFFLAGS
#if HAVE_SIOCGLIFNUM
  if (ioctl (sock, SIOCGLIFFLAGS, &ifr) == 0) {
#else
  if (ioctl (sock, SIOCGIFFLAGS, &ifr) == 0) {
#endif
    found = TRUE;

    if (ifr.CNAME(ifr_flags) & IFF_POINTOPOINT)
      is_p2p = TRUE;
  }
#endif

#if HAVE_SIOCGIFBRDADDR
#if HAVE_SIOCGLIFNUM
  if (!is_p2p && ioctl (sock, SIOCGLIFBRDADDR, &ifr) == 0) {
#else
  if (!is_p2p && ioctl (sock, SIOCGIFBRDADDR, &ifr) == 0) {
#endif
    found = TRUE;

    if (string_from_sockaddr ((struct sockaddr *)&ifr.CNAME(ifr_addr), buffer, sizeof (buffer)) == 0) {
      size_t braddr_len = strlen(buffer);
      char* braddr_str = (char *) malloc (braddr_len * sizeof(char));
      memcpy(braddr_str, buffer, braddr_len);
      address->broadcast = braddr_str;
    }
  }
#endif

#if HAVE_SIOCGIFDSTADDR
#if HAVE_SIOCGLIFNUM
  if (is_p2p && ioctl (sock, SIOCGLIFBRDADDR, &ifr) == 0) {
#else
  if (is_p2p && ioctl (sock, SIOCGIFBRDADDR, &ifr) == 0) {
#endif
    found = TRUE;

    if (string_from_sockaddr ((struct sockaddr *)&ifr.CNAME(ifr_addr), buffer, sizeof (buffer)) == 0) {
      size_t dstaddr_len = strlen(buffer);
      char* dstaddr_str = (char *) malloc (dstaddr_len * sizeof(char));
      memcpy(dstaddr_str, buffer, dstaddr_len);
      address->peer = braddr_str;
    }
  }
#endif

  address->family = AF_INET;

  append_address(results, address);
      
  close (sock);
#endif /* HAVE_SOCKET_IOCTLS */

  if (found)
    return 0;
  else {
    //You must specify a valid interface name
    return EINVAL;
  }
}

/* -- interfaces() ---------------------------------------------------------- */

      
static bool
append_interface(Interface** results, Interface* interface) {
    Interface* iface = *results;
    if (interface == NULL) {
        return false;
    }
    if (*results == NULL) {
        *results = interface;
        return true;
    }
    while(iface->next != NULL) {
        // as we have the information already, we do not need to add it again.
        if (iface == interface || (iface->name != NULL &&  interface->name != NULL && strcmp(iface->name, interface->name)) == 0) {
            return true;
        }

        iface = iface->next;
    }
    iface->next = interface;
    return true;
}

static void
release_interfaces(Interface** results) {
    Interface *head = *results, *tail = NULL;
    if (*results == NULL) {
        return;
    }
    // traseverse linked list
    while(head != NULL) {
        tail = head;
        head = head->next;
        
        if (tail->name != NULL) {
            free(tail->name);
        }

        if (tail->address != NULL) {
            release_addresses(&(tail->address));
        }

        tail->next = NULL;
        free(tail);
        tail = NULL;
    }
    *results = NULL;
}
      
static int
interfaces (Interface **results)
{
#if HAVE_GETIFADDRS
  /* .. UNIX, with getifaddrs() ............................................. */

  const char *prev_name = NULL;
  struct ifaddrs *addrs = NULL;
  struct ifaddrs *addr = NULL;

  if (getifaddrs (&addrs) < 0) {
    return ENETUNREACH;
  }

  for (addr = addrs; addr; addr = addr->ifa_next) {
    if (!prev_name || strncmp (addr->ifa_name, prev_name, IFNAMSIZ) != 0) {
      
      size_t ifname_len = strlen(addr->ifa_name);
      char* ifname_str = (char *) malloc (ifname_len * sizeof(char));
      memcpy(ifname_str, addr->ifa_name, ifname_len);
    
      Interface *interface = (Interface *)calloc(1, sizeof (Interface));
      interface->name = ifname_str;
      append_interface(results, interface);
        
      prev_name = addr->ifa_name;
    }
  }

  freeifaddrs (addrs);
#elif HAVE_SIOCGIFCONF
  /* .. UNIX, with SIOC ioctl()s ............................................ */

  const char *prev_name = NULL;
  int fd = socket (AF_INET, SOCK_DGRAM, 0);
  struct CNAME(ifconf) ifc;
  int len = -1;

  if (fd < 0) {
    return ENOMEM;
  }

  // Try to find out how much space we need
#if HAVE_SIOCGSIZIFCONF
  if (ioctl (fd, SIOCGSIZIFCONF, &len) < 0)
    len = -1;
#elif HAVE_SIOCGLIFNUM
  { struct lifnum lifn;
    lifn.lifn_family = AF_UNSPEC;
    lifn.lifn_flags = LIFC_NOXMIT | LIFC_TEMPORARY | LIFC_ALLZONES;
    ifc.lifc_family = AF_UNSPEC;
    ifc.lifc_flags = LIFC_NOXMIT | LIFC_TEMPORARY | LIFC_ALLZONES;
    if (ioctl (fd, SIOCGLIFNUM, (char *)&lifn) < 0)
      len = -1;
    else
      len = lifn.lifn_count;
  }
#endif

  // As a last resort, guess
  if (len < 0)
    len = 64;

  ifc.CNAME(ifc_len) = (int)(len * sizeof (struct CNAME(ifreq)));
  ifc.CNAME(ifc_buf) = malloc (ifc.CNAME(ifc_len));

  if (!ifc.CNAME(ifc_buf)) {
    // Not enough memory
    close (fd);
    return ENOMEM;
  }

#if HAVE_SIOCGLIFNUM
  if (ioctl (fd, SIOCGLIFCONF, &ifc) < 0) {
#else
  if (ioctl (fd, SIOCGIFCONF, &ifc) < 0) {
#endif
    free (ifc.CNAME(ifc_req));
    close (fd);
    return ENETUNREACH;
  }

  struct CNAME(ifreq) *pfreq = ifc.CNAME(ifc_req);
  struct CNAME(ifreq) *pfreqend = (struct CNAME(ifreq) *)((char *)pfreq
                                                          + ifc.CNAME(ifc_len));
  while (pfreq < pfreqend) {
    if (!prev_name || strncmp (prev_name, pfreq->CNAME(ifr_name), IFNAMSIZ) != 0) {

      size_t ifname_len = strlen(pfreq->CNAME(ifr_name));
      char* ifname_str = (char *) malloc (ifname_len * sizeof(char));
      memcpy(ifname_str, pfreq->CNAME(ifr_name), ifname_len);
      
      Interface *interface = (Interface *)calloc(1, sizeof (Interface));
      interface->name = ifname_str;
      append_interface(results, interface);
        
      prev_name = pfreq->CNAME(ifr_name);
    }

#if !HAVE_SOCKADDR_SA_LEN
    ++pfreq;
#else
    /* On some platforms, the ifreq struct can *grow*(!) if the socket address
       is very long.  Mac OS X is such a platform. */
    {
      size_t len = sizeof (struct CNAME(ifreq));
      if (pfreq->ifr_addr.sa_len > sizeof (struct sockaddr))
        len = len - sizeof (struct sockaddr) + pfreq->ifr_addr.sa_len;
        pfreq = (struct CNAME(ifreq) *)((char *)pfreq + len);
    }
#endif
  }

  free (ifc.CNAME(ifc_buf));
  close (fd);
#endif /* HAVE_SIOCGIFCONF */

  return 0;
}

int
find_system_interfaces(Interface **results) {
    Interface *iface = NULL;
    int err = interfaces(results);
    if (err != 0) {
        return err;
    }
    
    iface = *results;
    while(iface != NULL) {
        err = ifaddrs(&(iface->address), iface->name);
        if (err != 0) {
            return err;
        }
        iface = iface->next;
    }
    return 0;
}

void
release_interfaces_info(Interface **results) {
    Interface *iface = *results;
    while(iface != NULL) {
        release_addresses(&(iface->address));
        iface = iface->next;
    }
    release_interfaces(results);
}

/* -- gateways() ------------------------------------------------------------ */

static bool
append_gatway(Gateway** results, Gateway* gateway) {
    Gateway* node = *results;
    if (gateway == NULL) {
        return false;
    }
    if (*results == NULL) {
        *results = gateway;
        return true;
    }
    while(node->next != NULL) {
        node = node->next;
    }
    node->next = gateway;
    return true;
}
    
void
release_gateways_info(Gateway** results) {
    Gateway *head = *results, *tail = NULL;
    if (*results == NULL) {
        return;
    }
    // traseverse linked list
    while(head != NULL) {
        tail = head;
        head = head->next;
        
        if (tail->ifname != NULL) {
            free(tail->ifname);
        }
        if (tail->addr != NULL) {
            free(tail->addr);
        }
        tail->next = NULL;
        free(tail);
        tail = NULL;
    }
    *results = NULL;
}

static Gateway*
find_default_gateway_by_family(Gateway** results, unsigned char family) {
    Gateway *node = *results;
    if (*results == NULL) {
        return NULL;
    }
    while(node != NULL) {
        if (node->family == family && node->is_default) {
            return node;
        }
        node = node->next;
    }
    return NULL;
}
    
Gateway*
find_default_ip4_gw(Gateway** results) {
    return find_default_gateway_by_family(results, AF_INET);
}

Gateway*
find_default_ip6_gw(Gateway** results) {
    return find_default_gateway_by_family(results, AF_INET6);
}


    
static Gateway*
find_first_gateway_by_family_and_interface(Gateway** results, unsigned char family, char *interface) {
    Gateway *node = *results;
    if (*results == NULL) {
        return NULL;
    }
    while(node != NULL) {
        if (node->family == family && strcmp(node->ifname, interface) == 0) {
            return node;
        }
        node = node->next;
    }
    return NULL;
}

// find the ip4 gateway for an interface
Gateway*
find_ip4_gw_for_interface(Gateway** results, char *interface) {
    return find_first_gateway_by_family_and_interface(results, AF_INET, interface);
}

// find the ip4 gateway for an interface
Gateway*
find_ip6_gw_for_interface(Gateway** results, char *interface) {
    return find_first_gateway_by_family_and_interface(results, AF_INET6, interface);
}
    
int
find_system_gateways(Gateway** results)
{
#if defined(HAVE_PF_NETLINK)
  /* .. Linux (PF_NETLINK socket) ........................................... */

  /* PF_NETLINK is pretty poorly documented and it looks to be quite easy to
     get wrong.  This *appears* to be the right way to do it, even though a
     lot of the code out there on the 'Net is very different! */

  struct routing_msg {
    struct nlmsghdr hdr;
    struct rtmsg    rt;
    char            data[0];
  } *pmsg, *msgbuf;
  int s;
  int seq = 0;
  ssize_t ret;
  struct sockaddr_nl sanl;
  static const struct sockaddr_nl sanl_kernel = { .nl_family = AF_NETLINK };
  socklen_t sanl_len;
  int pagesize = getpagesize();
  int bufsize = pagesize < 8192 ? pagesize : 8192;
  int is_multi = 0;
  int interrupted = 0;
  int def_priorities[RTNL_FAMILY_MAX];

  memset(def_priorities, 0xff, sizeof(def_priorities));

  msgbuf = (struct routing_msg *)malloc (bufsize);
    
  if (!msgbuf) {
    return ENOMEM;
  }

  s = socket (PF_NETLINK, SOCK_RAW, NETLINK_ROUTE);

  if (s < 0) {
    free (msgbuf);
    return ENETUNREACH;
  }

  sanl.nl_family = AF_NETLINK;
  sanl.nl_groups = 0;
  sanl.nl_pid = 0;

  if (bind (s, (struct sockaddr *)&sanl, sizeof (sanl)) < 0) {
    free (msgbuf);
    close (s);
    return ENONET;
  }

  sanl_len = sizeof (sanl);
    
  if (getsockname (s, (struct sockaddr *)&sanl, &sanl_len) < 0) {
    free (msgbuf);
    close (s);
    return ENETUNREACH;
  }

  do {
    interrupted = 0;

    pmsg = msgbuf;
    memset (pmsg, 0, sizeof (struct routing_msg));
    pmsg->hdr.nlmsg_len = NLMSG_LENGTH(sizeof(struct rtmsg));
    pmsg->hdr.nlmsg_flags = NLM_F_DUMP | NLM_F_REQUEST;
    pmsg->hdr.nlmsg_seq = ++seq;
    pmsg->hdr.nlmsg_type = RTM_GETROUTE;
    pmsg->hdr.nlmsg_pid = 0;

    pmsg->rt.rtm_family = 0;

    if (sendto (s, pmsg, pmsg->hdr.nlmsg_len, 0,
                (struct sockaddr *)&sanl_kernel, sizeof(sanl_kernel)) < 0) {
      free (msgbuf);
      close (s);
      return ENETUNREACH;
    }

    do {
      struct sockaddr_nl sanl_from;
      struct iovec iov = { msgbuf, bufsize };
      struct msghdr msghdr = {
        &sanl_from,
        sizeof(sanl_from),
        &iov,
        1,
        NULL,
        0,
        0
      };
      int nllen;

      ret = recvmsg (s, &msghdr, 0);

      if (msghdr.msg_flags & MSG_TRUNC) {
        //"netlink message truncated"
        free (msgbuf);
        close (s);
        return ENETUNREACH;
      }

      if (ret < 0) {
        free (msgbuf);
        close (s);
        return ENETUNREACH;
      }

      nllen = ret;
      pmsg = msgbuf;
      while (NLMSG_OK (&pmsg->hdr, nllen)) {
        void *dst = NULL;
        void *gw = NULL;
        int ifndx = -1;
        struct rtattr *attrs, *attr;
        int len;
        int priority;

        /* Ignore messages not for us */
        if (pmsg->hdr.nlmsg_seq != seq || pmsg->hdr.nlmsg_pid != sanl.nl_pid)
          goto next;

        /* This is only defined on Linux kernel versions 3.1 and higher */
#ifdef NLM_F_DUMP_INTR
        if (pmsg->hdr.nlmsg_flags & NLM_F_DUMP_INTR) {
          /* The dump was interrupted by a signal; we need to go round again */
          interrupted = 1;
          is_multi = 0;
          break;
        }
#endif

        is_multi = pmsg->hdr.nlmsg_flags & NLM_F_MULTI;

        if (pmsg->hdr.nlmsg_type == NLMSG_DONE) {
          is_multi = interrupted = 0;
          break;
        }

        if (pmsg->hdr.nlmsg_type == NLMSG_ERROR) {
          struct nlmsgerr *perr = (struct nlmsgerr *)&pmsg->rt;
          errno = -perr->error;
          free (msgbuf);
          close (s);
          return errno;
        }

        attr = attrs = RTM_RTA(&pmsg->rt);
        len = RTM_PAYLOAD(&pmsg->hdr);
        priority = -1;
        while (RTA_OK(attr, len)) {
          switch (attr->rta_type) {
          case RTA_GATEWAY:
            gw = RTA_DATA(attr);
            break;
          case RTA_DST:
            dst = RTA_DATA(attr);
            break;
          case RTA_OIF:
            ifndx = *(int *)RTA_DATA(attr);
            break;
          case RTA_PRIORITY:
            priority = *(int *)RTA_DATA(attr);
            break;
          default:
            break;
          }

          attr = RTA_NEXT(attr, len);
        }

        /* We're looking for gateways with no destination */
        if (!dst && gw && ifndx >= 0) {
          char buffer[256];
          char ifnamebuf[IF_NAMESIZE];
          char *ifname;
          const char *addr;
          Gateway* gateway = NULL;
          bool is_default = false;
            
          ifname = if_indextoname (ifndx, ifnamebuf);

          if (!ifname)
            goto next;

          addr = inet_ntop (pmsg->rt.rtm_family, gw, buffer, sizeof (buffer));

          if (!addr)
            goto next;

          /* We set isdefault to True if this route came from the main table;
             this should correspond with the way most people set up alternate
             routing tables on Linux. */

          is_default = (bool)(pmsg->rt.rtm_table == RT_TABLE_MAIN);
            
          /* Try to pick the active default route based on priority (which
             is displayed in the UI as "metric", confusingly) */
          if (pmsg->rt.rtm_family < RTNL_FAMILY_MAX) {
            if (def_priorities[pmsg->rt.rtm_family] == -1)
              def_priorities[pmsg->rt.rtm_family] = priority;
            else {
              if (priority == -1
                  || priority > def_priorities[pmsg->rt.rtm_family])
                is_default = false;
            }
          }

          gateway = (Gateway *) calloc(1, sizeof (Gateway));
          
          gateway->is_default = is_default;
          gateway->family = pmsg->rt.rtm_family;
          
          size_t gw_addr_len = strlen(buffer);
          char* gw_addr = (char *) malloc (gw_addr_len * sizeof(char));
          memcpy(gw_addr, buffer, gw_addr_len);
          gateway->addr = gw_addr;
          
          size_t gw_ifname_len = strlen(ifname);
          char* gw_ifname = (char *) malloc (gw_ifname_len * sizeof(char));
          memcpy(gw_ifname, ifname, gw_ifname_len);
          gateway->ifname = gw_ifname;

          append_gatway(results, gateway);
        }

      next:
	pmsg = (struct routing_msg *)NLMSG_NEXT(&pmsg->hdr, nllen);
      }
    } while (is_multi);
  } while (interrupted);

  free (msgbuf);
  close (s);
#elif defined(HAVE_SYSCTL_CTL_NET)
  /* .. UNIX, via sysctl() .................................................. */

  int mib[] = { CTL_NET, PF_ROUTE, 0, 0, NET_RT_FLAGS,
                RTF_UP | RTF_GATEWAY };
  size_t len;
  char *buffer = NULL, *ptr, *end;
  int ret;
  char ifnamebuf[IF_NAMESIZE];
  char *ifname;

  /* Remembering that the routing table may change while we're reading it,
     we need to do this in a loop until we succeed. */
  do {
    if (sysctl (mib, 6, 0, &len, 0, 0) < 0) {
      free (buffer);
      return ENETUNREACH;
    }

    ptr = realloc(buffer, len);
    if (!ptr) {
      free (buffer);
      return ENOMEM;
    }

    buffer = ptr;

    ret = sysctl (mib, 6, buffer, &len, 0, 0);
  } while (ret != 0 || errno == ENOMEM || errno == EINTR);

  if (ret < 0) {
    free (buffer);
    return ENETUNREACH;
  }

  ptr = buffer;
  end = buffer + len;

  while (ptr + sizeof (struct rt_msghdr) <= end) {
    struct rt_msghdr *msg = (struct rt_msghdr *)ptr;
    char *msgend = (char *)msg + msg->rtm_msglen;
    int addrs = msg->rtm_addrs;
    int addr = RTA_DST;

    if (msgend > end)
      break;

    ifname = if_indextoname (msg->rtm_index, ifnamebuf);

    if (!ifname) {
      ptr = msgend;
      continue;
    }

    ptr = (char *)(msg + 1);
    while (ptr + sizeof (struct sockaddr) <= msgend && addrs) {
      struct sockaddr *sa = (struct sockaddr *)ptr;
      int len = SA_LEN(sa);

      if (!len)
        len = 4;
      else
        len = (len + 3) & ~3;

      if (ptr + len > msgend)
        break;

      while (!(addrs & addr))
        addr <<= 1;

      addrs &= ~addr;

      if (addr == RTA_DST) {
        if (sa->sa_family == AF_INET) {
          struct sockaddr_in *sin = (struct sockaddr_in *)sa;
          if (sin->sin_addr.s_addr != INADDR_ANY)
            break;
#ifdef AF_INET6
        } else if (sa->sa_family == AF_INET6) {
          struct sockaddr_in6 *sin6 = (struct sockaddr_in6 *)sa;
          if (memcmp (&sin6->sin6_addr, &in6addr_any, sizeof (in6addr_any)) != 0)
            break;
#endif
        } else {
          break;
        }
      }

      if (addr == RTA_GATEWAY) {
        char buffer[256];
        Gateway* gateway = NULL;

        if (string_from_sockaddr (sa, buffer, sizeof(buffer)) == 0) {
#ifdef RTF_IFSCOPE
          bool is_default = !(msg->rtm_flags & RTF_IFSCOPE);
#else
          bool is_default = true;
#endif
          gateway = (Gateway *) calloc(1, sizeof (Gateway));
          
          gateway->is_default = is_default;
          gateway->family = sa->sa_family;
            
          size_t gw_addr_len = strlen(buffer);
          char* gw_addr = (char *) malloc (gw_addr_len * sizeof(char));
          memcpy(gw_addr, buffer, gw_addr_len);
          gateway->addr = gw_addr;
          
          size_t gw_ifname_len = strlen(ifname);
          char* gw_ifname = (char *) malloc (gw_ifname_len * sizeof(char));
          memcpy(gw_ifname, ifname, gw_ifname_len);
          gateway->ifname = gw_ifname;
        }
          
        if (gateway != NULL) {
          append_gatway(results, gateway);
        }
      }

      /* These are aligned on a 4-byte boundary */
      ptr += len;
    }

    ptr = msgend;
  }

  free (buffer);
#elif defined(HAVE_PF_ROUTE)
  /* .. UNIX, via PF_ROUTE socket ........................................... */

  /* The PF_ROUTE code will only retrieve gateway information for AF_INET and
     AF_INET6.  This is because it would need to loop through all possible
     values, and the messages it needs to send in each case are potentially
     different.  It is also very likely to return a maximum of one gateway
     in each case (since we can't read the entire routing table this way, we
     can only ask about routes). */

  int pagesize = getpagesize();
  int bufsize = pagesize < 8192 ? 8192 : pagesize;
  struct rt_msghdr *pmsg;
  int s;
  int seq = 0;
  int pid = getpid();
  ssize_t ret;
  struct sockaddr_in *sin_dst, *sin_netmask;
  struct sockaddr_dl *sdl_ifp;
  struct sockaddr_in6 *sin6_dst;
  size_t msglen;
  char ifnamebuf[IF_NAMESIZE];
  char *ifname;
  int skip;

  pmsg = (struct rt_msghdr *)malloc (bufsize);
    
  if (!pmsg) {
    return ENOMEM;
  }

  s = socket (PF_ROUTE, SOCK_RAW, 0);

  if (s < 0) {
    free (pmsg);
    return ENETUNREACH;
  }

  msglen = (sizeof (struct rt_msghdr)
            + 2 * sizeof (struct sockaddr_in) 
            + sizeof (struct sockaddr_dl));
  memset (pmsg, 0, msglen);
  
  /* AF_INET first */
  pmsg->rtm_msglen = msglen;
  pmsg->rtm_type = RTM_GET;
  pmsg->rtm_index = 0;
  pmsg->rtm_flags = RTF_UP | RTF_GATEWAY;
  pmsg->rtm_version = RTM_VERSION;
  pmsg->rtm_seq = ++seq;
  pmsg->rtm_pid = 0;
  pmsg->rtm_addrs = RTA_DST | RTA_NETMASK | RTA_IFP;

  sin_dst = (struct sockaddr_in *)(pmsg + 1);
  sin_netmask = (struct sockaddr_in *)(sin_dst + 1);
  sdl_ifp = (struct sockaddr_dl *)(sin_netmask + 1);

  sin_dst->sin_family = AF_INET;
  sin_netmask->sin_family = AF_INET;
  sdl_ifp->sdl_family = AF_LINK;

#if HAVE_SOCKADDR_SA_LEN
  sin_dst->sin_len = sizeof (struct sockaddr_in);
  sin_netmask->sin_len = sizeof (struct sockaddr_in);
  sdl_ifp->sdl_len = sizeof (struct sockaddr_dl);
#endif

  skip = 0;
  if (send (s, pmsg, msglen, 0) < 0) {
    if (errno == ESRCH)
      skip = 1;
    else {
      close (s);
      free (pmsg);
      return ENETUNREACH;
    }
  }

  while (!skip && !(pmsg->rtm_flags & RTF_DONE)) {
    char *ptr;
    char *msgend;
    int addrs;
    int addr;
    struct sockaddr_in *dst = NULL;
    struct sockaddr_in *gw = NULL;
    struct sockaddr_dl *ifp = NULL;
    Gateway* gateway = NULL;

    do {
      ret = recv (s, pmsg, bufsize, 0);
    } while ((ret < 0 && errno == EINTR)
             || (ret > 0 && (pmsg->rtm_seq != seq || pmsg->rtm_pid != pid)));

    if (ret < 0) {
      close (s);
      free (pmsg);
      return ENETUNREACH;
    }

    if (pmsg->rtm_errno != 0) {
      if (pmsg->rtm_errno == ESRCH)
        skip = 1;
      else {
        errno = pmsg->rtm_errno;
        int err = pmsg->rtm_errno;
        close (s);
        free (pmsg);
        return err;
      }
    }

    if (skip)
      break;

    ptr = (char *)(pmsg + 1);
    msgend = (char *)pmsg + pmsg->rtm_msglen;
    addrs = pmsg->rtm_addrs;
    addr = RTA_DST;
    while (ptr + sizeof (struct sockaddr) <= msgend && addrs) {
      struct sockaddr *sa = (struct sockaddr *)ptr;
      int len = SA_LEN(sa);

      if (!len)
        len = 4;
      else
        len = (len + 3) & ~3;

      if (ptr + len > msgend)
        break;

      while (!(addrs & addr))
        addr <<= 1;

      addrs &= ~addr;

      switch (addr) {
      case RTA_DST:
        dst = (struct sockaddr_in *)sa;
        break;
      case RTA_GATEWAY:
        gw = (struct sockaddr_in *)sa;
        break;
      case RTA_IFP:
        ifp = (struct sockaddr_dl *)sa;
        break;
      }

      ptr += len;
    }

    if ((dst && dst->sin_family != AF_INET)
        || (gw && gw->sin_family != AF_INET)
        || (ifp && ifp->sdl_family != AF_LINK)) {
      dst = gw = NULL;
      ifp = NULL;
    }

    if (dst && dst->sin_addr.s_addr == INADDR_ANY)
        dst = NULL;

    if (!dst && gw && ifp) {
      char buffer[256];

      if (ifp->sdl_index)
        ifname = if_indextoname (ifp->sdl_index, ifnamebuf);
      else {
        memcpy (ifnamebuf, ifp->sdl_data, ifp->sdl_nlen);
        ifnamebuf[ifp->sdl_nlen] = '\0';
        ifname = ifnamebuf;
      }

      if (string_from_sockaddr ((struct sockaddr *)gw,
                                buffer, sizeof(buffer)) == 0) {
#ifdef RTF_IFSCOPE
        bool is_default = !(pmsg->rtm_flags & RTF_IFSCOPE);
#else
        bool is_default = true;
#endif

        gateway = (Gateway *) calloc(1, sizeof (Gateway));
        
        gateway->is_default = is_default;
        gateway->family = AF_INET;
        
        size_t gw_addr_len = strlen(buffer);
        char* gw_addr = (char *) malloc (gw_addr_len * sizeof(char));
        memcpy(gw_addr, buffer, gw_addr_len);
        gateway->addr = gw_addr;
        
        size_t gw_ifname_len = strlen(ifname);
        char* gw_ifname = (char *) malloc (gw_ifname_len * sizeof(char));
        memcpy(gw_ifname, ifname, gw_ifname_len);
        gateway->ifname = gw_ifname;
      }
        
      if (gateway != NULL) {
          append_gatway(results, gateway);
      }
    }
  }

  /* The code below is very similar to, but not identical to, the code above.
     We could probably refactor some of it, but take care---there are subtle
     differences! */

#ifdef AF_INET6
  /* AF_INET6 now */
  msglen = (sizeof (struct rt_msghdr)
            + sizeof (struct sockaddr_in6)
            + sizeof (struct sockaddr_dl));
  memset (pmsg, 0, msglen);

  pmsg->rtm_msglen = msglen;
  pmsg->rtm_type = RTM_GET;
  pmsg->rtm_index = 0;
  pmsg->rtm_flags = RTF_UP | RTF_GATEWAY;
  pmsg->rtm_version = RTM_VERSION;
  pmsg->rtm_seq = ++seq;
  pmsg->rtm_pid = 0;
  pmsg->rtm_addrs = RTA_DST | RTA_IFP;

  sin6_dst = (struct sockaddr_in6 *)(pmsg + 1);
  sdl_ifp = (struct sockaddr_dl *)(sin6_dst + 1);

  sin6_dst->sin6_family = AF_INET6;
  sin6_dst->sin6_addr = in6addr_any;
  sdl_ifp->sdl_family = AF_LINK;

#if HAVE_SOCKADDR_SA_LEN
  sin6_dst->sin6_len = sizeof (struct sockaddr_in6);
  sdl_ifp->sdl_len = sizeof (struct sockaddr_dl);
#endif

  skip = 0;
  if (send (s, pmsg, msglen, 0) < 0) {
    if (errno == ESRCH)
      skip = 1;
    else {
      close (s);
      free (pmsg);
      return ENETUNREACH;
    }
  }

  while (!skip && !(pmsg->rtm_flags & RTF_DONE)) {
    char *ptr;
    char *msgend;
    int addrs;
    int addr;
    struct sockaddr_in6 *dst = NULL;
    struct sockaddr_in6 *gw = NULL;
    struct sockaddr_dl *ifp = NULL;
    Gateway* gateway = NULL;

    do {
      ret = recv (s, pmsg, bufsize, 0);
    } while ((ret < 0 && errno == EINTR)
             || (ret > 0 && (pmsg->rtm_seq != seq || pmsg->rtm_pid != pid)));

    if (ret < 0) {
      close (s);
      free (pmsg);
      return ENETUNREACH;
    }

    if (pmsg->rtm_errno != 0) {
      if (pmsg->rtm_errno == ESRCH)
        skip = 1;
      else {
        errno = pmsg->rtm_errno;
        int err = pmsg->rtm_errno;
        close (s);
        free (pmsg);
        return err;
      }
    }

    if (skip)
      break;

    ptr = (char *)(pmsg + 1);
    msgend = (char *)pmsg + pmsg->rtm_msglen;
    addrs = pmsg->rtm_addrs;
    addr = RTA_DST;
    while (ptr + sizeof (struct sockaddr) <= msgend && addrs) {
      struct sockaddr *sa = (struct sockaddr *)ptr;
      int len = SA_LEN(sa);

      if (!len)
        len = 4;
      else
        len = (len + 3) & ~3;

      if (ptr + len > msgend)
        break;

      while (!(addrs & addr))
        addr <<= 1;

      addrs &= ~addr;

      switch (addr) {
      case RTA_DST:
        dst = (struct sockaddr_in6 *)sa;
        break;
      case RTA_GATEWAY:
        gw = (struct sockaddr_in6 *)sa;
        break;
      case RTA_IFP:
        ifp = (struct sockaddr_dl *)sa;
        break;
      }

      ptr += len;
    }

    if ((dst && dst->sin6_family != AF_INET6)
        || (gw && gw->sin6_family != AF_INET6)
        || (ifp && ifp->sdl_family != AF_LINK)) {
      dst = gw = NULL;
      ifp = NULL;
    }

    if (dst && memcmp (&dst->sin6_addr, &in6addr_any,
                       sizeof(struct in6_addr)) == 0)
        dst = NULL;

    if (!dst && gw && ifp) {
      char buffer[256];

      if (ifp->sdl_index)
        ifname = if_indextoname (ifp->sdl_index, ifnamebuf);
      else {
        memcpy (ifnamebuf, ifp->sdl_data, ifp->sdl_nlen);
        ifnamebuf[ifp->sdl_nlen] = '\0';
        ifname = ifnamebuf;
      }

      if (string_from_sockaddr ((struct sockaddr *)gw,
                                buffer, sizeof(buffer)) == 0) {
#ifdef RTF_IFSCOPE
        bool is_default = !(pmsg->rtm_flags & RTF_IFSCOPE);
#else
        bool is_default = true;
#endif
        gateway = (Gateway *) calloc(1, sizeof (Gateway));
        
        gateway->is_default = is_default;
        gateway->family = AF_INET6;
        
        size_t gw_addr_len = strlen(buffer);
        char* gw_addr = (char *) malloc (gw_addr_len * sizeof(char));
        memcpy(gw_addr, buffer, gw_addr_len);
        gateway->addr = gw_addr;
        
        size_t gw_ifname_len = strlen(ifname);
        char* gw_ifname = (char *) malloc (gw_ifname_len * sizeof(char));
        memcpy(gw_ifname, ifname, gw_ifname_len);
        gateway->ifname = gw_ifname;
      }
        
      if (gateway != NULL) {
          append_gatway(results, gateway);
      }
    }
  }
#endif /* AF_INET6 */

  free (pmsg);
#else
  /* If we don't know how to implement this on your platform, we raise an
     exception. */
#  error Unable to obtain gateway information on your platform
#endif

  return 0;
}

