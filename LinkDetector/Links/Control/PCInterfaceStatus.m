//
//  PCInterfaceList.m
//  NETUTIL
//
//  Created by Almighty Kim on 10/24/15.
//  Copyright © 2015 io.pocketcluster. All rights reserved.
//

#import <string.h>

#import "PCInterfaceStatus.h"
#import "SCNetworkInterfaces.h"
#import "LinkObserver.h"
#import "util.h"

static PCNetworkInterface**
_interface_status(unsigned int*, CFMutableArrayRef);

static void
_pc_interface_release(PCNetworkInterface*);

static void
_pc_interface_array_release(PCNetworkInterface**, unsigned int);

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


//TODO : 2) Interface status 3) Address Status 4) Gateway object 5) Async notification 6) leak check 7) isprimary? 8) interface type (thunberbolt?)
PCNetworkInterface**
_interface_status(unsigned int* pcIfaceCount, CFMutableArrayRef totalAddress) {
    
    PCNetworkInterface** pcIfaceArray = NULL;
    
    @autoreleasepool {
        // grap every *NETWORK* interface name
        SCDynamicStoreRef storeRef = SCDynamicStoreCreate(NULL, (CFStringRef)@"FindCurrentInterfaceIpMac", NULL, NULL);
        CFPropertyListRef global = SCDynamicStoreCopyValue(storeRef, CFSTR("State:/Network/Interface"));
        NSArray *netIfaceArray = [(__bridge NSArray *)global valueForKey:@"Interfaces"];
        NSUInteger netIfaceCount = [netIfaceArray count];
        
        *pcIfaceCount = (unsigned int) netIfaceCount;
        pcIfaceArray = (PCNetworkInterface **) malloc(sizeof(PCNetworkInterface*) * netIfaceCount);
        
        for (NSUInteger i = 0; i < netIfaceCount; i++) {
            NSString *bsdName = [netIfaceArray objectAtIndex:i];
            PCNetworkInterface *pcIface = (PCNetworkInterface *) calloc(1, sizeof(PCNetworkInterface));
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
                    
                    CFMutableArrayRef scniAddr = SCNIMutableAddressArray();
                    errno_t err = SCNetworkInterfaceAddresses(interface, scniAddr);
                    CFIndex addrCount = CFArrayGetCount(scniAddr);
                    
                    if (err == 0 && 0 < addrCount) {
                        SCNIAddress **address = (SCNIAddress **) malloc (sizeof (SCNIAddress*) * addrCount);
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
                }
            }
        }
    }
    
    return pcIfaceArray;
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

void
_pc_interface_array_release(PCNetworkInterface** interfaces, unsigned int length) {
    if (interfaces != NULL) {
        for (unsigned int i = 0; i < length; i++) {
            _pc_interface_release(*(interfaces + i));
        }
        free(interfaces);
    }
}

void
interface_status(pc_interface_callback callback) {
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
