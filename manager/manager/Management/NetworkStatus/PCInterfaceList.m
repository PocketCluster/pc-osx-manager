//
//  PCInterfaceList.m
//  NETUTIL
//
//  Created by Almighty Kim on 10/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import "PCInterfaceList.h"


#include <ifaddrs.h>
#include <net/if.h>
#include <netdb.h>

#import <SystemConfiguration/SystemConfiguration.h>

@implementation PCInterfaceList

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


+ (NSArray *)all
{
    
    __block NSMutableArray *result = [NSMutableArray new];
    struct ifaddrs *allInterfaces;
    
    @autoreleasepool {

        SCDynamicStoreRef storeRef = SCDynamicStoreCreate(NULL, (CFStringRef)@"FindCurrentInterfaceIpMac", NULL, NULL);
        CFPropertyListRef global = SCDynamicStoreCopyValue(storeRef, CFSTR("State:/Network/Interface"));
        NSArray *ifaceList = [(__bridge NSArray *)global valueForKey:@"Interfaces"];
        
        // grap every interface name
        for(NSString *iface in ifaceList)
        {
            LinkInterface* interface = [LinkInterface new];
            [interface setBSDName:iface];
            [result addObject:interface];
        }

        // match intefaces and idenfity their kind
        NSArray *scIfaceList = (NSArray*) CFBridgingRelease(SCNetworkInterfaceCopyAll());
        for (id iface in scIfaceList)
        {
            SCNetworkInterfaceRef interfaceRef = (__bridge SCNetworkInterfaceRef)iface;
            NSString *bsdName = (__bridge NSString*)SCNetworkInterfaceGetBSDName(interfaceRef);
            
            for (LinkInterface* _Nonnull rface in result)
            {
                if ([[rface BSDName] isEqualToString:bsdName])
                {
                    rface.displayName = (__bridge NSString*)SCNetworkInterfaceGetLocalizedDisplayName(interfaceRef);
                    rface.hardMAC = (__bridge NSString*)SCNetworkInterfaceGetHardwareAddressString(interfaceRef);
                    rface.kind = (__bridge NSString*)SCNetworkInterfaceGetInterfaceType(interfaceRef);
                    break;
                }
                
            }
        }
        
        // Get list of all interfaces on the local machine & match ip addresses:
        if (getifaddrs(&allInterfaces) == 0)
        {
            struct ifaddrs *interface;
            
            // For each interface ...
            for (interface = allInterfaces; interface != NULL; interface = interface->ifa_next)
            {
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
                            
                            for (LinkInterface* _Nonnull rface in result)
                            {
                                if([[rface BSDName] isEqualToString:bsdName]){
                                    [rface setIp4Address:ip4Address];
                                    break;
                                }
                                
                            }
                        }
                    }
                    
                    if (addr->sa_family == AF_INET6) {
                        
                        // Convert interface address to a human readable string:
                        char host[NI_MAXHOST];
                        getnameinfo(addr, addr->sa_len, host, sizeof(host), NULL, 0, NI_NUMERICHOST);
                        
                        
                        if(strlen(host) != 0){
                            
                            NSString *bsdName = [NSString stringWithCString:interface->ifa_name encoding:NSUTF8StringEncoding];
                            NSString *ip6Address = [NSString stringWithCString:host encoding:NSUTF8StringEncoding];
                            
                            for (LinkInterface* _Nonnull rface in result)
                            {
                                if([[rface BSDName] isEqualToString:bsdName]){
                                    [rface setIp4Address:ip6Address];
                                    break;
                                }
                                
                            }
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
@end
