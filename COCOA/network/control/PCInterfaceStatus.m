//
//  PCInterfaceList.m
//  NETUTIL
//
//  Created by Almighty Kim on 10/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <net/if.h>
#import <sys/socket.h>
#import <string.h>
#import <CoreWLAN/CWWiFiClient.h>
#import <CoreWLAN/CWInterface.h>

#import "PCInterfaceStatus.h"
#import "SCNetworkInterfaces.h"
#import "LinkObserver.h"
#import "util.h"

static const CFStringRef kPocketClusterPrimaryInterface = CFSTR("PocketClusterPrimary");

static PCNetworkInterface*
_pc_interface_new();

static void
_pc_interface_release(PCNetworkInterface*);

static PCNetworkInterface **
_pc_interface_array_new(unsigned int);

static void
_pc_interface_array_release(PCNetworkInterface**, unsigned int);

static SCNIAddress**
_SCNIAddressArrrayNew(unsigned int length);

static SCNIGateway**
_SCNIGatewayArrayNew(unsigned int length);

static void
_SCNIGatewayArrayRelease(SCNIGateway**);

static void
_primary_interface_address(const char**, const char**);

static PCNetworkInterface**
_interface_status(unsigned int*, CFMutableArrayRef);

static SCNIGateway**
_gateway_status(CFMutableArrayRef, unsigned int*);

@interface PCInterfaceStatus()<LinkObserverNotification>
@property (readonly) LinkObserver *linkObserver;
@end

@implementation PCInterfaceStatus {
    BOOL    _shouldStartMonitor;
}
@synthesize linkObserver;
-(instancetype)init {
    self = [super init];
    if (self) {
        _shouldStartMonitor = NO;
        [self.linkObserver setDelegate:self];
    }
    return self;
}

#pragma mark - PROPERTIES
- (LinkObserver*) linkObserver {
    if (linkObserver) {
        return linkObserver;
    }
    linkObserver = [LinkObserver new];
    return linkObserver;
}

#pragma mark - METHODS
- (void)startMonitoring {
    _shouldStartMonitor = YES;
}

- (void)stopMonitoring {
    _shouldStartMonitor = NO;
}

#pragma mark - LinkObserverNotification protocol
- (void)networkConfigurationDidChange:(LinkObserver *)observer configChanged:(NSDictionary *)configChanged {
    NSLog(@"Network configuration has changed");
}

- (void)networkConfigurationDidChange:(LinkObserver *)observer {
    NSLog(@"Network configuration has changed %@ - %@", observer, [[NSDate date] description]);
    if ([NSThread isMainThread]) {
        printf("!!! THIS IS M.A.I.N THREAD!!!\n\n");
    } else {
        printf("!!! this not main thread!!!\n\n");
    }
}

@end

#pragma mark - ALLOCATION / RELEASE
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
_SCNIAddressArrrayNew(unsigned int length) {
    if (length == 0) {
        return NULL;
    }
    return (SCNIAddress **) malloc (sizeof (SCNIAddress*) * length);
}

SCNIGateway**
_SCNIGatewayArrayNew(unsigned int length) {
    if (length == 0) {
        return NULL;
    }
    return (SCNIGateway **) malloc (sizeof (SCNIGateway*) * length);
}

void
_SCNIGatewayArrayRelease(SCNIGateway** gatewayArray) {
    if (gatewayArray != NULL) {
        free((void *)gatewayArray);
    }
}

