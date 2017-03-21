//
//  AppDelegate.m
//  LinkMon
//
//  Created by Almighty Kim on 10/26/16.
//  Copyright © 2016 PocketCluster. All rights reserved.
//

#import "AppDelegate.h"
#import "PCInterfaceStatus.h"


bool
pc_interface_list(PCNetworkInterface** interfaces, unsigned int count) {
    printf("\n\n---- total intefaces count %d ----\n\n", count);
    for (unsigned int i = 0; i < count; i++) {
        
        PCNetworkInterface *iface = *(interfaces + i);
        printf("wifiPowerOff : %d\n",iface->wifiPowerOff);
        printf("isActive : %d\n",iface->isActive);
        printf("isPrimary : %d\n",iface->isPrimary);
        printf("addrCount: %d\n",iface->addrCount);
        
        if (iface->addrCount != 0) {
            for (unsigned int i = 0; i < iface->addrCount; i++) {
                SCNIAddress *addr = *(iface->address + i);
                printf("\tflags  : %x\n", addr->flags);
                printf("\tfamily : %d\n", addr->family);
                printf("\tis_primary : %d\n", addr->is_primary);
                printf("\taddr : %s\n", addr->addr);
                printf("\tnetmask : %s\n", addr->netmask);
                printf("\tbroadcast : %s\n", addr->broadcast);
                printf("\tpeer : %s\n\t--------------------\n", addr->peer);
            }
        }
        
        printf("bsdName : %s\n",iface->bsdName);
        printf("displayName: %s\n",iface->displayName);
        printf("macAddress: %s\n",iface->macAddress);
        printf("mediaType: %s\n--------------------\n",iface->mediaType);
    }
    
    if ([NSThread isMainThread]) {
        printf("!!! THIS IS M.A.I.N THREAD!!!\n\n");
    } else {
        printf("!!! this not main thread!!!\n\n");
    }
    return true;
}

bool
gateway_list(SCNIGateway** gateways, unsigned int count) {
    printf("\n\n---- Total gateway count %d ----\n", count);
    for (unsigned int i = 0; i < count; i++) {
        SCNIGateway *gw = *(gateways + i);
        printf("family : %d\n",gw->family);
        printf("is_default : %d\n",gw->is_default);
        printf("ifname : %s\n",gw->ifname);
        printf("addr: %s\n",gw->addr);
    }
    if ([NSThread isMainThread]) {
        printf("!!! THIS IS M.A.I.N THREAD!!!\n\n");
    } else {
        printf("!!! this not main thread!!!\n\n");
    }
    return true;
}

@interface AppDelegate ()<PCInterfaceStatusNotification>

@property (weak) IBOutlet NSWindow *window;
@property (strong) PCInterfaceStatus *status;
@end

@implementation AppDelegate

-(void)PCInterfaceStatusChanged:(PCInterfaceStatus *)monitor interfaceStatus:(PCNetworkInterface**)status count:(unsigned int)count {
    pc_interface_list(status, count);
}

-(void)PCGatewayStatusChanged:(PCInterfaceStatus *)monitor gatewayStatus:(SCNIGateway**)status count:(unsigned int)count {
    gateway_list(status, count);
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
    
    interface_status_with_callback(&pc_interface_list);
    gateway_status_with_callback(&gateway_list);
    NSLog(@"\n--- --- --- CALLBACK C CALL ENDED --- --- ---");
    
    self.status = [[PCInterfaceStatus alloc] initWithStatusAudience:self];
    [self.status startMonitoring];
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
    [self.status stopMonitoring];
    self.status = nil;
}

@end
