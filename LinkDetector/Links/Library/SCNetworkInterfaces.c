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

#include "SCNetworkInterfaces.h"

// in general, it is recommended to copy struct data to CF- collection, but in
// in order to reduce memory ops, we'll not going to copy struct.
#  define SHOULD_COPY_STRUCT_DATA 0

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

#  include <libkern/OSAtomic.h>
#  include <net/if_media.h>
#  include <unistd.h>
#  include <stddef.h>
#  include <stdio.h>
#  include <string.h>
#  include <stdlib.h>

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

#pragma mark - Utility Functions

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


#pragma mark - MEDIA STATUS
CF_INLINE Boolean _SCNACStringIsStringStartWithString(const char *cStr1, const char *cStr2) {
    while(*cStr1 && *cStr2) {if(*cStr1++ != *cStr2++) { return false; }}; return true;
}

bool
SCNetworkInterfaceMediaStatus(SCNetworkInterfaceRef interface) {
    char devname[256];
    CFStringRef bsdName = SCNetworkInterfaceGetBSDName(interface);
    if (!CFStringGetCString(bsdName, devname, 256, kCFStringEncodingUTF8)) {
        return false;
    }
    struct ifmediareq ifm;
    memset(&ifm, 0, sizeof(struct ifmediareq));
    strncpy(ifm.ifm_name, devname, IFNAMSIZ);
    int s = socket(AF_INET, SOCK_DGRAM, 0);
    ioctl(s, SIOCGIFMEDIA, &ifm);
    bool status = false;
    
    switch (IFM_TYPE(ifm.ifm_active)) {
        case IFM_FDDI:
        case IFM_TOKEN:
            if (_SCNACStringIsStringStartWithString(devname, "fw")) {
                status = (ifm.ifm_status & IFM_ACTIVE) ? true : false;
            } else {
                status = (ifm.ifm_status & IFM_ACTIVE) ? true : false;
            }
            break;
        case IFM_IEEE80211:
            status = (ifm.ifm_status & IFM_ACTIVE) ? true : false;
            break;
        default:
            if (_SCNACStringIsStringStartWithString(devname, "en")) {
                status = (ifm.ifm_status & IFM_ACTIVE) ? true : false;
            } else {
                status = (ifm.ifm_status & IFM_ACTIVE) ? true : false;
            }
    }

    return status;
}

#pragma mark - INTERFACE
static void
_SCNIAddressRelease(SCNIAddress *address) {
    if (address != NULL) {
#if SHOULD_COPY_STRUCT_DATA
        if (address->addr != NULL) {
            free(address->addr);
            address->addr = NULL;
        }
        if (address->netmask != NULL) {
            free(address->netmask);
            address->netmask = NULL;
        }
        if (address->broadcast != NULL) {
            free(address->broadcast);
            address->broadcast = NULL;
        }
        if (address->peer != NULL) {
            free(address->peer);
            address->peer = NULL;
        }
#endif
        free(address);
    }
}

/* -- ifaddresses() --------------------------------------------------------- */
/*!
	@function ifaddresses
	@discussion Returns ipv4/ipv6 addresses (dotted format) linked to the interface.
	@param results List where address should go
            ifaName The network interface name.
	@result The list of ipv4 addresses linked the interface;
            NULL if no ipv4 addresses are supported or linked.
 */

static errno_t
ifaddresses (CFMutableArrayRef results, CFStringRef ifaName)
{
  int found = FALSE;
  char ifname[256];
#if HAVE_GETIFADDRS
  struct ifaddrs *addrs = NULL;
  struct ifaddrs *addr = NULL;
#endif
    
  if (!CFStringGetCString(ifaName, ifname, 256, kCFStringEncodingUTF8) || strlen(ifname) == 0) {
      return EINVAL;
  }

#if HAVE_GETIFADDRS
  /* .. UNIX, with getifaddrs() ............................................. */

  if (getifaddrs (&addrs) < 0) {
    return ENETUNREACH;
  }

  for (addr = addrs; addr; addr = addr->ifa_next) {
    char buffer[256];
    SCNIAddress *address = NULL;

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
    address = (SCNIAddress *) calloc(1, sizeof (SCNIAddress));
    
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
    
    CFArrayAppendValue(results, address);
    _SCNIAddressRelease(address);
  }

  freeifaddrs (addrs);
#endif /* HAVE_GETIFADDRS */

  if (found)
    return 0;
  else {
    //You must specify a valid interface name
    return EINVAL;
  }
}

