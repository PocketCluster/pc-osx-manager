//
//  AppDelegate+Netmonitor.m
//  manager
//
//  Created by Almighty Kim on 4/10/17.
//  Copyright © 2017 io.pocketcluster. All rights reserved.
//

#import "AppDelegate+Netmonitor.h"
#import "pc-core.h"

bool
PCUpdateInterfaceList(PCNetworkInterface** interfaces, unsigned int count) {
#ifdef COMPARE_NATIVE_GO_OUTPUT
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
#endif
    
    NetworkChangeNotificationInterface(interfaces, count);
    return true;
}

bool
PCUpdateGatewayList(SCNIGateway** gateways, unsigned int count) {
#ifdef COMPARE_NATIVE_GO_OUTPUT
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
#endif
    
    NetworkChangeNotificationGateway(gateways, count);
    return true;
}

@implementation AppDelegate(Netmonitor)
#pragma mark - PCInterfaceStatusNotification
-(void)PCInterfaceStatusChanged:(PCInterfaceStatus *)monitor interfaceStatus:(PCNetworkInterface**)status count:(unsigned int)count {
    PCUpdateInterfaceList(status, count);
}

-(void)PCGatewayStatusChanged:(PCInterfaceStatus *)monitor gatewayStatus:(SCNIGateway**)status count:(unsigned int)count {
    PCUpdateGatewayList(status, count);
}
@end