#pragma mark - INTERFACE STATUS ACQUISITION
// http://lists.apple.com/archives/macnetworkprog/2006/Oct/msg00007.html
void
_primary_interface_address(const char** primary_interface, const char** primary_address) {
    SCDynamicStoreRef       store = NULL;
    CFStringRef             globalKeys = NULL;
    CFDictionaryRef         ipv4State = NULL;
    CFStringRef             primaryService = NULL;
    CFStringRef             primaryInterface = NULL;
    CFStringRef             ipv4Key = NULL;
    CFDictionaryRef         serviceDict = NULL;
    CFArrayRef              addresses = NULL;
    CFStringRef             address = NULL;
    
    @autoreleasepool {
        store = SCDynamicStoreCreate(NULL, kPocketClusterPrimaryInterface, NULL, NULL);
        if (store != NULL) {
            globalKeys = SCDynamicStoreKeyCreateNetworkGlobalEntity(kCFAllocatorDefault, kSCDynamicStoreDomainState, kSCEntNetIPv4);
        }
        if (globalKeys != NULL) {
            ipv4State = (CFDictionaryRef) SCDynamicStoreCopyValue(store, globalKeys);
        }
        if (ipv4State != NULL) {
            primaryService = (CFStringRef) CFDictionaryGetValue(ipv4State, kSCDynamicStorePropNetPrimaryService);
            primaryInterface = (CFStringRef) CFDictionaryGetValue(ipv4State, kSCDynamicStorePropNetPrimaryInterface);
        }
        if (primaryInterface != NULL) {
            *primary_interface = CFStringCopyToCString(primaryInterface);
        }
        if (primaryService != NULL) {
            ipv4Key = SCDynamicStoreKeyCreateNetworkServiceEntity(NULL, kSCDynamicStoreDomainState, primaryService, kSCEntNetIPv4);
        }
        if (ipv4Key != NULL) {
            serviceDict = SCDynamicStoreCopyValue(store, ipv4Key);
        }
        if (serviceDict != NULL) {
            addresses = CFDictionaryGetValue(serviceDict, kSCPropNetIPv4Addresses);
        }
        if (addresses != NULL && CFArrayGetCount(addresses) != 0) {
            address = CFArrayGetValueAtIndex(addresses, 0);
        }
        if (address != NULL) {
            *primary_address = CFStringCopyToCString(address);
        }
        
        if (serviceDict != NULL)        CFRelease(serviceDict);
        if (ipv4Key != NULL)            CFRelease(ipv4Key);
        if (ipv4State != NULL)          CFRelease(ipv4State);
        if (globalKeys != NULL)         CFRelease(globalKeys);
        if (store != NULL)              CFRelease(store);
    }
}


/*!
	@function _interface_status
	@discussion Returns interfaces with ipv4 (dotted format) addresses (no ipv6) linked to the interfaces.
	@param pcIfaceCount interface count
           allAddresses A mutable CF array where addresses should be contained to be released later
	@result The list of interfaces
 */
