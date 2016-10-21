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

int main(int argc, const char * argv[]) {
    @autoreleasepool {
        interface_status_with_callback(&pc_interface_list);
    }
    
    sleep(5);
    return 0;
}
