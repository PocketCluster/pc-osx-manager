//
//  PCInterfaceList.m
//  NETUTIL
//
//  Created by Almighty Kim on 10/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <string.h>
#import <CoreWLAN/CWWiFiClient.h>
#import <CoreWLAN/CWInterface.h>

#import "PCInterfaceStatus.h"
#import "SCNetworkInterfaces.h"
#import "LinkObserver.h"
#import "util.h"

static PCNetworkInterface**
_interface_status(unsigned int*, CFMutableArrayRef);

static PCNetworkInterface*
_pc_interface_new();

static void
_pc_interface_release(PCNetworkInterface*);

static PCNetworkInterface **
_pc_interface_array_new(unsigned int);

static void
_pc_interface_array_release(PCNetworkInterface**, unsigned int);

SCNIAddress**
_SCNIAddressNewArrray(unsigned int length);

SCNIGateway**
_SCNIGatewayNewArray(unsigned int length);

void
interface_status(pc_interface_callback);

@interface PCInterfaceStatus()
@property (readonly) LinkObserver *linkObserver;
@end

@implementation PCInterfaceStatus
@synthesize linkObserver;

#pragma mark - PROPERTIES
- (LinkObserver*) linkObserver {
    if (linkObserver) return linkObserver;
    linkObserver = [LinkObserver new];
    return linkObserver;
}

#pragma mark - METHODS
- (void) interfacesDidChange:(NSNotification*)notififcation {
    NSLog(@"Interface change detected...");
}

- (void) startMonitoring {
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(interfacesDidChange:) name:@"State:/Network/Interface" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(interfacesDidChange:) name:@"State:/Network/Interface/en0/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(interfacesDidChange:) name:@"State:/Network/Interface/en1/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(interfacesDidChange:) name:@"State:/Network/Interface/en2/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] addObserver:self selector:@selector(interfacesDidChange:) name:@"State:/Network/Interface/en3/AirPort" object:self.linkObserver];
}

- (void) stopMonitoring {
    [[NSNotificationCenter defaultCenter] removeObserver:self name:@"State:/Network/Interface" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] removeObserver:self name:@"State:/Network/Interface/en0/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] removeObserver:self name:@"State:/Network/Interface/en1/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] removeObserver:self name:@"State:/Network/Interface/en2/AirPort" object:self.linkObserver];
    [[NSNotificationCenter defaultCenter] removeObserver:self name:@"State:/Network/Interface/en3/AirPort" object:self.linkObserver];
}
@end

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

//TODO : 2) Interface status 3) Address Status 4) Gateway object 5) Async notification 6) leak check 7) isprimary? 8) interface type (thunberbolt?)
PCNetworkInterface**
_interface_status(unsigned int* pcIfaceCount, CFMutableArrayRef totalAddress) {
    
    PCNetworkInterface** pcIfaceArray = NULL;
    
    @autoreleasepool {
        CWWiFiClient *wifiClent = [CWWiFiClient sharedWiFiClient];
        
        // grap every *NETWORK* interface name
        SCDynamicStoreRef storeRef = SCDynamicStoreCreate(NULL, (CFStringRef)@"FindCurrentInterfaceIpMac", NULL, NULL);
        CFPropertyListRef global = SCDynamicStoreCopyValue(storeRef, CFSTR("State:/Network/Interface"));
        NSArray *netIfaceArray = [(__bridge NSArray *)global valueForKey:@"Interfaces"];
        NSUInteger netIfaceCount = [netIfaceArray count];
        
        *pcIfaceCount = (unsigned int) netIfaceCount;
        pcIfaceArray = _pc_interface_array_new((unsigned int) netIfaceCount);
        
        for (NSUInteger i = 0; i < netIfaceCount; i++) {
            NSString *bsdName = [netIfaceArray objectAtIndex:i];
            PCNetworkInterface *pcIface = _pc_interface_new();
            pcIface->bsdName = copy_string([bsdName UTF8String]);
            *(pcIfaceArray + i) = pcIface;
        }
        
        // match intefaces and idenfity their kind
        CFArrayRef allIfaceArray = SCNetworkInterfaceCopyAll();
        for (CFIndex aai = 0; aai < CFArrayGetCount(allIfaceArray); aai++) {
            
            SCNetworkInterfaceRef interface = CFArrayGetValueAtIndex(allIfaceArray, aai);
            NSString *nsBsdName = (__bridge NSString*)SCNetworkInterfaceGetBSDName(interface);
            
            for (NSUInteger pci = 0; pci < netIfaceCount; pci++) {
                PCNetworkInterface *pcIface = *(pcIfaceArray + pci);
                
                if (pcIface->bsdName != NULL && nsBsdName != nil && strcmp(pcIface->bsdName, [nsBsdName UTF8String]) == 0) {
                    
                    if (SCNetworkInterfaceMediaStatus(interface)) {
                        pcIface->isActive = true;
                    } else {
                        pcIface->isActive = false;
                    }
                    
                    NSString *displayName = (__bridge NSString*)SCNetworkInterfaceGetLocalizedDisplayName(interface);
                    pcIface->displayName  = copy_string([displayName UTF8String]);
                    
                    NSString *hardMAC     = (__bridge NSString*)SCNetworkInterfaceGetHardwareAddressString(interface);
                    pcIface->macAddress   = copy_string([hardMAC UTF8String]);
                    
                    NSString *kind        = (__bridge NSString*)SCNetworkInterfaceGetInterfaceType(interface);
                    pcIface->mediaType    = copy_string([kind UTF8String]);
                    
                    if ([kind isEqualToString:@"IEEE80211"] || [displayName isEqualToString:@"Wi-Fi"]) {
                        CWInterface *wifiIface = [wifiClent interfaceWithName:nsBsdName];
                        if (wifiIface != nil) {
                            pcIface->wifiPowerOff = !wifiIface.powerOn;
                        }
                    }
                    
                    CFMutableArrayRef scniAddr = SCNIMutableAddressArray();
                    errno_t err = SCNetworkInterfaceAddresses(interface, scniAddr);
                    CFIndex addrCount = CFArrayGetCount(scniAddr);
                    
                    if (err == 0 && 0 < addrCount) {
                        SCNIAddress **address = _SCNIAddressNewArrray((unsigned int) addrCount);
                        for (CFIndex scai = 0; scai < addrCount; scai++) {
                            *(address + scai) = (SCNIAddress *) CFArrayGetValueAtIndex(scniAddr, scai);
                        }
                        pcIface->address = address;
                        pcIface->addrCount = (unsigned int) addrCount;
                        
                        CFArrayAppendValue(totalAddress, scniAddr);
                    } else {
                        // since this array is empty, we'll release now.
                        SCNetworkInterfaceAddressRelease(scniAddr);
                    }
                    
                    // break-out from interface searching iteration
                    break;
                }
            }
        }
    }
    
    return pcIfaceArray;
}


