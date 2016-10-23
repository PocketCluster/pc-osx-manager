//
//  main.m
//  LinkDetector
//
//  Created by Almighty Kim on 10/17/16.
//  Copyright (c) 2016 PocketCluster.io. All rights reserved.
//

#import "PCInterfaceStatus.h"

bool
pc_interface_list(PCNetworkInterface** interfaces, unsigned int count) {
    NSLog(@"total address count %d", count);
    for (unsigned int i = 0; i < count; i++) {
        
        PCNetworkInterface *iface = *(interfaces + i);
        printf("%s : %s\n",iface->bsdName, iface->isActive ? "active" : "in-active");
        if (iface->isPrimary) {
            printf("\t!!!THIS IS PRIMARY INTERFACE!!!\n");
        }
        printf("\t%s \n",iface->macAddress);
        printf("\t%s \n",iface->mediaType);
        printf("\t%s \n---------\n",iface->displayName);
    }
    return true;
}

bool
gateway_list(SCNIGateway** gateways, unsigned int count) {
    NSLog(@"Total gateway count %d", count);
    for (unsigned int i = 0; i < count; i++) {
        SCNIGateway *gw = *(gateways + i);
        printf("%s : %s\n",gw->addr, gw->is_default ? "DEFAULT" : "AUXILLARY");
        printf("\t%s \n---------\n", gw->ifname);
    }
    return true;
}

//TODO : 2) Interface status 3) Address Status 5) Async notification 6) leak check
int main(int argc, const char * argv[]) {
    @autoreleasepool {
        
        
        
        interface_status_with_callback(&pc_interface_list);
        gateway_status_with_callback(&gateway_list);
        
        PCInterfaceStatus *status = [PCInterfaceStatus new];
        [status startMonitoring];
        unsigned int counter = 0;
        while (counter < 120) {
            sleep(1);
            counter++;
        }
        [status stopMonitoring];
    }
    
    sleep(5);
    return 0;
}
