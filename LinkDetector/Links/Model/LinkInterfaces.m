/*
 Copyright (c) 2015 funkensturm. https://github.com/halo/LinkLiar
 
 Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the
 "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish,
 distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to
 the following conditions:
 
 The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 
 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
 LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
 WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

#include <ifaddrs.h>
#include <net/if.h>
#include <netdb.h>
#import <SystemConfiguration/SystemConfiguration.h>

#import "LinkInterfaces.h"
#import "LinkInterface.h"
#import "SCNetworkInterfaces.h"

@implementation LinkInterfaces

+ (LinkInterface*) interfaceByBSDNumber:(NSInteger)number {
  for (LinkInterface* interface in [self all]) {
    if (interface.BSDNumber == number) return interface;
  }
  return NULL;
}

+ (BOOL) leaking {
  for (LinkInterface* interface in [self all]) {
    if (interface.hasOriginalMAC) return YES;
  }
  return NO;
}

+ (NSArray*) all {
  NSMutableArray *result = [NSMutableArray new];
  @autoreleasepool {
    
    NSArray *interfaces = (NSArray*) CFBridgingRelease(SCNetworkInterfaceCopyAll());
   
    for (id interface_ in interfaces) {
      SCNetworkInterfaceRef interfaceRef = (__bridge SCNetworkInterfaceRef)interface_;

      LinkInterface* interface = [LinkInterface new];
      interface.BSDName = (__bridge NSString*)SCNetworkInterfaceGetBSDName(interfaceRef);
      interface.displayName = (__bridge NSString*)SCNetworkInterfaceGetLocalizedDisplayName(interfaceRef);
      interface.hardMAC = (__bridge NSString*)SCNetworkInterfaceGetHardwareAddressString(interfaceRef);
      interface.kind = (__bridge NSString*)SCNetworkInterfaceGetInterfaceType(interfaceRef);
      
      // You can only change MAC addresses of Ethernet and Wi-Fi adapters
      if (![interface.kind isEqualToString:@"Ethernet"] && ![interface.kind isEqualToString:@"IEEE80211"]) continue;
      // If there is no internal MAC this is to be ignored
      if (!interface.hardMAC) continue;
      // Bluetooth can also be filtered out
      if ([interface.displayName containsString:@"tooth"]) continue;
      // iPhones etc. are not spoofable either
      if ([interface.displayName containsString:@"iPhone"] || [interface.displayName containsString:@"iPad"] || [interface.displayName containsString:@"iPod"]) continue;
      // Internal Thunderbolt interfaces cannot be spoofed either
      if ([interface.displayName containsString:@"Thunderbolt 1"] || [interface.displayName containsString:@"Thunderbolt 2"] || [interface.displayName containsString:@"Thunderbolt 3"] || [interface.displayName containsString:@"Thunderbolt 4"] || [interface.displayName containsString:@"Thunderbolt 5"]) continue;
      // If this interface is not in ifconfig, it's probably Bluetooth
      if (!interface.softMAC) continue;
      
      [result addObject:interface];
    }
    
    return (NSArray*)result;
  }
}

#if 0
+ (NSArray*) allInterfaces
{
    __block NSMutableArray *result = [NSMutableArray new];
    @autoreleasepool {
        
        SCDynamicStoreRef storeRef = SCDynamicStoreCreate(NULL, (CFStringRef)@"FindCurrentInterfaceIpMac", NULL, NULL);
        CFPropertyListRef global = SCDynamicStoreCopyValue(storeRef, CFSTR("State:/Network/Interface"));
        NSArray *ifaceList = [(__bridge NSArray *)global valueForKey:@"Interfaces"];
        
        // grap every interface name
        for(NSString *iface in ifaceList) {
            LinkInterface* interface = [LinkInterface new];
            [interface setBSDName:iface];
            [result addObject:interface];
        }
        
        // match intefaces and idenfity their kind
        NSArray *scIfaceList = (NSArray*) CFBridgingRelease(SCNetworkInterfaceCopyAll());
        for (id iface in scIfaceList) {
            SCNetworkInterfaceRef interfaceRef = (__bridge SCNetworkInterfaceRef)iface;
            NSString *bsdName = (__bridge NSString*)SCNetworkInterfaceGetBSDName(interfaceRef);
            [result enumerateObjectsUsingBlock:^(LinkInterface *link, NSUInteger idx, BOOL *stop) {
                if ([[link BSDName] isEqualToString:bsdName]) {
                    link.displayName = (__bridge NSString*)SCNetworkInterfaceGetLocalizedDisplayName(interfaceRef);
                    link.hardMAC = (__bridge NSString*)SCNetworkInterfaceGetHardwareAddressString(interfaceRef);
                    link.kind = (__bridge NSString*)SCNetworkInterfaceGetInterfaceType(interfaceRef);
                    *stop = YES;
                }
            }];
        }
        
        // Get list of all interfaces on the local machine & match ip addresses:
        struct ifaddrs *allInterfaces;
        if (getifaddrs(&allInterfaces) == 0) {
            
            struct ifaddrs *interface;
            
            // For each interface ...
            for (interface = allInterfaces; interface != NULL; interface = interface->ifa_next) {
                unsigned int flags = interface->ifa_flags;
                struct sockaddr *addr = interface->ifa_addr;
                
                // Check for running IPv4, IPv6 interfaces. Skip the loopback interface.
                if ((flags & (IFF_UP|IFF_RUNNING|IFF_LOOPBACK)) == (IFF_UP|IFF_RUNNING)) {
                    
                    if (addr->sa_family == AF_INET) {
                        
                        // Convert interface address to a human readable string:
                        char host[NI_MAXHOST];
                        getnameinfo(addr, addr->sa_len, host, sizeof(host), NULL, 0, NI_NUMERICHOST);
                        
                        if(strlen(host) != 0){
                            
                            NSString *bsdName = [NSString stringWithCString:interface->ifa_name encoding:NSUTF8StringEncoding];
                            NSString *ip4Address = [NSString stringWithCString:host encoding:NSUTF8StringEncoding];
                            
                            [result enumerateObjectsUsingBlock:^(LinkInterface *link, NSUInteger idx, BOOL *stop) {
                                if([[link BSDName] isEqualToString:bsdName]){
                                    [link setIp4Address:ip4Address];
                                    *stop = YES;
                                }
                            }];
                        }
                    }
                    
                    if (addr->sa_family == AF_INET6) {
                        // Convert interface address to a human readable string:
                        char host[NI_MAXHOST];
                        getnameinfo(addr, addr->sa_len, host, sizeof(host), NULL, 0, NI_NUMERICHOST);
                        
                        if(strlen(host) != 0){
                            NSString *bsdName = [NSString stringWithCString:interface->ifa_name encoding:NSUTF8StringEncoding];
                            NSString *ip6Address = [NSString stringWithCString:host encoding:NSUTF8StringEncoding];
                            [result enumerateObjectsUsingBlock:^(LinkInterface *link, NSUInteger idx, BOOL *stop) {
                                if([[link BSDName] isEqualToString:bsdName]){
                                    [link setIp6Address:ip6Address];
                                    *stop = YES;
                                }
                            }];
                        }
                    }
                }
            }
            freeifaddrs(allInterfaces);
        }
        
        CFRelease(global);
        CFRelease(storeRef);
    }
    
    return (NSArray*)result;
}
#endif

+ (NSArray*) allInterfaces
{
    __block NSMutableArray *result = [NSMutableArray new];
    @autoreleasepool {
        
        
        
        // grap every *NETWORK* interface name
        SCDynamicStoreRef storeRef = SCDynamicStoreCreate(NULL, (CFStringRef)@"FindCurrentInterfaceIpMac", NULL, NULL);
        CFPropertyListRef global = SCDynamicStoreCopyValue(storeRef, CFSTR("State:/Network/Interface"));
        NSArray *netIfaceList = [(__bridge NSArray *)global valueForKey:@"Interfaces"];
        for (NSString *iface in netIfaceList) {
            LinkInterface* interface = [LinkInterface new];
            [interface setBSDName:iface];
            [result addObject:interface];
        }
        
        // match intefaces and idenfity their kind
        CFArrayRef ifaceList = SCNetworkInterfaceCopyAll();
        for (CFIndex i = 0; i < CFArrayGetCount(ifaceList); i++) {
            
            SCNetworkInterfaceRef interface = CFArrayGetValueAtIndex(ifaceList, i);
            CFStringRef cfBsdName = SCNetworkInterfaceGetBSDName(interface);
            NSString *nsBsdName = (__bridge NSString*)cfBsdName;


            if (SCNetworkInterfaceMediaStatus(interface)) {
                NSLog(@"Interface %@ is active", nsBsdName);
            } else {
                NSLog(@"Interface %@ is inactive", nsBsdName);
            }
            
            
            NSLog(@"SCNIAddress size %ld / SCNIGateway %ld\n", sizeof(SCNIAddress), sizeof(SCNIGateway));

//TODO : 2) Interface status 3) Address Status 4) Gateway object 5) Async notification 6) leak check 7) isprimary? 8) interface type (thunberbolt?)
            
            for (LinkInterface *link in result) {
                
                if ([[link BSDName] isEqualToString:nsBsdName]) {
                    link.displayName = (__bridge NSString*)SCNetworkInterfaceGetLocalizedDisplayName(interface);
                    link.hardMAC     = (__bridge NSString*)SCNetworkInterfaceGetHardwareAddressString(interface);
                    link.kind        = (__bridge NSString*)SCNetworkInterfaceGetInterfaceType(interface);
                    
                    if (true)
                    {
                        CFMutableArrayRef results = SCNIMutableAddressArray();
                        errno_t err = SCNetworkInterfaceAddresses(interface, results);
                        CFIndex addrCount = CFArrayGetCount(results);
                        
                        if (err == 0 && addrCount != 0) {
                            for (CFIndex a = 0; a < addrCount; a++) {
                                SCNIAddress *addr = (SCNIAddress *)CFArrayGetValueAtIndex(results, a);
                                printf("addr : %s | flag 0x%X\n",addr->addr, addr->flags);
                            }
                        } else {
                            printf("ERRNO %d / ADDR COUNT %ld\n", err, addrCount);
                        }
                        SCNetworkInterfaceAddressRelease(results);
                        break;
                    }
                }
            }
        }
        
        if (true)
        {
            CFMutableArrayRef gatewayList = SCNIMutableGatewayArray();
            errno_t err = SCNetworkGateways(gatewayList);
            CFIndex gatewayCount = CFArrayGetCount(gatewayList);
            
            if (err == 0 && gatewayCount != 0) {
                for (CFIndex a = 0; a < gatewayCount; a++) {
                    SCNIGateway *gw = (SCNIGateway *)CFArrayGetValueAtIndex(gatewayList, a);
                    printf("gw addr : %s\n",gw->addr);
                }
            } else {
                printf("ERRNO %d / ADDR COUNT %ld\n", err, gatewayCount);
            }
            SCNetworkGatewayRelease(gatewayList);
        }
        
#if 0
        // Get list of all interfaces on the local machine & match ip addresses:
        struct ifaddrs *allInterfaces;
        if (getifaddrs(&allInterfaces) == 0) {
            
            struct ifaddrs *interface;
            
            // For each interface ...
            for (interface = allInterfaces; interface != NULL; interface = interface->ifa_next) {
                unsigned int flags = interface->ifa_flags;
                struct sockaddr *addr = interface->ifa_addr;
                
                // Check for running IPv4, IPv6 interfaces. Skip the loopback interface.
                if ((flags & (IFF_UP|IFF_RUNNING|IFF_LOOPBACK)) == (IFF_UP|IFF_RUNNING)) {
                    
                    if (addr->sa_family == AF_INET) {
                        
                        // Convert interface address to a human readable string:
                        char host[NI_MAXHOST];
                        getnameinfo(addr, addr->sa_len, host, sizeof(host), NULL, 0, NI_NUMERICHOST);
                        
                        if(strlen(host) != 0){
                            
                            NSString *bsdName = [NSString stringWithCString:interface->ifa_name encoding:NSUTF8StringEncoding];
                            NSString *ip4Address = [NSString stringWithCString:host encoding:NSUTF8StringEncoding];
                            
                            [result enumerateObjectsUsingBlock:^(LinkInterface *link, NSUInteger idx, BOOL *stop) {
                                if([[link BSDName] isEqualToString:bsdName]){
                                    [link setIp4Address:ip4Address];
                                    *stop = YES;
                                }
                            }];
                        }
                    }
                    
                    if (addr->sa_family == AF_INET6) {
                        // Convert interface address to a human readable string:
                        char host[NI_MAXHOST];
                        getnameinfo(addr, addr->sa_len, host, sizeof(host), NULL, 0, NI_NUMERICHOST);
                        
                        if(strlen(host) != 0){
                            NSString *bsdName = [NSString stringWithCString:interface->ifa_name encoding:NSUTF8StringEncoding];
                            NSString *ip6Address = [NSString stringWithCString:host encoding:NSUTF8StringEncoding];
                            [result enumerateObjectsUsingBlock:^(LinkInterface *link, NSUInteger idx, BOOL *stop) {
                                if([[link BSDName] isEqualToString:bsdName]){
                                    [link setIp6Address:ip6Address];
                                    *stop = YES;
                                }
                            }];
                        }
                    }
                }
            }
            freeifaddrs(allInterfaces);
        }
#endif
        
        
        CFRelease(global);
        CFRelease(storeRef);
    }
    
    return (NSArray*)result;
}

@end