// CFMutableArray Callbacks for SCNIAddress
const void*
_address_retain(CFAllocatorRef allocator, const void *ptr) {
    SCNIAddress *add_src = (SCNIAddress *)ptr;
    SCNIAddress *address = (SCNIAddress *)CFAllocatorAllocate(allocator, sizeof(SCNIAddress), 0);
    
#if SHOULD_COPY_STRUCT_DATA
    address->flags = add_src->flags;
    address->family = add_src->family;
    
    if (add_src->addr != NULL && strlen(add_src->addr) != 0) {
        size_t str_len = strlen(add_src->addr);
        char *str = (char *) malloc(sizeof(char) * str_len);
        memcpy(str, add_src->addr, str_len);
        address->addr = str;
    } else {
        address->addr = NULL;
    }
    if (add_src->netmask != NULL && strlen(add_src->netmask) != 0) {
        size_t str_len = strlen(add_src->netmask);
        char *str = (char *) malloc(sizeof(char) * str_len);
        memcpy(str, add_src->netmask, str_len);
        address->netmask = str;
    } else {
        address->netmask = NULL;
    }
    if (add_src->broadcast != NULL && strlen(add_src->broadcast) != 0) {
        size_t str_len = strlen(add_src->broadcast);
        char *str = (char *) malloc(sizeof(char) * str_len);
        memcpy(str, add_src->broadcast, str_len);
        address->broadcast = str;
    } else {
        address->broadcast = NULL;
    }
    if (add_src->peer != NULL && strlen(add_src->peer) != 0) {
        size_t str_len = strlen(add_src->peer);
        char *str = (char *) malloc(sizeof(char) * str_len);
        memcpy(str, add_src->peer, str_len);
        address->peer = str;
    } else {
        address->peer = NULL;
    }
#else
    // as we are not copying the struct data, we'll save memory related ops.
    // but this is not a good practice of api design, nor safe to be recommended.
    memcpy(address, add_src, sizeof(SCNIAddress));
    
#endif
    
    return address;
}

void
_address_release(CFAllocatorRef allocator, const void *ptr) {
    SCNIAddress *address = (SCNIAddress *)ptr;
    if (address != NULL) {
        if (address->addr != NULL) {
            free(address->addr);
        }
        if (address->netmask != NULL) {
            free(address->netmask);
        }
        if (address->broadcast != NULL) {
            free(address->broadcast);
        }
        if (address->peer != NULL) {
            free(address->peer);
        }
    }
    CFAllocatorDeallocate(allocator, (SCNIAddress *)ptr);
}

CFStringRef
_address_copy_description(const void *ptr) {
    SCNIAddress *address = (SCNIAddress *)ptr;
    return CFStringCreateWithFormat(NULL, NULL, CFSTR("[%d, %s]"), address->family, address->addr);
}

Boolean
_address_equal(const void *ptr1, const void *ptr2) {
    SCNIAddress *addr1 = (SCNIAddress *)ptr1;
    SCNIAddress *addr2 = (SCNIAddress *)ptr2;

    bool flags = (addr1->flags == addr2->flags);
    bool family = (addr1->family == addr2->family);

    bool address = false;
    if (addr1->addr != NULL && addr2->addr != NULL && strcmp(addr1->addr, addr2->addr) == 0) {
        address = true;
    }
    
    bool netmask = false;
    if (addr1->netmask != NULL && addr2->netmask != NULL && strcmp(addr1->netmask, addr2->netmask) == 0) {
        netmask = true;
    }

    return (flags && family && address && netmask);
}

CFMutableArrayRef SCNIMutableAddressArray(void) {
    CFArrayCallBacks callbacks = {0, _address_retain, _address_release, _address_copy_description, _address_equal};
    return CFArrayCreateMutable(kCFAllocatorDefault, 0, &callbacks);
}

errno_t SCNetworkInterfaceAddresses(SCNetworkInterfaceRef interface, CFMutableArrayRef results) {
    CFStringRef ifaName = SCNetworkInterfaceGetBSDName(interface);
    return ifaddresses(results, ifaName);
}

void SCNetworkInterfaceAddressRelease(CFMutableArrayRef results) {
    CFRelease(results);
}

#pragma mark - GATEWAY