PCNetworkInterface**
_interface_status(unsigned int* pcIfaceCount, CFMutableArrayRef allAddresses) {
    
    static const CFStringRef kMediaTypeWIFI = CFSTR("Wi-Fi");
    static const CFStringRef kMediaTypeIEEE80211 = CFSTR("IEEE80211");
    static char ifaceBSDName[256];
    
    PCNetworkInterface** pcIfaceArray = NULL;
    const char* primaryAddress;
    const char* primaryInterface;
    
    _primary_interface_address(&primaryInterface, &primaryAddress);
    
    @autoreleasepool {
        CWWiFiClient *wifiClent = [CWWiFiClient sharedWiFiClient];
        
        // grap every *NETWORK* interface name
        SCDynamicStoreRef storeRef = SCDynamicStoreCreate(NULL, (CFStringRef)@"FindCurrentInterfaceIpMac", NULL, NULL);
        CFPropertyListRef global = SCDynamicStoreCopyValue(storeRef, CFSTR("State:/Network/Interface"));
        CFArrayRef netIfaceArray = (CFArrayRef) CFDictionaryGetValue(global, CFSTR("Interfaces"));
        CFIndex netIfaceCount = CFArrayGetCount(netIfaceArray);
        
        *pcIfaceCount = (unsigned int) netIfaceCount;
        pcIfaceArray = _pc_interface_array_new((unsigned int) netIfaceCount);
        
        for (CFIndex i = 0; i < netIfaceCount; i++) {
            CFStringRef bsdName = (CFStringRef) CFArrayGetValueAtIndex(netIfaceArray, i);
            PCNetworkInterface *pcIface = _pc_interface_new();
            pcIface->bsdName = CFStringCopyToCString(bsdName);
            if (primaryInterface != NULL && pcIface->bsdName != NULL && strcmp(primaryInterface, pcIface->bsdName) == 0) {
                pcIface->isPrimary = true;
            }
            *(pcIfaceArray + i) = pcIface;
        }
        
        // match intefaces and idenfity their kind
        CFArrayRef allIfaceArray = SCNetworkInterfaceCopyAll();
        for (CFIndex aai = 0; aai < CFArrayGetCount(allIfaceArray); aai++) {
            // to reduce memory pressure, we'll have a nested autorelease pool
            @autoreleasepool {
                SCNetworkInterfaceRef interface = CFArrayGetValueAtIndex(allIfaceArray, aai);
                CFStringRef bsdName = SCNetworkInterfaceGetBSDName(interface);
                CFStringGetCString(bsdName, ifaceBSDName, 256, kCFStringEncodingUTF8);
                
                for (NSUInteger pci = 0; pci < netIfaceCount; pci++) {
                    PCNetworkInterface *pcIface = *(pcIfaceArray + pci);
                    
                    if (pcIface->bsdName != NULL && strlen(ifaceBSDName) != 0 && strcmp(pcIface->bsdName, ifaceBSDName) == 0) {
                        
                        if (SCNetworkInterfaceMediaStatus(interface)) {
                            pcIface->isActive = true;
                        } else {
                            pcIface->isActive = false;
                        }
                        
                        CFStringRef displayName = SCNetworkInterfaceGetLocalizedDisplayName(interface);
                        pcIface->displayName    = CFStringCopyToCString(displayName);
                        
                        CFStringRef hardMAC     = SCNetworkInterfaceGetHardwareAddressString(interface);
                        pcIface->macAddress     = CFStringCopyToCString(hardMAC);
                        
                        CFStringRef mediaType   = SCNetworkInterfaceGetInterfaceType(interface);
                        pcIface->mediaType      = CFStringCopyToCString(mediaType);
                        
                        
                        if (CFStringCompare(mediaType, kMediaTypeIEEE80211, kCFCompareCaseInsensitive) == kCFCompareEqualTo ||
                            CFStringCompare(displayName, kMediaTypeWIFI, kCFCompareCaseInsensitive) == kCFCompareEqualTo) {
                            CWInterface *wifiIface = [wifiClent interfaceWithName:(__bridge NSString*)bsdName];
                            if (wifiIface != nil) {
                                pcIface->wifiPowerOff = !wifiIface.powerOn;
                            }
                        }
                        
                        CFMutableArrayRef scniAddr = SCNIMutableAddressArray();
                        errno_t err = SCNetworkInterfaceAddresses(interface, scniAddr);
                        CFIndex addrCount = CFArrayGetCount(scniAddr);
                        
                        if (err == 0 && 0 < addrCount) {

                            // TODO : this should be refactored into using realloc. Iterating twice on the same array is kind of ¯\_(ツ)_/¯
                            // we only take IPv4 and Valid (UP & RUNNING) addresses. IPv6 will be counted later
                            unsigned int validAddrCount = 0;
                            
                            for (CFIndex scai = 0; scai < addrCount; scai++) {
                                SCNIAddress *addr = (SCNIAddress *) CFArrayGetValueAtIndex(scniAddr, scai);
                                if ((addr->flags & (IFF_UP|IFF_RUNNING|IFF_LOOPBACK)) == (IFF_UP|IFF_RUNNING)) {
                                    if (addr->family == AF_INET) {
                                        validAddrCount++;
                                    }
                                }
                            }
                            
                            if (0 < validAddrCount) {
                                unsigned int addressIndex = 0;
                                SCNIAddress **address = _SCNIAddressArrrayNew(validAddrCount);
                                
                                for (CFIndex scai = 0; scai < addrCount; scai++) {
                                    SCNIAddress *addr = (SCNIAddress *) CFArrayGetValueAtIndex(scniAddr, scai);
                                    if ((addr->flags & (IFF_UP|IFF_RUNNING|IFF_LOOPBACK)) == (IFF_UP|IFF_RUNNING)) {
                                        if (addr->family == AF_INET) {
                                            if (primaryAddress != NULL && addr->addr != NULL && strcmp(primaryAddress, addr->addr) == 0) {
                                                addr->is_primary = true;
                                            }
                                            *(address + addressIndex) = addr;
                                            addressIndex++;
                                        }
                                    }
                                }
                                pcIface->address = address;
                                pcIface->addrCount = validAddrCount;
                                
                                CFArrayAppendValue(allAddresses, scniAddr);
                            } else {
                                // since there is no valid address, we'll release now.
                                SCNetworkInterfaceAddressRelease(scniAddr);
                            }

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
    }
    
    if (primaryAddress != NULL) {
        free((void *) primaryAddress);
    }
    if (primaryInterface != NULL) {
        free((void *) primaryInterface);
    }
    return pcIfaceArray;
}


void
interface_status_with_callback(pc_interface_callback callback) {
    if (callback == NULL) {
        return;
    }
    
    unsigned int interfacesCount = 0;
    CFMutableArrayRef allAddresses = CFArrayCreateMutable(kCFAllocatorDefault, 0, &kCFTypeArrayCallBacks);
    PCNetworkInterface** interfaces = _interface_status(&interfacesCount, allAddresses);
    
    // we don't mind if Golang has successfully received the interface array. We only care
    // if we can safely release all the memories
    if (interfaces != NULL && 0 < interfacesCount) {
        callback(interfaces, interfacesCount);
    } else {
        callback(NULL, 0);
    }
    
    // release phase
    _pc_interface_array_release(interfaces, interfacesCount);
    
    for (CFIndex i = CFArrayGetCount(allAddresses) - 1; 0 <= i; i--) {
        CFMutableArrayRef addr = (CFMutableArrayRef)CFArrayGetValueAtIndex(allAddresses, i);
        CFArrayRemoveValueAtIndex(allAddresses, i);
        SCNetworkInterfaceAddressRelease(addr);
    }
    CFRelease(allAddresses);
}

void
interface_status_with_gocall() {
    
    unsigned int interfacesCount = 0;
    CFMutableArrayRef allAddresses = CFArrayCreateMutable(kCFAllocatorDefault, 0, &kCFTypeArrayCallBacks);
    PCNetworkInterface** interfaces = _interface_status(&interfacesCount, allAddresses);
    
    /* PLACE GOCALL here */
    
    
    // release phase
    _pc_interface_array_release(interfaces, interfacesCount);
    
    for (CFIndex i = CFArrayGetCount(allAddresses) - 1; 0 <= i; i--) {
        CFMutableArrayRef addr = (CFMutableArrayRef)CFArrayGetValueAtIndex(allAddresses, i);
        CFArrayRemoveValueAtIndex(allAddresses, i);
        SCNetworkInterfaceAddressRelease(addr);
    }
    CFRelease(allAddresses);
}

#pragma mark - GATEWAY STATUS ACQUISITION

SCNIGateway**
_gateway_status(CFMutableArrayRef allGatways, unsigned int *gatewayCount) {
    SCNIGateway** scniGateways = NULL;
    @autoreleasepool {
        errno_t err = SCNetworkGateways(allGatways);
        CFIndex count = CFArrayGetCount(allGatways);
        if (err == 0 && 0 < count) {
            
            // TODO : this should be refactored into using realloc.
            // we only take IPv4 gateways. IPv6 will be counted later
            unsigned int validGateways = 0;
            for (CFIndex i = 0; i < count; i++) {
                SCNIGateway *gw = (SCNIGateway*)CFArrayGetValueAtIndex(allGatways, i);
                if (gw->family == AF_INET) {
                    validGateways++;
                }
            }

            if (0 < validGateways) {
                unsigned int gwIndex = 0;
                scniGateways = _SCNIGatewayArrayNew(validGateways);
                for (CFIndex i = 0; i < count; i++) {
                    SCNIGateway *gw = (SCNIGateway*)CFArrayGetValueAtIndex(allGatways, i);
                    if (gw->family == AF_INET) {
                        *(scniGateways + gwIndex) = gw;
                        gwIndex++;
                    }
                }
                *gatewayCount = validGateways;
            }
        }
    }
    return scniGateways;
}

CF_EXPORT void
gateway_status_with_callback(scni_gateway_callback callback) {
    if (callback == NULL) {
        return;
    }

    unsigned int gatewayCount = 0;
    CFMutableArrayRef allGateways = SCNIMutableGatewayArray();
    SCNIGateway** scniGateways = _gateway_status(allGateways, &gatewayCount);
    
    if (scniGateways != NULL && 0 < gatewayCount) {
        callback(scniGateways, gatewayCount);
    } else {
        callback(NULL, 0);
    }
    
    // release phase
    _SCNIGatewayArrayRelease(scniGateways);
    
    // release phase
    SCNetworkGatewayRelease(allGateways);
}

void
gateway_status_with_gocall() {

    return;
}