PCNetworkInterface*
_pc_interface_new() {
    return (PCNetworkInterface *) calloc(1, sizeof(PCNetworkInterface));
}

void
_pc_interface_release(PCNetworkInterface* interface) {
    if (interface != NULL) {
        if (interface->address != NULL) {
            free((void*)interface->address);
        }
        if (interface->bsdName != NULL) {
            free((void*)interface->bsdName);
        }
        if (interface->displayName != NULL) {
            free((void*)interface->displayName);
        }
        if (interface->macAddress != NULL) {
            free((void*)interface->macAddress);
        }
        if (interface->mediaType != NULL) {
            free((void*)interface->mediaType);
        }
        free(interface);
    }
}

PCNetworkInterface **
_pc_interface_array_new(unsigned int length) {
    if (length == 0) {
        return NULL;
    }
    return (PCNetworkInterface **) malloc (sizeof(PCNetworkInterface*) * length);
}

void
_pc_interface_array_release(PCNetworkInterface** interfaces, unsigned int length) {
    if (interfaces != NULL && length != 0) {
        for (unsigned int i = 0; i < length; i++) {
            _pc_interface_release(*(interfaces + i));
        }
        free(interfaces);
    }
}

SCNIAddress**
_SCNIAddressNewArrray(unsigned int length) {
    if (length == 0) {
        return NULL;
    }
    return (SCNIAddress **) malloc (sizeof (SCNIAddress*) * length);
}

SCNIGateway**
_SCNIGatewayNewArray(unsigned int length) {
    if (length == 0) {
        return NULL;
    }
    return (SCNIGateway **) malloc (sizeof (SCNIGateway*) * length);
}

void
interface_status_with_callback(pc_interface_callback callback) {
    if (callback == NULL) {
        return;
    }
    
    unsigned int interfacesCount = 0;
    CFMutableArrayRef totalAddress = CFArrayCreateMutable(kCFAllocatorDefault, 0, &kCFTypeArrayCallBacks);
    PCNetworkInterface** interfaces = _interface_status(&interfacesCount, totalAddress);
    
    // we don't mind if Golang has successfully received the interface array. We only care
    // if we can safely release all the memories
    callback(interfaces, interfacesCount);
    
    // release phase
    _pc_interface_array_release(interfaces, interfacesCount);
    
    for (CFIndex i = CFArrayGetCount(totalAddress) - 1; 0 <= i; i--) {
        CFMutableArrayRef addr = (CFMutableArrayRef)CFArrayGetValueAtIndex(totalAddress, i);
        CFArrayRemoveValueAtIndex(totalAddress, i);
        SCNetworkInterfaceAddressRelease(addr);
    }
    CFRelease(totalAddress);
}

void
interface_status_with_gocall() {
    
    unsigned int interfacesCount = 0;
    CFMutableArrayRef totalAddress = CFArrayCreateMutable(kCFAllocatorDefault, 0, &kCFTypeArrayCallBacks);
    PCNetworkInterface** interfaces = _interface_status(&interfacesCount, totalAddress);
    
    /* PLACE GOCALL here */
    
    
    // release phase
    _pc_interface_array_release(interfaces, interfacesCount);
    
    for (CFIndex i = CFArrayGetCount(totalAddress) - 1; 0 <= i; i--) {
        CFMutableArrayRef addr = (CFMutableArrayRef)CFArrayGetValueAtIndex(totalAddress, i);
        CFArrayRemoveValueAtIndex(totalAddress, i);
        SCNetworkInterfaceAddressRelease(addr);
    }
    CFRelease(totalAddress);
}