static void
_SCNIGatewayRelease(SCNIGateway *gateway) {
    if (gateway != NULL) {
#if SHOULD_COPY_STRUCT_DATA
        if (gateway->ifname != NULL) {
            free(gateway->ifname);
        }
        if (gateway->addr != NULL) {
            free(gateway->addr);
        }
#endif
        free(gateway);
    }
}

static errno_t
gateways(CFMutableArrayRef results)
{
#if defined(HAVE_SYSCTL_CTL_NET)
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
        SCNIGateway* gateway = NULL;

        if (string_from_sockaddr (sa, buffer, sizeof(buffer)) == 0) {
#ifdef RTF_IFSCOPE
          bool is_default = !(msg->rtm_flags & RTF_IFSCOPE);
#else
          bool is_default = true;
#endif
          gateway = (SCNIGateway *) calloc(1, sizeof (SCNIGateway));
          
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
          CFArrayAppendValue(results, gateway);
          _SCNIGatewayRelease(gateway);
        }
      }

      /* These are aligned on a 4-byte boundary */
      ptr += len;
    }

    ptr = msgend;
  }

  free (buffer);
#endif

  return 0;
}

// CFMutableArray Callbacks for SCNIGateway

const void*
_gateway_retain(CFAllocatorRef allocator, const void *ptr) {
    SCNIGateway *gw_src = (SCNIGateway *)ptr;
    SCNIGateway *gateway = (SCNIGateway *)CFAllocatorAllocate(allocator, sizeof(SCNIGateway), 0);
    
#if SHOULD_COPY_STRUCT_DATA
    gateway->family     = gw_src->family;
    gateway->is_default = gw_src->is_default;

    if (gw_src->ifname != NULL && strlen(gw_src->ifname) != 0) {
        size_t str_len = strlen(gw_src->ifname);
        char *str = (char *) malloc(sizeof(char) * str_len);
        memcpy(str, gw_src->ifname, str_len);
        gateway->ifname = str;
    } else {
        gateway->ifname = NULL;
    }
    if (gw_src->addr != NULL && strlen(gw_src->addr) != 0) {
        size_t str_len = strlen(gw_src->addr);
        char *str = (char *) malloc(sizeof(char) * str_len);
        memcpy(str, gw_src->addr, str_len);
        gateway->addr = str;
    } else {
        gateway->addr = NULL;
    }
#else
    memcpy(gateway, gw_src, sizeof(SCNIGateway));
#endif
    
    return gateway;
}

void
_gateway_release(CFAllocatorRef allocator, const void *ptr) {
    SCNIGateway *gateway = (SCNIGateway *)ptr;
    if (gateway != NULL) {
        if (gateway->ifname != NULL) {
            free(gateway->ifname);
        }
        if (gateway->addr != NULL) {
            free(gateway->addr);
        }
    }
    CFAllocatorDeallocate(allocator, (SCNIGateway *)ptr);
}

CFStringRef
_gateway_copy_description(const void *ptr) {
    SCNIGateway *gateway = (SCNIGateway *)ptr;
    return CFStringCreateWithFormat(NULL, NULL, CFSTR("[%d, %s : %s]"), gateway->family, gateway->ifname, gateway->addr);
}

Boolean
_gateway_equal(const void *ptr1, const void *ptr2) {
    
    SCNIGateway *gw1 = (SCNIGateway *)ptr1;
    SCNIGateway *gw2 = (SCNIGateway *)ptr2;
    bool family = (gw1->family == gw2->family);

    bool ifname = false;
    if (gw1->ifname != NULL && gw2->ifname != NULL && strcmp(gw1->ifname, gw2->ifname) == 0) {
        ifname = true;
    }
    
    bool addr = false;
    if (gw1->addr != NULL  && gw2->addr != NULL && strcmp(gw1->addr, gw2->addr) == 0) {
        addr = true;
    }
    
    return (family && ifname && addr);
}

CFMutableArrayRef SCNIMutableGatewayArray(void) {
    CFArrayCallBacks callbacks = {0, _gateway_retain, _gateway_release, _gateway_copy_description, _gateway_equal};
    return CFArrayCreateMutable(kCFAllocatorDefault, 0, &callbacks);
}

errno_t SCNetworkGateways(CFMutableArrayRef results) {
    return gateways(results);
}

void SCNetworkGatewayRelease(CFMutableArrayRef results) {
    CFRelease(results);
}
